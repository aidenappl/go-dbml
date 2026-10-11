package importer_test

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/aidenappl/go-dbml/importer"
)

func mysqldump(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "CREATE TABLE `t%d` (\n  `id` int NOT NULL AUTO_INCREMENT,\n  `p_id` int DEFAULT NULL,\n  PRIMARY KEY (`id`),\n  CONSTRAINT `fk%d` FOREIGN KEY (`p_id`) REFERENCES `t%d` (`id`)\n) ENGINE=InnoDB;\n", i, i, (i+1)%n)
	}
	return b.String()
}

// A 20k-table dump (20k FKs) imports in linear time.
func TestImport20kTables(t *testing.T) {
	if testing.Short() {
		t.Skip("slow")
	}
	src := mysqldump(20000)
	start := time.Now()
	db, err := importer.Import(src, importer.MySQL)
	if err != nil {
		t.Fatal(err)
	}
	el := time.Since(start)
	t.Logf("20k tables (%d bytes): %v", len(src), el)
	if len(db.Tables) != 20000 || len(db.Refs) != 20000 {
		t.Fatalf("%d tables, %d refs", len(db.Tables), len(db.Refs))
	}
	strictBudget(t, "20k tables", el, 1500*time.Millisecond)
	assertScales(t, "mysqldump import", 5000, 20000, 2.5, func(n int) { _, _ = importer.Import(mysqldump(n), importer.MySQL) })
}

// A 16 MB INSERT dump is skipped without tokenising it into memory when
// records are off (the default).
func TestImport16MBInsertDump(t *testing.T) {
	var b strings.Builder
	b.WriteString("CREATE TABLE t (a int);\nINSERT INTO t VALUES ")
	for b.Len() < importer.MaxInputBytes-100 {
		b.WriteString("(1),")
	}
	b.WriteString("(1);")
	src := b.String()
	runtime.GC()
	var m0, m1 runtime.MemStats
	runtime.ReadMemStats(&m0)
	start := time.Now()
	db, err := importer.Import(src, importer.MySQL)
	el := time.Since(start)
	runtime.ReadMemStats(&m1)
	if err != nil || len(db.Records) != 0 {
		t.Fatalf("%v, %d records", err, len(db.Records))
	}
	alloc := (m1.TotalAlloc - m0.TotalAlloc) >> 20
	t.Logf("16 MB INSERT dump, records skipped: %v, %d MB allocated", el, alloc)
	if alloc > 64 {
		t.Fatalf("allocated %d MB", alloc)
	}
	strictBudget(t, "16 MB dump, records off", el, 2*time.Second)
	dump := func(size int) string {
		var b strings.Builder
		b.WriteString("CREATE TABLE t (a int);\nINSERT INTO t VALUES ")
		for b.Len() < size {
			b.WriteString("(1),")
		}
		b.WriteString("(1);")
		return b.String()
	}
	small, big := dump(4<<20), dump(importer.MaxInputBytes-100)
	assertScales(t, "INSERT dump, records off", 1, 4, 2.5, func(n int) {
		s := small
		if n == 4 {
			s = big
		}
		_, _ = importer.Import(s, importer.MySQL)
	})
	// with records the row cap bounds the work
	start = time.Now()
	r, err := importer.ImportWith(src, importer.MySQL, importer.ImportOptions{IncludeRecords: true})
	if err != nil || r.Database == nil || !r.RecordsTruncated || len(r.Database.Records[0].Values) != 10000 {
		t.Fatalf("records on: %v %v", err, r.Diagnostics)
	}
	t.Logf("16 MB INSERT dump, records on (10k-row cap): %v", time.Since(start))
	strictBudget(t, "16 MB dump, records on", time.Since(start), 3*time.Second)
}

// Deeply nested defaults/values are handled without recursion in linear
// time in every dialect.
func TestDeepNesting200k(t *testing.T) {
	inputs := func(n int) map[importer.Format][]string {
		paren := strings.Repeat("(", n) + "1" + strings.Repeat(")", n)
		unary := strings.Repeat("- ", n) + "1"
		return map[importer.Format][]string{
			importer.MySQL:    {"CREATE TABLE t (a int DEFAULT " + paren + ");", "CREATE TABLE t (a int);\nINSERT INTO t VALUES (" + paren + ");"},
			importer.Postgres: {"CREATE TABLE t (a int DEFAULT " + paren + ");", "CREATE TABLE t (a int);\nINSERT INTO t VALUES (" + paren + ");"},
			importer.MSSQL: {"CREATE TABLE t (a int DEFAULT " + paren + ");", "CREATE TABLE t (a int DEFAULT " + unary + ");",
				"CREATE TABLE t (a int);\nINSERT INTO t VALUES (" + paren + ");"},
		}
	}
	for f, srcs := range inputs(200000) {
		for i, src := range srcs {
			start := time.Now()
			_, err := importer.ImportWith(src, f, importer.ImportOptions{IncludeRecords: true})
			el := time.Since(start)
			if err != nil {
				t.Errorf("%s/%d: %v", f, i, err)
			}
			strictBudget(t, fmt.Sprintf("%s/%d", f, i), el, time.Second)
			assertScales(t, fmt.Sprintf("%s/%d", f, i), 50000, 200000, 2.5, func(n int) {
				_, _ = importer.ImportWith(inputs(n)[f][i], f, importer.ImportOptions{IncludeRecords: true})
			})
		}
	}
}

// An SSMS-style T-SQL dump (one INSERT per line, no ';', a single GO at
// the end) streams statement by statement instead of being held as one
// batch.
func TestImport16MBSSMSDump(t *testing.T) {
	for _, sep := range []string{"\n", ";\n"} {
		var b strings.Builder
		b.WriteString("CREATE TABLE [dbo].[t] ([id] int NOT NULL, [s] nvarchar(100))\nGO\n")
		for i := 0; b.Len() < importer.MaxInputBytes-200; i++ {
			fmt.Fprintf(&b, "INSERT [dbo].[t] ([id], [s]) VALUES (%d, N'row %d it''s')%s", i, i, sep)
		}
		b.WriteString("GO\n")
		src := b.String()
		for _, recs := range []bool{false, true} {
			runtime.GC()
			var m0, m1 runtime.MemStats
			runtime.ReadMemStats(&m0)
			start := time.Now()
			r, err := importer.ImportWith(src, importer.MSSQL, importer.ImportOptions{IncludeRecords: recs})
			el := time.Since(start)
			runtime.ReadMemStats(&m1)
			if err != nil || r.Database == nil {
				t.Fatalf("%v %v", err, r.Diagnostics)
			}
			rows := 0
			for _, rec := range r.Database.Records {
				rows += len(rec.Values)
			}
			alloc := (m1.TotalAlloc - m0.TotalAlloc) >> 20
			t.Logf("16 MB SSMS dump sep=%q records=%v: %v, %d MB allocated, %d rows", sep, recs, el, alloc, rows)
			if (recs && rows != 10000) || (!recs && rows != 0) {
				t.Errorf("rows = %d", rows)
			}
			// the race detector inflates allocation ~10x: gate only without it
			if limit := map[bool]uint64{false: 50, true: 120}[recs]; alloc > limit && !raceEnabled {
				t.Errorf("allocated %d MB (> %d)", alloc, limit)
			}
		}
	}
}

func TestTSQLRoutineBatchesNotSplit(t *testing.T) {
	src := "CREATE TABLE t (a int)\nGO\nCREATE PROCEDURE p AS\nBEGIN\n  CREATE TABLE inner_t (b int);\n  INSERT INTO t VALUES (1);\nEND\nGO\nCREATE TABLE u (c int); INSERT INTO u VALUES (2)\nINSERT INTO u VALUES (3)\nGO"
	r, err := importer.ImportWith(src, importer.MSSQL, importer.ImportOptions{IncludeRecords: true})
	if err != nil || r.Database == nil {
		t.Fatal(err, r.Diagnostics)
	}
	var names []string
	for _, tb := range r.Database.Tables {
		names = append(names, tb.Name)
	}
	if strings.Join(names, ",") != "t,u" || len(r.Database.Records) != 1 || len(r.Database.Records[0].Values) != 2 {
		t.Errorf("tables %v records %+v", names, r.Database.Records)
	}
}

// schema.rb arguments nested 200k deep (hashes, arrays, lambdas) parse in
// linear time, without recursion, in clean and parity mode.
func TestSchemaRBDeepNesting(t *testing.T) {
	wrap := func(arg string) string {
		return "ActiveRecord::Schema[7.1].define(version: 2024_01_01_000000) do\n  create_table \"t\" do |t|\n    t.string \"a\", " + arg + "\n  end\nend\n"
	}
	cases := func(n int) map[string]string {
		return map[string]string{
			"hash":   wrap("x: " + strings.Repeat("{a: ", n) + "1" + strings.Repeat("}", n)),
			"array":  wrap("default: " + strings.Repeat("[", n) + "1" + strings.Repeat("]", n)),
			"lambda": wrap("default: " + strings.Repeat("-> { ", n) + "\"now()\"" + strings.Repeat(" }", n)),
			"mixed":  wrap("x: " + strings.Repeat("[{a: -> { ", n/3) + "1" + strings.Repeat(" }}]", n/3)),
			"open":   wrap("x: " + strings.Repeat("{a: [", n)),
			"braces": wrap("x: " + strings.Repeat("{", n*5/2) + strings.Repeat("}", n*5/2)), // 1 MiB at n = 200k
		}
	}
	for name, src := range cases(200000) {
		for _, mode := range []importer.Mode{importer.ModeClean, importer.ModeParity} {
			start := time.Now()
			_, err := importer.ImportWith(src, importer.SchemaRB, importer.ImportOptions{Mode: mode})
			el := time.Since(start)
			t.Logf("%s mode=%d: %v (%v)", name, mode, el, err)
			strictBudget(t, name, el, 200*time.Millisecond)
			assertScales(t, fmt.Sprintf("%s mode=%d", name, mode), 50000, 200000, 2.5, func(n int) {
				_, _ = importer.ImportWith(cases(n)[name], importer.SchemaRB, importer.ImportOptions{Mode: mode})
			})
		}
	}
}

func TestImportMaxInputBytes(t *testing.T) {
	src := "CREATE TABLE t (a int);" + strings.Repeat(" ", 2000)
	if _, err := importer.ImportWith(src, importer.MySQL, importer.ImportOptions{MaxInputBytes: 1000}); err != importer.ErrInputTooLarge {
		t.Errorf("err = %v", err)
	}
	if _, err := importer.ImportWith(src, importer.MySQL, importer.ImportOptions{MaxInputBytes: 64 << 20}); err != nil {
		t.Errorf("err = %v", err)
	}
}
