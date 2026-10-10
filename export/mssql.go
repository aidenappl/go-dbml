// Portions of this file are derived from holistics/dbml packages/dbml-core/src/export/SqlServerExporter.js
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

package export

import (
	"strconv"
	"strings"

	"github.com/aidenappl/go-dbml/model"
)

// sb quotes a SQL Server identifier (closing brackets doubled inside).
func sb(s string) string { return "[" + strings.ReplaceAll(s, "]", "]]") + "]" }

// msSchemaPrefix is "[schema]." for non-public schemas.
func msSchemaPrefix(s *cSchema) string {
	if s.print() {
		return sb(s.name) + "."
	}
	return ""
}

// exportMSSQL mirrors SqlServerExporter.export (GO-separated batches,
// descriptions via sp_addextendedproperty, NOCHECK around INSERTs).
func exportMSSQL(c *core) (string, error) {
	used := usedTableNames(c)
	var schemas, tables, indexes, comments, refs []string
	for _, s := range c.schemas {
		if s.print() {
			schemas = append(schemas, "CREATE SCHEMA "+sb(s.name)+"\nGO\n")
		}
		for _, t := range s.tables {
			tables = append(tables, msTable(t))
		}
		i := 0
		for _, t := range s.tables {
			for _, ix := range t.indexes {
				if truthy(ix.PK) {
					continue
				}
				indexes = append(indexes, msIndex(t, ix, i))
				i++
			}
		}
		for _, t := range s.tables {
			level0 := "dbo"
			if s.print() {
				level0 = s.name
			}
			if t.note != "" {
				comments = append(comments, "EXEC sp_addextendedproperty\n"+
					"@name = N'Table_Description',\n"+
					"@value = "+nLit(t.note)+",\n"+
					"@level0type = N'Schema', @level0name = "+nLit(level0)+",\n"+
					"@level1type = N'Table',  @level1name = "+nLit(t.name)+";\nGO\n")
			}
			for _, f := range t.fields {
				if f.note == "" {
					continue
				}
				comments = append(comments, "EXEC sp_addextendedproperty\n"+
					"@name = N'Column_Description',\n"+
					"@value = "+nLit(f.note)+",\n"+
					"@level0type = N'Schema', @level0name = "+nLit(level0)+",\n"+
					"@level1type = N'Table',  @level1name = "+nLit(t.name)+",\n"+
					"@level2type = N'Column', @level2name = "+nLit(f.name)+";\nGO\n")
			}
		}
		for _, r := range s.refs {
			refs = append(refs, msRef(r, used))
		}
	}
	recs, err := msRecords(c.db)
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
		// The reference lists every records block here, [example] ones included.
		var uniq []string
		seen := map[string]bool{}
		for _, r := range c.db.Records {
			ref := sb(r.TableName)
			if r.SchemaName != nil && *r.SchemaName != "" {
				ref = sb(*r.SchemaName) + "." + sb(r.TableName)
			}
			if !seen[ref] {
				seen[ref] = true
				uniq = append(uniq, ref)
			}
		}
		all = append(all, "-- Disable constraint checks for tables with data")
		for _, u := range uniq {
			all = append(all, "ALTER TABLE "+u+" NOCHECK CONSTRAINT ALL;\nGO")
		}
		all = append(all, "")
		all = append(all, recs...)
		all = append(all, "", "-- Re-enable constraint checks")
		for _, u := range uniq {
			all = append(all, "ALTER TABLE "+u+" WITH CHECK CHECK CONSTRAINT ALL;\nGO")
		}
	}
	return strings.Join(all, "\n"), nil
}

// msTable renders CREATE TABLE (enums become CHECK ... IN constraints).
func msTable(t *cTable) string {
	var lines []string
	for _, f := range t.fields {
		var line string
		if f.enum != nil {
			var vals []string
			for _, v := range f.enum.values {
				vals = append(vals, stdLit(v))
			}
			line = sb(f.name) + " nvarchar(255) CHECK (" + sb(f.name) + " IN (" + strings.Join(vals, ", ") + "))"
		} else {
			typ := f.typ.typeName
			if typ == "varchar" {
				typ = "nvarchar(255)"
			}
			line = sb(f.name) + " " + safeType(typ, MSSQL, "type of column "+sb(f.name))
		}
		if g := f.generated; g != nil {
			// a computed column has no type in T-SQL
			line = sb(f.name) + " AS (" + safeExpr(g.Expression, MSSQL, "computed column "+sb(f.name)) + ")"
			if g.Stored {
				line += " PERSISTED"
			}
		}
		if f.unique {
			line += " UNIQUE"
		}
		if f.pk {
			line += " PRIMARY KEY"
		}
		if f.notNull && (f.generated == nil || f.generated.Stored) {
			line += " NOT NULL"
		}
		if f.increment && f.generated == nil {
			line += " IDENTITY(1, 1)"
		}
		line += fieldCheck(f, sb, MSSQL)
		line += sqlDefault(genDefault(f), MSSQL, "default of column "+sb(f.name))
		lines = append(lines, line)
	}
	lines = append(lines, tableChecks(t, sb, MSSQL)...)
	lines = append(lines, compositePKs(t, sb, MSSQL)...)
	return "CREATE TABLE " + msSchemaPrefix(t.schema) + sb(t.name) + " (\n  " + strings.Join(lines, ",\n  ") + "\n)\nGO\n"
}

// msIndex renders CREATE INDEX (index columns double-quoted, as in the
// reference).
func msIndex(t *cTable, ix *model.Index, i int) string {
	line := "CREATE"
	if truthy(ix.Unique) {
		line += " UNIQUE"
	}
	name := sb(t.name + "_index_" + strconv.Itoa(i))
	if ix.Name != nil && *ix.Name != "" {
		name = sb(*ix.Name)
	}
	line += " INDEX " + name + " ON " + msSchemaPrefix(t.schema) + sb(t.name)
	// [x], not the reference's "x": a "..." name is a string literal
	// under QUOTED_IDENTIFIER OFF (sqlcmd's default)
	return line + " " + indexColumns(t, ix, sb, MSSQL) + "\nGO\n"
}

// msRef renders a FOREIGN KEY, or a junction table plus two FKs for a
// many-to-many ref.
func msRef(r *cRef, used map[string]bool) string {
	refEp, forEp, m2m := refSides(r)
	refT, forT := refEp.fields[0].table, forEp.fields[0].table
	refNames, forNames := fieldNames(refEp.fields, sb), fieldNames(forEp.fields, sb)
	if m2m {
		first := junctionFields1(refEp.fields)
		second := junctionFields2(forEp.fields, first)
		name := newTableName(refT.name, forT.name, used)
		var b strings.Builder
		b.WriteString("CREATE TABLE " + msSchemaPrefix(refT.schema) + sb(name) + " (\n")
		for _, k := range first.keys {
			b.WriteString("  " + sb(k) + " " + first.vals[k] + ",\n")
		}
		for _, k := range second.keys {
			b.WriteString("  " + sb(k) + " " + second.vals[k] + ",\n")
		}
		b.WriteString("  PRIMARY KEY (" + joinQ(first.keys, sb) + ", " + joinQ(second.keys, sb) + ")\n")
		b.WriteString(");\nGO\n\n")
		fk := func(m *orderedMap, foreignFields string, foreignTable string, foreignSchema *cSchema) {
			b.WriteString("ALTER TABLE " + msSchemaPrefix(refT.schema) + sb(name) + " ADD FOREIGN KEY (" +
				joinQ(m.keys, sb) + ") REFERENCES " + msSchemaPrefix(foreignSchema) + sb(foreignTable) + " " + foreignFields + ";\nGO\n\n")
		}
		fk(first, refNames, refT.name, refT.schema)
		fk(second, forNames, forT.name, forT.schema)
		return b.String()
	}
	line := "ALTER TABLE " + msSchemaPrefix(forT.schema) + sb(forT.name) + " ADD "
	if r.name != nil && *r.name != "" {
		line += "CONSTRAINT " + sb(*r.name) + " "
	}
	line += "FOREIGN KEY " + forNames + " REFERENCES " + msSchemaPrefix(refT.schema) + sb(refT.name) + " " + refNames
	if r.onDelete != nil && *r.onDelete != "" {
		line += " ON DELETE " + strings.ToUpper(*r.onDelete)
	}
	if r.onUpdate != nil && *r.onUpdate != "" {
		line += " ON UPDATE " + strings.ToUpper(*r.onUpdate)
	}
	return line + "\nGO\n"
}

// msRecords renders one INSERT per non-example records block.
func msRecords(db *model.Database) ([]string, error) {
	var out []string
	for _, r := range nonExampleRecords(db) {
		ref := sb(r.TableName)
		if r.SchemaName != nil && *r.SchemaName != "" {
			ref = sb(*r.SchemaName) + "." + sb(r.TableName)
		}
		cols := ""
		if len(r.Columns) > 0 {
			cols = "(" + joinQ(r.Columns, sb) + ")"
		}
		var rows []string
		for _, row := range r.Values {
			var vals []string
			for _, v := range row {
				s, err := msValue(v)
				if err != nil {
					return nil, err
				}
				vals = append(vals, s)
			}
			rows = append(rows, "("+strings.Join(vals, ", ")+")")
		}
		out = append(out, "INSERT INTO "+ref+" "+cols+"\nVALUES\n  "+strings.Join(rows, ",\n  ")+";\nGO")
	}
	return out, nil
}

// msValue renders one record value (booleans as 1/0, binary as 0x..).
func msValue(v *model.RecordValue) (string, error) {
	if !v.Undefined && v.Value == nil {
		return "NULL", nil
	}
	if v.Type == "expression" {
		return safeExpr(jsString(v.Value), MSSQL, "record expression"), nil
	}
	if isNumericType(v.Type) {
		if v.Undefined {
			return "undefined", nil
		}
		return safeNumber(jsString(v.Value), stdLit), nil
	}
	if isBooleanType(v.Type) {
		if !v.Undefined && strings.ToUpper(jsString(v.Value)) == "TRUE" {
			return "1", nil
		}
		return "0", nil
	}
	if isStringType(v.Type) || isDateTimeType(v.Type) {
		s, err := jsReplaceable(v)
		if err != nil {
			return "", err
		}
		return "'" + strings.ReplaceAll(s, "'", "''") + "'", nil
	}
	if isBinaryType(v.Type) {
		if v.Undefined {
			return "0xundefined", nil
		}
		if h := jsString(v.Value); hexDigits.MatchString(h) {
			return "0x" + h, nil
		}
		// not hex digits: a string literal instead of raw text (safety deviation)
		return stdLit(jsString(v.Value)), nil
	}
	s, err := jsReplaceable(v)
	if err != nil {
		return "", err
	}
	return "CAST('" + strings.ReplaceAll(s, "'", "''") + "' AS " + safeType(v.Type, MSSQL, "record value type") + ")", nil
}

// nLit renders a Unicode (N'...') string literal.
func nLit(s string) string { return "N" + stdLit(s) }
