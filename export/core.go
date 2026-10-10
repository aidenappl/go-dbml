// Portions of this file are derived from holistics/dbml packages/dbml-core/src/model_structure/
// {database.ts,schema.ts,table.ts,field.ts,ref.js,endpoint.js,depEdge.ts,tableGroup.ts}
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go; only the parts the SQL exporters observe are kept.

package export

import (
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/aidenappl/go-dbml/model"
)

const defaultSchema = "public"

// core mirrors the @dbml/core Database built from @dbml/parse JSON
// (model_structure/), reduced to what the SQL exporters read. Element order
// and schema creation order follow the reference exactly because they
// decide statement order in the output.
type core struct {
	db      *model.Database
	schemas []*cSchema
	// partial field -> shared type object (the reference mutates the raw
	// JSON type object, which is shared by every table the partial is
	// injected into)
	partialTypes map[*model.Column]*cType
	injected     []rawRef
	notes        []*model.StickyNote

	// lookup indexes (keep buildCore linear on 20k-table models)
	schemaIdx  map[string]*cSchema
	tableIdx   map[*cSchema]map[string]*cTable
	enumIdx    map[string]*cEnum
	aliasIdx   map[string]*model.Alias
	partialIdx map[string]*model.TablePartial
	bareTables map[string]bool
	refKeys    map[*cSchema]map[string]bool
	nextField  int
	// translated: SQLWith translated the model (MySQL schemas are then
	// created IF NOT EXISTS: a MySQL schema is a server-wide database)
	translated bool
}

type cSchema struct {
	name   string
	tables []*cTable
	enums  []*cEnum
	refs   []*cRef
	groups []*cGroup
	deps   []*model.Dep
}

type cGroup struct {
	g      *model.TableGroup
	tables []*cTable
}

type cTable struct {
	name    string
	schema  *cSchema
	fields  []*cField
	indexes []*model.Index
	checks  []cCheck
	note    string
	grouped bool
	src     *model.Table
	color   *string
	byName  map[string]*cField
}

type cCheck struct {
	name       *string
	expression string
}

type cType struct {
	schemaName       *string
	typeName         string
	originalTypeName *string
}

type cField struct {
	name                           string
	typ                            *cType
	unique, pk, notNull, increment bool
	note                           string
	dbdefault                      *model.Default
	checks                         []cCheck
	enum                           *cEnum
	table                          *cTable
	metadata                       map[string]string
	generated                      *model.GeneratedColumn
	id                             int
}

type cEnum struct {
	name       string
	schema     *cSchema
	values     []string
	valueNotes []string
}

type cEndpoint struct {
	relation string
	fields   []*cField
}

type cRef struct {
	name               *string
	onDelete, onUpdate *string
	inactive           bool
	schema             *cSchema
	endpoints          [2]*cEndpoint
}

type rawRef struct {
	endpoints [2]model.RefEndpoint
}

// ErrModel is wrapped by errors the reference model builder would throw
// (e.g. "Two endpoints are the same").
var ErrModel = errors.New("dbml model error")

type modelError struct{ msg string }

func (e *modelError) Error() string { return e.msg }
func (e *modelError) Unwrap() error { return ErrModel }
func fail(msg string) error         { return &modelError{msg} }
func str(p *string) string          { return deref(p, "") }
func deref(p *string, def string) string {
	if p == nil {
		return def
	}
	return *p
}
func truthy(b *bool) bool { return b != nil && *b }

// buildCore rebuilds @dbml/core's Database from db, in the reference's
// processing order (it decides schema and statement order).
func buildCore(db *model.Database) (*core, error) {
	c := &core{db: db, partialTypes: map[*model.Column]*cType{},
		schemaIdx: map[string]*cSchema{}, tableIdx: map[*cSchema]map[string]*cTable{}, enumIdx: map[string]*cEnum{},
		aliasIdx: map[string]*model.Alias{}, partialIdx: map[string]*model.TablePartial{}, bareTables: map[string]bool{},
		refKeys: map[*cSchema]map[string]bool{}}
	for _, a := range db.Aliases {
		if _, ok := c.aliasIdx[a.Name]; !ok {
			c.aliasIdx[a.Name] = a
		}
	}
	for _, p := range db.TablePartials {
		if _, ok := c.partialIdx[p.Name]; !ok {
			c.partialIdx[p.Name] = p
		}
	}
	for _, t := range db.Tables {
		if t.SchemaName == nil {
			c.bareTables[t.Name] = true
		}
	}
	for _, en := range db.Enums {
		s := c.findOrCreateSchema(deref(en.SchemaName, defaultSchema))
		e := &cEnum{name: en.Name, schema: s}
		for _, v := range en.Values {
			e.values = append(e.values, v.Name)
			note := ""
			if v.Note != nil {
				note = v.Note.Value
			}
			e.valueNotes = append(e.valueNotes, note)
		}
		s.enums = append(s.enums, e)
		if k := s.name + "\x00" + e.name; c.enumIdx[k] == nil {
			c.enumIdx[k] = e
		}
	}
	if err := validateTokens(db); err != nil {
		return nil, err
	}
	for _, t := range db.Tables {
		s := c.findOrCreateSchema(deref(t.SchemaName, defaultSchema))
		ct, err := c.table(t, s)
		if err != nil {
			return nil, err
		}
		s.tables = append(s.tables, ct)
		idx := c.tableIdx[s]
		if idx == nil {
			idx = map[string]*cTable{}
			c.tableIdx[s] = idx
		}
		if idx[ct.name] == nil {
			idx[ct.name] = ct
		}
	}
	noteNames := map[string]bool{}
	for _, n := range db.Notes {
		if noteNames[n.Name] {
			return nil, fail("Notes " + n.Name + " existed")
		}
		noteNames[n.Name] = true
		c.notes = append(c.notes, n)
		c.findOrCreateSchema(defaultSchema)
	}
	for _, r := range db.Refs {
		s := c.findOrCreateSchema(deref(r.SchemaName, defaultSchema))
		cr, err := c.ref(r.Endpoints, r)
		if err != nil {
			return nil, err
		}
		cr.schema = s
		if !c.addRefKey(s, cr) {
			return nil, fail("Reference with the same endpoints already exists")
		}
		s.refs = append(s.refs, cr)
	}
	for _, d := range db.Deps {
		ds := c.findOrCreateSchema(deref(d.SchemaName, defaultSchema))
		ds.deps = append(ds.deps, d)
		for _, e := range d.Edges {
			for _, ep := range []model.DepEndpoint{e.Upstream, e.Downstream} {
				sn := deref(ep.SchemaName, defaultSchema)
				if c.findTable(&sn, ep.TableName) == nil {
					return nil, fail("Can't find table \"" + ep.TableName + "\" referenced in Dep edge")
				}
			}
		}
	}
	for _, g := range db.TableGroups {
		gs := c.findOrCreateSchema(deref(g.SchemaName, defaultSchema))
		cg := &cGroup{g: g}
		gs.groups = append(gs.groups, cg)
		for _, it := range g.Tables {
			t := c.findTable(it.SchemaName, it.Name)
			if t == nil {
				return nil, fail("Table " + it.Name + " don't exist")
			}
			if t.grouped {
				return nil, fail("Table " + t.name + " is already in a group")
			}
			t.grouped = true
			cg.tables = append(cg.tables, t)
		}
	}
	for _, rr := range c.injected {
		s := c.findOrCreateSchema(defaultSchema)
		cr, err := c.ref(rr.endpoints, nil)
		if err != nil {
			return nil, err
		}
		cr.schema = s
		if c.addRefKey(s, cr) {
			s.refs = append(s.refs, cr)
		}
	}
	return c, nil
}

func (c *core) findOrCreateSchema(name string) *cSchema {
	if s := c.schemaIdx[name]; s != nil {
		return s
	}
	s := &cSchema{name: name}
	c.schemas = append(c.schemas, s)
	c.schemaIdx[name] = s
	return s
}

func (c *core) findEnum(schema, name string) *cEnum {
	return c.enumIdx[schema+"\x00"+name]
}

// endpointKey identifies an endpoint by its (unordered) field set, like
// Endpoint.equals (sorted field ids).
func endpointKey(e *cEndpoint) string {
	ids := make([]int, len(e.fields))
	for i, f := range e.fields {
		ids[i] = f.id
	}
	sort.Ints(ids)
	var b strings.Builder
	for _, id := range ids {
		b.WriteString(strconv.Itoa(id))
		b.WriteByte(',')
	}
	return b.String()
}

// addRefKey records r in its schema's ref set and reports whether it was
// new (Ref.equals: same endpoints in either order).
func (c *core) addRefKey(s *cSchema, r *cRef) bool {
	a, b := endpointKey(r.endpoints[0]), endpointKey(r.endpoints[1])
	if b < a {
		a, b = b, a
	}
	k := a + "|" + b
	m := c.refKeys[s]
	if m == nil {
		m = map[string]bool{}
		c.refKeys[s] = m
	}
	if m[k] {
		return false
	}
	m[k] = true
	return true
}

// findTable mirrors Database.findTable: an unqualified name is first
// looked up as a table alias; the schema is created if missing.
func (c *core) findTable(schema *string, name string) *cTable {
	if schema == nil || *schema == "" {
		if a := c.aliasIdx[name]; a != nil && a.Kind == "table" {
			if s := c.schemaIdx[deref(a.Value.SchemaName, defaultSchema)]; s != nil {
				if t := c.tableIdx[s][a.Value.TableName]; t != nil {
					return t
				}
			}
		}
	}
	sn := defaultSchema
	if schema != nil && *schema != "" {
		sn = *schema
	}
	return c.tableIdx[c.findOrCreateSchema(sn)][name]
}

func (t *cTable) findField(name string) *cField {
	if t.byName == nil {
		t.byName = make(map[string]*cField, len(t.fields))
		for _, f := range t.fields {
			if _, ok := t.byName[f.name]; !ok {
				t.byName[f.name] = f
			}
		}
	}
	return t.byName[name]
}

// newField builds a field and binds its enum type (Field.bindType).
func (c *core) newField(col *model.Column, t *cTable, typ *cType) *cField {
	c.nextField++
	f := &cField{
		id:   c.nextField,
		name: col.Name, typ: typ, table: t,
		unique: truthy(col.Unique), pk: truthy(col.PK), notNull: truthy(col.NotNull), increment: truthy(col.Increment),
		dbdefault: col.DBDefault, metadata: col.Metadata,
		generated: col.Generated,
	}
	if f.generated == nil {
		f.generated = model.GeneratedFromMetadata(col.Metadata)
	}
	if col.Note != nil {
		f.note = col.Note.Value
	}
	// Field.bindType
	if typ.schemaName != nil && *typ.schemaName != "" {
		if e := c.findEnum(*typ.schemaName, typ.typeName); e != nil {
			f.enum = e
		} else {
			orig := typ.typeName
			typ.typeName = *typ.schemaName + "." + orig
			typ.originalTypeName = &orig
		}
	} else if e := c.findEnum(defaultSchema, typ.typeName); e != nil {
		f.enum = e
	}
	for _, ch := range col.Checks {
		f.checks = append(f.checks, cCheck{name: ch.Name, expression: ch.Expression})
	}
	return f
}

// table builds one table, injecting TablePartials like Table.processPartials.
func (c *core) table(t *model.Table, s *cSchema) (*cTable, error) {
	ct := &cTable{name: t.Name, schema: s, src: t, color: t.HeaderColor}
	if t.Note != nil {
		ct.note = t.Note.Value
	}
	var fields []*cField
	existing := map[string]bool{}
	for _, col := range t.Fields {
		if existing[col.Name] {
			return nil, fail("Field \"" + col.Name + "\" existed in table \"" + t.Name + "\"")
		}
		existing[col.Name] = true
		fields = append(fields, c.newField(col, ct, &cType{schemaName: col.Type.SchemaName, typeName: col.Type.TypeName}))
	}
	// Table() runs processPartials before processIndexes/processChecks, so
	// injected indexes and checks come first.
	if len(t.Partials) > 0 {
		hasNote := t.Note != nil
		parts := append([]*model.PartialInject{}, t.Partials...)
		sort.SliceStable(parts, func(i, j int) bool { return parts[i].Order > parts[j].Order })
		for i := len(parts) - 1; i >= 0; i-- {
			at := clamp(parts[i].Order, len(fields))
			fields = append(fields[:at], append([]*cField{nil}, fields[at:]...)...)
		}
		for _, inj := range parts {
			p := c.partialIdx[inj.Name]
			if p == nil {
				return nil, fail("Table partial " + inj.Name + " not found")
			}
			at := clamp(inj.Order, len(fields)-1)
			var injected []*cField
			for _, col := range p.Fields {
				if existing[col.Name] {
					continue
				}
				existing[col.Name] = true
				for _, ir := range col.InlineRefs {
					target, targetSchema := ir.TableName, ir.SchemaName
					if targetSchema == nil && target == p.Name && !c.dbHasTable(p.Name) {
						target, targetSchema = t.Name, t.SchemaName
					}
					m, ok := relOps[ir.Relation]
					if !ok {
						m = [2]string{"*", "*"}
					}
					sn := s.name
					c.injected = append(c.injected, rawRef{endpoints: [2]model.RefEndpoint{
						{TableName: t.Name, SchemaName: &sn, FieldNames: []string{col.Name}, Relation: m[0]},
						{TableName: target, SchemaName: targetSchema, FieldNames: ir.FieldNames, Relation: m[1]},
					}})
				}
				typ := c.partialTypes[col]
				if typ == nil {
					typ = &cType{schemaName: col.Type.SchemaName, typeName: col.Type.TypeName}
					c.partialTypes[col] = typ
				}
				injected = append(injected, c.newField(col, ct, typ))
			}
			fields = append(fields[:at], append(injected, fields[at+1:]...)...)
			if !hasNote && p.Note != nil {
				ct.note = p.Note.Value
				hasNote = true
			}
			if ct.color == nil && p.HeaderColor != nil {
				ct.color = p.HeaderColor
			}
			ct.indexes = append(ct.indexes, p.Indexes...)
			for _, ch := range p.Checks {
				cc := cCheck{expression: ch.Expression}
				if ch.Name != nil && *ch.Name != "" {
					n := t.Name + "." + *ch.Name
					cc.name = &n
				} else {
					cc.name = ch.Name
				}
				ct.checks = append(ct.checks, cc)
			}
		}
	}
	ct.indexes = append(ct.indexes, t.Indexes...)
	for _, ch := range t.Checks {
		ct.checks = append(ct.checks, cCheck{name: ch.Name, expression: ch.Expression})
	}
	ct.fields = fields
	if len(ct.fields) == 0 {
		return nil, fail("Table must have at least one field")
	}
	return ct, nil
}

func (c *core) dbHasTable(name string) bool { return c.bareTables[name] }

// ref resolves a ref's endpoints to fields (PK fallback for empty field
// lists) and validates them like Ref/Endpoint.
func (c *core) ref(eps [2]model.RefEndpoint, r *model.Ref) (*cRef, error) {
	cr := &cRef{}
	if r != nil {
		cr.name, cr.onDelete, cr.onUpdate, cr.inactive = r.Name, r.OnDelete, r.OnUpdate, truthy(r.Inactive)
	}
	for i, ep := range eps {
		c.findOrCreateSchema(deref(ep.SchemaName, defaultSchema))
		t := c.findTable(ep.SchemaName, ep.TableName)
		if t == nil {
			return nil, fail("Can't find table \"" + ep.TableName + "\"")
		}
		ce := &cEndpoint{relation: ep.Relation}
		names := ep.FieldNames
		if len(names) == 0 {
			for _, f := range t.fields {
				if f.pk {
					names = []string{f.name}
					break
				}
			}
			if len(names) == 0 {
				for _, ix := range t.indexes {
					if truthy(ix.PK) {
						for _, col := range ix.Columns {
							names = append(names, col.Value)
						}
						break
					}
				}
			}
			if len(names) == 0 {
				return nil, fail("Can't find primary or composite key in table \"" + t.name + "\"")
			}
		}
		for _, n := range names {
			f := t.findField(n)
			if f == nil {
				return nil, fail("Can't find field \"" + n + "\" in table \"" + t.name + "\"")
			}
			ce.fields = append(ce.fields, f)
		}
		cr.endpoints[i] = ce
	}
	if endpointEquals(cr.endpoints[0], cr.endpoints[1]) {
		return nil, fail("Two endpoints are the same")
	}
	if len(cr.endpoints[0].fields) != len(cr.endpoints[1].fields) {
		return nil, fail("Two endpoints have unequal number of fields")
	}
	return cr, nil
}

func endpointEquals(a, b *cEndpoint) bool {
	if len(a.fields) != len(b.fields) {
		return false
	}
	count := map[*cField]int{}
	for _, f := range a.fields {
		count[f]++
	}
	for _, f := range b.fields {
		count[f]--
		if count[f] < 0 {
			return false
		}
	}
	return true
}

func clamp(i, n int) int {
	if i < 0 {
		return 0
	}
	if i > n {
		return n
	}
	return i
}

var relOps = map[string][2]string{
	"-": {"1", "1"}, "<>": {"*", "*"}, ">": {"*", "1"}, "<": {"1", "*"},
	"-?": {"1", "0..1"}, "?-": {"0..1", "1"}, "?-?": {"0..1", "0..1"},
	"?>": {"0..*", "1"}, ">?": {"*", "0..1"}, "?>?": {"0..*", "0..1"},
	"?<": {"0..1", "*"}, "<?": {"1", "0..*"}, "?<?": {"0..1", "0..*"},
	"?<>": {"0..*", "*"}, "<>?": {"*", "0..*"}, "?<>?": {"0..*", "0..*"},
}

// maxIsOne mirrors parseCardinality(rel).max === 1.
func maxIsOne(rel string) bool {
	switch rel {
	case "*":
		return false
	case "1", "0..1", "1..1":
		return true
	}
	// "<min>..<max>" or "<n>"
	for i := 0; i+1 < len(rel); i++ {
		if rel[i] == '.' && rel[i+1] == '.' {
			return rel[i+2:] == "1"
		}
	}
	return rel == "1"
}

// shouldPrintSchema: @dbml/core never sets hasDefaultSchema, so only
// non-public schemas are printed.
func (s *cSchema) print() bool { return s.name != defaultSchema }

var (
	actionRe    = regexp.MustCompile(`^[A-Za-z]+( [A-Za-z]+)*$`)
	indexTypeRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
)

// validateTokens rejects referential actions and index types that are not
// plain keywords: the exporters print them unquoted (ON DELETE x, USING x),
// so anything else could change the meaning of the generated DDL.
func validateTokens(db *model.Database) error {
	for _, r := range db.Refs {
		for _, a := range []*string{r.OnDelete, r.OnUpdate} {
			if a != nil && *a != "" && !actionRe.MatchString(*a) {
				return fail("Invalid referential action " + strconv.Quote(*a))
			}
		}
	}
	check := func(ixs []*model.Index) error {
		for _, ix := range ixs {
			if ix.Type != nil && *ix.Type != "" && !indexTypeRe.MatchString(*ix.Type) {
				return fail("Invalid index type " + strconv.Quote(*ix.Type))
			}
		}
		return nil
	}
	for _, t := range db.Tables {
		if err := check(t.Indexes); err != nil {
			return err
		}
	}
	for _, p := range db.TablePartials {
		if err := check(p.Indexes); err != nil {
			return err
		}
	}
	return nil
}
