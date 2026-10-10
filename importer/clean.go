package importer

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/aidenappl/go-dbml/model"
)

// Clean-mode helpers: undo SQL quoting so names and strings carry their
// real values, and make expressions safe inside DBML backticks.

// unescapeMySQLString decodes a MySQL string literal token: optional
// N / _charset prefix, '...' or "..." with doubled quotes and backslash
// escapes.
func unescapeMySQLString(text string) string {
	i := strings.IndexAny(text, `'"`)
	if i < 0 || len(text)-i < 2 {
		return text
	}
	q := text[i]
	body := text[i+1 : len(text)-1]
	var b strings.Builder
	for j := 0; j < len(body); j++ {
		c := body[j]
		switch {
		case c == '\\' && j+1 < len(body):
			j++
			switch body[j] {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case 'r':
				b.WriteByte('\r')
			case 'b':
				b.WriteByte('\b')
			case '0':
				b.WriteByte(0)
			case 'Z':
				b.WriteByte(26)
			case '%', '_':
				b.WriteByte('\\')
				b.WriteByte(body[j])
			default:
				b.WriteByte(body[j])
			}
		case c == q && j+1 < len(body) && body[j+1] == q:
			b.WriteByte(q)
			j++
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

// unescapeStdString decodes a PostgreSQL / SQL Server string literal:
// optional N / E / U& prefix, doubled quotes, E” backslash escapes,
// $tag$...$tag$ bodies.
func unescapeStdString(text string) string {
	if strings.HasPrefix(text, "$") {
		if i := strings.Index(text[1:], "$"); i >= 0 {
			tag := text[:i+2]
			if len(text) >= 2*len(tag) && strings.HasSuffix(text, tag) {
				return text[len(tag) : len(text)-len(tag)]
			}
		}
		return text
	}
	i := strings.IndexByte(text, '\'')
	if i < 0 || len(text)-i < 2 {
		return text
	}
	prefix := strings.ToUpper(text[:i])
	body := strings.ReplaceAll(text[i+1:len(text)-1], "''", "'")
	if prefix == "E" {
		if u, err := strconv.Unquote(`"` + strings.ReplaceAll(strings.ReplaceAll(body, `"`, `\"`), `\'`, `'`) + `"`); err == nil {
			return u
		}
	}
	return body
}

// unquoteIdent removes one level of identifier quoting: `a“b`, "a""b",
// [a]]b].
func unquoteIdent(text string) string {
	if len(text) < 2 {
		return text
	}
	switch text[0] {
	case '`':
		return strings.ReplaceAll(text[1:len(text)-1], "``", "`")
	case '"':
		return strings.ReplaceAll(text[1:len(text)-1], `""`, `"`)
	case '[':
		return strings.ReplaceAll(text[1:len(text)-1], "]]", "]")
	}
	return text
}

var simpleIdent = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_$]*$`)

// cleanModel fixes what would stop the generated DBML from compiling:
// records for unknown tables or columns, refs to unknown tables/columns,
// duplicate refs (same endpoints), duplicate ref names, and index columns
// that do not exist.
func cleanModel(db *model.Database) {
	type tinfo struct {
		t    *model.Table
		cols map[string]bool
	}
	tables := map[string]*tinfo{}
	for _, t := range db.Tables {
		k := tkey(t.SchemaName, t.Name)
		if tables[k] != nil {
			continue
		}
		ti := &tinfo{t: t, cols: map[string]bool{}}
		for _, f := range t.Fields {
			ti.cols[f.Name] = true
		}
		tables[k] = ti
		for _, ix := range t.Indexes {
			cols := ix.Columns[:0]
			for _, c := range ix.Columns {
				if c.Type == "expression" || ti.cols[c.Value] {
					cols = append(cols, c)
				}
			}
			ix.Columns = cols
		}
		ixs := t.Indexes[:0]
		for _, ix := range t.Indexes {
			if len(ix.Columns) > 0 {
				ixs = append(ixs, ix)
			}
		}
		t.Indexes = ixs
	}
	recs := db.Records[:0]
	for _, r := range db.Records {
		ti := tables[tkey(r.SchemaName, r.TableName)]
		if ti == nil {
			continue
		}
		keep := make([]int, 0, len(r.Columns))
		for i, c := range r.Columns {
			if ti.cols[c] {
				keep = append(keep, i)
			}
		}
		if len(keep) == 0 {
			continue
		}
		if len(keep) != len(r.Columns) {
			cols := make([]string, len(keep))
			for j, i := range keep {
				cols[j] = r.Columns[i]
			}
			for ri, row := range r.Values {
				nr := make([]*model.RecordValue, len(keep))
				for j, i := range keep {
					if i < len(row) {
						nr[j] = row[i]
					}
				}
				r.Values[ri] = nr
			}
			r.Columns = cols
		}
		recs = append(recs, r)
	}
	db.Records = recs
	epKey := func(e model.RefEndpoint) (string, bool) {
		ti := tables[tkey(e.SchemaName, e.TableName)]
		if ti == nil {
			return "", false
		}
		fs := append([]string(nil), e.FieldNames...)
		for _, f := range fs {
			if !ti.cols[f] {
				return "", false
			}
		}
		sort.Strings(fs)
		return tkey(e.SchemaName, e.TableName) + "\x00" + strings.Join(fs, "\x00"), true
	}
	pairs := map[string]bool{}
	names := map[string]int{}
	refs := db.Refs[:0]
	for _, r := range db.Refs {
		a, ok1 := epKey(r.Endpoints[0])
		b, ok2 := epKey(r.Endpoints[1])
		if !ok1 || !ok2 || a == b || len(r.Endpoints[0].FieldNames) != len(r.Endpoints[1].FieldNames) {
			continue
		}
		if b < a {
			a, b = b, a
		}
		if pairs[a+"|"+b] {
			continue
		}
		pairs[a+"|"+b] = true
		if r.Name != nil {
			k := str(r.SchemaName) + "\x00" + *r.Name
			names[k]++
			if n := names[k]; n > 1 {
				nn := *r.Name + "_" + strconv.Itoa(n)
				r.Name = &nn
			}
		}
		refs = append(refs, r)
	}
	db.Refs = refs
}

// tkey identifies a table (nil schema = default).
func tkey(schema *string, name string) string {
	if schema == nil {
		return "\x00" + name
	}
	return *schema + "\x00" + name
}

func str(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// setGenerated marks f as a generated / computed column (clean mode). The
// expression is normalised with the other expressions by normalizeModel.
func setGenerated(f *model.Column, expr string, stored bool) {
	if f.Metadata == nil {
		f.Metadata = map[string]string{}
	}
	f.Metadata[model.MetaGenerated] = strings.TrimSpace(expr)
	if stored {
		f.Metadata[model.MetaGeneratedStored] = "true"
	}
	f.DBDefault = nil
}

// unqualifyDefaultSchema removes the dialect's default schema (SQL
// Server's dbo) from names, so the model uses DBML's default schema and an
// export to SQL Server does not try to CREATE SCHEMA dbo.
func unqualifyDefaultSchema(db *model.Database, def string) {
	un := func(p **string) {
		if *p != nil && strings.EqualFold(**p, def) {
			*p = nil
		}
	}
	for _, t := range db.Tables {
		un(&t.SchemaName)
		for _, f := range t.Fields {
			un(&f.Type.SchemaName)
		}
	}
	for _, e := range db.Enums {
		un(&e.SchemaName)
	}
	for _, r := range db.Refs {
		un(&r.SchemaName)
		un(&r.Endpoints[0].SchemaName)
		un(&r.Endpoints[1].SchemaName)
	}
	for _, r := range db.Records {
		un(&r.SchemaName)
	}
}
