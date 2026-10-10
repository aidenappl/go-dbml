// Portions of this file are derived from holistics/dbml packages/dbml-parse/src/core/utils/validate.ts (aggregateSettingList), src/core/local_modules/{table,ref,indexes,dep,metadata}/validate.ts
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

package model

import (
	"strings"

	"github.com/aidenappl/go-dbml/ast"
	"github.com/aidenappl/go-dbml/diag"
)

// settings is an aggregated setting list: lower-cased name -> attributes,
// preserving first-seen key order.
type settings struct {
	keys []string
	vals [][]*ast.Attribute
	// idx indexes keys once there are many (settings lists are usually a
	// handful of entries; a linear scan avoids a map per column)
	idx map[string]int
}

// newSettings returns an empty setting map.
func newSettings() *settings { return &settings{} }

func (s *settings) find(name string) int {
	if s.idx != nil {
		if i, ok := s.idx[name]; ok {
			return i
		}
		return -1
	}
	for i, k := range s.keys {
		if k == name {
			return i
		}
	}
	return -1
}

// add appends an attribute under name, remembering first-seen key order.
func (s *settings) add(name string, a *ast.Attribute) {
	i := s.find(name)
	if i < 0 {
		i = len(s.keys)
		s.keys = append(s.keys, name)
		s.vals = append(s.vals, nil)
		if s.idx != nil {
			s.idx[name] = i
		} else if len(s.keys) > 16 {
			s.idx = make(map[string]int, len(s.keys))
			for j, k := range s.keys {
				s.idx[k] = j
			}
		}
	}
	s.vals[i] = append(s.vals[i], a)
}

// get returns every attribute named name.
func (s *settings) get(name string) []*ast.Attribute {
	if i := s.find(name); i >= 0 {
		return s.vals[i]
	}
	return nil
}

// has reports whether name appears at least once.
func (s *settings) has(name string) bool { return len(s.get(name)) > 0 }

// first returns the first attribute named name, or nil.
func (s *settings) first(name string) *ast.Attribute {
	if a := s.get(name); len(a) > 0 {
		return a[0]
	}
	return nil
}

// firstValue returns the value of the first attribute named name, or nil.
func (s *settings) firstValue(name string) ast.Node {
	if a := s.first(name); a != nil {
		return a.Value
	}
	return nil
}

// aggregateSettings mirrors aggregateSettingList.
func aggregateSettings(list *ast.ListExpression, file string) (*settings, []diag.Diagnostic) {
	s := newSettings()
	if list == nil {
		return s, nil
	}
	s.keys = make([]string, 0, len(list.Elements))
	s.vals = make([][]*ast.Attribute, 0, len(list.Elements))
	var errs []diag.Diagnostic
	for _, attr := range list.Elements {
		if attr == nil || ast.IsNil(attr.Name) {
			continue
		}
		if pe, ok := attr.Name.(*ast.PrimaryExpression); ok {
			errs = append(errs, nodeDiag(diag.INVALID_SETTINGS, "A setting name must be a stream of identifiers", pe, file))
			continue
		}
		is, ok := attr.Name.(*ast.IdentifierStream)
		if !ok {
			continue
		}
		name := strings.ToLower(is.String())
		if name == "" {
			continue
		}
		s.add(name, attr)
	}
	return s, errs
}

// attrErrNode is attr.value || attr.name (|| attr).
func attrErrNode(a *ast.Attribute) ast.Node {
	if !ast.IsNil(a.Value) {
		return a.Value
	}
	if !ast.IsNil(a.Name) {
		return a.Name
	}
	return a
}

type errSink struct {
	file  string
	diags []diag.Diagnostic
}

// n records a diagnostic over node (ignored when node is nil).
func (e *errSink) n(code diag.Code, msg string, node ast.Node) {
	if ast.IsNil(node) {
		return
	}
	e.diags = append(e.diags, nodeDiag(code, msg, node, e.file))
}

// dupes reports every attribute when a setting appears more than once.
func (e *errSink) dupes(attrs []*ast.Attribute, code diag.Code, msg string) {
	if len(attrs) > 1 {
		for _, a := range attrs {
			e.n(code, msg, a)
		}
	}
}

// customMetadata validates a free-form custom property: it may appear
// once and its value must be a string or a color literal.
func (e *errSink) customMetadata(name string, attrs []*ast.Attribute, dup, invalid diag.Code) {
	e.dupes(attrs, dup, "'"+name+"' can only appear once")
	for _, a := range attrs {
		if !isValidMetadataValue(a.Value) {
			e.n(invalid, "Custom setting '"+name+"' must be a string or a color literal", attrErrNode(a))
		}
	}
}

// legacyAttr wraps a legacy inline setting (`id int pk`) as an attribute
// without a value, keeping the primary expression's span.
func legacyAttr(n ast.Node) *ast.Attribute {
	return &ast.Attribute{NodeBase: *n.Base(), Name: n}
}

// columnSettings mirrors validateFieldSetting (table and partial columns).
// forPartial selects the TablePartial message variants and disables `dep`.
func columnSettings(parts []ast.Node, file string, forPartial bool) (*settings, []diag.Diagnostic) {
	sink := &errSink{file: file}
	ok := true
	if len(parts) > 0 {
		for _, p := range parts[:len(parts)-1] {
			if !isIdentifierNode(p) {
				ok = false
			}
		}
		last := parts[len(parts)-1]
		if _, isList := last.(*ast.ListExpression); !isIdentifierNode(last) && !isList {
			ok = false
		}
	}
	if !ok {
		for _, p := range parts {
			sink.n(diag.INVALID_COLUMN, "These fields must be some inline settings optionally ended with a setting list", p)
		}
		return newSettings(), sink.diags
	}
	if len(parts) == 0 {
		return newSettings(), nil
	}
	parts = append([]ast.Node{}, parts...)
	var list *ast.ListExpression
	if l, isList := parts[len(parts)-1].(*ast.ListExpression); isList {
		list = l
		parts = parts[:len(parts)-1]
	}
	s, errs := aggregateSettings(list, file)
	sink.diags = append(sink.diags, errs...)
	for _, p := range parts {
		name, _ := extractVariable(p)
		name = strings.ToLower(name)
		if name != settingPK && name != settingUnique {
			sink.n(diag.INVALID_SETTINGS, "Inline column settings can only be `pk` or `unique`", p)
			continue
		}
		s.add(name, legacyAttr(p))
	}
	if len(s.get(settingPK)) >= 1 && len(s.get(settingPrimaryKey)) >= 1 {
		for _, a := range append(append([]*ast.Attribute{}, s.get(settingPK)...), s.get(settingPrimaryKey)...) {
			sink.n(diag.DUPLICATE_COLUMN_SETTING, "Either one of 'primary key' and 'pk' can appear", a)
		}
	}
	q := func(name string) string {
		if forPartial {
			return "'" + name + "'"
		}
		return "'" + name + "'"
	}
	for _, name := range s.keys {
		attrs := s.get(name)
		noValue := func() {
			for _, a := range attrs {
				if !ast.IsNil(a.Value) {
					sink.n(diag.INVALID_COLUMN_SETTING_VALUE, q(name)+" must not have a value", attrErrNode(a))
				}
			}
		}
		switch name {
		case settingNote:
			sink.dupes(attrs, diag.DUPLICATE_COLUMN_SETTING, "note can only appear once")
			for _, a := range attrs {
				if !isQuotedString(a.Value) {
					sink.n(diag.INVALID_COLUMN_SETTING_VALUE, "'note' must be a quoted string", attrErrNode(a))
				}
			}
		case settingRef:
			for _, a := range attrs {
				if !isUnaryRelationship(a.Value) {
					sink.n(diag.INVALID_COLUMN_SETTING_VALUE, "'ref' must be a valid unary relationship", attrErrNode(a))
				}
			}
		case settingPrimaryKey, settingPK, settingNull, settingUnique, settingIncrement:
			sink.dupes(attrs, diag.DUPLICATE_COLUMN_SETTING, q(name)+" can only appear once")
			noValue()
		case settingNotNull:
			sink.dupes(attrs, diag.DUPLICATE_COLUMN_SETTING, "'not null' can only appear once")
			if nulls := s.get(settingNull); len(attrs) >= 1 && len(nulls) >= 1 {
				for _, a := range append(append([]*ast.Attribute{}, attrs...), nulls...) {
					sink.n(diag.CONFLICTING_SETTING, "'not null' and 'null' can not be set at the same time", a)
				}
			}
			noValue()
		case settingDefault:
			sink.dupes(attrs, diag.DUPLICATE_COLUMN_SETTING, "'default' can only appear once")
			for _, a := range attrs {
				if !isValidDefaultValue(a.Value) {
					sink.n(diag.INVALID_COLUMN_SETTING_VALUE, "'default' must be an enum value, a string literal, number literal, function expression, true, false or null", attrErrNode(a))
				}
			}
		case settingCheck:
			for _, a := range attrs {
				if _, ok := a.Value.(*ast.FunctionExpression); !ok {
					sink.n(diag.INVALID_COLUMN_SETTING_VALUE, "'check' must be a function expression", attrErrNode(a))
				}
			}
		case settingDep:
			if !forPartial {
				for _, a := range attrs {
					if !isUnaryDependency(a.Value) {
						sink.n(diag.INVALID_COLUMN_SETTING_VALUE, "'dep' must be a valid unary relationship", attrErrNode(a))
					}
				}
				continue
			}
			sink.customMetadata(name, attrs, diag.DUPLICATE_COLUMN_SETTING, diag.INVALID_COLUMN_SETTING_VALUE)
		default:
			sink.customMetadata(name, attrs, diag.DUPLICATE_COLUMN_SETTING, diag.INVALID_COLUMN_SETTING_VALUE)
		}
	}
	return s, sink.diags
}

// tableSettings validates a Table header setting list (headercolor, note,
// dep, custom properties) (validateTableSettings).
func tableSettings(list *ast.ListExpression, file string) (*settings, []diag.Diagnostic) {
	s, errs := aggregateSettings(list, file)
	sink := &errSink{file: file, diags: errs}
	for _, name := range s.keys {
		attrs := s.get(name)
		switch name {
		case settingHeaderColor:
			sink.dupes(attrs, diag.DUPLICATE_TABLE_SETTING, "'headercolor' can only appear once")
			for _, a := range attrs {
				if !isValidHexColor(a.Value) {
					sink.n(diag.INVALID_TABLE_SETTING_VALUE, "'headercolor' must be a color literal", attrErrNode(a))
				}
			}
		case settingNote:
			sink.dupes(attrs, diag.DUPLICATE_TABLE_SETTING, "'note' can only appear once")
			for _, a := range attrs {
				if !isQuotedString(a.Value) {
					sink.n(diag.INVALID_TABLE_SETTING_VALUE, "'note' must be a string literal", attrErrNode(a))
				}
			}
		case settingDep:
			for _, a := range attrs {
				if !isUnaryDependency(a.Value) {
					sink.n(diag.INVALID_TABLE_SETTING_VALUE, "'dep' must be `-> target` or `<- source`", attrErrNode(a))
				}
			}
		default:
			sink.customMetadata(name, attrs, diag.DUPLICATE_TABLE_SETTING, diag.INVALID_TABLE_SETTING_VALUE)
		}
	}
	return s, sink.diags
}

// isValidPolicy accepts the referential actions cascade, restrict,
// set null, set default and no action (identifiers, not quoted strings).
func isValidPolicy(n ast.Node) bool {
	var str string
	if is, ok := n.(*ast.IdentifierStream); ok && is != nil {
		str = is.String()
	} else if isIdentifierNode(n) {
		str = variableOf(n).Value
	} else {
		return false
	}
	switch strings.ToLower(str) {
	case "cascade", "no action", "set null", "set default", "restrict":
		return true
	}
	return false
}

// refFieldSettings validates a Ref field setting list (delete, update,
// color, inactive).
func refFieldSettings(list *ast.ListExpression, file string) (*settings, []diag.Diagnostic) {
	s, errs := aggregateSettings(list, file)
	sink := &errSink{file: file, diags: errs}
	for _, name := range s.keys {
		attrs := s.get(name)
		switch name {
		case settingDelete, settingUpdate:
			sink.dupes(attrs, diag.DUPLICATE_REF_SETTING, "'"+name+"' can only appear once")
			for _, a := range attrs {
				if !isValidPolicy(a.Value) {
					sink.n(diag.INVALID_REF_SETTING_VALUE, "'"+name+"' can only have values \"cascade\", \"no action\", \"set null\", \"set default\" or \"restrict\"", a)
				}
			}
		case settingColor:
			sink.dupes(attrs, diag.DUPLICATE_REF_SETTING, "'color' can only appear once")
			for _, a := range attrs {
				if !isValidHexColor(a.Value) {
					sink.n(diag.INVALID_REF_SETTING_VALUE, "'color' must be a color literal", a)
				}
			}
		case settingInactive:
			sink.dupes(attrs, diag.DUPLICATE_REF_SETTING, "'inactive' can only appear once")
			for _, a := range attrs {
				if !ast.IsNil(a.Value) {
					sink.n(diag.INVALID_REF_SETTING_VALUE, "'inactive' cannot have a value", a)
				}
			}
		default:
			for _, a := range attrs {
				sink.n(diag.UNKNOWN_REF_SETTING, "Unknown ref setting '"+name+"'", a)
			}
		}
	}
	return s, sink.diags
}

// indexFieldSettings validates an index setting list (name, note, unique,
// pk, type).
func indexFieldSettings(list *ast.ListExpression, file string) (*settings, []diag.Diagnostic) {
	s, errs := aggregateSettings(list, file)
	sink := &errSink{file: file, diags: errs}
	for _, name := range s.keys {
		attrs := s.get(name)
		switch name {
		case settingNote, settingName:
			sink.dupes(attrs, diag.DUPLICATE_INDEX_SETTING, "'"+name+"' can only appear once")
			for _, a := range attrs {
				if !isQuotedString(a.Value) {
					sink.n(diag.INVALID_INDEX_SETTING_VALUE, "'"+name+"' must be a string", a)
				}
			}
		case settingUnique, settingPK:
			sink.dupes(attrs, diag.DUPLICATE_INDEX_SETTING, "'"+name+"' can only appear once")
			for _, a := range attrs {
				if !ast.IsNil(a.Value) {
					sink.n(diag.INVALID_INDEX_SETTING_VALUE, "'"+name+"' must not have a value", a)
				}
			}
		case settingType:
			sink.dupes(attrs, diag.DUPLICATE_INDEX_SETTING, "'type' can only appear once")
			for _, a := range attrs {
				if !isVariableNode(a.Value) {
					sink.n(diag.INVALID_INDEX_SETTING_VALUE, "'type' must be \"btree\" or \"hash\"", a)
				}
			}
		default:
			for _, a := range attrs {
				sink.n(diag.UNKNOWN_INDEX_SETTING, "Unknown index setting '"+name+"'", a)
			}
		}
	}
	return s, sink.diags
}

// depSettings validates a Dep setting list: color, note, and custom keys
// whose values are strings, numbers, colors or identifiers.
func depSettings(list *ast.ListExpression, file string) (*settings, []diag.Diagnostic) {
	s, errs := aggregateSettings(list, file)
	sink := &errSink{file: file, diags: errs}
	for _, name := range s.keys {
		for _, a := range s.get(name) {
			switch name {
			case settingColor:
				if !isValidHexColor(a.Value) {
					sink.n(diag.INVALID_SETTINGS, "Invalid color value. Expected a hex color (e.g. #fff or #aabbcc)", a)
				}
			case settingNote:
				if !isQuotedString(a.Value) {
					sink.n(diag.INVALID_SETTINGS, "Invalid note value. Expected a quoted string", a)
				}
			default:
				if !ast.IsNil(a.Value) && !isQuotedString(a.Value) && !isValidHexColor(a.Value) && !isSignedNumber(a.Value) && !isSimpleName(a.Value) {
					sink.n(diag.INVALID_SETTINGS, "Invalid value for setting '"+name+"'. Expected a string, number, color, or identifier", a)
				}
			}
		}
	}
	return s, sink.diags
}

// fieldSettingList returns the trailing setting list of a column-like field.
func trailingList(args []ast.Node) *ast.ListExpression {
	if len(args) == 0 {
		return nil
	}
	l, _ := args[len(args)-1].(*ast.ListExpression)
	return l
}
