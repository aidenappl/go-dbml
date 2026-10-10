package export_test

import (
	"errors"
	"strings"
	"testing"

	dbml "github.com/aidenappl/go-dbml"
	"github.com/aidenappl/go-dbml/export"
)

func TestGeneratedColumns(t *testing.T) {
	src := "Table t {\n  a int\n  b int [generated: 'a * 2', generated_stored: 'true', not null]\n  c int [generated: 'a + 1', default: 5]\n}\n"
	db, d := dbml.Parse(src)
	if db == nil {
		t.Fatal(d)
	}
	if g := db.Tables[0].Fields[1].Generated; g == nil || g.Expression != "a * 2" || !g.Stored {
		t.Fatalf("Generated = %+v", g)
	}
	want := map[export.Dialect][]string{
		export.MySQL:    {"`b` int GENERATED ALWAYS AS (a * 2) STORED NOT NULL", "`c` int GENERATED ALWAYS AS (a + 1) VIRTUAL\n"},
		export.Postgres: {`"b" int GENERATED ALWAYS AS (a * 2) STORED NOT NULL`, `"c" int GENERATED ALWAYS AS (a + 1) STORED` + "\n"},
		export.MSSQL:    {"[b] AS (a * 2) PERSISTED NOT NULL", "[c] AS (a + 1)\n"},
	}
	for dl, ws := range want {
		out, err := export.SQL(db, dl)
		if err != nil {
			t.Fatal(err)
		}
		for _, w := range ws {
			if !strings.Contains(out, w) {
				t.Errorf("%s: missing %q in\n%s", dl, w, out)
			}
		}
		if strings.Contains(out, "DEFAULT") {
			t.Errorf("%s: a generated column got a DEFAULT\n%s", dl, out)
		}
	}
}

func TestUnsafeExpressionRefused(t *testing.T) {
	for _, src := range []string{
		"Table t {\n  a int [check: `a > 0); DROP TABLE v; --`]\n}\n",
		"Table t {\n  a int [default: `1); DROP TABLE v`]\n}\n",
		"Table t {\n  a \"int;drop/**/table/**/v\"\n}\n",
		"Table t {\n  a int\n  indexes {\n    `a); DROP TABLE v; --`\n  }\n}\n",
		"Table t {\n  a int [generated: 'a); DROP TABLE v']\n}\n",
		"Table t {\n  a int\n}\nRecords t(a) {\n  `1); DROP TABLE v`\n}\n",
	} {
		db, d := dbml.Parse(src)
		if db == nil {
			t.Fatal(d)
		}
		for _, dl := range export.Dialects {
			out, err := export.SQL(db, dl)
			var ue *export.UnsafeExpressionError
			if !errors.Is(err, export.ErrUnsafeExpression) || !errors.As(err, &ue) || out != "" {
				t.Errorf("%s %q: err = %v, out %q", dl, src, err, out)
			}
		}
	}
}

func TestStringLiteralEscapers(t *testing.T) {
	db, _ := dbml.Parse("Table t {\n  a text [note: 'x\\\\y', default: 'z\\\\']\n}\n")
	pg, _ := export.SQL(db, export.Postgres)
	if !strings.Contains(pg, `DEFAULT E'z\\'`) || !strings.Contains(pg, `IS E'x\\y'`) {
		t.Errorf("postgres backslash literals not E'...':\n%s", pg)
	}
	ms, _ := export.SQL(db, export.MSSQL)
	if !strings.Contains(ms, `@value = N'x\y'`) || !strings.Contains(ms, `@level1name = N't'`) {
		t.Errorf("mssql extended property not N'...':\n%s", ms)
	}
}

func TestPostgresMultiWordTypes(t *testing.T) {
	db, _ := dbml.Parse("Table t {\n  a \"character varying(255)\"\n  b \"timestamp(6) with time zone\"\n  c \"double precision\"\n  d \"bit varying(8)\"\n  e \"VARCHAR(20)\"\n  f \"MyType\"\n}\n")
	out, _ := export.SQL(db, export.Postgres)
	for _, w := range []string{`"a" character varying(255)`, `"b" timestamp(6) with time zone`, `"c" double precision`, `"d" bit varying(8)`, `"e" VARCHAR(20)`, `"f" "MyType"`} {
		if !strings.Contains(out, w) {
			t.Errorf("missing %q in\n%s", w, out)
		}
	}
}
