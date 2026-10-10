// Portions of this file are derived from holistics/dbml packages/dbml-parse/src/core/types/nodes.ts
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

// Package ast defines the concrete syntax tree produced by the DBML parser.
//
// The tree mirrors @dbml/parse's node set: DBML is parsed generically as
// element declarations whose bodies contain function applications
// (space-separated expressions) or nested element declarations. Meaning is
// assigned later by the model package. Every node records its span; spans are
// computed from the first and last child that has a valid position.
package ast

import (
	"reflect"
	"strings"

	"github.com/aidenappl/go-dbml/lexer"
)

// Kind classifies a node.
type Kind uint8

// Node kinds.
const (
	KindProgram Kind = iota
	KindElementDeclaration
	KindUseDeclaration
	KindUseSpecifier
	KindUseSpecifierList
	KindAttribute
	KindIdentifierStream
	KindLiteral
	KindVariable
	KindPrefixExpression
	KindInfixExpression
	KindPostfixExpression
	KindFunctionExpression
	KindFunctionApplication
	KindBlockExpression
	KindListExpression
	KindTupleExpression
	KindCallExpression
	KindPrimaryExpression
	KindGroupExpression
	KindCommaExpression
	KindEmpty
	KindArray
	KindWildcard
)

var kindNames = [...]string{
	KindProgram: "<program>", KindElementDeclaration: "<element-declaration>", KindUseDeclaration: "<use-declaration>",
	KindUseSpecifier: "<use-specifier>", KindUseSpecifierList: "<use-specifier-list>", KindAttribute: "<attribute>",
	KindIdentifierStream: "<identifier-stream>", KindLiteral: "<literal>", KindVariable: "<variable>",
	KindPrefixExpression: "<prefix-expression>", KindInfixExpression: "<infix-expression>",
	KindPostfixExpression: "<postfix-expression>", KindFunctionExpression: "<function-expression>",
	KindFunctionApplication: "<function-application>", KindBlockExpression: "<block-expression>",
	KindListExpression: "<list-expression>", KindTupleExpression: "<tuple-expression>",
	KindCallExpression: "<call-expression>", KindPrimaryExpression: "<primary-expression>",
	KindGroupExpression: "<group-expression>", KindCommaExpression: "<comma-expression>", KindEmpty: "<dummy>",
	KindArray: "<array>", KindWildcard: "<wildcard>",
}

// String returns the @dbml/parse name of the kind.
func (k Kind) String() string { return kindNames[k] }

// Node is implemented by every syntax node.
type Node interface {
	Element
	Base() *NodeBase
}

// NodeBase holds the data common to all nodes.
type NodeBase struct {
	ID   int32
	Kind Kind
	// HasPos is false when no child carried a valid position (the JavaScript
	// implementation represents this with NaN).
	HasPos    bool
	FullStart int32
	FullEnd   int32
	Filepath  string
	Start     lexer.Position
	End       lexer.Position
	Parent    Node
}

// Base implements Node.
func (b *NodeBase) Base() *NodeBase { return b }

// IsNil reports whether n is nil or a typed nil pointer.
func IsNil(n Node) bool {
	if n == nil {
		return true
	}
	v := reflect.ValueOf(n)
	return v.Kind() == reflect.Ptr && v.IsNil()
}

// Element is a child of a node: a *lexer.Token or a Node. Both report
// their range through Span.
type Element interface {
	Span() (start, end lexer.Position, fullStart, fullEnd int, ok bool)
}

type span = Element

// Span implements Element for nodes.
func (b *NodeBase) Span() (start, end lexer.Position, fullStart, fullEnd int, ok bool) {
	if b == nil || !b.HasPos {
		return
	}
	return b.Start, b.End, int(b.FullStart), int(b.FullEnd), true
}

func validSpan(s Element) (start, end lexer.Position, fullStart, fullEnd int, ok bool) {
	if s == nil {
		return
	}
	if n, isNode := s.(Node); isNode && IsNil(n) {
		return
	}
	return s.Span()
}

// Init computes the node's span from its children, given in source order.
func (b *NodeBase) Init(id int, kind Kind, filepath string, children ...span) {
	b.ID = int32(id)
	b.Kind = kind
	b.Filepath = filepath
	first := -1
	for i, c := range children {
		if s, _, fs, _, ok := validSpan(c); ok {
			b.Start = s
			b.FullStart = int32(fs)
			first = i
			break
		}
	}
	if first < 0 {
		return
	}
	for i := len(children) - 1; i >= 0; i-- {
		if _, e, _, fe, ok := validSpan(children[i]); ok {
			b.End = e
			b.FullEnd = int32(fe)
			break
		}
	}
	b.HasPos = true
}

// Program is the root node.
type Program struct {
	NodeBase
	Body   []Node // *UseDeclaration | *ElementDeclaration
	EOF    *lexer.Token
	Source string
}

// Declarations returns the top-level element declarations.
func (p *Program) Declarations() []*ElementDeclaration {
	var out []*ElementDeclaration
	for _, n := range p.Body {
		if e, ok := n.(*ElementDeclaration); ok && e != nil {
			out = append(out, e)
		}
	}
	return out
}

// Uses returns the top-level use/reuse declarations.
func (p *Program) Uses() []*UseDeclaration {
	var out []*UseDeclaration
	for _, n := range p.Body {
		if u, ok := n.(*UseDeclaration); ok && u != nil {
			out = append(out, u)
		}
	}
	return out
}

// UseDeclaration is `use { ... } from 'path'`, `use * from 'path'` or the reuse variants.
type UseDeclaration struct {
	NodeBase
	UseKeyword  *lexer.Token
	Specifiers  Node // *UseSpecifierList | *Wildcard
	FromKeyword *lexer.Token
	ImportPath  *lexer.Token
}

// IsReuse reports whether this is a `reuse` declaration.
func (u *UseDeclaration) IsReuse() bool {
	return u.UseKeyword != nil && strings.EqualFold(u.UseKeyword.Value, "reuse")
}

// UseSpecifier is `<kind> <name> [as <alias>]`.
type UseSpecifier struct {
	NodeBase
	ImportKind *lexer.Token
	Name       Node
	AsKeyword  *lexer.Token
	Alias      Node
}

// UseSpecifierList is the braced list of specifiers.
type UseSpecifierList struct {
	NodeBase
	OpenBrace  *lexer.Token
	Specifiers []*UseSpecifier
	CloseBrace *lexer.Token
}

// ElementDeclaration is `<type> [<target-kind>] [<name>] [as <alias>] [<settings>] (: <body> | { <body> })`.
type ElementDeclaration struct {
	NodeBase
	Type          *lexer.Token
	TargetKind    *lexer.Token // only for Metadata blocks
	Name          Node
	As            *lexer.Token
	Alias         Node
	AttributeList *ListExpression
	BodyColon     *lexer.Token
	Body          Node // *FunctionApplication | *BlockExpression
}

// IsKind reports whether the element type matches one of kinds (case-insensitive).
func (e *ElementDeclaration) IsKind(kinds ...string) bool {
	if e == nil || e.Type == nil {
		return false
	}
	t := strings.ToLower(e.Type.Value)
	for _, k := range kinds {
		if t == k {
			return true
		}
	}
	return false
}

// Attribute is a `name: value` or bare `name` inside a setting list.
type Attribute struct {
	NodeBase
	Name  Node // *IdentifierStream | *PrimaryExpression
	Colon *lexer.Token
	Value Node
}

// IdentifierStream is a run of identifiers such as `primary key` or `no action`.
type IdentifierStream struct {
	NodeBase
	Identifiers []*lexer.Token
}

// String joins the identifiers with single spaces.
func (s *IdentifierStream) String() string {
	parts := make([]string, len(s.Identifiers))
	for i, t := range s.Identifiers {
		parts[i] = t.Value
	}
	return strings.Join(parts, " ")
}

// Literal wraps a number, string or color token.
type Literal struct {
	NodeBase
	Literal *lexer.Token
}

// Variable wraps an identifier or double-quoted name.
type Variable struct {
	NodeBase
	Variable *lexer.Token
}

// PrefixExpression is `<op> <expr>`.
type PrefixExpression struct {
	NodeBase
	Op         *lexer.Token
	Expression Node
}

// InfixExpression is `<expr> <op> <expr>`.
type InfixExpression struct {
	NodeBase
	Left  Node
	Op    *lexer.Token
	Right Node
}

// PostfixExpression is `<expr> <op>`.
type PostfixExpression struct {
	NodeBase
	Expression Node
	Op         *lexer.Token
}

// FunctionExpression is a backtick expression.
type FunctionExpression struct {
	NodeBase
	Value *lexer.Token
}

// FunctionApplication is a space-separated sequence: callee followed by args.
type FunctionApplication struct {
	NodeBase
	Callee Node
	Args   []Node
}

// BlockExpression is `{ ... }`.
type BlockExpression struct {
	NodeBase
	Open  *lexer.Token
	Body  []Node // *FunctionApplication | *ElementDeclaration
	Close *lexer.Token
}

// ListExpression is a `[ ... ]` setting list.
type ListExpression struct {
	NodeBase
	Open     *lexer.Token
	Elements []*Attribute
	Commas   []*lexer.Token
	Close    *lexer.Token
}

// TupleExpression is `( a, b, ... )`. A parenthesised single expression is a
// TupleExpression with IsGroup set (the GroupExpression of @dbml/parse).
type TupleExpression struct {
	NodeBase
	Open     *lexer.Token
	Elements []Node
	Commas   []*lexer.Token
	Close    *lexer.Token
	IsGroup  bool
}

// CallExpression is `<callee>(<args>)`.
type CallExpression struct {
	NodeBase
	Callee       Node
	ArgumentList *TupleExpression
}

// PrimaryExpression wraps a Literal or Variable.
type PrimaryExpression struct {
	NodeBase
	Expression Node // *Literal | *Variable
}

// CommaExpression is a CSV row `a, b, , c`.
type CommaExpression struct {
	NodeBase
	Elements []Node
	Commas   []*lexer.Token
}

// Empty is a zero-width placeholder (e.g. an empty CSV field).
type Empty struct {
	NodeBase
}

// Array is `<expr>[<indexer>]`.
type Array struct {
	NodeBase
	Array   Node
	Indexer *ListExpression
}

// Wildcard is `*`.
type Wildcard struct {
	NodeBase
	Token *lexer.Token
}
