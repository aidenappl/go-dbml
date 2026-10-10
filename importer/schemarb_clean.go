package importer

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/aidenappl/go-dbml/model"
)

// Clean-mode schema.rb importer. Unlike the reference's PEG grammar it
// keeps column options (null, default, limit, precision, scale, comment),
// table comments, primary keys (implicit id typed bigint for Rails >= 5.1
// or integer before, id: :uuid / id: false / primary_key:), t.index and
// add_index, t.timestamps, and add_foreign_key with any option; unknown
// lines are skipped instead of dropping the table. Foreign keys become
// many-to-one refs (the referenced side required when the FK column is
// NOT NULL), plus the reference's implicit <singular table>_id refs.

// rbVal is one Ruby literal argument.
type rbVal struct {
	kind string // str, sym, num, bool, nil, arr, hash, lambda, raw
	s    string
	arr  []rbVal
	hash map[string]rbVal
}

// rbArgs is a parsed argument list: positional values and key: value
// options (keys without the colon).
type rbArgs struct {
	pos  []rbVal
	opts map[string]rbVal
}

func (a rbArgs) str(i int) (string, bool) {
	if i < len(a.pos) && (a.pos[i].kind == "str" || a.pos[i].kind == "sym") {
		return a.pos[i].s, true
	}
	return "", false
}

// parseRubyArgs parses `"a", :b, key: "v", k2: [..], k3: -> { "x" }, h: {a: 1}`
// in one linear pass with an explicit stack (no recursion, no re-scan):
// nesting depth costs nothing but a frame.
func parseRubyArgs(s string) rbArgs {
	type frame struct {
		kind    byte // 0 top / hash, '[' array, 'l' lambda
		hash    bool
		pos     []rbVal
		opts    map[string]rbVal
		key     string
		haveKey bool
	}
	newFrame := func(kind byte) *frame { return &frame{kind: kind} }
	stack := []*frame{newFrame(0)}
	add := func(v rbVal) {
		f := stack[len(stack)-1]
		if f.haveKey {
			if f.opts == nil {
				f.opts = map[string]rbVal{}
			}
			f.opts[f.key] = v
			f.haveKey = false
			return
		}
		f.pos = append(f.pos, v)
	}
	closeTop := func() rbVal {
		f := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch {
		case f.kind == '[':
			return rbVal{kind: "arr", arr: f.pos}
		case f.kind == 'l':
			if len(f.pos) > 0 && f.pos[0].kind == "str" {
				return rbVal{kind: "lambda", s: f.pos[0].s}
			}
			return rbVal{kind: "raw"}
		}
		return rbVal{kind: "hash", hash: f.opts}
	}
	lambdaNext := false
	n := len(s)
	for i := 0; i < n; {
		c := s[i]
		switch {
		case c == ' ' || c == '\t' || c == '\r' || c == '\n' || c == ',':
			i++
		case c == '"' || c == '\'':
			j := i + 1
			var b strings.Builder
			for j < n && s[j] != c {
				if s[j] == '\\' && j+1 < n {
					j++
					switch s[j] {
					case 'n':
						b.WriteByte('\n')
					case 't':
						b.WriteByte('\t')
					default:
						b.WriteByte(s[j])
					}
					j++
					continue
				}
				b.WriteByte(s[j])
				j++
			}
			i = min(j+1, n) // an unterminated string ends the input
			// "str" => value (old hash syntax with string keys)
			if k := skipSpaces(s, i); strings.HasPrefix(s[k:], "=>") {
				f := stack[len(stack)-1]
				f.key, f.haveKey = b.String(), true
				i = k + 2
				continue
			}
			add(rbVal{kind: "str", s: b.String()})
		case c == ':' && i+1 < n && s[i+1] != ':':
			j := i + 1
			for j < n && isRbChar(s[j]) {
				j++
			}
			name := s[i+1 : j]
			if k := skipSpaces(s, j); strings.HasPrefix(s[k:], "=>") || strings.HasPrefix(s[k:], "->") {
				f := stack[len(stack)-1]
				f.key, f.haveKey = name, true
				i = k + 2
				continue
			}
			add(rbVal{kind: "sym", s: name})
			i = j
		case c == '[':
			stack = append(stack, newFrame('['))
			i++
		case c == '{':
			k := byte(0)
			if lambdaNext {
				k = 'l'
			}
			lambdaNext = false
			stack = append(stack, newFrame(k))
			i++
		case c == ']' || c == '}':
			if len(stack) > 1 {
				add2 := closeTop()
				add(add2)
			}
			i++
		case c == '-' && i+1 < n && s[i+1] == '>':
			lambdaNext = true
			i += 2
		default:
			// key: value
			j := i
			for j < n && (s[j] == '_' || s[j] >= '0' && s[j] <= '9' || s[j]|0x20 >= 'a' && s[j]|0x20 <= 'z') {
				j++
			}
			if j > i && j < n && s[j] == ':' && (j+1 >= n || s[j+1] != ':') && !(s[i] >= '0' && s[i] <= '9') {
				f := stack[len(stack)-1]
				f.key, f.haveKey = s[i:j], true
				i = j + 1
				continue
			}
			// a bare word / number / call: up to a delimiter, parentheses balanced
			j = i
			depth := 0
			for j < n {
				d := s[j]
				if depth == 0 && (d == ',' || d == ']' || d == '}' || d == '[' || d == '{') {
					break
				}
				if d == '(' {
					depth++
				} else if d == ')' && depth > 0 {
					depth--
				} else if d == '"' || d == '\'' {
					for j++; j < n && s[j] != d; j++ {
						if s[j] == '\\' {
							j++
						}
					}
				}
				j++
			}
			if j == i {
				j = i + 1
			}
			t := strings.TrimSpace(s[i:min(j, n)])
			i = j
			switch {
			case t == "":
			case t == "true" || t == "false":
				add(rbVal{kind: "bool", s: t})
			case t == "nil":
				add(rbVal{kind: "nil"})
			case rbNumber.MatchString(t):
				add(rbVal{kind: "num", s: strings.ReplaceAll(t, "_", "")})
			default:
				add(rbVal{kind: "raw", s: t})
			}
		}
	}
	for len(stack) > 1 {
		v := closeTop()
		add(v)
	}
	if stack[0].opts == nil {
		stack[0].opts = map[string]rbVal{}
	}
	return rbArgs{pos: stack[0].pos, opts: stack[0].opts}
}

func skipSpaces(s string, i int) int {
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	return i
}

var (
	rbNumber    = regexp.MustCompile(`^-?[0-9][0-9_]*(\.[0-9_]+)?$`)
	rbKey       = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*):(?:[^:]|$)`)
	rbRocket    = regexp.MustCompile(`^:([A-Za-z_][A-Za-z0-9_]*)\s*[=-]>\s*`)
	rbCall      = regexp.MustCompile(`^(\w+)\.(\w+)\s*(.*?)\s*$`)
	rbDefine    = regexp.MustCompile(`ActiveRecord::Schema(?:\[(\d+)\.(\d+)\])?\.define\s*\(\s*version:\s*([0-9_]+)`)
	rbTopCall   = regexp.MustCompile(`^(create_table|create_enum|add_index|add_foreign_key|add_check_constraint|change_column_comment|change_table_comment)\s*\(?\s*(.*?)\)?\s*(do\s*\|\s*\w+\s*\|)?\s*$`)
	rbRailsDate = 20170427 // Rails 5.1 (bigint primary keys)
)

func importSchemaRBClean(src string) (*rawDB, error) {
	raw := &rawDB{}
	my := rbIsMySQL(src)
	idType := "bigint"
	if m := rbDefine.FindStringSubmatch(src); m != nil {
		if m[1] != "" {
			major, _ := strconv.Atoi(m[1])
			minor, _ := strconv.Atoi(m[2])
			if major < 5 || (major == 5 && minor < 1) {
				idType = "integer"
			}
		} else if v := strings.ReplaceAll(m[3], "_", ""); len(v) >= 8 {
			if d, err := strconv.Atoi(v[:8]); err == nil && d < rbRailsDate && d > 19000101 {
				idType = "integer"
			}
		}
	}
	type fk struct {
		from, to string
		args     rbArgs
		col      string // set for t.references ... foreign_key:
	}
	var fks []fk
	var cur *model.Table
	var curTable rbCleanTable
	var curVar string
	lines := strings.Split(src, "\n")
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		if cur != nil {
			if l == "end" {
				raw.addTable(cur)
				cur = nil
				continue
			}
			m := splitRbCall(l)
			if m == nil || m[1] != curVar {
				continue
			}
			curTable.column(m[2], parseRubyArgs(m[3]), func(from, col, to string, o map[string]rbVal) {
				fks = append(fks, fk{from: from, to: to, args: rbArgs{opts: o}, col: col})
			})
			continue
		}
		m := rbTopCall.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		args := parseRubyArgs(m[2])
		switch m[1] {
		case "create_table":
			name, ok := args.str(0)
			if !ok {
				continue
			}
			cur = &model.Table{Name: name}
			curTable = rbCleanTable{t: cur, my: my}
			curVar = "t"
			if mv := regexp.MustCompile(`\|\s*(\w+)\s*\|`).FindStringSubmatch(m[3]); mv != nil {
				curVar = mv[1]
			}
			if c, ok := args.opts["comment"]; ok && c.kind == "str" {
				cur.Note = &model.NoteValue{Value: c.s}
			}
			idv, hasID := args.opts["id"]
			pkName := "id"
			if pk, ok := args.opts["primary_key"]; ok && (pk.kind == "str" || pk.kind == "sym") {
				pkName = pk.s
			}
			if !(hasID && idv.kind == "bool" && idv.s == "false") {
				typ := idType
				if hasID && (idv.kind == "sym" || idv.kind == "str") {
					typ = idv.s
				}
				col := &model.Column{Name: pkName, Type: model.ColumnType{TypeName: rbConcreteType(typ, rbArgs{opts: map[string]rbVal{}}, my)}, PK: boolp(true)}
				if typ == "bigint" || typ == "integer" || typ == "serial" || typ == "bigserial" || typ == "primary_key" {
					col.Increment = boolp(true)
				}
				if d, ok := args.opts["default"]; ok {
					col.DBDefault = rbDefault(d)
				}
				cur.Fields = append(cur.Fields, col)
			}
			if m[3] == "" { // create_table "x" without a block
				raw.addTable(cur)
				cur = nil
			}
		case "create_enum":
			name, ok := args.str(0)
			if ok && len(args.pos) > 1 && args.pos[1].kind == "arr" {
				e := &model.Enum{Name: name}
				for _, v := range args.pos[1].arr {
					if v.kind == "str" || v.kind == "sym" {
						e.Values = append(e.Values, &model.EnumValue{Name: v.s})
					}
				}
				raw.enums = append(raw.enums, e)
			}
		case "add_index":
			table, ok := args.str(0)
			if t := raw.findTable(nil, table); ok && t != nil {
				if ix := rbIndex(args, 1); ix != nil {
					t.Indexes = append(t.Indexes, ix)
				}
			}
		case "add_foreign_key":
			from, ok1 := args.str(0)
			to, ok2 := args.str(1)
			if ok1 && ok2 {
				fks = append(fks, fk{from: from, to: to, args: args})
			}
		case "add_check_constraint":
			table, ok := args.str(0)
			expr, ok2 := args.str(1)
			if t := raw.findTable(nil, table); ok && ok2 && t != nil {
				ch := &model.Check{Expression: expr}
				if n, ok := args.opts["name"]; ok && n.kind == "str" {
					ch.Name = strp(n.s)
				}
				t.Checks = append(t.Checks, ch)
			}
		case "change_table_comment":
			table, ok := args.str(0)
			if t := raw.findTable(nil, table); ok && t != nil {
				if v, ok := args.str(1); ok {
					t.Note = &model.NoteValue{Value: v}
				}
			}
		}
	}
	if cur != nil {
		raw.addTable(cur)
	}
	// explicit foreign keys, then the reference's implicit <singular>_id refs
	seen := map[string]bool{}
	addRef := func(from, col, to, pk string, name *string, del, upd *string) {
		ft, tt := raw.findTable(nil, from), raw.findTable(nil, to)
		if ft == nil || tt == nil || findField(ft, col) == nil || findField(tt, pk) == nil {
			return
		}
		k := from + "." + col + ">" + to + "." + pk
		if seen[k] || (from == to && col == pk) {
			return
		}
		seen[k] = true
		// a foreign key column typed unlike the key it references (a
		// bigint t.references to a uuid id) cannot be created: follow the key
		if fc, pc := findField(ft, col), findField(tt, pk); fc.Type.TypeName != pc.Type.TypeName && pc.Type.TypeName != "" {
			fc.Type = model.ColumnType{TypeName: pc.Type.TypeName}
		}
		r := &model.Ref{Name: name, OnDelete: del, OnUpdate: upd}
		r.Endpoints[0] = model.RefEndpoint{TableName: from, FieldNames: []string{col}, Relation: "0..*"}
		r.Endpoints[1] = model.RefEndpoint{TableName: to, FieldNames: []string{pk}, Relation: "0..1"}
		raw.refs = append(raw.refs, r)
	}
	for _, f := range fks {
		col := singularize(f.to) + "_id"
		if f.col != "" {
			col = f.col
		}
		if c, ok := f.args.opts["column"]; ok && (c.kind == "str" || c.kind == "sym") {
			col = c.s
		}
		pk := "id"
		if c, ok := f.args.opts["primary_key"]; ok && (c.kind == "str" || c.kind == "sym") {
			pk = c.s
		}
		var name, del, upd *string
		if n, ok := f.args.opts["name"]; ok && n.kind == "str" {
			name = strp(n.s)
		}
		if v, ok := f.args.opts["on_delete"]; ok && v.kind == "sym" {
			del = strp(rbAction(v.s))
		}
		if v, ok := f.args.opts["on_update"]; ok && v.kind == "sym" {
			upd = strp(rbAction(v.s))
		}
		addRef(f.from, col, f.to, pk, name, del, upd)
	}
	byField := map[string]string{}
	for _, t := range raw.tables {
		k := singularize(t.Name) + "_id"
		if _, ok := byField[k]; !ok {
			byField[k] = t.Name
		}
	}
	for _, t := range raw.tables {
		for _, f := range t.Fields {
			if to, ok := byField[f.Name]; ok {
				addRef(t.Name, f.Name, to, "id", nil, nil, nil)
			}
		}
	}
	return raw, nil
}

// rbCleanTable is the table of the current create_table block.
type rbCleanTable struct {
	t  *model.Table
	my bool // the schema targets MySQL (charset:/collation:/ENGINE options)
}

// column handles one `t.<method> ...` line of a create_table block; fk
// records a foreign key declared with references ... foreign_key:.
func (rt rbCleanTable) column(method string, a rbArgs, fk func(from, col, to string, opts map[string]rbVal)) {
	t := rt.t
	switch method {
	case "index":
		if ix := rbIndex(a, 0); ix != nil {
			t.Indexes = append(t.Indexes, ix)
		}
		return
	case "timestamps":
		for _, n := range []string{"created_at", "updated_at"} {
			col := &model.Column{Name: n, Type: model.ColumnType{TypeName: rbConcreteType("datetime", a, rt.my)}, NotNull: boolp(true)}
			if v, ok := a.opts["null"]; ok && v.kind == "bool" && v.s == "true" {
				col.NotNull = nil
			}
			t.Fields = append(t.Fields, col)
		}
		return
	case "check_constraint":
		if e, ok := a.str(0); ok {
			ch := &model.Check{Expression: e}
			if n, ok := a.opts["name"]; ok && n.kind == "str" {
				ch.Name = strp(n.s)
			}
			t.Checks = append(t.Checks, ch)
		}
		return
	case "references", "belongs_to":
		for _, v := range a.pos {
			if v.kind != "str" && v.kind != "sym" {
				continue
			}
			typ := "bigint"
			if ty, ok := a.opts["type"]; ok && (ty.kind == "sym" || ty.kind == "str") {
				typ = ty.s
			}
			col := &model.Column{Name: v.s + "_id", Type: model.ColumnType{TypeName: rbConcreteType(typ, rbArgs{opts: map[string]rbVal{}}, rt.my)}}
			if n, ok := a.opts["null"]; ok && n.kind == "bool" && n.s == "false" {
				col.NotNull = boolp(true)
			}
			t.Fields = append(t.Fields, col)
			poly := false
			if p, ok := a.opts["polymorphic"]; ok && p.kind == "bool" && p.s == "true" {
				poly = true
				t.Fields = append(t.Fields, &model.Column{Name: v.s + "_type", Type: model.ColumnType{TypeName: rbConcreteType("string", rbArgs{opts: map[string]rbVal{}}, rt.my)}})
			}
			// Rails indexes a reference unless index: false
			if ix, ok := a.opts["index"]; !(ok && ix.kind == "bool" && ix.s == "false") {
				cols := []*model.IndexColumn{{Type: "column", Value: col.Name}}
				if poly {
					cols = []*model.IndexColumn{{Type: "column", Value: v.s + "_type"}, {Type: "column", Value: col.Name}}
				}
				idx := &model.Index{Columns: cols, Name: strp("index_" + t.Name + "_on_" + col.Name)}
				if ix.kind == "hash" {
					if u, ok := ix.hash["unique"]; ok && u.kind == "bool" && u.s == "true" {
						idx.Unique = boolp(true)
					}
					if n, ok := ix.hash["name"]; ok && n.kind == "str" {
						idx.Name = strp(n.s)
					}
				}
				t.Indexes = append(t.Indexes, idx)
			}
			if f, ok := a.opts["foreign_key"]; ok && !poly {
				switch {
				case f.kind == "bool" && f.s == "true":
					fk(t.Name, col.Name, pluralize(v.s), map[string]rbVal{})
				case f.kind == "hash":
					to := pluralize(v.s)
					if tt, ok := f.hash["to_table"]; ok && (tt.kind == "sym" || tt.kind == "str") {
						to = tt.s
					}
					fk(t.Name, col.Name, to, f.hash)
				}
			}
		}
		return
	case "virtual":
		for _, v := range a.pos {
			if v.kind != "str" && v.kind != "sym" {
				continue
			}
			typ := "string"
			if ty, ok := a.opts["type"]; ok && (ty.kind == "sym" || ty.kind == "str") {
				typ = ty.s
			}
			col := &model.Column{Name: v.s, Type: model.ColumnType{TypeName: rbConcreteType(typ, a, rt.my)}}
			if as, ok := a.opts["as"]; ok && as.kind == "str" {
				st := a.opts["stored"]
				setGenerated(col, as.s, st.kind == "bool" && st.s == "true")
			}
			t.Fields = append(t.Fields, col)
		}
		return
	}
	for _, v := range a.pos {
		if v.kind != "str" && v.kind != "sym" {
			continue
		}
		typ := method
		if method == "column" && len(a.pos) > 1 {
			if tv := a.pos[1]; tv.kind == "sym" || tv.kind == "str" {
				typ = tv.s
			}
		}
		if method == "enum" {
			if et, ok := a.opts["enum_type"]; ok && (et.kind == "str" || et.kind == "sym") {
				typ = et.s
			}
		}
		col := &model.Column{Name: v.s, Type: model.ColumnType{TypeName: rbConcreteType(typ, a, rt.my)}}
		if n, ok := a.opts["null"]; ok && n.kind == "bool" && n.s == "false" {
			col.NotNull = boolp(true)
		}
		if d, ok := a.opts["default"]; ok {
			col.DBDefault = rbDefault(d)
		}
		if c, ok := a.opts["comment"]; ok && c.kind == "str" {
			col.Note = &model.NoteValue{Value: c.s}
		}
		if p, ok := a.opts["primary_key"]; ok && p.kind == "bool" && p.s == "true" {
			col.PK = boolp(true)
		}
		// an explicit column named like the implicit primary key replaces it
		if len(t.Fields) > 0 && t.Fields[0].Name == col.Name && truthy(t.Fields[0].PK) {
			col.PK = boolp(true)
			t.Fields[0] = col
			continue
		}
		t.Fields = append(t.Fields, col)
		if method == "column" {
			break // t.column "name", :type
		}
	}
}

// rbConcreteType maps a Rails abstract column type (with limit /
// precision / scale / array / unsigned / size) to the SQL type the Rails
// adapter creates: PostgreSQL, or MySQL when the schema shows MySQL table
// options. Unknown types (database-specific ones such as inet, citext,
// hstore, or an enum type name) are kept as written.
func rbConcreteType(typ string, a rbArgs, my bool) string {
	num := func(k string) (string, bool) {
		v, ok := a.opts[k]
		return v.s, ok && v.kind == "num"
	}
	limit, hasL := num("limit")
	prec, hasP := num("precision")
	scale, hasS := num("scale")
	dims := func() string {
		switch {
		case hasP && hasS:
			return "(" + prec + "," + scale + ")"
		case hasP:
			return "(" + prec + ")"
		}
		return ""
	}
	var out string
	switch typ {
	case "string":
		out = "varchar"
		if hasL {
			out += "(" + limit + ")"
		} else if my {
			out += "(255)"
		}
	case "text":
		out = "text"
		if sz, ok := a.opts["size"]; my && ok && sz.kind == "sym" {
			switch sz.s {
			case "tiny":
				out = "tinytext"
			case "medium":
				out = "mediumtext"
			case "long":
				out = "longtext"
			}
		}
	case "integer":
		n, _ := strconv.Atoi(limit)
		switch {
		case hasL && n <= 1:
			out = map[bool]string{true: "tinyint", false: "smallint"}[my]
		case hasL && n == 2:
			out = "smallint"
		case hasL && n == 3 && my:
			out = "mediumint"
		case hasL && n > 4:
			out = "bigint"
		default:
			out = map[bool]string{true: "int", false: "integer"}[my]
		}
	case "bigint", "primary_key":
		out = "bigint"
	case "float":
		out = map[bool]string{true: "float", false: "double precision"}[my]
	case "decimal", "numeric":
		out = typ + dims()
	case "datetime":
		out = map[bool]string{true: "datetime", false: "timestamp"}[my] + dims()
	case "timestamp", "time":
		out = typ + dims()
	case "timestamptz":
		out = map[bool]string{true: "timestamp", false: "timestamptz"}[my] + dims()
	case "date":
		out = "date"
	case "binary", "blob":
		out = map[bool]string{true: "blob", false: "bytea"}[my]
	case "boolean":
		out = map[bool]string{true: "tinyint(1)", false: "boolean"}[my]
	case "json":
		out = "json"
	case "jsonb":
		out = map[bool]string{true: "json", false: "jsonb"}[my]
	case "uuid":
		out = map[bool]string{true: "char(36)", false: "uuid"}[my]
	default:
		out = typ
		if hasP || hasS {
			out += dims()
		} else if hasL {
			out += "(" + limit + ")"
		}
	}
	if u, ok := a.opts["unsigned"]; my && ok && u.kind == "bool" && u.s == "true" {
		out += " unsigned"
	}
	if arr, ok := a.opts["array"]; ok && arr.kind == "bool" && arr.s == "true" {
		out += "[]"
	}
	return out
}

func rbDefault(v rbVal) *model.Default {
	switch v.kind {
	case "arr", "hash":
		if len(v.arr) == 0 && len(v.hash) == 0 {
			return &model.Default{Type: "string", Value: "{}"} // [] / {} (PostgreSQL array / json)
		}
		return nil
	case "str", "sym":
		return &model.Default{Type: "string", Value: v.s}
	case "num":
		return &model.Default{Type: "number", Value: v.s}
	case "bool":
		return &model.Default{Type: "boolean", Value: v.s}
	case "lambda":
		return &model.Default{Type: "expression", Value: v.s}
	}
	return nil
}

// rbIndex builds an index from `["a", "b"], name: "x", unique: true`
// starting at positional argument i.
func rbIndex(a rbArgs, i int) *model.Index {
	if i >= len(a.pos) {
		return nil
	}
	ix := &model.Index{}
	switch v := a.pos[i]; v.kind {
	case "arr":
		for _, c := range v.arr {
			if c.kind == "str" || c.kind == "sym" {
				ix.Columns = append(ix.Columns, &model.IndexColumn{Type: "column", Value: c.s})
			}
		}
	case "str", "sym":
		ix.Columns = append(ix.Columns, &model.IndexColumn{Type: "column", Value: v.s})
	}
	if len(ix.Columns) == 0 {
		return nil
	}
	if n, ok := a.opts["name"]; ok && n.kind == "str" {
		ix.Name = strp(n.s)
	}
	if u, ok := a.opts["unique"]; ok && u.kind == "bool" && u.s == "true" {
		ix.Unique = boolp(true)
	}
	return ix
}

// rbAction maps a Rails on_delete / on_update symbol to SQL.
func rbAction(s string) string {
	switch s {
	case "nullify":
		return "set null"
	case "restrict":
		return "restrict"
	}
	return strings.ReplaceAll(s, "_", " ")
}

// splitRbCall splits `recv.method rest` (rbCall without a regexp: lines can
// be megabytes long and the regexp engine's per-byte cost dominated).
func splitRbCall(l string) []string {
	word := func(i int) int {
		for i < len(l) && (l[i] == '_' || l[i] >= '0' && l[i] <= '9' || l[i]|0x20 >= 'a' && l[i]|0x20 <= 'z') {
			i++
		}
		return i
	}
	i := word(0)
	if i == 0 || i >= len(l) || l[i] != '.' {
		return nil
	}
	j := word(i + 1)
	if j == i+1 {
		return nil
	}
	if j < len(l) && l[j] != ' ' && l[j] != '\t' && l[j] != '(' {
		return nil
	}
	return []string{l, l[:i], l[i+1 : j], strings.TrimSpace(l[j:])}
}

// rbIsMySQL reports whether the schema shows MySQL table options
// (charset: "…", collation: "utf8…", options: "ENGINE=…", unsigned: true).
// A plain scan: a regexp over a multi-megabyte schema was the importer's
// largest cost.
func rbIsMySQL(src string) bool {
	for _, k := range []struct{ key, val string }{{"charset:", `"`}, {"collation:", `"utf8`}, {"options:", `"ENGINE=`}, {"unsigned:", "true"}} {
		for i := 0; ; {
			j := strings.Index(src[i:], k.key)
			if j < 0 {
				break
			}
			i += j + len(k.key)
			if strings.HasPrefix(src[skipSpaces(src, i):], k.val) {
				return true
			}
		}
	}
	return false
}
