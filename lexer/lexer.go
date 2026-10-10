// Portions of this file are derived from holistics/dbml packages/dbml-parse/src/core/lexer/lexer.ts
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

package lexer

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/aidenappl/go-dbml/diag"
)

// Lexer scans DBML source. Use Lex for the common case.
type Lexer struct {
	text     string
	runes    []rune
	idx      int // index into runes
	start    Position
	current  Position
	filepath string
	tokens   []*Token
	errs     []diag.Diagnostic
	// maxErrs bounds recorded errors (0 = unlimited); one extra is kept so
	// callers can tell the list was truncated.
	maxErrs int
	// maxTokens bounds the non-trivia tokens (0 = unlimited)
	maxTokens int
	stopped   bool
	cancel    <-chan struct{}
	// raw keeps trivia as tokens (re-lexing a trivia range)
	raw bool

	// trivia attachment state (lexer.ts gatherTrivia, done while scanning)
	src           *string
	startOfLine   bool
	lastNonTrivia *Token
	pend          bool
	pendStart     Position
	pendEnd       Position
	pendFlags     uint8
}

// Lex tokenizes src. The returned slice always ends with an EOF token and
// contains no trivia or invalid tokens at the top level (they are attached to
// neighbouring tokens).
func Lex(src, filepath string) ([]*Token, []diag.Diagnostic) {
	return LexWithLimit(src, filepath, 0)
}

// LexWithLimit is Lex with at most maxErrors+1 recorded errors (0 = no limit).
func LexWithLimit(src, filepath string, maxErrors int) ([]*Token, []diag.Diagnostic) {
	return LexWithOptions(src, filepath, Options{MaxErrors: maxErrors})
}

// Options bounds the lexer.
type Options struct {
	// MaxErrors bounds recorded errors (0 = unlimited; one extra is kept).
	MaxErrors int
	// MaxTokens bounds the non-trivia tokens (0 = unlimited). Past it the
	// lexer stops with one TOO_MANY_TOKENS error and the rest of the input
	// is not tokenised.
	MaxTokens int
	// Cancel stops lexing when closed (one CANCELED error).
	Cancel <-chan struct{}
}

// LexWithOptions is Lex with limits.
func LexWithOptions(src, filepath string, o Options) ([]*Token, []diag.Diagnostic) {
	l := &Lexer{text: src, runes: []rune(src), filepath: filepath, maxErrs: o.MaxErrors, maxTokens: o.MaxTokens, src: &src, startOfLine: true, cancel: o.Cancel}
	l.tokens = make([]*Token, 0, len(src)/6+16)
	return l.lex()
}

// lexTrivia re-lexes the trivia in src[start.Byte:end] as tokens.
func lexTrivia(src string, start Position, end int32) []*Token {
	l := &Lexer{text: src, runes: []rune(src[start.Byte:end]), start: start, current: start, raw: true}
	l.scanTokens()
	return l.tokens
}

// ToPos converts a 0-based lexer position to a 1-based diagnostic position.
func ToPos(p Position) diag.Pos {
	return diag.Pos{Offset: int(p.Offset), Line: int(p.Line) + 1, Column: int(p.Column) + 1}
}

// isAtEnd reports whether every rune has been consumed.
func (l *Lexer) isAtEnd() bool { return l.idx >= len(l.runes) }

// peek returns the rune lookahead positions ahead without consuming it.
func (l *Lexer) peek(lookahead int) (rune, bool) {
	if l.idx+lookahead >= len(l.runes) {
		return 0, false
	}
	return l.runes[l.idx+lookahead], true
}

// peekIs reports whether the next rune is c.
func (l *Lexer) peekIs(c rune) bool {
	r, ok := l.peek(0)
	return ok && r == c
}

// utf16Len is the number of UTF-16 code units r occupies (2 outside the BMP).
func utf16Len(r rune) int {
	if r >= 0x10000 && r <= unicode.MaxRune {
		return 2
	}
	return 1
}

// utf8Len is the number of bytes r occupies in UTF-8.
func utf8Len(r rune) int {
	return len(string(r))
}

// advance consumes one rune, updating the line/column (UTF-16), offset
// (UTF-16) and byte positions. Invalid UTF-8 bytes decode to U+FFFD but
// still advance the byte offset by exactly one.
func (l *Lexer) advance() rune {
	c := l.runes[l.idx]
	l.idx++
	if c == '\n' {
		l.current.Line++
		l.current.Column = 0
	} else {
		l.current.Column += int32(utf16Len(c))
	}
	l.current.Offset += int32(utf16Len(c))
	if c == utf8RuneError && !strings.HasPrefix(l.text[l.current.Byte:], string(utf8RuneError)) {
		// invalid UTF-8 byte decoded as U+FFFD; consumes one byte
		l.current.Byte++
	} else {
		l.current.Byte += int32(utf8Len(c))
	}
	return c
}

const utf8RuneError = '�'

// check reports whether the upcoming runes spell seq (without consuming).
func (l *Lexer) check(seq string) bool {
	i := 0
	for _, c := range seq {
		r, ok := l.peek(i)
		if !ok || r != c {
			return false
		}
		i++
	}
	return true
}

// match consumes seq if it is next, reporting whether it did.
func (l *Lexer) match(seq string) bool {
	if !l.check(seq) {
		return false
	}
	for range seq {
		l.advance()
	}
	return true
}

// createToken builds a token spanning start..current with the raw source text.
func (l *Lexer) createToken(kind Kind, invalid bool) *Token {
	return &Token{Kind: kind, Start: l.start, End: l.current, Value: l.text[l.start.Byte:l.current.Byte], Invalid: invalid}
}

// addToken appends a token spanning start..current.
func (l *Lexer) addToken(kind Kind, invalid bool) {
	if !l.raw && isTriviaKind(kind) {
		l.trivia(kind, l.start, l.current)
		return
	}
	l.push(l.createToken(kind, invalid))
}

func isTriviaKind(k Kind) bool {
	switch k {
	case Newline, Space, Tab, SingleLineComment, MultilineComment:
		return true
	}
	return false
}

// trivia records a whitespace/comment span: trivia after a token up to and
// including the newline is its trailing trivia; trivia before the first
// token of a line is that token's leading trivia (lexer.ts gatherTrivia).
func (l *Lexer) trivia(kind Kind, start, end Position) {
	if !l.pend {
		l.pend, l.pendStart, l.pendFlags = true, start, 0
	}
	l.pendEnd = end
	switch kind {
	case Newline:
		l.pendFlags |= flagTrailNewline
	case Space, Tab:
		l.pendFlags |= flagTrailSpace
	}
	if kind == Newline && l.lastNonTrivia != nil {
		l.attachTrailing(l.lastNonTrivia)
		l.startOfLine = true
		l.lastNonTrivia = nil
	}
}

func (l *Lexer) attachTrailing(t *Token) {
	t.trailLen = l.pendEnd.Byte - t.End.Byte
	t.trailEndOff = l.pendEnd.Offset
	t.flags |= l.pendFlags & (flagTrailNewline | flagTrailSpace)
	t.src = l.src
	l.pend = false
}

// push appends a non-trivia token (trivia tokens only in raw mode),
// attaching pending trivia.
func (l *Lexer) push(t *Token) {
	if l.raw {
		l.tokens = append(l.tokens, t)
		return
	}
	if t.IsTrivia() { // an unterminated comment
		l.trivia(t.Kind, t.Start, t.End)
		return
	}
	if l.pend {
		if l.startOfLine {
			t.leadStart, t.leadLen, t.src = l.pendStart, l.pendEnd.Byte-l.pendStart.Byte, l.src
			if l.pendFlags&flagTrailNewline != 0 {
				t.flags |= flagLeadNewline
			}
			l.pend = false
		} else {
			l.attachTrailing(l.lastNonTrivia)
		}
	}
	l.lastNonTrivia = t
	l.startOfLine = false
	l.tokens = append(l.tokens, t)
	if l.maxTokens > 0 && len(l.tokens) > l.maxTokens && !l.stopped && t.Kind != EOF {
		l.stopped = true
		l.errs = append(l.errs, diag.New(diag.TOO_MANY_TOKENS,
			"File has more than "+strconv.Itoa(l.maxTokens)+" tokens; the rest was not parsed", l.filepath,
			diag.Range{Start: ToPos(t.Start), End: ToPos(t.End)}))
	}
}

// errorAt records a lexing error over t's span (bounded by maxErrs).
func (l *Lexer) errorAt(code diag.Code, msg string, t *Token) {
	if l.maxErrs > 0 && len(l.errs) > l.maxErrs {
		return
	}
	l.errs = append(l.errs, diag.New(code, msg, l.filepath, diag.Range{Start: ToPos(t.Start), End: ToPos(t.End)}))
}

// lex scans the whole input, appends EOF, then attaches trivia and invalid
// tokens to their neighbours (lexer.ts lex).
func (l *Lexer) lex() ([]*Token, []diag.Diagnostic) {
	l.scanTokens()
	l.push(&Token{Kind: EOF, Start: l.start, End: l.current})
	l.gatherInvalid()
	return l.tokens, l.errs
}

// scanTokens is the main dispatch loop: one token (or trivia) per iteration.
func (l *Lexer) scanTokens() {
	for !l.isAtEnd() && !l.stopped {
		if l.cancel != nil && l.idx&4095 == 0 {
			select {
			case <-l.cancel:
				l.stopped = true
				l.errs = append(l.errs, diag.New(diag.CANCELED, "Compilation was canceled", l.filepath, diag.Range{Start: ToPos(l.current), End: ToPos(l.current)}))
				return
			default:
			}
		}
		c := l.advance()
		switch c {
		case ' ':
			l.addToken(Space, false)
		case '\r':
		case '\n':
			l.addToken(Newline, false)
		case '\t':
			l.addToken(Tab, false)
		case ',':
			l.addToken(Comma, false)
		case '(':
			l.addToken(LParen, false)
		case ')':
			l.addToken(RParen, false)
		case '[':
			l.addToken(LBracket, false)
		case ']':
			l.addToken(RBracket, false)
		case '{':
			l.addToken(LBrace, false)
		case '}':
			l.addToken(RBrace, false)
		case ';':
			l.addToken(Semicolon, false)
		case ':':
			l.addToken(Colon, false)
		case '\'':
			if l.match("''") {
				l.consumeUntil(StringLiteral, "'''", true, false, false, true)
			} else {
				l.consumeUntil(StringLiteral, "'", false, false, false, true)
			}
		case '"':
			l.consumeUntil(QuotedString, "\"", false, false, false, true)
		case '`':
			l.consumeUntil(FunctionExpression, "`", true, false, true, true)
		case '#':
			l.colorLiteral()
		case '/':
			if l.match("/") {
				l.consumeUntil(SingleLineComment, "\n", true, true, true, false)
			} else if l.match("*") {
				l.consumeUntil(MultilineComment, "*/", true, false, true, true)
			} else {
				l.operator(c)
			}
		case '*':
			l.addToken(Wildcard, false)
		default:
			switch {
			case IsOpChar(c):
				l.operator(c)
			case isAlphaOrUnderscore(c):
				l.identifier()
			case isDigit(c):
				l.numericLiteralOrIdentifier()
			default:
				l.addToken(Op, true)
				l.errorAt(diag.UNKNOWN_SYMBOL, "Unexpected token '"+string(c)+"'", l.createToken(Op, true))
			}
		}
		l.start = l.current
	}
}

// gatherInvalid removes invalid tokens from the stream, attaching each to
// the previous valid token (or the first valid one, when leading).
func (l *Lexer) gatherInvalid() {
	i := 0
	var leading []*Token
	for ; i < len(l.tokens) && l.tokens[i].Invalid; i++ {
		leading = append(leading, l.tokens[i])
	}
	prev := l.tokens[i]
	prev.PrependLeadingInvalid(leading)
	newList := l.tokens[:0] // in place: invalid tokens only ever shrink the list
	for ; i < len(l.tokens); i++ {
		tok := l.tokens[i]
		if tok.Invalid {
			prev.AppendTrailingInvalid(tok)
		} else {
			prev = tok
			newList = append(newList, tok)
		}
	}
	l.tokens = newList
}

// consumeUntil consumes characters until stop is encountered.
func (l *Lexer) consumeUntil(kind Kind, stop string, allowNewline, allowEOF, raw, consumeStop bool) {
	var sb strings.Builder
	for !l.isAtEnd() && (allowNewline || !l.peekIs('\n')) && !l.check(stop) {
		if l.peekIs('\\') && !raw {
			l.advance()
			sb.WriteString(l.escapedString())
		} else {
			sb.WriteRune(l.advance())
		}
	}
	if l.isAtEnd() && !allowEOF {
		tok := l.createToken(kind, true)
		l.push(tok)
		l.errorAt(diag.UNEXPECTED_EOF, "EOF reached while parsing", tok)
		return
	}
	if l.peekIs('\n') && !allowNewline {
		tok := l.createToken(kind, true)
		l.push(tok)
		l.errorAt(diag.UNEXPECTED_NEWLINE, "Invalid newline encountered while parsing", tok)
		return
	}
	if consumeStop {
		l.match(stop)
	}
	if !l.raw && isTriviaKind(kind) {
		l.trivia(kind, l.start, l.current)
		return
	}
	l.push(&Token{Kind: kind, Start: l.start, End: l.current, Value: sb.String()})
}

// identifier consumes the rest of an identifier (letters, marks, digits, _).
func (l *Lexer) identifier() {
	for !l.isAtEnd() {
		r, _ := l.peek(0)
		if !isAlphaNumeric(r) {
			break
		}
		l.advance()
	}
	l.addToken(Identifier, false)
}

// operator consumes a (possibly multi-character) operator starting with c:
// comparison ops, relationship ops with optional `?` on either side
// (`<>?`, `?-?`, ...) and the dependency arrows `->` / `<-`.
func (l *Lexer) operator(c rune) {
	pk := func(want rune) bool { return l.peekIs(want) }
	switch c {
	case '<':
		switch {
		case pk('='):
			l.advance()
		case pk('>'):
			l.advance()
			if pk('?') {
				l.advance()
			}
		case pk('?'):
			l.advance()
		case pk('-'):
			l.advance()
		}
	case '>':
		if pk('=') {
			l.advance()
		} else if pk('?') {
			l.advance()
		}
	case '-':
		if pk('>') {
			l.advance()
		} else if pk('?') {
			l.advance()
		}
	case '?':
		switch {
		case pk('-'):
			l.advance()
			if pk('?') {
				l.advance()
			}
		case pk('<'):
			l.advance()
			if pk('>') {
				l.advance()
				if pk('?') {
					l.advance()
				}
			} else if pk('?') {
				l.advance()
			}
		case pk('>'):
			l.advance()
			if pk('?') {
				l.advance()
			}
		}
	case '=':
		if pk('=') {
			l.advance()
		}
	case '!':
		if pk('=') {
			l.advance()
		}
	}
	l.addToken(Op, false)
}

// numericLiteralOrIdentifier accepts identifiers starting with digits (they
// must contain a letter or underscore) and numbers with scientific notation.
func (l *Lexer) numericLiteralOrIdentifier() {
	nDots := 0
	if l.isAtEnd() {
		l.addToken(NumericLiteral, false)
		return
	}
	for !l.isAtEnd() {
		isDot := l.check(".")
		if isDot {
			nDots++
		}
		if nDots > 1 {
			break
		}
		if l.check("e") || l.check("E") {
			after, okAfter := l.peek(1)
			hasSign := okAfter && (after == '+' || after == '-')
			digit, okDigit := after, okAfter
			if hasSign {
				digit, okDigit = l.peek(2)
			}
			if okDigit && isDigit(digit) {
				l.advance()
				if hasSign {
					l.advance()
				}
				for !l.isAtEnd() {
					r, _ := l.peek(0)
					if !isDigit(r) {
						break
					}
					l.advance()
				}
				if l.isAtEnd() {
					l.addToken(NumericLiteral, false)
					return
				}
				if r, _ := l.peek(0); !isAlphaNumeric(r) {
					l.addToken(NumericLiteral, false)
					return
				}
				break
			}
			break
		}
		if !isDot && l.idx == len(l.runes)-1 {
			l.advance()
			l.addToken(NumericLiteral, false)
			return
		}
		r, _ := l.peek(0)
		if !isDot && !isAlphaNumeric(r) {
			l.addToken(NumericLiteral, false)
			return
		}
		if !isDot && !isDigit(r) {
			break
		}
		l.advance()
	}
	if nDots > 0 {
		for !l.isAtEnd() {
			r, _ := l.peek(0)
			if r != '.' && !isAlphaNumeric(r) {
				break
			}
			l.advance()
		}
		tok := l.createToken(NumericLiteral, true)
		l.push(tok)
		l.errorAt(diag.UNKNOWN_TOKEN, "Invalid number", tok)
		return
	}
	for !l.isAtEnd() {
		r, _ := l.peek(0)
		if !isAlphaNumeric(r) {
			break
		}
		l.advance()
	}
	l.push(l.createToken(Identifier, false))
}

// colorLiteral consumes `#` followed by alphanumerics; validity (3 or 6 hex
// digits) is checked later by the model.
func (l *Lexer) colorLiteral() {
	for !l.isAtEnd() {
		r, _ := l.peek(0)
		if !isAlphaNumeric(r) {
			break
		}
		l.advance()
	}
	l.addToken(ColorLiteral, false)
}

// escapedString decodes one escape sequence after a backslash. Unknown
// escapes yield the character itself; `\u` needs exactly four hex digits.
func (l *Lexer) escapedString() string {
	prev := Position{Column: l.current.Column - 1, Offset: l.current.Offset - 1, Line: l.current.Line, Byte: l.current.Byte - 1}
	if l.isAtEnd() {
		return "\\"
	}
	c := l.advance()
	switch c {
	case '\r':
		if l.peekIs('\n') {
			l.advance()
		}
		return ""
	case '\n':
		return ""
	case 't':
		return "\t"
	case 'n':
		return "\n"
	case '\\':
		return "\\"
	case 'r':
		return "\r"
	case '\'':
		return "'"
	case '"':
		return "\""
	case '`':
		return "`"
	case '0':
		return "\x00"
	case 'b':
		return "\b"
	case 'v':
		return "\v"
	case 'f':
		return "\f"
	case ' ':
		return "\\ "
	case 'u':
		hex := ""
		for i := 0; i <= 3; i++ {
			r, ok := l.peek(0)
			if !ok || !isAlphaNumeric(r) {
				tok := &Token{Kind: StringLiteral, Start: prev, End: l.current, Value: "\\u" + hex, Invalid: true}
				l.errorAt(diag.INVALID_ESCAPE_SEQUENCE, "Invalid unicode escape sequence '\\u"+hex+"', only unicode escape sequences of the form '\\uHHHH' where H is a hexadecimal number are allowed", tok)
				return "\\u" + hex
			}
			hex += string(l.advance())
		}
		v, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			// JS parseInt parses the leading hex prefix; String.fromCharCode(NaN) is "\u0000".
			v = parseHexPrefix(hex)
		}
		return string(rune(v))
	default:
		return string(c)
	}
}

// parseHexPrefix mimics JavaScript parseInt(s, 16): it parses the longest
// hexadecimal prefix and returns 0 when there is none.
func parseHexPrefix(s string) uint64 {
	var v uint64
	for _, c := range s {
		var d uint64
		switch {
		case c >= '0' && c <= '9':
			d = uint64(c - '0')
		case c >= 'a' && c <= 'f':
			d = uint64(c-'a') + 10
		case c >= 'A' && c <= 'F':
			d = uint64(c-'A') + 10
		default:
			return v
		}
		v = v*16 + d
	}
	return v
}

// isAlphaOrUnderscore matches \p{L}, \p{M} and `_` (identifier start).
func isAlphaOrUnderscore(c rune) bool {
	return c == '_' || unicode.IsLetter(c) || unicode.Is(unicode.M, c)
}

// isDigit matches ASCII digits only.
func isDigit(c rune) bool { return c >= '0' && c <= '9' }

// isAlphaNumeric matches identifier continuation characters.
func isAlphaNumeric(c rune) bool { return isAlphaOrUnderscore(c) || isDigit(c) }

// IsHexChar reports whether c is a hexadecimal digit.
func IsHexChar(c rune) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
