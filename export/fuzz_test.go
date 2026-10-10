package export_test

import (
	"errors"
	"testing"

	dbml "github.com/aidenappl/go-dbml"
	"github.com/aidenappl/go-dbml/export"
)

// FuzzTranslate: SQLWith never panics, and every SQL it emits passes the
// fragment check (refusals are ErrUnsafeExpression / ErrModel only).
func FuzzTranslate(f *testing.F) {
	f.Add("Project {\n  database_type: 'MySQL'\n}\nTable t {\n  id \"int unsigned\" [pk, increment]\n  a tinyint(1) [default: 1, check: `a = 1`]\n  b datetime [default: `now()`]\n  indexes {\n    `lower(a)`\n  }\n}\n")
	f.Add("Project {\n  database_type: 'PostgreSQL'\n}\nTable t {\n  id serial [pk]\n  a \"text[]\" [default: `'{}'::text[]`]\n  b jsonb [generated: 'data->>\\'x\\'']\n}\nTable u {\n  t_id int [ref: > t.a]\n}\n")
	f.Fuzz(func(t *testing.T, src string) {
		db, _ := dbml.Parse(src)
		if db == nil {
			return
		}
		for _, from := range export.Dialects {
			for _, to := range export.Dialects {
				_, err := export.SQLWith(db, to, export.Options{TranslateTypes: true, SourceDialect: from})
				var ie *export.InternalError
				if errors.As(err, &ie) {
					t.Fatalf("%s -> %s: %v", from, to, err)
				}
			}
		}
	})
}
