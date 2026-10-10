// Portions of this file are derived from holistics/dbml packages/dbml-parse/src/core/global_modules/*/interpret.ts, src/core/global_modules/metadata/utils.ts
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

package model

import (
	"strings"

	"github.com/aidenappl/go-dbml/ast"
	"github.com/aidenappl/go-dbml/diag"
)

// metaKey identifies a Metadata target: a table/group/note, or a column.
type metaKey struct {
	s     *sym
	isCol bool
	col   string
}

type metaBlock struct {
	decl *ast.ElementDeclaration
	unit *unit
}

// metadataFor returns the Metadata blocks targeting (s, cl) in precedence
// order (imported files first, then the entry file; source order within a
// file). The index is built once per compile.
func (in *interpreter) metadataFor(s *sym, cl *col) []metaBlock {
	if in.metaIndex == nil {
		in.metaIndex = map[metaKey][]metaBlock{}
		files := in.c.reachable(in.u)
		for i := len(files) - 1; i >= 0; i-- {
			f := files[i]
			for _, e := range f.prog.Declarations() {
				if !e.IsKind(kMetadata) {
					continue
				}
				t, ok := in.c.resolveMetadataTarget(f.scope, e)
				if !ok {
					continue
				}
				k := metaKey{s: t.sym}
				if t.col != nil {
					k.isCol, k.col = true, t.col.name
				}
				in.metaIndex[k] = append(in.metaIndex[k], metaBlock{e, f})
			}
		}
	}
	k := metaKey{s: s}
	if cl != nil {
		k.isCol, k.col = true, cl.name
	}
	return in.metaIndex[k]
}

type interpreter struct {
	metaIndex map[metaKey][]metaBlock

	c     *compiler
	u     *unit
	sc    *scope
	diags []diag.Diagnostic
	// rs/rfile are the scope and path of the file whose refs, deps and
	// records are being interpreted (the entry file or one it imports).
	rs    *scope
	rfile string
	// records budget (Options.MaxRecordRows / MaxRecordBytes)
	recordRows, recordBytes int
	recordsCapped           bool
}

// visibleSym reports whether s is visible from the entry file.
func (in *interpreter) visibleSym(s *sym) bool {
	if s == nil {
		return false
	}
	if s.unit == in.u {
		return true
	}
	_, ok := in.sc.bySym[s]
	return ok
}

func (in *interpreter) allVisible(rs ...resolved) bool {
	for _, r := range rs {
		if r.table == nil || !in.visibleSym(r.table.sym) {
			return false
		}
	}
	return true
}

func (in *interpreter) n(code diag.Code, msg string, node ast.Node, file string) {
	in.diags = append(in.diags, nodeDiag(code, msg, node, file))
}

func (in *interpreter) tok(n ast.Node) *Token { return nodeToken(n, n.Base().Filepath) }

func (in *interpreter) interpret() *Database {
	in.sc = in.u.scope
	db := &Database{
		Schemas: []Schema{}, Tables: []*Table{}, Notes: []*StickyNote{}, Refs: []*Ref{}, Deps: []*Dep{}, Enums: []*Enum{},
		TableGroups: []*TableGroup{}, Aliases: []*Alias{}, TablePartials: []*TablePartial{}, Records: []*TableRecord{},
		DiagramViews: []*DiagramView{},
		Externals: Externals{Tables: []*ElementRef{}, Enums: []*ElementRef{}, TableGroups: []*ElementRef{},
			TablePartials: []*ElementRef{}, Notes: []*ElementRef{}},
		Token: in.tok(in.u.prog),
	}
	// 1. local declarations, in source order
	for _, s := range in.u.decls {
		in.c.poll(false)
		in.pushSym(db, s)
	}
	// 2. imported elements
	done := map[*sym]bool{}
	for _, v := range in.sc.entries {
		if v.sym.unit == in.u || done[v.sym] {
			continue
		}
		done[v.sym] = true
		in.pushSym(db, v.sym)
		in.pushExternal(db, v)
	}
	// 3. metadata-like elements: refs, deps, records, project
	in.interpretMetadata(db)
	// 4. aliases
	for _, s := range in.u.decls {
		if s.kind == symTable && s.alias != "" {
			schema, name := in.sc.interpretedName(s)
			db.Aliases = append(db.Aliases, &Alias{Name: s.alias, Kind: "table", Value: AliasValue{TableName: name, SchemaName: schema}})
		}
	}
	in.validateRefCardinalities(db)
	deriveGenerated(db)
	return db
}

func (in *interpreter) pushSym(db *Database, s *sym) {
	switch s.kind {
	case symTable:
		db.Tables = append(db.Tables, in.table(s))
	case symEnum:
		db.Enums = append(db.Enums, in.enum(s))
	case symGroup:
		db.TableGroups = append(db.TableGroups, in.tableGroup(s))
	case symPartial:
		db.TablePartials = append(db.TablePartials, in.partial(s))
	case symNote:
		db.Notes = append(db.Notes, in.stickyNote(s))
	case symView:
		db.DiagramViews = append(db.DiagramViews, in.diagramView(s))
	}
}

func (in *interpreter) pushExternal(db *Database, v *visible) {
	s := v.sym
	ref := &ElementRef{Name: s.name, SchemaName: schemaPtr(s.schema), Filepath: s.unit.path}
	if s.kind == symTable {
		sch, name := in.sc.interpretedName(s)
		ref.VisibleNames = []VisibleName{{SchemaName: sch, Name: name}}
	} else {
		ref.SchemaName = nil
		ref.VisibleNames = []VisibleName{{Name: v.name}}
	}
	switch s.kind {
	case symTable:
		db.Externals.Tables = append(db.Externals.Tables, ref)
	case symEnum:
		db.Externals.Enums = append(db.Externals.Enums, ref)
	case symGroup:
		db.Externals.TableGroups = append(db.Externals.TableGroups, ref)
	case symPartial:
		db.Externals.TablePartials = append(db.Externals.TablePartials, ref)
	case symNote:
		db.Externals.Notes = append(db.Externals.Notes, ref)
	}
}

// ---------------------------------------------------------------- helpers

func noteFromAttr(a *ast.Attribute, normalize bool) *NoteValue {
	if a == nil {
		return nil
	}
	s, ok := extractQuotedString(a.Value)
	if !ok {
		return nil
	}
	if normalize {
		s = normalizeNote(s)
	}
	return &NoteValue{Value: s, Token: nodeToken(a, a.Base().Filepath)}
}

func noteFromSub(sub *ast.ElementDeclaration) (*NoteValue, bool) {
	var callee ast.Node
	switch b := sub.Body.(type) {
	case *ast.BlockExpression:
		if len(b.Body) > 0 {
			if f, ok := b.Body[0].(*ast.FunctionApplication); ok {
				callee = f.Callee
			}
		}
	case *ast.FunctionApplication:
		callee = b.Callee
	}
	s, ok := extractQuotedString(callee)
	if !ok {
		return nil, false
	}
	return &NoteValue{Value: normalizeNote(s), Token: nodeToken(sub, sub.Base().Filepath)}, true
}

func customMetadata(s *settings, builtins ...string) map[string]string {
	var out map[string]string
outer:
	for _, k := range s.keys {
		for _, b := range builtins {
			if k == b {
				continue outer
			}
		}
		if v, ok := extractMetadataValue(s.first(k).Value); ok {
			if out == nil {
				out = map[string]string{}
			}
			out[k] = v
		}
	}
	return out
}

func emptyToNil(m map[string]string) map[string]string {
	if len(m) == 0 {
		return nil
	}
	return m
}

var recognizedColumnSettings = []string{settingNote, settingRef, settingPrimaryKey, settingPK, settingNotNull, settingNull,
	settingUnique, settingIncrement, settingDefault, settingCheck}

func processColumnType(n ast.Node) (ColumnType, bool) {
	suffix := ""
	var args *string
	argText := func(e ast.Node) string {
		if isSignedNumber(e) {
			return numberText(e)
		}
		if s, ok := extractQuotedString(e); ok {
			return s
		}
		v, _ := extractVariable(e)
		return v
	}
	if call, ok := n.(*ast.CallExpression); ok && call != nil {
		var parts []string
		if call.ArgumentList != nil {
			for _, e := range call.ArgumentList.Elements {
				parts = append(parts, argText(e))
			}
		}
		a := strings.Join(parts, ",")
		args = &a
		suffix = "(" + a + ")"
		if ast.IsNil(call.Callee) {
			return ColumnType{TypeName: "", Args: args}, true
		}
		n = call.Callee
	}
loop:
	for {
		switch v := n.(type) {
		case *ast.CallExpression:
			var parts []string
			if v.ArgumentList != nil {
				for _, e := range v.ArgumentList.Elements {
					parts = append(parts, argText(e))
				}
			}
			suffix = "(" + strings.Join(parts, ",") + ")" + suffix
			if ast.IsNil(v.Callee) {
				break loop
			}
			n = v.Callee
		case *ast.Array:
			var parts []string
			if v.Indexer != nil {
				for _, a := range v.Indexer.Elements {
					val := ""
					if a != nil {
						if t := literalOf(a.Name); t != nil {
							val = t.Value
						}
					}
					parts = append(parts, val)
				}
			}
			suffix = "[" + strings.Join(parts, ",") + "]" + suffix
			if ast.IsNil(v.Array) {
				break loop
			}
			n = v.Array
		default:
			break loop
		}
	}
	schema, name := elementNameParts(n)
	ct := ColumnType{TypeName: name + suffix, Args: args}
	if len(schema) > 0 {
		ct.SchemaName = strPtr(schema[0])
	}
	return ct, len(schema) <= 1
}

func processDefault(n ast.Node) *Default {
	if ast.IsNil(n) {
		return nil
	}
	if s, ok := extractQuotedString(n); ok {
		return &Default{Type: "string", Value: s}
	}
	if isSignedNumber(n) {
		return &Default{Type: "number", Value: extractNumber(n)}
	}
	if isIdentifierNode(n) {
		return &Default{Type: "boolean", Value: strings.ToLower(variableOf(n).Value)}
	}
	if fe, ok := n.(*ast.FunctionExpression); ok && fe.Value != nil {
		return &Default{Type: "expression", Value: fe.Value.Value}
	}
	if isDotDelimitedIdentifier(n) {
		frags := destructureMemberAccess(n)
		if v, ok := extractVariable(frags[len(frags)-1]); ok {
			return &Default{Type: "string", Value: v}
		}
	}
	return nil
}

// ---------------------------------------------------------------- table

func (in *interpreter) table(s *sym) *Table {
	e := s.decl
	file := s.unit.path
	schema, name := in.sc.interpretedName(s)
	t := &Table{Name: name, SchemaName: schema, Fields: []*Column{}, Checks: []*Check{}, Partials: []*PartialInject{},
		Indexes: []*Index{}, Token: nodeToken(e, file)}
	if s.nested {
		in.n(diag.UNSUPPORTED, "Nested schema is not supported", e.Name, file)
	}
	if s.alias != "" {
		t.Alias = strPtr(s.alias)
	}
	set, _ := aggregateSettings(e.AttributeList, file)
	if a := set.first(settingHeaderColor); a != nil {
		if c, ok := extractColor(a.Value); ok {
			t.HeaderColor = &c
		}
	}
	t.Note = noteFromAttr(set.first(settingNote), true)
	t.Metadata = customMetadata(set, settingNote, settingHeaderColor)

	fields, subs := partitionBody(elementBody(e))
	var pkCols []*Column
	for order, f := range fields {
		if isValidPartialInjection(f.Callee) {
			pre := f.Callee.(*ast.PrefixExpression)
			pname, _ := extractVariable(pre.Expression)
			t.Partials = append(t.Partials, &PartialInject{Name: pname, Order: order, Token: nodeToken(pre, file)})
			continue
		}
		col := in.column(f, s, false)
		t.Fields = append(t.Fields, col)
		if col.PK != nil && *col.PK {
			pkCols = append(pkCols, col)
		}
	}
	if len(in.c.columns(s, s.unit.scope)) == 0 {
		in.n(diag.EMPTY_TABLE, "A Table must have at least one column", e, file)
	}
	for _, sub := range subs {
		switch {
		case sub.IsKind(kNote):
			if nv, ok := noteFromSub(sub); ok {
				t.Note = nv
			}
		case sub.IsKind(kIndexes):
			t.Indexes = append(t.Indexes, interpretIndexes(sub, file)...)
		case sub.IsKind(kChecks):
			t.Checks = append(t.Checks, interpretChecks(sub, file)...)
		}
	}
	t.Indexes = appendCompositePK(t.Indexes, pkCols, file)
	in.attachMetadata(s, nil, func(k, v string, tok *Token) bool {
		switch k {
		case settingNote:
			t.Note = &NoteValue{Value: v, Token: tok}
			return true
		case settingHeaderColor:
			t.HeaderColor = strPtr(v)
			return true
		}
		return false
	}, &t.Metadata)
	t.Metadata = emptyToNil(t.Metadata)
	// column-level metadata blocks
	for _, col := range t.Fields {
		cl := in.c.lookupColumn(s, s.unit.scope, col.Name)
		if cl == nil || cl.partial != nil {
			continue
		}
		c := col
		in.attachMetadata(s, cl, func(k, v string, tok *Token) bool {
			if k == settingNote {
				c.Note = &NoteValue{Value: v, Token: tok}
				return true
			}
			return false
		}, &c.Metadata)
		c.Metadata = emptyToNil(c.Metadata)
	}
	return t
}

func appendCompositePK(idx []*Index, pkCols []*Column, file string) []*Index {
	if len(pkCols) < 2 {
		return idx
	}
	neg := Position{Offset: -1, Line: -1, Column: -1}
	ix := &Index{Token: &Token{Start: neg, End: neg, Filepath: file}, PK: boolPtr(true)}
	for _, c := range pkCols {
		ix.Columns = append(ix.Columns, &IndexColumn{Value: c.Name, Type: "column", Token: c.Token})
		c.PK = boolPtr(false)
	}
	return append(idx, ix)
}

func (in *interpreter) column(f *ast.FunctionApplication, owner *sym, partial bool) *Column {
	file := owner.unit.path
	name, _ := extractVariable(f.Callee)
	col := &Column{Name: name, Token: nodeToken(f, file), InlineRefs: []*InlineRef{}, InlineDeps: []*InlineDep{}, Checks: []*Check{}}
	if len(f.Args) > 0 {
		ct, ok := processColumnType(f.Args[0])
		col.Type = ct
		if !ok {
			in.n(diag.UNSUPPORTED, "Nested schema is not supported", f.Args[0], file)
		}
	}
	var rest []ast.Node
	if len(f.Args) > 1 {
		rest = f.Args[1:]
	}
	set, errs := columnSettings(rest, file, partial)
	if len(errs) > 0 {
		set = newSettings()
	}
	pk := set.has(settingPK) || set.has(settingPrimaryKey)
	unique := set.has(settingUnique)
	if partial {
		if pk {
			col.PK = boolPtr(true)
		}
		if unique {
			col.Unique = boolPtr(true)
		}
	} else {
		col.PK = boolPtr(pk)
		col.Unique = boolPtr(unique)
	}
	if set.has(settingIncrement) {
		col.Increment = boolPtr(true)
	}
	if set.has(settingNotNull) {
		col.NotNull = boolPtr(true)
	} else if set.has(settingNull) {
		col.NotNull = boolPtr(false)
	}
	col.DBDefault = processDefault(set.firstValue(settingDefault))
	if nv := noteFromAttr(set.first(settingNote), true); nv != nil && nv.Value != "" {
		col.Note = nv
	}
	col.Metadata = emptyToNil(customMetadata(set, recognizedColumnSettings...))

	for _, a := range set.get(settingRef) {
		pre, ok := a.Value.(*ast.PrefixExpression)
		if !ok || pre.Op == nil || !isRelationshipOp(pre.Op.Value) {
			continue
		}
		res, _ := in.c.resolveEndpoint(in.sc, pre.Expression, modeInline, owner)
		if !res.ok || len(res.cols) == 0 {
			continue
		}
		schema, tname := in.sc.interpretedName(res.table.sym)
		col.InlineRefs = append(col.InlineRefs, &InlineRef{SchemaName: schema, TableName: tname, FieldNames: colNames(res.cols),
			Relation: pre.Op.Value, Token: nodeToken(a, file)})
	}
	if !partial {
		for _, a := range set.get(settingDep) {
			pre, ok := a.Value.(*ast.PrefixExpression)
			if !ok || pre.Op == nil || !isDependencyOp(pre.Op.Value) {
				continue
			}
			res, _ := in.c.resolveEndpoint(in.sc, pre.Expression, modeInline, owner)
			if !res.ok {
				continue
			}
			schema, tname := in.sc.interpretedName(res.table.sym)
			col.InlineDeps = append(col.InlineDeps, &InlineDep{SchemaName: schema, TableName: tname, FieldNames: colNames(res.cols),
				Direction: pre.Op.Value, Token: nodeToken(a, file)})
		}
	}
	for _, a := range set.get(settingCheck) {
		fe, ok := a.Value.(*ast.FunctionExpression)
		if !ok || fe.Value == nil {
			continue
		}
		col.Checks = append(col.Checks, &Check{Token: nodeToken(a, file), Expression: fe.Value.Value})
	}
	return col
}

func colNames(cols []*col) []string {
	out := make([]string, len(cols))
	for i, c := range cols {
		out[i] = c.name
	}
	return out
}

func interpretIndexes(e *ast.ElementDeclaration, file string) []*Index {
	var out []*Index
	b, ok := e.Body.(*ast.BlockExpression)
	if !ok {
		return nil
	}
	for _, it := range b.Body {
		f, ok := it.(*ast.FunctionApplication)
		if !ok || ast.IsNil(f.Callee) {
			continue
		}
		ix := &Index{Token: nodeToken(f, file), Columns: []*IndexColumn{}}
		args := append([]ast.Node{f.Callee}, f.Args...)
		if l := trailingList(args); l != nil {
			args = args[:len(args)-1]
			set, _ := aggregateSettings(l, file)
			if set.has(settingPK) {
				ix.PK = boolPtr(true)
			}
			if set.has(settingUnique) {
				ix.Unique = boolPtr(true)
			}
			if s, ok := extractQuotedString(set.firstValue(settingName)); ok {
				ix.Name = &s
			}
			if a := set.first(settingNote); a != nil {
				ix.Note = noteFromAttr(a, false)
			}
			if t, ok := extractVariable(set.firstValue(settingType)); ok {
				ix.Type = &t
			}
		}
		var flat []ast.Node
		for _, a := range args {
			call, ok := a.(*ast.CallExpression)
			if !ok {
				flat = append(flat, a)
				continue
			}
			var cur ast.Node = call
			for {
				ce, ok := cur.(*ast.CallExpression)
				if !ok {
					break
				}
				if ce.ArgumentList != nil {
					flat = append(flat, ce.ArgumentList)
				}
				if ast.IsNil(ce.Callee) {
					cur = nil
					break
				}
				cur = ce.Callee
			}
			if !ast.IsNil(cur) {
				flat = append(flat, cur)
			}
		}
		for _, a := range flat {
			parts, ok := destructureIndexNode(a)
			if !ok {
				continue
			}
			for _, fe := range parts.functional {
				v := ""
				if fe.Value != nil {
					v = fe.Value.Value
				}
				ix.Columns = append(ix.Columns, &IndexColumn{Value: v, Type: "expression", Token: nodeToken(fe, file)})
			}
			for _, n := range parts.nonFunctional {
				v, _ := extractVariable(n)
				ix.Columns = append(ix.Columns, &IndexColumn{Value: v, Type: "column", Token: nodeToken(n, file)})
			}
		}
		out = append(out, ix)
	}
	return out
}

func interpretChecks(e *ast.ElementDeclaration, file string) []*Check {
	var out []*Check
	b, ok := e.Body.(*ast.BlockExpression)
	if !ok {
		return nil
	}
	for _, it := range b.Body {
		f, ok := it.(*ast.FunctionApplication)
		if !ok {
			continue
		}
		c := &Check{Token: nodeToken(f, file)}
		if fe, ok := f.Callee.(*ast.FunctionExpression); ok && fe.Value != nil {
			c.Expression = fe.Value.Value
		}
		for _, a := range f.Args {
			if l, ok := a.(*ast.ListExpression); ok {
				set, _ := aggregateSettings(l, file)
				if s, ok := extractQuotedString(set.firstValue(settingName)); ok {
					c.Name = &s
				}
				break
			}
		}
		out = append(out, c)
	}
	return out
}

// ---------------------------------------------------------------- partial

func (in *interpreter) partial(s *sym) *TablePartial {
	e := s.decl
	file := s.unit.path
	_, name := in.sc.interpretedName(s)
	p := &TablePartial{Name: name, Fields: []*Column{}, Indexes: []*Index{}, Checks: []*Check{}, Token: nodeToken(e, file)}
	set, _ := aggregateSettings(e.AttributeList, file)
	if a := set.first(settingHeaderColor); a != nil {
		if c, ok := extractColor(a.Value); ok {
			p.HeaderColor = &c
		}
	}
	p.Note = noteFromAttr(set.first(settingNote), true)
	fields, subs := partitionBody(elementBody(e))
	var pkCols []*Column
	for _, f := range fields {
		col := in.column(f, s, true)
		p.Fields = append(p.Fields, col)
		if col.PK != nil && *col.PK {
			pkCols = append(pkCols, col)
		}
	}
	for _, sub := range subs {
		switch {
		case sub.IsKind(kNote):
			if nv, ok := noteFromSub(sub); ok {
				p.Note = nv
			}
		case sub.IsKind(kIndexes):
			p.Indexes = append(p.Indexes, interpretIndexes(sub, file)...)
		case sub.IsKind(kChecks):
			p.Checks = append(p.Checks, interpretChecks(sub, file)...)
		}
	}
	p.Indexes = appendCompositePK(p.Indexes, pkCols, file)
	return p
}

// ---------------------------------------------------------------- enum / group / note

func (in *interpreter) enum(s *sym) *Enum {
	file := s.unit.path
	schema, name := in.sc.interpretedName(s)
	en := &Enum{Name: name, SchemaName: schema, Token: nodeToken(s.decl, file), Values: []*EnumValue{}}
	if s.nested {
		in.n(diag.UNSUPPORTED, "Nested schema is not supported", s.decl.Name, file)
	}
	fields, _ := partitionBody(elementBody(s.decl))
	for _, f := range fields {
		name, _ := extractVariable(f.Callee)
		ev := &EnumValue{Name: name, Token: nodeToken(f, file)}
		if len(f.Args) > 0 {
			if l, ok := f.Args[0].(*ast.ListExpression); ok {
				set, _ := aggregateSettings(l, file)
				ev.Note = noteFromAttr(set.first(settingNote), true)
			}
		}
		en.Values = append(en.Values, ev)
	}
	return en
}

func (in *interpreter) tableGroup(s *sym) *TableGroup {
	e := s.decl
	file := s.unit.path
	g := &TableGroup{Tables: []*TableGroupItem{}, Token: nodeToken(e, file)}
	schemaParts, name := elementNameParts(e.Name)
	g.Name = strPtr(name)
	if len(schemaParts) > 0 {
		g.SchemaName = strPtr(schemaParts[0])
	}
	if len(schemaParts) >= 2 {
		in.n(diag.UNSUPPORTED, "Nested schema is not supported", e.Name, file)
	}
	set, _ := aggregateSettings(e.AttributeList, file)
	if a := set.first(settingColor); a != nil {
		if c, ok := extractColor(a.Value); ok {
			g.Color = &c
		}
	}
	g.Note = noteFromAttr(set.first(settingNote), true)
	g.Metadata = customMetadata(set, settingNote, settingColor)
	fields, subs := partitionBody(elementBody(e))
	for _, f := range fields {
		// Like @dbml/parse, only a bare name is resolved through the symbol
		// table (so aliases map to the real table); a dotted name such as
		// `public.users` is reported exactly as written.
		frags := destructureMemberAccess(f.Callee)
		res, _ := in.c.resolveChain(s.unit.scope, frags, nil, modeTable, nil)
		if len(frags) == 1 && res.ok && res.table != nil {
			sch, n := in.sc.interpretedName(res.table.sym)
			g.Tables = append(g.Tables, &TableGroupItem{Name: n, SchemaName: sch})
			continue
		}
		parts, _ := destructureComplexVariable(f.Callee)
		if len(parts) == 0 {
			continue
		}
		item := &TableGroupItem{Name: parts[len(parts)-1]}
		if len(parts) > 1 {
			item.SchemaName = strPtr(strings.Join(parts[:len(parts)-1], "."))
		}
		g.Tables = append(g.Tables, item)
	}
	for _, sub := range subs {
		if sub.IsKind(kNote) {
			if nv, ok := noteFromSub(sub); ok {
				g.Note = nv
			}
		}
	}
	in.attachMetadata(s, nil, func(k, v string, tok *Token) bool {
		switch k {
		case settingNote:
			g.Note = &NoteValue{Value: v, Token: tok}
			return true
		case settingColor:
			g.Color = strPtr(v)
			return true
		}
		return false
	}, &g.Metadata)
	g.Metadata = emptyToNil(g.Metadata)
	return g
}

func (in *interpreter) stickyNote(s *sym) *StickyNote {
	e := s.decl
	file := s.unit.path
	_, name := in.sc.interpretedName(s)
	nt := &StickyNote{Name: name, Token: nodeToken(e, file)}
	set, _ := aggregateSettings(e.AttributeList, file)
	if a := set.first(settingColor); a != nil {
		if c, ok := extractColor(a.Value); ok {
			nt.Color = &c
		}
	}
	nt.Metadata = customMetadata(set, settingColor)
	fields, subs := partitionBody(elementBody(e))
	if b, ok := e.Body.(*ast.BlockExpression); ok && (len(fields) != 1 || len(subs) > 0) {
		in.n(diag.INVALID_NOTE, "Invalid note syntax", b, file)
	} else if len(fields) == 1 {
		content, _ := extractQuotedString(fields[0].Callee)
		nt.Content = normalizeNote(content)
	}
	in.attachMetadata(s, nil, func(k, v string, tok *Token) bool {
		if k == settingColor {
			nt.Color = strPtr(v)
			return true
		}
		return false
	}, &nt.Metadata)
	nt.Metadata = emptyToNil(nt.Metadata)
	return nt
}

// ---------------------------------------------------------------- metadata blocks

// attachMetadata merges every Metadata block targeting (s, col) into the
// element. builtin assigns built-in keys (note, color, ...) and reports
// whether it handled the key. Precedence (lowest to highest): inline
// settings, Metadata blocks in imported files (later imports win), Metadata
// blocks in the current file (later blocks win).
func (in *interpreter) attachMetadata(s *sym, cl *col, builtin func(key, value string, tok *Token) bool, meta *map[string]string) {
	type entry struct{ orig, value string }
	merged := map[string]entry{}
	var order []string
	if *meta != nil {
		for _, k := range sortedKeys(*meta) {
			merged[strings.ToLower(k)] = entry{k, (*meta)[k]}
			order = append(order, strings.ToLower(k))
		}
	}
	found := false
	for _, mb := range in.metadataFor(s, cl) {
		e, f := mb.decl, mb.unit
		{
			found = true
			_, subs := partitionBody(elementBody(e))
			for _, sub := range subs {
				if sub.Type == nil {
					continue
				}
				fa, ok := sub.Body.(*ast.FunctionApplication)
				if !ok {
					continue
				}
				v, ok := extractMetadataValue(fa.Callee)
				if !ok || v == "" {
					continue
				}
				k := sub.Type.Value
				lk := strings.ToLower(k)
				if lk == settingNote {
					v = normalizeNote(v)
				}
				if builtin(lk, v, nodeToken(sub, f.path)) {
					continue
				}
				if _, exists := merged[lk]; !exists {
					order = append(order, lk)
				}
				merged[lk] = entry{k, v}
			}
		}
	}
	if !found && *meta == nil {
		return
	}
	out := map[string]string{}
	for _, lk := range order {
		e := merged[lk]
		out[e.orig] = e.value
	}
	*meta = out
}

// ---------------------------------------------------------------- refs / deps / records / project

func (in *interpreter) interpretMetadata(db *Database) {
	seenRefs := map[string]bool{}
	seenDeps := map[string]bool{}
	recordCount := map[*sym]int{}
	type recordDecl struct {
		e     *ast.ElementDeclaration
		owner *sym
		unit  *unit
	}
	var records []recordDecl
	var blocks []tableRecords

	addRef := func(r *Ref, errNode ast.Node, file string) {
		if r == nil {
			return
		}
		ep1, ep2 := r.Endpoints[0], r.Endpoints[1]
		k := func(a, b RefEndpoint) string {
			return strings.Join([]string{str(a.SchemaName), a.TableName, strings.Join(a.FieldNames, ","), str(b.SchemaName), b.TableName, strings.Join(b.FieldNames, ",")}, "|")
		}
		if seenRefs[k(ep1, ep2)] || seenRefs[k(ep2, ep1)] {
			in.n(diag.SAME_ENDPOINT, "Ref with same endpoints already exists", errNode, file)
			return
		}
		seenRefs[k(ep1, ep2)] = true
		db.Refs = append(db.Refs, r)
	}
	addDep := func(d *Dep, edgeNodes []ast.Node, errNode ast.Node, block bool, file string) {
		if d == nil {
			return
		}
		for i, e := range d.Edges {
			key := strings.Join([]string{str(e.Upstream.SchemaName), e.Upstream.TableName, strings.Join(e.Upstream.FieldNames, ","),
				str(e.Downstream.SchemaName), e.Downstream.TableName, strings.Join(e.Downstream.FieldNames, ",")}, "|")
			if seenDeps[key] {
				target := errNode
				if i < len(edgeNodes) {
					target = edgeNodes[i]
				}
				in.n(diag.SAME_ENDPOINT, "Dep with same endpoints already exists", target, file)
				return
			}
		}
		for _, e := range d.Edges {
			key := strings.Join([]string{str(e.Upstream.SchemaName), e.Upstream.TableName, strings.Join(e.Upstream.FieldNames, ","),
				str(e.Downstream.SchemaName), e.Downstream.TableName, strings.Join(e.Downstream.FieldNames, ",")}, "|")
			seenDeps[key] = true
		}
		if block && len(d.Edges) > 0 {
			first := tableKey(d.Edges[0].Downstream.SchemaName, d.Edges[0].Downstream.TableName)
			for _, e := range d.Edges[1:] {
				cur := tableKey(e.Downstream.SchemaName, e.Downstream.TableName)
				if cur != first {
					in.n(diag.DEP_MIXED_DOWNSTREAM_TABLES, "All edges in a Dep block must target the same downstream table, but found '"+first+"' and '"+cur+"'", errNode, file)
					return
				}
			}
			lvl := func(e *DepEdge) bool { return len(e.Upstream.FieldNames) > 0 || len(e.Downstream.FieldNames) > 0 }
			for _, e := range d.Edges[1:] {
				if lvl(e) != lvl(d.Edges[0]) {
					in.n(diag.DEP_MIXED_LEVEL, "All edges in a Dep block must be at the same level (all table-level or all column-level)", errNode, file)
					return
				}
			}
		}
		db.Deps = append(db.Deps, d)
	}

	for _, u := range in.c.reachable(in.u) {
		in.rs, in.rfile = u.scope, u.path
		for _, e := range u.prog.Declarations() {
			in.c.poll(false)
			switch {
			case e.IsKind(kRef):
				addRef(in.standaloneRef(e), e, u.path)
			case e.IsKind(kDep):
				d, edges := in.standaloneDep(e)
				_, isBlock := e.Body.(*ast.BlockExpression)
				addDep(d, edges, e, isBlock, u.path)
			case e.IsKind(kTable):
				t := u.symOfDecl(e)
				if t == nil || !in.visibleSym(t) {
					continue
				}
				if e.AttributeList != nil {
					set, _ := aggregateSettings(e.AttributeList, u.path)
					for _, a := range set.get(settingDep) {
						addDep(in.inlineDep(a, t, ""), nil, a, false, u.path)
					}
				}
				fields, subs := partitionBody(elementBody(e))
				for _, f := range fields {
					cname, ok := extractVariable(f.Callee)
					if !ok {
						continue
					}
					l := trailingList(f.Args)
					if l == nil {
						continue
					}
					set, _ := aggregateSettings(l, u.path)
					for _, a := range set.get(settingRef) {
						addRef(in.inlineRef(a, t, f, cname), a, u.path)
					}
					for _, a := range set.get(settingDep) {
						addDep(in.inlineDep(a, t, cname), nil, a, false, u.path)
					}
				}
				for _, sub := range subs {
					if sub.IsKind(kRecords) {
						records = append(records, recordDecl{sub, t, u})
						recordCount[t]++
					}
				}
			case e.IsKind(kRecords):
				if tv, _, _ := in.c.recordsTarget(in.rs, e, nil); tv != nil {
					recordCount[tv.sym]++
				}
				records = append(records, recordDecl{e, nil, u})
			case e.IsKind(kProject):
				if u == in.u {
					db.Project = in.project(e)
				}
			}
		}
	}
	for _, r := range records {
		in.rs, in.rfile = r.unit.scope, r.unit.path
		tv, _, _ := in.c.recordsTarget(in.rs, r.e, r.owner)
		if tv != nil && recordCount[tv.sym] > 1 {
			in.n(diag.DUPLICATE_RECORDS_FOR_TABLE, "Duplicate Records blocks for the same Table '"+tv.sym.name+"' - A Table can only have one Records block", r.e, r.unit.path)
			continue
		}
		if rec, t := in.records(r.e, r.owner); rec != nil {
			db.Records = append(db.Records, rec)
			blocks = append(blocks, tableRecords{sym: t, rec: rec, block: r.e, file: r.unit.path})
		}
	}
	in.validateRecordConstraints(db, blocks)
}

func str(p *string) string {
	if p == nil {
		return "null"
	}
	return *p
}

func tableKey(schema *string, table string) string {
	if schema == nil {
		return table
	}
	return *schema + "." + table
}

func (in *interpreter) endpoint(res resolved, rel string, tok *Token) RefEndpoint {
	ep := RefEndpoint{FieldNames: []string{}, Relation: rel, Token: tok}
	if res.table != nil {
		ep.SchemaName, ep.TableName = in.sc.interpretedName(res.table.sym)
	}
	if res.cols != nil {
		ep.FieldNames = colNames(res.cols)
	}
	return ep
}

func (in *interpreter) standaloneRef(e *ast.ElementDeclaration) *Ref {
	file := in.rfile
	in2 := refInfix(e)
	if in2 == nil || in2.Op == nil {
		return nil
	}
	mult, ok := relationshipOps[in2.Op.Value]
	if !ok {
		return nil
	}
	l, _ := in.c.resolveEndpoint(in.rs, in2.Left, modeRef, nil)
	r, _ := in.c.resolveEndpoint(in.rs, in2.Right, modeRef, nil)
	if !l.ok || !r.ok || !in.allVisible(l, r) {
		return nil
	}
	if sameCols(l, r) {
		in.n(diag.SAME_ENDPOINT, "Two endpoints are the same", e, file)
		return nil
	}
	ref := &Ref{Token: nodeToken(e, file)}
	if parts, ok := destructureComplexVariable(e.Name); ok && len(parts) > 0 {
		ref.Name = strPtr(parts[len(parts)-1])
		if len(parts) > 1 {
			ref.SchemaName = strPtr(strings.Join(parts[:len(parts)-1], "."))
		}
		if len(parts) > 2 {
			in.n(diag.UNSUPPORTED, "Nested schema is not supported", e.Name, file)
		}
	}
	ref.Endpoints[0] = in.endpoint(l, mult[0], nodeToken(in2.Left, file))
	ref.Endpoints[1] = in.endpoint(r, mult[1], nodeToken(in2.Right, file))
	field := elementBody(e)[0].(*ast.FunctionApplication)
	if len(field.Args) > 0 {
		if lst, ok := field.Args[0].(*ast.ListExpression); ok {
			set, _ := aggregateSettings(lst, file)
			ref.OnDelete = policyValue(set.firstValue(settingDelete))
			ref.OnUpdate = policyValue(set.firstValue(settingUpdate))
			if a := set.first(settingColor); a != nil {
				if c, ok := extractColor(a.Value); ok {
					ref.Color = &c
				}
			}
			if set.has(settingInactive) {
				ref.Inactive = boolPtr(true)
			}
		}
	}
	return ref
}

func policyValue(n ast.Node) *string {
	if is, ok := n.(*ast.IdentifierStream); ok && is != nil {
		if s := is.String(); s != "" {
			return &s
		}
		return nil
	}
	if v, ok := extractVariable(n); ok {
		return &v
	}
	return nil
}

func sameCols(l, r resolved) bool {
	if l.table == nil || r.table == nil || l.table.sym != r.table.sym {
		return false
	}
	n := len(l.cols)
	if len(r.cols) > n {
		n = len(r.cols)
	}
	for i := 0; i < n; i++ {
		if i >= len(l.cols) || i >= len(r.cols) || l.cols[i].name != r.cols[i].name {
			return false
		}
	}
	return true
}

func (in *interpreter) inlineRef(a *ast.Attribute, t *sym, f *ast.FunctionApplication, cname string) *Ref {
	file := in.rfile
	pre, ok := a.Value.(*ast.PrefixExpression)
	if !ok || pre.Op == nil {
		return nil
	}
	mult, ok := relationshipOps[pre.Op.Value]
	if !ok {
		return nil
	}
	r, _ := in.c.resolveEndpoint(in.rs, pre.Expression, modeInline, t)
	if !r.ok || !in.allVisible(r) {
		return nil
	}
	cl := in.c.lookupColumn(t, t.unit.scope, cname)
	if cl == nil {
		return nil
	}
	l := resolved{table: in.sc.visibleOf(t), cols: []*col{cl}, ok: true}
	if sameCols(l, r) {
		in.n(diag.SAME_ENDPOINT, "Two endpoints are the same", a, file)
		return nil
	}
	ref := &Ref{Token: nodeToken(a, file)}
	ref.Endpoints[0] = in.endpoint(r, mult[1], nodeToken(a, file))
	ref.Endpoints[1] = in.endpoint(l, mult[0], nodeToken(f, file))
	return ref
}

func (in *interpreter) inlineDep(a *ast.Attribute, t *sym, cname string) *Dep {
	file := in.rfile
	pre, ok := a.Value.(*ast.PrefixExpression)
	if !ok || pre.Op == nil || !isDependencyOp(pre.Op.Value) {
		return nil
	}
	mode := modeDep
	if cname != "" {
		mode = modeInline
	}
	other, _ := in.c.resolveEndpoint(in.rs, pre.Expression, mode, t)
	if !other.ok || !in.allVisible(other) {
		return nil
	}
	host := resolved{table: in.sc.visibleOf(t), ok: true}
	if cname != "" {
		if cl := in.c.lookupColumn(t, t.unit.scope, cname); cl != nil {
			host.cols = []*col{cl}
		}
	}
	tok := nodeToken(a, file)
	d := &Dep{Token: tok}
	up, down := host, other
	if pre.Op.Value == depUpstream {
		up, down = other, host
	}
	d.Edges = []*DepEdge{{Upstream: in.depEndpoint(up, tok), Downstream: in.depEndpoint(down, tok), Token: tok}}
	return d
}

func (in *interpreter) depEndpoint(r resolved, tok *Token) DepEndpoint {
	ep := DepEndpoint{FieldNames: []string{}, Token: tok}
	if r.table != nil {
		ep.SchemaName, ep.TableName = in.sc.interpretedName(r.table.sym)
	}
	if len(r.cols) > 0 {
		ep.FieldNames = colNames(r.cols)
	}
	return ep
}

func (in *interpreter) standaloneDep(e *ast.ElementDeclaration) (*Dep, []ast.Node) {
	file := in.rfile
	tok := nodeToken(e, file)
	d := &Dep{Token: tok, Edges: []*DepEdge{}}
	if parts, ok := destructureComplexVariable(e.Name); ok && len(parts) > 0 {
		d.Name = strPtr(parts[len(parts)-1])
		if len(parts) > 1 {
			d.SchemaName = strPtr(strings.Join(parts[:len(parts)-1], "."))
		}
	}
	fields, subs := partitionBody(elementBody(e))
	var edgeNodes []ast.Node
	for _, f := range fields {
		inx, ok := f.Callee.(*ast.InfixExpression)
		if !ok || inx.Op == nil {
			continue
		}
		edgeNodes = append(edgeNodes, inx)
		upN, downN := inx.Left, inx.Right
		if inx.Op.Value == depUpstream {
			upN, downN = inx.Right, inx.Left
		}
		up, _ := in.c.resolveEndpoint(in.rs, upN, modeDep, nil)
		down, _ := in.c.resolveEndpoint(in.rs, downN, modeDep, nil)
		if !in.allVisible(up, down) {
			return nil, nil
		}
		d.Edges = append(d.Edges, &DepEdge{Upstream: in.depEndpoint(up, tok), Downstream: in.depEndpoint(down, tok), Token: tok})
	}
	meta := map[string]any{}
	consume := func(l *ast.ListExpression) {
		set, _ := aggregateSettings(l, file)
		for _, name := range set.keys {
			a := set.first(name)
			if name == settingColor {
				if c, ok := extractColor(a.Value); ok {
					d.Color = &c
				}
				continue
			}
			v, ok := depAttrValue(a.Value)
			if !ok {
				continue
			}
			if name == settingNote {
				d.Note = &NoteValue{Value: normalizeNote(anyString(v)), Token: nodeToken(a, file)}
			} else {
				meta[name] = v
			}
		}
	}
	if e.AttributeList != nil {
		consume(e.AttributeList)
	}
	for _, f := range fields {
		for _, a := range f.Args {
			if l, ok := a.(*ast.ListExpression); ok {
				consume(l)
				break
			}
		}
	}
	if _, isBlock := e.Body.(*ast.BlockExpression); isBlock {
		for _, sub := range subs {
			if sub.Type == nil {
				continue
			}
			k := strings.ToLower(sub.Type.Value)
			if sub.IsKind(kNote) {
				if b, ok := sub.Body.(*ast.BlockExpression); ok {
					if len(b.Body) > 0 {
						if inner, ok := b.Body[0].(*ast.FunctionApplication); ok && !ast.IsNil(inner.Callee) {
							if v, ok := depAttrValue(inner.Callee); ok {
								d.Note = &NoteValue{Value: normalizeNote(anyString(v)), Token: nodeToken(sub, file)}
							}
						}
					}
					continue
				}
			}
			fa, ok := sub.Body.(*ast.FunctionApplication)
			if !ok {
				continue
			}
			if k == settingColor {
				if c, ok := extractColor(fa.Callee); ok {
					d.Color = &c
				}
				continue
			}
			v, ok := depAttrValue(fa.Callee)
			if !ok {
				continue
			}
			if sub.IsKind(kNote) {
				d.Note = &NoteValue{Value: normalizeNote(anyString(v)), Token: nodeToken(sub, file)}
			} else {
				meta[k] = v
			}
		}
	}
	if len(meta) > 0 {
		d.Metadata = meta
	}
	return d, edgeNodes
}

func depAttrValue(n ast.Node) (any, bool) {
	if is, ok := n.(*ast.IdentifierStream); ok && is != nil {
		s := is.String()
		if s == "" {
			return nil, false
		}
		return s, true
	}
	if v, ok := extractNumericLiteral(n); ok {
		return v, true
	}
	if s, ok := extractQuotedString(n); ok {
		return s, true
	}
	if v, ok := extractVariable(n); ok {
		return v, true
	}
	return nil, false
}

func anyString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return formatJSNumber(t)
	}
	return ""
}

func (in *interpreter) project(e *ast.ElementDeclaration) *Project {
	file := in.u.path
	p := &Project{Tables: []*Table{}, Refs: []*Ref{}, Enums: []*Enum{}, TableGroups: []*TableGroup{}, TablePartials: []*TablePartial{},
		Token: nodeToken(e, file)}
	if !ast.IsNil(e.Name) {
		_, name := elementNameParts(e.Name)
		p.Name = strPtr(name)
	}
	_, subs := partitionBody(elementBody(e))
	for _, sub := range subs {
		switch {
		case sub.IsKind(kTable), sub.IsKind(kTableGroup), sub.IsKind(kEnum), sub.IsKind(kTablePartial):
			// Elements nested in a Project are not symbols (they are rejected by
			// validation); nothing to interpret.
		case sub.IsKind(kNote):
			if nv, ok := noteFromSub(sub); ok && nv.Value != "" {
				p.Note = nv
			} else if fa, ok := sub.Body.(*ast.FunctionApplication); ok {
				if s, ok := extractQuotedString(fa.Callee); ok && s != "" {
					p.Note = &NoteValue{Value: normalizeNote(s), Token: nodeToken(sub, file)}
				}
			}
		default:
			if sub.Type == nil {
				continue
			}
			if p.Extra == nil {
				p.Extra = map[string]string{}
			}
			if fa, ok := sub.Body.(*ast.FunctionApplication); ok {
				if s, ok := extractQuotedString(fa.Callee); ok {
					p.Extra[strings.ToLower(sub.Type.Value)] = s
				}
			}
		}
	}
	return p
}

// ---------------------------------------------------------------- diagram view

func (in *interpreter) diagramView(s *sym) *DiagramView {
	e := s.decl
	file := s.unit.path
	schema, name := in.sc.interpretedName(s)
	dv := &DiagramView{Name: name, SchemaName: schema, Token: nodeToken(e, file)}
	ve := &dv.VisibleEntities
	b, ok := e.Body.(*ast.BlockExpression)
	if !ok {
		return dv
	}
	wild := map[string]bool{}
	if len(b.Body) == 1 {
		if f, ok := b.Body[0].(*ast.FunctionApplication); ok && isWildcard(f.Callee) {
			ve.Tables, ve.StickyNotes, ve.TableGroups, ve.Schemas = []DiagramViewTable{}, []NamedItem{}, []NamedItem{}, []NamedItem{}
			return dv
		}
	}
	_, subs := partitionBody(b.Body)
	explicit := map[string]bool{}
	for _, sub := range subs {
		if sub.Type == nil {
			continue
		}
		bt := strings.ToLower(sub.Type.Value)
		explicit[bt] = true
		sb, ok := sub.Body.(*ast.BlockExpression)
		if !ok {
			continue
		}
		hasWild := false
		for _, it := range sb.Body {
			if f, ok := it.(*ast.FunctionApplication); ok && isWildcard(f.Callee) {
				hasWild = true
			}
		}
		if hasWild {
			switch bt {
			case kDVTables:
				ve.Tables = []DiagramViewTable{}
			case kDVNotes:
				ve.StickyNotes = []NamedItem{}
			case kDVTableGroups:
				ve.TableGroups = []NamedItem{}
			case kDVSchemas:
				ve.Schemas = []NamedItem{}
			}
			wild[bt] = true
			continue
		}
		if len(sb.Body) == 0 {
			continue
		}
		var items []DiagramViewTable
		for _, it := range sb.Body {
			f, ok := it.(*ast.FunctionApplication)
			if !ok {
				continue
			}
			if bt == kDVTables {
				if rv := getRightmostVariable(f.Callee); rv != nil {
					parts, _ := destructureComplexVariable(f.Callee)
					schema := ""
					if len(parts) > 1 {
						schema = strings.Join(parts[:len(parts)-1], ".")
					}
					name, _ := extractVariable(rv)
					if t := in.sc.lookup(symTable, schema, name); t != nil {
						items = append(items, DiagramViewTable{Name: t.sym.name, SchemaName: t.sym.schema})
						continue
					}
				}
			}
			parts, ok := destructureComplexVariable(f.Callee)
			if !ok || len(parts) == 0 {
				continue
			}
			item := DiagramViewTable{Name: parts[len(parts)-1], SchemaName: defaultSchema}
			if len(parts) > 1 {
				item.SchemaName = parts[0]
			}
			items = append(items, item)
		}
		named := func() []NamedItem {
			if len(items) == 0 {
				return nil
			}
			out := make([]NamedItem, len(items))
			for i, it := range items {
				out[i] = NamedItem{Name: it.Name}
			}
			return out
		}
		switch bt {
		case kDVTables:
			if len(items) > 0 {
				ve.Tables = items
			}
		case kDVNotes:
			ve.StickyNotes = named()
		case kDVTableGroups:
			ve.TableGroups = named()
		case kDVSchemas:
			ve.Schemas = named()
		}
	}
	nonNull := (explicit[kDVTables] && ve.Tables != nil) || (explicit[kDVTableGroups] && ve.TableGroups != nil) || (explicit[kDVSchemas] && ve.Schemas != nil)
	if nonNull {
		if !explicit[kDVTables] {
			ve.Tables = []DiagramViewTable{}
		}
		if !explicit[kDVTableGroups] {
			ve.TableGroups = []NamedItem{}
		}
		if !explicit[kDVSchemas] {
			ve.Schemas = []NamedItem{}
		}
	}
	if wild[kDVTables] || wild[kDVSchemas] {
		ve.Tables, ve.TableGroups, ve.Schemas = []DiagramViewTable{}, []NamedItem{}, []NamedItem{}
		return dv
	}
	if wild[kDVTableGroups] && ve.TableGroups != nil && len(ve.TableGroups) == 0 {
		seen := map[*sym]bool{}
		ve.TableGroups = []NamedItem{}
		for _, v := range in.sc.entries {
			if v.sym.kind == symGroup && !seen[v.sym] {
				seen[v.sym] = true
				ve.TableGroups = append(ve.TableGroups, NamedItem{Name: v.name})
			}
		}
	}
	return dv
}

// deriveGenerated sets Column.Generated from the generated /
// generated_stored custom properties (after Metadata blocks applied).
func deriveGenerated(db *Database) {
	for _, t := range db.Tables {
		for _, c := range t.Fields {
			c.Generated = GeneratedFromMetadata(c.Metadata)
		}
	}
	for _, p := range db.TablePartials {
		for _, c := range p.Fields {
			c.Generated = GeneratedFromMetadata(c.Metadata)
		}
	}
}
