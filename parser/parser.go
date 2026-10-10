// Portions of this file are derived from holistics/dbml packages/dbml-parse/src/core/parser/parser.ts
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

// Package parser builds an ast.Program from DBML tokens.
//
// It is a hand-written recursive-descent parser (with Pratt parsing for
// expressions) ported from @dbml/parse. Syntax errors never abort parsing:
// each parse method either recovers locally (panic-mode synchronisation) or
// returns a partial node to a caller whose context can recover, so the tree
// is always complete enough for editor features and further diagnostics.
package parser

import (
	"strconv"
	"strings"

	"github.com/aidenappl/go-dbml/ast"
	"github.com/aidenappl/go-dbml/diag"
	"github.com/aidenappl/go-dbml/lexer"
)

type ctxKind int

const (
	ctxNone ctxKind = iota - 1
	ctxList
	ctxGroup
	ctxBlock
)

// pErr is a parsing failure carrying the partially parsed node and the
// context that should handle it (ctxNone = top level).
type pErr struct {
	token   *lexer.Token
	partial ast.Node
	handler ctxKind
}

// Result is the output of Parse.
type Result struct {
	Program *ast.Program
	Tokens  []*lexer.Token
	Errors  []diag.Diagnostic
}

// Parser holds parsing state. Use Parse.
type Parser struct {
	// depth is the current recursion depth; maxDepth bounds it (see enter).
	depth         int
	maxDepth      int
	depthReported bool
	maxErrs       int
	// maxElems bounds the elements of one block / list / tuple / row
	maxElems     int
	elemReported bool
	// maxRecordBytes bounds the source bytes of all Records bodies
	maxRecordBytes int
	recordBytes    int
	maxRecordRows  int
	recordRows     int
	recReported    bool
	cancel         <-chan struct{}
	polls          int
	canceled       bool

	tokens   []*lexer.Token
	current  int
	errs     []diag.Diagnostic
	ctx      contextStack
	f        *ast.Factory
	source   string
	filepath string
}

// DefaultMaxDepth bounds nesting of brackets, blocks, prefix operators and
// chained operators. Real schemas nest a handful of levels; the limit keeps
// every recursive walker (parser and later passes) on a small, bounded stack
// for hostile input such as a megabyte of "(".
const DefaultMaxDepth = 256

// Options tunes the parser.
type Options struct {
	// MaxDepth is the maximum nesting depth (DefaultMaxDepth when <= 0).
	// Deeper input is skipped with one NESTING_TOO_DEEP diagnostic.
	MaxDepth int
	// MaxErrors bounds the syntax errors recorded (0 = unlimited); at most
	// MaxErrors+1 are kept so callers can detect truncation.
	MaxErrors int
	// MaxTokens bounds the non-trivia tokens (0 = unlimited); see
	// lexer.Options.
	MaxTokens int
	// MaxElements bounds the elements of one block, list, tuple or comma
	// row (0 = unlimited); the rest is skipped with one TOO_MANY_ELEMENTS.
	MaxElements int
	// MaxRecordBytes bounds the source bytes of all Records bodies (0 =
	// unlimited); a body past it is skipped, unparsed, with one
	// RECORDS_TOO_LARGE (the interpreter enforces the same limit on what
	// it reads).
	MaxRecordBytes int
	// MaxRecordRows bounds the rows of all Records bodies (0 = unlimited);
	// once exceeded, later Records bodies are skipped unparsed with one
	// RECORDS_TOO_LARGE.
	MaxRecordRows int
	// Cancel stops lexing and parsing when closed: the lexer stops with a
	// CANCELED error and the parser treats the rest as end of input.
	Cancel <-chan struct{}
}

// Parse lexes and parses src with default options.
func Parse(src, filepath string) Result { return ParseWithOptions(src, filepath, Options{}) }

// ParseWithOptions lexes and parses src.
func ParseWithOptions(src, filepath string, opts Options) Result {
	toks, lexErrs := lexer.LexWithOptions(src, filepath, lexer.Options{MaxErrors: opts.MaxErrors, MaxTokens: opts.MaxTokens, Cancel: opts.Cancel})
	res := parseTokens(src, toks, filepath, opts)
	res.Errors = append(append([]diag.Diagnostic{}, lexErrs...), res.Errors...)
	return res
}

// ParseTokens parses an already-lexed token stream with default options.
func ParseTokens(src string, toks []*lexer.Token, filepath string) Result {
	return parseTokens(src, toks, filepath, Options{})
}

func parseTokens(src string, toks []*lexer.Token, filepath string, opts Options) Result {
	max := opts.MaxDepth
	if max <= 0 {
		max = DefaultMaxDepth
	}
	p := &Parser{tokens: toks, f: &ast.Factory{Filepath: filepath}, source: src, filepath: filepath, maxDepth: max, maxErrs: opts.MaxErrors,
		maxElems: opts.MaxElements, maxRecordBytes: opts.MaxRecordBytes, maxRecordRows: opts.MaxRecordRows, cancel: opts.Cancel}
	body := p.program()
	if p.canceled {
		p.canceled = false // let advance reach EOF
		p.current = len(p.tokens) - 1
		p.errs = append(p.errs, diag.New(diag.CANCELED, "Compilation was canceled", p.filepath, diag.Range{Start: lexer.ToPos(p.last().Start), End: lexer.ToPos(p.last().End)}))
	}
	eof := p.advance()
	prog := p.f.NewProgram(body, eof, src)
	p.gatherInvalid()
	ast.AssignParents(prog)
	return Result{Program: prog, Tokens: p.tokens, Errors: p.errs}
}

// ---------------------------------------------------------------- helpers

func (p *Parser) isAtEnd() bool {
	if p.cancel != nil && !p.canceled {
		p.polls++
		if p.polls&1023 == 0 {
			select {
			case <-p.cancel:
				p.canceled = true
			default:
			}
		}
	}
	return p.canceled || p.current >= len(p.tokens) || p.tokens[p.current].Kind == lexer.EOF
}

func (p *Parser) last() *lexer.Token { return p.tokens[len(p.tokens)-1] }

func (p *Parser) advance() *lexer.Token {
	if p.isAtEnd() {
		return p.last()
	}
	t := p.tokens[p.current]
	p.current++
	return t
}

func (p *Parser) peek(n int) *lexer.Token {
	if p.current+n >= len(p.tokens) {
		return p.last()
	}
	return p.tokens[p.current+n]
}

func (p *Parser) check(kinds ...lexer.Kind) bool {
	k := p.peek(0).Kind
	for _, want := range kinds {
		if k == want {
			return true
		}
	}
	return false
}

func (p *Parser) match(kinds ...lexer.Kind) bool {
	if p.check(kinds...) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) previous() *lexer.Token {
	if p.current-1 < 0 {
		return nil
	}
	return p.tokens[p.current-1]
}

func (p *Parser) canHandle(e *pErr) bool {
	return e.handler == ctxNone || e.handler == p.ctx.top()
}

func (p *Parser) newErr(tok *lexer.Token, partial ast.Node) *pErr {
	return &pErr{token: tok, partial: partial, handler: p.ctx.findHandler(p.tokens, p.current)}
}

func (p *Parser) consume(msg string, kinds ...lexer.Kind) *pErr {
	if !p.match(kinds...) {
		p.logErrorTok(p.peek(0), diag.UNEXPECTED_TOKEN, msg)
		return p.newErr(p.peek(0), nil)
	}
	return nil
}

func (p *Parser) discardUntil(msg string, kinds ...lexer.Kind) bool {
	if p.isAtEnd() || !p.check(kinds...) {
		ast.MarkInvalid(p.peek(0))
		p.logErrorTok(p.advance(), diag.UNEXPECTED_TOKEN, msg)
		for !p.isAtEnd() && !p.check(kinds...) {
			ast.MarkInvalid(p.advance())
		}
		return !p.isAtEnd()
	}
	return true
}

func (p *Parser) gatherInvalid() {
	i := 0
	var first []*lexer.Token
	for ; i < len(p.tokens); i++ {
		if p.tokens[i].Invalid {
			first = append(first, p.tokens[i])
		} else {
			break
		}
	}
	if i >= len(p.tokens) {
		return
	}
	cur := p.tokens[i]
	cur.PrependLeadingInvalid(first)
	toks := make([]*lexer.Token, 0, len(p.tokens))
	toks = append(toks, cur)
	for i++; i < len(p.tokens); i++ {
		t := p.tokens[i]
		if t.Invalid {
			cur.AppendTrailingInvalid(t)
		} else {
			cur = t
			toks = append(toks, cur)
		}
	}
	p.tokens = toks
}

func (p *Parser) logErrorTok(t *lexer.Token, code diag.Code, msg string) {
	if p.maxErrs > 0 && len(p.errs) > p.maxErrs {
		return
	}
	p.errs = append(p.errs, diag.New(code, msg, p.filepath, diag.Range{Start: lexer.ToPos(t.Start), End: lexer.ToPos(t.End)}))
}

func (p *Parser) logErrorNode(n ast.Node, code diag.Code, msg string) {
	if p.maxErrs > 0 && len(p.errs) > p.maxErrs {
		return
	}
	b := n.Base()
	p.errs = append(p.errs, diag.New(code, msg, p.filepath, diag.Range{Start: lexer.ToPos(b.Start), End: lexer.ToPos(b.End)}))
}

// ---------------------------------------------------------------- depth limit

// enter increments the nesting depth, or returns false when the limit is
// reached. Callers that get false must not recurse further.
func (p *Parser) enter() bool {
	if p.depth >= p.maxDepth {
		return false
	}
	p.depth++
	return true
}

func (p *Parser) leave() { p.depth-- }

// reportTooDeep logs NESTING_TOO_DEEP once per file.
func (p *Parser) reportTooDeep() {
	if p.depthReported {
		return
	}
	p.depthReported = true
	p.logErrorTok(p.peek(0), diag.NESTING_TOO_DEEP, "Nesting is too deep (more than "+strconv.Itoa(p.maxDepth)+" levels); the nested part was skipped")
}

func isOpener(k lexer.Kind) bool {
	return k == lexer.LParen || k == lexer.LBracket || k == lexer.LBrace
}
func isCloser(k lexer.Kind) bool {
	return k == lexer.RParen || k == lexer.RBracket || k == lexer.RBrace
}

// tooMany reports whether a container already holds the maximum number of
// elements; the first time, it records TOO_MANY_ELEMENTS.
func (p *Parser) tooMany(n int) bool {
	if p.maxElems <= 0 || n < p.maxElems {
		return false
	}
	if !p.elemReported {
		p.elemReported = true
		p.logErrorTok(p.peek(0), diag.TOO_MANY_ELEMENTS, "More than "+strconv.Itoa(p.maxElems)+" elements in one block or list; the rest was skipped")
	}
	return true
}

// skipTo consumes tokens up to (not including) the next closer of kind k
// at nesting depth 0, without building nodes.
func (p *Parser) skipTo(k lexer.Kind) {
	depth := 0
	for !p.isAtEnd() {
		t := p.peek(0).Kind
		if depth == 0 && t == k {
			return
		}
		switch {
		case isOpener(t):
			depth++
		case isCloser(t):
			if depth == 0 {
				return // an unbalanced closer of another kind
			}
			depth--
		}
		p.advance()
	}
}

// skipBalanced consumes an opener and everything up to and including its
// matching closer (or EOF), iteratively.
func (p *Parser) skipBalanced() {
	n := 0
	for !p.isAtEnd() {
		k := p.peek(0).Kind
		ast.MarkInvalid(p.advance())
		switch {
		case isOpener(k):
			n++
		case isCloser(k):
			n--
		}
		if n <= 0 {
			return
		}
	}
}

// terminator reports whether t ends the innermost open context, so a
// too-deep expression may stop there without consuming anything.
func (p *Parser) terminator(t *lexer.Token) bool {
	switch p.ctx.top() {
	case ctxBlock:
		return t.Kind == lexer.RBrace
	case ctxGroup:
		return t.Kind == lexer.RParen || t.Kind == lexer.Comma
	case ctxList:
		return t.Kind == lexer.RBracket || t.Kind == lexer.Comma
	}
	return false
}

// skipExpression discards the rest of a too-deep expression: balanced
// bracket runs are skipped whole, and it stops (unconsumed) at a closer, a
// comma or a new line. It always consumes at least one token unless the
// current token terminates the enclosing context, guaranteeing progress.
func (p *Parser) skipExpression() {
	consumed := false
	for !p.isAtEnd() {
		t := p.peek(0)
		if consumed || p.terminator(t) {
			if isCloser(t.Kind) || t.Kind == lexer.Comma || (consumed && lexer.IsAtStartOfLine(p.previous(), t)) {
				return
			}
		}
		if isOpener(t.Kind) {
			p.skipBalanced()
		} else {
			ast.MarkInvalid(p.advance())
		}
		consumed = true
	}
}

func isKeyword(t *lexer.Token, kw string) bool {
	return t != nil && t.Kind == lexer.Identifier && strings.EqualFold(t.Value, kw)
}

// ---------------------------------------------------------------- program

func (p *Parser) program() []ast.Node {
	var stmts []ast.Node
	for !p.isAtEnd() {
		if isKeyword(p.peek(0), "use") || isKeyword(p.peek(0), "reuse") {
			stmts = append(stmts, p.useDeclaration())
			continue
		}
		el, err := p.elementDeclaration()
		if err != nil {
			if !ast.IsNil(err.partial) {
				stmts = append(stmts, err.partial)
			}
			p.synchronizeProgram()
			continue
		}
		stmts = append(stmts, el)
	}
	return stmts
}

func (p *Parser) synchronizeProgram() {
	t := p.peek(0)
	if t.Kind != lexer.EOF {
		ast.MarkInvalid(p.advance())
	} else {
		p.logErrorTok(t, diag.UNEXPECTED_EOF, "Unexpected EOF")
	}
}

func (p *Parser) useDeclaration() *ast.UseDeclaration {
	useKw := p.advance()
	var specifiers ast.Node
	if p.peek(0).Kind == lexer.Wildcard {
		p.advance()
		specifiers = p.f.NewWildcard(p.previous())
	} else {
		list, err := p.useSpecifierList()
		if err != nil {
			if l, ok := err.partial.(*ast.UseSpecifierList); ok && l != nil {
				specifiers = l
			}
		} else {
			specifiers = list
		}
	}
	afterWhat := "specifier list"
	if _, ok := specifiers.(*ast.Wildcard); ok {
		afterWhat = "'*'"
	}
	var fromKw, path *lexer.Token
	if isKeyword(p.peek(0), "from") {
		fromKw = p.advance()
	} else {
		p.logErrorTok(p.peek(0), diag.UNEXPECTED_TOKEN, "Expect 'from' after "+afterWhat)
	}
	if p.match(lexer.StringLiteral) {
		path = p.previous()
	} else {
		p.logErrorTok(p.peek(0), diag.UNEXPECTED_TOKEN, "Expect a string literal path")
	}
	return p.f.NewUseDeclaration(useKw, specifiers, fromKw, path)
}

func (p *Parser) useSpecifierList() (*ast.UseSpecifierList, *pErr) {
	p.ctx.push(ctxBlock)
	defer p.ctx.pop()

	var open, closeTok *lexer.Token
	var specs []*ast.UseSpecifier
	build := func() *ast.UseSpecifierList { return p.f.NewUseSpecifierList(open, specs, closeTok) }

	if err := p.consume("Expect an opening brace '{'", lexer.LBrace); err != nil {
		if !p.canHandle(err) {
			return nil, &pErr{err.token, build(), err.handler}
		}
		p.synchronizeUseSpecifierList()
	} else {
		open = p.previous()
	}

	for !p.isAtEnd() && !p.check(lexer.RBrace) {
		spec, err := p.useSpecifier()
		if err != nil {
			if s, ok := err.partial.(*ast.UseSpecifier); ok && s != nil {
				specs = append(specs, s)
			}
			if !p.canHandle(err) {
				return nil, &pErr{err.token, build(), err.handler}
			}
			p.synchronizeUseSpecifierList()
		} else {
			specs = append(specs, spec)
		}
		if p.check(lexer.Identifier) && isKeyword(p.peek(0), "from") {
			break
		}
		if !p.previous().HasTrailingNewline() && !p.check(lexer.RBrace) {
			p.logErrorTok(p.peek(0), diag.MISSING_SPACES, "Expect a newline between use specifiers")
		}
	}

	if err := p.consume("Expect a closing brace '}'", lexer.RBrace); err != nil {
		if !p.canHandle(err) {
			return nil, &pErr{err.token, build(), err.handler}
		}
		p.synchronizeUseSpecifierList()
	} else {
		closeTok = p.previous()
	}
	return build(), nil
}

func (p *Parser) synchronizeUseSpecifierList() {
	if p.check(lexer.RBrace) {
		return
	}
	ast.MarkInvalid(p.advance())
	for !p.isAtEnd() {
		t := p.peek(0)
		if p.check(lexer.RBrace) || p.check(lexer.StringLiteral) || (p.check(lexer.Identifier) && isKeyword(t, "from")) || lexer.IsAtStartOfLine(p.previous(), t) {
			break
		}
		ast.MarkInvalid(t)
		p.advance()
	}
}

func (p *Parser) useSpecifier() (*ast.UseSpecifier, *pErr) {
	var kind, asKw *lexer.Token
	var name, alias ast.Node
	build := func() *ast.UseSpecifier { return p.f.NewUseSpecifier(kind, name, asKw, alias) }

	if err := p.consume("Expect an import kind specifier (e.g. table, enum)", lexer.Identifier); err != nil {
		return nil, &pErr{err.token, build(), err.handler}
	}
	kind = p.previous()

	if p.check(lexer.Identifier, lexer.QuotedString) {
		n, err := p.normalExpression()
		if err != nil {
			if !ast.IsNil(err.partial) {
				name = err.partial
			}
			return nil, &pErr{err.token, build(), err.handler}
		}
		name = n
	} else {
		p.logErrorTok(p.peek(0), diag.UNEXPECTED_TOKEN, "Expect an element name")
	}

	if isKeyword(p.peek(0), "as") {
		asKw = p.advance()
		if !p.check(lexer.Identifier, lexer.QuotedString) {
			p.logErrorTok(p.peek(0), diag.UNEXPECTED_TOKEN, "Expect an alias name after 'as'")
			err := p.newErr(p.peek(0), build())
			if !p.canHandle(err) {
				return nil, &pErr{err.token, build(), err.handler}
			}
		} else {
			a, err := p.normalExpression()
			if err != nil {
				if !p.canHandle(err) {
					return nil, &pErr{err.token, build(), err.handler}
				}
			} else {
				alias = a
			}
		}
	}
	return build(), nil
}

// ---------------------------------------------------------------- elements

func (p *Parser) elementDeclaration() (*ast.ElementDeclaration, *pErr) {
	var e ast.ElementDeclaration
	build := func() *ast.ElementDeclaration { return p.f.NewElementDeclaration(e) }

	if err := p.consume("Expect an identifier", lexer.Identifier); err != nil {
		return nil, &pErr{err.token, build(), err.handler}
	}
	e.Type = p.previous()

	if isKeyword(e.Type, "metadata") {
		if err := p.consume("Expect a metadata target kind", lexer.Identifier); err != nil {
			return nil, &pErr{err.token, build(), err.handler}
		}
		e.TargetKind = p.previous()
	}

	if !p.check(lexer.Colon, lexer.LBrace, lexer.LBracket) {
		n, err := p.normalExpression()
		if err != nil {
			e.Name = err.partial
			if !p.canHandle(err) {
				return nil, &pErr{err.token, build(), err.handler}
			}
			p.synchronizeElementDeclarationName()
		} else {
			e.Name = n
		}
	}

	if isKeyword(p.peek(0), "as") {
		e.As = p.advance()
		if !p.check(lexer.Colon, lexer.LBrace, lexer.LBracket) {
			n, err := p.normalExpression()
			if err != nil {
				e.Alias = err.partial
				if !p.canHandle(err) {
					return nil, &pErr{err.token, build(), err.handler}
				}
				p.synchronizeElementDeclarationAlias()
			} else {
				e.Alias = n
			}
		} else {
			p.logErrorTok(p.peek(0), diag.UNEXPECTED_TOKEN, "Expect an alias")
		}
	}

	if p.check(lexer.LBracket) {
		l, err := p.listExpression()
		if err != nil {
			if pl, ok := err.partial.(*ast.ListExpression); ok {
				e.AttributeList = pl
			}
			return nil, &pErr{err.token, build(), err.handler}
		}
		e.AttributeList = l
	}

	if !p.discardUntil("Expect an opening brace '{' or a colon ':'", lexer.LBrace, lexer.Colon) {
		return build(), nil
	}

	if p.match(lexer.Colon) {
		e.BodyColon = p.previous()
		expr, err := p.expression()
		if err != nil {
			e.Body = err.partial
			return nil, &pErr{err.token, build(), err.handler}
		}
		if ed, ok := expr.(*ast.ElementDeclaration); ok {
			ast.MarkInvalid(ed)
			p.logErrorNode(ed, diag.UNEXPECTED_ELEMENT_DECLARATION, "An element's simple body must not be an element declaration")
		} else {
			e.Body = expr
		}
	} else {
		records := isKeyword(e.Type, "records")
		if records && p.skipLargeRecords() {
			e.Body = p.f.NewBlock(nil, nil, nil)
			return build(), nil
		}
		b, err := p.blockExpression()
		if records && b != nil {
			p.recordRows += len(b.Body)
		}
		if err != nil {
			e.Body = err.partial
			return nil, &pErr{err.token, build(), err.handler}
		}
		e.Body = b
	}
	return build(), nil
}

// skipLargeRecords measures the Records body at the current '{' and, when
// the bodies seen so far exceed MaxRecordBytes, skips it unparsed with one
// RECORDS_TOO_LARGE error (bounding memory before any node is built).
func (p *Parser) skipLargeRecords() bool {
	if !p.check(lexer.LBrace) {
		return false
	}
	if p.maxRecordRows > 0 && p.recordRows > p.maxRecordRows {
		return p.skipRecordsBody("Records exceed the " + strconv.Itoa(p.maxRecordRows) + "-row limit; the block was not parsed")
	}
	if p.maxRecordBytes <= 0 {
		return false
	}
	open := p.peek(0)
	depth, i := 0, p.current
	for ; i < len(p.tokens); i++ {
		switch k := p.tokens[i].Kind; {
		case k == lexer.EOF:
			i = len(p.tokens) - 1
			depth = -1
		case isOpener(k):
			depth++
		case isCloser(k):
			depth--
		}
		if depth <= 0 {
			break
		}
	}
	if i >= len(p.tokens) {
		i = len(p.tokens) - 1
	}
	p.recordBytes += int(p.tokens[i].End.Byte - open.Start.Byte)
	if p.recordBytes <= p.maxRecordBytes {
		return false
	}
	return p.skipRecordsBody("Records exceed the " + strconv.Itoa(p.maxRecordBytes) + "-byte limit; the block was not parsed")
}

// skipRecordsBody skips the Records body at the current '{' (through its
// '}') and reports RECORDS_TOO_LARGE once.
func (p *Parser) skipRecordsBody(msg string) bool {
	if !p.recReported {
		p.recReported = true
		p.logErrorTok(p.peek(0), diag.RECORDS_TOO_LARGE, msg)
	}
	// skipped without building nodes or marking tokens (no allocation)
	p.advance() // {
	p.skipTo(lexer.RBrace)
	if p.check(lexer.RBrace) {
		p.advance()
	}
	return true
}

func (p *Parser) synchronizeElementDeclarationName() {
	for !p.isAtEnd() {
		t := p.peek(0)
		if isKeyword(t, "as") || p.check(lexer.Colon, lexer.LBrace, lexer.LBracket) {
			break
		}
		ast.MarkInvalid(t)
		p.advance()
	}
}

func (p *Parser) synchronizeElementDeclarationAlias() {
	for !p.isAtEnd() {
		t := p.peek(0)
		if p.check(lexer.Colon, lexer.LBrace, lexer.LBracket) {
			break
		}
		ast.MarkInvalid(t)
		p.advance()
	}
}

// fieldDeclaration parses `ident: <expr>` inside a block (e.g. `Note: '...'`).
func (p *Parser) fieldDeclaration() (*ast.ElementDeclaration, *pErr) {
	var e ast.ElementDeclaration
	build := func() *ast.ElementDeclaration { return p.f.NewElementDeclaration(e) }

	if err := p.consume("Expect an identifier", lexer.Identifier); err != nil {
		return nil, &pErr{err.token, build(), err.handler}
	}
	e.Type = p.previous()
	if err := p.consume("Expect a colon ':'", lexer.Colon); err != nil {
		return nil, &pErr{err.token, build(), err.handler}
	}
	e.BodyColon = p.previous()

	expr, err := p.expression()
	if err != nil {
		e.Body = err.partial
		return nil, &pErr{err.token, build(), err.handler}
	}
	if ed, ok := expr.(*ast.ElementDeclaration); ok {
		p.logErrorNode(ed, diag.INVALID_ELEMENT_IN_SIMPLE_BODY, "Simple body cannot be an element declaration")
	} else {
		e.Body = expr
	}
	return build(), nil
}

// ---------------------------------------------------------------- expressions

// expression parses a function application (the most generic form), which
// may turn out to be a nested element declaration.
func (p *Parser) expression() (ast.Node, *pErr) {
	var callee ast.Node
	var args []ast.Node
	build := func() ast.Node {
		if el := convertFuncAppToElem(callee, args, p.f); el != nil {
			return el
		}
		return p.f.NewFunctionApplication(callee, args)
	}

	c, err := p.commaExpression()
	if err != nil {
		callee = err.partial
		return nil, &pErr{err.token, build(), err.handler}
	}
	callee = c

	if p.shouldStopFunctionApplication() {
		return build(), nil
	}

	prevNode := callee
	for !p.shouldStopFunctionApplication() {
		if !p.previous().HasTrailingSpaces() {
			p.logErrorNode(prevNode, diag.MISSING_SPACES, "Expect a following space")
		}
		n, err := p.commaExpression()
		if err != nil {
			prevNode = err.partial
			args = append(args, prevNode)
			return nil, &pErr{err.token, build(), err.handler}
		}
		prevNode = n
		args = append(args, n)
	}
	return build(), nil
}

func (p *Parser) shouldStopFunctionApplication() bool {
	if p.isAtEnd() || p.previous().HasTrailingNewline() {
		return true
	}
	switch p.peek(0).Kind {
	case lexer.RBrace, lexer.RBracket, lexer.RParen, lexer.Comma, lexer.Colon:
		return true
	}
	return false
}

func (p *Parser) commaExpression() (ast.Node, *pErr) {
	var first ast.Node
	if p.check(lexer.Comma) {
		first = p.f.NewEmpty(p.previous())
	} else {
		n, err := p.normalExpression()
		if err != nil {
			return nil, err
		}
		first = n
	}
	if !p.check(lexer.Comma) {
		return first, nil
	}
	elems := []ast.Node{first}
	var commas []*lexer.Token
	for {
		if p.tooMany(len(elems)) {
			for !p.shouldStopCommaExpression() {
				p.advance()
			}
			break
		}
		commas = append(commas, p.advance())
		if p.shouldStopCommaExpression() {
			elems = append(elems, p.f.NewEmpty(p.previous()))
			break
		}
		if p.check(lexer.Comma) {
			elems = append(elems, p.f.NewEmpty(p.previous()))
			continue
		}
		n, err := p.normalExpression()
		if err != nil {
			if !ast.IsNil(err.partial) {
				elems = append(elems, err.partial)
			}
			return nil, &pErr{err.token, p.f.NewComma(elems, commas), err.handler}
		}
		elems = append(elems, n)
		if p.shouldStopCommaExpression() || !p.check(lexer.Comma) {
			break
		}
	}
	return p.f.NewComma(elems, commas), nil
}

func (p *Parser) shouldStopCommaExpression() bool {
	if p.isAtEnd() || p.previous().HasTrailingNewline() {
		return true
	}
	switch p.peek(0).Kind {
	case lexer.RBrace, lexer.LBrace, lexer.RBracket, lexer.LBracket, lexer.RParen, lexer.LParen, lexer.Colon:
		return true
	}
	return false
}

func (p *Parser) normalExpression() (ast.Node, *pErr) { return p.expressionBP(0) }

// expressionBP implements Pratt parsing with minimum binding power mbp.
func (p *Parser) expressionBP(mbp int) (ast.Node, *pErr) {
	if !p.enter() {
		p.reportTooDeep()
		prev := p.previous()
		p.skipExpression()
		return p.f.NewEmpty(prev), nil
	}
	defer p.leave()
	left, err := p.leftExpressionBP()
	if err != nil {
		return nil, err
	}
	// links counts operators folded into `left` by this loop; each makes the
	// tree one level deeper, so it counts against the depth limit too.
	links := 0
	for !p.isAtEnd() {
		tok := p.peek(0)
		if p.depth+links >= p.maxDepth {
			if tok.Kind == lexer.Op || tok.Kind == lexer.LParen || tok.Kind == lexer.LBracket {
				p.reportTooDeep()
				p.skipExpression()
			}
			return left, nil
		}
		links++
		switch {
		case tok.Kind == lexer.LParen:
			if precPostfix < mbp {
				return left, nil
			}
			if lexer.IsAtStartOfLine(p.previous(), tok) && !p.ctx.withinGroup() && !p.ctx.withinList() {
				return left, nil
			}
			tup, err := p.tupleExpression()
			if err != nil {
				pt, _ := err.partial.(*ast.TupleExpression)
				left = p.f.NewCall(left, pt)
				return nil, &pErr{err.token, left, err.handler}
			}
			left = p.f.NewCall(left, tup)
		case tok.Kind == lexer.LBracket:
			if p.previous().HasTrailingSpaces() {
				return left, nil
			}
			l, err := p.listExpression()
			if err != nil {
				pl, _ := err.partial.(*ast.ListExpression)
				left = p.f.NewArray(left, pl)
				return nil, &pErr{err.token, left, err.handler}
			}
			left = p.f.NewArray(left, l)
		case tok.Kind != lexer.Op:
			return left, nil
		default:
			op := tok
			// The only postfix operator is '(' which is not an Op token.
			il, ir, ok := infixBP(op.Value)
			if !ok || il <= mbp {
				return left, nil
			}
			p.advance()
			var right ast.Node
			var rerr *pErr
			if op.Value == "." {
				right, rerr = p.extractOperand()
			} else {
				right, rerr = p.expressionBP(ir)
			}
			if rerr != nil {
				left = p.f.NewInfix(left, op, rerr.partial)
				return nil, &pErr{rerr.token, left, rerr.handler}
			}
			left = p.f.NewInfix(left, op, right)
		}
	}
	return left, nil
}

func (p *Parser) leftExpressionBP() (ast.Node, *pErr) {
	if p.peek(0).Kind == lexer.Op {
		op := p.peek(0)
		rbp, ok := prefixBP(op.Value)
		if !ok {
			p.logErrorTok(op, diag.UNKNOWN_PREFIX_OP, "Unexpected '"+op.Value+"' in an expression")
			return nil, p.newErr(op, p.f.NewEmpty(op))
		}
		p.advance()
		expr, err := p.expressionBP(rbp)
		if err != nil {
			return nil, &pErr{err.token, p.f.NewPrefix(op, err.partial), err.handler}
		}
		return p.f.NewPrefix(op, expr), nil
	}
	operand, err := p.extractOperand()
	if err != nil {
		return nil, err
	}
	if _, ok := operand.(*ast.Empty); ok {
		return nil, p.newErr(p.peek(0), p.f.NewEmpty(p.peek(0)))
	}
	return operand, nil
}

// extractOperand parses a primary expression, list, block, tuple, function
// expression or wildcard. Errors from nested list/block/tuple parsing that
// cannot be handled locally propagate to the caller; an unexpected token is
// logged and yields an Empty node.
func (p *Parser) extractOperand() (ast.Node, *pErr) {
	if p.check(lexer.Wildcard) {
		p.advance()
		return p.f.NewWildcard(p.previous()), nil
	}
	if p.check(lexer.NumericLiteral, lexer.StringLiteral, lexer.ColorLiteral, lexer.QuotedString, lexer.Identifier) {
		return p.primaryExpression()
	}
	if p.check(lexer.FunctionExpression) {
		p.advance()
		return p.f.NewFunctionExpression(p.previous()), nil
	}
	if p.check(lexer.LBracket) {
		l, err := p.listExpression()
		if err != nil {
			return nil, err
		}
		return l, nil
	}
	if p.check(lexer.LBrace) {
		b, err := p.blockExpression()
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	if p.check(lexer.LParen) {
		t, err := p.tupleExpression()
		if err != nil {
			return nil, err
		}
		return t, nil
	}
	if p.peek(0).Kind == lexer.EOF {
		p.logErrorTok(p.peek(0), diag.UNEXPECTED_EOF, "Unexpected EOF")
	} else {
		p.logErrorTok(p.peek(0), diag.INVALID_OPERAND, "Invalid start of operand \""+p.peek(0).Value+"\"")
	}
	return p.f.NewEmpty(p.previous()), nil
}

func (p *Parser) primaryExpression() (ast.Node, *pErr) {
	if p.match(lexer.ColorLiteral, lexer.StringLiteral, lexer.NumericLiteral) {
		return p.f.NewPrimary(p.f.NewLiteral(p.previous())), nil
	}
	if p.match(lexer.QuotedString, lexer.Identifier) {
		return p.f.NewPrimary(p.f.NewVariable(p.previous())), nil
	}
	p.logErrorTok(p.peek(0), diag.UNEXPECTED_TOKEN, "Expect a variable or literal")
	return nil, p.newErr(p.peek(0), nil)
}

// ---------------------------------------------------------------- block

func (p *Parser) blockExpression() (*ast.BlockExpression, *pErr) {
	if !p.enter() {
		p.reportTooDeep()
		p.skipBalanced()
		return p.f.NewBlock(nil, nil, nil), nil
	}
	defer p.leave()
	p.ctx.push(ctxBlock)
	defer p.ctx.pop()

	var open, closeTok *lexer.Token
	var body []ast.Node
	build := func() *ast.BlockExpression { return p.f.NewBlock(open, body, closeTok) }

	if err := p.consume("Expect an opening brace '{'", lexer.LBrace); err != nil {
		if !p.canHandle(err) {
			return nil, &pErr{err.token, build(), err.handler}
		}
		p.synchronizeBlock()
	} else {
		open = p.previous()
	}

	for !p.isAtEnd() && !p.check(lexer.RBrace) {
		if p.tooMany(len(body)) {
			p.skipTo(lexer.RBrace)
			break
		}
		var n ast.Node
		var err *pErr
		if p.canBeField() {
			var fd *ast.ElementDeclaration
			fd, err = p.fieldDeclaration()
			if err == nil {
				n = fd
			}
		} else {
			n, err = p.expression()
		}
		if err != nil {
			if !ast.IsNil(err.partial) {
				body = append(body, err.partial)
			}
			if !p.canHandle(err) {
				return nil, &pErr{err.token, build(), err.handler}
			}
			p.synchronizeBlock()
			continue
		}
		body = append(body, n)
	}

	if err := p.consume("Expect a closing brace '}'", lexer.RBrace); err != nil {
		if !p.canHandle(err) {
			return nil, &pErr{err.token, build(), err.handler}
		}
		p.synchronizeBlock()
	} else {
		closeTok = p.previous()
	}
	return build(), nil
}

func (p *Parser) canBeField() bool {
	return p.peek(0).Kind == lexer.Identifier && p.peek(1).Kind == lexer.Colon
}

func (p *Parser) synchronizeBlock() {
	if p.check(lexer.RBrace) {
		return
	}
	ast.MarkInvalid(p.advance())
	for !p.isAtEnd() {
		t := p.peek(0)
		if p.check(lexer.RBrace) || lexer.IsAtStartOfLine(p.previous(), t) {
			break
		}
		ast.MarkInvalid(t)
		p.advance()
	}
}

// ---------------------------------------------------------------- tuple

func (p *Parser) tupleExpression() (*ast.TupleExpression, *pErr) {
	if !p.enter() {
		p.reportTooDeep()
		p.skipBalanced()
		return p.f.NewTuple(nil, nil, nil, nil), nil
	}
	defer p.leave()
	p.ctx.push(ctxGroup)
	defer p.ctx.pop()

	var open, closeTok *lexer.Token
	var elems []ast.Node
	var commas []*lexer.Token
	buildGroup := func() *ast.TupleExpression {
		var e ast.Node
		if len(elems) > 0 {
			e = elems[0]
		}
		return p.f.NewGroup(open, e, closeTok)
	}
	buildTuple := func() *ast.TupleExpression { return p.f.NewTuple(open, elems, commas, closeTok) }

	if err := p.consume("Expect an opening parenthesis '('", lexer.LParen); err != nil {
		if !p.canHandle(err) {
			return nil, &pErr{err.token, buildTuple(), err.handler}
		}
		p.synchronizeTuple()
	} else {
		open = p.previous()
	}

	if !p.isAtEnd() && !p.check(lexer.RParen) {
		n, err := p.normalExpression()
		if err != nil {
			elems = append(elems, err.partial)
			if !p.canHandle(err) {
				return nil, &pErr{err.token, buildGroup(), err.handler}
			}
			p.synchronizeTuple()
		} else {
			elems = append(elems, n)
		}
	}

	for !p.isAtEnd() && !p.check(lexer.RParen) {
		if p.tooMany(len(elems)) {
			p.skipTo(lexer.RParen)
			break
		}
		err := p.consume("Expect a comma ','", lexer.Comma)
		if err == nil {
			commas = append(commas, p.previous())
			var n ast.Node
			n, err = p.normalExpression()
			if err == nil {
				elems = append(elems, n)
			}
		}
		if err != nil {
			if !ast.IsNil(err.partial) {
				elems = append(elems, err.partial)
			}
			if !p.canHandle(err) {
				return nil, &pErr{err.token, buildTuple(), err.handler}
			}
			p.synchronizeTuple()
		}
	}

	if err := p.consume("Expect a closing parenthesis ')'", lexer.RParen); err != nil {
		if !p.canHandle(err) {
			return nil, &pErr{err.token, buildTuple(), err.handler}
		}
		p.synchronizeTuple()
	} else {
		closeTok = p.previous()
	}
	if len(elems) == 1 {
		return buildGroup(), nil
	}
	return buildTuple(), nil
}

func (p *Parser) synchronizeTuple() {
	for !p.isAtEnd() {
		t := p.peek(0)
		if p.check(lexer.RParen, lexer.Comma) {
			break
		}
		ast.MarkInvalid(t)
		p.advance()
	}
}

// ---------------------------------------------------------------- list

func (p *Parser) listExpression() (*ast.ListExpression, *pErr) {
	if !p.enter() {
		p.reportTooDeep()
		p.skipBalanced()
		return p.f.NewList(nil, nil, nil, nil), nil
	}
	defer p.leave()
	p.ctx.push(ctxList)
	defer p.ctx.pop()

	var open, closeTok *lexer.Token
	var elems []*ast.Attribute
	var commas []*lexer.Token
	build := func() *ast.ListExpression { return p.f.NewList(open, elems, commas, closeTok) }

	if err := p.consume("Expect an opening bracket '['", lexer.LBracket); err != nil {
		if !p.canHandle(err) {
			return nil, &pErr{err.token, build(), err.handler}
		}
		p.synchronizeList()
	} else {
		open = p.previous()
	}

	if !p.isAtEnd() && !p.check(lexer.RBracket) {
		a, err := p.attribute()
		if err != nil {
			if pa, ok := err.partial.(*ast.Attribute); ok && pa != nil {
				elems = append(elems, pa)
			}
			if !p.canHandle(err) {
				return nil, &pErr{err.token, build(), err.handler}
			}
			p.synchronizeList()
		} else {
			elems = append(elems, a)
		}
	}

	for !p.isAtEnd() && !p.check(lexer.RBracket) {
		if p.tooMany(len(elems)) {
			p.skipTo(lexer.RBracket)
			break
		}
		err := p.consume("Expect a comma ','", lexer.Comma)
		if err == nil {
			commas = append(commas, p.previous())
			var a *ast.Attribute
			a, err = p.attribute()
			if err == nil {
				elems = append(elems, a)
			}
		}
		if err != nil {
			if pa, ok := err.partial.(*ast.Attribute); ok && pa != nil {
				elems = append(elems, pa)
			}
			if !p.canHandle(err) {
				return nil, &pErr{err.token, build(), err.handler}
			}
			p.synchronizeList()
		}
	}

	if err := p.consume("Expect a closing bracket ']'", lexer.RBracket); err != nil {
		if !p.canHandle(err) {
			return nil, &pErr{err.token, build(), err.handler}
		}
		p.synchronizeList()
	} else {
		closeTok = p.previous()
	}
	return build(), nil
}

func (p *Parser) synchronizeList() {
	for !p.isAtEnd() {
		t := p.peek(0)
		if p.check(lexer.Comma, lexer.RBracket) {
			break
		}
		ast.MarkInvalid(t)
		p.advance()
	}
}

func (p *Parser) attribute() (*ast.Attribute, *pErr) {
	var name, value ast.Node
	var colon *lexer.Token

	if p.check(lexer.Colon, lexer.RBracket, lexer.Comma) {
		p.logErrorTok(p.peek(0), diag.EMPTY_ATTRIBUTE_NAME, "Expect a non-empty attribute name")
		name = p.f.NewIdentifierStream(nil)
	} else {
		n, err := p.attributeName()
		if err != nil {
			name = err.partial
			if !p.canHandle(err) {
				return nil, &pErr{err.token, p.f.NewAttribute(name, colon, value), err.handler}
			}
			p.synchronizeAttributeName()
		} else {
			name = n
		}
	}

	if p.match(lexer.Colon) {
		colon = p.previous()
		v, err := p.attributeValue()
		if err != nil {
			// The reference rethrows the value's own partial (not an
			// attribute); the enclosing list then holds no attribute for it.
			return nil, &pErr{err.token, nil, err.handler}
		}
		value = v
	}
	return p.f.NewAttribute(name, colon, value), nil
}

func (p *Parser) synchronizeAttributeName() {
	for !p.isAtEnd() {
		t := p.peek(0)
		if p.check(lexer.Comma, lexer.RBracket, lexer.Colon) {
			break
		}
		ast.MarkInvalid(t)
		p.advance()
	}
}

func (p *Parser) attributeValue() (ast.Node, *pErr) {
	var v ast.Node
	var err *pErr
	if p.peek(0).Kind == lexer.Identifier && p.peek(1).Kind == lexer.Identifier {
		v, err = p.attributeName()
	} else {
		v, err = p.normalExpression()
	}
	if err != nil {
		if !p.canHandle(err) {
			return err.partial, err
		}
		v = err.partial
		p.synchronizeAttributeValue()
	}
	return v, nil
}

func (p *Parser) synchronizeAttributeValue() {
	for !p.isAtEnd() {
		t := p.peek(0)
		if p.check(lexer.Comma, lexer.RBracket) {
			break
		}
		ast.MarkInvalid(t)
		p.advance()
	}
}

func (p *Parser) attributeName() (ast.Node, *pErr) {
	if p.peek(0).Kind != lexer.Identifier {
		return p.primaryExpression()
	}
	var ids []*lexer.Token
	for !p.isAtEnd() && !p.check(lexer.Colon, lexer.Comma, lexer.RBracket) {
		if err := p.consume("Expect an identifier", lexer.Identifier); err != nil {
			return nil, &pErr{err.token, p.f.NewIdentifierStream(ids), err.handler}
		}
		ids = append(ids, p.previous())
	}
	return p.f.NewIdentifierStream(ids), nil
}

// ---------------------------------------------------------------- binding powers

const (
	precAssign       = 2
	precEquality     = 4
	precComparison   = 7
	precAdditive     = 9
	precMultiplicity = 11
	precPostfix      = 14
	precPrefix       = 15
	precMember       = 16
)

var infixLeft = map[string]int{
	".": precMember, "/": precMultiplicity, "%": precMultiplicity, "+": precAdditive, "=": precAssign,
	"==": precEquality, "!=": precEquality,
	"-": precAdditive, "-?": precAdditive, "?-": precAdditive, "?-?": precAdditive,
	"<": precComparison, "<=": precComparison, ">": precComparison, ">=": precComparison, "<>": precComparison,
	"?>": precComparison, ">?": precComparison, "?<": precComparison, "<?": precComparison, "?>?": precComparison,
	"?<?": precComparison, "?<>": precComparison, "<>?": precComparison, "?<>?": precComparison,
	"->": precComparison, "<-": precComparison,
}

func infixBP(op string) (left, right int, ok bool) {
	l, ok := infixLeft[op]
	return l, l + 1, ok
}

var prefixOps = map[string]bool{
	"+": true, "-": true, "<": true, ">": true, "<>": true, "-?": true, "?-": true, "?-?": true, "?>": true, ">?": true,
	"?<": true, "<?": true, "?>?": true, "?<?": true, "?<>": true, "<>?": true, "?<>?": true, "->": true, "<-": true,
	"!": true, "~": true,
}

func prefixBP(op string) (int, bool) {
	if prefixOps[op] {
		return precPrefix, true
	}
	return 0, false
}
