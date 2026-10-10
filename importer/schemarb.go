// Portions of this file are derived from holistics/dbml
// packages/dbml-core/src/parse/deprecated/schemarb/parser.pegjs
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: the PEG grammar and its actions are hand-translated to a Go
// backtracking matcher with the same rule order and side effects.

package importer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aidenappl/go-dbml/model"
)

// The reference schema.rb importer reads only: create_table "name" ... do
// blocks whose every body line is `t.<type> "col"`, `t.index ...`,
// `t.references :x` or `t.check_constraint "expr"`; add_foreign_key; and
// add_check_constraint. Column options (null:, default:, limit: ...) are
// ignored, a missing "id" column is added as varchar, and every column
// named <singular table>_id becomes a one-to-one ref to that table's id.

type rbTable struct {
	name   string
	fields []*model.Column
	checks []*model.Check
}

type rbRef struct {
	name               string
	endpoints          [2]rbEndpoint
	onDelete, onUpdate *string
}

type rbEndpoint struct {
	table  string
	fields []string // nil = unset
}

type rbParser struct {
	s      string
	tables []*rbTable
	refs   []*rbRef
	err    error
}

func importSchemaRB(src string, cfg *config) (*rawDB, error) {
	if cfg.clean {
		return importSchemaRBClean(src)
	}
	p := &rbParser{s: src}
	pos := 0
	for pos < len(p.s) {
		n, ok := p.lineRule(pos)
		if p.err != nil {
			return nil, p.err
		}
		if !ok || n == pos {
			line := 1 + strings.Count(p.s[:pos], "\n")
			return nil, fmt.Errorf("import: %w: schema.rb syntax error at line %d", ErrUnsupported, line)
		}
		pos = n
	}
	p.implicitRefs()
	raw := &rawDB{}
	for _, t := range p.tables {
		mt := &model.Table{Name: t.name, Fields: t.fields, Checks: t.checks}
		raw.tables = append(raw.tables, mt)
	}
	for _, r := range p.refs {
		mr := &model.Ref{Name: strp(r.name), OnDelete: r.onDelete, OnUpdate: r.onUpdate}
		for i, ep := range r.endpoints {
			mr.Endpoints[i] = model.RefEndpoint{TableName: ep.table, FieldNames: ep.fields, Relation: "1"}
		}
		raw.refs = append(raw.refs, mr)
	}
	return raw, nil
}

// ---------------------------------------------------------------- terminals

func (p *rbParser) at(pos int) byte {
	if pos < len(p.s) {
		return p.s[pos]
	}
	return 0
}

func isRbChar(c byte) bool {
	return c == '_' || c == '.' || isDigit(c) || (c|0x20 >= 'a' && c|0x20 <= 'z')
}

func (p *rbParser) lit(pos int, s string, fold bool) (int, bool) {
	if pos+len(s) > len(p.s) {
		return pos, false
	}
	got := p.s[pos : pos+len(s)]
	if got == s || (fold && strings.EqualFold(got, s)) {
		return pos + len(s), true
	}
	return pos, false
}

func (p *rbParser) sps(pos int) int {
	for p.at(pos) == ' ' {
		pos++
	}
	return pos
}

func (p *rbParser) whitespaces(pos int) int {
	for {
		switch p.at(pos) {
		case ' ', '\t', '\r', '\n':
			pos++
		default:
			return pos
		}
	}
}

// whateters: [^\t\r\n]*
func (p *rbParser) whateters(pos int) int {
	for pos < len(p.s) {
		switch p.s[pos] {
		case '\t', '\r', '\n':
			return pos
		}
		pos++
	}
	return pos
}

// endline: sp* newline
func (p *rbParser) endline(pos int) (int, bool) {
	q := p.sps(pos)
	if n, ok := p.lit(q, "\r\n", false); ok {
		return n, true
	}
	if p.at(q) == '\n' {
		return q + 1, true
	}
	return pos, false
}

// name: "..." | '...' (no newline, no escapes)
func (p *rbParser) name(pos int) (string, int, bool) {
	q := p.at(pos)
	if q != '"' && q != '\'' {
		return "", pos, false
	}
	i := pos + 1
	for i < len(p.s) && p.s[i] != q && p.s[i] != '\n' {
		i++
	}
	if p.at(i) != q {
		return "", pos, false
	}
	return p.s[pos+1 : i], i + 1, true
}

// variable: [a-z0-9_.]i+
func (p *rbParser) variable(pos int) (string, int, bool) {
	i := pos
	for i < len(p.s) && isRbChar(p.s[i]) {
		i++
	}
	return p.s[pos:i], i, i > pos
}

// symbol: ":" character*
func (p *rbParser) symbol(pos int) (string, int, bool) {
	if p.at(pos) != ':' {
		return "", pos, false
	}
	i := pos + 1
	for i < len(p.s) && isRbChar(p.s[i]) {
		i++
	}
	return p.s[pos+1 : i], i, true
}

// ---------------------------------------------------------------- rules

func (p *rbParser) lineRule(pos int) (int, bool) {
	if n, ok := p.rule(p.whitespaces(pos)); ok || p.err != nil {
		return n, ok
	}
	// comment_line: whitespace* "#" whateters endline?
	if q := p.whitespaces(pos); p.at(q) == '#' {
		q = p.whateters(q + 1)
		if n, ok := p.endline(q); ok {
			q = n
		}
		return q, true
	}
	// end_line
	if n, ok := p.endLine(pos); ok {
		return n, true
	}
	// __ = (comment / whitespace)+  (comment = "//" [^\n]?)
	q := pos
	for {
		if n, ok := p.lit(q, "//", false); ok {
			if n < len(p.s) && p.s[n] != '\n' {
				n++
			}
			q = n
			continue
		}
		if c := p.at(q); c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			q++
			continue
		}
		break
	}
	return q, q > pos
}

// end_line: whitespace* "end" endline?
func (p *rbParser) endLine(pos int) (int, bool) {
	q := p.whitespaces(pos)
	n, ok := p.lit(q, "end", false)
	if !ok {
		return pos, false
	}
	if m, ok := p.endline(n); ok {
		n = m
	}
	return n, true
}

func (p *rbParser) rule(pos int) (int, bool) {
	if n, ok := p.createTable(pos); ok || p.err != nil {
		return n, ok
	}
	if n, ok := p.addForeignKey(pos); ok || p.err != nil {
		return n, ok
	}
	if n, ok := p.addCheckConstraint(pos); ok || p.err != nil {
		return n, ok
	}
	// other_class_prop: variable whateters endline?
	if _, n, ok := p.variable(pos); ok {
		n = p.whateters(n)
		if m, ok := p.endline(n); ok {
			n = m
		}
		return n, true
	}
	return pos, false
}

type rbField struct {
	kind  string // field | reference | check | index
	field *model.Column
	ref   string
	check *model.Check
}

func (p *rbParser) createTable(pos int) (int, bool) {
	q, ok := p.lit(pos, "create_table", true)
	if !ok {
		return pos, false
	}
	q = p.sps(q)
	tname, q, ok := p.name(q)
	if !ok {
		return pos, false
	}
	q = p.whateters(q)
	if q, ok = p.endline(q); !ok {
		return pos, false
	}
	var body []rbField
	for {
		f, n, ok := p.field(q)
		if !ok {
			break
		}
		body = append(body, f)
		q = n
	}
	if q, ok = p.endLine(q); !ok {
		return pos, false
	}
	// action: build the table, refs from t.references, pushTable, pushRefs
	t := &rbTable{name: tname}
	var references []string
	for _, f := range body {
		switch f.kind {
		case "field":
			t.fields = append(t.fields, f.field)
		case "reference":
			references = append(references, f.ref)
		case "check":
			t.checks = append(t.checks, f.check)
		}
	}
	refs := p.refsFromReferences(t, references)
	for _, x := range p.tables {
		if x.name == t.name {
			p.err = fmt.Errorf("import: %w: duplicated table name %q", ErrUnsupported, t.name)
			return pos, false
		}
	}
	hasID := false
	for _, f := range t.fields {
		if f.Name == "id" {
			hasID = true
		}
	}
	if !hasID {
		t.fields = append([]*model.Column{{Name: "id", Type: model.ColumnType{TypeName: "varchar"}}}, t.fields...)
	}
	p.tables = append(p.tables, t)
	for _, r := range refs {
		p.pushRef(r)
		if p.err != nil {
			return pos, false
		}
	}
	return q, true
}

// field: whitespace* table_field_syntax whatever* endline
func (p *rbParser) field(pos int) (rbField, int, bool) {
	q := p.whitespaces(pos)
	var f rbField
	ok := false
	if c := p.at(q); isRbChar(c) {
		// field_index_syntax: character ".index" sp+ whateters
		if n, k := p.lit(q+1, ".index", false); k && p.at(n) == ' ' {
			f, ok, q = rbField{kind: "index"}, true, p.whateters(p.sps(n))
		} else if n, k := p.lit(q+1, ".references", false); k && p.at(n) == ' ' {
			n = p.sps(n)
			if v, m, k := p.symbolVar(n); k {
				f, ok, q = rbField{kind: "reference", ref: v}, true, m
			} else if v, m, k := p.name(n); k {
				f, ok, q = rbField{kind: "reference", ref: v}, true, m
			}
		}
		if !ok {
			if n, k := p.lit(q+1, ".check_constraint", false); k && p.at(n) == ' ' {
				if expr, m, k := p.name(p.sps(n)); k {
					ch := &model.Check{Expression: expr}
					for {
						cn, mm, k := p.checkProp(m)
						if !k {
							break
						}
						ch.Name = strp(cn)
						m = mm
					}
					f, ok, q = rbField{kind: "check", check: ch}, true, m
				}
			}
		}
		if !ok && p.at(q+1) == '.' {
			// field_type_syntax: character "." character+ sp+ name
			i := q + 2
			for i < len(p.s) && isRbChar(p.s[i]) {
				i++
			}
			if i > q+2 && p.at(i) == ' ' {
				typ := p.s[q+2 : i]
				if col, m, k := p.name(p.sps(i)); k {
					f, ok, q = rbField{kind: "field", field: &model.Column{Name: col, Type: model.ColumnType{TypeName: typ}}}, true, m
				}
			}
		}
	}
	if !ok {
		return f, pos, false
	}
	q = p.whateters(q)
	n, k := p.endline(q)
	if !k {
		return f, pos, false
	}
	return f, n, true
}

// reference_value ":" variable
func (p *rbParser) symbolVar(pos int) (string, int, bool) {
	if p.at(pos) != ':' {
		return "", pos, false
	}
	return p.variable(pos + 1)
}

// check_constraint_props: "," sp* "name"i ":" sp* name
func (p *rbParser) checkProp(pos int) (string, int, bool) {
	if p.at(pos) != ',' {
		return "", pos, false
	}
	q := p.sps(pos + 1)
	q, ok := p.lit(q, "name", true)
	if !ok || p.at(q) != ':' {
		return "", pos, false
	}
	n, q, ok := p.name(p.sps(q + 1))
	if !ok {
		return "", pos, false
	}
	return n, q, true
}

func (p *rbParser) addForeignKey(pos int) (int, bool) {
	q := p.sps(pos)
	q, ok := p.lit(q, "add_foreign_key", true)
	if !ok {
		return pos, false
	}
	q = p.sps(q)
	from, q, ok := p.name(q)
	if !ok || p.at(q) != ',' {
		return pos, false
	}
	to, q, ok := p.name(p.sps(q + 1))
	if !ok {
		return pos, false
	}
	r := &rbRef{name: "fk_rails_" + from + "_" + to}
	r.endpoints[0].table, r.endpoints[1].table = from, to
	for p.at(q) == ',' {
		n := p.sps(q + 1)
		if m, k := p.lit(n, "column:", false); k {
			if v, mm, k := p.name(p.sps(m)); k {
				r.endpoints[0].fields = []string{v}
				q = mm
				continue
			}
		}
		if m, k := p.lit(n, "primary_key:", false); k {
			if v, mm, k := p.name(p.sps(m)); k {
				r.endpoints[1].fields = []string{v}
				q = mm
				continue
			}
		}
		var which string
		if m, k := p.lit(n, "on_delete", true); k && p.at(m) == ':' {
			which, n = "delete", m+1
		} else if m, k := p.lit(n, "on_update", true); k && p.at(m) == ':' {
			which, n = "update", m+1
		}
		if which != "" {
			if v, mm, k := p.symbol(p.sps(n)); k {
				act := strings.Join(strings.Split(v, "_"), " ")
				if which == "delete" {
					r.onDelete = &act
				} else {
					r.onUpdate = &act
				}
				q = mm
				continue
			}
		}
		break
	}
	if fr := p.refactorForeign(r); fr != nil {
		p.pushRef(fr)
	}
	return q, true
}

func (p *rbParser) addCheckConstraint(pos int) (int, bool) {
	q := p.sps(pos)
	q, ok := p.lit(q, "add_check_constraint", true)
	if !ok {
		return pos, false
	}
	q = p.sps(q)
	table, n, ok := p.symbol(q)
	if !ok {
		if table, n, ok = p.name(q); !ok {
			return pos, false
		}
	}
	if p.at(n) != ',' {
		return pos, false
	}
	expr, n, ok := p.name(p.sps(n + 1))
	if !ok {
		return pos, false
	}
	var cname *string
	for {
		v, m, k := p.checkProp(n)
		if !k {
			break
		}
		if cname == nil {
			cname = strp(v)
		}
		n = m
	}
	t := p.findTable(table)
	if t == nil {
		p.err = fmt.Errorf("import: %w: table %s not found", ErrUnsupported, table)
		return pos, false
	}
	t.checks = append(t.checks, &model.Check{Expression: expr, Name: cname})
	return n, true
}

func (p *rbParser) findTable(name string) *rbTable {
	for _, t := range p.tables {
		if t.name == name {
			return t
		}
	}
	return nil
}

func (p *rbParser) refactorForeign(r *rbRef) *rbRef {
	from := p.findTable(r.endpoints[0].table)
	if from == nil || p.findTable(r.endpoints[1].table) == nil {
		return nil
	}
	if r.endpoints[0].fields == nil {
		col := singularize(r.endpoints[1].table) + "_id"
		r.endpoints[0].fields = []string{col}
		found := false
		for _, f := range from.fields {
			if f.Name == col {
				found = true
			}
		}
		if !found {
			return nil
		}
	}
	if r.endpoints[1].fields == nil {
		r.endpoints[1].fields = []string{"id"}
	}
	return r
}

func (p *rbParser) refsFromReferences(t *rbTable, references []string) []*rbRef {
	var refs []*rbRef
	for _, ref := range references {
		refTable := pluralize(ref)
		to := p.findTable(refTable)
		if to == nil {
			continue
		}
		col := ref + "_id"
		has := false
		for _, f := range t.fields {
			if f.Name == col {
				has = true
			}
		}
		if !has {
			typ := "varchar"
			for _, f := range to.fields {
				if f.Name == "id" {
					typ = f.Type.TypeName
					break
				}
			}
			t.fields = append(t.fields, &model.Column{Name: col, Type: model.ColumnType{TypeName: typ}})
		}
		r := &rbRef{name: "fk_rails_" + t.name + "_" + refTable}
		r.endpoints[0] = rbEndpoint{table: t.name, fields: []string{col}}
		r.endpoints[1] = rbEndpoint{table: refTable, fields: []string{"id"}}
		refs = append(refs, r)
	}
	return refs
}

func sameEndpoint(a, b rbEndpoint) bool {
	if a.table != b.table || len(a.fields) != len(b.fields) {
		return false
	}
	x, y := append([]string{}, a.fields...), append([]string{}, b.fields...)
	sort.Strings(x)
	sort.Strings(y)
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func sameRef(a, b *rbRef) bool {
	return (sameEndpoint(a.endpoints[0], b.endpoints[0]) && sameEndpoint(a.endpoints[1], b.endpoints[1])) ||
		(sameEndpoint(a.endpoints[0], b.endpoints[1]) && sameEndpoint(a.endpoints[1], b.endpoints[0]))
}

func (p *rbParser) pushRef(r *rbRef) {
	for _, x := range p.refs {
		if sameRef(x, r) {
			p.err = fmt.Errorf("import: %w: duplicated references", ErrUnsupported)
			return
		}
	}
	p.refs = append(p.refs, r)
}

// implicitRefs mirrors implicityRef: every <singular(table)>_id column
// references that table's id.
func (p *rbParser) implicitRefs() {
	type tf struct{ name, field string }
	var names []tf
	for _, t := range p.tables {
		names = append(names, tf{t.name, singularize(t.name) + "_id"})
	}
	for _, t := range p.tables {
		for _, f := range t.fields {
			for _, n := range names {
				if n.field != f.Name {
					continue
				}
				r := &rbRef{name: "fk_rails_" + t.name + "_" + n.name}
				r.endpoints[0] = rbEndpoint{table: t.name, fields: []string{f.Name}}
				r.endpoints[1] = rbEndpoint{table: n.name, fields: []string{"id"}}
				dup := false
				for _, x := range p.refs {
					if sameRef(x, r) {
						dup = true
						break
					}
				}
				if !dup {
					p.refs = append(p.refs, r)
				}
				break
			}
		}
	}
}
