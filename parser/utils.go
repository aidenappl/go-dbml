// Portions of this file are derived from holistics/dbml packages/dbml-parse/src/core/parser/contextStack.ts, src/core/parser/utils.ts
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

package parser

import (
	"github.com/aidenappl/go-dbml/ast"
	"github.com/aidenappl/go-dbml/lexer"
)

type contextStack struct {
	stack    []ctxKind
	nParens  int
	nBracket int
	nBraces  int
}

func (c *contextStack) push(k ctxKind) {
	c.stack = append(c.stack, k)
	switch k {
	case ctxList:
		c.nBracket++
	case ctxGroup:
		c.nParens++
	case ctxBlock:
		c.nBraces++
	}
}

func (c *contextStack) pop() {
	if len(c.stack) == 0 {
		return
	}
	k := c.stack[len(c.stack)-1]
	c.stack = c.stack[:len(c.stack)-1]
	switch k {
	case ctxList:
		c.nBracket--
	case ctxGroup:
		c.nParens--
	case ctxBlock:
		c.nBraces--
	}
}

func (c *contextStack) top() ctxKind {
	if len(c.stack) == 0 {
		return ctxNone
	}
	return c.stack[len(c.stack)-1]
}

func (c *contextStack) withinGroup() bool { return c.nParens > 0 }
func (c *contextStack) withinList() bool  { return c.nBracket > 0 }
func (c *contextStack) withinBlock() bool { return c.nBraces > 0 }

// findHandler returns the context that should handle an error at token index
// cur, by scanning ahead for the first closing delimiter that matches an open
// context.
func (c *contextStack) findHandler(tokens []*lexer.Token, cur int) ctxKind {
	if c.nBraces <= 0 && c.nBracket <= 0 && c.nParens <= 0 {
		return ctxNone
	}
	for i := cur; i < len(tokens)-1; i++ {
		switch tokens[i].Kind {
		case lexer.Comma:
			if c.withinGroup() || c.withinList() {
				for j := len(c.stack) - 1; j >= 0; j-- {
					if c.stack[j] == ctxGroup || c.stack[j] == ctxList {
						return c.stack[j]
					}
				}
			}
		case lexer.RParen:
			if c.withinGroup() {
				return ctxGroup
			}
		case lexer.RBrace:
			if c.withinBlock() {
				return ctxBlock
			}
		case lexer.RBracket:
			if c.withinList() {
				return ctxList
			}
		}
	}
	return ctxNone
}

func isIdentifierNode(n ast.Node) bool {
	pe, ok := n.(*ast.PrimaryExpression)
	if !ok || pe == nil {
		return false
	}
	v, ok := pe.Expression.(*ast.Variable)
	return ok && v != nil && v.Variable != nil && v.Variable.Kind == lexer.Identifier
}

func variableToken(n ast.Node) *lexer.Token {
	pe, ok := n.(*ast.PrimaryExpression)
	if !ok || pe == nil {
		return nil
	}
	v, ok := pe.Expression.(*ast.Variable)
	if !ok || v == nil {
		return nil
	}
	return v.Variable
}

// convertFuncAppToElem interprets `type [name [as alias]] [settings] { ... }`
// written as a function application inside a block as an element declaration.
func convertFuncAppToElem(callee ast.Node, args []ast.Node, f *ast.Factory) *ast.ElementDeclaration {
	if call, ok := callee.(*ast.CallExpression); ok && call != nil && call.ArgumentList != nil {
		args = append([]ast.Node{call.ArgumentList}, args...)
		callee = call.Callee
	}
	if ast.IsNil(callee) || !isIdentifierNode(callee) || len(args) == 0 {
		return nil
	}
	cp := append([]ast.Node{}, args...)
	typ := variableToken(callee)

	body, ok := cp[len(cp)-1].(*ast.BlockExpression)
	if !ok || body == nil {
		return nil
	}
	cp = cp[:len(cp)-1]

	var attrs *ast.ListExpression
	if len(cp) > 0 {
		if l, ok := cp[len(cp)-1].(*ast.ListExpression); ok && l != nil {
			attrs = l
			cp = cp[:len(cp)-1]
		}
	}

	switch len(cp) {
	case 3:
		asKw := variableToken(cp[1])
		if asKw == nil || !isKeyword(asKw, "as") {
			return nil
		}
		return f.NewElementDeclaration(ast.ElementDeclaration{Type: typ, Name: cp[0], As: asKw, Alias: cp[2], AttributeList: attrs, Body: body})
	case 1:
		return f.NewElementDeclaration(ast.ElementDeclaration{Type: typ, Name: cp[0], AttributeList: attrs, Body: body})
	case 0:
		return f.NewElementDeclaration(ast.ElementDeclaration{Type: typ, AttributeList: attrs, Body: body})
	}
	return nil
}
