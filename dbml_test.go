package dbml

import (
	"encoding/json"
	"testing"
	"testing/fstest"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantDB  bool
		wantErr bool
	}{
		{"valid", "Table users {\n id int [pk]\n}", true, false},
		{"empty", "", true, false},
		{"error", "Table {", false, true},
		{"info only", "Table a { id int [pk] }\nTable b { aid int [ref: > a.id] }", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, diags := Parse(tt.src)
			if (db != nil) != tt.wantDB || diags.HasErrors() != tt.wantErr {
				t.Fatalf("db=%v errors=%v", db != nil, diags)
			}
		})
	}
}

func TestJSONShape(t *testing.T) {
	db, _ := Parse("Table users {\n id int [pk]\n}")
	b, err := json.Marshal(db)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"schemas", "tables", "notes", "refs", "deps", "enums", "tableGroups", "aliases", "tablePartials", "records", "externals", "diagramViews", "token"} {
		if _, ok := m[k]; !ok {
			t.Fatalf("missing key %q", k)
		}
	}
	field := m["tables"].([]any)[0].(map[string]any)["fields"].([]any)[0].(map[string]any)
	if field["pk"] != true || field["type"].(map[string]any)["type_name"] != "int" {
		t.Fatalf("bad field %v", field)
	}
}

func TestParseFS(t *testing.T) {
	fsys := fstest.MapFS{
		"schema/base.dbml": {Data: []byte("Table users { id int [pk] }")},
		"schema/main.dbml": {Data: []byte("use * from './base'\nTable posts { uid int [ref: > users.id] }")},
	}
	r := ParseFS(fsys, "schema/main.dbml")
	if r.Database == nil {
		t.Fatalf("errors: %v", r.Diagnostics.Errors())
	}
	if len(r.Database.Tables) != 2 || len(r.Database.Refs) != 1 {
		t.Fatalf("want 2 tables and 1 ref, got %d/%d", len(r.Database.Tables), len(r.Database.Refs))
	}
}

func TestMonacoMarkers(t *testing.T) {
	_, diags := Parse("Table a {\n  id int [ref: > missing.id]\n}")
	ms := diags.Markers()
	if len(ms) == 0 || ms[0].StartLineNumber != 2 || ms[0].Severity != 8 {
		t.Fatalf("bad markers %+v", ms)
	}
}

func FuzzParse(f *testing.F) {
	seeds := []string{
		"", "Table users {\n id int [pk]\n}",
		"Ref: a.(x, y) <> b.(x, y) [delete: cascade]",
		"Table t {\n ~p\n records { 1, 'a', , `x` }\n}",
		"use {\n table a as b\n} from './x'",
		"Metadata Column t.c {\n k: 'v'\n}",
		"DiagramView v {\n Tables { * }\n}",
		"Dep d {\n a -> b\n note: 'x'\n}",
		"'''unterminated", "Table [ { ( \"", "enum e { a [note: 'x'] }",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	// hostile nesting just past the depth limit
	for _, s := range deepInputs(300) {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, src string) {
		db, diags := Parse(src)
		if db == nil && !diags.HasErrors() {
			t.Fatal("nil database without errors")
		}
		if db != nil {
			if _, err := json.Marshal(db); err != nil {
				t.Fatal(err)
			}
		}
		_ = ParseSyntax(src)
	})
}
