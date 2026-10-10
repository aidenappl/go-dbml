// Portions of this file are derived from holistics/dbml packages/dbml-parse/src/core/local_modules/*/validate.ts
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

package model

import (
	"strings"

	"github.com/aidenappl/go-dbml/ast"
	"github.com/aidenappl/go-dbml/diag"
)

// validator performs the per-element syntactic checks of @dbml/parse's
// local modules (context, name, alias, settings and body shape).
type validator struct {
	c     *compiler
	u     *unit
	diags []diag.Diagnostic
}

// n records a diagnostic over node (ignored when node is nil).
func (v *validator) n(code diag.Code, msg string, node ast.Node) {
	if ast.IsNil(node) {
		return
	}
	v.diags = append(v.diags, nodeDiag(code, msg, node, v.u.path))
}

// add records several diagnostics.
func (v *validator) add(d []diag.Diagnostic) { v.diags = append(v.diags, d...) }

// validateProgram validates every top-level declaration of the file.
func (v *validator) validateProgram() {
	for _, el := range v.u.prog.Body {
		v.c.poll(false)
		v.validateNode(el)
	}
}

// validateNode dispatches on the element keyword, like the chain of local
// modules in @dbml/parse; unknown keywords are "custom" elements.
func (v *validator) validateNode(n ast.Node) {
	switch e := n.(type) {
	case *ast.UseDeclaration:
		v.validateUse(e)
	case *ast.ElementDeclaration:
		if e.Type == nil {
			return
		}
		switch {
		case e.IsKind(kTable):
			v.validateTable(e)
		case e.IsKind(kEnum):
			v.validateEnum(e)
		case e.IsKind(kRecords):
			v.validateRecords(e)
		case e.IsKind(kIndexes):
			v.validateIndexes(e)
		case e.IsKind(kChecks):
			v.validateChecks(e)
		case e.IsKind(kRef):
			v.validateRef(e)
		case e.IsKind(kDep):
			v.validateDep(e)
		case e.IsKind(kProject):
			v.validateProject(e)
		case e.IsKind(kTableGroup):
			v.validateTableGroup(e)
		case e.IsKind(kTablePartial):
			v.validateTablePartial(e)
		case e.IsKind(kNote):
			v.validateNote(e)
		case e.IsKind(kMetadata):
			v.validateMetadata(e)
		case e.IsKind(kDiagramView):
			v.validateDiagramView(e)
		default:
			v.validateCustom(e)
		}
	}
}

// validateSubs validates nested element declarations.
func (v *validator) validateSubs(subs []*ast.ElementDeclaration) {
	for _, s := range subs {
		if s.Type != nil {
			v.validateNode(s)
		}
	}
}

// noteRedefined reports every Note when an element has more than one.
func (v *validator) noteRedefined(subs []*ast.ElementDeclaration) {
	var notes []*ast.ElementDeclaration
	for _, s := range subs {
		if s.IsKind(kNote) {
			notes = append(notes, s)
		}
	}
	if len(notes) > 1 {
		for _, n := range notes {
			v.n(diag.NOTE_REDEFINED, "Duplicate notes are defined", n)
		}
	}
}

// topLevelOnly reports e when it is nested inside another element.
func (v *validator) topLevelOnly(e *ast.ElementDeclaration, code diag.Code, msg string) {
	if _, ok := parentElem(e).(*ast.ElementDeclaration); ok {
		v.n(code, msg, e)
	}
}

// noAlias reports an alias on elements that cannot have one.
func (v *validator) noAlias(e *ast.ElementDeclaration, msg string) {
	if !ast.IsNil(e.Alias) {
		v.n(diag.UNEXPECTED_ALIAS, msg, e.Alias)
	}
}

// noSettings reports a setting list on elements that cannot have one.
func (v *validator) noSettings(e *ast.ElementDeclaration, msg string) {
	if e.AttributeList != nil {
		v.n(diag.UNEXPECTED_SETTINGS, msg, e.AttributeList)
	}
}

// noName reports a name on elements that cannot have one.
func (v *validator) noName(e *ast.ElementDeclaration, msg string) {
	if !ast.IsNil(e.Name) {
		v.n(diag.UNEXPECTED_NAME, msg, e.Name)
	}
}

// blockBody returns (fields, subs, ok) for a block body, reporting a simple
// body with msg when msg != "".
func (v *validator) blockBody(e *ast.ElementDeclaration, simpleMsg string) ([]*ast.FunctionApplication, []*ast.ElementDeclaration, bool) {
	if ast.IsNil(e.Body) {
		return nil, nil, false
	}
	if fa, ok := e.Body.(*ast.FunctionApplication); ok {
		if simpleMsg != "" {
			v.n(diag.UNEXPECTED_SIMPLE_BODY, simpleMsg, fa)
		}
		return nil, nil, false
	}
	f, s := partitionBody(elementBody(e))
	return f, s, true
}

// nameKindCheck validates an element name: required-ness, wildcard and a
// shape predicate (simple name or schema-qualified name).
func (v *validator) nameKindCheck(e *ast.ElementDeclaration, required bool, what, wildcardMsg string, valid func(ast.Node) bool, invalidMsg string) {
	if ast.IsNil(e.Name) {
		if required {
			v.n(diag.NAME_NOT_FOUND, what+" must have a name", e)
		}
		return
	}
	if isWildcard(e.Name) && wildcardMsg != "" {
		v.n(diag.INVALID_NAME, wildcardMsg, e.Name)
		return
	}
	if !valid(e.Name) {
		v.n(diag.INVALID_NAME, invalidMsg, e.Name)
	}
}

// ---------------------------------------------------------------- table

// validateTable checks a Table: context, name, alias, settings, columns,
// partial injections and sub-elements.
func (v *validator) validateTable(e *ast.ElementDeclaration) {
	v.topLevelOnly(e, diag.INVALID_TABLE_CONTEXT, "Table must appear top-level")
	if ast.IsNil(e.Name) {
		v.n(diag.NAME_NOT_FOUND, "A Table must have a name", e)
	} else if _, ok := e.Name.(*ast.Array); ok {
		v.n(diag.INVALID_NAME, "Invalid array as Table name, maybe you forget to add a space between the name and the setting list?", e.Name)
	} else if isWildcard(e.Name) {
		v.n(diag.INVALID_NAME, "Wildcard (*) is not allowed as a Table name", e.Name)
	} else if !isValidName(e.Name) {
		v.n(diag.INVALID_NAME, "A Table name must be of the form <table> or <schema>.<table>", e.Name)
	}
	if !ast.IsNil(e.Alias) && !isSimpleName(e.Alias) {
		v.n(diag.INVALID_ALIAS, "Table aliases can only contain alphanumeric and underscore unless surrounded by double quotes", e.Alias)
	}
	_, errs := tableSettings(e.AttributeList, v.u.path)
	v.add(errs)
	fields, subs, ok := v.blockBody(e, "A Table's body must be a block")
	if !ok {
		return
	}
	for _, f := range fields {
		if ast.IsNil(f.Callee) {
			continue
		}
		if pre, ok := f.Callee.(*ast.PrefixExpression); ok && pre.Op != nil && pre.Op.Value == "~" {
			if !isValidPartialInjection(f.Callee) {
				v.n(diag.INVALID_TABLE_PARTIAL_INJECTION, "A partial injection should be of the form ~<table-partial>", f.Callee)
			}
			for _, a := range f.Args {
				v.n(diag.INVALID_TABLE_PARTIAL_INJECTION, "A partial injection does not have any trailing attributes", a)
			}
			continue
		}
		v.validateColumn(f, false)
	}
	v.validateSubs(subs)
	v.noteRedefined(subs)
}

// validateColumn checks `name type [settings]`, including legacy inline
// `pk`/`unique` settings.
func (v *validator) validateColumn(f *ast.FunctionApplication, partial bool) {
	if len(f.Args) == 0 {
		v.n(diag.INVALID_COLUMN, "A column must have a type", f.Callee)
	}
	if !isVariableNode(f.Callee) {
		v.n(diag.INVALID_COLUMN_NAME, "A column name must be an identifier or a quoted identifier", f.Callee)
	}
	if len(f.Args) > 0 && !isValidColumnType(f.Args[0]) {
		v.n(diag.INVALID_COLUMN_TYPE, "Invalid column type", f.Args[0])
	}
	var rest []ast.Node
	if len(f.Args) > 1 {
		rest = f.Args[1:]
	}
	_, errs := columnSettings(rest, v.u.path, partial)
	v.add(errs)
}

// ---------------------------------------------------------------- partial

// validateTablePartial checks a TablePartial (no dep settings on columns).
func (v *validator) validateTablePartial(e *ast.ElementDeclaration) {
	v.topLevelOnly(e, diag.INVALID_TABLE_PARTIAL_CONTEXT, "TablePartial must appear top-level")
	v.nameKindCheck(e, true, "A TablePartial", "Wildcard (*) is not allowed as a TablePartial name", isSimpleName, "A TablePartial name must be an identifier or a quoted identifer")
	v.noAlias(e, "A TablePartial shouldn't have an alias")
	s, errs := aggregateSettings(e.AttributeList, v.u.path)
	v.add(errs)
	for _, name := range s.keys {
		attrs := s.get(name)
		switch name {
		case settingHeaderColor:
			if len(attrs) > 1 {
				for _, a := range attrs {
					v.n(diag.DUPLICATE_TABLE_PARTIAL_SETTING, "'"+name+"' can only appear once", a)
				}
			}
			for _, a := range attrs {
				if !isValidHexColor(a.Value) {
					v.n(diag.INVALID_TABLE_PARTIAL_SETTING_VALUE, "'"+name+"' must be a color literal", attrErrNode(a))
				}
			}
		case settingNote:
			if len(attrs) > 1 {
				for _, a := range attrs {
					v.n(diag.DUPLICATE_TABLE_PARTIAL_SETTING, "'"+name+"' can only appear once", a)
				}
			}
			for _, a := range attrs {
				if !isQuotedString(a.Value) {
					v.n(diag.INVALID_TABLE_PARTIAL_SETTING_VALUE, "'"+name+"' must be a string literal", attrErrNode(a))
				}
			}
		default:
			for _, a := range attrs {
				v.n(diag.UNKNOWN_TABLE_PARTIAL_SETTING, "Unknown '"+name+"' setting", a)
			}
		}
	}
	fields, subs, ok := v.blockBody(e, "A TablePartial's body must be a block")
	if !ok {
		return
	}
	for _, f := range fields {
		if ast.IsNil(f.Callee) {
			continue
		}
		v.validateColumn(f, true)
	}
	v.validateSubs(subs)
	v.noteRedefined(subs)
}

// ---------------------------------------------------------------- enum

// validateEnum checks an Enum and its members.
func (v *validator) validateEnum(e *ast.ElementDeclaration) {
	v.topLevelOnly(e, diag.INVALID_PROJECT_CONTEXT, "An Enum can only appear top-level")
	v.nameKindCheck(e, true, "An Enum", "Wildcard (*) is not allowed as an Enum name", isValidName, "An Enum name must be of the form <enum> or <schema>.<enum>")
	v.noAlias(e, "An Enum shouldn't have an alias")
	v.noSettings(e, "An Enum shouldn't have a setting list")
	if ast.IsNil(e.Body) {
		return
	}
	fields, subs := partitionBody(elementBody(e))
	if len(fields) == 0 {
		v.n(diag.EMPTY_ENUM, "An Enum must have at least one element", e)
	}
	for _, f := range fields {
		if !ast.IsNil(f.Callee) && !isVariableNode(f.Callee) {
			v.n(diag.INVALID_ENUM_ELEMENT_NAME, "An enum field must be an identifier or a quoted identifier", f.Callee)
		}
		args := append([]ast.Node{}, f.Args...)
		if l := trailingList(args); l != nil {
			v.enumFieldSettings(l)
			args = args[:len(args)-1]
		} else if len(args) > 0 {
			if l, ok := args[0].(*ast.ListExpression); ok {
				v.enumFieldSettings(l)
				args = args[1:]
			}
		}
		for _, a := range args {
			v.n(diag.INVALID_ENUM_ELEMENT, "An Enum must have only a field and optionally a setting list", a)
		}
	}
	v.validateSubs(subs)
}

// enumFieldSettings checks an enum member setting list (only note).
func (v *validator) enumFieldSettings(l *ast.ListExpression) {
	s, errs := aggregateSettings(l, v.u.path)
	v.add(errs)
	for _, name := range s.keys {
		attrs := s.get(name)
		if name == settingNote {
			if len(attrs) > 1 {
				for _, a := range attrs {
					v.n(diag.DUPLICATE_ENUM_ELEMENT_SETTING, "'note' can only appear once", a)
				}
			}
			for _, a := range attrs {
				if !isQuotedString(a.Value) {
					v.n(diag.INVALID_ENUM_ELEMENT_SETTING, "'note' must be a string", a)
				}
			}
			continue
		}
		for _, a := range attrs {
			v.n(diag.UNKNOWN_ENUM_ELEMENT_SETTING, "Unknown enum field setting '"+name+"'", a)
		}
	}
}

// ---------------------------------------------------------------- ref

// validateRef checks a short or long Ref: exactly one binary relationship
// between `table.column` (or `table.(a, b)`) endpoints of equal arity.
func (v *validator) validateRef(e *ast.ElementDeclaration) {
	if !isProgram(parentElem(e)) {
		v.n(diag.INVALID_REF_CONTEXT, "A Ref must appear top-level", e)
	}
	if !ast.IsNil(e.Name) {
		if isWildcard(e.Name) {
			v.n(diag.INVALID_NAME, "Wildcard (*) is not allowed as a Ref name", e.Name)
		} else if !isSimpleName(e.Name) {
			v.n(diag.INVALID_NAME, "A Ref's name is optional or must be an identifier or a quoted identifer", e.Name)
		}
	}
	v.noAlias(e, "A Ref shouldn't have an alias")
	v.noSettings(e, "A Ref shouldn't have a setting list")
	if ast.IsNil(e.Body) {
		return
	}
	fields, subs := partitionBody(elementBody(e))
	if len(fields) == 0 {
		v.n(diag.EMPTY_REF, "A Ref must have at least one field", e)
	}
	for i, f := range fields {
		if i > 0 {
			v.n(diag.REF_REDEFINED, "A Ref can only contain one binary relationship", f)
		}
	}
	for _, f := range fields {
		if !ast.IsNil(f.Callee) && !isBinaryRelationship(f.Callee) {
			v.n(diag.INVALID_REF_FIELD, "A Ref field must be a binary relationship", f.Callee)
		}
		if !ast.IsNil(f.Callee) && isBinaryRelationship(f.Callee) {
			in := f.Callee.(*ast.InfixExpression)
			l, _ := destructureComplexVariableTuple(in.Left)
			r, _ := destructureComplexVariableTuple(in.Right)
			if len(l.variables)+min(len(l.tupleElements), 1) < 2 {
				v.n(diag.INVALID_REF_FIELD, "Invalid column reference", in.Left)
			}
			if len(r.variables)+min(len(r.tupleElements), 1) < 2 {
				v.n(diag.INVALID_REF_FIELD, "Invalid column reference", in.Right)
			}
		}
		if !ast.IsNil(f.Callee) && !isEqualTupleOperands(f.Callee) {
			v.n(diag.UNEQUAL_FIELDS_BINARY_REF, "Unequal fields in ref endpoints", f.Callee)
		}
		args := append([]ast.Node{}, f.Args...)
		if l := trailingList(args); l != nil {
			_, errs := refFieldSettings(l, v.u.path)
			v.add(errs)
			args = args[:len(args)-1]
		} else if len(args) > 0 {
			if l, ok := args[0].(*ast.ListExpression); ok {
				_, errs := refFieldSettings(l, v.u.path)
				v.add(errs)
				args = args[1:]
			}
		}
		for _, a := range args {
			v.n(diag.INVALID_REF_FIELD, "A Ref field should only have a single binary relationship", a)
		}
	}
	v.validateSubs(subs)
}

// ---------------------------------------------------------------- dep

// validateDep checks a Dep: `->`/`<-` edges and block settings.
func (v *validator) validateDep(e *ast.ElementDeclaration) {
	if !isProgram(parentElem(e)) {
		v.n(diag.INVALID_DEP_CONTEXT, "A Dep must appear top-level", e)
	}
	if !ast.IsNil(e.Name) {
		if isWildcard(e.Name) {
			v.n(diag.INVALID_NAME, "Wildcard (*) is not allowed as a Dep name", e.Name)
		} else if !isSimpleName(e.Name) {
			v.n(diag.INVALID_NAME, "A Dep's name is optional or must be an identifier or a quoted identifier", e.Name)
		}
	}
	v.noAlias(e, "A Dep shouldn't have an alias")
	if e.AttributeList != nil {
		_, errs := depSettings(e.AttributeList, v.u.path)
		v.add(errs)
	}
	if fa, ok := e.Body.(*ast.FunctionApplication); ok {
		for _, a := range fa.Args {
			if l, ok := a.(*ast.ListExpression); ok {
				_, errs := depSettings(l, v.u.path)
				v.add(errs)
				break
			}
		}
	}
	if ast.IsNil(e.Body) {
		return
	}
	fields, subs := partitionBody(elementBody(e))
	for _, f := range fields {
		if ast.IsNil(f.Callee) {
			v.n(diag.INVALID_DEP_FIELD, "A Dep field must be a binary relationship", f)
			continue
		}
		in, ok := f.Callee.(*ast.InfixExpression)
		if !ok {
			v.n(diag.INVALID_DEP_FIELD, "A Dep field must be a binary relationship", f.Callee)
			continue
		}
		if in.Op == nil || !isDependencyOp(in.Op.Value) {
			v.n(diag.INVALID_DEP_FIELD, "Dep edges must use the '->' or '<-' operator", f.Callee)
		}
		l, _ := destructureComplexVariableTuple(in.Left)
		r, _ := destructureComplexVariableTuple(in.Right)
		if len(l.variables) == 0 && len(l.tupleElements) == 0 {
			v.n(diag.INVALID_DEP_FIELD, "Invalid Dep endpoint", orNode(in.Left, f.Callee))
		}
		if len(r.variables) == 0 && len(r.tupleElements) == 0 {
			v.n(diag.INVALID_DEP_FIELD, "Invalid Dep endpoint", orNode(in.Right, f.Callee))
		}
		args := append([]ast.Node{}, f.Args...)
		if trailingList(args) != nil {
			args = args[:len(args)-1]
		} else if len(args) > 0 {
			if _, ok := args[0].(*ast.ListExpression); ok {
				args = args[1:]
			}
		}
		for _, a := range args {
			v.n(diag.INVALID_DEP_FIELD, "A Dep field should only have a single binary relationship", a)
		}
	}
	for _, sub := range subs {
		if sub.Type == nil {
			continue
		}
		k := strings.ToLower(sub.Type.Value)
		fa, simple := sub.Body.(*ast.FunctionApplication)
		switch k {
		case settingColor:
			if !ast.IsNil(sub.Name) || !simple || ast.IsNil(fa.Callee) || !isValidHexColor(fa.Callee) {
				v.n(diag.INVALID_SETTINGS, "Invalid color value. Expected a hex color (e.g. #fff or #aabbcc)", sub)
			}
		case settingNote:
			var val ast.Node
			if simple {
				val = fa.Callee
			} else if b, ok := sub.Body.(*ast.BlockExpression); ok && len(b.Body) > 0 {
				if f0, ok := b.Body[0].(*ast.FunctionApplication); ok {
					val = f0.Callee
				}
			}
			if ast.IsNil(val) || !isQuotedString(val) {
				v.n(diag.INVALID_SETTINGS, "Invalid note value. Expected a quoted string", sub)
			}
		default:
			if _, isBlock := sub.Body.(*ast.BlockExpression); !ast.IsNil(sub.Name) || isBlock {
				v.n(diag.INVALID_SETTINGS, "'"+k+"' must be a simple field (e.g. "+k+": value). Block syntax is not allowed", sub)
				continue
			}
			if !simple || ast.IsNil(fa.Callee) {
				continue
			}
			if !isQuotedString(fa.Callee) && !isValidHexColor(fa.Callee) && !isSignedNumber(fa.Callee) && !isSimpleName(fa.Callee) {
				v.n(diag.INVALID_SETTINGS, "Invalid value for '"+k+"'. Expected a string, number, color, or identifier", sub)
			}
		}
	}
}

// orNode returns a unless it is nil, else b.
func orNode(a, b ast.Node) ast.Node {
	if !ast.IsNil(a) {
		return a
	}
	return b
}

// ---------------------------------------------------------------- project

// validateProject checks a Project block (only sub-elements allowed).
func (v *validator) validateProject(e *ast.ElementDeclaration) {
	v.topLevelOnly(e, diag.INVALID_PROJECT_CONTEXT, "A Project can only appear top-level")
	if !ast.IsNil(e.Name) && !isSimpleName(e.Name) {
		v.n(diag.INVALID_NAME, "A Project's name is optional or must be an identifier or a quoted identifer", e.Name)
	}
	v.noAlias(e, "A Project shouldn't have an alias")
	v.noSettings(e, "A Project shouldn't have a setting list")
	fields, subs, ok := v.blockBody(e, "A Project's body must be a block")
	if !ok {
		return
	}
	for _, f := range fields {
		v.n(diag.INVALID_PROJECT_FIELD, "A Project can not have inline fields", f)
	}
	v.validateSubs(subs)
}

// ---------------------------------------------------------------- tablegroup

// validateTableGroup checks a TableGroup: settings and table names.
func (v *validator) validateTableGroup(e *ast.ElementDeclaration) {
	v.topLevelOnly(e, diag.INVALID_TABLEGROUP_CONTEXT, "TableGroup must appear top-level")
	v.nameKindCheck(e, true, "A TableGroup", "Wildcard (*) is not allowed as a TableGroup name", isSimpleName, "A TableGroup name must be a single identifier")
	v.noAlias(e, "A TableGroup shouldn't have an alias")
	s, errs := aggregateSettings(e.AttributeList, v.u.path)
	v.add(errs)
	for _, name := range s.keys {
		attrs := s.get(name)
		switch name {
		case settingColor, settingNote:
			if len(attrs) > 1 {
				for _, a := range attrs {
					v.n(diag.DUPLICATE_TABLE_SETTING, "'"+name+"' can only appear once", a)
				}
			}
			for _, a := range attrs {
				if name == settingColor && !isValidHexColor(a.Value) {
					v.n(diag.INVALID_TABLE_SETTING_VALUE, "'color' must be a color literal", attrErrNode(a))
				}
				if name == settingNote && !isQuotedString(a.Value) {
					v.n(diag.INVALID_TABLE_SETTING_VALUE, "'note' must be a string literal", attrErrNode(a))
				}
			}
		default:
			sink := &errSink{file: v.u.path}
			sink.customMetadata(name, attrs, diag.DUPLICATE_TABLE_SETTING, diag.INVALID_TABLE_SETTING_VALUE)
			v.add(sink.diags)
		}
	}
	fields, subs, ok := v.blockBody(e, "A TableGroup's body must be a block")
	if !ok {
		return
	}
	for _, f := range fields {
		if !ast.IsNil(f.Callee) {
			if _, ok := destructureComplexVariable(f.Callee); !ok {
				v.n(diag.INVALID_TABLEGROUP_FIELD, "A TableGroup field must be of the form <table> or <schema>.<table>", f.Callee)
			}
		}
		for _, a := range f.Args {
			v.n(diag.INVALID_TABLEGROUP_FIELD, "A TableGroup field should only have a single Table name", a)
		}
	}
	v.validateSubs(subs)
	v.noteRedefined(subs)
}

// ---------------------------------------------------------------- note

// validateNote checks element notes and sticky (top-level) notes.
func (v *validator) validateNote(e *ast.ElementDeclaration) {
	parent := parentElem(e)
	sticky := isProgram(parent)
	if !sticky && !isElem(parent, kTable, kTableGroup, kTablePartial, kProject, kDep) {
		v.n(diag.INVALID_NOTE_CONTEXT, "A Note can only appear inside a Table, a TableGroup, a TablePartial or a Project. Sticky note can only appear at the global scope.", e)
	}
	if !sticky {
		v.noName(e, "A Note shouldn't have a name")
	} else if ast.IsNil(e.Name) {
		v.n(diag.INVALID_NAME, "Sticky note must have a name", e)
	} else if isWildcard(e.Name) {
		v.n(diag.INVALID_NAME, "Wildcard (*) is not allowed as a Note name", e.Name)
	} else if _, ok := destructureComplexVariable(e.Name); !ok {
		v.n(diag.INVALID_NAME, "Invalid name for sticky note ", e)
	}
	v.noAlias(e, "A Note shouldn't have an alias")
	if e.AttributeList != nil {
		if !sticky {
			v.n(diag.UNEXPECTED_SETTINGS, "A Note shouldn't have a setting list", e.AttributeList)
		} else {
			s, errs := aggregateSettings(e.AttributeList, v.u.path)
			v.add(errs)
			for _, name := range s.keys {
				attrs := s.get(name)
				if name == settingColor {
					if len(attrs) > 1 {
						for _, a := range attrs {
							v.n(diag.DUPLICATE_NOTE_SETTING, "'color' can only appear once", a)
						}
					}
					for _, a := range attrs {
						if !isValidColorOrNone(a.Value) {
							v.n(diag.INVALID_NOTE_SETTING_VALUE, "'color' must be a color literal or 'none'", attrErrNode(a))
						}
					}
					continue
				}
				sink := &errSink{file: v.u.path}
				sink.customMetadata(name, attrs, diag.DUPLICATE_NOTE_SETTING, diag.INVALID_NOTE_SETTING_VALUE)
				v.add(sink.diags)
			}
		}
	}
	if ast.IsNil(e.Body) {
		return
	}
	fields, subs := partitionBody(elementBody(e))
	if len(fields) == 0 {
		v.n(diag.EMPTY_NOTE, "A Note must have a content", e)
	} else {
		for _, f := range fields[1:] {
			v.n(diag.NOTE_CONTENT_REDEFINED, "A Note can only contain one string", f)
		}
		if !isQuotedString(fields[0].Callee) {
			v.n(diag.INVALID_NOTE, "A Note content must be a quoted string", fields[0])
		}
		for _, a := range fields[0].Args {
			v.n(diag.INVALID_NOTE, "A Note can only contain one quoted string", a)
		}
	}
	v.validateSubs(subs)
}

// ---------------------------------------------------------------- indexes / checks

// inTableOrPartial reports e unless it is nested in a Table or TablePartial.
func (v *validator) inTableOrPartial(e *ast.ElementDeclaration, code diag.Code, msg string) {
	p := parentElem(e)
	if isProgram(p) || !isElem(p, kTable, kTablePartial) {
		v.n(code, msg, e)
	}
}

// validateIndexes checks an indexes block: each field is a column, a
// backtick expression or a tuple of those, plus an optional setting list.
func (v *validator) validateIndexes(e *ast.ElementDeclaration) {
	v.inTableOrPartial(e, diag.INVALID_INDEXES_CONTEXT, "An Indexes can only appear inside a Table or a TablePartial")
	v.noName(e, "An Indexes shouldn't have a name")
	v.noAlias(e, "An Indexes shouldn't have an alias")
	v.noSettings(e, "An Indexes shouldn't have a setting list")
	fields, subs, ok := v.blockBody(e, "An Indexes must have a complex body")
	if !ok {
		return
	}
	for _, f := range fields {
		if ast.IsNil(f.Callee) {
			continue
		}
		args := append([]ast.Node{f.Callee}, f.Args...)
		if l := trailingList(args); l != nil {
			_, errs := indexFieldSettings(l, v.u.path)
			v.add(errs)
			args = args[:len(args)-1]
		}
		for _, sub := range args {
			for {
				call, ok := sub.(*ast.CallExpression)
				if !ok {
					break
				}
				if call.ArgumentList != nil {
					if _, ok := destructureIndexNode(call.ArgumentList); !ok {
						v.n(diag.INVALID_INDEXES_FIELD, "An index field must be an identifier, a quoted identifier, a functional expression or a tuple of such", call.ArgumentList)
					}
				}
				sub = call.Callee
			}
			if ast.IsNil(sub) {
				continue
			}
			if _, ok := destructureIndexNode(sub); !ok {
				v.n(diag.INVALID_INDEXES_FIELD, "An index field must be an identifier, a quoted identifier, a functional expression or a tuple of such", sub)
			}
		}
	}
	v.validateSubs(subs)
}

// validateChecks checks a checks block: backtick expressions with an
// optional [name: ...].
func (v *validator) validateChecks(e *ast.ElementDeclaration) {
	v.inTableOrPartial(e, diag.INVALID_CHECKS_CONTEXT, "A Checks can only appear inside a Table or a TablePartial")
	v.noName(e, "A Checks shouldn't have a name")
	v.noAlias(e, "A Checks shouldn't have an alias")
	v.noSettings(e, "A Checks shouldn't have a setting list")
	fields, subs, ok := v.blockBody(e, "A Checks must have a complex body")
	if !ok {
		return
	}
	for _, f := range fields {
		if ast.IsNil(f.Callee) {
			continue
		}
		args := append([]ast.Node{f.Callee}, f.Args...)
		if l := trailingList(args); l != nil {
			s, errs := aggregateSettings(l, v.u.path)
			v.add(errs)
			for _, name := range s.keys {
				attrs := s.get(name)
				if name == settingName {
					if len(attrs) > 1 {
						for _, a := range attrs {
							v.n(diag.DUPLICATE_CHECK_SETTING, "'name' can only appear once", a)
						}
					}
					for _, a := range attrs {
						if !isQuotedString(a.Value) {
							v.n(diag.INVALID_CHECK_SETTING_VALUE, "'name' must be a string", a)
						}
					}
					continue
				}
				for _, a := range attrs {
					v.n(diag.UNKNOWN_CHECK_SETTING, "Unknown check setting '"+name+"'", a)
				}
			}
			args = args[:len(args)-1]
		}
		if _, isFn := args[0].(*ast.FunctionExpression); len(args) > 1 || !isFn {
			v.n(diag.INVALID_CHECKS_FIELD, "A check field must be a function expression", f)
		}
	}
	v.validateSubs(subs)
}

// ---------------------------------------------------------------- records

// isValidRecordValue accepts the literal forms allowed in a records row.
func isValidRecordValue(n ast.Node) bool {
	switch n.(type) {
	case *ast.Empty, *ast.FunctionExpression:
		return true
	}
	if isSignedNumber(n) || isQuotedString(n) {
		return true
	}
	if isVariableNode(n) {
		switch strings.ToLower(variableOf(n).Value) {
		case "null", "true", "false":
			return true
		}
		return false
	}
	if isAccess(n) {
		f, ok := destructureComplexVariable(n)
		return ok && len(f) > 0
	}
	return false
}

// validateRecords checks a records block (top-level needs table(cols),
// nested may only list columns) and its rows.
func (v *validator) validateRecords(e *ast.ElementDeclaration) {
	parent := parentElem(e)
	top := isProgram(parent)
	if !top && !isElem(parent, kTable) {
		v.n(diag.INVALID_RECORDS_CONTEXT, "Records can only appear at top-level or inside a Table", e)
	}
	if isWildcard(e.Name) {
		v.n(diag.INVALID_NAME, "Wildcard (*) is not allowed as a Records name", e.Name)
	} else if top {
		call, ok := e.Name.(*ast.CallExpression)
		if !ok {
			var target ast.Node = e
			if !ast.IsNil(e.Name) {
				target = e.Name
			}
			if ast.IsNil(e.Name) && e.Type != nil {
				v.diags = append(v.diags, tokDiag(diag.INVALID_RECORDS_NAME, "Records at top-level must have a name in the form of table(col1, col2, ...) or schema.table(col1, col2, ...)", e.Type, v.u.path))
			} else {
				v.n(diag.INVALID_RECORDS_NAME, "Records at top-level must have a name in the form of table(col1, col2, ...) or schema.table(col1, col2, ...)", target)
			}
		} else {
			if ast.IsNil(call.Callee) || !isValidName(call.Callee) {
				v.n(diag.INVALID_RECORDS_NAME, "Records table reference must be a valid table name", orNode(call.Callee, call))
			}
			if call.ArgumentList == nil || !isTupleOfVariables(call.ArgumentList) {
				var target ast.Node = call
				if call.ArgumentList != nil {
					target = call.ArgumentList
				}
				v.n(diag.INVALID_RECORDS_NAME, "Records column list must be simple column names", target)
			}
		}
	} else if !ast.IsNil(e.Name) && !isTupleOfVariables(e.Name) {
		v.n(diag.INVALID_RECORDS_NAME, "Records inside a Table can only have a column list like (col1, col2, ...)", e.Name)
	}
	v.noAlias(e, "Records cannot have an alias")
	s, errs := aggregateSettings(e.AttributeList, v.u.path)
	v.add(errs)
	for _, name := range s.keys {
		attrs := s.get(name)
		if name == settingExample {
			if len(attrs) > 1 {
				for _, a := range attrs {
					v.n(diag.DUPLICATE_RECORDS_SETTING, "'example' can only appear once", a)
				}
			}
			for _, a := range attrs {
				if !ast.IsNil(a.Value) {
					v.n(diag.INVALID_RECORDS_SETTING_VALUE, "'example' cannot have a value", a)
				}
			}
			continue
		}
		for _, a := range attrs {
			v.n(diag.UNKNOWN_RECORDS_SETTING, "Unknown '"+name+"' setting", a)
		}
	}
	if ast.IsNil(e.Body) {
		return
	}
	fields, subs := partitionBody(elementBody(e))
	for _, row := range fields {
		if ast.IsNil(row.Callee) || len(row.Args) > 0 {
			v.n(diag.INVALID_RECORDS_FIELD, "Invalid record row structure", row)
			continue
		}
		vals := []ast.Node{row.Callee}
		if ce, ok := row.Callee.(*ast.CommaExpression); ok {
			vals = ce.Elements
		}
		for _, val := range vals {
			if !isValidRecordValue(val) {
				v.n(diag.INVALID_RECORDS_FIELD, "Records can only contain simple values (literals, null, true, false, or enum references). Complex expressions are not allowed.", val)
			}
		}
	}
	v.validateSubs(subs)
}

// ---------------------------------------------------------------- metadata

var metadataTargetKinds = []string{"Table", "Column", "TableGroup", "Note"}

// metadataTargetKind returns the canonical target kind of a Metadata block
// ("Table", "Column", "TableGroup", "Note") or "".
func metadataTargetKind(e *ast.ElementDeclaration) string {
	if e.TargetKind == nil {
		return ""
	}
	for _, k := range metadataTargetKinds {
		if strings.EqualFold(k, e.TargetKind.Value) {
			return k
		}
	}
	return ""
}

// validateMetadata checks a `Metadata <Kind> <name> { key: value }` block.
func (v *validator) validateMetadata(e *ast.ElementDeclaration) {
	kind := metadataTargetKind(e)
	if kind == "" {
		if e.TargetKind != nil {
			v.diags = append(v.diags, tokDiag(diag.INVALID_METADATA_TARGET_KIND, "A Metadata target kind must be one of: Table, Column, TableGroup, Note", e.TargetKind, v.u.path))
		} else {
			v.n(diag.INVALID_METADATA_TARGET_KIND, "A Metadata target kind must be one of: Table, Column, TableGroup, Note", e)
		}
	}
	if ast.IsNil(e.Name) {
		v.n(diag.NAME_NOT_FOUND, "A Metadata target must have a name", e)
	} else if isWildcard(e.Name) {
		v.n(diag.INVALID_NAME, "Wildcard (*) is not allowed as a Metadata target's name", e.Name)
	} else if !isValidName(e.Name) {
		v.n(diag.INVALID_NAME, "A Metadata target's name must be of the form <schema>.<table>, <table>.<column>, <schema>.<table>.<column> or a single identifier", e.Name)
	}
	fields, subs, ok := v.blockBody(e, "A Metadata's body must be a block")
	if !ok {
		return
	}
	for _, f := range fields {
		v.n(diag.INVALID_METADATA_FIELD, "A Metadata field must use the 'key: value' syntax", f)
	}
	byKey := map[string][]*ast.ElementDeclaration{}
	var order []string
	for _, sub := range subs {
		if sub.Type == nil {
			continue
		}
		if b, ok := sub.Body.(*ast.BlockExpression); ok {
			v.n(diag.UNEXPECTED_COMPLEX_BODY, "A Metadata field can only have an inline value", b)
			continue
		}
		k := strings.ToLower(sub.Type.Value)
		if _, ok := byKey[k]; !ok {
			order = append(order, k)
		}
		byKey[k] = append(byKey[k], sub)
		var callee ast.Node
		if fa, ok := sub.Body.(*ast.FunctionApplication); ok {
			callee = fa.Callee
		}
		if spec, ok := builtinMetadataSpec(kind, k); ok {
			if !spec.valid(callee) {
				v.n(diag.INVALID_METADATA_FIELD, spec.message, sub)
			}
			continue
		}
		if !isValidMetadataValue(callee) {
			var target ast.Node = sub
			if !ast.IsNil(sub.Body) {
				target = sub.Body
			}
			v.n(diag.INVALID_METADATA_FIELD, "A Metadata field value must be a string or a color literal", target)
		}
	}
	for _, k := range order {
		if len(byKey[k]) > 1 {
			for _, s := range byKey[k] {
				v.n(diag.DUPLICATE_METADATA_FIELD, "Duplicate Metadata fields are defined", s)
			}
		}
	}
}

type metadataSpec struct {
	valid   func(ast.Node) bool
	message string
}

// builtinMetadataSpec returns the validation rule for a built-in key of a
// metadata target kind (e.g. a Table's note or headercolor).
func builtinMetadataSpec(kind, key string) (metadataSpec, bool) {
	switch kind {
	case "Table":
		switch key {
		case settingNote:
			return metadataSpec{isQuotedString, "'note' must be a string literal"}, true
		case settingHeaderColor:
			return metadataSpec{isValidHexColor, "'headercolor' must be a color literal"}, true
		}
	case "Column":
		if key == settingNote {
			return metadataSpec{isQuotedString, "'note' must be a quoted string"}, true
		}
	case "TableGroup":
		switch key {
		case settingNote:
			return metadataSpec{isQuotedString, "'note' must be a string literal"}, true
		case settingColor:
			return metadataSpec{isValidHexColor, "'color' must be a color literal"}, true
		}
	case "Note":
		if key == settingColor {
			return metadataSpec{isValidColorOrNone, "'color' must be a color literal or 'none'"}, true
		}
	}
	return metadataSpec{}, false
}

// ---------------------------------------------------------------- diagram view

// validateDiagramView checks a DiagramView and its Tables/Notes/
// TableGroups/Schemas blocks.
func (v *validator) validateDiagramView(e *ast.ElementDeclaration) {
	v.topLevelOnly(e, diag.INVALID_DIAGRAMVIEW_CONTEXT, "DiagramView must appear top-level")
	if ast.IsNil(e.Name) {
		v.n(diag.NAME_NOT_FOUND, "A DiagramView must have a name", e)
	} else if !isSimpleName(e.Name) {
		v.n(diag.INVALID_NAME, "A DiagramView name must be a single identifier", e.Name)
	}
	v.noAlias(e, "A DiagramView shouldn't have an alias")
	s, errs := aggregateSettings(e.AttributeList, v.u.path)
	v.add(errs)
	for _, name := range s.keys {
		for _, a := range s.get(name) {
			v.n(diag.UNEXPECTED_SETTINGS, "Unknown '"+name+"' setting for DiagramView", a)
		}
	}
	fields, subs, ok := v.blockBody(e, "A DiagramView's body must be a block")
	if !ok {
		return
	}
	for _, f := range fields {
		if isWildcard(f.Callee) {
			for _, a := range f.Args {
				v.n(diag.INVALID_DIAGRAMVIEW_FIELD, "DiagramView field should only have a single name", a)
			}
			continue
		}
		v.n(diag.INVALID_DIAGRAMVIEW_FIELD, "Fields are not allowed at DiagramView level. Use Tables, Notes, TableGroups, or Schemas blocks instead.", f)
	}
	for _, sub := range subs {
		if sub.Type == nil {
			continue
		}
		switch strings.ToLower(sub.Type.Value) {
		case kDVTables, kDVNotes, kDVTableGroups, kDVSchemas:
		default:
			v.n(diag.INVALID_DIAGRAMVIEW_FIELD, "Unknown block type \""+sub.Type.Value+"\" in DiagramView. Allowed: Tables, Notes, TableGroups, Schemas", sub)
			continue
		}
		if sub.BodyColon != nil {
			v.diags = append(v.diags, tokDiag(diag.UNEXPECTED_SIMPLE_BODY, sub.Type.Value+" body must be a block", sub.BodyColon, v.u.path))
		}
	}
}

// ---------------------------------------------------------------- custom / use

// validateCustom checks a free-form `key: 'value'` element, allowed only
// inside Project and Dep.
func (v *validator) validateCustom(e *ast.ElementDeclaration) {
	if !isElem(parentElem(e), kProject, kDep) {
		v.n(diag.INVALID_CUSTOM_CONTEXT, "A Custom element can only appear in a Project or a Dep", e)
	}
	if !ast.IsNil(e.Name) {
		v.n(diag.UNEXPECTED_NAME, "A Custom element shouldn't have a name", e.Name)
	}
	if !ast.IsNil(e.Alias) {
		v.n(diag.UNEXPECTED_NAME, "A Custom element shouldn't have an alias", e.Alias)
	}
	v.noSettings(e, "A Custom element shouldn't have a setting list")
	switch b := e.Body.(type) {
	case *ast.BlockExpression:
		v.n(diag.UNEXPECTED_COMPLEX_BODY, "A Custom element can only have an inline field", b)
	case *ast.FunctionApplication:
		if !isElem(parentElem(e), kDep) && !isQuotedString(b.Callee) {
			v.n(diag.INVALID_CUSTOM_ELEMENT_VALUE, "A Custom element value can only be a string", b)
		}
		for _, a := range b.Args {
			v.n(diag.INVALID_CUSTOM_ELEMENT_VALUE, "A Custom element value can only be a string", a)
		}
	}
}

var importKinds = []string{"table", "enum", "tablepartial", "note", "schema", "tablegroup"}

// validateUse checks a use/reuse declaration and its specifiers.
func (v *validator) validateUse(u *ast.UseDeclaration) {
	if !isProgram(u.Base().Parent) {
		v.n(diag.INVALID_USE_CONTEXT, "'use' must appear at the top level", u)
	}
	l, ok := u.Specifiers.(*ast.UseSpecifierList)
	if !ok || l == nil {
		return
	}
	for _, s := range l.Specifiers {
		if s.ImportKind == nil {
			v.n(diag.INVALID_USE_SPECIFIER_KIND, "A use specifier must have a type (e.g. table, enum)", s)
		} else {
			valid := false
			for _, k := range importKinds {
				if strings.EqualFold(k, s.ImportKind.Value) {
					valid = true
				}
			}
			if !valid {
				v.diags = append(v.diags, tokDiag(diag.INVALID_USE_SPECIFIER_KIND, "'"+s.ImportKind.Value+"' is not a valid specifier type", s.ImportKind, v.u.path))
			}
		}
		// nodeFullname checks (local_modules/use)
		parts, ok := destructureComplexVariable(s.Name)
		kind := ""
		if s.ImportKind != nil {
			kind = strings.ToLower(s.ImportKind.Value)
		}
		switch {
		case !ok || len(parts) == 0:
			v.n(diag.INVALID_NAME, "Use specifiers must have a valid name (simple name or schema-qualified name)", s)
		case kind == "tablegroup" && len(parts) > 1:
			v.n(diag.INVALID_USE_SPECIFIER_NAME, "A TableGroup name must be a simple name", s)
		case kind == "tablepartial" && len(parts) > 1:
			v.n(diag.INVALID_USE_SPECIFIER_NAME, "A TablePartial name must be a simple name", s)
		case kind == "note" && len(parts) > 1:
			v.n(diag.INVALID_USE_SPECIFIER_NAME, "A Sticky note name must be a simple name", s)
		}
		// nodeAlias check
		if !ast.IsNil(s.Alias) && !isSimpleName(s.Alias) {
			v.n(diag.INVALID_ALIAS, "Use aliases can only contain alphanumeric and underscore unless surrounded by double quotes", s.Alias)
		}
	}
}
