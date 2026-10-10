package emit_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	dbml "github.com/aidenappl/go-dbml"
	"github.com/aidenappl/go-dbml/emit"
	"github.com/aidenappl/go-dbml/model"
)

// stripTokens returns the JSON value of db without any "token" keys.
func stripTokens(t *testing.T, db *model.Database) any {
	t.Helper()
	b, err := json.Marshal(db)
	if err != nil {
		t.Fatal(err)
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	return dropTokens(v)
}

func dropTokens(v any) any {
	switch x := v.(type) {
	case map[string]any:
		out := map[string]any{}
		for k, y := range x {
			if k == "token" {
				continue
			}
			out[k] = dropTokens(y)
		}
		return out
	case []any:
		out := make([]any, len(x))
		for i, y := range x {
			out[i] = dropTokens(y)
		}
		return out
	}
	return v
}

func corpus(t *testing.T) []string {
	t.Helper()
	var inputs []string
	for _, root := range []string{"../testdata/reference", "../testdata/holistics"} {
		_ = filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() && d.Name() == "multi" {
				return filepath.SkipDir
			}
			if !d.IsDir() && strings.HasSuffix(p, ".dbml") {
				inputs = append(inputs, p)
			}
			return nil
		})
	}
	return inputs
}

func firstDiff(path string, a, b any) string {
	switch x := a.(type) {
	case map[string]any:
		y, ok := b.(map[string]any)
		if !ok {
			return path + ": type differs"
		}
		for k := range x {
			if !reflect.DeepEqual(x[k], y[k]) {
				return firstDiff(path+"."+k, x[k], y[k])
			}
		}
		for k := range y {
			if _, ok := x[k]; !ok {
				return fmt.Sprintf("%s.%s: extra %v", path, k, y[k])
			}
		}
	case []any:
		y, ok := b.([]any)
		if !ok || len(x) != len(y) {
			return fmt.Sprintf("%s: len %d vs %v", path, len(x), b)
		}
		for i := range x {
			if !reflect.DeepEqual(x[i], y[i]) {
				return firstDiff(fmt.Sprintf("%s[%d]", path, i), x[i], y[i])
			}
		}
	}
	return fmt.Sprintf("%s: %v vs %v", path, a, b)
}

// TestRoundTrip: parse -> emit -> parse yields an equal model (ignoring
// tokens), and emit is idempotent, for every compiling input in testdata.
func TestRoundTrip(t *testing.T) {
	n, pass, idem := 0, 0, 0
	var fails []string
	for _, p := range corpus(t) {
		src, _ := os.ReadFile(p)
		db, _ := dbml.Parse(string(src))
		if db == nil {
			continue
		}
		n++
		out, err := emit.DBML(db)
		if err != nil {
			fails = append(fails, fmt.Sprintf("%s: emit: %v", p, err))
			continue
		}
		db2, diags := dbml.Parse(out)
		if db2 == nil {
			fails = append(fails, fmt.Sprintf("%s: emitted DBML does not compile: %v\n%s", p, diags, out))
			continue
		}
		a, b := stripTokens(t, db), stripTokens(t, db2)
		if !reflect.DeepEqual(a, b) {
			fails = append(fails, fmt.Sprintf("%s: model differs: %s\n%s", p, firstDiff("db", a, b), out))
			continue
		}
		pass++
		if again, _ := emit.DBML(db2); again == out {
			idem++
		} else {
			fails = append(fails, fmt.Sprintf("%s: not idempotent", p))
		}
	}
	t.Logf("round trip: %d/%d model-identical, %d/%d idempotent", pass, n, idem, n)
	for i, f := range fails {
		if i >= 10 {
			break
		}
		t.Errorf("%s", f)
	}
	if n < 500 {
		t.Fatalf("only %d inputs compiled", n)
	}
}

func TestEmitBasics(t *testing.T) {
	src := `Project shop {
  database_type: 'PostgreSQL'
  Note: 'Shop'
}

Enum status {
  active [note: 'on']
  "in progress"
}

Table users as U [headerColor: #3498DB] {
  id int [pk, increment]
  "full name" varchar(255) [not null, default: 'n\'a']
  status status [default: 'active']
  created timestamp [default: ` + "`now()`" + `]

  indexes {
    (id, "full name") [unique, name: 'ix']
    ` + "`lower(\"full name\")`" + `
  }

  Note: '''
    Users table
      indented'''
}

Table posts {
  uid int
}

Ref fk: U.id < posts.uid [delete: cascade]
`
	db, diags := dbml.Parse(src)
	if db == nil {
		t.Fatal(diags)
	}
	got, err := emit.DBML(db)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"Project shop {\n  database_type: 'PostgreSQL'\n  Note: 'Shop'\n}",
		`"in progress"`,
		"Table users as U [headerColor: #3498DB] {",
		`"full name" varchar(255) [not null, default: 'n\'a']`,
		"Note: '''\n    Users table\n      indented'''",
		"Ref fk: users.id < posts.uid [delete: cascade]",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in\n%s", want, got)
		}
	}
}

func TestIdentQuoting(t *testing.T) {
	for in, want := range map[string]string{
		"users": "users", "my table": `"my table"`, "note": `"note"`, "2fa": `"2fa"`, `a"b`: `"a\"b"`,
	} {
		if got := emit.Ident(in); got != want {
			t.Errorf("Ident(%q) = %s, want %s", in, got, want)
		}
	}
}

func TestEmitErrors(t *testing.T) {
	if _, err := emit.DBML(nil); !errors.Is(err, emit.ErrNilDatabase) {
		t.Errorf("nil db: %v", err)
	}
	db, _ := dbml.Parse("Table t {\n  a int [default: `now()`]\n}\n")
	db.Tables[0].Fields[0].DBDefault.Value = "a`b"
	if _, err := emit.DBML(db); !errors.Is(err, emit.ErrUnrepresentable) {
		t.Errorf("backtick expression: %v", err)
	}
	for _, j := range []string{`{"tables":[{"name":"x","fields":[null]}]}`, `{"tables":[null],"refs":[{"endpoints":[]}]}`,
		`{"refs":[{"endpoints":[null,null]}],"tableGroups":[{"tables":[null]}]}`} {
		var m model.Database
		if err := json.Unmarshal([]byte(j), &m); err != nil {
			continue
		}
		if _, err := emit.DBML(&m); !errors.Is(err, emit.ErrInvalidModel) {
			t.Errorf("%s: want ErrInvalidModel, got %v", j, err)
		}
	}
}
