package model

import (
	"testing"

	"github.com/aidenappl/go-dbml/diag"
)

func mustCompile(t *testing.T, src string) *Database {
	t.Helper()
	r := Compile(src, "/main.dbml")
	if errs := r.Diagnostics.Errors(); len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	return r.Database
}

func deref(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}

func TestConstructs(t *testing.T) {
	tests := []struct {
		name  string
		src   string
		check func(t *testing.T, db *Database)
	}{
		{"project", "Project p {\n database_type: 'PostgreSQL'\n Note: 'x'\n}", func(t *testing.T, db *Database) {
			if db.Project == nil || deref(db.Project.Name) != "p" || db.Project.Extra["database_type"] != "PostgreSQL" || db.Project.Note.Value != "x" {
				t.Fatalf("bad project %+v", db.Project)
			}
		}},
		{"schema default public", "Table core.users { id int }\nTable posts { id int }", func(t *testing.T, db *Database) {
			if deref(db.Tables[0].SchemaName) != "core" || db.Tables[1].SchemaName != nil {
				t.Fatal("bad schemas")
			}
		}},
		{"column settings", "Table t {\n id int [pk, increment]\n n varchar(255) [not null, unique, default: 'x', note: 'hi']\n c int [null, default: -1.5]\n e timestamp [default: `now()`]\n}", func(t *testing.T, db *Database) {
			f := db.Tables[0].Fields
			if !*f[0].PK || f[0].Increment == nil || f[1].Type.TypeName != "varchar(255)" || deref(f[1].Type.Args) != "255" {
				t.Fatal("bad id/name")
			}
			if !*f[1].NotNull || !*f[1].Unique || f[1].DBDefault.Value != "x" || f[1].Note.Value != "hi" {
				t.Fatal("bad name settings")
			}
			if *f[2].NotNull || f[2].DBDefault.Type != "number" || f[2].DBDefault.Value != -1.5 {
				t.Fatalf("bad c %+v", f[2].DBDefault)
			}
			if f[3].DBDefault.Type != "expression" || f[3].DBDefault.Value != "now()" {
				t.Fatal("bad expression default")
			}
		}},
		{"legacy inline settings", "Table t {\n id int pk unique\n}", func(t *testing.T, db *Database) {
			if !*db.Tables[0].Fields[0].PK || !*db.Tables[0].Fields[0].Unique {
				t.Fatal("legacy pk/unique not applied")
			}
		}},
		{"alias", "Table very_long as U { id int }\nTable p { uid int [ref: > U.id] }", func(t *testing.T, db *Database) {
			if len(db.Aliases) != 1 || db.Aliases[0].Value.TableName != "very_long" {
				t.Fatal("bad alias")
			}
			if db.Refs[0].Endpoints[0].TableName != "very_long" {
				t.Fatal("alias not resolved in ref")
			}
		}},
		{"refs all forms", "Table a { id int [pk]\n x int\n y int }\nTable b { aid int [ref: > a.id]\n x int\n y int }\nRef r1: a.id < b.x [delete: cascade, update: no action, color: #fff, inactive]\nRef {\n a.(x, y) - b.(x, y)\n}", func(t *testing.T, db *Database) {
			if len(db.Refs) != 3 {
				t.Fatalf("want 3 refs, got %d", len(db.Refs))
			}
			var named *Ref
			for _, r := range db.Refs {
				if deref(r.Name) == "r1" {
					named = r
				}
			}
			if named == nil || deref(named.OnDelete) != "cascade" || deref(named.OnUpdate) != "no action" || deref(named.Color) != "#fff" || named.Inactive == nil {
				t.Fatalf("bad settings %+v", named)
			}
			if named.Endpoints[0].Relation != "1" || named.Endpoints[1].Relation != "*" {
				t.Fatal("bad relation")
			}
		}},
		{"optional ref ops", "Table a { id int [pk] }\nTable b { aid int [ref: >? a.id]\n bid int [ref: ?<> a.id] }", func(t *testing.T, db *Database) {
			f := db.Tables[1].Fields
			if f[0].InlineRefs[0].Relation != ">?" || f[1].InlineRefs[0].Relation != "?<>" {
				t.Fatal("bad inline relations")
			}
			for _, r := range db.Refs {
				if r.Endpoints[1].FieldNames[0] == "aid" && (r.Endpoints[0].Relation != "0..1" || r.Endpoints[1].Relation != "*") {
					t.Fatalf("bad >? cardinalities %s %s", r.Endpoints[0].Relation, r.Endpoints[1].Relation)
				}
			}
		}},
		{"cross schema ref", "Table core.users { id int }\nTable blog.posts { uid int [ref: > core.users.id] }", func(t *testing.T, db *Database) {
			ep := db.Refs[0].Endpoints[0]
			if deref(ep.SchemaName) != "core" || ep.TableName != "users" {
				t.Fatalf("bad endpoint %+v", ep)
			}
		}},
		{"enum", "enum v2.job_status {\n created [note: 'waiting']\n \"A+\"\n}\nTable j { s v2.job_status }", func(t *testing.T, db *Database) {
			e := db.Enums[0]
			if deref(e.SchemaName) != "v2" || len(e.Values) != 2 || e.Values[0].Note.Value != "waiting" || e.Values[1].Name != "A+" {
				t.Fatalf("bad enum %+v", e)
			}
			if deref(db.Tables[0].Fields[0].Type.SchemaName) != "v2" {
				t.Fatal("bad enum column type")
			}
		}},
		{"indexes", "Table t {\n id int\n c varchar\n indexes {\n  (id, c) [pk]\n  c [name: 'c_idx', type: hash, unique, note: 'n']\n  (`lower(c)`, id)\n }\n}", func(t *testing.T, db *Database) {
			ix := db.Tables[0].Indexes
			if len(ix) != 3 || !*ix[0].PK || deref(ix[1].Name) != "c_idx" || deref(ix[1].Type) != "hash" || ix[1].Note.Value != "n" {
				t.Fatal("bad indexes")
			}
			if ix[2].Columns[0].Type != "expression" || ix[2].Columns[1].Type != "column" {
				t.Fatal("bad expression index")
			}
		}},
		{"composite pk via columns", "Table t {\n a int [pk]\n b int [pk]\n}", func(t *testing.T, db *Database) {
			tb := db.Tables[0]
			if *tb.Fields[0].PK || len(tb.Indexes) != 1 || !*tb.Indexes[0].PK || tb.Indexes[0].Token.Start.Offset != -1 {
				t.Fatal("composite pk not synthesised")
			}
		}},
		{"checks", "Table t {\n w int [check: `w > 0`]\n checks {\n  `w < 10` [name: 'small']\n }\n}", func(t *testing.T, db *Database) {
			tb := db.Tables[0]
			if tb.Fields[0].Checks[0].Expression != "w > 0" || deref(tb.Checks[0].Name) != "small" {
				t.Fatal("bad checks")
			}
		}},
		{"table group", "Table a { id int }\nTable b { id int }\nTableGroup g [color: #abc, note: 'n'] {\n a\n b\n}", func(t *testing.T, db *Database) {
			g := db.TableGroups[0]
			if deref(g.Name) != "g" || deref(g.Color) != "#abc" || g.Note.Value != "n" || len(g.Tables) != 2 {
				t.Fatal("bad group")
			}
		}},
		{"sticky note", "Note n [color: none] {\n '''\n   multi\n   line\n '''\n}", func(t *testing.T, db *Database) {
			n := db.Notes[0]
			if n.Name != "n" || deref(n.Color) != "none" || n.Content != "multi\nline\n" {
				t.Fatalf("bad note %q", n.Content)
			}
		}},
		{"table partial injection", "TablePartial base [headercolor: #f00] {\n id int [pk]\n created timestamp\n}\nTable users {\n ~base\n name varchar\n created date\n}", func(t *testing.T, db *Database) {
			tb := db.Tables[0]
			if len(tb.Partials) != 1 || tb.Partials[0].Order != 0 {
				t.Fatal("bad partial injection")
			}
			rt := db.ResolveTable(tb)
			var names []string
			for _, c := range rt.Columns {
				names = append(names, c.Name+"/"+c.InjectedFrom)
			}
			want := []string{"id/base", "name/", "created/"}
			for i := range want {
				if names[i] != want[i] {
					t.Fatalf("resolved columns %v, want %v", names, want)
				}
			}
			if deref(rt.HeaderColor) != "#f00" {
				t.Fatal("headercolor not injected")
			}
		}},
		{"partial precedence last wins", "TablePartial a { x int }\nTablePartial b { x varchar }\nTable t {\n ~a\n ~b\n id int\n}", func(t *testing.T, db *Database) {
			rt := db.ResolveTable(db.Tables[0])
			if len(rt.Columns) != 2 || rt.Columns[0].Type.TypeName != "varchar" || rt.Columns[0].InjectedFrom != "b" {
				t.Fatalf("want x from b, got %+v", rt.Columns[0])
			}
		}},
		{"records", "enum S { on\n off }\nTable t {\n id int\n ok bool\n s S\n at date\n records [example] {\n  1, true, S.on, '2024-01-15'\n  2, 'N', 'off', null\n }\n}", func(t *testing.T, db *Database) {
			r := db.Records[0]
			if r.Example == nil || len(r.Values) != 2 {
				t.Fatal("bad records")
			}
			row := r.Values[1]
			if row[0].Value != 2.0 || row[0].Type != "integer" || row[1].Value != false || row[2].Value != "off" || row[3].Value != nil {
				t.Fatalf("bad row %+v %+v %+v %+v", row[0], row[1], row[2], row[3])
			}
			if r.Values[0][3].Value != "2024-01-15" || r.Values[0][3].Type != "date" {
				t.Fatal("bad date")
			}
		}},
		{"top-level records with expression", "Table t { id int\n n varchar }\nrecords t(id, n) {\n 1, `now()`\n}", func(t *testing.T, db *Database) {
			v := db.Records[0].Values[0][1]
			if v.Type != "expression" || v.Value != "now()" {
				t.Fatal("bad expression value")
			}
		}},
		{"deps", "Table a { id int }\nTable b { id int [dep: <- a.id] }\nTable c [dep: <- b] { id int }\nTable e { id int }\nDep: a -> b\nDep d [color: #79AD51] {\n a -> e\n b -> e\n note: 'n'\n owner: 'x'\n priority: 2\n}", func(t *testing.T, db *Database) {
			if len(db.Deps) != 4 {
				t.Fatalf("want 4 deps, got %d", len(db.Deps))
			}
			col := db.Tables[1].Fields[0]
			if len(col.InlineDeps) != 1 || col.InlineDeps[0].Direction != "<-" || col.InlineDeps[0].TableName != "a" {
				t.Fatal("bad inline dep")
			}
			var block *Dep
			for _, d := range db.Deps {
				if deref(d.Name) == "d" {
					block = d
				}
			}
			if block == nil || deref(block.Color) != "#79AD51" || block.Note.Value != "n" || block.Metadata["owner"] != "x" || block.Metadata["priority"] != 2.0 {
				t.Fatalf("bad dep block %+v", block)
			}
		}},
		{"diagram view", "Table a { id int }\nTableGroup g { a }\nNote n { 'x' }\nDiagramView v {\n Tables { a }\n TableGroups { * }\n}", func(t *testing.T, db *Database) {
			ve := db.DiagramViews[0].VisibleEntities
			if len(ve.Tables) != 1 || len(ve.TableGroups) != 1 || ve.TableGroups[0].Name != "g" || ve.StickyNotes != nil || ve.Schemas == nil {
				t.Fatalf("bad view %+v", ve)
			}
		}},
		{"custom properties and metadata block", "Table users [owner: 'a'] {\n id int [pii: 'true']\n}\nMetadata Table users {\n owner: 'b'\n note: 'meta note'\n}\nMetadata Column users.id {\n masking: 'partial'\n}", func(t *testing.T, db *Database) {
			tb := db.Tables[0]
			if tb.Metadata["owner"] != "b" || tb.Note.Value != "meta note" {
				t.Fatalf("bad table metadata %+v", tb.Metadata)
			}
			if tb.Fields[0].Metadata["pii"] != "true" || tb.Fields[0].Metadata["masking"] != "partial" {
				t.Fatalf("bad column metadata %+v", tb.Fields[0].Metadata)
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, mustCompile(t, tt.src))
		})
	}
}

func TestErrors(t *testing.T) {
	tests := []struct {
		name string
		src  string
		code diag.Code
	}{
		{"unknown table in ref", "Table a { id int }\nRef: a.id > b.id", diag.BINDING_ERROR},
		{"unknown column", "Table a { id int }\nTable b { x int [ref: > a.nope] }", diag.BINDING_ERROR},
		{"duplicate table", "Table a { id int }\nTable a { id int }", diag.DUPLICATE_NAME},
		{"duplicate column", "Table a {\n id int\n id int\n}", diag.DUPLICATE_COLUMN_NAME},
		{"empty table", "Table a {\n}", diag.EMPTY_TABLE},
		{"circular ref", "Table a { id int }\nTable b { id int }\nRef: a.id > b.id\nRef: b.id > a.id", diag.CIRCULAR_REF},
		{"same endpoint", "Table a { id int }\nRef: a.id > a.id", diag.SAME_ENDPOINT},
		{"bad ref setting", "Table a { id int }\nTable b { id int }\nRef: a.id > b.id [delete: explode]", diag.INVALID_REF_SETTING_VALUE},
		{"unknown partial", "Table a {\n ~nope\n id int\n}", diag.BINDING_ERROR},
		{"table in two groups", "Table a { id int }\nTableGroup g1 { a }\nTableGroup g2 { a }", diag.TABLE_REAPPEAR_IN_TABLEGROUP},
		{"two projects", "Project a {}\nProject b {}", diag.PROJECT_REDEFINED},
		{"syntax error", "Table a { id int [pk }", diag.UNEXPECTED_TOKEN},
		{"record arity", "Table a { id int\n n int }\nrecords a(id, n) {\n 1\n}", diag.INVALID_RECORDS_FIELD},
		{"duplicate records", "Table a { id int }\nrecords a(id) { 1 }\nrecords a(id) { 2 }", diag.DUPLICATE_RECORDS_FOR_TABLE},
		{"dep mixed downstream", "Table a { id int }\nTable b { id int }\nTable c { id int }\nDep {\n a -> b\n a -> c\n}", diag.DEP_MIXED_DOWNSTREAM_TABLES},
		{"metadata missing target", "Table a { id int }\nMetadata Table nope {\n owner: 'x'\n}", diag.BINDING_ERROR},
		{"missing import", "use * from './nope'\nTable a { id int }", diag.NONEXISTENT_MODULE},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Compile(tt.src, "/main.dbml")
			if r.Database != nil {
				t.Fatal("want nil database on error")
			}
			for _, d := range r.Diagnostics.Errors() {
				if d.Code == tt.code {
					if d.Range.Start.Line < 1 || d.Range.Start.Column < 1 {
						t.Fatalf("bad range %+v", d.Range)
					}
					return
				}
			}
			t.Fatalf("want %s, got %v", tt.code, r.Diagnostics)
		})
	}
}

func TestMultiFile(t *testing.T) {
	files := map[string]string{
		"/lib/users.dbml": "Table users {\n id int [pk]\n}\nTable hidden { id int }",
		"/lib/index.dbml": "reuse { table users } from './users'",
		"/main.dbml":      "use * from './lib/index'\nTable posts {\n uid int [ref: > users.id]\n}",
	}
	r := CompileProject("/main.dbml", func(p string) (string, bool) { s, ok := files[p]; return s, ok })
	if r.Database == nil {
		t.Fatalf("errors: %v", r.Diagnostics.Errors())
	}
	if len(r.Database.Tables) != 2 || len(r.Database.Externals.Tables) != 1 {
		t.Fatalf("want posts + users (hidden not re-exported), got %d tables", len(r.Database.Tables))
	}
	if r.Database.Externals.Tables[0].Filepath != "/lib/users.dbml" {
		t.Fatal("bad external filepath")
	}
}

func TestMultiFileImportedRefs(t *testing.T) {
	tests := []struct {
		name     string
		main     string
		wantRefs int
	}{
		{"both tables visible", "use * from './base'\nTable x { id int }", 1},
		{"one table hidden", "use {\n table users\n} from './base'\nTable x { id int }", 0},
	}
	base := "Table users {\n id int [pk]\n rid int\n}\nTable roles {\n id int [pk]\n}\nRef: users.rid > roles.id"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := map[string]string{"/base.dbml": base, "/main.dbml": tt.main}
			r := CompileProject("/main.dbml", func(p string) (string, bool) { s, ok := files[p]; return s, ok })
			if r.Database == nil {
				t.Fatalf("errors: %v", r.Diagnostics.Errors())
			}
			if len(r.Database.Refs) != tt.wantRefs {
				t.Fatalf("want %d refs, got %d", tt.wantRefs, len(r.Database.Refs))
			}
		})
	}
}

func TestResolveTableIndexOrder(t *testing.T) {
	r := Compile("TablePartial p {\n  a int\n  indexes {\n    a [name: 'from_partial']\n  }\n}\nTable t {\n  b int\n  ~p\n  indexes {\n    b [name: 'own']\n  }\n}\n", "/m.dbml")
	if r.Database == nil {
		t.Fatal(r.Diagnostics)
	}
	rt := r.Database.ResolveTable(r.Database.Tables[0])
	if len(rt.AllIndexes) != 2 || *rt.AllIndexes[0].Name != "from_partial" || *rt.AllIndexes[1].Name != "own" {
		t.Fatalf("want injected index first (@dbml/core order), got %v %v", *rt.AllIndexes[0].Name, *rt.AllIndexes[1].Name)
	}
}
