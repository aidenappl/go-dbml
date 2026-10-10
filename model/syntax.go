// Portions of this file are derived from holistics/dbml packages/dbml-parse/src/core/utils/{expression,validate,interpret}.ts, src/core/types/relation.ts
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

package model

import (
	"strconv"
	"strings"

	"github.com/aidenappl/go-dbml/ast"
	"github.com/aidenappl/go-dbml/lexer"
)

// Element kinds (lower-cased element type keywords).
const (
	kTable              = "table"
	kEnum               = "enum"
	kRef                = "ref"
	kDep                = "dep"
	kNote               = "note"
	kProject            = "project"
	kIndexes            = "indexes"
	kTableGroup         = "tablegroup"
	kTablePartial       = "tablepartial"
	kChecks             = "checks"
	kRecords            = "records"
	kDiagramView        = "diagramview"
	kDVTables           = "tables"
	kDVNotes            = "notes"
	kDVTableGroups      = "tablegroups"
	kDVSchemas          = "schemas"
	kMetadata           = "metadata"
	defaultSchema       = "public"
	noneColor           = "none"
	settingColor        = "color"
	settingHeaderColor  = "headercolor"
	settingNote         = "note"
	settingPK           = "pk"
	settingPrimaryKey   = "primary key"
	settingUnique       = "unique"
	settingRef          = "ref"
	settingDep          = "dep"
	settingNotNull      = "not null"
	settingNull         = "null"
	settingIncrement    = "increment"
	settingDefault      = "default"
	settingName         = "name"
	settingType         = "type"
	settingCheck        = "check"
	settingUpdate       = "update"
	settingDelete       = "delete"
	settingInactive     = "inactive"
	settingExample      = "example"
	depDownstream       = "->"
	depUpstream         = "<-"
	elementKindNotFound = ""
)

var elementKinds = []string{kTable, kEnum, kRef, kDep, kNote, kProject, kIndexes, kTableGroup, kTablePartial, kChecks,
	kRecords, kDiagramView, kDVTables, kDVNotes, kDVTableGroups, kDVSchemas, kMetadata}

func elementKind(e *ast.ElementDeclaration) string {
	if e == nil || e.Type == nil {
		return elementKindNotFound
	}
	t := strings.ToLower(e.Type.Value)
	for _, k := range elementKinds {
		if t == k {
			return k
		}
	}
	return elementKindNotFound
}

func isElem(n ast.Node, kinds ...string) bool {
	e, ok := n.(*ast.ElementDeclaration)
	return ok && e != nil && e.IsKind(kinds...)
}

func isProgram(n ast.Node) bool {
	_, ok := n.(*ast.Program)
	return ok
}

// parentElem returns the nearest enclosing element declaration or program.
func parentElem(n ast.Node) ast.Node { return ast.ParentElement(n) }

func parentDecl(n ast.Node) *ast.ElementDeclaration { return ast.ParentElementDecl(n) }

// ---------------------------------------------------------------- predicates

func variableOf(n ast.Node) *lexer.Token {
	pe, ok := n.(*ast.PrimaryExpression)
	if !ok || pe == nil {
		return nil
	}
	v, ok := pe.Expression.(*ast.Variable)
	if !ok || v == nil || v.Variable == nil {
		return nil
	}
	return v.Variable
}

func literalOf(n ast.Node) *lexer.Token {
	pe, ok := n.(*ast.PrimaryExpression)
	if !ok || pe == nil {
		return nil
	}
	l, ok := pe.Expression.(*ast.Literal)
	if !ok || l == nil || l.Literal == nil {
		return nil
	}
	return l.Literal
}

// isVariableNode: a primary expression with a variable (identifier or quoted name).
func isVariableNode(n ast.Node) bool { return variableOf(n) != nil }

// isIdentifierNode: a primary expression with an identifier (not quoted).
func isIdentifierNode(n ast.Node) bool {
	t := variableOf(n)
	return t != nil && t.Kind == lexer.Identifier
}

// isSimpleName mirrors isSimpleName (primary expression with a variable node).
func isSimpleName(n ast.Node) bool {
	pe, ok := n.(*ast.PrimaryExpression)
	if !ok || pe == nil {
		return false
	}
	v, ok := pe.Expression.(*ast.Variable)
	return ok && v != nil
}

func isQuotedString(n ast.Node) bool {
	if t := variableOf(n); t != nil && t.Kind == lexer.QuotedString {
		return true
	}
	if t := literalOf(n); t != nil && t.Kind == lexer.StringLiteral {
		return true
	}
	return false
}

func extractQuotedString(n ast.Node) (string, bool) {
	if !isQuotedString(n) {
		return "", false
	}
	if t := variableOf(n); t != nil {
		return t.Value, true
	}
	return literalOf(n).Value, true
}

func extractVariable(n ast.Node) (string, bool) {
	if t := variableOf(n); t != nil {
		return t.Value, true
	}
	return "", false
}

func isAccess(n ast.Node) bool {
	in, ok := n.(*ast.InfixExpression)
	return ok && in != nil && !ast.IsNil(in.Left) && !ast.IsNil(in.Right) && in.Op != nil && in.Op.Value == "."
}

func asInfix(n ast.Node) *ast.InfixExpression {
	in, _ := n.(*ast.InfixExpression)
	return in
}

func isDotDelimitedIdentifier(n ast.Node) bool {
	if isVariableNode(n) {
		return true
	}
	if !isAccess(n) {
		return false
	}
	in := asInfix(n)
	return isVariableNode(in.Right) && isDotDelimitedIdentifier(in.Left)
}

func isSignedNumber(n ast.Node) bool {
	if pre, ok := n.(*ast.PrefixExpression); ok && pre != nil {
		if pre.Op == nil || (pre.Op.Value != "-" && pre.Op.Value != "+") {
			return false
		}
		return isSignedNumber(pre.Expression)
	}
	t := literalOf(n)
	return t != nil && t.Kind == lexer.NumericLiteral
}

func numberText(n ast.Node) string {
	if pre, ok := n.(*ast.PrefixExpression); ok && pre != nil {
		op := ""
		if pre.Op != nil {
			op = pre.Op.Value
		}
		return op + numberText(pre.Expression)
	}
	if t := literalOf(n); t != nil {
		return t.Value
	}
	return ""
}

// jsParseFloat mirrors Number.parseFloat on a numeric literal.
func jsParseFloat(s string) float64 {
	end := 0
	seenDigit, seenDot, seenExp := false, false, false
	for end < len(s) {
		c := s[end]
		switch {
		case c >= '0' && c <= '9':
			seenDigit = true
		case c == '.' && !seenDot && !seenExp:
			seenDot = true
		case (c == 'e' || c == 'E') && seenDigit && !seenExp:
			if end+1 < len(s) && (s[end+1] == '+' || s[end+1] == '-') {
				end++
			}
			seenExp = true
		default:
			goto done
		}
		end++
	}
done:
	v, err := strconv.ParseFloat(strings.TrimRight(s[:end], "eE+-"), 64)
	if err != nil {
		return 0
	}
	return v
}

// jsNumber mirrors Number(str) for numeric literal text.
func jsNumber(s string) (float64, bool) {
	t := strings.TrimSpace(s)
	if t == "" {
		return 0, true
	}
	v, err := strconv.ParseFloat(t, 64)
	if err != nil {
		if ne, ok := err.(*strconv.NumError); ok && ne.Err == strconv.ErrRange {
			return v, true
		}
		return 0, false
	}
	return v, true
}

func extractNumber(n ast.Node) float64 {
	if pre, ok := n.(*ast.PrefixExpression); ok && pre != nil {
		if pre.Op != nil && pre.Op.Value == "-" {
			return -extractNumber(pre.Expression)
		}
		return extractNumber(pre.Expression)
	}
	if t := literalOf(n); t != nil {
		return jsParseFloat(t.Value)
	}
	return 0
}

func extractNumericLiteral(n ast.Node) (float64, bool) {
	if t := literalOf(n); t != nil && t.Kind == lexer.NumericLiteral {
		v, ok := jsNumber(t.Value)
		return v, ok
	}
	return 0, false
}

func isValidHexColor(n ast.Node) bool {
	t := literalOf(n)
	if t == nil || t.Kind != lexer.ColorLiteral {
		return false
	}
	c := []rune(t.Value)
	if len(c) != 4 && len(c) != 7 {
		return false
	}
	if c[0] != '#' {
		return false
	}
	for _, r := range c[1:] {
		if !lexer.IsHexChar(r) {
			return false
		}
	}
	return true
}

func isValidColorOrNone(n ast.Node) bool {
	if isValidHexColor(n) {
		return true
	}
	t := variableOf(n)
	return t != nil && t.Kind == lexer.Identifier && strings.ToLower(t.Value) == noneColor
}

func extractColor(n ast.Node) (string, bool) {
	if isValidHexColor(n) {
		return literalOf(n).Value, true
	}
	if isIdentifierNode(n) {
		if strings.ToLower(variableOf(n).Value) == noneColor {
			return noneColor, true
		}
	}
	return "", false
}

func isValidMetadataValue(n ast.Node) bool { return isQuotedString(n) || isValidHexColor(n) }

func extractMetadataValue(n ast.Node) (string, bool) {
	if ast.IsNil(n) {
		return "", false
	}
	if s, ok := extractQuotedString(n); ok {
		return s, true
	}
	return extractColor(n)
}

func isValidPartialInjection(n ast.Node) bool {
	pre, ok := n.(*ast.PrefixExpression)
	return ok && pre != nil && pre.Op != nil && pre.Op.Value == "~" && isVariableNode(pre.Expression)
}

func isWildcard(n ast.Node) bool {
	w, ok := n.(*ast.Wildcard)
	return ok && w != nil
}

// destructureMemberAccess splits a.b.c into [a, b, c].
func destructureMemberAccess(n ast.Node) []ast.Node {
	if ast.IsNil(n) {
		return nil
	}
	if !isAccess(n) {
		return []ast.Node{n}
	}
	in := asInfix(n)
	frags := destructureMemberAccess(in.Left)
	if frags == nil {
		return nil
	}
	return append(frags, in.Right)
}

// destructureComplexVariable returns the variable names of a.b.c, or ok=false.
func destructureComplexVariable(n ast.Node) ([]string, bool) {
	if ast.IsNil(n) {
		return nil, false
	}
	frags := destructureMemberAccess(n)
	if frags == nil {
		return nil, false
	}
	out := make([]string, 0, len(frags))
	for _, f := range frags {
		v, ok := extractVariable(f)
		if !ok {
			return nil, false
		}
		out = append(out, v)
	}
	return out, true
}

func isValidName(n ast.Node) bool {
	v, ok := destructureComplexVariable(n)
	return ok && len(v) > 0
}

func isTupleOfVariables(n ast.Node) bool {
	t, ok := n.(*ast.TupleExpression)
	if !ok || t == nil {
		return false
	}
	for _, e := range t.Elements {
		if !isVariableNode(e) {
			return false
		}
	}
	return true
}

// complexVariableTuple is the result of destructuring `a.b.(c, d)` or `a.b.c`.
type complexVariableTuple struct {
	variables     []ast.Node
	tupleElements []ast.Node
}

func destructureComplexVariableTuple(n ast.Node) (complexVariableTuple, bool) {
	if ast.IsNil(n) {
		return complexVariableTuple{}, false
	}
	frags := destructureMemberAccess(n)
	if len(frags) == 0 {
		return complexVariableTuple{}, false
	}
	var tuple []ast.Node
	if !isVariableNode(frags[len(frags)-1]) {
		top := frags[len(frags)-1]
		frags = frags[:len(frags)-1]
		if !isTupleOfVariables(top) {
			return complexVariableTuple{}, false
		}
		tuple = top.(*ast.TupleExpression).Elements
	}
	for _, f := range frags {
		if !isVariableNode(f) {
			return complexVariableTuple{}, false
		}
	}
	return complexVariableTuple{variables: frags, tupleElements: tuple}, true
}

// destructureCallExpression splits schema.table(col1, col2).
func destructureCallExpression(n ast.Node) (vars []ast.Node, args []ast.Node, ok bool) {
	call, isCall := n.(*ast.CallExpression)
	if !isCall || call == nil || ast.IsNil(call.Callee) {
		return nil, nil, false
	}
	frags := destructureMemberAccess(call.Callee)
	if len(frags) == 0 {
		return nil, nil, false
	}
	for _, f := range frags {
		if !isVariableNode(f) {
			return nil, nil, false
		}
	}
	if call.ArgumentList != nil && isTupleOfVariables(call.ArgumentList) {
		args = append(args, call.ArgumentList.Elements...)
	}
	return frags, args, true
}

func isValidIndexName(n ast.Node) bool {
	if isSimpleName(n) {
		return true
	}
	_, ok := n.(*ast.FunctionExpression)
	return ok
}

type indexNodeParts struct {
	functional    []*ast.FunctionExpression
	nonFunctional []ast.Node
}

func destructureIndexNode(n ast.Node) (indexNodeParts, bool) {
	if isValidIndexName(n) {
		if fe, ok := n.(*ast.FunctionExpression); ok {
			return indexNodeParts{functional: []*ast.FunctionExpression{fe}}, true
		}
		return indexNodeParts{nonFunctional: []ast.Node{n}}, true
	}
	if t, ok := n.(*ast.TupleExpression); ok && t != nil {
		for _, e := range t.Elements {
			if !isValidIndexName(e) {
				return indexNodeParts{}, false
			}
		}
		var parts indexNodeParts
		for _, e := range t.Elements {
			if fe, ok := e.(*ast.FunctionExpression); ok {
				parts.functional = append(parts.functional, fe)
			}
		}
		for _, e := range t.Elements {
			if isVariableNode(e) {
				parts.nonFunctional = append(parts.nonFunctional, e)
			}
		}
		return parts, true
	}
	return indexNodeParts{}, false
}

var relationshipOps = map[string][2]string{
	"-":    {"1", "1"},
	"<>":   {"*", "*"},
	">":    {"*", "1"},
	"<":    {"1", "*"},
	"-?":   {"1", "0..1"},
	"?-":   {"0..1", "1"},
	"?-?":  {"0..1", "0..1"},
	"?>":   {"0..*", "1"},
	">?":   {"*", "0..1"},
	"?>?":  {"0..*", "0..1"},
	"?<":   {"0..1", "*"},
	"<?":   {"1", "0..*"},
	"?<?":  {"0..1", "0..*"},
	"?<>":  {"0..*", "*"},
	"<>?":  {"*", "0..*"},
	"?<>?": {"0..*", "0..*"},
}

func isRelationshipOp(op string) bool {
	_, ok := relationshipOps[op]
	return ok
}

func isDependencyOp(op string) bool { return op == depUpstream || op == depDownstream }

func isUnaryRelationship(n ast.Node) bool {
	pre, ok := n.(*ast.PrefixExpression)
	if !ok || pre == nil || pre.Op == nil || !isRelationshipOp(pre.Op.Value) {
		return false
	}
	v, ok := destructureComplexVariable(pre.Expression)
	return ok && len(v) > 0
}

func isUnaryDependency(n ast.Node) bool {
	pre, ok := n.(*ast.PrefixExpression)
	if !ok || pre == nil || pre.Op == nil || !isDependencyOp(pre.Op.Value) {
		return false
	}
	v, ok := destructureComplexVariable(pre.Expression)
	return ok && len(v) > 0
}

func isBinaryRelationship(n ast.Node) bool {
	in, ok := n.(*ast.InfixExpression)
	if !ok || in == nil || in.Op == nil || !isRelationshipOp(in.Op.Value) {
		return false
	}
	_, lok := destructureComplexVariableTuple(in.Left)
	_, rok := destructureComplexVariableTuple(in.Right)
	return lok && rok
}

func isEqualTupleOperands(n ast.Node) bool {
	in, ok := n.(*ast.InfixExpression)
	if !ok || in == nil {
		return false
	}
	l, lok := destructureComplexVariableTuple(in.Left)
	r, rok := destructureComplexVariableTuple(in.Right)
	if !lok || !rok {
		return false
	}
	return len(l.tupleElements) == len(r.tupleElements)
}

func isValidDefaultValue(n ast.Node) bool {
	if ast.IsNil(n) {
		return false
	}
	if isQuotedString(n) || isSignedNumber(n) {
		return true
	}
	if isIdentifierNode(n) {
		switch strings.ToLower(variableOf(n).Value) {
		case "true", "false", "null":
			return true
		}
	}
	if pre, ok := n.(*ast.PrefixExpression); ok && pre != nil && pre.Op != nil && (pre.Op.Value == "-" || pre.Op.Value == "+") && isSignedNumber(pre.Expression) {
		return true
	}
	if _, ok := n.(*ast.FunctionExpression); ok {
		return true
	}
	if !isDotDelimitedIdentifier(n) {
		return false
	}
	frags := destructureMemberAccess(n)
	return len(frags) == 2 || len(frags) == 3
}

func isValidColumnType(n ast.Node) bool {
	switch n.(type) {
	case *ast.CallExpression, *ast.PrimaryExpression, *ast.Array:
	default:
		if !isAccess(n) {
			return false
		}
	}
	for {
		switch v := n.(type) {
		case *ast.CallExpression:
			if ast.IsNil(v.Callee) || v.ArgumentList == nil {
				return false
			}
			for _, e := range v.ArgumentList.Elements {
				if !(isSignedNumber(e) || isQuotedString(e) || isIdentifierNode(e)) {
					return false
				}
			}
			n = v.Callee
			continue
		case *ast.Array:
			if ast.IsNil(v.Array) || v.Indexer == nil {
				return false
			}
			for _, a := range v.Indexer.Elements {
				if a.Colon != nil || !ast.IsNil(a.Value) || !isSignedNumber(a.Name) {
					return false
				}
			}
			n = v.Array
			continue
		}
		break
	}
	vars, ok := destructureComplexVariable(n)
	return ok && len(vars) > 0
}

func getRightmostVariable(n ast.Node) ast.Node {
	if isVariableNode(n) {
		return n
	}
	if isAccess(n) {
		r := asInfix(n).Right
		if isVariableNode(r) {
			return r
		}
	}
	return nil
}

// elementBody returns the body items of a declaration.
func elementBody(e *ast.ElementDeclaration) []ast.Node {
	if e == nil || ast.IsNil(e.Body) {
		return nil
	}
	if b, ok := e.Body.(*ast.BlockExpression); ok {
		return b.Body
	}
	return []ast.Node{e.Body}
}

func hasSimpleBody(e *ast.ElementDeclaration) bool { return e.BodyColon != nil }

// partitionBody splits body items into function applications and sub-elements.
func partitionBody(items []ast.Node) (fields []*ast.FunctionApplication, subs []*ast.ElementDeclaration) {
	for _, it := range items {
		switch v := it.(type) {
		case *ast.FunctionApplication:
			fields = append(fields, v)
		case *ast.ElementDeclaration:
			subs = append(subs, v)
		}
	}
	return
}

// elementNameParts mirrors extractElementName.
func elementNameParts(n ast.Node) (schema []string, name string) {
	frags, ok := destructureComplexVariable(n)
	if !ok || len(frags) == 0 {
		return nil, ""
	}
	return frags[:len(frags)-1], frags[len(frags)-1]
}

// normalizeNote strips leading blank lines and the common indentation.
func normalizeNote(content string) string {
	lines := strings.Split(content, "\n")
	first := -1
	for i, l := range lines {
		if jsTrimStart(l) != "" {
			first = i
			break
		}
	}
	var trimmed []string
	if first < 0 {
		// JS: lines.slice(-1) when findIndex returns -1
		trimmed = lines[len(lines)-1:]
	} else {
		trimmed = lines[first:]
	}
	var nonEmpty []string
	for _, l := range trimmed {
		if jsTrimStart(l) != "" {
			nonEmpty = append(nonEmpty, l)
		}
	}
	if len(nonEmpty) == 0 {
		return strings.Join(trimmed, "\n")
	}
	minIndent := -1
	for _, l := range nonEmpty {
		ind := utf16Len(l) - utf16Len(jsTrimStart(l))
		if minIndent < 0 || ind < minIndent {
			minIndent = ind
		}
	}
	out := make([]string, len(trimmed))
	for i, l := range trimmed {
		out[i] = utf16Slice(l, minIndent)
	}
	return strings.Join(out, "\n")
}

func jsTrimStart(s string) string {
	return strings.TrimLeftFunc(s, isJSWhitespace)
}

func isJSWhitespace(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0xA0, 0x1680, 0x2028, 0x2029, 0x202F, 0x205F, 0x3000, 0xFEFF:
		return true
	}
	return r >= 0x2000 && r <= 0x200A
}

func utf16Len(s string) int {
	n := 0
	for _, r := range s {
		if r >= 0x10000 {
			n += 2
		} else {
			n++
		}
	}
	return n
}

// utf16Slice returns s[from:] where from counts UTF-16 units.
func utf16Slice(s string, from int) string {
	n := 0
	for i, r := range s {
		if n >= from {
			return s[i:]
		}
		if r >= 0x10000 {
			n += 2
		} else {
			n++
		}
	}
	return ""
}
