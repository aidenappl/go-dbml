// Portions of this file mirror behaviour of holistics/dbml packages/dbml-core
// src/parse/ANTLR (the MySQL/PostgreSQL/T-SQL lexers' token rules, MySQL
// DELIMITER handling, ANTLR getText/getOriginalText semantics)
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: re-implemented in Go as a hand-written streaming tokenizer.

package importer

import (
	"strings"

	"github.com/aidenappl/go-dbml/model"
)

// cursor walks the tokens of one statement.
type cursor struct {
	src   string
	toks  []token // always ends with an EOF token
	i     int
	match []int32 // "(" index -> matching ")" index (lazy, see matchOf)
	// truncated: the splitter stopped storing this statement's tokens
	truncated bool
}

func newCursor(src string, toks []token) *cursor {
	if len(toks) == 0 || toks[len(toks)-1].kind != tEOF {
		end := len(src)
		if len(toks) > 0 {
			end = toks[len(toks)-1].end
		}
		toks = append(append([]token{}, toks...), token{kind: tEOF, start: end, end: end})
	}
	return &cursor{src: src, toks: toks}
}

func (c *cursor) peek() token { return c.at(0) }

func (c *cursor) at(n int) token {
	if c.i+n < len(c.toks) {
		return c.toks[c.i+n]
	}
	return c.toks[len(c.toks)-1]
}

func (c *cursor) eof() bool { return c.peek().kind == tEOF }

func (c *cursor) next() token {
	t := c.peek()
	if c.i < len(c.toks)-1 {
		c.i++
	}
	return t
}

// accept consumes the next token if it is the bare word kw.
func (c *cursor) accept(kw string) bool {
	if c.peek().is(kw) {
		c.next()
		return true
	}
	return false
}

// acceptSeq consumes the words kws if they all follow.
func (c *cursor) acceptSeq(kws ...string) bool {
	for i, kw := range kws {
		if !c.at(i).is(kw) {
			return false
		}
	}
	c.i += len(kws)
	return true
}

func (c *cursor) acceptPunct(p string) bool {
	if c.peek().isPunct(p) {
		c.next()
		return true
	}
	return false
}

// skipGroup skips a balanced (...) group starting at the current "(".
func (c *cursor) skipGroup() (from, to int) {
	from = c.i
	if c.peek().isPunct("(") {
		if m := c.matchOf(c.i); m >= 0 {
			c.i = m + 1
			return from, c.i
		}
		c.i = len(c.toks) - 1 // unclosed: up to EOF
		return from, c.i
	}
	depth := 0
	for !c.eof() {
		t := c.next()
		if t.isPunct("(") {
			depth++
		} else if t.isPunct(")") {
			depth--
			if depth <= 0 {
				break
			}
		}
	}
	return from, c.i
}

// matchOf returns the index of the ")" matching the "(" at i, or -1. The
// match table is built once per cursor (O(n), no recursion), so deeply
// nested input never costs more than linear time.
func (c *cursor) matchOf(i int) int {
	if c.match == nil {
		c.match = make([]int32, len(c.toks))
		var stack []int
		for j, t := range c.toks {
			c.match[j] = -1
			switch {
			case t.isPunct("("):
				stack = append(stack, j)
			case t.isPunct(")") && len(stack) > 0:
				o := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				c.match[o] = int32(j)
			}
		}
	}
	if i < 0 || i >= len(c.match) {
		return -1
	}
	return int(c.match[i])
}

// skipUntil advances to the first token (at paren depth 0) for which stop
// returns true, or EOF. It returns the skipped range [from, to).
func (c *cursor) skipUntil(stop func(t token) bool) (from, to int) {
	from = c.i
	depth := 0
	for !c.eof() {
		t := c.peek()
		if depth == 0 && stop(t) {
			break
		}
		if t.isPunct("(") || t.isPunct("[") && t.kind == tPunct {
			depth++
		} else if t.isPunct(")") || t.isPunct("]") {
			if depth == 0 {
				break
			}
			depth--
		}
		c.next()
	}
	return from, c.i
}

// inner returns the tokens strictly inside a group [from, to) returned by
// skipGroup (empty when the group is unclosed or empty).
func (c *cursor) inner(from, to int) []token {
	if to-1 <= from+1 || to > len(c.toks) {
		return nil
	}
	if !c.toks[to-1].isPunct(")") {
		return c.toks[from+1 : to]
	}
	return c.toks[from+1 : to-1]
}

// getText mirrors ANTLR ctx.getText(): token texts concatenated.
func (c *cursor) getText(from, to int) string {
	var b strings.Builder
	for i := from; i < to && i < len(c.toks); i++ {
		b.WriteString(c.toks[i].text)
	}
	return b.String()
}

// origText mirrors getOriginalText(ctx): the source span of the tokens.
func (c *cursor) origText(from, to int) string {
	if from >= to || from >= len(c.toks) {
		return ""
	}
	if to > len(c.toks) {
		to = len(c.toks)
	}
	return c.src[c.toks[from].start:c.toks[to-1].end]
}

// atTop reports a comma or closing paren (a definition boundary).
func boundary(t token) bool { return t.isPunct(",") || t.isPunct(")") || t.kind == tTerm }

// statement is one statement's tokens; truncated is set when the splitter
// stopped storing tokens (see splitter.maxTokens).
type statement struct {
	toks      []token
	truncated bool
}

// splitter groups a token stream into statements at terminators. For
// MySQL-style procedural bodies (CREATE ... PROCEDURE|FUNCTION|TRIGGER|
// EVENT with BEGIN ... END) a ";" inside BEGIN/END does not end the
// statement, and a top-level CREATE/ALTER starts a new statement (the
// reference grammars make the terminator optional). It needs no lookahead,
// so the importers stream: a statement is handed to emit as soon as it
// ends, and statements rejected by skip are never stored.
type splitter struct {
	compound bool
	// tsql: T-SQL batches (GO-separated) are split further into statements
	// at ';' and before CREATE / INSERT / EXEC at paren depth 0 — the
	// points where the batch walker starts a statement anyway — so an
	// SSMS-style dump of INSERTs streams instead of being held as one
	// batch. A batch that defines a routine or view (CREATE [OR ALTER]
	// PROCEDURE|PROC|FUNCTION|TRIGGER|VIEW, ALTER ...) is never split (its
	// body holds statements that are not schema) and never stored.
	tsql      bool
	batchN    int // tokens since the batch started
	batchKind int // 0 undecided/normal, -1 deciding (after CREATE/ALTER), 1 routine
	// skip is asked at a statement's first token whether to drop it.
	skip func(first token) bool
	// maxTokens caps the tokens stored for one statement (0 = no cap).
	maxTokens int
	emit      func(statement)

	cur       []token
	n         int // tokens seen in the current statement (stored or not)
	skipping  bool
	truncated bool
	depth     int
	paren     int
	routine   int // -1 undecided, 0 no, 1 yes
	pendBegin bool
	pendEnd   bool
}

func (sp *splitter) flush() {
	if sp.n > 0 && !sp.skipping {
		sp.emit(statement{toks: sp.cur, truncated: sp.truncated})
	}
	sp.cur, sp.n, sp.skipping, sp.truncated = nil, 0, false, false
	sp.depth, sp.paren, sp.routine, sp.pendBegin, sp.pendEnd = 0, 0, 0, false, false
}

func (sp *splitter) push(t token) {
	if t.kind == tEOF {
		sp.flush()
		return
	}
	if sp.tsql {
		sp.pushTSQL(t)
		return
	}
	if sp.pendBegin {
		sp.pendBegin = false
		if !(t.is("TRANSACTION") || t.is("WORK")) {
			sp.depth++
		}
	}
	if sp.pendEnd {
		sp.pendEnd = false
		if !(t.is("IF") || t.is("LOOP") || t.is("WHILE") || t.is("CASE") || t.is("REPEAT")) && sp.depth > 0 {
			sp.depth--
		}
	}
	if sp.compound && sp.n > 0 && sp.paren == 0 && sp.depth == 0 && (t.is("CREATE") || t.is("ALTER")) {
		sp.flush()
	}
	if t.kind == tTerm && (sp.depth == 0 || sp.routine != 1 || t.text != ";") {
		sp.flush()
		return
	}
	if t.isPunct("(") {
		sp.paren++
	} else if t.isPunct(")") && sp.paren > 0 {
		sp.paren--
	}
	if sp.compound {
		switch {
		case sp.n == 0 && t.is("CREATE"):
			sp.routine = -1
		case sp.routine == -1:
			switch {
			case t.is("PROCEDURE") || t.is("FUNCTION") || t.is("TRIGGER") || t.is("EVENT"):
				sp.routine = 1
			case t.is("TABLE") || t.is("INDEX") || t.is("VIEW") || sp.n >= 12:
				sp.routine = 0
			}
		case sp.routine == 1 && t.is("BEGIN"):
			sp.pendBegin = true
		case sp.routine == 1 && t.is("END"):
			sp.pendEnd = true
		}
	}
	if sp.n == 0 && sp.skip != nil && sp.skip(t) {
		sp.skipping = true
	}
	sp.n++
	if sp.skipping {
		return
	}
	if sp.maxTokens > 0 && len(sp.cur) >= sp.maxTokens {
		sp.truncated = true
		return
	}
	sp.cur = append(sp.cur, t)
}

// pushTSQL is push for T-SQL (see splitter.tsql).
func (sp *splitter) pushTSQL(t token) {
	if t.kind == tTerm { // GO
		sp.flush()
		sp.batchN, sp.batchKind = 0, 0
		return
	}
	switch {
	case sp.batchN == 0 && (t.is("CREATE") || t.is("ALTER")):
		sp.batchKind = -1
	case sp.batchKind == -1 && (t.is("OR") || t.is("ALTER")):
	case sp.batchKind == -1:
		switch strings.ToUpper(t.text) {
		case "PROCEDURE", "PROC", "FUNCTION", "TRIGGER", "VIEW":
			sp.batchKind = 1
			sp.skipping = true
		default:
			sp.batchKind = 0
		}
	}
	sp.batchN++
	if sp.batchKind != 1 && sp.paren == 0 {
		if t.isPunct(";") {
			sp.flush()
			return
		}
		if sp.n > 0 && sp.batchKind == 0 && (t.is("CREATE") || t.is("INSERT") || t.is("EXEC") || t.is("EXECUTE")) {
			sp.flush()
		}
	}
	if t.isPunct("(") {
		sp.paren++
	} else if t.isPunct(")") && sp.paren > 0 {
		sp.paren--
	}
	if sp.n == 0 && sp.skip != nil && sp.skip(t) {
		sp.skipping = true
	}
	sp.n++
	if sp.skipping {
		return
	}
	if sp.maxTokens > 0 && len(sp.cur) >= sp.maxTokens {
		sp.truncated = true
		return
	}
	sp.cur = append(sp.cur, t)
}

// forEachStatement lexes src and streams its statements to fn.
func forEachStatement(src string, d dialect, compound bool, skip func(token) bool, maxTokens int, fn func(statement)) {
	sp := &splitter{compound: compound, tsql: d == dMSSQL, skip: skip, maxTokens: maxTokens, emit: fn}
	lexTo(src, d, sp.push)
}

// splitStatements is forEachStatement collected into a slice (tests).
func splitStatements(toks []token, compound bool) [][]token {
	var out [][]token
	sp := &splitter{compound: compound, emit: func(st statement) { out = append(out, st.toks) }}
	for _, t := range toks {
		sp.push(t)
	}
	sp.push(token{kind: tEOF})
	return out
}

// unquote strips one layer of surrounding quotes/brackets (no unescaping,
// like the reference's slice(1, -1)).
func stripQuotes(s string) string {
	if len(s) >= 2 {
		return s[1 : len(s)-1]
	}
	return s
}

func strp(s string) *string { return &s }

func boolp(b bool) *bool { return &b }

// ---------------------------------------------------------------- raw model

// rawDB accumulates the reference's "raw database" JSON (tables, refs,
// enums, records) in model types; see toModel.
type rawDB struct {
	tables  []*model.Table
	refs    []*model.Ref
	enums   []*model.Enum
	records []*model.TableRecord
	idx     map[string]*model.Table
}

func (r *rawDB) findTable(schema *string, name string) *model.Table {
	return r.lookup("public", schema, name)
}

// addTable appends t and indexes it for the lookups below (first table
// per key wins, like the reference's Array.find).
func (r *rawDB) addTable(t *model.Table) {
	r.tables = append(r.tables, t)
	if r.idx == nil {
		r.idx = map[string]*model.Table{}
	}
	for _, k := range []string{rawKey(t.SchemaName, t.Name), "public\x00" + defKey("public", t.SchemaName, t.Name), "dbo\x00" + defKey("dbo", t.SchemaName, t.Name)} {
		if _, ok := r.idx[k]; !ok {
			r.idx[k] = t
		}
	}
}

func rawKey(schema *string, name string) string {
	if schema == nil {
		return "\x01\x00" + name
	}
	return "\x02" + *schema + "\x00" + name
}

func defKey(def string, schema *string, name string) string {
	if schema != nil {
		def = *schema
	}
	return def + "\x00" + name
}

// lookup finds a table, a nil schema meaning def ("public" or "dbo").
func (r *rawDB) lookup(def string, schema *string, name string) *model.Table {
	return r.idx[def+"\x00"+defKey(def, schema, name)]
}

// strict finds a table whose schema is exactly schema (nil only matches nil).
func (r *rawDB) strict(schema *string, name string) *model.Table {
	return r.idx[rawKey(schema, name)]
}

func (r *rawDB) toModel() *model.Database {
	return &model.Database{Tables: r.tables, Refs: r.refs, Enums: r.enums, Records: r.records}
}

func findField(t *model.Table, name string) *model.Column {
	for _, f := range t.Fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// schemaOf returns names[len-2] or nil (the reference's getTableNames).
func lastTwo(names []string) (schema *string, name string) {
	if len(names) == 0 {
		return nil, ""
	}
	name = names[len(names)-1]
	if len(names) > 1 {
		schema = strp(names[len(names)-2])
	}
	return schema, name
}
