package parser

import (
	"testing"

	"github.com/aidenappl/go-dbml/ast"
)

func TestParseElements(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		elements int
		errs     int
		check    func(t *testing.T, p *ast.Program)
	}{
		{"table with alias and settings", "Table users as U [headercolor: #fff] {\n id int [pk]\n}", 1, 0, func(t *testing.T, p *ast.Program) {
			e := p.Declarations()[0]
			if !e.IsKind("table") || e.As == nil || e.AttributeList == nil {
				t.Fatalf("bad table: %+v", e)
			}
			if _, ok := e.Body.(*ast.BlockExpression); !ok {
				t.Fatal("want block body")
			}
		}},
		{"short ref", "Ref: a.id > b.id", 1, 0, func(t *testing.T, p *ast.Program) {
			e := p.Declarations()[0]
			fa, ok := e.Body.(*ast.FunctionApplication)
			if !ok {
				t.Fatal("want simple body")
			}
			in, ok := fa.Callee.(*ast.InfixExpression)
			if !ok || in.Op.Value != ">" {
				t.Fatalf("want infix >, got %T", fa.Callee)
			}
		}},
		{"composite ref", "Ref: a.(x, y) <> b.(x, y)", 1, 0, nil},
		{"nested element", "Table t {\n id int\n indexes {\n  (id) [pk]\n }\n}", 1, 0, func(t *testing.T, p *ast.Program) {
			b := p.Declarations()[0].Body.(*ast.BlockExpression)
			if _, ok := b.Body[1].(*ast.ElementDeclaration); !ok {
				t.Fatalf("want nested element, got %T", b.Body[1])
			}
		}},
		{"note field", "Table t {\n id int\n Note: 'hi'\n}", 1, 0, nil},
		{"use declaration", "use * from './x'\nuse {\n table a as b\n} from './y'", 2, 0, func(t *testing.T, p *ast.Program) {
			if len(p.Uses()) != 2 || !isKeyword(p.Uses()[0].UseKeyword, "use") {
				t.Fatal("want 2 use declarations")
			}
		}},
		{"metadata block", "Metadata Table users {\n owner: 'x'\n}", 1, 0, func(t *testing.T, p *ast.Program) {
			if p.Declarations()[0].TargetKind == nil {
				t.Fatal("want target kind")
			}
		}},
		{"records with rows", "records t(a, b) {\n 1, 'x'\n 2, \n}", 1, 0, nil},
		{"missing brace recovers", "Table t {\n id int\n", 1, 1, nil},
		{"garbage recovers", "Table t { id int [pk }\nTable u { id int }", 2, 1, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Parse(tt.src, "/main.dbml")
			if len(res.Errors) != tt.errs {
				t.Fatalf("errors: want %d, got %d: %v", tt.errs, len(res.Errors), res.Errors)
			}
			if len(res.Program.Body) != tt.elements {
				t.Fatalf("elements: want %d, got %d", tt.elements, len(res.Program.Body))
			}
			if tt.check != nil {
				tt.check(t, res.Program)
			}
		})
	}
}

func TestParentsAssigned(t *testing.T) {
	res := Parse("Table t {\n id int [ref: > u.id]\n}", "/main.dbml")
	ast.Walk(res.Program, func(n ast.Node) bool {
		if n == ast.Node(res.Program) {
			return true
		}
		if ast.IsNil(n.Base().Parent) {
			t.Fatalf("node %s has no parent", n.Base().Kind)
		}
		return true
	})
}

func TestSpans(t *testing.T) {
	res := Parse("Table users {\n  id int\n}", "/main.dbml")
	e := res.Program.Declarations()[0]
	b := e.Base()
	if b.Start.Offset != 0 || b.End.Offset != 24 || b.End.Line != 2 {
		t.Fatalf("unexpected span %+v - %+v", b.Start, b.End)
	}
}
