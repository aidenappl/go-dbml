// Portions of this file are derived from holistics/dbml packages/dbml-core/src/export/PostgresExporter.js
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

package export

import (
	"regexp"
	"strings"

	"github.com/aidenappl/go-dbml/model"
)

// dq quotes a PostgreSQL identifier (double quotes doubled inside); the
// SQL Server exporter uses it for index columns (a reference quirk).
func dq(s string) string { return `"` + strings.ReplaceAll(s, `"`, `""`) + `"` }

// stdLit renders a standard SQL string literal (single quotes doubled), as
// PostgreSQL (standard_conforming_strings) and SQL Server read it.
func stdLit(s string) string { return "'" + strings.ReplaceAll(s, "'", "''") + "'" }

// pgLit renders a PostgreSQL string literal. A value containing a
// backslash is written as E'...' with the backslashes doubled, so it reads
// the same whether standard_conforming_strings is on or off (a plain
// '...' would let a trailing backslash escape the closing quote when it is
// off).
func pgLit(s string) string {
	if strings.IndexByte(s, '\\') < 0 {
		return stdLit(s)
	}
	return "E'" + strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), "'", "''") + "'"
}

// pgSchemaPrefix is "\"schema\"." for non-public schemas.
func pgSchemaPrefix(s *cSchema) string {
	if s.print() {
		return dq(s.name) + "."
	}
	return ""
}

var postgresBuiltinTypes = setOf(
	"SMALLINT", "INTEGER", "INT", "BIGINT", "DECIMAL", "NUMERIC", "REAL", "DOUBLE PRECISION", "SMALLSERIAL", "SERIAL",
	"BIGSERIAL", "FLOAT", "MONEY", "CHARACTER", "CHAR", "CHARACTER VARYING", "VARCHAR", "TEXT", "NAME", "BPCHAR",
	"BYTEA", "TIMESTAMP", "TIMESTAMP WITH TIME ZONE", "TIMESTAMP WITHOUT TIME ZONE", "DATE", "TIME",
	"TIME WITH TIME ZONE", "TIME WITHOUT TIME ZONE", "INTERVAL", "BOOLEAN", "BOOL", "POINT", "LINE", "LSEG", "BOX",
	"PATH", "POLYGON", "CIRCLE", "CIDR", "INET", "MACADDR", "MACADDR8", "BIT", "BIT VARYING", "VARBIT", "UUID", "XML",
	"JSON", "JSONB", "INT4RANGE", "INT8RANGE", "NUMRANGE", "TSRANGE", "TSTZRANGE", "DATERANGE", "OID", "REGCLASS",
	"REGCOLLATION", "REGCONFIG", "REGDICTIONARY", "REGNAMESPACE", "REGOPER", "REGOPERATOR", "REGPROC", "REGPROCEDURE",
	"REGROLE", "REGTYPE", "PG_LSN", "VOID", "RECORD", "TRIGGER", "EVENT_TRIGGER", "PG_DDL_COMMAND", "UNKNOWN",
	"ANYELEMENT", "ANYARRAY", "ANYNONARRAY", "ANYENUM", "ANYRANGE", "ANY", "CSTRING", "INTERNAL", "LANGUAGE_HANDLER",
	"FDW_HANDLER", "INDEX_AM_HANDLER", "TSM_HANDLER", "OPAQUE",
)

// postgresTypeAliases are built-in aliases the reference list lacks; they
// are written unquoted too (go-dbml deviation, see pgBaseType).
var postgresTypeAliases = setOf("TIMESTAMPTZ", "TIMETZ", "INT2", "INT4", "INT8", "FLOAT4", "FLOAT8", "SERIAL2", "SERIAL4", "SERIAL8",
	"CHARACTER VARYING", "CHAR VARYING", "DEC", "DOUBLE PRECISION", "TSVECTOR", "TSQUERY", "TXID_SNAPSHOT", "PG_SNAPSHOT")

var postgresReserved = setOf("USER")

var (
	whiteSpaceRe          = regexp.MustCompile(`\s`)
	whiteSpaceOrUpperRe   = regexp.MustCompile(`[\sA-Z]`)
	hasWhiteSpace         = whiteSpaceRe.MatchString
	hasWhiteSpaceOrUpper  = whiteSpaceOrUpperRe.MatchString
	typicalIntegers       = setOf("BIGINT", "INT", "INTEGER", "SMALLINT")
	incrementIntegerTypes = setOf("SMALLSERIAL", "SERIAL", "BIGSERIAL")
)

// exportPostgres mirrors PostgresExporter.export: schemas, enum types,
// tables, indexes, COMMENT ON, refs (DEFERRABLE), then INSERTs inside a
// deferred-constraints transaction.
func exportPostgres(c *core) (string, error) {
	used := usedTableNames(c)
	var schemas, enums, tables, indexes, comments, refs []string
	enumSet := map[string]bool{}
	for _, s := range c.schemas {
		if s.print() {
			schemas = append(schemas, "CREATE SCHEMA "+dq(s.name)+";\n")
		}
		for _, e := range s.enums {
			name := pgSchemaPrefix(s) + dq(e.name)
			var vals []string
			for _, v := range e.values {
				vals = append(vals, "  "+pgLit(v))
			}
			enums = append(enums, "CREATE TYPE "+name+" AS ENUM (\n"+strings.Join(vals, ",\n")+"\n);\n")
			enumSet[name] = true
		}
	}
	for _, s := range c.schemas {
		for _, t := range s.tables {
			tables = append(tables, pgTable(t, enumSet))
		}
		for _, t := range s.tables {
			for _, ix := range t.indexes {
				if !truthy(ix.PK) {
					indexes = append(indexes, pgIndex(t, ix))
				}
			}
		}
		for _, t := range s.tables {
			if t.note != "" {
				comments = append(comments, "COMMENT ON TABLE "+pgSchemaPrefix(s)+dq(t.name)+" IS "+pgLit(t.note)+";\n")
			}
			for _, f := range t.fields {
				if f.note != "" {
					comments = append(comments, "COMMENT ON COLUMN "+pgSchemaPrefix(s)+dq(t.name)+"."+dq(f.name)+" IS "+pgLit(f.note)+";\n")
				}
			}
		}
		for _, r := range s.refs {
			refs = append(refs, pgRef(r, used))
		}
	}
	recs := pgRecords(c.db)
	var all []string
	all = append(all, schemas...)
	all = append(all, enums...)
	all = append(all, tables...)
	all = append(all, indexes...)
	all = append(all, comments...)
	all = append(all, refs...)
	if len(recs) > 0 {
		all = append(all, "-- Defer constraint checking for INSERT", "BEGIN;", "SET CONSTRAINTS ALL DEFERRED;", "")
		all = append(all, recs...)
		all = append(all, "", "SET CONSTRAINTS ALL IMMEDIATE;", "COMMIT;")
	}
	return strings.Join(all, "\n"), nil
}

// pgFieldType renders a column type: identity/serial for increment,
// quoting user-defined and schema-qualified types like the reference.
func pgFieldType(f *cField, enumSet map[string]bool) string {
	t := f.typ
	if f.increment {
		raw := strings.ToUpper(t.typeName)
		switch {
		case typicalIntegers[raw]:
			return raw + " GENERATED BY DEFAULT AS IDENTITY"
		case incrementIntegerTypes[raw]:
			return raw
		}
		return "SERIAL"
	}
	if t.schemaName == nil || *t.schemaName == "" || *t.schemaName == defaultSchema {
		orig := t.typeName
		up := strings.ToUpper(orig)
		if !postgresBuiltinTypes[up] && !pgUnquoted(up) && (hasWhiteSpaceOrUpper(orig) || postgresReserved[up]) {
			return dq(orig)
		}
		return orig
	}
	if t.originalTypeName != nil {
		return dq(*t.schemaName) + "." + dq(*t.originalTypeName)
	}
	schema := *t.schemaName + "."
	if hasWhiteSpaceOrUpper(*t.schemaName) {
		schema = dq(*t.schemaName) + "."
	}
	name := t.typeName
	if hasWhiteSpaceOrUpper(name) {
		name = dq(name)
	}
	withSchema := schema + name
	asEnum := dq(*t.schemaName) + "." + dq(t.typeName)
	if !enumSet[asEnum] && !hasWhiteSpace(asEnum) {
		withSchema = strings.ReplaceAll(withSchema, `"`, "")
	}
	return withSchema
}

// pgTable renders CREATE TABLE.
func pgTable(t *cTable, enumSet map[string]bool) string {
	var lines []string
	for _, f := range t.fields {
		line := dq(f.name) + " " + safeType(pgFieldType(f, enumSet), Postgres, "type of column "+dq(f.name))
		if g := f.generated; g != nil {
			// PostgreSQL before 18 only has STORED generated columns.
			line += " GENERATED ALWAYS AS (" + safeExpr(g.Expression, Postgres, "generated column "+dq(f.name)) + ") STORED"
		}
		if f.unique {
			line += " UNIQUE"
		}
		if f.pk {
			line += " PRIMARY KEY"
		}
		if f.notNull {
			line += " NOT NULL"
		}
		line += fieldCheck(f, dq, Postgres)
		line += sqlDefault(genDefault(f), Postgres, "default of column "+dq(f.name))
		lines = append(lines, line)
	}
	lines = append(lines, tableChecks(t, dq, Postgres)...)
	lines = append(lines, compositePKs(t, dq, Postgres)...)
	return "CREATE TABLE " + pgSchemaPrefix(t.schema) + dq(t.name) + " (\n  " + strings.Join(lines, ",\n  ") + "\n);\n"
}

// pgIndex renders CREATE INDEX (USING type before the column list).
func pgIndex(t *cTable, ix *model.Index) string {
	line := "CREATE"
	if truthy(ix.Unique) {
		line += " UNIQUE"
	}
	line += " INDEX"
	if ix.Name != nil && *ix.Name != "" {
		line += " " + dq(*ix.Name)
	}
	line += " ON " + pgSchemaPrefix(t.schema) + dq(t.name)
	if ix.Type != nil && *ix.Type != "" {
		line += " USING " + strings.ToUpper(*ix.Type)
	}
	return line + " " + indexColumns(t, ix, dq, Postgres) + ";\n"
}

// pgRef renders a FOREIGN KEY, or a junction table plus two FKs for a
// many-to-many ref.
func pgRef(r *cRef, used map[string]bool) string {
	refEp, forEp, m2m := refSides(r)
	refT, forT := refEp.fields[0].table, forEp.fields[0].table
	refNames, forNames := fieldNames(refEp.fields, dq), fieldNames(forEp.fields, dq)
	if m2m {
		first := junctionFields1(refEp.fields)
		second := junctionFields2(forEp.fields, first)
		name := newTableName(refT.name, forT.name, used)
		var b strings.Builder
		b.WriteString("CREATE TABLE " + pgSchemaPrefix(refT.schema) + dq(name) + " (\n")
		for _, k := range first.keys {
			b.WriteString("  " + dq(k) + " " + first.vals[k] + ",\n")
		}
		for _, k := range second.keys {
			b.WriteString("  " + dq(k) + " " + second.vals[k] + ",\n")
		}
		b.WriteString("  PRIMARY KEY (" + joinQ(first.keys, dq) + ", " + joinQ(second.keys, dq) + ")\n")
		b.WriteString(");\n\n")
		fk := func(m *orderedMap, foreignFields string, foreignTable string, foreignSchema *cSchema) {
			b.WriteString("ALTER TABLE " + pgSchemaPrefix(refT.schema) + dq(name) + " ADD FOREIGN KEY (" +
				joinQ(m.keys, dq) + ") REFERENCES " + pgSchemaPrefix(foreignSchema) + dq(foreignTable) + " " + foreignFields +
				" DEFERRABLE INITIALLY IMMEDIATE;\n\n")
		}
		fk(first, refNames, refT.name, refT.schema)
		fk(second, forNames, forT.name, forT.schema)
		return b.String()
	}
	line := "ALTER TABLE " + pgSchemaPrefix(forT.schema) + dq(forT.name) + " ADD "
	if r.name != nil && *r.name != "" {
		line += "CONSTRAINT " + dq(*r.name) + " "
	}
	line += "FOREIGN KEY " + forNames + " REFERENCES " + pgSchemaPrefix(refT.schema) + dq(refT.name) + " " + refNames
	if r.onDelete != nil && *r.onDelete != "" {
		line += " ON DELETE " + strings.ToUpper(*r.onDelete)
	}
	if r.onUpdate != nil && *r.onUpdate != "" {
		line += " ON UPDATE " + strings.ToUpper(*r.onUpdate)
	}
	return line + " DEFERRABLE INITIALLY IMMEDIATE;\n"
}

// pgRecords renders one INSERT per non-empty, non-example records block.
func pgRecords(db *model.Database) []string {
	var out []string
	for _, r := range nonExampleRecords(db) {
		if len(r.Values) == 0 {
			continue
		}
		ref := dq(r.TableName)
		if r.SchemaName != nil && *r.SchemaName != "" {
			ref = dq(*r.SchemaName) + "." + dq(r.TableName)
		}
		cols := ""
		if len(r.Columns) > 0 {
			var qs []string
			for _, c := range r.Columns {
				qs = append(qs, dq(c))
			}
			cols = "(" + strings.Join(qs, ", ") + ")"
		}
		var rows []string
		for _, row := range r.Values {
			var vals []string
			for _, v := range row {
				vals = append(vals, pgValue(v))
			}
			rows = append(rows, "("+strings.Join(vals, ", ")+")")
		}
		out = append(out, "INSERT INTO "+ref+" "+cols+"\nVALUES\n  "+strings.Join(rows, ",\n  ")+";")
	}
	return out
}

// pgValue renders one record value.
func pgValue(v *model.RecordValue) string {
	if !v.Undefined && v.Value == nil {
		return "NULL"
	}
	s := jsString(v.Value)
	if v.Undefined {
		s = "undefined"
	}
	if v.Type == "expression" {
		return safeExpr(s, Postgres, "record expression")
	}
	if isNumericType(v.Type) {
		if v.Undefined {
			return s
		}
		return safeNumber(s, stdLit)
	}
	if isBooleanType(v.Type) {
		if strings.ToUpper(s) == "TRUE" {
			return "TRUE"
		}
		return "FALSE"
	}
	if isStringType(v.Type) || isDateTimeType(v.Type) || isBinaryType(v.Type) {
		return pgLit(s)
	}
	return "CAST(" + pgLit(s) + " AS " + safeType(v.Type, Postgres, "record value type") + ")"
}

var (
	pgTypeArgs  = regexp.MustCompile(`\s*\([^()]*\)`)
	pgTypeArray = regexp.MustCompile(`(\s*\[[0-9]*\])+$`)
	pgSpaces    = regexp.MustCompile(`\s+`)
)

// pgBaseType reduces an upper-cased type such as
// "TIMESTAMP(6) WITH TIME ZONE" or "CHARACTER VARYING(255)[]" to its
// built-in name ("TIMESTAMP WITH TIME ZONE"), so multi-word and
// parameterised built-ins are not quoted as one identifier (the reference
// quotes them, which yields invalid SQL; a documented deviation).
func pgBaseType(up string) string {
	t := pgTypeArray.ReplaceAllString(up, "")
	t = pgTypeArgs.ReplaceAllString(t, "")
	return strings.TrimSpace(pgSpaces.ReplaceAllString(t, " "))
}

// pgUnquoted reports whether an upper-cased type is a (parameterised /
// array / aliased) built-in that go-dbml writes unquoted.
func pgUnquoted(up string) bool {
	b := pgBaseType(up)
	return postgresBuiltinTypes[b] || postgresTypeAliases[b] || pgInterval.MatchString(b)
}

// pgInterval matches INTERVAL with its optional fields (DAY TO SECOND, ...).
var pgInterval = regexp.MustCompile(`^INTERVAL( (YEAR|MONTH|DAY|HOUR|MINUTE|SECOND)( TO (MONTH|HOUR|MINUTE|SECOND))?)?$`)
