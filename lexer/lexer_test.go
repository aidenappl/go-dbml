package lexer

import (
	"testing"
)

func kinds(toks []*Token) []Kind {
	out := make([]Kind, len(toks))
	for i, t := range toks {
		out[i] = t.Kind
	}
	return out
}

func TestLexTokens(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		kinds  []Kind
		values []string
		errs   int
	}{
		{"identifier", "users", []Kind{Identifier, EOF}, []string{"users", ""}, 0},
		{"digit-leading identifier", "1st_col", []Kind{Identifier, EOF}, []string{"1st_col", ""}, 0},
		{"number", "12.5", []Kind{NumericLiteral, EOF}, []string{"12.5", ""}, 0},
		{"scientific", "1.5e10", []Kind{NumericLiteral, EOF}, []string{"1.5e10", ""}, 0},
		{"invalid number", "1.2.3", []Kind{EOF}, []string{""}, 1},
		{"single string", `'it\'s'`, []Kind{StringLiteral, EOF}, []string{"it's", ""}, 0},
		{"multiline string", "'''a\nb'''", []Kind{StringLiteral, EOF}, []string{"a\nb", ""}, 0},
		{"quoted name", `"my table"`, []Kind{QuotedString, EOF}, []string{"my table", ""}, 0},
		{"function expression", "`now()`", []Kind{FunctionExpression, EOF}, []string{"now()", ""}, 0},
		{"color", "#3498DB", []Kind{ColorLiteral, EOF}, []string{"#3498DB", ""}, 0},
		{"ref ops", "< > - <> >? ?< <>?", []Kind{Op, Op, Op, Op, Op, Op, Op, EOF}, []string{"<", ">", "-", "<>", ">?", "?<", "<>?", ""}, 0},
		{"dep ops", "-> <-", []Kind{Op, Op, EOF}, []string{"->", "<-", ""}, 0},
		{"wildcard", "*", []Kind{Wildcard, EOF}, []string{"*", ""}, 0},
		{"comments are trivia", "a // c\n/* d */ b", []Kind{Identifier, Identifier, EOF}, []string{"a", "b", ""}, 0},
		{"unclosed string", "'abc", []Kind{EOF}, []string{""}, 1},
		{"unknown symbol", "a @ b", []Kind{Identifier, Identifier, EOF}, []string{"a", "b", ""}, 1},
		{"punctuation", "{}[](),:;", []Kind{LBrace, RBrace, LBracket, RBracket, LParen, RParen, Comma, Colon, Semicolon, EOF}, nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toks, errs := Lex(tt.src, "/main.dbml")
			if len(errs) != tt.errs {
				t.Fatalf("errors: want %d, got %d (%v)", tt.errs, len(errs), errs)
			}
			got := kinds(toks)
			if len(got) != len(tt.kinds) {
				t.Fatalf("kinds: want %v, got %v", tt.kinds, got)
			}
			for i := range got {
				if got[i] != tt.kinds[i] {
					t.Fatalf("kind %d: want %v, got %v", i, tt.kinds[i], got[i])
				}
				if tt.values != nil && toks[i].Value != tt.values[i] {
					t.Fatalf("value %d: want %q, got %q", i, tt.values[i], toks[i].Value)
				}
			}
		})
	}
}

func TestLexPositionsUTF16(t *testing.T) {
	// "😀" is two UTF-16 units; positions must follow JavaScript semantics.
	toks, _ := Lex("'😀' x", "/main.dbml")
	x := toks[1]
	if x.Value != "x" || x.Start.Offset != 5 || x.Start.Column != 5 || x.Start.Line != 0 {
		t.Fatalf("unexpected position %+v", x.Start)
	}
	toks, _ = Lex("a\n  b", "/main.dbml")
	if b := toks[1]; b.Start.Line != 1 || b.Start.Column != 2 || b.Start.Offset != 4 {
		t.Fatalf("unexpected position %+v", b.Start)
	}
}

func TestTrivia(t *testing.T) {
	toks, _ := Lex("  a // note\nb", "/main.dbml")
	if len(toks[0].LeadingTrivia()) != 2 {
		t.Fatalf("want 2 leading trivia, got %d", len(toks[0].LeadingTrivia()))
	}
	if !toks[0].HasTrailingNewline() {
		t.Fatal("want trailing newline on a")
	}
}
