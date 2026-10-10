// Portions of this file are derived from holistics/dbml
// packages/dbml-parse/src/core/global_modules/records/interpret.ts,
// records/utils/data/{sqlTypes,values}.ts and records/utils/constraints/*.ts
// (Apache-2.0, Copyright Holistics Software). Modified: ported to Go.

package model

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/aidenappl/go-dbml/ast"
	"github.com/aidenappl/go-dbml/diag"
)

// ---------------------------------------------------------------- records block

// records interprets one records block (top-level `records t(a, b) {}` or
// nested `records (a, b) {}` / `records {}` inside a table). Type-mismatch
// problems are reported as INVALID_RECORDS_FIELD infos, never errors, and
// the offending value is kept with a best-effort type.
func (in *interpreter) records(e *ast.ElementDeclaration, owner *sym) (*TableRecord, *sym) {
	file := in.rfile
	tv, colNodes, _ := in.c.recordsTarget(in.rs, e, owner)
	if tv == nil || !in.visibleSym(tv.sym) {
		return nil, nil
	}
	all := in.c.columns(tv.sym, tv.sym.unit.scope)
	var cols []*col
	if owner != nil && ast.IsNil(e.Name) {
		// implicit column list: every merged column in order
		cols = all
	} else {
		for _, cn := range colNodes {
			name, _ := extractVariable(cn)
			for _, c := range all {
				if c.name == name {
					cols = append(cols, c)
					break
				}
			}
		}
	}
	if len(cols) == 0 {
		return nil, nil
	}
	schema, name := in.sc.interpretedName(tv.sym)
	rec := &TableRecord{SchemaName: schema, TableName: name, Columns: colNames(cols), Values: [][]*RecordValue{}, Token: nodeToken(e, file)}
	set, _ := aggregateSettings(e.AttributeList, file)
	if set.has(settingExample) {
		rec.Example = boolPtr(true)
	}
	if b := e.Base(); b.HasPos {
		size := int(b.End.Byte - b.Start.Byte)
		if in.recordBytes+size > in.c.opts.MaxRecordBytes {
			in.recordsLimit(e, file, fmt.Sprintf("Records exceed the %d-byte limit; this block was skipped", in.c.opts.MaxRecordBytes))
			return rec, tv.sym
		}
		in.recordBytes += size
	}
	fields, _ := partitionBody(elementBody(e))
	for _, row := range fields {
		in.c.poll(false)
		if in.recordRows >= in.c.opts.MaxRecordRows {
			in.recordsLimit(row, file, fmt.Sprintf("Records exceed the %d-row limit; the remaining rows were dropped", in.c.opts.MaxRecordRows))
			break
		}
		in.recordRows++
		args := []ast.Node{row.Callee}
		if ce, ok := row.Callee.(*ast.CommaExpression); ok {
			args = ce.Elements
		}
		if len(args) != len(cols) {
			in.n(diag.INVALID_RECORDS_FIELD, "Expected "+strconv.Itoa(len(cols))+" values but got "+strconv.Itoa(len(args)), row, file)
			continue
		}
		vals := make([]*RecordValue, 0, len(cols))
		for i, a := range args {
			vals = append(vals, in.recordValue(a, cols[i], tv.sym, file))
		}
		rec.Values = append(rec.Values, vals)
	}
	return rec, tv.sym
}

// recordsLimit reports the records budget being exhausted (once per compile).
func (in *interpreter) recordsLimit(n ast.Node, file, msg string) {
	if in.recordsCapped {
		return
	}
	in.recordsCapped = true
	in.n(diag.RECORDS_TOO_LARGE, msg, n, file)
}

// ---------------------------------------------------------------- SQL type families

// The type families union every dialect @dbml/parse knows (it checks "any
// dialect has this type").
var (
	integerTypes  = setOf("int", "integer", "smallint", "bigint", "tinyint", "mediumint", "serial", "bigserial", "smallserial")
	floatTypes    = setOf("decimal", "numeric", "float", "double", "real", "double precision", "number")
	boolTypes     = setOf("bool", "boolean", "bit")
	stringTypes   = setOf("varchar", "char", "text", "tinytext", "mediumtext", "longtext", "string", "character", "character varying", "nvarchar", "nchar", "ntext", "varchar2", "nvarchar2")
	datetimeTypes = setOf("date", "datetime", "timestamp", "time", "timestamptz", "timestamp with time zone", "timestamp without time zone", "timetz", "time with time zone", "time without time zone", "datetime2", "smalldatetime", "timestamp with local time zone")
	serialTypes   = setOf("serial", "smallserial", "bigserial")
)

// setOf builds a string set.
func setOf(xs ...string) map[string]bool {
	m := map[string]bool{}
	for _, x := range xs {
		m[x] = true
	}
	return m
}

// normalizeTypeName lower-cases, trims and collapses internal whitespace.
func normalizeTypeName(t string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(t))), " ")
}

// recordValueType maps an SQL type to the record value type reported in
// RecordValue.Type; unknown types are reported as written.
func recordValueType(sqlType string, isEnum bool) string {
	t := normalizeTypeName(sqlType)
	switch {
	case isEnum:
		return "string"
	case integerTypes[t]:
		return "integer"
	case floatTypes[t]:
		return "real"
	case boolTypes[t]:
		return "bool"
	case stringTypes[t]:
		return "string"
	case t == "date":
		return "date"
	case t == "time" || t == "timetz" || t == "time with time zone" || t == "time without time zone":
		return "time"
	case datetimeTypes[t]:
		return "datetime"
	}
	return sqlType
}

// colTypeInfo mirrors ColumnSymbol.type(): the base name without array
// suffix, numeric args, and the enum it resolves to. Like @dbml/parse only a
// bare type name is resolved to an enum (`status`, not `app.status`).
type colTypeInfo struct {
	base  string
	args  []string
	enum  *sym
	valid bool
}

// columnType computes colTypeInfo for a column (see colTypeInfo).
func (in *interpreter) columnType(c *col, owner *sym) colTypeInfo {
	if c.fa == nil || len(c.fa.Args) == 0 {
		return colTypeInfo{}
	}
	ct, _ := processColumnType(c.fa.Args[0])
	info := colTypeInfo{valid: true}
	base := ct.TypeName
	// strip `[...]` array parts, then everything from `(`
	for {
		i := strings.Index(base, "[")
		if i < 0 {
			break
		}
		j := strings.Index(base[i:], "]")
		if j < 0 {
			break
		}
		base = base[:i] + base[i+j+1:]
	}
	info.base = strings.SplitN(base, "(", 2)[0]
	if ct.Args != nil && *ct.Args != "" {
		info.args = strings.Split(*ct.Args, ",")
	}
	if isVariableNode(c.fa.Args[0]) {
		sc := owner.unit.scope
		if c.partial != nil {
			sc = c.partial.unit.scope
		}
		name, _ := extractVariable(c.fa.Args[0])
		if v := sc.lookup(symEnum, "", name); v != nil {
			info.enum = v.sym
		}
	}
	return info
}

// ---------------------------------------------------------------- value typing

// info records an info-severity diagnostic over n.
func (in *interpreter) info(code diag.Code, msg string, n ast.Node, file string) {
	in.diags = append(in.diags, diag.NewWithSeverity(diag.SeverityInfo, code, msg, file, nodeRange(n)))
}

// recordValue mirrors extractValue in records/interpret.ts.
func (in *interpreter) recordValue(n ast.Node, c *col, owner *sym, file string) *RecordValue {
	tok := nodeToken(n, file)
	if fe, ok := n.(*ast.FunctionExpression); ok {
		v := ""
		if fe.Value != nil {
			v = fe.Value.Value
		}
		return &RecordValue{Value: v, Type: "expression", Token: tok}
	}
	ti := in.columnType(c, owner)
	typeName := ti.base
	norm := normalizeTypeName(typeName)
	set := in.colSettings(c)
	increment := set.has(settingIncrement) || serialTypes[norm]
	notNull := set.has(settingNotNull)
	dbdefault := processDefault(set.firstValue(settingDefault))
	vt := recordValueType(typeName, ti.enum != nil)
	raw, hasRaw := tryExtractString(n)
	fallback := func() *RecordValue {
		if hasRaw {
			return &RecordValue{Value: raw, Type: vt, Token: tok}
		}
		return &RecordValue{Value: "", Type: "expression", Token: tok}
	}
	colName := c.name

	// NULL literal (or '' for a non-string column, or an empty CSV field)
	if s, isStr := extractQuotedString(n); isNullish(n) || (isStr && s == "" && !stringTypes[norm]) {
		hasDefault := dbdefault != nil && strings.ToLower(defaultString(dbdefault.Value)) != "null"
		if notNull && !hasDefault && !increment && !serialTypes[norm] {
			in.info(diag.INVALID_RECORDS_FIELD, "NULL not allowed for NOT NULL column `"+colName+"` without a default value or auto-increment", n, file)
		}
		return &RecordValue{Value: nil, Type: vt, Token: tok}
	}

	if ti.enum != nil {
		v, ok := extractQuotedString(n)
		if !ok {
			v, ok = extractVariable(n)
		}
		if !ok || !enumHasField(ti.enum, v) {
			in.info(diag.INVALID_RECORDS_FIELD, "Invalid enum value for column `"+colName+"`", n, file)
		}
		if !ok {
			// `Enum.member` is not a bare variable: @dbml/parse leaves the
			// value undefined (the key is omitted from the JSON).
			return &RecordValue{Undefined: true, Type: vt, Token: tok}
		}
		return &RecordValue{Value: v, Type: vt, Token: tok}
	}

	switch {
	case integerTypes[norm] || floatTypes[norm] || serialTypes[norm]:
		num, ok := extractSignedNumber(n)
		if !ok {
			in.info(diag.INVALID_RECORDS_FIELD, "Invalid numeric value for column `"+colName+"`", n, file)
			return fallback()
		}
		if integerTypes[norm] && num != math.Trunc(num) {
			in.info(diag.INVALID_RECORDS_FIELD, "Invalid integer value `"+formatJSNumber(num)+"` for column `"+colName+"`: expected integer, got decimal", n, file)
			return &RecordValue{Value: math.Floor(num), Type: vt, Token: tok}
		}
		if floatTypes[norm] && len(ti.args) > 0 {
			precision, _ := jsNumber(ti.args[0])
			scale := 0.0
			if len(ti.args) >= 2 {
				scale, _ = jsNumber(ti.args[1])
			}
			intPart, decPart, _ := strings.Cut(formatJSNumber(num), ".")
			intPart = strings.TrimPrefix(intPart, "-")
			total := len(intPart) + len(decPart)
			if float64(total) > precision {
				in.info(diag.INVALID_RECORDS_FIELD, "Numeric value `"+formatJSNumber(num)+"` for column `"+colName+"` exceeds precision: expected at most "+formatJSNumber(precision)+" total digits, got "+strconv.Itoa(total), n, file)
			} else if float64(len(decPart)) > scale {
				in.info(diag.INVALID_RECORDS_FIELD, "Numeric value `"+formatJSNumber(num)+"` for column `"+colName+"` exceeds scale: expected at most "+formatJSNumber(scale)+" decimal digits, got "+strconv.Itoa(len(decPart)), n, file)
			}
		}
		return &RecordValue{Value: num, Type: vt, Token: tok}
	case boolTypes[norm]:
		b, ok := tryExtractBool(n)
		if !ok {
			in.info(diag.INVALID_RECORDS_FIELD, "Invalid boolean value for column `"+colName+"`", n, file)
			return fallback()
		}
		return &RecordValue{Value: b, Type: vt, Token: tok}
	case datetimeTypes[norm]:
		d, ok := tryExtractDateTime(n)
		if !ok {
			in.info(diag.INVALID_RECORDS_FIELD, "Invalid datetime value for column `"+colName+"`", n, file)
			return fallback()
		}
		return &RecordValue{Value: d, Type: vt, Token: tok}
	case stringTypes[norm]:
		if !hasRaw {
			in.info(diag.INVALID_RECORDS_FIELD, "Invalid string value for column `"+colName+"`", n, file)
			return fallback()
		}
		if len(ti.args) > 0 {
			if max, ok := jsNumber(ti.args[0]); ok && max > 0 {
				if l := len(raw); float64(l) > max {
					in.info(diag.INVALID_RECORDS_FIELD, "String value for column `"+colName+"` exceeds maximum length: expected at most "+formatJSNumber(max)+" bytes (UTF-8), got "+strconv.Itoa(l)+" bytes", n, file)
				}
			}
		}
		return &RecordValue{Value: raw, Type: "string", Token: tok}
	}
	return fallback()
}

// defaultString renders a default value like JavaScript String(value).
func defaultString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return formatJSNumber(t)
	}
	return ""
}

// isNullish matches the `null` identifier and empty CSV fields.
func isNullish(n ast.Node) bool {
	if isIdentifierNode(n) {
		return strings.ToLower(variableOf(n).Value) == "null"
	}
	_, ok := n.(*ast.Empty)
	return ok
}

// extractSignedNumber accepts 42, -42, +3.14 (numeric literals only).
func extractSignedNumber(n ast.Node) (float64, bool) {
	if v, ok := extractNumericLiteral(n); ok {
		return v, true
	}
	if isSignedNumber(n) {
		if pre, ok := n.(*ast.PrefixExpression); ok {
			if v, ok := extractNumericLiteral(pre.Expression); ok {
				if pre.Op != nil && pre.Op.Value == "-" {
					return -v, true
				}
				return v, true
			}
		}
	}
	return 0, false
}

var truthy = setOf("true", "yes", "y", "t", "1")
var falsy = setOf("false", "no", "n", "f", "0")

// tryExtractBool accepts true/false identifiers, the literals 0/1 and the
// quoted strings yes/no/y/n/t/f/1/0/true/false (case-insensitive).
func tryExtractBool(n ast.Node) (bool, bool) {
	if isIdentifierNode(n) {
		switch strings.ToLower(variableOf(n).Value) {
		case "true":
			return true, true
		case "false":
			return false, true
		}
	}
	if v, ok := extractNumericLiteral(n); ok {
		if v == 0 {
			return false, true
		}
		if v == 1 {
			return true, true
		}
	}
	if s, ok := extractQuotedString(n); ok && s != "" {
		l := strings.ToLower(s)
		if truthy[l] {
			return true, true
		}
		if falsy[l] {
			return false, true
		}
	}
	return false, false
}

// tryExtractNumeric accepts signed numeric literals and quoted numbers.
func tryExtractNumeric(n ast.Node) (float64, bool) {
	if v, ok := extractSignedNumber(n); ok {
		return v, true
	}
	if s, ok := extractQuotedString(n); ok {
		if v, ok := jsNumber(s); ok {
			return v, true
		}
	}
	return 0, false
}

// tryExtractString: quoted string, else a number, else a boolean, as text.
func tryExtractString(n ast.Node) (string, bool) {
	if s, ok := extractQuotedString(n); ok {
		return s, true
	}
	if v, ok := tryExtractNumeric(n); ok {
		return formatJSNumber(v), true
	}
	if b, ok := tryExtractBool(n); ok {
		return strconv.FormatBool(b), true
	}
	return "", false
}

// formatJSNumber formats like JavaScript's Number.prototype.toString.
func formatJSNumber(f float64) string {
	if math.IsNaN(f) {
		return "NaN"
	}
	if math.IsInf(f, 0) {
		if f > 0 {
			return "Infinity"
		}
		return "-Infinity"
	}
	abs := math.Abs(f)
	if abs == 0 || (abs >= 1e-6 && abs < 1e21) {
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

// Supported date/time formats (luxon tokens in @dbml/parse).
var (
	timeLayouts     = []string{"15:04:05", "15:04:05Z07:00", "15:04:05.000", "15:04:05.000Z07:00"}
	dateLayouts     = []string{"2006-01-02", "1/2/2006", "2 Jan 2006", "Jan 2, 2006"}
	datetimeLayouts = []string{"2006-01-02 15:04:05", "2006-01-02 15:04:05.000", "2006-01-02 15:04:05.000Z07:00", "2006-01-02 15:04:05Z07:00"}
	isoLayouts      = []string{time.RFC3339Nano, "2006-01-02T15:04:05", "2006-01-02T15:04:05.000", "2006-01-02T15:04", "2006-01-02T15"}
)

// tryExtractDateTime normalises a quoted date/time to ISO 8601. It follows
// the reference's order (time-only, date-only, datetime, ISO) and its
// output (milliseconds suppressed when zero, offset only when explicit), but
// uses Go layouts instead of luxon, so exotic ISO forms may differ.
func tryExtractDateTime(n ast.Node) (string, bool) {
	s, ok := extractQuotedString(n)
	if !ok {
		return "", false
	}
	return NormalizeDateTime(s)
}

// NormalizeDateTime normalises a date, time or datetime string to ISO 8601
// the way @dbml/parse's tryExtractDateTime does (see tryExtractDateTime for
// the supported formats); ok is false when s is not a recognised value.
func NormalizeDateTime(s string) (string, bool) {
	zoned := func(layout string) bool { return strings.Contains(layout, "Z07") }
	frac := func(t time.Time) string {
		if t.Nanosecond() != 0 {
			return ".000"
		}
		return ""
	}
	for _, l := range timeLayouts {
		if t, err := time.Parse(l, s); err == nil {
			out := t.Format("15:04:05" + frac(t))
			if zoned(l) {
				out += t.Format("Z07:00")
			}
			return out, true
		}
	}
	for _, l := range dateLayouts {
		if t, err := time.Parse(l, s); err == nil {
			return t.Format("2006-01-02"), true
		}
	}
	for _, l := range append(append([]string{}, datetimeLayouts...), isoLayouts...) {
		if t, err := time.Parse(l, s); err == nil {
			out := t.Format("2006-01-02T15:04:05" + frac(t))
			if zoned(l) {
				out += t.Format("Z07:00")
			}
			return out, true
		}
	}
	return "", false
}

// ---------------------------------------------------------------- record constraints

// tableRecords pairs a table with its (single) records block.
type tableRecords struct {
	sym   *sym
	rec   *TableRecord
	block ast.Node
	file  string
}

type keyedRow map[string]*RecordValue

// keyedRows converts positional rows to column-name keyed rows.
func keyedRows(r *TableRecord) []keyedRow {
	out := make([]keyedRow, len(r.Values))
	for i, row := range r.Values {
		k := keyedRow{}
		for j, c := range r.Columns {
			if j < len(row) {
				k[c] = row[j]
			}
		}
		out[i] = k
	}
	return out
}

// jsonValue is JSON.stringify(value) ("" for a missing/undefined value).
func jsonValue(v *RecordValue) string {
	if v == nil || v.Undefined {
		return ""
	}
	b, _ := json.Marshal(v.Value)
	return string(b)
}

// isNullValue reports a missing, undefined or null value.
func isNullValue(v *RecordValue) bool { return v == nil || v.Undefined || v.Value == nil }

// colFacts are the column properties the constraint checks need.
type colFacts struct {
	c         *col
	pk        bool
	unique    bool
	increment bool
	def       *Default
}

// facts collects a column's pk/unique/increment/default.
func (in *interpreter) facts(c *col, owner *sym) colFacts {
	s := in.colSettings(c)
	f := colFacts{c: c, pk: s.has(settingPK) || s.has(settingPrimaryKey), unique: s.has(settingUnique), def: processDefault(s.firstValue(settingDefault))}
	f.increment = s.has(settingIncrement) || serialTypes[normalizeTypeName(in.columnType(c, owner).base)]
	return f
}

// keyWithDefault builds a constraint key, substituting column defaults for
// null values (extractKeyValueWithDefault).
func keyWithDefault(row keyedRow, cols []colFacts) string {
	parts := make([]string, len(cols))
	for i, c := range cols {
		v := row[c.c.name]
		if isNullValue(v) && c.def != nil {
			b, _ := json.Marshal(c.def.Value)
			parts[i] = string(b)
			continue
		}
		parts[i] = jsonValue(v)
	}
	return strings.Join(parts, "|")
}

// nullWithoutDefault reports a null key column that has no default.
func nullWithoutDefault(row keyedRow, cols []colFacts) bool {
	for _, c := range cols {
		if isNullValue(row[c.c.name]) && c.def == nil {
			return true
		}
	}
	return false
}

// formatValues renders key values for messages: `1` or `(1, "US")`.
func formatValues(row keyedRow, cols []colFacts) string {
	parts := make([]string, len(cols))
	for i, c := range cols {
		if v, ok := row[c.c.name]; ok {
			parts[i] = jsonValue(v)
		} else if c.def != nil {
			b, _ := json.Marshal(c.def.Value)
			parts[i] = string(b)
		}
	}
	if len(cols) == 1 {
		return parts[0]
	}
	return "(" + strings.Join(parts, ", ") + ")"
}

// anchors returns the row values to attach a diagnostic to: the named
// columns if present, otherwise every value in the row (in column order).
func anchors(row keyedRow, rec *TableRecord, names []string) []*RecordValue {
	var out []*RecordValue
	for _, n := range names {
		if v, ok := row[n]; ok && v != nil {
			out = append(out, v)
		}
	}
	if len(out) > 0 {
		return out
	}
	for _, c := range rec.Columns {
		if v := row[c]; v != nil {
			out = append(out, v)
		}
	}
	return out
}

// infoAt records an INVALID_RECORDS_FIELD info at a record value.
func (in *interpreter) infoAt(v *RecordValue, msg string) {
	in.diags = append(in.diags, diag.NewWithSeverity(diag.SeverityInfo, diag.INVALID_RECORDS_FIELD, msg, v.Token.Filepath,
		diag.Range{Start: v.Token.Start, End: v.Token.End}))
}

// columnRef renders `[schema.]table.col` or a parenthesised list.
func (in *interpreter) columnRef(t *sym, cols []colFacts) string {
	schema, name := in.sc.interpretedName(t)
	var names []string
	for _, c := range cols {
		q := name + "." + c.c.name
		if schema != nil {
			q = *schema + "." + q
		}
		names = append(names, q)
	}
	if len(names) == 1 {
		return names[0]
	}
	return "(" + strings.Join(names, ", ") + ")"
}

// label prefixes "Composite " for multi-column constraints.
func label(cols []colFacts, kind string) string {
	if len(cols) > 1 {
		return "Composite " + kind
	}
	return kind
}

// constraintSets returns single-column constraints from column settings and
// composite ones from indexes (own and injected), as @dbml/parse does.
func (in *interpreter) constraintSets(t *sym, pk bool) [][]colFacts {
	cols := in.c.columns(t, t.unit.scope)
	byName := map[string]colFacts{}
	var out [][]colFacts
	for _, c := range cols {
		f := in.facts(c, t)
		byName[c.name] = f
		if (pk && f.pk) || (!pk && f.unique) {
			out = append(out, []colFacts{f})
		}
	}
	for _, ix := range in.tableIndexes(t) {
		if (pk && (ix.PK == nil || !*ix.PK)) || (!pk && (ix.Unique == nil || !*ix.Unique)) {
			continue
		}
		var set []colFacts
		for _, ic := range ix.Columns {
			if f, ok := byName[ic.Value]; ok {
				set = append(set, f)
			}
		}
		out = append(out, set)
	}
	return out
}

// names returns the column names of a constraint.
func names(cols []colFacts) []string {
	out := make([]string, len(cols))
	for i, c := range cols {
		out[i] = c.c.name
	}
	return out
}

// groupDuplicates returns, in first-seen order, the groups of rows sharing
// a key, keeping only groups with more than one row.
func groupDuplicates(rows []keyedRow, key func(keyedRow) string) [][]keyedRow {
	idx := map[string]int{}
	var groups [][]keyedRow
	for _, r := range rows {
		k := key(r)
		if i, ok := idx[k]; ok {
			groups[i] = append(groups[i], r)
			continue
		}
		idx[k] = len(groups)
		groups = append(groups, []keyedRow{r})
	}
	var out [][]keyedRow
	for _, g := range groups {
		if len(g) > 1 {
			out = append(out, g)
		}
	}
	return out
}

// validateRecordConstraints reports PK, UNIQUE and FK violations in records
// (validateRecords in program/utils/records.ts).
func (in *interpreter) validateRecordConstraints(db *Database, blocks []tableRecords) {
	byTable := map[*sym]*tableRecords{}
	for i := range blocks {
		b := &blocks[i]
		if len(b.rec.Values) == 0 {
			if _, ok := byTable[b.sym]; !ok {
				byTable[b.sym] = b
			}
			continue
		}
		in.validatePK(b)
		in.validateUnique(b)
		byTable[b.sym] = b
	}
	refs := append([]*Ref{}, db.Refs...)
	for _, t := range db.Tables {
		refs = append(refs, db.ResolveTable(t).InjectedRefs...)
	}
	lookup := map[string]*tableRecords{}
	for _, v := range in.sc.entries {
		if v.sym.kind != symTable || v.sym.placeholder {
			continue
		}
		schema, name := in.sc.interpretedName(v.sym)
		k := tableKey(orPublic(schema), name)
		if b, ok := byTable[v.sym]; ok {
			lookup[k] = b
		} else if _, ok := lookup[k]; !ok {
			lookup[k] = &tableRecords{sym: v.sym}
		}
	}
	for _, r := range refs {
		k1 := tableKey(orPublic(r.Endpoints[0].SchemaName), r.Endpoints[0].TableName)
		k2 := tableKey(orPublic(r.Endpoints[1].SchemaName), r.Endpoints[1].TableName)
		t1, t2 := lookup[k1], lookup[k2]
		if t1 == nil || t2 == nil {
			continue
		}
		if (t1.rec == nil || len(t1.rec.Values) == 0) && (t2.rec == nil || len(t2.rec.Values) == 0) {
			continue
		}
		e1 := in.endpointFacts(t1.sym, r.Endpoints[0].FieldNames)
		e2 := in.endpointFacts(t2.sym, r.Endpoints[1].FieldNames)
		c1, c2 := parseCard(r.Endpoints[0].Relation), parseCard(r.Endpoints[1].Relation)
		oneToOne := !c1.many && !c2.many
		skip1 := (!c1.many && c2.min == 0) || oneToOne
		skip2 := !c2.many && c1.min == 0
		if !skip1 {
			in.validateFKEndpoint(t1, e1, t2, e2, c1, c2)
		}
		if !skip2 {
			in.validateFKEndpoint(t2, e2, t1, e1, c2, c1)
		}
	}
}

// orPublic maps a nil schema to "public" (makeTableKey).
func orPublic(s *string) *string {
	if s == nil {
		return strPtr(defaultSchema)
	}
	return s
}

// endpointFacts resolves ref endpoint field names to column facts.
func (in *interpreter) endpointFacts(t *sym, fieldNames []string) []colFacts {
	cols := in.c.columns(t, t.unit.scope)
	var out []colFacts
	for _, fn := range fieldNames {
		for _, c := range cols {
			if c.name == fn {
				out = append(out, in.facts(c, t))
				break
			}
		}
	}
	return out
}

// validatePK reports missing, null and duplicate primary keys
// (constraints/pk.ts). Auto-increment keys are never duplicates.
func (in *interpreter) validatePK(b *tableRecords) {
	rows := keyedRows(b.rec)
	for _, pk := range in.constraintSets(b.sym, true) {
		present := map[string]bool{}
		for _, c := range b.rec.Columns {
			present[c] = true
		}
		var missing []colFacts
		for _, c := range pk {
			if !present[c.c.name] && !c.increment && c.def == nil {
				missing = append(missing, c)
			}
		}
		anyMissing := false
		for _, c := range pk {
			if !present[c.c.name] {
				anyMissing = true
			}
		}
		if anyMissing && len(missing) > 0 {
			msg := label(missing, "PK") + ": Column `" + in.columnRef(b.sym, missing) + "` is missing from record and has no default value"
			for range b.rec.Values {
				in.diags = append(in.diags, diag.NewWithSeverity(diag.SeverityInfo, diag.INVALID_RECORDS_FIELD, msg, b.file, nodeRange(b.block)))
			}
			continue
		}
		var nullCheck []colFacts
		autoInc := false
		for _, c := range pk {
			if !c.increment && c.def == nil {
				nullCheck = append(nullCheck, c)
			}
			if c.increment {
				autoInc = true
			}
		}
		var withoutNull []keyedRow
		for _, row := range rows {
			if nullWithoutDefault(row, nullCheck) {
				for _, v := range anchors(row, b.rec, names(pk)) {
					in.infoAt(v, "NULL in "+label(nullCheck, "PK")+": `"+in.columnRef(b.sym, nullCheck)+"` cannot be NULL")
				}
				continue
			}
			withoutNull = append(withoutNull, row)
		}
		if autoInc {
			continue
		}
		for _, g := range groupDuplicates(withoutNull, func(r keyedRow) string { return keyWithDefault(r, pk) }) {
			for _, row := range g {
				msg := "Duplicate " + label(pk, "PK") + ": `" + in.columnRef(b.sym, pk) + "` = `" + formatValues(row, pk) + "`"
				for _, v := range anchors(row, b.rec, names(pk)) {
					in.infoAt(v, msg)
				}
			}
		}
	}
}

// validateUnique reports duplicate unique keys; nulls never conflict
// (constraints/unique.ts).
func (in *interpreter) validateUnique(b *tableRecords) {
	rows := keyedRows(b.rec)
	for _, uq := range in.constraintSets(b.sym, false) {
		var candidates []keyedRow
		for _, row := range rows {
			if !nullWithoutDefault(row, uq) {
				candidates = append(candidates, row)
			}
		}
		for _, g := range groupDuplicates(candidates, func(r keyedRow) string { return keyWithDefault(r, uq) }) {
			for _, row := range g {
				msg := "Duplicate " + label(uq, "UNIQUE") + ": `" + in.columnRef(b.sym, uq) + "` = `" + formatValues(row, uq) + "`"
				for _, v := range anchors(row, b.rec, names(uq)) {
					in.infoAt(v, msg)
				}
			}
		}
	}
}

// validateFKEndpoint checks left's rows against right's cardinality.
func (in *interpreter) validateFKEndpoint(left *tableRecords, le []colFacts, right *tableRecords, re []colFacts, lc, rc card) {
	if left.rec == nil || len(left.rec.Values) == 0 {
		return
	}
	valid := map[string]bool{}
	if right.rec != nil {
		for _, row := range keyedRows(right.rec) {
			valid[keyWithDefault(row, re)] = true
		}
	}
	leftRef, rightRef := in.columnRef(left.sym, le), in.columnRef(right.sym, re)
	_ = rightRef
	for _, row := range keyedRows(left.rec) {
		if nullWithoutDefault(row, le) {
			if rc.min == 0 {
				continue
			}
			msg := "FK violation: `" + leftRef + "` = `" + formatValues(row, le) + "` must not be NULL"
			for _, v := range anchors(row, left.rec, names(le)) {
				in.infoAt(v, msg)
			}
			continue
		}
		if rc.many && lc.min == 0 && lc.many {
			continue
		}
		if valid[keyWithDefault(row, le)] {
			continue
		}
		msg := "FK violation: `" + leftRef + "` = `" + formatValues(row, le) + "` does not exist in `" + rightRef + "`"
		for _, v := range anchors(row, left.rec, names(le)) {
			in.infoAt(v, msg)
		}
	}
}
