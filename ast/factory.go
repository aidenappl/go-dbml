// Portions of this file are derived from holistics/dbml packages/dbml-parse/src/core/types/nodes.ts, src/core/parser/factory.ts
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

package ast

import "github.com/aidenappl/go-dbml/lexer"

// Factory creates nodes with sequential ids and computed spans.
type Factory struct {
	Filepath string
	nextID   int
}

func (f *Factory) id() int {
	id := f.nextID
	f.nextID++
	return id
}

func nodesToSpans[T Node](nodes []T) []span {
	out := make([]span, len(nodes))
	for i, n := range nodes {
		out[i] = n
	}
	return out
}

func tokensToSpans(toks []*lexer.Token) []span {
	out := make([]span, len(toks))
	for i, t := range toks {
		out[i] = t
	}
	return out
}

func alternate(a, b []span) []span {
	out := make([]span, 0, len(a)+len(b))
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		out = append(out, a[i], b[i])
	}
	out = append(out, a[n:]...)
	out = append(out, b[n:]...)
	return out
}

func cat(parts ...[]span) []span {
	var out []span
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

// NewProgram creates a Program.
func (f *Factory) NewProgram(body []Node, eof *lexer.Token, source string) *Program {
	n := &Program{Body: body, EOF: eof, Source: source}
	n.Init(f.id(), KindProgram, f.Filepath, cat(nodesToSpans(body), []span{eof})...)
	return n
}

// NewUseDeclaration creates a UseDeclaration.
func (f *Factory) NewUseDeclaration(useKw *lexer.Token, specifiers Node, fromKw, path *lexer.Token) *UseDeclaration {
	n := &UseDeclaration{UseKeyword: useKw, Specifiers: specifiers, FromKeyword: fromKw, ImportPath: path}
	n.Init(f.id(), KindUseDeclaration, f.Filepath, useKw, specifiers, fromKw, path)
	return n
}

// NewUseSpecifier creates a UseSpecifier.
func (f *Factory) NewUseSpecifier(kind *lexer.Token, name Node, asKw *lexer.Token, alias Node) *UseSpecifier {
	n := &UseSpecifier{ImportKind: kind, Name: name, AsKeyword: asKw, Alias: alias}
	n.Init(f.id(), KindUseSpecifier, f.Filepath, kind, name, asKw, alias)
	return n
}

// NewUseSpecifierList creates a UseSpecifierList.
func (f *Factory) NewUseSpecifierList(open *lexer.Token, specs []*UseSpecifier, closeTok *lexer.Token) *UseSpecifierList {
	n := &UseSpecifierList{OpenBrace: open, Specifiers: specs, CloseBrace: closeTok}
	n.Init(f.id(), KindUseSpecifierList, f.Filepath, cat([]span{open}, nodesToSpans(specs), []span{closeTok})...)
	return n
}

// NewElementDeclaration creates an ElementDeclaration.
func (f *Factory) NewElementDeclaration(e ElementDeclaration) *ElementDeclaration {
	n := &ElementDeclaration{Type: e.Type, TargetKind: e.TargetKind, Name: e.Name, As: e.As, Alias: e.Alias,
		AttributeList: e.AttributeList, BodyColon: e.BodyColon, Body: e.Body}
	var attr span
	if n.AttributeList != nil {
		attr = n.AttributeList
	}
	n.Init(f.id(), KindElementDeclaration, f.Filepath, n.Type, n.TargetKind, n.Name, n.As, n.Alias, attr, n.BodyColon, n.Body)
	return n
}

// NewAttribute creates an Attribute.
func (f *Factory) NewAttribute(name Node, colon *lexer.Token, value Node) *Attribute {
	n := &Attribute{Name: name, Colon: colon, Value: value}
	n.Init(f.id(), KindAttribute, f.Filepath, name, colon, value)
	return n
}

// NewIdentifierStream creates an IdentifierStream.
func (f *Factory) NewIdentifierStream(ids []*lexer.Token) *IdentifierStream {
	n := &IdentifierStream{Identifiers: ids}
	n.Init(f.id(), KindIdentifierStream, f.Filepath, tokensToSpans(ids)...)
	return n
}

// NewLiteral creates a Literal.
func (f *Factory) NewLiteral(t *lexer.Token) *Literal {
	n := &Literal{Literal: t}
	n.Init(f.id(), KindLiteral, f.Filepath, t)
	return n
}

// NewVariable creates a Variable.
func (f *Factory) NewVariable(t *lexer.Token) *Variable {
	n := &Variable{Variable: t}
	n.Init(f.id(), KindVariable, f.Filepath, t)
	return n
}

// NewPrefix creates a PrefixExpression.
func (f *Factory) NewPrefix(op *lexer.Token, expr Node) *PrefixExpression {
	n := &PrefixExpression{Op: op, Expression: expr}
	n.Init(f.id(), KindPrefixExpression, f.Filepath, op, expr)
	return n
}

// NewInfix creates an InfixExpression.
func (f *Factory) NewInfix(left Node, op *lexer.Token, right Node) *InfixExpression {
	n := &InfixExpression{Left: left, Op: op, Right: right}
	n.Init(f.id(), KindInfixExpression, f.Filepath, left, op, right)
	return n
}

// NewPostfix creates a PostfixExpression.
func (f *Factory) NewPostfix(expr Node, op *lexer.Token) *PostfixExpression {
	n := &PostfixExpression{Expression: expr, Op: op}
	n.Init(f.id(), KindPostfixExpression, f.Filepath, expr, op)
	return n
}

// NewFunctionExpression creates a FunctionExpression.
func (f *Factory) NewFunctionExpression(v *lexer.Token) *FunctionExpression {
	n := &FunctionExpression{Value: v}
	n.Init(f.id(), KindFunctionExpression, f.Filepath, v)
	return n
}

// NewFunctionApplication creates a FunctionApplication.
func (f *Factory) NewFunctionApplication(callee Node, args []Node) *FunctionApplication {
	n := &FunctionApplication{Callee: callee, Args: args}
	n.Init(f.id(), KindFunctionApplication, f.Filepath, cat([]span{callee}, nodesToSpans(args))...)
	return n
}

// NewBlock creates a BlockExpression.
func (f *Factory) NewBlock(open *lexer.Token, body []Node, closeTok *lexer.Token) *BlockExpression {
	n := &BlockExpression{Open: open, Body: body, Close: closeTok}
	n.Init(f.id(), KindBlockExpression, f.Filepath, cat([]span{open}, nodesToSpans(body), []span{closeTok})...)
	return n
}

// NewList creates a ListExpression.
func (f *Factory) NewList(open *lexer.Token, elems []*Attribute, commas []*lexer.Token, closeTok *lexer.Token) *ListExpression {
	n := &ListExpression{Open: open, Elements: elems, Commas: commas, Close: closeTok}
	n.Init(f.id(), KindListExpression, f.Filepath, cat([]span{open}, alternate(nodesToSpans(elems), tokensToSpans(commas)), []span{closeTok})...)
	return n
}

// NewTuple creates a TupleExpression.
func (f *Factory) NewTuple(open *lexer.Token, elems []Node, commas []*lexer.Token, closeTok *lexer.Token) *TupleExpression {
	n := &TupleExpression{Open: open, Elements: elems, Commas: commas, Close: closeTok}
	n.Init(f.id(), KindTupleExpression, f.Filepath, cat([]span{open}, alternate(nodesToSpans(elems), tokensToSpans(commas)), []span{closeTok})...)
	return n
}

// NewGroup creates a single-element TupleExpression marked as a group.
func (f *Factory) NewGroup(open *lexer.Token, expr Node, closeTok *lexer.Token) *TupleExpression {
	var elems []Node
	if !IsNil(expr) {
		elems = []Node{expr}
	}
	n := &TupleExpression{Open: open, Elements: elems, Close: closeTok, IsGroup: true}
	n.Init(f.id(), KindGroupExpression, f.Filepath, cat([]span{open}, nodesToSpans(elems), []span{closeTok})...)
	return n
}

// NewCall creates a CallExpression.
func (f *Factory) NewCall(callee Node, args *TupleExpression) *CallExpression {
	n := &CallExpression{Callee: callee, ArgumentList: args}
	var a span
	if args != nil {
		a = args
	}
	n.Init(f.id(), KindCallExpression, f.Filepath, callee, a)
	return n
}

// NewPrimary creates a PrimaryExpression.
func (f *Factory) NewPrimary(expr Node) *PrimaryExpression {
	n := &PrimaryExpression{Expression: expr}
	n.Init(f.id(), KindPrimaryExpression, f.Filepath, expr)
	return n
}

// NewComma creates a CommaExpression.
func (f *Factory) NewComma(elems []Node, commas []*lexer.Token) *CommaExpression {
	n := &CommaExpression{Elements: elems, Commas: commas}
	n.Init(f.id(), KindCommaExpression, f.Filepath, alternate(nodesToSpans(elems), tokensToSpans(commas))...)
	return n
}

// NewEmpty creates a zero-width Empty node positioned at the end of prev
// (a *lexer.Token or Node).
func (f *Factory) NewEmpty(prev span) *Empty {
	n := &Empty{}
	var end lexer.Position
	ok := false
	switch v := prev.(type) {
	case *lexer.Token:
		if v != nil {
			end, ok = v.End, true
		}
	case Node:
		if !IsNil(v) && v.Base().HasPos {
			end, ok = v.Base().End, true
		}
	}
	if ok {
		fake := &lexer.Token{Kind: lexer.Space, Start: end, End: end, Value: " "}
		n.Init(f.id(), KindEmpty, f.Filepath, fake)
	} else {
		n.Init(f.id(), KindEmpty, f.Filepath)
	}
	return n
}

// NewArray creates an Array node.
func (f *Factory) NewArray(expr Node, indexer *ListExpression) *Array {
	n := &Array{Array: expr, Indexer: indexer}
	var ix span
	if indexer != nil {
		ix = indexer
	}
	n.Init(f.id(), KindArray, f.Filepath, expr, ix)
	return n
}

// NewWildcard creates a Wildcard node.
func (f *Factory) NewWildcard(t *lexer.Token) *Wildcard {
	n := &Wildcard{Token: t}
	n.Init(f.id(), KindWildcard, f.Filepath, t)
	return n
}
