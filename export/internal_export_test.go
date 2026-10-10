package export

import "strings"

// Test-only hooks for package export_test.

// PgUnquotedBuiltin reports whether go-dbml writes the PostgreSQL type t
// unquoted although @dbml/core quotes it (a parameterised or multi-word
// built-in such as "VARCHAR(20)" or "character varying(255)").
func PgUnquotedBuiltin(t string) bool {
	return pgUnquoted(strings.ToUpper(t))
}
