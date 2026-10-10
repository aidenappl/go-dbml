// Portions of this file are derived from holistics/dbml packages/dbml-core/src/export/MysqlExporter.js
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

package export

import (
	"strconv"
	"strings"

	"github.com/aidenappl/go-dbml/model"
)

// bq quotes a MySQL identifier (backticks doubled inside).
func bq(s string) string { return "`" + strings.ReplaceAll(s, "`", "``") + "`" }

// myLit renders a MySQL string literal (backslash and quote escaped).
func myLit(s string) string { return "'" + mysqlEscape(s) + "'" }

// mysqlSchemaPrefix is "`schema`." for non-public schemas.
func mysqlSchemaPrefix(s *cSchema) string {
	if s.print() {
		return bq(s.name) + "."
	}
	return ""
}

// exportMySQL mirrors MySQLExporter.export: schemas, tables, indexes,
// table comments, refs, then INSERTs for non-example records.
func exportMySQL(c *core) (string, error) {
	used := usedTableNames(c)
	var schemas, tables, indexes, comments, refs []string
	for _, s := range c.schemas {
		if s.print() {
			if c.translated {
				schemas = append(schemas, "CREATE SCHEMA IF NOT EXISTS "+bq(s.name)+";\n")
			} else {
				schemas = append(schemas, "CREATE SCHEMA "+bq(s.name)+";\n")
			}
		}
		for _, t := range s.tables {
			tables = append(tables, mysqlTable(t))
		}
		i := 0
		for _, t := range s.tables {
			for _, ix := range t.indexes {
				if truthy(ix.PK) {
					continue
				}
				indexes = append(indexes, mysqlIndex(t, ix, i))
				i++
			}
		}
		for _, t := range s.tables {
			if t.note != "" {
				comments = append(comments, "ALTER TABLE "+mysqlSchemaPrefix(s)+bq(t.name)+" COMMENT = "+myLit(t.note)+";\n")
			}
		}
		for _, r := range s.refs {
			refs = append(refs, mysqlRef(r, used))
		}
	}
	recs, err := mysqlRecords(c.db)
	if err != nil {
		return "", err
	}
	var all []string
	all = append(all, schemas...)
	all = append(all, tables...)
	all = append(all, indexes...)
	all = append(all, comments...)
	all = append(all, refs...)
	if len(recs) > 0 {
		all = append(all, "-- Disable foreign key checks for INSERT", "SET FOREIGN_KEY_CHECKS = 0;", "")
		all = append(all, recs...)
		all = append(all, "", "-- Re-enable foreign key checks", "SET FOREIGN_KEY_CHECKS = 1;")
	}
	return strings.Join(all, "\n"), nil
}

// mysqlTable renders CREATE TABLE (inline ENUMs, column checks,
// defaults, COMMENTs, table checks, composite primary keys).
func mysqlTable(t *cTable) string {
	var lines []string
	for _, f := range t.fields {
		var line string
		if f.enum != nil {
			var vals []string
			for _, v := range f.enum.values {
				vals = append(vals, myLit(v))
			}
			line = bq(f.name) + " ENUM (" + strings.Join(vals, ", ") + ")"
		} else {
			typ := f.typ.typeName
			if typ == "varchar" {
				typ = "varchar(255)"
			}
			line = bq(f.name) + " " + safeType(typ, MySQL, "type of column "+bq(f.name))
		}
		if g := f.generated; g != nil {
			line += " GENERATED ALWAYS AS (" + safeExpr(g.Expression, MySQL, "generated column "+bq(f.name)) + ")"
			if g.Stored {
				line += " STORED"
			} else {
				line += " VIRTUAL"
			}
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
		if f.increment && f.generated == nil {
			line += " AUTO_INCREMENT"
		}
		line += fieldCheck(f, bq, MySQL)
		line += sqlDefault(genDefault(f), MySQL, "default of column "+bq(f.name))
		if f.note != "" {
			line += " COMMENT " + myLit(f.note)
		}
		lines = append(lines, line)
	}
	lines = append(lines, tableChecks(t, bq, MySQL)...)
	lines = append(lines, compositePKs(t, bq, MySQL)...)
	return "CREATE TABLE " + mysqlSchemaPrefix(t.schema) + bq(t.name) + " (\n  " + strings.Join(lines, ",\n  ") + "\n);\n"
}

// mysqlIndex renders CREATE INDEX; unnamed indexes are numbered per
// schema (i), with the reference's quirky schema prefix inside the quotes.
func mysqlIndex(t *cTable, ix *model.Index, i int) string {
	line := "CREATE"
	if truthy(ix.Unique) {
		line += " UNIQUE"
	}
	// FULLTEXT / SPATIAL are index kinds in MySQL, not USING methods
	kind := ""
	if ix.Type != nil {
		if k := strings.ToUpper(*ix.Type); k == "FULLTEXT" || k == "SPATIAL" {
			kind = k
			line += " " + k
		}
	}
	var name string
	if ix.Name != nil && *ix.Name != "" {
		name = bq(*ix.Name)
	} else {
		// reference quirk: the schema prefix ends up inside the backticks
		name = "`" + mysqlSchemaPrefix(t.schema) + strings.ReplaceAll(t.name+"_index_"+strconv.Itoa(i), "`", "``") + "`"
	}
	line += " INDEX " + name + " ON " + mysqlSchemaPrefix(t.schema) + bq(t.name)
	line += " " + indexColumns(t, ix, bq, MySQL)
	if ix.Type != nil && *ix.Type != "" && kind == "" {
		line += " USING " + strings.ToUpper(*ix.Type)
	}
	return line + ";\n"
}

// mysqlRef renders a FOREIGN KEY, or a junction table plus two FKs for a
// many-to-many ref.
func mysqlRef(r *cRef, used map[string]bool) string {
	refEp, forEp, m2m := refSides(r)
	refT, forT := refEp.fields[0].table, forEp.fields[0].table
	refNames, forNames := fieldNames(refEp.fields, bq), fieldNames(forEp.fields, bq)
	if m2m {
		first := junctionFields1(refEp.fields)
		second := junctionFields2(forEp.fields, first)
		name := newTableName(refT.name, forT.name, used)
		var b strings.Builder
		b.WriteString("CREATE TABLE " + mysqlSchemaPrefix(refT.schema) + bq(name) + " (\n")
		for _, k := range first.keys {
			b.WriteString("  " + bq(k) + " " + first.vals[k] + ",\n")
		}
		for _, k := range second.keys {
			b.WriteString("  " + bq(k) + " " + second.vals[k] + ",\n")
		}
		b.WriteString("  PRIMARY KEY (" + joinQ(first.keys, bq) + ", " + joinQ(second.keys, bq) + ")\n")
		b.WriteString(");\n\n")
		fk := func(m *orderedMap, foreignFields string, foreignTable string, foreignSchema *cSchema) {
			b.WriteString(mysqlFinish("ALTER TABLE "+mysqlSchemaPrefix(refT.schema)+bq(name)+" ADD FOREIGN KEY ("+
				joinQ(m.keys, bq)+") REFERENCES "+mysqlTableRef(foreignSchema, foreignTable, refT.schema)+" "+foreignFields) + "\n")
		}
		fk(first, refNames, refT.name, refT.schema)
		fk(second, forNames, forT.name, forT.schema)
		return b.String()
	}
	line := "ALTER TABLE " + mysqlSchemaPrefix(forT.schema) + bq(forT.name) + " ADD "
	if r.name != nil && *r.name != "" {
		line += "CONSTRAINT " + bq(*r.name) + " "
	}
	line += "FOREIGN KEY " + forNames + " REFERENCES " + mysqlTableRef(refT.schema, refT.name, forT.schema) + " " + refNames
	if r.onDelete != nil && *r.onDelete != "" {
		line += " ON DELETE " + strings.ToUpper(*r.onDelete)
	}
	if r.onUpdate != nil && *r.onUpdate != "" {
		line += " ON UPDATE " + strings.ToUpper(*r.onUpdate)
	}
	return mysqlFinish(line)
}

// curDB marks where mysqlTableRef needs the current database's name.
const curDB = "\x00dbml-current-database\x00"

// mysqlTableRef names a table referenced from a table in schema from. A
// table in the default schema lives in the current database; referenced
// from a table in another database (schema) it must be qualified with the
// current database's name, which only the server knows (the reference
// leaves it unqualified, so MySQL looks for it in the wrong database).
func mysqlTableRef(s *cSchema, name string, from *cSchema) string {
	if !s.print() && from.print() {
		return "`" + curDB + "`." + bq(name)
	}
	return mysqlSchemaPrefix(s) + bq(name)
}

// mysqlFinish terminates stmt; a statement naming the current database
// becomes a prepared statement built with DATABASE().
func mysqlFinish(stmt string) string {
	if !strings.Contains(stmt, curDB) {
		return stmt + ";\n"
	}
	var args []string
	for i, part := range strings.Split(stmt, curDB) {
		if i > 0 {
			args = append(args, "REPLACE(DATABASE(), '`', '``')")
		}
		args = append(args, myLit(part))
	}
	return "SET @dbml_stmt = CONCAT(" + strings.Join(args, ", ") + ");\nPREPARE dbml_stmt FROM @dbml_stmt;\nEXECUTE dbml_stmt;\nDEALLOCATE PREPARE dbml_stmt;\n"
}

// mysqlRecords renders one INSERT per non-example records block.
func mysqlRecords(db *model.Database) ([]string, error) {
	var out []string
	for _, r := range nonExampleRecords(db) {
		ref := bq(r.TableName)
		if r.SchemaName != nil && *r.SchemaName != "" {
			ref = bq(*r.SchemaName) + "." + bq(r.TableName)
		}
		cols := ""
		if len(r.Columns) > 0 {
			cols = "(" + joinQ(r.Columns, bq) + ")"
		}
		var rows []string
		for _, row := range r.Values {
			var vals []string
			for _, v := range row {
				s, err := mysqlValue(v)
				if err != nil {
					return nil, err
				}
				vals = append(vals, s)
			}
			rows = append(rows, "("+strings.Join(vals, ", ")+")")
		}
		out = append(out, "INSERT INTO "+ref+" "+cols+"\nVALUES\n  "+strings.Join(rows, ",\n  ")+";")
	}
	return out, nil
}

// mysqlEscape escapes a MySQL string literal body (quotes doubled,
// backslashes escaped).
func mysqlEscape(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "'", "''"), `\`, `\\`)
}

// mysqlValue renders one record value.
func mysqlValue(v *model.RecordValue) (string, error) {
	if !v.Undefined && v.Value == nil {
		return "NULL", nil
	}
	if v.Type == "expression" {
		return safeExpr(jsString(v.Value), MySQL, "record expression"), nil
	}
	if isNumericType(v.Type) {
		if v.Undefined {
			return "undefined", nil
		}
		return safeNumber(jsString(v.Value), myLit), nil
	}
	if isBooleanType(v.Type) {
		if !v.Undefined && strings.ToUpper(jsString(v.Value)) == "TRUE" {
			return "TRUE", nil
		}
		return "FALSE", nil
	}
	s, err := jsReplaceable(v)
	if err != nil {
		return "", err
	}
	if isStringType(v.Type) || isBinaryType(v.Type) || isDateTimeType(v.Type) {
		return "'" + mysqlEscape(s) + "'", nil
	}
	return "CAST('" + mysqlEscape(s) + "' AS " + safeType(v.Type, MySQL, "record value type") + ")", nil
}
