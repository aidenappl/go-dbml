// Portions of this file are derived from holistics/dbml packages/dbml-parse/src/core/global_modules/*/{bind,index}.ts, src/core/global_modules/program/bind.ts, src/core/global_modules/use/index.ts
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

package model

import (
	"sort"
	"strings"

	"github.com/aidenappl/go-dbml/ast"
	"github.com/aidenappl/go-dbml/diag"
)

// ---------------------------------------------------------------- resolution

// chainMode selects how the first fragment of a dotted reference resolves.
type chainMode int

const (
	modeRef chainMode = iota
	modeInline
	modeDep
	modeTable
)

// resolved is the outcome of resolving `schema.table.column` style references.
type resolved struct {
	table *visible
	cols  []*col
	// ok is false when any fragment failed to resolve.
	ok bool
}

type resolveErr struct {
	node ast.Node
	msg  string
}

func tableFullname(s *sym) string { return s.schema + "." + s.name }

// resolveChain resolves frags (variables) plus optional tuple column names.
// container is the enclosing table/partial for bare column references.
func (c *compiler) resolveChain(sc *scope, frags, tuple []ast.Node, mode chainMode, container *sym) (resolved, []resolveErr) {
	var errs []resolveErr
	fail := func(n ast.Node, msg string) (resolved, []resolveErr) {
		return resolved{}, append(errs, resolveErr{n, msg})
	}
	if len(frags) == 0 {
		return resolved{}, nil
	}
	total := len(frags)
	if len(tuple) > 0 {
		total++
	}
	name0, _ := extractVariable(frags[0])

	// Bare `col` in an inline ref / column dep: column of the container.
	if (mode == modeInline) && len(frags) == 1 && len(tuple) == 0 {
		if container == nil {
			return resolved{}, nil
		}
		cl := c.containerColumn(container, sc, name0)
		if cl == nil {
			kind := "Table"
			if container.kind == symPartial {
				return fail(frags[0], "Column '"+name0+"' does not exist in TablePartial '"+container.name+"'")
			}
			return fail(frags[0], "Column '"+name0+"' does not exist in "+kind+" 'public."+container.name+"'")
		}
		return resolved{table: sc.visibleOf(container), cols: []*col{cl}, ok: true}, nil
	}

	var curSchema string
	var curTable *visible
	i := 0
	switch {
	case mode == modeTable:
		// table references: `t` or `s.t` or `s.x.t`
		if len(frags) >= 2 {
			if !sc.hasSchema(name0) {
				return fail(frags[0], "Schema '"+name0+"' does not exist")
			}
			curSchema = name0
		} else {
			curTable = sc.lookup(symTable, "", name0)
			if curTable == nil {
				return fail(frags[0], "Table '"+name0+"' does not exist in Schema 'public'")
			}
		}
	case mode == modeDep && total <= 2 && len(frags) == 1 && len(tuple) == 0:
		curTable = sc.lookup(symTable, "", name0)
		if curTable == nil {
			return fail(frags[0], "Table '"+name0+"' does not exist in Schema 'public'")
		}
	case mode == modeDep && total == 2:
		if sc.hasSchema(name0) {
			curSchema = name0
		} else {
			curTable = sc.lookup(symTable, "", name0)
			if curTable == nil {
				return fail(frags[0], "Schema or Table '"+name0+"' does not exist")
			}
		}
	case total >= 3:
		if !sc.hasSchema(name0) {
			if mode == modeRef {
				return fail(frags[0], "Schema '"+name0+"' does not exist")
			}
			return fail(frags[0], "Schema '"+name0+"' does not exist in Schema 'public'")
		}
		curSchema = name0
	default:
		curTable = sc.lookup(symTable, "", name0)
		if curTable == nil {
			return fail(frags[0], "Table '"+name0+"' does not exist in Schema 'public'")
		}
	}
	i = 1
	var cols []*col
	for ; i < len(frags); i++ {
		name, _ := extractVariable(frags[i])
		switch {
		case curTable != nil:
			if len(cols) > 0 {
				return resolved{}, errs
			}
			cl := c.lookupColumn(curTable.sym, curTable.sym.unit.scope, name)
			if cl == nil {
				return fail(frags[i], "Column '"+name+"' does not exist in Table '"+tableFullname(curTable.sym)+"'")
			}
			cols = append(cols, cl)
		case curSchema != "":
			if mode == modeInline && i == len(frags)-1 && len(tuple) == 0 {
				// terminal fragment of an inline ref must be a column
				return fail(frags[i], "Column '"+name+"' does not exist")
			}
			if t := sc.lookup(symTable, curSchema, name); t != nil {
				curTable = t
				curSchema = ""
			} else if mode != modeDep && sc.hasSchema(curSchema+"."+name) {
				curSchema = curSchema + "." + name
			} else {
				switch mode {
				case modeDep:
					return fail(frags[i], "Table '"+name+"' does not exist in Schema '"+lastSeg(curSchema)+"'")
				case modeTable:
					return fail(frags[i], "Table '"+name+"' does not exist in Schema '"+curSchema+"'")
				}
				return fail(frags[i], "Table or schema '"+name+"' does not exist")
			}
		}
	}
	if curTable == nil {
		return resolved{}, errs
	}
	for _, t := range tuple {
		name, _ := extractVariable(t)
		cl := c.lookupColumn(curTable.sym, curTable.sym.unit.scope, name)
		if cl == nil {
			errs = append(errs, resolveErr{t, "Column '" + name + "' does not exist in Table '" + tableFullname(curTable.sym) + "'"})
			continue
		}
		cols = append(cols, cl)
	}
	if len(errs) > 0 {
		return resolved{}, errs
	}
	return resolved{table: curTable, cols: cols, ok: true}, nil
}

func mustVar(frags []ast.Node, i int) string {
	if i < len(frags) {
		v, _ := extractVariable(frags[i])
		return v
	}
	return ""
}

func lastSeg(s string) string {
	if i := strings.LastIndex(s, "."); i >= 0 {
		return s[i+1:]
	}
	return s
}

func (c *compiler) containerColumn(container *sym, sc *scope, name string) *col {
	if container.kind == symPartial {
		return c.lookupColumn(container, container.unit.scope, name)
	}
	return c.lookupColumn(container, sc, name)
}

// resolveEndpoint resolves a ref/dep endpoint expression.
func (c *compiler) resolveEndpoint(sc *scope, expr ast.Node, mode chainMode, container *sym) (resolved, []resolveErr) {
	t, ok := destructureComplexVariableTuple(expr)
	if !ok {
		return resolved{}, nil
	}
	return c.resolveChain(sc, t.variables, t.tupleElements, mode, container)
}

// symOfDecl finds the symbol declared by e in unit u.
func (u *unit) symOfDecl(e *ast.ElementDeclaration) *sym {
	if u.declIndex == nil {
		u.declIndex = make(map[*ast.ElementDeclaration]*sym, len(u.decls))
		for _, s := range u.decls {
			if _, ok := u.declIndex[s.decl]; !ok {
				u.declIndex[s.decl] = s
			}
		}
	}
	return u.declIndex[e]
}

// ---------------------------------------------------------------- binder

type binder struct {
	c     *compiler
	u     *unit
	diags []diag.Diagnostic
}

func (b *binder) n(code diag.Code, msg string, node ast.Node) {
	if ast.IsNil(node) {
		return
	}
	b.diags = append(b.diags, nodeDiag(code, msg, node, b.u.path))
}

func (b *binder) report(errs []resolveErr) {
	for _, e := range errs {
		b.n(diag.BINDING_ERROR, e.msg, e.node)
	}
}

func (b *binder) bind() {
	sc := b.u.scope
	var projects []*ast.ElementDeclaration
	for _, e := range b.u.prog.Declarations() {
		b.c.poll(false)
		if e.Type == nil {
			continue
		}
		switch {
		case e.IsKind(kProject):
			projects = append(projects, e)
			b.bindProject(e)
		case e.IsKind(kTable):
			b.bindTable(e, b.u.symOfDecl(e))
		case e.IsKind(kTablePartial):
			b.bindPartial(e, b.u.symOfDecl(e))
		case e.IsKind(kRef):
			b.bindRef(e)
		case e.IsKind(kDep):
			b.bindDep(e)
		case e.IsKind(kTableGroup):
			b.bindTableGroup(e)
		case e.IsKind(kEnum):
			b.bindEnum(e)
		case e.IsKind(kRecords):
			b.bindRecords(e, nil)
		case e.IsKind(kMetadata):
			b.bindMetadata(e)
		case e.IsKind(kDiagramView):
			b.bindDiagramView(e)
		}
	}
	if len(projects) > 1 {
		for _, p := range projects {
			b.n(diag.PROJECT_REDEFINED, "Only one project can exist", p)
		}
	}
	b.bindUses()
	b.checkDuplicateMembers(sc)
	b.checkDuplicateRefs()
	b.checkTableGroups()
}

func (b *binder) bindProject(e *ast.ElementDeclaration) {
	_, subs := partitionBody(elementBody(e))
	for _, s := range subs {
		switch {
		case s.IsKind(kTable):
			b.bindTable(s, nil)
		case s.IsKind(kRef):
			b.bindRef(s)
		case s.IsKind(kDep):
			b.bindDep(s)
		}
	}
}

func (b *binder) bindTable(e *ast.ElementDeclaration, t *sym) {
	sc := b.u.scope
	if e.AttributeList != nil {
		s, _ := aggregateSettings(e.AttributeList, b.u.path)
		for _, a := range s.get(settingDep) {
			if pre, ok := a.Value.(*ast.PrefixExpression); ok && isUnaryDependency(pre) {
				_, errs := b.c.resolveEndpoint(sc, pre.Expression, modeDep, nil)
				b.report(errs)
			}
		}
	}
	fields, subs := partitionBody(elementBody(e))
	seen := map[string]ast.Node{}
	for _, f := range fields {
		if ast.IsNil(f.Callee) {
			continue
		}
		if isValidPartialInjection(f.Callee) {
			expr := f.Callee.(*ast.PrefixExpression).Expression
			name, _ := extractVariable(expr)
			if b.c.lookupPartial(sc, name) == nil {
				b.n(diag.BINDING_ERROR, "TablePartial '"+name+"' does not exist in Schema 'public'", expr)
			}
			k := "~" + name
			if first, dup := seen[k]; dup {
				b.n(diag.DUPLICATE_TABLE_PARTIAL_INJECTION_NAME, "Duplicate table partial injection '"+name+"'", first)
				b.n(diag.DUPLICATE_TABLE_PARTIAL_INJECTION_NAME, "Duplicate table partial injection '"+name+"'", f)
			} else {
				seen[k] = f
			}
			continue
		}
		if name, ok := extractVariable(f.Callee); ok {
			if first, dup := seen[name]; dup {
				b.n(diag.DUPLICATE_COLUMN_NAME, "Duplicate column '"+name+"'", first)
				b.n(diag.DUPLICATE_COLUMN_NAME, "Duplicate column '"+name+"'", f)
			} else {
				seen[name] = f
			}
		}
		b.bindColumn(f, t, false)
	}
	for _, s := range subs {
		switch {
		case s.IsKind(kIndexes):
			b.bindIndexes(s, t)
		case s.IsKind(kRecords):
			b.bindRecords(s, t)
		}
	}
}

func (b *binder) bindColumn(f *ast.FunctionApplication, container *sym, partial bool) {
	sc := b.u.scope
	args := append([]ast.Node{f.Callee}, f.Args...)
	if l := trailingList(args); l != nil {
		s, _ := aggregateSettings(l, b.u.path)
		for _, a := range s.get(settingRef) {
			pre, ok := a.Value.(*ast.PrefixExpression)
			if !ok || ast.IsNil(pre.Expression) {
				continue
			}
			_, errs := b.c.resolveEndpoint(sc, pre.Expression, modeInline, container)
			b.report(errs)
		}
		if !partial {
			for _, a := range s.get(settingDep) {
				pre, ok := a.Value.(*ast.PrefixExpression)
				if !ok || ast.IsNil(pre.Expression) {
					continue
				}
				_, errs := b.c.resolveEndpoint(sc, pre.Expression, modeInline, container)
				b.report(errs)
			}
		}
		for _, a := range s.get(settingDefault) {
			b.bindEnumDefault(a.Value)
		}
		args = args[:len(args)-1]
	}
	if len(args) > 1 {
		b.bindEnumType(args[1])
	}
}

// bindEnumType is a no-op for diagnostics: like @dbml/parse, an unresolved
// type is treated as a raw SQL type, never an error.
func (b *binder) bindEnumType(typ ast.Node) {}

func (b *binder) bindEnumDefault(val ast.Node) {
	if ast.IsNil(val) || isQuotedString(val) {
		return
	}
	if isVariableNode(val) {
		return
	}
	t, ok := destructureComplexVariableTuple(val)
	if !ok || len(t.variables) < 2 || len(t.tupleElements) > 0 {
		return
	}
	sc := b.u.scope
	vars := t.variables
	schema := ""
	for _, v := range vars[:len(vars)-2] {
		name, _ := extractVariable(v)
		if !sc.hasSchema(name) {
			b.n(diag.BINDING_ERROR, "Schema '"+name+"' does not exist in Schema 'public'", v)
			return
		}
		schema = name
	}
	enumName, _ := extractVariable(vars[len(vars)-2])
	e := sc.lookup(symEnum, schema, enumName)
	if e == nil {
		msg := "Enum '" + enumName + "' does not exist in Schema 'public'"
		if schema != "" {
			msg = "Enum or schema '" + enumName + "' does not exist"
		}
		b.n(diag.BINDING_ERROR, msg, vars[len(vars)-2])
		return
	}
	field, _ := extractVariable(vars[len(vars)-1])
	if !enumHasField(e.sym, field) {
		b.n(diag.BINDING_ERROR, "Enum field '"+field+"' does not exist in Enum 'public."+e.sym.name+"'", vars[len(vars)-1])
	}
}

func enumFields(e *sym) []string {
	fields, _ := partitionBody(elementBody(e.decl))
	var out []string
	for _, f := range fields {
		if v, ok := extractVariable(f.Callee); ok {
			out = append(out, v)
		}
	}
	return out
}

func enumHasField(e *sym, name string) bool {
	for _, f := range enumFields(e) {
		if f == name {
			return true
		}
	}
	return false
}

func (b *binder) bindPartial(e *ast.ElementDeclaration, p *sym) {
	fields, subs := partitionBody(elementBody(e))
	seen := map[string]ast.Node{}
	for _, f := range fields {
		if ast.IsNil(f.Callee) {
			continue
		}
		if name, ok := extractVariable(f.Callee); ok {
			if first, dup := seen[name]; dup {
				b.n(diag.DUPLICATE_COLUMN_NAME, "Duplicate column '"+name+"'", first)
				b.n(diag.DUPLICATE_COLUMN_NAME, "Duplicate column '"+name+"'", f)
			} else {
				seen[name] = f
			}
		}
		if p != nil {
			b.bindColumn(f, p, true)
		}
	}
	for _, s := range subs {
		if s.IsKind(kIndexes) && p != nil {
			b.bindIndexes(s, p)
		}
	}
}

func (b *binder) bindIndexes(e *ast.ElementDeclaration, owner *sym) {
	if owner == nil {
		return
	}
	fields, _ := partitionBody(elementBody(e))
	ownerKind := "Table"
	if owner.kind == symPartial {
		ownerKind = "TablePartial"
	}
	ownerName := owner.name
	if parts, ok := destructureComplexVariable(owner.decl.Name); ok {
		ownerName = strings.Join(parts, ".")
	}
	for _, f := range fields {
		if ast.IsNil(f.Callee) {
			continue
		}
		args := append([]ast.Node{f.Callee}, f.Args...)
		for _, a := range args {
			for _, bindee := range scanForBinding(a) {
				var targets []ast.Node
				if len(bindee.variables) > 1 {
					continue
				}
				if len(bindee.variables) == 1 {
					targets = bindee.variables
				} else {
					targets = bindee.tupleElements
				}
				for _, t := range targets {
					name, ok := extractVariable(t)
					if !ok {
						continue
					}
					if b.c.containerColumn(owner, b.u.scope, name) == nil {
						b.n(diag.BINDING_ERROR, "No column named '"+name+"' inside "+ownerKind+" '"+ownerName+"'", t)
					}
				}
			}
		}
	}
}

// scanForBinding mirrors scanNonListNodeForBinding.
func scanForBinding(n ast.Node) []complexVariableTuple {
	if ast.IsNil(n) {
		return nil
	}
	if isVariableNode(n) {
		return []complexVariableTuple{{variables: []ast.Node{n}}}
	}
	switch v := n.(type) {
	case *ast.InfixExpression:
		if t, ok := destructureComplexVariableTuple(v); ok {
			return []complexVariableTuple{t}
		}
		return append(scanForBinding(v.Left), scanForBinding(v.Right)...)
	case *ast.PrefixExpression:
		return scanForBinding(v.Expression)
	case *ast.PostfixExpression:
		return scanForBinding(v.Expression)
	case *ast.TupleExpression:
		if t, ok := destructureComplexVariableTuple(v); ok {
			return []complexVariableTuple{t}
		}
		var out []complexVariableTuple
		for _, e := range v.Elements {
			out = append(out, scanForBinding(e)...)
		}
		return out
	case *ast.CallExpression:
		// (a, b)(c, d) in indexes is a call expression: scan both sides.
		out := scanForBinding(v.Callee)
		if v.ArgumentList != nil {
			out = append(out, scanForBinding(v.ArgumentList)...)
		}
		return out
	}
	return nil
}

func refInfix(e *ast.ElementDeclaration) *ast.InfixExpression {
	body := elementBody(e)
	if len(body) == 0 {
		return nil
	}
	f, ok := body[0].(*ast.FunctionApplication)
	if !ok {
		return nil
	}
	in, _ := f.Callee.(*ast.InfixExpression)
	return in
}

func (b *binder) bindRef(e *ast.ElementDeclaration) {
	if !isProgram(parentElem(e)) && !isElem(parentElem(e), kProject) {
		return
	}
	fields, _ := partitionBody(elementBody(e))
	for _, f := range fields {
		if ast.IsNil(f.Callee) {
			continue
		}
		// RefBinder: scan every non-list argument for `[schema.]table.col`
		// or `table.(a, b)` references, even inside malformed expressions.
		for _, a := range append([]ast.Node{f.Callee}, f.Args...) {
			for _, bindee := range scanForBinding(a) {
				vars := bindee.variables
				if len(bindee.tupleElements) == 0 && len(vars) < 2 || len(vars) == 0 {
					continue
				}
				_, errs := b.c.resolveChain(b.u.scope, vars, bindee.tupleElements, modeRef, nil)
				b.report(errs)
			}
		}
	}
}

func (b *binder) bindDep(e *ast.ElementDeclaration) {
	if !isProgram(parentElem(e)) && !isElem(parentElem(e), kProject) {
		return
	}
	fields, _ := partitionBody(elementBody(e))
	for _, f := range fields {
		in, ok := f.Callee.(*ast.InfixExpression)
		if !ok || in == nil {
			continue
		}
		var l, r resolved
		for i, side := range []ast.Node{in.Left, in.Right} {
			if ast.IsNil(side) {
				continue
			}
			res, errs := b.c.resolveEndpoint(b.u.scope, side, modeDep, nil)
			b.report(errs)
			if i == 0 {
				l = res
			} else {
				r = res
			}
		}
		if l.ok && r.ok && len(l.cols) == 1 && len(r.cols) == 1 && l.table.sym == r.table.sym && l.cols[0].name == r.cols[0].name {
			b.n(diag.DEP_SELF_LOOP, "A column cannot depend on itself", in)
		}
	}
}

func (b *binder) bindTableGroup(e *ast.ElementDeclaration) {
	fields, _ := partitionBody(elementBody(e))
	seen := map[*sym]bool{}
	for _, f := range fields {
		if ast.IsNil(f.Callee) {
			continue
		}
		frags := destructureMemberAccess(f.Callee)
		if len(frags) == 0 {
			continue
		}
		for _, fr := range frags {
			if !isVariableNode(fr) {
				frags = nil
			}
		}
		if frags == nil {
			continue
		}
		res, errs := b.c.resolveChain(b.u.scope, frags, nil, modeTable, nil)
		b.report(errs)
		if res.ok && res.table != nil {
			if seen[res.table.sym] {
				name, _ := destructureComplexVariable(f.Callee)
				b.n(diag.DUPLICATE_NAME, "Duplicate TableGroupField '"+strings.Join(name, ".")+"'", f)
			}
			seen[res.table.sym] = true
		}
	}
}

func (b *binder) bindEnum(e *ast.ElementDeclaration) {
	fields, _ := partitionBody(elementBody(e))
	seen := map[string]ast.Node{}
	for _, f := range fields {
		name, ok := extractVariable(f.Callee)
		if !ok {
			continue
		}
		if first, dup := seen[name]; dup {
			b.n(diag.DUPLICATE_COLUMN_NAME, "Duplicate enum field '"+name+"'", first)
			b.n(diag.DUPLICATE_COLUMN_NAME, "Duplicate enum field '"+name+"'", f)
		} else {
			seen[name] = f
		}
	}
}

// recordsTarget returns the table and listed column nodes of a records block.
func (c *compiler) recordsTarget(sc *scope, e *ast.ElementDeclaration, owner *sym) (*visible, []ast.Node, []resolveErr) {
	if owner != nil {
		var cols []ast.Node
		if t, ok := e.Name.(*ast.TupleExpression); ok && t != nil {
			cols = t.Elements
		}
		return sc.visibleOf(owner), cols, nil
	}
	vars, args, ok := destructureCallExpression(e.Name)
	if !ok {
		return nil, nil, nil
	}
	res, errs := c.resolveChain(sc, vars, nil, modeTable, nil)
	if !res.ok {
		if len(errs) == 0 {
			return nil, nil, nil
		}
		// records/bind.ts only binds the table fragment: a missing schema
		// or table is reported once, on the table name.
		last := vars[len(vars)-1]
		name, _ := extractVariable(last)
		msg := "Table '" + name + "' does not exist in Schema 'public'"
		if len(vars) > 1 {
			msg = "Table or schema '" + name + "' does not exist"
		}
		return nil, nil, []resolveErr{{last, msg}}
	}
	return res.table, args, nil
}

func (b *binder) bindRecords(e *ast.ElementDeclaration, owner *sym) {
	b.bindRecordsName(e, owner)
	// Enum references in the data rows (records/bind.ts bindDataRow).
	fields, _ := partitionBody(elementBody(e))
	for _, row := range fields {
		if ast.IsNil(row.Callee) {
			continue
		}
		vals := []ast.Node{row.Callee}
		if ce, ok := row.Callee.(*ast.CommaExpression); ok {
			vals = ce.Elements
		}
		for _, v := range vals {
			for _, bindee := range scanForBinding(v) {
				if len(bindee.variables) >= 2 {
					b.bindEnumValue(bindee.variables)
				}
			}
		}
	}
}

func (b *binder) bindRecordsName(e *ast.ElementDeclaration, owner *sym) {
	t, cols, errs := b.c.recordsTarget(b.u.scope, e, owner)
	b.report(errs)
	if t == nil {
		return
	}
	// Only top-level records check for duplicate column references; each
	// duplicate is reported against the previous occurrence and itself.
	seen := map[string]ast.Node{}
	tname := t.sym.name
	for _, cn := range cols {
		name, ok := extractVariable(cn)
		if !ok {
			continue
		}
		if b.c.lookupColumn(t.sym, t.sym.unit.scope, name) == nil {
			b.n(diag.BINDING_ERROR, "Column '"+name+"' does not exist in Table '"+tname+"'", cn)
			continue
		}
		if owner != nil {
			continue
		}
		if prev, dup := seen[name]; dup {
			b.n(diag.DUPLICATE_COLUMN_REFERENCES_IN_RECORDS, "Column '"+name+"' is referenced more than once in a Records for Table '"+tname+"'", prev)
			b.n(diag.DUPLICATE_COLUMN_REFERENCES_IN_RECORDS, "Column '"+name+"' is referenced more than once in a Records for Table '"+tname+"'", cn)
		}
		seen[name] = cn
	}
}

// bindEnumValue resolves `[schema.]enum.field` in a records row, reporting
// the first fragment that fails (nodeRefereeOfEnumValue).
func (b *binder) bindEnumValue(vars []ast.Node) {
	sc := b.u.scope
	name0, _ := extractVariable(vars[0])
	var schema string
	var enum *visible
	if len(vars) >= 3 {
		if !sc.hasSchema(name0) {
			b.n(diag.BINDING_ERROR, "Schema '"+name0+"' does not exist in Schema 'public'", vars[0])
			return
		}
		schema = name0
	} else {
		enum = sc.lookup(symEnum, "", name0)
		if enum == nil {
			b.n(diag.BINDING_ERROR, "Enum '"+name0+"' does not exist in Schema 'public'", vars[0])
			return
		}
	}
	for _, v := range vars[1:] {
		name, _ := extractVariable(v)
		switch {
		case enum != nil:
			if !enumHasField(enum.sym, name) {
				b.n(diag.BINDING_ERROR, "Enum field '"+name+"' does not exist in Enum 'public."+enum.sym.name+"'", v)
			}
			return
		case schema != "":
			if e := sc.lookup(symEnum, schema, name); e != nil {
				enum, schema = e, ""
			} else if sc.hasSchema(schema + "." + name) {
				schema = schema + "." + name
			} else {
				b.n(diag.BINDING_ERROR, "Enum or schema '"+name+"' does not exist", v)
				return
			}
		}
	}
}

// metadataTarget resolves a Metadata block's target element.
type metaTarget struct {
	sym *sym
	col *col
}

func (c *compiler) resolveMetadataTarget(sc *scope, e *ast.ElementDeclaration) (metaTarget, bool) {
	kind := metadataTargetKind(e)
	parts, ok := destructureComplexVariable(e.Name)
	if !ok || len(parts) == 0 || kind == "" {
		return metaTarget{}, false
	}
	lookupSym := func(k symKind, parts []string) *visible {
		if len(parts) == 1 {
			return sc.lookup(k, "", parts[0])
		}
		return sc.lookup(k, strings.Join(parts[:len(parts)-1], "."), parts[len(parts)-1])
	}
	switch kind {
	case "Table":
		if v := lookupSym(symTable, parts); v != nil {
			return metaTarget{sym: v.sym}, true
		}
	case "TableGroup":
		if v := lookupSym(symGroup, parts); v != nil {
			return metaTarget{sym: v.sym}, true
		}
	case "Note":
		if v := lookupSym(symNote, parts); v != nil {
			return metaTarget{sym: v.sym}, true
		}
	case "Column":
		if len(parts) < 2 {
			return metaTarget{}, false
		}
		if v := lookupSym(symTable, parts[:len(parts)-1]); v != nil {
			if cl := c.lookupColumn(v.sym, v.sym.unit.scope, parts[len(parts)-1]); cl != nil && cl.partial == nil {
				return metaTarget{sym: v.sym, col: cl}, true
			}
		}
	}
	return metaTarget{}, false
}

func (b *binder) bindMetadata(e *ast.ElementDeclaration) {
	if metadataTargetKind(e) == "" || ast.IsNil(e.Name) || !isValidName(e.Name) {
		return
	}
	if _, ok := b.c.resolveMetadataTarget(b.u.scope, e); !ok {
		b.n(diag.BINDING_ERROR, "cannot find metadata target element", e.Name)
	}
}

func (b *binder) bindDiagramView(e *ast.ElementDeclaration) {
	_, subs := partitionBody(elementBody(e))
	sc := b.u.scope
	for _, sub := range subs {
		if _, ok := sub.Body.(*ast.BlockExpression); !ok {
			continue
		}
		fields, _ := partitionBody(elementBody(sub))
		for _, f := range fields {
			if ast.IsNil(f.Callee) || isWildcard(f.Callee) {
				continue
			}
			parts, ok := destructureComplexVariable(f.Callee)
			if !ok || len(parts) == 0 {
				continue
			}
			name := parts[len(parts)-1]
			schema := ""
			if len(parts) > 1 {
				schema = strings.Join(parts[:len(parts)-1], ".")
			}
			switch {
			case sub.IsKind(kDVTables):
				if schema != "" && !sc.hasSchema(schema) {
					b.n(diag.BINDING_ERROR, "Could not find schema '"+schema+"'", f.Callee)
				} else if sc.lookup(symTable, schema, name) == nil {
					b.n(diag.BINDING_ERROR, "Could not find table '"+name+"'", getRightmostVariable(f.Callee))
				}
			case sub.IsKind(kDVNotes):
				if sc.lookup(symNote, schema, name) == nil {
					b.n(diag.BINDING_ERROR, "Could not find note '"+name+"'", f.Callee)
				}
			case sub.IsKind(kDVTableGroups):
				if sc.lookup(symGroup, schema, name) == nil {
					b.n(diag.BINDING_ERROR, "Could not find tablegroup '"+name+"'", f.Callee)
				}
			case sub.IsKind(kDVSchemas):
				if !sc.hasSchema(strings.Join(parts, ".")) {
					b.n(diag.BINDING_ERROR, "Could not find schema '"+name+"'", f.Callee)
				}
			}
		}
	}
}

// checkDuplicateMembers reports elements declared twice in the same schema.
func (b *binder) checkDuplicateMembers(sc *scope) {
	seen := map[string]bool{}
	for _, s := range b.u.decls {
		k := key(s.kind, s.schema, s.name)
		if seen[k] {
			b.n(diag.DUPLICATE_NAME, duplicateMessage(s), s.decl.Name)
			continue
		}
		seen[k] = true
	}
	// aliases are Table members of the public schema
	for _, s := range b.u.decls {
		if s.kind != symTable || s.alias == "" {
			continue
		}
		if s.alias == s.name && s.schema == defaultSchema {
			continue
		}
		k := key(symTable, defaultSchema, s.alias)
		if seen[k] {
			b.n(diag.DUPLICATE_NAME, "Table '"+s.alias+"' already exists in schema 'public'", s.decl.Name)
			continue
		}
		seen[k] = true
	}
	_ = sc
}

func duplicateMessage(s *sym) string {
	switch s.kind {
	case symTable:
		return "Table '" + s.name + "' already exists in schema '" + s.schema + "'"
	case symEnum:
		return "Enum '" + s.name + "' already exists in schema '" + s.schema + "'"
	case symPartial:
		return "TablePartial '" + s.name + "' already exists in schema '" + s.schema + "'"
	case symGroup:
		return "TableGroup '" + s.name + "' already exists in schema '" + s.schema + "'"
	case symView:
		return "DiagramView '" + s.name + "' already exists in schema '" + s.schema + "'"
	}
	return "Duplicate StickyNote '" + s.name + "' in schema '" + s.schema + "'"
}

func colKey(t *sym, name string) string {
	return t.unit.path + "#" + t.qualified() + "#" + name
}

func (b *binder) checkDuplicateRefs() {
	sc := b.u.scope
	seen := map[string]ast.Node{}
	check := func(left, right []string, errNode ast.Node) {
		sort.Strings(left)
		sort.Strings(right)
		lk, rk := strings.Join(left, ","), strings.Join(right, ",")
		if lk == rk {
			b.n(diag.SAME_ENDPOINT, "Two endpoints are the same", errNode)
			return
		}
		id := lk + "-" + rk
		if rk < lk {
			id = rk + "-" + lk
		}
		if existing, ok := seen[id]; ok {
			b.n(diag.CIRCULAR_REF, "References with same endpoints exist", errNode)
			b.n(diag.CIRCULAR_REF, "References with same endpoints exist", existing)
			return
		}
		seen[id] = errNode
	}
	keysOf := func(r resolved) []string {
		out := make([]string, len(r.cols))
		for i, cl := range r.cols {
			out[i] = colKey(colOwner(r.table.sym, cl), cl.name)
		}
		return out
	}
	for _, e := range b.u.prog.Declarations() {
		switch {
		case e.IsKind(kRef):
			for _, f := range func() []*ast.FunctionApplication { fs, _ := partitionBody(elementBody(e)); return fs }() {
				if !isBinaryRelationship(f.Callee) {
					continue
				}
				in := f.Callee.(*ast.InfixExpression)
				l, _ := b.c.resolveEndpoint(sc, in.Left, modeRef, nil)
				r, _ := b.c.resolveEndpoint(sc, in.Right, modeRef, nil)
				if !l.ok || !r.ok {
					continue
				}
				check(keysOf(l), keysOf(r), e)
			}
		case e.IsKind(kTable):
			t := b.u.symOfDecl(e)
			if t == nil {
				continue
			}
			fields, _ := partitionBody(elementBody(e))
			for _, f := range fields {
				name, ok := extractVariable(f.Callee)
				if !ok {
					continue
				}
				l := trailingList(f.Args)
				if l == nil {
					continue
				}
				s, _ := aggregateSettings(l, b.u.path)
				for _, a := range s.get(settingRef) {
					pre, ok := a.Value.(*ast.PrefixExpression)
					if !ok || pre.Op == nil || !isRelationshipOp(pre.Op.Value) || ast.IsNil(pre.Expression) {
						continue
					}
					r, _ := b.c.resolveEndpoint(sc, pre.Expression, modeInline, t)
					if !r.ok {
						continue
					}
					check([]string{colKey(t, name)}, keysOf(r), a)
				}
			}
		case e.IsKind(kTablePartial):
			fields, _ := partitionBody(elementBody(e))
			for _, f := range fields {
				name, ok := extractVariable(f.Callee)
				if !ok {
					continue
				}
				l := trailingList(f.Args)
				if l == nil {
					continue
				}
				s, _ := aggregateSettings(l, b.u.path)
				for _, a := range s.get(settingRef) {
					pre, ok := a.Value.(*ast.PrefixExpression)
					if !ok || pre.Op == nil || !isRelationshipOp(pre.Op.Value) {
						continue
					}
					if target, ok := extractVariable(pre.Expression); ok && target == name {
						b.n(diag.SAME_ENDPOINT, "Two endpoints are the same", a)
					}
				}
			}
		}
	}
}

// colOwner returns the table that owns a resolved column for identity
// purposes (injected columns belong to the table they were injected into).
func colOwner(t *sym, _ *col) *sym { return t }

func (b *binder) checkTableGroups() {
	owner := map[*sym]string{}
	for _, e := range b.u.prog.Declarations() {
		if !e.IsKind(kTableGroup) {
			continue
		}
		groupName := "<unknown>"
		if parts, ok := destructureComplexVariable(e.Name); ok {
			groupName = strings.Join(parts, ".")
		}
		fields, _ := partitionBody(elementBody(e))
		for _, f := range fields {
			frags := destructureMemberAccess(f.Callee)
			if len(frags) == 0 {
				continue
			}
			res, _ := b.c.resolveChain(b.u.scope, frags, nil, modeTable, nil)
			if !res.ok || res.table == nil {
				continue
			}
			if g, ok := owner[res.table.sym]; ok {
				disp := "<unknown>"
				if parts, ok := destructureComplexVariable(f.Callee); ok {
					disp = strings.Join(parts, ".")
				}
				b.n(diag.TABLE_REAPPEAR_IN_TABLEGROUP, "Table \""+disp+"\" already appears in group \""+g+"\"", f)
				continue
			}
			owner[res.table.sym] = groupName
		}
	}
}

var symKindDisplay = map[string]string{"table": "Table", "enum": "Enum", "tablepartial": "TablePartial", "note": "Note", "tablegroup": "TableGroup", "schema": "Schema"}

// bindUses reports unresolvable `use`/`reuse` declarations: a missing or
// non-relative file on the path string, and per specifier a missing file or
// element on the specifier name (mirrors use/index.ts).
func (b *binder) bindUses() {
	for _, imp := range b.u.imports {
		ud := imp.decl
		if ud.ImportPath != nil && imp.target == nil {
			b.diags = append(b.diags, tokDiag(diag.NONEXISTENT_MODULE, "Failed to resolve the non-existent file '"+ud.ImportPath.Value+"'", ud.ImportPath, b.u.path))
		}
		for _, spec := range imp.specs {
			if spec.ImportKind == nil || ast.IsNil(spec.Name) {
				continue
			}
			kind := strings.ToLower(spec.ImportKind.Value)
			display, ok := symKindDisplay[kind]
			if !ok {
				continue
			}
			parts, ok := destructureComplexVariable(spec.Name)
			if !ok || len(parts) == 0 {
				continue
			}
			full := strings.Join(parts, ".")
			switch {
			case imp.path == "":
				b.n(diag.BINDING_ERROR, "Import path must be relative, got '"+ud.ImportPath.Value+"'", spec.Name)
			case imp.target == nil:
				b.n(diag.NONEXISTENT_MODULE, display+" '"+full+"' does not exist in file "+imp.path+". Does the file exist?", spec.Name)
			case !b.c.importExists(imp.target, kind, parts):
				b.n(diag.BINDING_ERROR, display+" '"+parts[len(parts)-1]+"' does not exist in file "+imp.path, spec.Name)
			}
		}
	}
	for _, un := range b.u.unresolved {
		kind := strings.ToLower(un.spec.ImportKind.Value)
		b.n(diag.BINDING_ERROR, "Failed to resolve the import of "+symKindDisplay[kind]+" '"+un.display+"'", un.spec)
	}
}

// importExists reports whether target exposes an element of kind named parts.
func (c *compiler) importExists(target *unit, kind string, parts []string) bool {
	avail := c.exported(target, map[*unit]bool{})
	full := strings.Join(parts, ".")
	if kind == "schema" {
		for _, v := range avail {
			if v.schema == full || strings.HasPrefix(v.schema, full+".") {
				return true
			}
		}
		return false
	}
	k := importSymKinds[kind]
	name := parts[len(parts)-1]
	schema := defaultSchema
	if len(parts) > 1 {
		schema = strings.Join(parts[:len(parts)-1], ".")
	}
	for _, v := range avail {
		if v.sym.kind == k && v.name == name && v.schema == schema {
			return true
		}
	}
	return false
}
