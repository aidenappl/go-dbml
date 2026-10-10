package emit_test

import (
	"errors"
	"reflect"
	"testing"

	dbml "github.com/aidenappl/go-dbml"
	"github.com/aidenappl/go-dbml/emit"
)

// FuzzEmitRoundTrip: whenever DBML compiles, its emitted form must compile
// to the same model (tokens ignored) and emitting is idempotent.
func FuzzEmitRoundTrip(f *testing.F) {
	for _, s := range []string{
		"Table a {\n  id int [pk]\n  n \"varchar(2)\" [note: 'x\\ny']\n}\n",
		"Enum e {\n  \"a b\"\n}\nTable t {\n  e e [default: 'a b']\n}\n",
		"Table a {\n  id int\n}\nTable b {\n  a_id int\n}\nRef r: a.id < b.a_id [delete: cascade]\nDep: a -> b\n",
		"TablePartial p {\n  x int\n}\nTable t {\n  ~p\n  y int\n  Note: '''\n  multi\n    line'''\n}\n",
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, src string) {
		db, _ := dbml.Parse(src)
		if db == nil {
			return
		}
		out, err := emit.DBML(db)
		if err != nil {
			if errors.Is(err, emit.ErrUnrepresentable) {
				return
			}
			t.Fatalf("emit: %v\nsrc: %q", err, src)
		}
		db2, diags := dbml.Parse(out)
		if db2 == nil {
			t.Fatalf("emitted DBML does not compile: %v\nsrc: %q\nout: %q", diags, src, out)
		}
		a, b := stripTokens(t, db), stripTokens(t, db2)
		if !reflect.DeepEqual(a, b) {
			t.Fatalf("model differs: %s\nsrc: %q\nout: %q", firstDiff("db", a, b), src, out)
		}
		if again, _ := emit.DBML(db2); again != out {
			t.Fatalf("not idempotent\nsrc: %q\n1: %q\n2: %q", src, out, again)
		}
	})
}
