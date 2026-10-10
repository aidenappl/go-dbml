package importer

import (
	"strings"
	"testing"

	"github.com/aidenappl/go-dbml/model"
)

func TestNormalizeExpr(t *testing.T) {
	cases := []struct {
		d        dialect
		in, want string // want "" = dropped
		warn     bool
	}{
		{dMySQL, "`a` > 0", "a > 0", false},
		{dMySQL, "`order` <> 'x'", `"order" <> 'x'`, false}, // reserved: stays quoted
		{dMySQL, "`a b` = \"it\\\"s\"", `"a b" = 'it"s'`, false},
		{dMySQL, `a <> 'x\')); DROP TABLE v; -- '`, `a <> 'x'')); DROP TABLE v; -- '`, false},
		{dMySQL, `a REGEXP '^\\d'`, `a REGEXP '^\\d'`, true},
		{dMySQL, "a > 0 # c", "", false},
		{dMySQL, "_utf8mb4'x'", "_utf8mb4'x'", false},
		{dPostgres, `a <> E'x\'y'`, `a <> 'x''y'`, false},
		{dPostgres, `a <> E'\\'`, `a <> E'\\'`, true},
		{dPostgres, `a <> 'p\'`, `a <> E'p\\'`, true},
		{dPostgres, `d <> $$')); DROP TABLE v; -- $$`, `d <> ''')); DROP TABLE v; -- '`, false},
		{dPostgres, `d <> $t$x$t$`, `d <> 'x'`, false},
		{dPostgres, `"A" > 0`, `"A" > 0`, false},
		{dPostgres, `a > 0 -- c`, "", false},
		{dMSSQL, "[a] > 0 AND [b c] = N'x'", `a > 0 AND "b c" = N'x'`, false},
		{dMSSQL, `a <> N'x\'`, `a <> N'x\'`, true},
		{dMSSQL, "(a > 0", "", false},
		{dRails, "a <> 'x'", "a <> 'x'", false},
		{dRails, `a <> 'y\'`, "", false},
	}
	for _, c := range cases {
		got, w, ok := normalizeExpr(c.in, c.d)
		if !ok {
			got = ""
		}
		if got != c.want || (ok && (w != "") != c.warn) {
			t.Errorf("normalizeExpr(%q, %d) = %q (ok=%v warn=%q), want %q warn=%v", c.in, c.d, got, ok, w, c.want, c.warn)
		}
	}
}

func TestImportGeneratedColumns(t *testing.T) {
	cases := []struct {
		f      Format
		src    string
		expr   string
		stored bool
	}{
		{MySQL, "CREATE TABLE t (a int, b int GENERATED ALWAYS AS (`a` * 2) STORED);", "a * 2", true},
		{MySQL, "CREATE TABLE t (a int, b int AS (a + 1));", "a + 1", false},
		{Postgres, "CREATE TABLE t (a int, b int GENERATED ALWAYS AS (a * 2) STORED);", "a * 2", true},
		{MSSQL, "CREATE TABLE t (a int, b AS ([a] * 2) PERSISTED)\nGO", "a * 2", true},
		{MSSQL, "CREATE TABLE t (a int, b AS a + 1)\nGO", "a + 1", false},
	}
	for _, c := range cases {
		r, err := ImportWith(c.src, c.f, ImportOptions{})
		if err != nil || r.Database == nil {
			t.Fatalf("%s: %v %v", c.src, err, r.Diagnostics)
		}
		f := r.Database.Tables[0].Fields[1]
		if g := f.Generated; g == nil || g.Expression != c.expr || g.Stored != c.stored {
			t.Errorf("%s: Generated = %+v, want %q stored=%v\n%s", c.src, g, c.expr, c.stored, r.DBML)
		}
		if strings.Contains(f.Type.TypeName, "AS") {
			t.Errorf("%s: expression smuggled through the type %q", c.src, f.Type.TypeName)
		}
	}
}

func TestImportWithNeverRaisesModelRowLimit(t *testing.T) {
	var b strings.Builder
	b.WriteString("CREATE TABLE t (a int);\nINSERT INTO t VALUES ")
	for i := 0; i < 50; i++ {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString("(1)")
	}
	b.WriteString(";")
	r, err := ImportWith(b.String(), MySQL, ImportOptions{IncludeRecords: true, MaxRecordRows: 40, Model: model.Options{MaxRecordRows: 10}})
	if err != nil || r.Database == nil {
		t.Fatal(err, r.Diagnostics)
	}
	if n := len(r.Database.Records[0].Values); n != 10 || !r.RecordsTruncated {
		t.Errorf("rows = %d truncated=%v, want 10 (the stricter Model limit)", n, r.RecordsTruncated)
	}
}

func TestMSSQLIdentityIsIncrement(t *testing.T) {
	r, err := ImportWith("CREATE TABLE t (id int IDENTITY(1,1) NOT NULL PRIMARY KEY)\nGO", MSSQL, ImportOptions{})
	if err != nil || r.Database == nil {
		t.Fatal(err)
	}
	f := r.Database.Tables[0].Fields[0]
	if f.Type.TypeName != "int" || f.Increment == nil || !*f.Increment {
		t.Errorf("type %q increment %v", f.Type.TypeName, f.Increment)
	}
}
