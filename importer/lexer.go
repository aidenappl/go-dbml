// Portions of this file mirror behaviour of holistics/dbml packages/dbml-core
// src/parse/ANTLR (the MySQL/PostgreSQL/T-SQL lexers' token rules, MySQL
// DELIMITER handling, ANTLR getText/getOriginalText semantics)
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: re-implemented in Go as a hand-written streaming tokenizer.

package importer

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// tokKind classifies SQL tokens.
type tokKind int

const (
	tEOF    tokKind = iota
	tWord           // bare identifier or keyword
	tQIdent         // quoted identifier: `x`, "x" (postgres/mssql), [x] (mssql)
	tString         // string literal, including prefixes: 'x', "x" (mysql), N'x', E'x', $$x$$, _utf8'x'
	tNumber         // numeric literal
	tVar            // @x, @@x, $1
	tPunct          // operators and punctuation
	tTerm           // statement terminator (`;` or a MySQL DELIMITER, or a T-SQL GO)
)

// token is one lexical token. Text is the exact source text; ANTLR's
// ctx.getText() of a rule is the concatenation of its token texts (no
// whitespace, no comments), which getText reproduces.
type token struct {
	kind       tokKind
	text       string
	start, end int // byte offsets into the source
}

func (t token) upper() string { return strings.ToUpper(t.text) }

// is reports whether t is the bare word kw (case-insensitive).
func (t token) is(kw string) bool { return t.kind == tWord && strings.EqualFold(t.text, kw) }

func (t token) isPunct(p string) bool { return t.kind == tPunct && t.text == p }

type dialect int

const (
	dMySQL dialect = iota
	dPostgres
	dMSSQL
)

// lex tokenizes src for the dialect. Comments and whitespace are dropped.
func lex(src string, d dialect) []token {
	l := &lexer{src: src, d: d, delim: ";"}
	l.run()
	return l.toks
}

// lexTo streams tokens of src to out (ending with EOF) without keeping
// them, so callers can drop statements they do not need.
func lexTo(src string, d dialect, out func(token)) {
	l := &lexer{src: src, d: d, delim: ";", out: out}
	l.run()
}

type lexer struct {
	src   string
	pos   int
	d     dialect
	toks  []token
	delim string
	// startOfLine is true when only whitespace was seen since the last newline.
	startOfLine bool
	// out, when set, receives tokens instead of toks.
	out func(token)
}

func (l *lexer) push(t token) {
	if l.out != nil {
		l.out(t)
		return
	}
	l.toks = append(l.toks, t)
}

func (l *lexer) peek(off int) byte {
	if l.pos+off < len(l.src) {
		return l.src[l.pos+off]
	}
	return 0
}

func (l *lexer) emit(k tokKind, start int) {
	l.push(token{kind: k, text: l.src[start:l.pos], start: start, end: l.pos})
}

func (l *lexer) run() {
	l.startOfLine = true
	for l.pos < len(l.src) {
		c := l.src[l.pos]
		start := l.pos
		switch {
		case c == '\n':
			l.pos++
			l.startOfLine = true
			continue
		case c == ' ' || c == '\t' || c == '\r' || c == '\f' || c == '\v':
			l.pos++
			continue
		}
		// MySQL client DELIMITER command (only at the start of a line).
		if l.d == dMySQL && l.startOfLine && l.hasWordAt(l.pos, "DELIMITER") {
			l.delimiterCommand()
			continue
		}
		// T-SQL batch separator GO on its own line.
		if l.d == dMSSQL && l.startOfLine && l.hasWordAt(l.pos, "GO") && l.restOfLineBlankAfter(l.pos+2) {
			l.pos += 2
			l.emit(tTerm, start)
			continue
		}
		l.startOfLine = false
		// custom MySQL delimiter
		if l.d == dMySQL && l.delim != ";" && strings.HasPrefix(l.src[l.pos:], l.delim) {
			l.pos += len(l.delim)
			l.emit(tTerm, start)
			continue
		}
		switch {
		case c == '-' && l.peek(1) == '-':
			l.lineComment()
		case c == '#' && l.d == dMySQL:
			l.lineComment()
		case c == '/' && l.peek(1) == '*':
			l.blockComment()
		case c == ';':
			l.pos++
			if (l.d == dMySQL && l.delim != ";") || l.d == dMSSQL {
				l.emit(tPunct, start)
			} else {
				l.emit(tTerm, start)
			}
		case c == '\'':
			l.quoted('\'', l.d == dMySQL)
			l.emit(tString, start)
		case c == '"':
			l.quoted('"', l.d == dMySQL)
			if l.d == dMySQL {
				l.emit(tString, start)
			} else {
				l.emit(tQIdent, start)
			}
		case c == '`' && l.d == dMySQL:
			l.quoted('`', false)
			l.emit(tQIdent, start)
		case c == '[' && l.d == dMSSQL:
			l.bracketIdent()
			l.emit(tQIdent, start)
		case c == '$' && l.d == dPostgres && l.dollarQuote():
			l.emit(tString, start)
		case c == '$' && l.d == dPostgres && isDigit(l.peek(1)):
			l.pos++
			for isDigit(l.peek(0)) {
				l.pos++
			}
			l.emit(tVar, start)
		case c == '@' && (l.d == dMySQL || l.d == dMSSQL):
			l.pos++
			if l.peek(0) == '@' {
				l.pos++
			}
			if l.peek(0) == '`' || l.peek(0) == '\'' || l.peek(0) == '"' {
				l.quoted(l.peek(0), false)
			} else {
				l.word()
			}
			l.emit(tVar, start)
		case isDigit(c) || (c == '.' && isDigit(l.peek(1))):
			l.number()
			// MySQL identifiers may start with digits (e.g. 1col)
			if l.d == dMySQL && l.pos < len(l.src) {
				r, _ := utf8.DecodeRuneInString(l.src[l.pos:])
				if r == '_' || unicode.IsLetter(r) {
					l.word()
					l.emit(tWord, start)
					continue
				}
			}
			l.emit(tNumber, start)
		default:
			r, size := utf8.DecodeRuneInString(l.src[l.pos:])
			if r == '_' || unicode.IsLetter(r) || r >= 0x80 {
				// string prefixes: N'..', E'..', B'..', X'..', U&'..', _charset'..'
				if l.stringPrefix() {
					l.emit(tString, start)
					continue
				}
				l.word()
				l.emit(tWord, start)
				continue
			}
			l.pos += size
			// characters the reference MySQL lexer has no token for are
			// dropped (its lexer errors are not reported)
			if l.d == dMySQL && size == 1 && strings.IndexByte("[]{}?\\", c) >= 0 {
				continue
			}
			l.punct(c)
			l.emit(tPunct, start)
		}
	}
	l.push(token{kind: tEOF, start: len(l.src), end: len(l.src)})
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

func (l *lexer) hasWordAt(p int, w string) bool {
	if p+len(w) > len(l.src) || !strings.EqualFold(l.src[p:p+len(w)], w) {
		return false
	}
	if p+len(w) == len(l.src) {
		return true
	}
	c := l.src[p+len(w)]
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func (l *lexer) restOfLineBlankAfter(p int) bool {
	for p < len(l.src) {
		switch l.src[p] {
		case '\n':
			return true
		case ' ', '\t', '\r':
			p++
		case '-':
			return p+1 < len(l.src) && l.src[p+1] == '-'
		default:
			// "GO 2" (repeat count) is also a batch separator
			if isDigit(l.src[p]) {
				p++
				continue
			}
			return false
		}
	}
	return true
}

// delimiterCommand handles `DELIMITER x` (rest of line, quotes stripped).
func (l *lexer) delimiterCommand() {
	l.pos += len("DELIMITER")
	for l.pos < len(l.src) && (l.src[l.pos] == ' ' || l.src[l.pos] == '\t') {
		l.pos++
	}
	start := l.pos
	for l.pos < len(l.src) && l.src[l.pos] != '\n' && l.src[l.pos] != ' ' && l.src[l.pos] != '\t' && l.src[l.pos] != '\r' {
		l.pos++
	}
	d := l.src[start:l.pos]
	if len(d) >= 2 && (d[0] == '\'' || d[0] == '"' || d[0] == '`') && d[len(d)-1] == d[0] {
		d = d[1 : len(d)-1]
	}
	if d != "" {
		l.delim = d
	}
	for l.pos < len(l.src) && l.src[l.pos] != '\n' {
		l.pos++
	}
	// the end of a DELIMITER command terminates any pending statement
	l.push(token{kind: tTerm, text: "", start: l.pos, end: l.pos})
}

func (l *lexer) lineComment() {
	for l.pos < len(l.src) && l.src[l.pos] != '\n' {
		l.pos++
	}
}

func (l *lexer) blockComment() {
	l.pos += 2
	for l.pos < len(l.src) {
		if l.src[l.pos] == '*' && l.peek(1) == '/' {
			l.pos += 2
			return
		}
		l.pos++
	}
}

// quoted consumes a quoted run starting at l.pos (the opening quote). A
// doubled quote is an escaped quote; with backslash, \x escapes too.
func (l *lexer) quoted(q byte, backslash bool) {
	l.pos++
	for l.pos < len(l.src) {
		c := l.src[l.pos]
		if backslash && c == '\\' {
			l.pos += 2
			continue
		}
		if c == q {
			if l.peek(1) == q {
				l.pos += 2
				continue
			}
			l.pos++
			return
		}
		l.pos++
	}
	if l.pos > len(l.src) {
		l.pos = len(l.src)
	}
}

func (l *lexer) bracketIdent() {
	l.pos++
	for l.pos < len(l.src) {
		if l.src[l.pos] == ']' {
			if l.peek(1) == ']' {
				l.pos += 2
				continue
			}
			l.pos++
			return
		}
		l.pos++
	}
}

// dollarQuote consumes a PostgreSQL $tag$...$tag$ string.
func (l *lexer) dollarQuote() bool {
	end := l.pos + 1
	for end < len(l.src) && (l.src[end] == '_' || isDigit(l.src[end]) || (l.src[end]|0x20 >= 'a' && l.src[end]|0x20 <= 'z')) {
		end++
	}
	if end >= len(l.src) || l.src[end] != '$' {
		return false
	}
	tag := l.src[l.pos : end+1]
	if i := strings.Index(l.src[end+1:], tag); i >= 0 {
		l.pos = end + 1 + i + len(tag)
	} else {
		l.pos = len(l.src)
	}
	return true
}

func (l *lexer) stringPrefix() bool {
	rest := l.src[l.pos:]
	head := rest
	if len(head) > 3 {
		head = head[:3]
	}
	up := strings.ToUpper(head)
	switch {
	case l.d == dPostgres && strings.HasPrefix(up, "U&'"):
		l.pos += 2
		l.quoted('\'', false)
		return true
	case len(rest) >= 2 && rest[1] == '\'':
		switch up[0] {
		case 'N', 'X', 'B':
			l.pos++
			l.quoted('\'', l.d == dMySQL)
			return true
		case 'E':
			if l.d == dPostgres {
				l.pos++
				l.quoted('\'', true)
				return true
			}
		}
	case l.d == dMySQL && rest[0] == '_':
		// _utf8mb4'...'
		i := 1
		for i < len(rest) && (rest[i] == '_' || isDigit(rest[i]) || (rest[i]|0x20 >= 'a' && rest[i]|0x20 <= 'z')) {
			i++
		}
		if i < len(rest) && rest[i] == '\'' {
			l.pos += i
			l.quoted('\'', true)
			return true
		}
	}
	return false
}

func (l *lexer) word() {
	start := l.pos
	defer func() {
		if l.pos == start && l.pos < len(l.src) {
			_, size := utf8.DecodeRuneInString(l.src[l.pos:])
			l.pos += size
		}
	}()
	for l.pos < len(l.src) {
		r, size := utf8.DecodeRuneInString(l.src[l.pos:])
		if r == '_' || r == '$' || unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.Is(unicode.M, r) || (r == '#' && l.d == dMSSQL) {
			l.pos += size
			continue
		}
		break
	}
}

func (l *lexer) number() {
	if l.peek(0) == '0' && (l.peek(1) == 'x' || l.peek(1) == 'X') {
		l.pos += 2
		for isHex(l.peek(0)) {
			l.pos++
		}
		return
	}
	if l.d == dPostgres && l.peek(0) == '0' && (l.peek(1)|0x20 == 'o' || l.peek(1)|0x20 == 'b') && isDigit(l.peek(2)) {
		l.pos += 2
		for isDigit(l.peek(0)) || l.peek(0) == '_' {
			l.pos++
		}
		return
	}
	for isDigit(l.peek(0)) {
		l.pos++
	}
	if l.peek(0) == '.' && l.peek(1) != '.' {
		l.pos++
		for isDigit(l.peek(0)) {
			l.pos++
		}
	}
	if (l.peek(0) == 'e' || l.peek(0) == 'E') && (isDigit(l.peek(1)) || ((l.peek(1) == '+' || l.peek(1) == '-') && isDigit(l.peek(2)))) {
		l.pos += 2
		for isDigit(l.peek(0)) {
			l.pos++
		}
	}
}

func isHex(c byte) bool { return isDigit(c) || (c|0x20 >= 'a' && c|0x20 <= 'f') }

var multiPunct = []string{"::", "<=>", "<>", "!=", "<=", ">=", "||", "&&", "->>", "->", ":=", "<<", ">>", "#>>", "#>", "@>", "<@", "~*", "!~*", "!~", ".."}

func (l *lexer) punct(c byte) {
	start := l.pos - 1
	for _, p := range multiPunct {
		if strings.HasPrefix(l.src[start:], p) {
			l.pos = start + len(p)
			return
		}
	}
}
