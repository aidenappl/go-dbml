// Package emit pretty-prints a model.Database back to DBML source.
//
// The output is canonical: two-space indentation, identifiers quoted only
// when they must be, settings in a fixed order, one blank line between
// top-level elements. Elements are written in source order when the model
// carries tokens (a parsed model) and in a fixed kind order otherwise (a
// model built by an importer).
//
// Round-trip property (enforced by emit_test.go over every compiling input
// in testdata/): Parse(DBML(Parse(src))) equals Parse(src) ignoring token
// positions, and DBML is idempotent: DBML(Parse(DBML(db))) == DBML(db).
package emit

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/aidenappl/go-dbml/model"
)

// Options tunes the output.
type Options struct {
	// OmitRecords drops Records blocks.
	OmitRecords bool
}

// ErrUnrepresentable is wrapped by errors for models DBML cannot express
// faithfully (an expression containing a backtick: DBML has no escape for
// it inside a backtick expression).
var ErrUnrepresentable = errors.New("emit: value cannot be represented in DBML")

// ErrNilDatabase is returned for a nil model.
var ErrNilDatabase = errors.New("emit: nil database")

// InternalError reports a recovered panic inside the emitter (a bug).
type InternalError struct{ Value any }

func (e *InternalError) Error() string { return fmt.Sprintf("emit: internal error: %v", e.Value) }

// DBML renders db as DBML source with default options.
func DBML(db *model.Database) (string, error) { return DBMLWith(db, Options{}) }

// DBMLWith renders db as DBML source. It never panics: a nil model,
// an unrepresentable value (ErrUnrepresentable) or an internal error is
// returned as an error.
func DBMLWith(db *model.Database, opts Options) (out string, err error) {
	if db == nil {
		return "", ErrNilDatabase
	}
	defer func() {
		if r := recover(); r != nil {
			if u, ok := r.(unrepresentable); ok {
				out, err = "", fmt.Errorf("%w: %s", ErrUnrepresentable, string(u))
				return
			}
			out, err = "", &InternalError{Value: r}
		}
	}()
	if err := validate(db); err != nil {
		return "", err
	}
	e := &emitter{db: db, opts: opts}
	return e.run(), nil
}

// unrepresentable is panicked by expr and turned into ErrUnrepresentable.
type unrepresentable string

type emitter struct {
	db   *model.Database
	opts Options
	sb   strings.Builder
	// db.Refs / db.Deps entries that are emitted inline (on a column or a
	// table header) instead of as standalone elements.
	inlineRef map[*model.Ref]bool
	inlineDep map[*model.Dep]bool
	// table-header deps per table
	headerDeps map[*model.Table][]*model.Dep
	// records nested in a table are emitted as top-level Records blocks
}

type item struct {
	kind   int
	offset int
	seq    int
	write  func()
}

// kind order used when the model has no tokens
const (
	kProject = iota
	kEnum
	kPartial
	kTable
	kGroup
	kRef
	kDep
	kNote
	kRecords
	kView
)

func tokOffset(t *model.Token) int {
	if t == nil {
		return -1
	}
	return t.Start.Offset
}

func (e *emitter) run() string {
	db := e.db
	e.inlineRef = map[*model.Ref]bool{}
	e.inlineDep = map[*model.Dep]bool{}
	e.headerDeps = map[*model.Table][]*model.Dep{}
	e.matchInline()

	external := map[string]bool{}
	for _, x := range db.Externals.Tables {
		external["t|"+ptr(x.SchemaName)+"|"+x.Name] = true
	}
	for _, x := range db.Externals.Enums {
		external["e|"+ptr(x.SchemaName)+"|"+x.Name] = true
	}
	for _, x := range db.Externals.TablePartials {
		external["p|"+ptr(x.SchemaName)+"|"+x.Name] = true
	}
	for _, x := range db.Externals.TableGroups {
		external["g|"+ptr(x.SchemaName)+"|"+x.Name] = true
	}
	for _, x := range db.Externals.Notes {
		external["n|"+ptr(x.SchemaName)+"|"+x.Name] = true
	}

	var items []item
	add := func(kind int, tok *model.Token, w func()) {
		items = append(items, item{kind: kind, offset: tokOffset(tok), seq: len(items), write: w})
	}
	if db.Project != nil {
		p := db.Project
		add(kProject, p.Token, func() { e.project(p) })
	}
	for _, en := range db.Enums {
		if external["e|"+ptr(en.SchemaName)+"|"+en.Name] {
			continue
		}
		en := en
		add(kEnum, en.Token, func() { e.enum(en) })
	}
	for _, p := range db.TablePartials {
		if external["p||"+p.Name] {
			continue
		}
		p := p
		add(kPartial, p.Token, func() { e.partial(p) })
	}
	for _, t := range db.Tables {
		if external["t|"+ptr(t.SchemaName)+"|"+t.Name] {
			continue
		}
		t := t
		add(kTable, t.Token, func() { e.table(t) })
	}
	for _, g := range db.TableGroups {
		if g.Name != nil && external["g|"+ptr(g.SchemaName)+"|"+*g.Name] {
			continue
		}
		g := g
		add(kGroup, g.Token, func() { e.group(g) })
	}
	for _, r := range db.Refs {
		if e.inlineRef[r] {
			continue
		}
		r := r
		add(kRef, r.Token, func() { e.ref(r) })
	}
	for _, d := range db.Deps {
		if e.inlineDep[d] {
			continue
		}
		d := d
		add(kDep, d.Token, func() { e.dep(d) })
	}
	for _, n := range db.Notes {
		if external["n||"+n.Name] {
			continue
		}
		n := n
		add(kNote, n.Token, func() { e.stickyNote(n) })
	}
	if !e.opts.OmitRecords {
		for _, r := range db.Records {
			r := r
			add(kRecords, r.Token, func() { e.records(r) })
		}
	}
	for _, v := range db.DiagramViews {
		v := v
		add(kView, v.Token, func() { e.view(v) })
	}

	allTokens := true
	for _, it := range items {
		if it.offset < 0 {
			allTokens = false
			break
		}
	}
	sort.SliceStable(items, func(i, j int) bool {
		a, b := items[i], items[j]
		if allTokens && a.offset != b.offset {
			return a.offset < b.offset
		}
		if !allTokens && a.kind != b.kind {
			return a.kind < b.kind
		}
		return a.seq < b.seq
	})
	for i, it := range items {
		if i > 0 {
			e.sb.WriteString("\n")
		}
		it.write()
	}
	return e.sb.String()
}

// matchInline marks db.Refs/db.Deps entries that correspond to inline
// `[ref: ...]` / `[dep: ...]` settings so they are not emitted twice.
func (e *emitter) matchInline() {
	db := e.db
	refKey := func(a, b model.RefEndpoint) string {
		return ptr(a.SchemaName) + "|" + a.TableName + "|" + strings.Join(a.FieldNames, ",") + "|" + a.Relation + "||" +
			ptr(b.SchemaName) + "|" + b.TableName + "|" + strings.Join(b.FieldNames, ",") + "|" + b.Relation
	}
	plainRefs := map[string][]*model.Ref{}
	for _, r := range db.Refs {
		if r.Name != nil || r.Color != nil || r.OnDelete != nil || r.OnUpdate != nil || r.Inactive != nil {
			continue
		}
		k := refKey(r.Endpoints[0], r.Endpoints[1])
		plainRefs[k] = append(plainRefs[k], r)
	}
	depKey := func(up, down model.DepEndpoint) string {
		return ptr(up.SchemaName) + "|" + up.TableName + "|" + strings.Join(up.FieldNames, ",") + "||" +
			ptr(down.SchemaName) + "|" + down.TableName + "|" + strings.Join(down.FieldNames, ",")
	}
	plainDeps := map[string][]*model.Dep{}
	for _, d := range db.Deps {
		if d.Name != nil || d.Color != nil || d.Note != nil || len(d.Metadata) > 0 || len(d.Edges) != 1 {
			continue
		}
		k := depKey(d.Edges[0].Upstream, d.Edges[0].Downstream)
		plainDeps[k] = append(plainDeps[k], d)
	}
	take := func(m map[string][]*model.Ref, k string) {
		if l := m[k]; len(l) > 0 {
			e.inlineRef[l[0]] = true
			m[k] = l[1:]
		}
	}
	takeDep := func(k string) bool {
		if l := plainDeps[k]; len(l) > 0 {
			e.inlineDep[l[0]] = true
			plainDeps[k] = l[1:]
			return true
		}
		return false
	}
	for _, t := range db.Tables {
		for _, c := range t.Fields {
			for _, ir := range c.InlineRefs {
				m, ok := relOps[ir.Relation]
				if !ok {
					continue
				}
				target := model.RefEndpoint{SchemaName: ir.SchemaName, TableName: ir.TableName, FieldNames: ir.FieldNames, Relation: m[1]}
				owner := model.RefEndpoint{SchemaName: t.SchemaName, TableName: t.Name, FieldNames: []string{c.Name}, Relation: m[0]}
				take(plainRefs, refKey(target, owner))
			}
			for _, id := range c.InlineDeps {
				host := model.DepEndpoint{SchemaName: t.SchemaName, TableName: t.Name, FieldNames: []string{c.Name}}
				other := model.DepEndpoint{SchemaName: id.SchemaName, TableName: id.TableName, FieldNames: id.FieldNames}
				if id.Direction == "<-" {
					takeDep(depKey(other, host))
				} else {
					takeDep(depKey(host, other))
				}
			}
		}
	}
	// Table-header deps: only recognisable through tokens (a plain dep
	// whose token lies inside the table's declaration).
	for _, d := range db.Deps {
		if e.inlineDep[d] || d.Token == nil || d.Name != nil || d.Color != nil || d.Note != nil || len(d.Metadata) > 0 || len(d.Edges) != 1 {
			continue
		}
		ed := d.Edges[0]
		if len(ed.Upstream.FieldNames) > 0 || len(ed.Downstream.FieldNames) > 0 {
			continue
		}
		for _, t := range db.Tables {
			if t.Token == nil || t.Token.Filepath != d.Token.Filepath || d.Token.Start.Offset <= t.Token.Start.Offset || d.Token.End.Offset > t.Token.End.Offset {
				continue
			}
			isT := func(ep model.DepEndpoint) bool {
				return ep.TableName == t.Name && ptr(ep.SchemaName) == ptr(t.SchemaName)
			}
			if isT(ed.Upstream) || isT(ed.Downstream) {
				e.inlineDep[d] = true
				e.headerDeps[t] = append(e.headerDeps[t], d)
			}
			break
		}
	}
}

// relOps maps a relationship operator to its endpoint cardinalities.
var relOps = map[string][2]string{
	"-": {"1", "1"}, "<>": {"*", "*"}, ">": {"*", "1"}, "<": {"1", "*"},
	"-?": {"1", "0..1"}, "?-": {"0..1", "1"}, "?-?": {"0..1", "0..1"},
	"?>": {"0..*", "1"}, ">?": {"*", "0..1"}, "?>?": {"0..*", "0..1"},
	"?<": {"0..1", "*"}, "<?": {"1", "0..*"}, "?<?": {"0..1", "0..*"},
	"?<>": {"0..*", "*"}, "<>?": {"*", "0..*"}, "?<>?": {"0..*", "0..*"},
}

// relationOp returns the DBML operator for a pair of endpoint
// cardinalities ("1", "*", "0..1", "0..*"), or "" if there is none.
func relationOp(left, right string) string {
	for op, m := range relOps {
		if m[0] == left && m[1] == right {
			return op
		}
	}
	return ""
}

// ---------------------------------------------------------------- writers

func (e *emitter) w(s ...string) {
	for _, x := range s {
		e.sb.WriteString(x)
	}
}

func (e *emitter) project(p *model.Project) {
	e.w("Project")
	if p.Name != nil {
		e.w(" ", Ident(*p.Name))
	}
	e.w(" {\n")
	keys := make([]string, 0, len(p.Extra))
	for k := range p.Extra {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		e.w("  ", Ident(k), ": ", String(p.Extra[k]), "\n")
	}
	if p.Note != nil {
		e.w("  Note: ", noteLiteral(p.Note.Value, "  "), "\n")
	}
	e.w("}\n")
}

func (e *emitter) enum(en *model.Enum) {
	e.w("Enum ", qualified(en.SchemaName, en.Name), " {\n")
	for _, v := range en.Values {
		e.w("  ", Ident(v.Name))
		if v.Note != nil {
			e.w(" [note: ", String(v.Note.Value), "]")
		}
		e.w("\n")
	}
	e.w("}\n")
}

func (e *emitter) partial(p *model.TablePartial) {
	e.w("TablePartial ", Ident(p.Name))
	var set []string
	if p.HeaderColor != nil {
		set = append(set, "headerColor: "+color(*p.HeaderColor))
	}
	e.settings(set)
	e.w(" {\n")
	for _, c := range p.Fields {
		e.column(c)
	}
	e.indexes(p.Indexes)
	e.checks(p.Checks)
	if p.Note != nil {
		e.w("\n  Note: ", noteLiteral(p.Note.Value, "  "), "\n")
	}
	e.w("}\n")
}

func (e *emitter) table(t *model.Table) {
	e.w("Table ", qualified(t.SchemaName, t.Name))
	if t.Alias != nil {
		e.w(" as ", Ident(*t.Alias))
	}
	var set []string
	if t.HeaderColor != nil {
		set = append(set, "headerColor: "+color(*t.HeaderColor))
	}
	for _, d := range e.headerDeps[t] {
		ed := d.Edges[0]
		if ed.Downstream.TableName == t.Name && ptr(ed.Downstream.SchemaName) == ptr(t.SchemaName) &&
			!(ed.Upstream.TableName == t.Name && ptr(ed.Upstream.SchemaName) == ptr(t.SchemaName)) {
			set = append(set, "dep: <- "+qualified(ed.Upstream.SchemaName, ed.Upstream.TableName))
		} else {
			set = append(set, "dep: -> "+qualified(ed.Downstream.SchemaName, ed.Downstream.TableName))
		}
	}
	set = append(set, metadata(t.Metadata)...)
	e.settings(set)
	e.w(" {\n")
	parts := append([]*model.PartialInject{}, t.Partials...)
	sort.SliceStable(parts, func(i, j int) bool { return parts[i].Order < parts[j].Order })
	fi, pi := 0, 0
	for pos := 0; fi < len(t.Fields) || pi < len(parts); pos++ {
		if pi < len(parts) && (parts[pi].Order <= pos || fi >= len(t.Fields)) {
			e.w("  ~", Ident(parts[pi].Name), "\n")
			pi++
			continue
		}
		e.column(t.Fields[fi])
		fi++
	}
	e.indexes(t.Indexes)
	e.checks(t.Checks)
	if t.Note != nil {
		e.w("\n  Note: ", noteLiteral(t.Note.Value, "  "), "\n")
	}
	e.w("}\n")
}

func (e *emitter) column(c *model.Column) {
	e.w("  ", Ident(c.Name), " ", typeString(c.Type))
	var set []string
	if c.PK != nil && *c.PK {
		set = append(set, "pk")
	}
	if c.Increment != nil && *c.Increment {
		set = append(set, "increment")
	}
	if c.NotNull != nil {
		if *c.NotNull {
			set = append(set, "not null")
		} else {
			set = append(set, "null")
		}
	}
	if c.Unique != nil && *c.Unique {
		set = append(set, "unique")
	}
	if c.DBDefault != nil {
		if d := defaultLiteral(c.DBDefault); d != "" {
			set = append(set, "default: "+d)
		}
	}
	for _, r := range c.InlineRefs {
		set = append(set, "ref: "+r.Relation+" "+endpointString(r.SchemaName, r.TableName, r.FieldNames))
	}
	for _, d := range c.InlineDeps {
		set = append(set, "dep: "+d.Direction+" "+endpointString(d.SchemaName, d.TableName, d.FieldNames))
	}
	for _, ch := range c.Checks {
		set = append(set, "check: "+expr(ch.Expression))
	}
	if c.Note != nil {
		set = append(set, "note: "+String(c.Note.Value))
	}
	set = append(set, metadata(columnMetadata(c))...)
	e.settings(set)
	e.w("\n")
}

func (e *emitter) indexes(ixs []*model.Index) {
	if len(ixs) == 0 {
		return
	}
	e.w("\n  indexes {\n")
	for _, ix := range ixs {
		var cols []string
		for _, c := range ix.Columns {
			if c.Type == "expression" {
				cols = append(cols, expr(c.Value))
			} else {
				cols = append(cols, Ident(c.Value))
			}
		}
		if len(cols) == 1 {
			e.w("    ", cols[0])
		} else {
			e.w("    (", strings.Join(cols, ", "), ")")
		}
		var set []string
		if ix.PK != nil && *ix.PK {
			set = append(set, "pk")
		}
		if ix.Unique != nil && *ix.Unique {
			set = append(set, "unique")
		}
		if ix.Name != nil {
			set = append(set, "name: "+String(*ix.Name))
		}
		if ix.Type != nil {
			set = append(set, "type: "+bareOrQuoted(*ix.Type))
		}
		if ix.Note != nil {
			set = append(set, "note: "+String(ix.Note.Value))
		}
		e.settings(set)
		e.w("\n")
	}
	e.w("  }\n")
}

func (e *emitter) checks(cs []*model.Check) {
	if len(cs) == 0 {
		return
	}
	e.w("\n  checks {\n")
	for _, c := range cs {
		e.w("    ", expr(c.Expression))
		if c.Name != nil {
			e.w(" [name: ", String(*c.Name), "]")
		}
		e.w("\n")
	}
	e.w("  }\n")
}

func (e *emitter) group(g *model.TableGroup) {
	name := ""
	if g.Name != nil {
		name = *g.Name
	}
	e.w("TableGroup ", qualified(g.SchemaName, name))
	var set []string
	if g.Color != nil {
		set = append(set, "color: "+color(*g.Color))
	}
	set = append(set, metadata(g.Metadata)...)
	e.settings(set)
	e.w(" {\n")
	for _, it := range g.Tables {
		e.w("  ", qualified(it.SchemaName, it.Name), "\n")
	}
	if g.Note != nil {
		if len(g.Tables) > 0 {
			e.w("\n")
		}
		e.w("  Note: ", noteLiteral(g.Note.Value, "  "), "\n")
	}
	e.w("}\n")
}

func (e *emitter) ref(r *model.Ref) {
	e.w("Ref")
	if r.Name != nil {
		e.w(" ", qualified(r.SchemaName, *r.Name))
	}
	a, b := r.Endpoints[0], r.Endpoints[1]
	op := relationOp(a.Relation, b.Relation)
	if op == "" {
		op = ">"
	}
	e.w(": ", endpointString(a.SchemaName, a.TableName, a.FieldNames), " ", op, " ", endpointString(b.SchemaName, b.TableName, b.FieldNames))
	var set []string
	if r.OnDelete != nil {
		set = append(set, "delete: "+*r.OnDelete)
	}
	if r.OnUpdate != nil {
		set = append(set, "update: "+*r.OnUpdate)
	}
	if r.Color != nil {
		set = append(set, "color: "+color(*r.Color))
	}
	if r.Inactive != nil && *r.Inactive {
		set = append(set, "inactive")
	}
	e.settings(set)
	e.w("\n")
}

func depEnd(ep model.DepEndpoint) string {
	if len(ep.FieldNames) == 0 {
		return qualified(ep.SchemaName, ep.TableName)
	}
	return endpointString(ep.SchemaName, ep.TableName, ep.FieldNames)
}

func (e *emitter) dep(d *model.Dep) {
	e.w("Dep")
	if d.Name != nil {
		e.w(" ", qualified(d.SchemaName, *d.Name))
	}
	var set []string
	if d.Color != nil {
		set = append(set, "color: "+color(*d.Color))
	}
	if d.Note != nil {
		set = append(set, "note: "+String(d.Note.Value))
	}
	keys := make([]string, 0, len(d.Metadata))
	for k := range d.Metadata {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		set = append(set, Ident(k)+": "+anyLiteral(d.Metadata[k]))
	}
	if len(d.Edges) == 1 {
		ed := d.Edges[0]
		e.w(": ", depEnd(ed.Upstream), " -> ", depEnd(ed.Downstream))
		e.settings(set)
		e.w("\n")
		return
	}
	e.settings(set)
	e.w(" {\n")
	for _, ed := range d.Edges {
		e.w("  ", depEnd(ed.Upstream), " -> ", depEnd(ed.Downstream), "\n")
	}
	e.w("}\n")
}

func (e *emitter) stickyNote(n *model.StickyNote) {
	e.w("Note ", Ident(n.Name))
	var set []string
	if n.Color != nil {
		set = append(set, "color: "+color(*n.Color))
	}
	set = append(set, metadata(n.Metadata)...)
	e.settings(set)
	e.w(" {\n  ", noteLiteral(n.Content, "  "), "\n}\n")
}

func (e *emitter) records(r *model.TableRecord) {
	var cols []string
	for _, c := range r.Columns {
		cols = append(cols, Ident(c))
	}
	e.w("Records ", qualified(r.SchemaName, r.TableName), "(", strings.Join(cols, ", "), ")")
	if r.Example != nil && *r.Example {
		e.w(" [example]")
	}
	e.w(" {\n")
	for _, row := range r.Values {
		var vals []string
		for i, v := range row {
			vals = append(vals, e.recordValue(r, i, v))
		}
		e.w("  ", strings.Join(vals, ", "), "\n")
	}
	e.w("}\n")
}

func (e *emitter) recordValue(r *model.TableRecord, i int, v *model.RecordValue) string {
	if v == nil {
		return "null"
	}
	if v.Undefined {
		// An `Enum.member` reference: rebuild one from the column's enum type.
		if en := e.enumOfColumn(r, i); en != nil && len(en.Values) > 0 {
			return qualified(en.SchemaName, en.Name) + "." + Ident(en.Values[0].Name)
		}
		return "null"
	}
	switch x := v.Value.(type) {
	case nil:
		return "null"
	case bool:
		if x {
			return "true"
		}
		return "false"
	case float64:
		return number(x)
	case string:
		if v.Type == "expression" {
			return expr(x)
		}
		return String(x)
	}
	return "null"
}

func (e *emitter) enumOfColumn(r *model.TableRecord, i int) *model.Enum {
	if i >= len(r.Columns) {
		return nil
	}
	for _, t := range e.db.Tables {
		if t.Name != r.TableName || ptr(t.SchemaName) != ptr(r.SchemaName) {
			continue
		}
		for _, c := range t.Fields {
			if c.Name != r.Columns[i] {
				continue
			}
			for _, en := range e.db.Enums {
				if en.Name == c.Type.TypeName && ptr(en.SchemaName) == ptr(c.Type.SchemaName) {
					return en
				}
			}
		}
	}
	return nil
}

func (e *emitter) view(v *model.DiagramView) {
	e.w("DiagramView ", qualified(v.SchemaName, v.Name), " {\n")
	ve := v.VisibleEntities
	allEmpty := ve.Tables != nil && len(ve.Tables) == 0 && ve.TableGroups != nil && len(ve.TableGroups) == 0 &&
		ve.Schemas != nil && len(ve.Schemas) == 0
	if allEmpty && ve.StickyNotes != nil && len(ve.StickyNotes) == 0 {
		e.w("  *\n}\n")
		return
	}
	block := func(kw string, items []string) {
		if len(items) == 0 {
			e.w("  ", kw, " {}\n")
			return
		}
		e.w("  ", kw, " {\n")
		for _, it := range items {
			e.w("    ", it, "\n")
		}
		e.w("  }\n")
	}
	named := func(xs []model.NamedItem) []string {
		var out []string
		for _, x := range xs {
			out = append(out, Ident(x.Name))
		}
		return out
	}
	if allEmpty {
		e.w("  Tables {\n    *\n  }\n")
	} else {
		// A non-nil list in any of tables/groups/schemas makes the omitted
		// ones default to "show all" (empty); nil ones must be written as
		// explicit empty blocks.
		anyList := len(ve.Tables) > 0 || len(ve.TableGroups) > 0 || len(ve.Schemas) > 0
		if len(ve.Tables) > 0 {
			var items []string
			for _, t := range ve.Tables {
				var sch *string
				if t.SchemaName != "" && t.SchemaName != "public" {
					s := t.SchemaName
					sch = &s
				}
				items = append(items, qualified(sch, t.Name))
			}
			block("Tables", items)
		} else if ve.Tables == nil && anyList {
			block("Tables", nil)
		}
		if len(ve.TableGroups) > 0 {
			block("TableGroups", named(ve.TableGroups))
		} else if ve.TableGroups == nil && anyList {
			block("TableGroups", nil)
		}
		if len(ve.Schemas) > 0 {
			block("Schemas", named(ve.Schemas))
		} else if ve.Schemas == nil && anyList {
			block("Schemas", nil)
		}
	}
	if ve.StickyNotes != nil {
		if len(ve.StickyNotes) == 0 {
			e.w("  Notes {\n    *\n  }\n")
		} else {
			block("Notes", named(ve.StickyNotes))
		}
	}
	e.w("}\n")
}

func (e *emitter) settings(set []string) {
	if len(set) == 0 {
		return
	}
	e.w(" [", strings.Join(set, ", "), "]")
}

// ---------------------------------------------------------------- literals

var (
	identRe   = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	numberLit = regexp.MustCompile(`^-?([0-9]+\.?[0-9]*|\.[0-9]+)([eE][-+]?[0-9]+)?$`)
)

// keywords that must be quoted when used as a bare name.
var keywords = map[string]bool{
	"table": true, "enum": true, "ref": true, "dep": true, "note": true, "project": true, "indexes": true,
	"tablegroup": true, "tablepartial": true, "checks": true, "records": true, "diagramview": true,
	"metadata": true, "as": true, "null": true, "true": true, "false": true, "use": true, "reuse": true,
	"tables": true, "notes": true, "tablegroups": true, "schemas": true,
}

// Ident renders a name, quoting it with double quotes unless it is a plain
// identifier.
func Ident(s string) string {
	if identRe.MatchString(s) && !keywords[strings.ToLower(s)] {
		return s
	}
	return quote(s, '"')
}

func bareOrQuoted(s string) string {
	if identRe.MatchString(s) {
		return s
	}
	return String(s)
}

func qualified(schema *string, name string) string {
	if schema != nil && *schema != "" {
		return Ident(*schema) + "." + Ident(name)
	}
	return Ident(name)
}

func endpointString(schema *string, table string, fields []string) string {
	t := qualified(schema, table)
	if len(fields) == 1 {
		return t + "." + Ident(fields[0])
	}
	var fs []string
	for _, f := range fields {
		fs = append(fs, Ident(f))
	}
	return t + ".(" + strings.Join(fs, ", ") + ")"
}

// String renders a single-line single-quoted DBML string literal.
func String(s string) string { return quote(s, '\'') }

func quote(s string, q rune) string {
	var b strings.Builder
	b.WriteRune(q)
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case q:
			b.WriteRune('\\')
			b.WriteRune(r)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			if r < 0x20 || r == 0x7f || r == utf8.RuneError {
				b.WriteString(`\u` + hex4(r))
			} else {
				b.WriteRune(r)
			}
		}
	}
	b.WriteRune(q)
	return b.String()
}

func hex4(r rune) string {
	s := strconv.FormatInt(int64(r), 16)
	return strings.Repeat("0", 4-len(s)) + s
}

// expr renders a backtick expression. A backtick inside one cannot be
// written in DBML, so such a value aborts the emit (ErrUnrepresentable)
// rather than being silently changed.
func expr(s string) string {
	if strings.Contains(s, "`") {
		panic(unrepresentable("expression contains a backtick: " + s))
	}
	return "`" + s + "`"
}

// noteLiteral renders a note. Multi-line values whose normalised form is
// themselves become indented ”' blocks; anything else a single-line string.
func noteLiteral(s, indent string) string {
	if !strings.Contains(s, "\n") || !multilineSafe(s) {
		return String(s)
	}
	inner := indent + "  "
	var b strings.Builder
	b.WriteString("'''\n")
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		if l != "" {
			b.WriteString(inner)
		}
		b.WriteString(escapeMultiline(l, i == len(lines)-1))
		if i < len(lines)-1 {
			b.WriteString("\n")
		}
	}
	b.WriteString("'''")
	return b.String()
}

// multilineSafe reports whether DBML's note normalisation (drop leading
// blank lines, strip the common indentation) leaves s unchanged.
func multilineSafe(s string) bool {
	lines := strings.Split(s, "\n")
	if strings.TrimLeft(lines[0], " \t") == "" {
		return false
	}
	for _, l := range lines {
		if strings.TrimLeft(l, " \t") == "" {
			if l != "" && !isAllSpaceTab(l) {
				return false
			}
			continue
		}
		if l[0] != ' ' && l[0] != '\t' {
			return true
		}
	}
	return false
}

func isAllSpaceTab(s string) bool { return strings.Trim(s, " \t") == "" }

func escapeMultiline(l string, last bool) string {
	var b strings.Builder
	rs := []rune(l)
	for i, r := range rs {
		switch {
		case r == '\\':
			b.WriteString(`\\`)
		case r == '\'' && ((i+1 < len(rs) && rs[i+1] == '\'') || (last && i == len(rs)-1)):
			b.WriteString(`\'`)
		case r == '\r':
			b.WriteString(`\r`)
		case r == '\t':
			b.WriteRune(r)
		case r < 0x20 || r == 0x7f:
			b.WriteString(`\u` + hex4(r))
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func number(f float64) string {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return "null"
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func anyLiteral(v any) string {
	switch x := v.(type) {
	case float64:
		return number(x)
	case string:
		return String(x)
	case bool:
		if x {
			return "true"
		}
		return "false"
	}
	return "null"
}

func defaultLiteral(d *model.Default) string {
	switch d.Type {
	case "number":
		switch x := d.Value.(type) {
		case float64:
			return number(x)
		case string:
			if numberLit.MatchString(x) {
				return x
			}
			return String(x)
		}
	case "boolean":
		if s, ok := d.Value.(string); ok {
			if identRe.MatchString(s) {
				return s
			}
		}
		if b, ok := d.Value.(bool); ok {
			if b {
				return "true"
			}
			return "false"
		}
	case "expression":
		if s, ok := d.Value.(string); ok {
			return expr(s)
		}
	default:
		if s, ok := d.Value.(string); ok {
			return String(s)
		}
	}
	if d.Value == nil {
		return "null"
	}
	return ""
}

func color(c string) string {
	if strings.HasPrefix(c, "#") && identRe.MatchString("x"+c[1:]) || strings.EqualFold(c, "none") {
		return c
	}
	return String(c)
}

func metadata(m map[string]string) []string {
	if len(m) == 0 {
		return nil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, settingKey(k)+": "+String(m[k]))
	}
	return out
}

// settingKey renders a custom-property key. Keys come from identifier
// streams (settings) or element names (Metadata blocks) and DBML has no
// quoted form for them, so they are written as is.
func settingKey(k string) string { return k }

var (
	bareTypeRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*(\[[0-9]*\])*$`)
	argRe      = regexp.MustCompile(`^(-?[0-9]+(\.[0-9]+)?|[A-Za-z_][A-Za-z0-9_]*)$`)
)

func typeString(t model.ColumnType) string {
	name := t.TypeName
	args := ""
	hasArgs := false
	if t.Args != nil && strings.HasSuffix(name, "("+*t.Args+")") {
		name = strings.TrimSuffix(name, "("+*t.Args+")")
		hasArgs = true
		var parts []string
		if *t.Args != "" {
			for _, a := range strings.Split(*t.Args, ",") {
				if argRe.MatchString(a) {
					parts = append(parts, a)
				} else {
					parts = append(parts, String(a))
				}
			}
		}
		args = "(" + strings.Join(parts, ",") + ")"
	}
	var base string
	switch {
	case name == "" && hasArgs:
		base = `""`
	case bareTypeRe.MatchString(name) && !keywords[strings.ToLower(name)]:
		base = name
	default:
		base = quote(name, '"')
	}
	if t.SchemaName != nil && *t.SchemaName != "" {
		base = Ident(*t.SchemaName) + "." + base
	}
	return base + args
}

func ptr(p *string) string {
	if p == nil {
		return "\x00"
	}
	return *p
}

// columnMetadata is c's custom properties plus the generated /
// generated_stored keys for a Generated column built without them.
func columnMetadata(c *model.Column) map[string]string {
	g := c.Generated
	if g == nil || c.Metadata[model.MetaGenerated] != "" {
		return c.Metadata
	}
	m := make(map[string]string, len(c.Metadata)+2)
	for k, v := range c.Metadata {
		m[k] = v
	}
	m[model.MetaGenerated] = g.Expression
	if g.Stored {
		m[model.MetaGeneratedStored] = "true"
	}
	return m
}
