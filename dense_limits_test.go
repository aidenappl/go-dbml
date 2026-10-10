package dbml

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/aidenappl/go-dbml/diag"
)

// TestDenseInputsBounded: 2 MiB inputs built to maximise tokens/elements
// per byte (re-verifier probes, Oct 2026) are refused by MaxTokens /
// MaxElements / MaxRecordBytes with bounded allocation.
func TestDenseInputsBounded(t *testing.T) {
	const L = 2<<20 - 64
	var b strings.Builder
	gen := func(f func(i int)) string {
		b.Reset()
		for i := 0; b.Len() < L-32; i++ {
			f(i)
		}
		return b.String()
	}
	cases := []struct {
		name string
		src  string
		code diag.Code
	}{
		{"columns", "Table t {\n" + gen(func(i int) { fmt.Fprintf(&b, "  c%d int\n", i) }) + "}\n", diag.TOO_MANY_ELEMENTS},
		{"enum", "Enum e {\n" + gen(func(i int) { fmt.Fprintf(&b, "  v%d\n", i) }) + "}\n", diag.TOO_MANY_ELEMENTS},
		{"index", "Table t {\n  a int\n  indexes {\n    (" + strings.Repeat("a,", L/2-40) + "a)\n  }\n}\n", diag.TOO_MANY_TOKENS},
		{"row", "Table t {\n  a int\n}\nRecords t(a) {\n  " + strings.Repeat("1,", L/2-40) + "1\n}\n", diag.TOO_MANY_TOKENS},
		{"rows", "Table t {\n  a int\n}\nRecords t(a) {\n" + strings.Repeat("  1\n", L/4-20) + "}\n", diag.TOO_MANY_ELEMENTS},
		{"records_bytes", "Table t {\n  a varchar\n}\nRecords t(a) {\n" + strings.Repeat("  'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'\n", 9000) + "}\n", diag.RECORDS_TOO_LARGE},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			opts := Options{}
			if c.name == "records_bytes" {
				opts.MaxRecordBytes = 512 << 10
			}
			runtime.GC()
			var a, z runtime.MemStats
			runtime.ReadMemStats(&a)
			r := ParseFileWith(DefaultFilepath, c.src, opts)
			runtime.ReadMemStats(&z)
			alloc := float64(z.TotalAlloc-a.TotalAlloc) / 1e6
			found := false
			for _, d := range r.Diagnostics {
				found = found || d.Code == c.code
			}
			if !found || r.Database != nil {
				t.Errorf("want %v; got %v", c.code, r.Diagnostics.Errors()[:1])
			}
			t.Logf("%d bytes: %.0f MB allocated", len(c.src), alloc)
			if alloc > 300 && !raceEnabled { // -race inflates allocation
				t.Errorf("allocated %.0f MB (> 300 MB)", alloc)
			}
		})
	}
}
