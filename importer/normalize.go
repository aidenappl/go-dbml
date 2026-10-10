package importer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aidenappl/go-dbml/export"
	"github.com/aidenappl/go-dbml/model"
)

var uuidOssp = regexp.MustCompile(`(?i)\b(?:public\.)?uuid_generate_v4\s*\(\s*\)`)

// dRails marks expressions from schema.rb: SQL text whose dialect is not
// known (a Rails app may run on MySQL or PostgreSQL).
const dRails dialect = 100

// Clean-mode expression normalisation.
//
// Every raw SQL expression the importers keep (expression defaults, CHECK
// expressions, expression index parts, expression record values,
// generated-column expressions) is re-lexed with the source dialect's rules
// and rewritten into the canonical form DBML expressions use:
//
//   - string literals become '...' with ' doubled ("..." MySQL strings,
//     E'...' and $tag$...$tag$ PostgreSQL strings included). A value
//     containing a backslash cannot be written so that every dialect reads
//     it the same: it keeps the source convention in a form whose extent is
//     the same in every lexical mode of the source (MySQL: '\\' doubled,
//     never \'; PostgreSQL: E'...' with \\ doubled; SQL Server: verbatim)
//     and the import reports a warning; from schema.rb (dialect unknown)
//     such an expression is dropped;
//   - `x` and [x] identifiers become bare names when they are plain
//     identifiers and "x" otherwise; "x" is kept;
//   - comments, unterminated quotes, unbalanced parentheses and ';' make
//     the expression unusable: it is dropped with a warning.
//
// The exporters re-check every fragment for the target dialect (see
// export.ErrUnsafeExpression); this pass makes imported expressions mean
// what they meant in the source and pass that check.

// normalizeExpr rewrites expr from source dialect d. ok=false means the
// expression must be dropped; warn is a non-fatal note.
func normalizeExpr(expr string, d dialect) (out string, warn string, ok bool) {
	var b strings.Builder
	depth := 0
	backslash := false
	n := len(expr)
	for i := 0; i < n; i++ {
		c := expr[i]
		switch {
		case c == '\'' || (c == '"' && d == dMySQL):
			// MySQL: '...' and "..." strings with backslash escapes.
			// Others: '...' with doubling (PostgreSQL E'...' below).
			prefixE := d == dPostgres && i > 0 && (expr[i-1] == 'E' || expr[i-1] == 'e') && (i < 2 || !isIdentChar(expr[i-2]))
			bs := d == dMySQL || prefixE
			j := scanQuote(expr, i, c, bs)
			if j < 0 {
				return "", "unterminated string literal", false
			}
			var v string
			switch {
			case d == dMySQL:
				v = unescapeMySQLString(expr[i : j+1])
			case prefixE:
				v = unescapeStdString("E" + expr[i:j+1])
			default:
				v = strings.ReplaceAll(expr[i+1:j], "''", "'")
			}
			if prefixE {
				// drop the E already copied
				s := b.String()
				b.Reset()
				b.WriteString(s[:len(s)-1])
			}
			if strings.IndexByte(v, 0) >= 0 {
				return "", "NUL byte in a string literal", false
			}
			if strings.IndexByte(v, '\\') >= 0 {
				switch d {
				case dRails:
					return "", "string literal with a backslash (dialect unknown, meaning ambiguous)", false
				case dMySQL:
					b.WriteString("'" + strings.ReplaceAll(strings.ReplaceAll(v, `\`, `\\`), "'", "''") + "'")
				case dPostgres:
					if b.Len() > 0 && isIdentChar(b.String()[b.Len()-1]) {
						return "", "backslash in a prefixed string literal", false
					}
					b.WriteString("E'" + strings.ReplaceAll(strings.ReplaceAll(v, `\`, `\\`), "'", "''") + "'")
				default:
					b.WriteString(expr[i : j+1])
				}
				backslash = true
			} else {
				b.WriteString("'" + strings.ReplaceAll(v, "'", "''") + "'")
			}
			i = j
		case c == '$' && d == dPostgres:
			k := i + 1
			for k < n && (expr[k] == '_' || isAlpha(expr[k]) || k > i+1 && isDigit(expr[k])) {
				k++
			}
			if k < n && expr[k] == '$' && !(i > 0 && isIdentChar(expr[i-1])) {
				tag := expr[i : k+1]
				end := strings.Index(expr[k+1:], tag)
				if end < 0 {
					return "", "unterminated dollar-quoted string", false
				}
				v := expr[k+1 : k+1+end]
				if strings.IndexByte(v, '\\') >= 0 {
					b.WriteString("E'" + strings.ReplaceAll(strings.ReplaceAll(v, `\`, `\\`), "'", "''") + "'")
					backslash = true
				} else {
					b.WriteString("'" + strings.ReplaceAll(v, "'", "''") + "'")
				}
				i = k + end + len(tag)
				continue
			}
			b.WriteByte(c)
		case c == '"' || c == '`' || (c == '[' && d == dMSSQL):
			closeCh := c
			if c == '[' {
				closeCh = ']'
			}
			j := scanQuote(expr, i, closeCh, false)
			if j < 0 {
				return "", "unterminated quoted identifier", false
			}
			name := unquoteIdent(expr[i : j+1])
			if c == '"' {
				// already standard (and case-sensitive in PostgreSQL)
				b.WriteString(expr[i : j+1])
			} else if simpleIdent.MatchString(name) && !strings.Contains(name, "$") && !export.IsMySQLReserved(name) {
				b.WriteString(name)
			} else {
				b.WriteString(`"` + strings.ReplaceAll(name, `"`, `""`) + `"`)
			}
			i = j
		case c == '-' && i+1 < n && expr[i+1] == '-', c == '/' && i+1 < n && expr[i+1] == '*', c == '#' && d == dMySQL:
			return "", "comment inside the expression", false
		case c == ';':
			return "", "';' inside the expression", false
		case c == 0:
			return "", "NUL byte", false
		case c == '\\':
			return "", "backslash outside a string literal", false
		case c == '(':
			depth++
			b.WriteByte(c)
		case c == ')':
			depth--
			if depth < 0 {
				return "", "unbalanced parentheses", false
			}
			b.WriteByte(c)
		default:
			b.WriteByte(c)
		}
	}
	if depth != 0 {
		return "", "unbalanced parentheses", false
	}
	if d == dPostgres || d == dRails {
		// uuid-ossp's uuid_generate_v4() is gen_random_uuid() since
		// PostgreSQL 13 (no extension needed)
		b2 := uuidOssp.ReplaceAllString(b.String(), "gen_random_uuid()")
		b.Reset()
		b.WriteString(b2)
	}
	if backslash {
		warn = "string literal with a backslash keeps the source dialect's escaping; another dialect may read it differently or refuse it"
	}
	return strings.TrimSpace(b.String()), warn, true
}

// scanQuote returns the index of the quote closing the token at s[i], or
// -1 (doubling escapes the quote; bs makes a backslash escape a byte).
func scanQuote(s string, i int, q byte, bs bool) int {
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

func isAlpha(c byte) bool { return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= 0x80 }
func isIdentChar(c byte) bool {
	return c == '_' || c == '$' || c == '&' || isAlpha(c) || isDigit(c)
}

// normalizeModel normalises every expression in db (clean mode) and
// returns the warnings.
func normalizeModel(db *model.Database, d dialect) []string {
	var warns []string
	norm := func(where, e string) (string, bool) {
		out, w, ok := normalizeExpr(e, d)
		if !ok {
			warns = append(warns, fmt.Sprintf("%s: dropped (%s): %s", where, w, clip(e)))
			return "", false
		}
		if w != "" {
			warns = append(warns, where+": "+w)
		}
		return out, true
	}
	for _, t := range db.Tables {
		tn := t.Name
		if t.SchemaName != nil {
			tn = *t.SchemaName + "." + tn
		}
		for _, f := range t.Fields {
			where := "column " + tn + "." + f.Name
			if df := f.DBDefault; df != nil && df.Type == "expression" {
				if s, ok := df.Value.(string); ok {
					if out, ok := norm("default of "+where, s); ok {
						df.Value = out
					} else {
						f.DBDefault = nil
					}
				}
			}
			f.Checks = normChecks(f.Checks, "check on "+where, norm)
			if g := f.Metadata[model.MetaGenerated]; g != "" {
				if out, ok := norm("generated "+where, g); ok {
					f.Metadata[model.MetaGenerated] = out
				} else {
					delete(f.Metadata, model.MetaGenerated)
					delete(f.Metadata, model.MetaGeneratedStored)
				}
			}
		}
		t.Checks = normChecks(t.Checks, "check on table "+tn, norm)
		for _, ix := range t.Indexes {
			cols := ix.Columns[:0]
			for _, c := range ix.Columns {
				if c.Type == "expression" {
					out, ok := norm("index expression on table "+tn, c.Value)
					if !ok {
						continue
					}
					c.Value = out
				}
				cols = append(cols, c)
			}
			ix.Columns = cols
		}
	}
	for _, r := range db.Records {
		for _, row := range r.Values {
			for _, v := range row {
				if v == nil || v.Type != "expression" {
					continue
				}
				if s, ok := v.Value.(string); ok {
					if out, ok := norm("record value in "+r.TableName, s); ok {
						v.Value = out
					} else {
						v.Value, v.Type = nil, "expression"
					}
				}
			}
		}
	}
	return warns
}

func normChecks(cs []*model.Check, where string, norm func(string, string) (string, bool)) []*model.Check {
	out := cs[:0]
	for _, c := range cs {
		if e, ok := norm(where, c.Expression); ok {
			c.Expression = e
			out = append(out, c)
		}
	}
	return out
}

func clip(s string) string {
	if len(s) > 60 {
		return s[:60] + "…"
	}
	return s
}
