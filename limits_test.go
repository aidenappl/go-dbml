package dbml

import (
	"strings"
	"testing"
	"time"

	"github.com/aidenappl/go-dbml/diag"
)

// deepInputs are hostile shapes that used to overflow the stack (each is
// ~1-4 MB at n = 10^6).
func deepInputs(n int) map[string]string {
	return map[string]string{
		"paren":          "Table a { id int [default: `x`, note: " + strings.Repeat("(", n) + "1" + strings.Repeat(")", n) + "] }",
		"paren-expr":     "Ref: " + strings.Repeat("(", n) + "a.b" + strings.Repeat(")", n) + " > c.d",
		"unclosed-paren": "Ref: " + strings.Repeat("(", n),
		"prefix":         "Table a { id int [default: " + strings.Repeat("-", n) + "1] }",
		"bracket":        "Table a { id int " + strings.Repeat("[", n) + " }",
		"nested-setting": "Table a { id int " + strings.Repeat("[a: ", n) + strings.Repeat("]", n) + " }",
		"block":          "Table a " + strings.Repeat("{ ", n) + strings.Repeat("} ", n),
		"unclosed-block": "Table a " + strings.Repeat("{ x ", n),
		"dots":           "Ref: a" + strings.Repeat(".a", n) + " > b.c",
		"binary":         "Table a { id int [default: 1" + strings.Repeat(" + 1", n) + "] }",
		"call":           "Table a { id int [default: f" + strings.Repeat("(", n) + strings.Repeat(")", n) + "] }",
	}
}

func TestDeepNestingIsBounded(t *testing.T) {
	for name, src := range deepInputs(1_000_000) {
		t.Run(name, func(t *testing.T) {
			start := time.Now()
			db, diags := Parse(src)
			if db != nil {
				t.Fatal("want errors for hostile input")
			}
			found := false
			for _, d := range diags {
				if d.Code == diag.NESTING_TOO_DEEP || d.Code == diag.TOO_MANY_DIAGNOSTICS {
					found = true
				}
			}
			if !found && name != "unclosed-block" && name != "bracket" {
				t.Fatalf("want NESTING_TOO_DEEP, got %d diagnostics", len(diags))
			}
			el := time.Since(start)
			strictBudget(t, name, el, 10*time.Second)
			// linear: 4x the nesting costs < 6x the time
			assertScales(t, name, 125_000, 500_000, 2.5, func(n int) { Parse(deepInputs(n)[name]) })
		})
	}
}

func TestShallowNestingStillParses(t *testing.T) {
	src := "Table a { id int [note: '" + strings.Repeat("x", 10) + "', default: " + strings.Repeat("(", 100) + "1" + strings.Repeat(")", 100) + "] }"
	_, diags := Parse(src)
	for _, d := range diags {
		if d.Code == diag.NESTING_TOO_DEEP {
			t.Fatal("100 levels must stay under the default limit")
		}
	}
}

func wideTable(cols int) string {
	var sb strings.Builder
	sb.WriteString("Table wide {\n")
	for i := 0; i < cols; i++ {
		sb.WriteString("  c")
		sb.WriteString(itoa(i))
		sb.WriteString(" int [note: 'n']\n")
	}
	sb.WriteString("}\n")
	return sb.String()
}

func manyTables(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString("Table t" + itoa(i) + " {\n  id int [pk]\n  parent_id int [ref: > t" + itoa((i+1)%n) + ".id]\n  name varchar\n}\n")
		sb.WriteString("Ref: t" + itoa(i) + ".name - t" + itoa((i+2)%n) + ".name\n")
	}
	sb.WriteString("Metadata Table t0 {\n owner: 'x'\n}\n")
	return sb.String()
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b []byte
	for i > 0 {
		b = append([]byte{byte('0' + i%10)}, b...)
		i /= 10
	}
	return string(b)
}

func TestLargeSchemasCompileFast(t *testing.T) {
	tests := []struct {
		name string
		src  string
	}{
		{"20k-column table", wideTable(20_000)},
		{"4k tables with refs", manyTables(4_000)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			// above the default MaxElements (10,000): timing only
			r := ParseFileWith(DefaultFilepath, tt.src, Options{MaxElements: 50_000})
			db, diags := r.Database, r.Diagnostics
			el := time.Since(start)
			if db == nil {
				t.Fatalf("unexpected errors: %v", diags.Errors()[:1])
			}
			t.Logf("%s: %v", tt.name, el)
			strictBudget(t, tt.name, el, time.Second)
			gen := map[string]func(int) string{"20k-column table": wideTable, "4k tables with refs": manyTables}[tt.name]
			n := map[string]int{"20k-column table": 20_000, "4k tables with refs": 4_000}[tt.name]
			assertScales(t, tt.name, n/4, n, 2.5, func(k int) {
				ParseFileWith(DefaultFilepath, gen(k), Options{MaxElements: 50_000})
			})
		})
	}
}

func BenchmarkWideTable(b *testing.B) {
	src := wideTable(20_000)
	for i := 0; i < b.N; i++ {
		Parse(src)
	}
}

func BenchmarkManyTables(b *testing.B) {
	src := manyTables(4_000)
	for i := 0; i < b.N; i++ {
		Parse(src)
	}
}

func TestOptionsLimits(t *testing.T) {
	t.Run("max bytes", func(t *testing.T) {
		r := ParseFileWith(DefaultFilepath, "Table a { id int }", Options{MaxBytes: 5})
		if r.Database != nil || r.Diagnostics[0].Code != diag.INPUT_TOO_LARGE {
			t.Fatalf("want INPUT_TOO_LARGE, got %v", r.Diagnostics)
		}
	})
	t.Run("max files", func(t *testing.T) {
		files := map[string]string{"/main.dbml": "use * from './a'\nuse * from './b'\nTable m { id int }", "/a.dbml": "Table a { id int }", "/b.dbml": "Table b { id int }"}
		r := ParseProjectWith("/main.dbml", files, Options{MaxFiles: 2})
		found := false
		for _, d := range r.Diagnostics {
			found = found || d.Code == diag.TOO_MANY_FILES
		}
		if !found {
			t.Fatalf("want TOO_MANY_FILES, got %v", r.Diagnostics)
		}
	})
	t.Run("max diagnostics", func(t *testing.T) {
		src := strings.Repeat("Table t {\n}\n", 50)
		r := ParseFileWith(DefaultFilepath, src, Options{MaxDiagnostics: 10})
		if len(r.Diagnostics) != 11 || r.Diagnostics[10].Code != diag.TOO_MANY_DIAGNOSTICS {
			t.Fatalf("want 10 + summary, got %d", len(r.Diagnostics))
		}
	})
	t.Run("huge garbage stays bounded", func(t *testing.T) {
		r := ParseFile(DefaultFilepath, strings.Repeat("@", 2_000_000))
		if len(r.Diagnostics) > DefaultOptions().MaxDiagnostics+1 {
			t.Fatalf("got %d diagnostics", len(r.Diagnostics))
		}
	})
}

func TestPartialModel(t *testing.T) {
	r := ParseFile(DefaultFilepath, "Table users {\n id int [pk]\n}\nTable posts {\n uid int [ref: > nope.id]\n}")
	if r.Database != nil {
		t.Fatal("want nil Database on error")
	}
	if r.Partial == nil || len(r.Partial.Tables) != 2 {
		t.Fatalf("want best-effort Partial with 2 tables, got %+v", r.Partial)
	}
	ok := ParseFile(DefaultFilepath, "Table a { id int }")
	if ok.Partial == nil || ok.Partial != ok.Database {
		t.Fatal("Partial must equal Database when there are no errors")
	}
}
