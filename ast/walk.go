// Portions of this file are derived from holistics/dbml packages/dbml-parse/src/core/parser/utils.ts (getMemberChain, markInvalid)
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

package ast

import "github.com/aidenappl/go-dbml/lexer"

// Children returns the direct children of n (tokens and nodes) in source
// order, skipping absent ones.
func Children(n Node) []Element {
	var out []Element
	eachChild(n, func(e Element) { out = append(out, e) })
	return out
}

// eachChild calls yield for every direct child of n (tokens and nodes) in
// source order, skipping absent ones, without allocating.
func eachChild(n Node, yield func(Element)) {
	addT := func(t *lexer.Token) {
		if t != nil {
			yield(t)
		}
	}
	addN := func(c Node) {
		if !IsNil(c) {
			yield(c)
		}
	}
	switch v := n.(type) {
	case *Program:
		for _, c := range v.Body {
			addN(c)
		}
		addT(v.EOF)
	case *ElementDeclaration:
		addT(v.Type)
		addT(v.TargetKind)
		addN(v.Name)
		addT(v.As)
		addN(v.Alias)
		if v.AttributeList != nil {
			addN(v.AttributeList)
		}
		addT(v.BodyColon)
		addN(v.Body)
	case *Attribute:
		addN(v.Name)
		addT(v.Colon)
		addN(v.Value)
	case *IdentifierStream:
		for _, t := range v.Identifiers {
			addT(t)
		}
	case *Literal:
		addT(v.Literal)
	case *Variable:
		addT(v.Variable)
	case *PrefixExpression:
		addT(v.Op)
		addN(v.Expression)
	case *InfixExpression:
		addN(v.Left)
		addT(v.Op)
		addN(v.Right)
	case *PostfixExpression:
		addN(v.Expression)
		addT(v.Op)
	case *FunctionExpression:
		addT(v.Value)
	case *FunctionApplication:
		addN(v.Callee)
		for _, a := range v.Args {
			addN(a)
		}
	case *BlockExpression:
		addT(v.Open)
		for _, b := range v.Body {
			addN(b)
		}
		addT(v.Close)
	case *ListExpression:
		addT(v.Open)
		alternateAdd(len(v.Elements), len(v.Commas), func(i int) { addN(v.Elements[i]) }, func(i int) { addT(v.Commas[i]) })
		addT(v.Close)
	case *TupleExpression:
		addT(v.Open)
		alternateAdd(len(v.Elements), len(v.Commas), func(i int) { addN(v.Elements[i]) }, func(i int) { addT(v.Commas[i]) })
		addT(v.Close)
	case *CommaExpression:
		alternateAdd(len(v.Elements), len(v.Commas), func(i int) { addN(v.Elements[i]) }, func(i int) { addT(v.Commas[i]) })
	case *CallExpression:
		addN(v.Callee)
		if v.ArgumentList != nil {
			addN(v.ArgumentList)
		}
	case *PrimaryExpression:
		addN(v.Expression)
	case *Array:
		addN(v.Array)
		if v.Indexer != nil {
			addN(v.Indexer)
		}
	case *UseDeclaration:
		addT(v.UseKeyword)
		addN(v.Specifiers)
		addT(v.FromKeyword)
		addT(v.ImportPath)
	case *UseSpecifierList:
		addT(v.OpenBrace)
		for _, s := range v.Specifiers {
			addN(s)
		}
		addT(v.CloseBrace)
	case *UseSpecifier:
		addT(v.ImportKind)
		addN(v.Name)
		addT(v.AsKeyword)
		addN(v.Alias)
	case *Wildcard:
		addT(v.Token)
	case *Empty:
	}
}

func alternateAdd(na, nb int, a, b func(int)) {
	n := na
	if nb < n {
		n = nb
	}
	for i := 0; i < n; i++ {
		a(i)
		b(i)
	}
	for i := n; i < na; i++ {
		a(i)
	}
	for i := n; i < nb; i++ {
		b(i)
	}
}

// AssignParents sets Parent on every node under root.
func AssignParents(root Node) {
	eachChild(root, func(c Element) {
		if cn, ok := c.(Node); ok {
			cn.Base().Parent = root
			AssignParents(cn)
		}
	})
}

// MarkInvalid flags every token under n (or n itself if it is a token) as
// invalid. EOF tokens are never marked.
func MarkInvalid(x Element) {
	switch v := x.(type) {
	case nil:
		return
	case *lexer.Token:
		if v != nil && v.Kind != lexer.EOF {
			v.Invalid = true
		}
	case Node:
		if IsNil(v) {
			return
		}
		eachChild(v, MarkInvalid)
	}
}

// Walk calls fn for n and every descendant node in pre-order. If fn returns
// false the node's children are skipped.
func Walk(n Node, fn func(Node) bool) {
	if IsNil(n) || !fn(n) {
		return
	}
	eachChild(n, func(c Element) {
		if cn, ok := c.(Node); ok {
			Walk(cn, fn)
		}
	})
}

// ParentElement returns the nearest enclosing ElementDeclaration or Program.
func ParentElement(n Node) Node {
	cur := n.Base().Parent
	for !IsNil(cur) {
		switch cur.(type) {
		case *ElementDeclaration, *Program:
			return cur
		}
		cur = cur.Base().Parent
	}
	return nil
}

// ParentElementDecl returns the nearest enclosing ElementDeclaration, or nil.
func ParentElementDecl(n Node) *ElementDeclaration {
	if e, ok := ParentElement(n).(*ElementDeclaration); ok {
		return e
	}
	return nil
}

// ContainsEq reports whether inner's span lies within outer's span.
func ContainsEq(outer, inner Node) bool {
	if IsNil(outer) || IsNil(inner) {
		return false
	}
	o, i := outer.Base(), inner.Base()
	if !o.HasPos || !i.HasPos {
		return false
	}
	return o.Start.Offset <= i.Start.Offset && o.End.Offset >= i.End.Offset
}
