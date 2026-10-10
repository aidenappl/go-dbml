// Portions of this file are derived from holistics/dbml packages/dbml-core/src/export/utils.ts and
// packages/dbml-parse/src/core/global_modules/records/utils/data/sqlTypes.ts
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

// Package export renders a model.Database as SQL DDL (plus INSERTs for
// records) for MySQL/MariaDB, PostgreSQL and SQL Server.
//
// The output is byte-for-byte what @dbml/core 10.3.1's
// exporter.export(dbml, dialect) produces, quirks included (see AGENTS.md
// §6 for the measured match rate and the known differences): it is a port
// of MysqlExporter.js, PostgresExporter.js and SqlServerExporter.js over a
// port of the @dbml/core model (partial injection, enum binding, schema
// ordering, many-to-many junction tables).
package export

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/aidenappl/go-dbml/model"
)

// Dialect is a SQL dialect.
type Dialect string

// Supported dialects.
const (
	MySQL    Dialect = "mysql" // also MariaDB
	Postgres Dialect = "postgres"
	MSSQL    Dialect = "mssql"
)

// Dialects lists the supported dialects.
var Dialects = []Dialect{MySQL, Postgres, MSSQL}

// ParseDialect maps a name (mysql, mariadb, postgres, postgresql, pg,
// mssql, sqlserver) to a Dialect.
func ParseDialect(s string) (Dialect, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "mysql", "mariadb":
		return MySQL, nil
	case "postgres", "postgresql", "pg":
		return Postgres, nil
	case "mssql", "sqlserver", "sql server", "tsql":
		return MSSQL, nil
	}
	return "", fmt.Errorf("unsupported SQL dialect %q", s)
}

// InternalError reports a recovered panic inside an exporter (a bug).
type InternalError struct{ Value any }

func (e *InternalError) Error() string { return fmt.Sprintf("export: internal error: %v", e.Value) }

// SQL exports db in the given dialect. It returns an error (wrapping
// ErrModel) where @dbml/core would throw while building its model, e.g.
// when a TablePartial injects a ref whose two endpoints are the same.
func SQL(db *model.Database, d Dialect) (out string, err error) {
	if db == nil {
		return "", fmt.Errorf("export: nil database")
	}
	defer func() {
		if r := recover(); r != nil {
			if ue, ok := r.(*UnsafeExpressionError); ok {
				out, err = "", ue
				return
			}
			out, err = "", &InternalError{Value: r}
		}
	}()
	c, err := buildCore(db)
	if err != nil {
		return "", err
	}
	switch d {
	case MySQL:
		return exportMySQL(c)
	case Postgres:
		return exportPostgres(c)
	case MSSQL:
		return exportMSSQL(c)
	}
	return "", fmt.Errorf("unsupported SQL dialect %q", d)
}

// ---------------------------------------------------------------- helpers

// jsNumber formats a float64 like JavaScript String(number).
func jsNumber(f float64) string {
	if math.IsNaN(f) {
		return "NaN"
	}
	if math.IsInf(f, 0) {
		if f > 0 {
			return "Infinity"
		}
		return "-Infinity"
	}
	if f == 0 {
		return "0"
	}
	abs := math.Abs(f)
	if abs >= 1e-6 && abs < 1e21 {
		return strconv.FormatFloat(f, 'f', -1, 64)
	}
	s := strconv.FormatFloat(f, 'e', -1, 64)
	mant, exp, _ := strings.Cut(s, "e")
	sign := exp[:1]
	e := strings.TrimLeft(exp[1:], "0")
	if e == "" {
		e = "0"
	}
	return mant + "e" + sign + e
}

// jsString renders a JSON-ish model value like JavaScript String(v).
func jsString(v any) string {
	switch x := v.(type) {
	case nil:
		return "null"
	case string:
		return x
	case float64:
		return jsNumber(x)
	case bool:
		if x {
			return "true"
		}
		return "false"
	}
	return fmt.Sprint(v)
}

func setOf(xs ...string) map[string]bool {
	m := map[string]bool{}
	for _, x := range xs {
		m[x] = true
	}
	return m
}

// Union over every dialect of @dbml/parse's sqlTypes.ts sets (the
// exporters call the helpers without a dialect).
var (
	integerTypes = setOf("int", "integer", "smallint", "bigint", "tinyint", "mediumint", "serial", "bigserial", "smallserial")
	floatTypes   = setOf("decimal", "numeric", "float", "double", "real", "double precision", "number")
	boolTypes    = setOf("bool", "boolean", "bit")
	stringTypes  = setOf("varchar", "char", "text", "tinytext", "mediumtext", "longtext", "string", "character",
		"character varying", "nvarchar", "nchar", "ntext", "varchar2", "nvarchar2")
	binaryTypes   = setOf("binary", "varbinary", "blob", "tinyblob", "mediumblob", "longblob", "bytea", "raw")
	datetimeTypes = setOf("date", "datetime", "timestamp", "time", "timestamptz", "timestamp with time zone",
		"timestamp without time zone", "timetz", "time with time zone", "time without time zone", "datetime2",
		"smalldatetime", "timestamp with local time zone")
)

var spaces = regexp.MustCompile(`\s+`)

func normType(t string) string {
	return spaces.ReplaceAllString(strings.TrimSpace(strings.ToLower(t)), " ")
}

func isNumericType(t string) bool {
	n := normType(t)
	// MySQL "int unsigned" / "decimal zerofill" (go-dbml clean imports)
	for _, sfx := range []string{" zerofill", " unsigned", " signed"} {
		n = strings.TrimSuffix(n, sfx)
	}
	return integerTypes[n] || floatTypes[n]
}
func isBooleanType(t string) bool  { return boolTypes[normType(t)] }
func isStringType(t string) bool   { return stringTypes[normType(t)] }
func isBinaryType(t string) bool   { return binaryTypes[normType(t)] }
func isDateTimeType(t string) bool { return datetimeTypes[normType(t)] }

// errJSCrash marks inputs on which the reference exporter throws a
// TypeError (calling .replace on a non-string record value).
func errJSCrash(what string) error {
	return fail("Cannot read properties of " + what + " (reading 'replace')")
}

// jsReplaceable returns v as a string when JavaScript `v.replace(...)`
// would work, mirroring the reference's crash otherwise.
func jsReplaceable(v *model.RecordValue) (string, error) {
	if v.Undefined {
		return "", errJSCrash("undefined")
	}
	s, ok := v.Value.(string)
	if !ok {
		return "", fail("val.value.replace is not a function")
	}
	return s, nil
}

// orderedMap is a JavaScript Map<string,string> (insertion-ordered keys).
type orderedMap struct {
	keys []string
	vals map[string]string
}

func newOrderedMap() *orderedMap { return &orderedMap{vals: map[string]string{}} }

func (m *orderedMap) set(k, v string) {
	if _, ok := m.vals[k]; !ok {
		m.keys = append(m.keys, k)
	}
	m.vals[k] = v
}

func (m *orderedMap) has(k string) bool { _, ok := m.vals[k]; return ok }

// junctionFields1 mirrors buildJunctionFields1.
func junctionFields1(fs []*cField) *orderedMap {
	m := newOrderedMap()
	for _, f := range fs {
		m.set(f.table.name+"_"+f.name, f.typ.typeName)
	}
	return m
}

// junctionFields2 mirrors buildJunctionFields2.
func junctionFields2(fs []*cField, first *orderedMap) *orderedMap {
	m := newOrderedMap()
	for _, f := range fs {
		name := f.table.name + "_" + f.name
		for n := 1; first.has(name); n++ {
			name = f.table.name + "_" + f.name + "(" + strconv.Itoa(n) + ")"
		}
		m.set(name, f.typ.typeName)
	}
	return m
}

// newTableName mirrors buildNewTableName (mutates used).
func newTableName(a, b string, used map[string]bool) string {
	name := a + "_" + b
	for n := 1; used[name]; n++ {
		name = a + "_" + b + "(" + strconv.Itoa(n) + ")"
	}
	used[name] = true
	return name
}

func usedTableNames(c *core) map[string]bool {
	used := map[string]bool{}
	for _, s := range c.schemas {
		for _, t := range s.tables {
			used[t.name] = true
		}
	}
	return used
}

// refSides mirrors the ref/foreign endpoint selection of exportRefs.
func refSides(r *cRef) (ref, foreign *cEndpoint, manyToMany bool) {
	one := -1
	for i, ep := range r.endpoints {
		if maxIsOne(ep.relation) {
			one = i
			break
		}
	}
	idx := one
	if one == -1 {
		idx = 0
	}
	return r.endpoints[idx], r.endpoints[1-idx], one == -1
}

func isNullDefault(d *model.Default) bool {
	if d.Type != "boolean" {
		return false
	}
	if d.Value == nil {
		return true
	}
	s, ok := d.Value.(string)
	return ok && strings.ToLower(s) == "null"
}

func nonExampleRecords(db *model.Database) []*model.TableRecord {
	var out []*model.TableRecord
	for _, r := range db.Records {
		if !truthy(r.Example) {
			out = append(out, r)
		}
	}
	return out
}

func indexColumns(t *cTable, ix *model.Index, quote func(string) string, d Dialect) string {
	var cols []string
	for _, c := range ix.Columns {
		if m := prefixPart.FindStringSubmatch(c.Value); d == MySQL && c.Type == "expression" && m != nil && t.findField(m[1]) != nil {
			// a MySQL prefix key part `col(n)` (a column of the table)
			cols = append(cols, quote(m[1])+"("+m[2]+")")
			continue
		}
		if c.Type == "expression" {
			cols = append(cols, "("+safeExpr(c.Value, d, "index expression")+")")
		} else {
			cols = append(cols, quote(c.Value))
		}
	}
	return "(" + strings.Join(cols, ", ") + ")"
}

func fieldCheck(f *cField, quote func(string) string, d Dialect) string {
	where := "check on column " + quote(f.name)
	if len(f.checks) == 0 {
		return ""
	}
	if len(f.checks) == 1 {
		s := ""
		if f.checks[0].name != nil && *f.checks[0].name != "" {
			s += " CONSTRAINT " + quote(*f.checks[0].name)
		}
		return s + " CHECK (" + safeExpr(f.checks[0].expression, d, where) + ")"
	}
	var parts []string
	for _, ch := range f.checks {
		parts = append(parts, "("+safeExpr(ch.expression, d, where)+")")
	}
	return " CHECK (" + strings.Join(parts, " AND ") + ")"
}

func tableChecks(t *cTable, quote func(string) string, d Dialect) []string {
	var out []string
	for _, ch := range t.checks {
		line := ""
		if ch.name != nil && *ch.name != "" {
			line = "CONSTRAINT " + quote(*ch.name) + " "
		}
		out = append(out, line+"CHECK ("+safeExpr(ch.expression, d, "check on table "+quote(t.name))+")")
	}
	return out
}

func compositePKs(t *cTable, quote func(string) string, d Dialect) []string {
	var out []string
	for _, ix := range t.indexes {
		if truthy(ix.PK) {
			out = append(out, "PRIMARY KEY "+indexColumns(t, ix, quote, d))
		}
	}
	return out
}

func fieldNames(fs []*cField, quote func(string) string) string {
	var parts []string
	for _, f := range fs {
		parts = append(parts, quote(f.name))
	}
	return "(" + strings.Join(parts, ", ") + ")"
}

var (
	numberText = regexp.MustCompile(`^[-+]?([0-9]+\.?[0-9]*|\.[0-9]+)([eE][-+]?[0-9]+)?$`)
	hexDigits  = regexp.MustCompile(`^[0-9A-Fa-f]*$`)
	boolText   = regexp.MustCompile(`(?i)^(true|false|null|0|1)$`)
	prefixPart = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_$]*)\s*\(\s*([0-9]+)\s*\)\s*$`)
)

// safeNumber returns s when it is a numeric literal and a quoted string
// literal otherwise. The reference prints a numeric column's record value
// raw even when it is an arbitrary string; go-dbml never emits such text
// unquoted (a deliberate, documented deviation).
func safeNumber(s string, lit func(string) string) string {
	if numberText.MatchString(s) {
		return s
	}
	return lit(s)
}

// joinQ quotes each name and joins them with ", ".
func joinQ(names []string, quote func(string) string) string {
	out := make([]string, len(names))
	for i, n := range names {
		out[i] = quote(n)
	}
	return strings.Join(out, ", ")
}

// sqlDefault renders a column DEFAULT. String defaults go through the
// dialect's literal escaper; expression defaults through safeExpr; a
// number or boolean default that is not a number/boolean literal is
// quoted.
func sqlDefault(dflt *model.Default, d Dialect, where string) string {
	if dflt == nil || isNullDefault(dflt) {
		return ""
	}
	lit := litFor(d)
	v := jsString(dflt.Value)
	switch dflt.Type {
	case "expression":
		return " DEFAULT (" + safeExpr(v, d, where) + ")"
	case "string":
		return " DEFAULT " + lit(v)
	case "boolean":
		if !boolText.MatchString(v) {
			v = lit(v)
		}
	default:
		v = safeNumber(v, lit)
	}
	if d == MSSQL {
		return " DEFAULT (" + v + ")"
	}
	return " DEFAULT " + v
}

// litFor returns the string literal escaper of d.
func litFor(d Dialect) func(string) string {
	switch d {
	case MySQL:
		return myLit
	case Postgres:
		return pgLit
	}
	return stdLit
}

// genDefault is f's DEFAULT, or nil for a generated column (which cannot
// have one).
func genDefault(f *cField) *model.Default {
	if f.generated != nil {
		return nil
	}
	return f.dbdefault
}
