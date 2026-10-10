package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestStripTokens(t *testing.T) {
	src := "Enum e {\n  a\n}\nTable t {\n  id int [pk]\n  e e\n  indexes {\n    id\n  }\n}\nTable u {\n  t_id int [ref: > t.id, note: 'n']\n}\nNote n {\n  'x'\n}\nTableGroup g {\n  t\n}\n"
	r := CompileWithOptions(src, "/main.dbml", Options{})
	if r.Database == nil {
		t.Fatal(r.Diagnostics)
	}
	full, _ := json.Marshal(r.Database)
	StripTokens(r.Database)
	StripTokens(r.Database) // idempotent
	compact, _ := json.Marshal(r.Database)
	if strings.Contains(string(compact), `"token"`) {
		t.Errorf("token left in %s", compact)
	}
	if len(compact)*2 > len(full) {
		t.Errorf("compact %d bytes vs %d", len(compact), len(full))
	}
	db := r.Database
	if db.Tables[0].Line != 4 || db.Tables[0].Fields[1].Line != 6 || db.Tables[1].Line != 11 || db.Refs[0].Line != 12 || db.Enums[0].Line != 1 || db.Notes[0].Line != 14 || db.TableGroups[0].Line != 17 {
		t.Errorf("lines: table %d col %d table %d ref %d enum %d note %d group %d", db.Tables[0].Line, db.Tables[0].Fields[1].Line, db.Tables[1].Line, db.Refs[0].Line, db.Enums[0].Line, db.Notes[0].Line, db.TableGroups[0].Line)
	}
	// Line is only set by StripTokens
	if r2 := CompileWithOptions(src, "/main.dbml", Options{}); r2.Database.Tables[0].Line != 0 {
		t.Error("Line set without StripTokens")
	}
}
