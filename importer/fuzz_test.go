package importer

import (
	"errors"
	"testing"

	"github.com/aidenappl/go-dbml/export"
)

// fuzzImport runs an importer without the public entry points' recover (a
// panic fails the fuzz run) in parity and clean mode, with records, and
// checks the rendered output is produced without an internal error.
func fuzzImport(t *testing.T, src string, imp func(string, *config) (*rawDB, error)) {
	for _, clean := range []bool{false, true} {
		raw, err := imp(src, &config{clean: clean, records: true, maxRows: 100})
		if err != nil {
			continue
		}
		db := raw.toModel()
		if clean {
			cleanModel(db)
		}
		if _, err := export.CoreDBML(db, export.CoreDBMLOptions{OverrideCardinality: true, Escape: clean}); err != nil {
			var ie *export.InternalError
			if !errors.As(err, &ie) {
				continue
			}
			t.Fatalf("export (clean=%v): %v", clean, err)
		}
	}
}

func seeds(f *testing.F, extra ...string) {
	for _, s := range append([]string{
		"CREATE TABLE a (id int PRIMARY KEY, b varchar(10) NOT NULL DEFAULT 'x', FOREIGN KEY (id) REFERENCES c(d));",
		"ALTER TABLE a ADD CONSTRAINT x FOREIGN KEY (b) REFERENCES a (id) ON DELETE CASCADE;",
		"INSERT INTO a (id, b) VALUES (1, 'x'), (2, NULL);",
		"CREATE INDEX i ON a (b);",
		"create_table \"users\" do |t|\n  t.string \"name\"\nend\n",
	}, extra...) {
		f.Add(s)
	}
}

func FuzzImportMySQL(f *testing.F) {
	seeds(f, "CREATE TABLE t (e ENUM('a','b'), s SET('x'), c INT CHECK (c > 0)) COMMENT='n';", "DELIMITER //\nCREATE PROCEDURE p() BEGIN SELECT 1; END //")
	f.Fuzz(func(t *testing.T, s string) { fuzzImport(t, s, importMySQL) })
}

func FuzzImportPostgres(f *testing.F) {
	seeds(f, "CREATE TYPE e AS ENUM ('a'); COMMENT ON TABLE a IS 'x'; CREATE TABLE b (x e[] DEFAULT '{}'::e[]);", "$$ $tag$ x $tag$")
	f.Fuzz(func(t *testing.T, s string) { fuzzImport(t, s, importPostgres) })
}

func FuzzImportMSSQL(f *testing.F) {
	seeds(f, "CREATE TABLE [a] ([id] int IDENTITY(1,1) PRIMARY KEY, [n] nvarchar(max))\nGO\nEXEC sp_addextendedproperty @name = N'Table_Description', @value = 'x', @level0type = N'Schema', @level0name = 'dbo', @level1type = N'Table', @level1name = 'a';")
	f.Fuzz(func(t *testing.T, s string) { fuzzImport(t, s, importMSSQL) })
}

func FuzzImportSchemaRB(f *testing.F) {
	seeds(f, "add_foreign_key \"a\", \"b\", column: \"x\", on_delete: :cascade\n", "add_check_constraint :a, \"x > 0\", name: \"c\"\n")
	f.Fuzz(func(t *testing.T, s string) { fuzzImport(t, s, importSchemaRB) })
}
