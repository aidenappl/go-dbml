package export

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Raw SQL fragments: column types, expression defaults, CHECK expressions,
// expression index parts, expression record values and generated-column
// expressions are SQL text, not names or strings, so they cannot be
// escaped and are emitted as written (like @dbml/core). Before emitting
// one, the exporter lexes it with the target dialect's own rules — in
// every lexical mode the server may be in: MySQL with and without
// NO_BACKSLASH_ESCAPES / ANSI_QUOTES, PostgreSQL with
// standard_conforming_strings on and off — and refuses it unless it is one
// balanced expression in all of them: no ';', no comment, no unterminated
// string / quoted identifier / dollar quote, balanced parentheses, no
// backslash outside a string literal, no client batch or command line (GO,
// DELIMITER), no NUL byte. A string such as 'a\' (MySQL-style escape)
// therefore passes for MySQL only when its end is the same with and
// without backslash escapes. A refused fragment fails the whole export
// with an *UnsafeExpressionError; nothing is ever silently rewritten.

// ErrUnsafeExpression is wrapped by every *UnsafeExpressionError.
var ErrUnsafeExpression = errors.New("unsafe SQL expression")

// UnsafeExpressionError reports a raw SQL fragment that is not a single
// balanced expression in the target dialect.
type UnsafeExpressionError struct {
	Dialect  Dialect
	Where    string // e.g. `check on "t"."a"`
	Fragment string
	Reason   string
}

func (e *UnsafeExpressionError) Error() string {
	f := e.Fragment
	if len(f) > 80 {
		f = f[:80] + "…"
	}
	return fmt.Sprintf("export: unsafe %s for %s (%s): %q", e.Where, e.Dialect, e.Reason, f)
}

// Unwrap returns ErrUnsafeExpression.
func (e *UnsafeExpressionError) Unwrap() error { return ErrUnsafeExpression }

// unsafe aborts an export; SQL's recover turns it into the error.
func unsafe(d Dialect, where, frag, reason string) {
	panic(&UnsafeExpressionError{Dialect: d, Where: where, Fragment: frag, Reason: reason})
}

// safeExpr validates a DBML expression for dialect d or aborts the export.
func safeExpr(expr string, d Dialect, where string) string {
	if d == MySQL {
		expr = mysqlQuoteReserved(expr)
	}
	if reason := checkFragment(expr, d); reason != "" {
		unsafe(d, where, expr, reason)
	}
	return expr
}

// safeType validates an already rendered column type or aborts the export.
func safeType(typ string, d Dialect, where string) string {
	if reason := checkFragment(typ, d); reason != "" {
		unsafe(d, where, typ, reason)
	}
	return typ
}

func isIdentByte(c byte) bool {
	return c == '_' || c == '$' || c == '&' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= 0x80
}

// scanQuoted returns the index of the closing quote q of a quoted token
// starting at s[i] (doubling escapes the quote; with bs, a backslash
// escapes the next byte), or -1 when it is unterminated.
func scanQuoted(s string, i int, q byte, bs bool) int {
	for j := i + 1; j < len(s); j++ {
		switch {
		case bs && s[j] == '\\':
			j++
		case s[j] == q:
			if j+1 < len(s) && s[j+1] == q {
				j++
				continue
			}
			return j
		}
	}
	return -1
}

var (
	goLine        = regexp.MustCompile(`(?im)^[ \t]*go([ \t]+[0-9]+)?[ \t]*\r?$`)
	delimiterLine = regexp.MustCompile(`(?im)^[ \t]*delimiter\b`)
)

// checkFragment lexes text with d's rules (in every lexical mode the
// server may run in) and returns why it is not one balanced expression,
// or "".
func checkFragment(text string, d Dialect) string {
	if strings.IndexByte(text, 0) >= 0 {
		return "NUL byte"
	}
	switch d {
	case MySQL:
		if delimiterLine.MatchString(text) {
			return "DELIMITER line"
		}
		if r := lexCheck(text, d, true); r != "" {
			return r
		}
		return lexCheck(text, d, false)
	case Postgres:
		if r := lexCheck(text, d, false); r != "" {
			return r
		}
		return lexCheck(text, d, true)
	case MSSQL:
		if goLine.MatchString(text) {
			return "GO batch separator line"
		}
		return lexCheck(text, d, false)
	}
	return "unknown dialect"
}

// lexCheck is one lexical pass. alt selects the alternative mode: MySQL
// backslash escapes on (default) vs off; PostgreSQL
// standard_conforming_strings on vs off.
func lexCheck(s string, d Dialect, alt bool) string {
	depth := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\'':
			bs := false
			switch d {
			case MySQL:
				bs = alt
			case Postgres:
				// E'..' always escapes; plain '..' when
				// standard_conforming_strings is off.
				e := i > 0 && (s[i-1] == 'E' || s[i-1] == 'e') && (i < 2 || !isIdentByte(s[i-2]))
				bs = e || alt
			}
			j := scanQuoted(s, i, '\'', bs)
			if j < 0 {
				return "unterminated string literal"
			}
			i = j
		case c == '"':
			bs := d == MySQL && alt // a MySQL "..." is a string
			j := scanQuoted(s, i, '"', bs)
			if j < 0 {
				return "unterminated quoted identifier"
			}
			i = j
		case c == '`' && d == MySQL:
			j := scanQuoted(s, i, '`', false)
			if j < 0 {
				return "unterminated quoted identifier"
			}
			i = j
		case c == '[' && d == MSSQL:
			j := scanQuoted(s, i, ']', false)
			if j < 0 {
				return "unterminated bracketed identifier"
			}
			i = j
		case c == '$' && d == Postgres:
			// $tag$ ... $tag$ dollar quoting ($1 is a parameter).
			k := i + 1
			for k < len(s) && (s[k] == '_' || s[k] >= 'a' && s[k] <= 'z' || s[k] >= 'A' && s[k] <= 'Z' || s[k] >= 0x80 || k > i+1 && s[k] >= '0' && s[k] <= '9') {
				k++
			}
			if k < len(s) && s[k] == '$' && !(i > 0 && isIdentByte(s[i-1])) {
				tag := s[i : k+1]
				end := strings.Index(s[k+1:], tag)
				if end < 0 {
					return "unterminated dollar-quoted string"
				}
				i = k + 1 + end + len(tag) - 1
			}
		case c == '-' && i+1 < len(s) && s[i+1] == '-',
			c == '/' && i+1 < len(s) && s[i+1] == '*',
			c == '#' && d == MySQL:
			return "comment"
		case c == ';':
			return "statement separator ';'"
		case c == '\\':
			return "backslash outside a string literal"
		case c == '(':
			depth++
		case c == ')':
			depth--
			if depth < 0 {
				return "unbalanced parentheses"
			}
		case c == '{' && d == MSSQL:
			// ODBC escape sequences; never needed in DDL fragments
			return "ODBC escape"
		}
	}
	if depth != 0 {
		return "unbalanced parentheses"
	}
	return ""
}
