package export_test

import (
	"strings"
	"testing"

	dbml "github.com/aidenappl/go-dbml"
	"github.com/aidenappl/go-dbml/emit"
	"github.com/aidenappl/go-dbml/export"
)

// adversarialCases are DBML models whose names/notes/defaults/enum values
// try to break out of the generated DDL (reviewer probes, Oct 2026).
func adversarialCases() map[string]string {
	I, S := emit.Ident, emit.String
	return map[string]string{
		"mysql_tablename":  "Table " + I("x` (a int); DROP TABLE victim; CREATE TABLE `y") + " {\n  id int\n}\n",
		"pg_tablename":     "Table " + I(`x" (a int); DROP TABLE victim; CREATE TABLE "y`) + " {\n  id int\n}\n",
		"ms_tablename":     "Table " + I("x] (a int); DROP TABLE victim; CREATE TABLE [y") + " {\n  id int\n}\n",
		"mysql_colname":    "Table t {\n  " + I("a` int); DROP TABLE victim; CREATE TABLE y (`b") + " int\n}\n",
		"pg_colname":       "Table t {\n  " + I(`a" int); DROP TABLE victim; CREATE TABLE y ("b`) + " int\n}\n",
		"tablenote":        "Table t {\n  id int\n  Note: " + S(`\'; DROP TABLE victim; -- `) + "\n}\n",
		"colnote":          "Table t {\n  id int [note: " + S(`\'; DROP TABLE victim; -- `) + "]\n}\n",
		"enumval":          "Enum e {\n  " + I("a'); DROP TABLE victim; CREATE TABLE e2 (x int); --") + "\n}\nTable t {\n  s e\n}\n",
		"strdefault":       "Table t {\n  a varchar(200) [default: " + S("x'); DROP TABLE victim; CREATE TABLE dummy (q varchar(9) DEFAULT 'y") + "]\n}\n",
		"backslashdefault": "Table t {\n  a varchar(200) [default: " + S(`x\'); DROP TABLE victim; -- `) + "]\n}\n",
		"mysql_idxname":    "Table t {\n  a int\n  indexes {\n    a [name: " + S("i` ON t (a); DROP TABLE victim; CREATE INDEX `j") + "]\n  }\n}\n",
		"pg_idxname":       "Table t {\n  a int\n  indexes {\n    a [name: " + S(`i" ON t (a); DROP TABLE victim; CREATE INDEX "j`) + "]\n  }\n}\n",
		"mysql_refname":    "Table p {\n  id int [pk]\n}\nTable t {\n  pid int\n}\nRef " + I("r` FOREIGN KEY (pid) REFERENCES p (id); DROP TABLE victim;#") + ": t.pid > p.id\n",
		"pg_refname":       "Table p {\n  id int [pk]\n}\nTable t {\n  pid int\n}\nRef " + I(`r" FOREIGN KEY (pid) REFERENCES p (id); DROP TABLE victim;--`) + ": t.pid > p.id\n",
		"junction":         "Table " + I("a`\"]") + " {\n  id int [pk]\n}\nTable b {\n  id int [pk]\n}\nRef: " + I("a`\"]") + ".id <> b.id\n",
		"record_text":      "Table t {\n  a varchar(50)\n}\nRecords t(a) {\n  " + S(`\'); DROP TABLE victim; -- `) + "\n}\n",
		"checkname":        "Table t {\n  a int\n  checks {\n    `a > 0` [name: " + S("c`\"]); DROP TABLE victim; --") + "]\n  }\n}\n",
	}
}

// TestAdversarialNamesStayQuoted is the DB-less check: no exported
// statement may contain the injected DROP outside a quoted token. We strip
// every quoted identifier/literal and require that "DROP" never survives.
func TestAdversarialNamesStayQuoted(t *testing.T) {
	for name, src := range adversarialCases() {
		db, diags := dbml.Parse(src)
		if db == nil {
			t.Fatalf("%s: %v", name, diags)
		}
		for _, d := range export.Dialects {
			sql, err := export.SQL(db, d)
			if err != nil {
				t.Fatalf("%s/%s: %v", name, d, err)
			}
			if bare := stripQuoted(sql, d); strings.Contains(strings.ToUpper(bare), "DROP TABLE") {
				t.Errorf("%s/%s: injected statement outside quotes:\n%s", name, d, sql)
			}
		}
	}
}

func TestInvalidActionsAndIndexTypesRejected(t *testing.T) {
	for _, src := range []string{
		"Table a {\n  id int [pk]\n}\nTable b {\n  a int\n}\nRef: b.a > a.id [delete: \"cascade; DROP TABLE x\"]\n",
		"Table a {\n  id int\n  indexes {\n    id [type: \"hash; DROP TABLE x\"]\n  }\n}\n",
	} {
		db, _ := dbml.Parse(src)
		if db == nil {
			continue // the parser already refuses it
		}
		if _, err := export.SQL(db, export.Postgres); err == nil {
			t.Errorf("want an error for %q", src)
		}
	}
}

// stripQuoted removes quoted identifiers and string literals following the
// dialect's lexical rules (doubled quotes, MySQL backslash escapes).
func stripQuoted(sql string, d export.Dialect) string {
	var b strings.Builder
	for i := 0; i < len(sql); i++ {
		c := sql[i]
		var end byte
		switch {
		case c == '\'':
			end = '\''
		case c == '`' && d == export.MySQL, c == '"' && d != export.MySQL:
			end = c
		case c == '"' && d == export.MSSQL:
			end = '"'
		case c == '[' && d == export.MSSQL:
			end = ']'
		default:
			b.WriteByte(c)
			continue
		}
		i++
		for ; i < len(sql); i++ {
			if d == export.MySQL && end == '\'' && sql[i] == '\\' {
				i++
				continue
			}
			if sql[i] == end {
				if i+1 < len(sql) && sql[i+1] == end {
					i++
					continue
				}
				break
			}
		}
		b.WriteString("Q")
	}
	return b.String()
}
