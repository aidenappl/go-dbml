package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func decode(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("invalid JSON %q: %v", s, err)
	}
	if _, ok := m["ok"]; !ok {
		t.Fatalf("no ok key: %s", s)
	}
	if m["ok"] == false {
		if e, _ := m["error"].(string); e == "" {
			t.Fatalf("ok=false without error: %s", s)
		}
	}
	return m
}

func TestAPI(t *testing.T) {
	src := "Table users {\n  id int [pk]\n}\nTable posts {\n  uid int [ref: > users.id]\n}\n"
	m := decode(t, apiParse(src, ""))
	if m["ok"] != true || m["database"] == nil || m["partial"] != nil {
		t.Errorf("parse: %v", m)
	}
	m = decode(t, apiParse("Table {", ""))
	if m["ok"] != false || len(m["diagnostics"].([]any)) == 0 || m["partial"] == nil {
		t.Errorf("parse error: %v", m)
	}
	d := m["diagnostics"].([]any)[0].(map[string]any)
	if _, ok := d["startLineNumber"]; !ok {
		t.Errorf("not a Monaco marker: %v", d)
	}
	if m := decode(t, apiExportSQL(src, "postgres", "")); m["ok"] != true || !strings.Contains(m["sql"].(string), `CREATE TABLE "users"`) {
		t.Errorf("exportSQL: %v", m)
	}
	if m := decode(t, apiExportSQL(src, "oracle", "")); m["ok"] != false {
		t.Errorf("exportSQL oracle: %v", m)
	}
	m = decode(t, apiImportSQL("CREATE TABLE a (id int PRIMARY KEY, b_id int REFERENCES a(id));", "postgres", ""))
	if m["ok"] != true || !strings.Contains(m["dbml"].(string), "Table a {") {
		t.Errorf("importSQL: %v", m)
	}
	m = decode(t, apiImportSQL("CREATE TABLE t (a int);\nINSERT INTO t VALUES (1),(2),(3);", "mysql", `{"includeRecords":true,"maxRecordRows":2}`))
	if m["ok"] != true || m["recordsTruncated"] != true {
		t.Errorf("importSQL records cap: %v", m)
	}
	db := decode(t, apiParse(src, ""))["database"]
	b, _ := json.Marshal(db)
	if m := decode(t, apiEmit(string(b))); m["ok"] != true || !strings.Contains(m["dbml"].(string), "uid int [ref: > users.id]") {
		t.Errorf("emit: %v", m)
	}
	if m := decode(t, apiEmit("{")); m["ok"] != false {
		t.Errorf("emit bad json: %v", m)
	}
	if m := decode(t, apiParse(src, "{bad")); m["ok"] != false {
		t.Errorf("bad options: %v", m)
	}
}

func TestBrowserLimits(t *testing.T) {
	var b strings.Builder
	b.WriteString("Table t {\n  a int\n}\nRecords t(a) {\n")
	for i := 0; i < 5000; i++ {
		b.WriteString("  1\n")
	}
	b.WriteString("}\n")
	m := decode(t, apiParse(b.String(), ""))
	if m["ok"] != false || !strings.Contains(apiParse(b.String(), ""), "RECORDS_TOO_LARGE") {
		t.Errorf("browser record default not enforced: %v", m["error"])
	}
	big := strings.Repeat("// x\n", (3<<19)/5) // 1.5 MiB > the 1 MiB browser default
	if m := decode(t, apiParse(big, "")); m["ok"] != false {
		t.Errorf("browser size default not enforced")
	}
}

func TestCompactTokensAndWarnings(t *testing.T) {
	src := "Table t {\n  id int [pk]\n}\n"
	full := apiParse(src, "")
	compact := apiParse(src, `{"tokens":false}`)
	if !strings.Contains(full, `"token"`) || strings.Contains(compact, `"token"`) || !strings.Contains(compact, `"line":1`) {
		t.Errorf("full %s\ncompact %s", full, compact)
	}
	imp := apiImportSQL("CREATE TABLE t (a varchar(9) CHECK (a <> 'x\\\\y'));", "mysql", `{"tokens":false}`)
	var r ImportResponse
	if err := json.Unmarshal([]byte(imp), &r); err != nil || !r.OK || len(r.Warnings) != 1 || strings.Contains(imp, `"token"`) {
		t.Errorf("import %s", imp)
	}
}

func TestTranslateTypesDefault(t *testing.T) {
	imp := apiImportSQL("CREATE TABLE t (id int unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY, ok tinyint(1) DEFAULT 1, at datetime DEFAULT CURRENT_TIMESTAMP);", "mysql", "")
	var r ImportResponse
	if err := json.Unmarshal([]byte(imp), &r); err != nil || !r.OK || !strings.Contains(r.DBML, "database_type: 'MySQL'") {
		t.Fatalf("import %s", imp)
	}
	var e ExportResponse
	_ = json.Unmarshal([]byte(apiExportSQL(r.DBML, "postgres", "")), &e)
	if !e.OK || !e.Translated || e.SourceDialect != "mysql" || !strings.Contains(e.SQL, `"ok" boolean DEFAULT true`) || !strings.Contains(e.SQL, `"at" timestamp DEFAULT (now())`) {
		t.Errorf("translated export %+v", e)
	}
	e = ExportResponse{}
	_ = json.Unmarshal([]byte(apiExportSQL(r.DBML, "postgres", `{"translateTypes":false}`)), &e)
	if !e.OK || e.Translated || !strings.Contains(e.SQL, "tinyint(1)") {
		t.Errorf("untranslated export %+v", e)
	}
	r = ImportResponse{}
	_ = json.Unmarshal([]byte(apiImportSQL("CREATE TABLE t (a int);", "mysql", `{"omitProject":true}`)), &r)
	if strings.Contains(r.DBML, "Project") {
		t.Errorf("omitProject ignored: %s", r.DBML)
	}
}

func TestPerFormatInputCap(t *testing.T) {
	sql := "CREATE TABLE t (a int);" + strings.Repeat(" ", 2000)
	if m := decode(t, apiImportSQL(sql, "mysql", `{"maxInputBytes":{"mysql":1000}}`)); m["ok"] != false || !strings.Contains(m["error"].(string), "1000") {
		t.Errorf("per-format cap not applied: %v", m)
	}
	if m := decode(t, apiImportSQL(sql, "postgres", `{"maxInputBytes":{"mysql":1000}}`)); m["ok"] != true {
		t.Errorf("cap applied to another format: %v", m)
	}
}
