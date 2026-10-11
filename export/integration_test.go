//go:build integration

package export_test

// Executes exported DDL against real MariaDB 11.4, PostgreSQL 16 (also
// with standard_conforming_strings=off) and SQL Server 2022 (Docker
// required; SQL Server runs as linux/amd64):
//
//	go test -tags integration -run 'Against' -v ./export
//
// TestDDLAgainstRealDatabases: adversarial names/strings (adversarialCases)
// must execute without error and leave the sentinel table `victim` plus
// exactly the model's tables.
//
// TestSmugglingAgainstRealDatabases: raw-expression smuggling — every
// testdata/smuggle input (MySQL, PostgreSQL, T-SQL, schema.rb), imported in
// clean mode and exported to every dialect, plus hand-written DBML
// fragments (smugglingCases). The export may be refused
// (ErrUnsafeExpression) or may fail to execute; it must never drop
// `victim` or create a table the model does not have.

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
	"unicode/utf16"

	dbml "github.com/aidenappl/go-dbml"
	"github.com/aidenappl/go-dbml/export"
	"github.com/aidenappl/go-dbml/importer"
	"github.com/aidenappl/go-dbml/model"
)

func docker(stdin string, args ...string) (string, error) {
	cmd := exec.Command("docker", args...)
	cmd.Stdin = strings.NewReader(stdin)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("%v: %s %s", err, errb.String(), out.String())
	}
	return out.String(), nil
}

func startContainer(t *testing.T, args ...string) string {
	t.Helper()
	id, err := docker("", append([]string{"run", "-d", "--rm"}, args...)...)
	if err != nil {
		t.Skipf("docker unavailable: %v", err)
	}
	id = strings.TrimSpace(id)
	t.Cleanup(func() { _, _ = docker("", "rm", "-f", id) })
	return id
}

func waitFor(t *testing.T, probe func() error) {
	t.Helper()
	deadline := time.Now().Add(180 * time.Second)
	for {
		err := probe()
		if err == nil {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("database did not start: %v", err)
		}
		time.Sleep(time.Second)
	}
}

func hexLines(out string) []string {
	var names []string
	for _, l := range strings.Fields(out) {
		b, err := hex.DecodeString(strings.TrimSpace(l))
		if err == nil {
			names = append(names, string(b))
		}
	}
	sort.Strings(names)
	return names
}

// engine is one live database.
type engine struct {
	name    string
	dialect export.Dialect
	// reset recreates database t holding only the sentinel table victim
	reset func() error
	// exec runs a script in t
	exec func(script string) error
	// tables lists the table names in t
	tables func() ([]string, error)
}

func mariadbEngine(t *testing.T) *engine {
	id := startContainer(t, "-e", "MARIADB_ROOT_PASSWORD=pw", "mariadb:11.4")
	sql := func(db, stdin string) (string, error) {
		// TCP for the same reason (the init-time server skips networking)
		args := []string{"exec", "-i", id, "mariadb", "-h", "127.0.0.1", "--protocol=tcp", "-uroot", "-ppw", "-N", "-B"}
		if db != "" {
			args = append(args, db)
		}
		return docker(stdin, args...)
	}
	waitFor(t, func() error { _, err := sql("", "SELECT 1"); return err })
	return &engine{
		name: "mariadb", dialect: export.MySQL,
		reset: func() error {
			// schemas are databases: drop every non-system one (FKs may
			// cross them)
			out, err := sql("", "SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('mysql','information_schema','performance_schema','sys')")
			if err != nil {
				return err
			}
			script := "SET FOREIGN_KEY_CHECKS = 0;"
			for _, n := range strings.Fields(out) {
				script += " DROP DATABASE `" + strings.ReplaceAll(n, "`", "``") + "`;"
			}
			_, err = sql("", script+" SET FOREIGN_KEY_CHECKS = 1; CREATE DATABASE t; CREATE TABLE t.victim (id int);")
			return err
		},
		exec: func(s string) error { _, err := sql("t", s); return err },
		tables: func() ([]string, error) {
			out, err := sql("", "SELECT HEX(table_name) FROM information_schema.tables WHERE table_schema = 't'")
			return hexLines(out), err
		},
	}
}

func postgresEngine(t *testing.T, scsOff bool) *engine {
	id := startContainer(t, "-e", "POSTGRES_PASSWORD=pw", "postgres:16")
	sql := func(db, stdin string) (string, error) {
		args := []string{"exec", "-i"}
		if scsOff {
			args = append(args, "-e", "PGOPTIONS=-c standard_conforming_strings=off -c escape_string_warning=off")
		}
		// TCP, not the socket: the image's init-time temporary server
		// listens only on the socket, so a probe through it can succeed
		// just before that server is restarted
		args = append(args, id, "psql", "-h", "127.0.0.1", "-U", "postgres", "-d", db, "-v", "ON_ERROR_STOP=1", "-qtA")
		return docker(stdin, args...)
	}
	waitFor(t, func() error { _, err := sql("postgres", "SELECT 1"); return err })
	name := "postgres"
	if scsOff {
		name = "postgres_scs_off"
		if out, _ := sql("postgres", "SHOW standard_conforming_strings;"); strings.TrimSpace(out) != "off" {
			t.Fatalf("standard_conforming_strings = %q, want off", out)
		}
	}
	return &engine{
		name: name, dialect: export.Postgres,
		reset: func() error {
			if _, err := sql("postgres", "DROP DATABASE IF EXISTS t;"); err != nil {
				return err
			}
			if _, err := sql("postgres", "CREATE DATABASE t;"); err != nil {
				return err
			}
			_, err := sql("t", "CREATE TABLE victim (id int);")
			return err
		},
		exec: func(s string) error { _, err := sql("t", s); return err },
		tables: func() ([]string, error) {
			out, err := sql("t", "SELECT encode(convert_to(tablename, 'UTF8'), 'hex') FROM pg_tables WHERE schemaname NOT IN ('pg_catalog', 'information_schema');")
			return hexLines(out), err
		},
	}
}

// mssqlEngine starts SQL Server. quotedIdent runs scripts with
// QUOTED_IDENTIFIER ON (sqlcmd -I, like SSMS); off is sqlcmd's default and
// the worst case for smuggling ("..." is then a string).
func mssqlEngine(t *testing.T, quotedIdent bool) *engine {
	const pw = "Str0ng!Passw0rd"
	id := startContainer(t, "--platform", "linux/amd64", "-e", "ACCEPT_EULA=Y", "-e", "MSSQL_SA_PASSWORD="+pw,
		"mcr.microsoft.com/mssql/server:2022-latest")
	sqlcmd := []string{"exec", "-i", id, "/opt/mssql-tools18/bin/sqlcmd", "-C", "-S", "localhost", "-U", "sa", "-P", pw}
	q := func(db, query string) (string, error) {
		args := append(append([]string{}, sqlcmd...), "-b", "-h", "-1", "-W")
		if db != "" {
			args = append(args, "-d", db)
		}
		return docker("", append(args, "-Q", query)...)
	}
	waitFor(t, func() error { _, err := q("", "SELECT 1"); return err })
	return &engine{
		name: "mssql", dialect: export.MSSQL,
		reset: func() error {
			if _, err := q("", "IF DB_ID('t') IS NOT NULL BEGIN ALTER DATABASE t SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE t; END; CREATE DATABASE t;"); err != nil {
				return err
			}
			_, err := q("t", "CREATE TABLE victim (id int)")
			return err
		},
		exec: func(s string) error {
			if _, err := docker(s, "exec", "-i", id, "sh", "-c", "cat > /tmp/ddl.sql"); err != nil {
				return err
			}
			// sqlcmd splits batches at GO lines and, like SSMS, carries on
			// after a failed batch (the worst case for smuggling); report
			// any error message.
			args := append(append([]string{}, sqlcmd...), "-d", "t", "-i", "/tmp/ddl.sql")
			if quotedIdent {
				args = append(args, "-I")
			}
			out, err := docker("", args...)
			if err == nil && strings.Contains(out, "Msg ") {
				err = errors.New(out)
			}
			return err
		},
		tables: func() ([]string, error) {
			// sysname is nvarchar: the hex is UTF-16LE
			out, err := q("t", "SET NOCOUNT ON; SELECT CONVERT(varchar(max), CAST(o.name AS varbinary(max)), 2) FROM sys.tables o")
			var names []string
			for _, n := range hexLines(out) {
				u := make([]uint16, len(n)/2)
				for i := range u {
					u[i] = uint16(n[2*i]) | uint16(n[2*i+1])<<8
				}
				names = append(names, string(utf16.Decode(u)))
			}
			sort.Strings(names)
			return names, err
		},
	}
}

func modelTables(db *model.Database) []string {
	var names []string
	for _, tb := range db.Tables {
		names = append(names, tb.Name)
	}
	for _, r := range db.Refs {
		if r.Endpoints[0].Relation == "*" && r.Endpoints[1].Relation == "*" {
			names = append(names, r.Endpoints[0].TableName+"_"+r.Endpoints[1].TableName)
		}
	}
	names = append(names, "victim")
	sort.Strings(names)
	return names
}

func TestDDLAgainstRealDatabases(t *testing.T) {
	es := []*engine{mariadbEngine(t), postgresEngine(t, false), postgresEngine(t, true), mssqlEngine(t, false)}
	for _, e := range es {
		t.Run(e.name, func(t *testing.T) {
			for name, src := range adversarialCases() {
				db, _ := dbml.Parse(src)
				ddl, err := export.SQL(db, e.dialect)
				if err != nil {
					t.Fatal(err)
				}
				if err := e.reset(); err != nil {
					t.Fatal(err)
				}
				if err := e.exec(ddl); err != nil {
					t.Errorf("%s: DDL failed: %v\n%s", name, err, ddl)
					continue
				}
				got, err := e.tables()
				if err != nil {
					t.Fatal(err)
				}
				if want := modelTables(db); strings.Join(got, "\x00") != strings.Join(want, "\x00") {
					t.Errorf("%s: tables %q, want %q", name, got, want)
				}
			}
		})
	}
}

// smugglingCases are hand-written DBML models whose raw fragments try to
// end the statement they sit in.
func smugglingCases() map[string]string {
	return map[string]string{
		"check_close":     "Table t {\n  a int [check: `a > 0)); DROP TABLE victim; CREATE TABLE pwned (i int); -- `]\n}\n",
		"check_mystr":     "Table t {\n  a varchar(9) [check: `a <> 'x\\' OR a = ')); DROP TABLE victim; CREATE TABLE pwned (i int); -- '`]\n}\n",
		"default_expr":    "Table t {\n  a int [default: `1); DROP TABLE victim; CREATE TABLE pwned (i int); --`]\n}\n",
		"default_dollar":  "Table t {\n  a text [default: `$$`]\n  b text [default: `$$); DROP TABLE victim; CREATE TABLE pwned (i int); -- $$`]\n}\n",
		"default_bracket": "Table t {\n  a int [default: `[x`]\n  b int [default: `1]); DROP TABLE victim; CREATE TABLE pwned (i int); -- `]\n}\n",
		"type_smuggle":    "Table t {\n  a \"int); DROP TABLE victim; CREATE TABLE pwned (i int); --\"\n}\n",
		"type_nospace":    "Table t {\n  a s.\"int);DROP/**/TABLE/**/victim;CREATE/**/TABLE/**/pwned(i/**/int);--\"\n}\n",
		"index_expr":      "Table t {\n  a int\n  indexes {\n    `a); DROP TABLE victim; CREATE TABLE pwned (i int); --`\n  }\n}\n",
		"go_line":         "Table t {\n  a int [check: `a > 0\nGO\nDROP TABLE victim\nCREATE TABLE pwned (i int)\nGO\nSELECT (1`]\n}\n",
		"generated":       "Table t {\n  a int\n  b int [generated: 'a); DROP TABLE victim; CREATE TABLE pwned (i int); --']\n}\n",
		"generated_str":   "Table t {\n  a varchar(9)\n  b varchar(9) [generated: 'concat(a, \\'\\\\\\') ; DROP TABLE victim; CREATE TABLE pwned (i int); -- \\')']\n}\n",
		"record_expr":     "Table t {\n  a int\n}\nRecords t(a) {\n  `1); DROP TABLE victim; CREATE TABLE pwned (i int); --`\n}\n",
		"pg_estring":      "Table t {\n  a text [check: `a <> e'\\\\'')); DROP TABLE victim; CREATE TABLE pwned (i int); -- '`]\n}\n",
		"note_backslash":  "Table t {\n  a text [note: 'x\\\\', default: 'y\\\\']\n}\nRecords t(a) {\n  'z\\\\'\n}\n",
	}
}

func TestSmugglingAgainstRealDatabases(t *testing.T) {
	type input struct {
		name string
		db   *model.Database
	}
	var inputs []input
	files, _ := filepath.Glob("testdata/smuggle/*")
	for _, p := range files {
		f := map[string]importer.Format{".sql": "", ".rb": importer.SchemaRB}[filepath.Ext(p)]
		switch {
		case strings.HasPrefix(filepath.Base(p), "my"):
			f = importer.MySQL
		case strings.HasPrefix(filepath.Base(p), "pg"):
			f = importer.Postgres
		case strings.HasPrefix(filepath.Base(p), "ms"):
			f = importer.MSSQL
		}
		src, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		r, err := importer.ImportWith(string(src), f, importer.ImportOptions{})
		if err != nil || r.Database == nil {
			t.Fatalf("%s: %v %v", p, err, r)
		}
		inputs = append(inputs, input{"import_" + filepath.Base(p), r.Database})
	}
	for name, src := range smugglingCases() {
		db, d := dbml.Parse(src)
		if db == nil {
			t.Fatalf("%s: %v", name, d)
		}
		inputs = append(inputs, input{"dbml_" + name, db})
	}
	if len(inputs) < 15 {
		t.Fatalf("only %d inputs", len(inputs))
	}
	es := []*engine{mariadbEngine(t), postgresEngine(t, false), postgresEngine(t, true), mssqlEngine(t, false)}
	for _, e := range es {
		t.Run(e.name, func(t *testing.T) {
			refused, failed, ran := 0, 0, 0
			for _, in := range inputs {
				ddl, err := export.SQL(in.db, e.dialect)
				if err != nil {
					if !errors.Is(err, export.ErrUnsafeExpression) {
						t.Errorf("%s: unexpected export error %v", in.name, err)
					}
					refused++
					continue
				}
				if err := e.reset(); err != nil {
					t.Fatal(err)
				}
				if e.exec(ddl) != nil {
					failed++
				} else {
					ran++
				}
				got, err := e.tables()
				if err != nil {
					t.Fatal(err)
				}
				allowed := map[string]bool{}
				for _, n := range modelTables(in.db) {
					allowed[n] = true
				}
				victim := false
				for _, n := range got {
					victim = victim || n == "victim"
					if !allowed[n] {
						t.Errorf("%s: SMUGGLED table %q created\n%s", in.name, n, ddl)
					}
				}
				if !victim {
					t.Errorf("%s: SMUGGLED: victim dropped\n%s", in.name, ddl)
				}
			}
			t.Logf("%s: %d inputs: %d refused at export, %d executed cleanly, %d failed to execute; sentinel intact in all", e.name, len(inputs), refused, ran, failed)
		})
	}
}

// TestCleanRoundTripAgainstRealDatabases: every testdata/roundtrip input
// (real-world-shaped dumps: mysqldump, pg_dump, SSMS, Rails schema.rb) is
// imported in clean mode (records on), exported to the same dialect
// (schema.rb: PostgreSQL, or MySQL for *mysql* files) and executed on a
// fresh database. MySQL and PostgreSQL must execute cleanly for >= 90%.
func TestCleanRoundTripAgainstRealDatabases(t *testing.T) {
	engines := map[export.Dialect]*engine{
		export.MySQL:    mariadbEngine(t),
		export.Postgres: postgresEngine(t, false),
		export.MSSQL:    mssqlEngine(t, true),
	}
	// the dumps use extensions (citext, uuid-ossp): the target database has
	// them, as the source database had
	pgReset := engines[export.Postgres].reset
	engines[export.Postgres].reset = func() error {
		if err := pgReset(); err != nil {
			return err
		}
		return engines[export.Postgres].exec("CREATE EXTENSION IF NOT EXISTS citext; CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"; CREATE EXTENSION IF NOT EXISTS hstore; CREATE EXTENSION IF NOT EXISTS pgcrypto;")
	}
	for _, src := range []string{"mysql", "postgres", "mssql", "schemarb"} {
		files, _ := filepath.Glob(filepath.Join("testdata/roundtrip", src, "*"))
		pass, total := 0, 0
		var fails, invalid []string
		for _, p := range files {
			b, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			f := importer.Format(src)
			target := map[string]export.Dialect{"mysql": export.MySQL, "postgres": export.Postgres, "mssql": export.MSSQL, "schemarb": export.Postgres}[src]
			if src == "schemarb" && strings.Contains(filepath.Base(p), "mysql") {
				target = export.MySQL
			}
			name := filepath.Base(p)
			if src != "schemarb" {
				// a source the database itself rejects is not counted
				e := engines[target]
				if err := e.reset(); err != nil {
					t.Fatal(err)
				}
				if err := e.exec(string(b)); err != nil {
					msg := err.Error()
					if len(msg) > 200 {
						msg = msg[:200]
					}
					invalid = append(invalid, name+": "+strings.ReplaceAll(msg, "\n", " "))
				}
			}
			total++
			r, err := importer.ImportWith(string(b), f, importer.ImportOptions{IncludeRecords: true})
			if err != nil || r.Database == nil {
				fails = append(fails, fmt.Sprintf("%s: import: %v %v", name, err, r))
				continue
			}
			ddl, err := export.SQL(r.Database, target)
			if err != nil {
				fails = append(fails, fmt.Sprintf("%s: export: %v", name, err))
				continue
			}
			e := engines[target]
			if err := e.reset(); err != nil {
				t.Fatal(err)
			}
			if err := e.exec(ddl); err != nil {
				msg := err.Error()
				if len(msg) > 300 {
					msg = msg[:300]
				}
				fails = append(fails, fmt.Sprintf("%s -> %s: %s", name, target, strings.ReplaceAll(msg, "\n", " ")))
				continue
			}
			pass++
		}
		rate := float64(pass) / float64(total)
		t.Logf("%s: %d/%d round trips execute cleanly (%.0f%%); %d of the source files fail on a fresh database themselves (roles, USE, missing schemas/sequences, invalid DDL)", src, pass, total, rate*100, len(invalid))
		for _, f := range invalid {
			t.Logf("  SOURCE FAILS RAW %s", f)
		}
		for _, f := range fails {
			t.Logf("  FAIL %s", f)
		}
		if (src == "mysql" || src == "postgres") && rate < 0.9 {
			t.Errorf("%s round-trip rate %.0f%% < 90%%", src, rate*100)
		}
	}
	// cross-dialect: the same imports exported with TranslateTypes to the
	// other dialects (the import records its dialect in Project)
	pass, total := 0, 0
	var fails []string
	for _, src := range []string{"mysql", "postgres", "mssql", "schemarb"} {
		files, _ := filepath.Glob(filepath.Join("testdata/roundtrip", src, "*"))
		for _, p := range files {
			b, _ := os.ReadFile(p)
			r, err := importer.ImportWith(string(b), importer.Format(src), importer.ImportOptions{IncludeRecords: true})
			if err != nil || r.Database == nil {
				continue
			}
			from := export.SourceDialectOf(r.Database)
			for _, target := range export.Dialects {
				if target == from {
					continue
				}
				total++
				res, err := export.SQLWith(r.Database, target, export.Options{TranslateTypes: true})
				if err != nil {
					fails = append(fails, fmt.Sprintf("%s -> %s: export: %v", filepath.Base(p), target, err))
					continue
				}
				e := engines[target]
				if err := e.reset(); err != nil {
					t.Fatal(err)
				}
				if err := e.exec(res.SQL); err != nil {
					msg := strings.ReplaceAll(err.Error(), "\n", " ")
					if len(msg) > 250 {
						msg = msg[:250]
					}
					fails = append(fails, fmt.Sprintf("%s -> %s: %s", filepath.Base(p), target, msg))
					continue
				}
				pass++
			}
		}
	}
	rate := float64(pass) / float64(total)
	t.Logf("cross-dialect (TranslateTypes): %d/%d execute cleanly (%.0f%%)", pass, total, rate*100)
	for _, f := range fails {
		t.Logf("  FAIL %s", f)
	}
	if rate < 0.85 {
		t.Errorf("cross-dialect rate %.0f%% < 85%%", rate*100)
	}
}

// TestDBMLCorpusAgainstRealDatabases exports every testdata/roundtrip/dbml
// model to each dialect and executes it on a fresh database.
func TestDBMLCorpusAgainstRealDatabases(t *testing.T) {
	engines := map[export.Dialect]*engine{
		export.MySQL:    mariadbEngine(t),
		export.Postgres: postgresEngine(t, false),
		export.MSSQL:    mssqlEngine(t, true),
	}
	var files []string
	all, _ := filepath.Glob("testdata/roundtrip/dbml/*.dbml")
	for _, p := range all {
		b, _ := os.ReadFile(p)
		if db, diags := dbml.Parse(string(b)); db == nil {
			t.Logf("skipped (does not compile, by design of the corpus): %s: %v", filepath.Base(p), diags.Errors()[:1])
			continue
		}
		files = append(files, p)
	}
	for _, d := range export.Dialects {
		pass := 0
		var fails []string
		for _, p := range files {
			b, _ := os.ReadFile(p)
			db, _ := dbml.Parse(string(b))
			ddl, err := export.SQL(db, d)
			if err != nil {
				fails = append(fails, filepath.Base(p)+": export: "+err.Error())
				continue
			}
			e := engines[d]
			if err := e.reset(); err != nil {
				t.Fatal(err)
			}
			if err := e.exec(ddl); err != nil {
				msg := strings.ReplaceAll(err.Error(), "\n", " ")
				if len(msg) > 250 {
					msg = msg[:250]
				}
				fails = append(fails, filepath.Base(p)+": "+msg)
				continue
			}
			pass++
		}
		t.Logf("%s: %d/%d DBML corpus exports execute cleanly", d, pass, len(files))
		for _, f := range fails {
			t.Logf("  FAIL %s", f)
		}
	}
}
