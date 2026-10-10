// Portions of this file are derived from holistics/dbml packages/dbml-core/src/export/DbmlExporter.ts,
// packages/dbml-core/src/model_structure/database.ts (adjustCardinalitiesFromColumns) and
// packages/dbml-parse/src/compiler/queries/utils.ts (formatRecordValue, escapeString)
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

package export

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/aidenappl/go-dbml/model"
)

// CoreDBMLOptions tunes CoreDBML (the @dbml/core DbmlExporter port).
type CoreDBMLOptions struct {
	// OmitRecords drops Records blocks (@dbml/core includeRecords: false).
	OmitRecords bool
	// OverrideCardinality applies @dbml/core's adjustCardinalitiesFromColumns
	// (SQL importers mark every FK target optional; a target whose FK
	// columns are all NOT NULL/PK/increment becomes required). The SQL
	// importers set it.
	OverrideCardinality bool
	// Escape quotes and escapes names, strings and expressions so the
	// output always lexes (DbmlExporter escapes nothing). Off by default to
	// stay byte-identical with @dbml/core.
	Escape bool
}

// CoreDBML renders db the way @dbml/core 10.3.1's DbmlExporter does
// (exporter.export(src, "dbml") / importer.import(sql, format)): every name
// double-quoted, refs written with the "one" side on the left, records
// grouped per table. For go-dbml's own canonical formatting use package
// emit instead.
func CoreDBML(db *model.Database, opts CoreDBMLOptions) (out string, err error) {
	if db == nil {
		return "", fmt.Errorf("export: nil database")
	}
	defer func() {
		if r := recover(); r != nil {
			out, err = "", &InternalError{Value: r}
		}
	}()
	c, err := buildCore(db)
	if err != nil {
		return "", err
	}
	if opts.OverrideCardinality {
		c.adjustCardinalities()
	}
	return exportDBML(c, opts), nil
}

func maxIsMany(rel string) bool {
	if rel == "*" {
		return true
	}
	return strings.HasSuffix(rel, "..*")
}

func makeRequired(rel string) string {
	min := rel
	if i := strings.Index(rel, ".."); i >= 0 {
		min = rel[:i]
	}
	if rel == "*" || (min != "0" && min != "") {
		return rel
	}
	if maxIsMany(rel) {
		return "*"
	}
	return "1"
}

// adjustCardinalities mirrors Database.adjustCardinalitiesFromColumns.
func (c *core) adjustCardinalities() {
	for _, s := range c.schemas {
		for _, r := range s.refs {
			var many, one *cEndpoint
			for _, ep := range r.endpoints {
				if many == nil && maxIsMany(ep.relation) {
					many = ep
				}
				if one == nil && maxIsOne(ep.relation) {
					one = ep
				}
			}
			if one == nil {
				continue
			}
			src, dst := many, one
			if many == nil {
				src, dst = r.endpoints[1], r.endpoints[0]
			}
			all := true
			for _, f := range src.fields {
				if !(f.notNull || f.pk || f.increment) {
					all = false
					break
				}
			}
			if all {
				dst.relation = makeRequired(dst.relation)
			}
		}
	}
}

var (
	reWhiteSpace     = regexp.MustCompile(`\s`)
	reSquareBracket  = regexp.MustCompile(`\[|\]`)
	reNoteNeedsTripe = regexp.MustCompile(`[\n\r']`)
)

func escapeNote(s string) string {
	n := strings.ReplaceAll(s, `\`, `\\`)
	if !reNoteNeedsTripe.MatchString(n) {
		return "'" + n + "'"
	}
	n = strings.ReplaceAll(n, "'", `\'`)
	n = strings.ReplaceAll(n, "\r\n", "\n")
	return "'''" + n + "'''"
}

// escapeString mirrors @dbml/parse escapeString.
func escapeString(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\'':
			b.WriteString(`\'`)
		case '\n':
			b.WriteString(`\n`)
		case '\t':
			b.WriteString(`\t`)
		case '\r':
			b.WriteString(`\r`)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		case '\v':
			b.WriteString(`\v`)
		case 0:
			b.WriteString(`\0`)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func isValidIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 && r >= '0' && r <= '9' {
			return false
		}
		if !(r == '_' || unicode.IsLetter(r) || unicode.Is(unicode.M, r) || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

func addDoubleQuoteIfNeeded(s string) string {
	if isValidIdentifier(s) {
		return s
	}
	return `"` + escapeString(s) + `"`
}

// dbmlWriter renders DbmlExporter's layout. With esc set (the importers'
// clean mode) names, strings and expressions are escaped so the output
// always lexes; without it the text is byte-identical to @dbml/core.
type dbmlWriter struct{ esc bool }

var plainType = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*(\([A-Za-z0-9_, ]*\))?$`)

func (w *dbmlWriter) q(s string) string {
	if w.esc {
		return `"` + escapeString(s) + `"`
	}
	return rawdq(s)
}

func (w *dbmlWriter) lit(s string) string {
	if w.esc {
		return "'" + escapeString(s) + "'"
	}
	return "'" + s + "'"
}

// expr renders a backtick expression. DBML has no escape for a backtick
// inside one; in escaped mode backticks become double quotes (the clean
// importers already rewrite quoted identifiers before this).
func (w *dbmlWriter) expr(s string) string {
	if w.esc {
		s = strings.ReplaceAll(s, "`", `"`)
	}
	return "`" + s + "`"
}

func exportDBML(c *core, opts CoreDBMLOptions) string {
	w := &dbmlWriter{esc: opts.Escape}
	var elems []string
	for _, s := range c.schemas {
		if len(s.enums) > 0 {
			elems = append(elems, w.dbmlEnums(s))
		}
		if len(s.tables) > 0 {
			elems = append(elems, w.dbmlTables(s))
		}
		if len(s.groups) > 0 {
			elems = append(elems, w.dbmlGroups(s))
		}
		if len(s.refs) > 0 {
			elems = append(elems, w.dbmlRefs(s))
		}
		if len(s.deps) > 0 {
			elems = append(elems, w.dbmlDeps(s))
		}
	}
	if len(c.notes) > 0 {
		var parts []string
		for _, n := range c.notes {
			set := ""
			if n.Color != nil && *n.Color != "" {
				set = " [color: " + *n.Color + "]"
			}
			name := n.Name
			if w.esc {
				name = addDoubleQuoteIfNeeded(name)
			}
			parts = append(parts, "Note "+name+set+" {\n  "+escapeNote(n.Content)+"\n}\n")
		}
		elems = append(elems, strings.Join(parts, "\n"))
	}
	if !opts.OmitRecords && len(c.db.Records) > 0 {
		elems = append(elems, dbmlRecords(c.db.Records))
	}
	if m := w.dbmlMetadata(c); m != "" {
		elems = append(elems, m)
	}
	return strings.Join(elems, "\n")
}

func (w *dbmlWriter) schemaQ(s *cSchema) string {
	if s.print() {
		return w.q(s.name) + "."
	}
	return ""
}

func (w *dbmlWriter) dbmlEnums(s *cSchema) string {
	var parts []string
	for _, e := range s.enums {
		var vals []string
		for i, v := range e.values {
			line := "  " + w.q(v)
			if e.valueNotes[i] != "" {
				line += " [note: " + escapeNote(e.valueNotes[i]) + "]"
			}
			vals = append(vals, line)
		}
		parts = append(parts, "Enum "+w.schemaQ(s)+w.q(e.name)+" {\n"+strings.Join(vals, "\n")+"\n}\n")
	}
	return strings.Join(parts, "\n")
}

func (w *dbmlWriter) dbmlTables(s *cSchema) string {
	var parts []string
	for _, t := range s.tables {
		name := w.q(t.name)
		if s.print() {
			name = w.q(s.name) + "." + w.q(t.name)
		}
		alias := ""
		if t.src != nil && t.src.Alias != nil && *t.src.Alias != "" {
			alias = " as " + addDoubleQuoteIfNeeded(*t.src.Alias)
		}
		settings := ""
		if t.color != nil && *t.color != "" {
			settings = " [headerColor: " + *t.color + "]"
		}
		var b strings.Builder
		b.WriteString("Table " + name + alias + settings + " {\n")
		for _, f := range t.fields {
			b.WriteString("  " + w.dbmlField(f) + "\n")
		}
		if len(t.checks) > 0 {
			b.WriteString("\n  Checks {\n")
			for _, ch := range t.checks {
				line := w.expr(ch.expression)
				if ch.name != nil && *ch.name != "" {
					line += " [name: " + w.lit(*ch.name) + "]"
				}
				b.WriteString("    " + line + "\n")
			}
			b.WriteString("  }\n")
		}
		if len(t.indexes) > 0 {
			b.WriteString("\n  Indexes {\n")
			for _, ix := range t.indexes {
				b.WriteString("    " + w.dbmlIndex(ix) + "\n")
			}
			b.WriteString("  }\n")
		}
		if t.note != "" {
			b.WriteString("  Note: " + escapeNote(t.note) + "\n")
		}
		b.WriteString("}\n")
		parts = append(parts, b.String())
	}
	return strings.Join(parts, "\n")
}

func (w *dbmlWriter) dbmlField(f *cField) string {
	schema := ""
	if f.typ.schemaName != nil && *f.typ.schemaName != "" && *f.typ.schemaName != defaultSchema {
		if reWhiteSpace.MatchString(*f.typ.schemaName) {
			schema = w.q(*f.typ.schemaName) + "."
		} else {
			schema = *f.typ.schemaName + "."
			if w.esc && !isValidIdentifier(*f.typ.schemaName) {
				schema = w.q(*f.typ.schemaName) + "."
			}
		}
	}
	tn := f.typ.typeName
	if reWhiteSpace.MatchString(tn) || reSquareBracket.MatchString(tn) || (w.esc && !plainType.MatchString(tn)) {
		tn = w.q(tn)
	}
	line := w.q(f.name) + " " + schema + tn
	var cons []string
	if f.unique {
		cons = append(cons, "unique")
	}
	if f.pk {
		cons = append(cons, "pk")
	}
	if f.notNull {
		cons = append(cons, "not null")
	}
	if f.increment {
		cons = append(cons, "increment")
	}
	for _, ch := range f.checks {
		cons = append(cons, "check: "+w.expr(ch.expression))
	}
	if d := f.dbdefault; d != nil && !isNullDefault(d) {
		v := jsString(d.Value)
		switch d.Type {
		case "boolean", "number":
			cons = append(cons, "default: "+v)
		case "string":
			if w.esc {
				cons = append(cons, "default: "+w.lit(v))
				break
			}
			q := "'"
			if strings.Contains(v, "\n") {
				q = "'''"
			}
			cons = append(cons, "default: "+q+v+q)
		default:
			cons = append(cons, "default: "+w.expr(v))
		}
	}
	if f.note != "" {
		cons = append(cons, "note: "+escapeNote(f.note))
	}
	if len(cons) > 0 {
		line += " [" + strings.Join(cons, ", ") + "]"
	}
	return line
}

func (w *dbmlWriter) dbmlIndex(ix *model.Index) string {
	col := func(c *model.IndexColumn) string {
		if c.Type == "expression" {
			return w.expr(c.Value)
		}
		if w.esc {
			return addDoubleQuoteIfNeeded(c.Value)
		}
		return c.Value
	}
	line := ""
	if len(ix.Columns) > 1 {
		var cs []string
		for _, c := range ix.Columns {
			cs = append(cs, col(c))
		}
		line = "(" + strings.Join(cs, ", ") + ")"
	} else if len(ix.Columns) == 1 {
		line = col(ix.Columns[0])
	}
	var set []string
	if truthy(ix.PK) {
		set = append(set, "pk")
	}
	if ix.Type != nil && *ix.Type != "" {
		set = append(set, "type: "+strings.ToLower(*ix.Type))
	}
	if truthy(ix.Unique) {
		set = append(set, "unique")
	}
	if ix.Name != nil && *ix.Name != "" {
		set = append(set, "name: "+w.q(*ix.Name))
	}
	if len(set) > 0 {
		line += " [" + strings.Join(set, ", ") + "]"
	}
	return line
}

func (w *dbmlWriter) dbmlGroups(s *cSchema) string {
	var parts []string
	for _, g := range s.groups {
		set := ""
		if g.g.Color != nil && *g.g.Color != "" {
			set = " [color: " + *g.g.Color + "]"
		}
		note := ""
		if g.g.Note != nil && g.g.Note.Value != "" {
			note = "  Note: " + escapeNote(g.g.Note.Value) + "\n"
		}
		var names strings.Builder
		for _, t := range g.tables {
			names.WriteString("  " + w.schemaQ(t.schema) + w.q(t.name) + "\n")
		}
		parts = append(parts, "TableGroup "+w.q(str(g.g.Name))+set+" {\n"+names.String()+note+"}\n")
	}
	return strings.Join(parts, "\n")
}

func (w *dbmlWriter) dbmlFieldNames(fs []*cField) string {
	var ns []string
	for _, f := range fs {
		ns = append(ns, w.q(f.name))
	}
	if len(fs) == 1 {
		return ns[0]
	}
	return "(" + strings.Join(ns, ", ") + ")"
}

// relationOp mirrors getRelationshipOp (unknown pairs fall back to "<>").
func relationOp(l, r string) string {
	for op, m := range relOps {
		if m[0] == l && m[1] == r {
			return op
		}
	}
	return "<>"
}

func (w *dbmlWriter) dbmlRefs(s *cSchema) string {
	var parts []string
	for _, r := range s.refs {
		left := 0
		for i, ep := range r.endpoints {
			if maxIsOne(ep.relation) {
				left = i
				break
			}
		}
		l, rt := r.endpoints[left], r.endpoints[1-left]
		line := "Ref"
		if r.name != nil && *r.name != "" {
			line += " " + w.schemaQ(r.schema) + w.q(*r.name)
		}
		line += ":"
		lt := l.fields[0].table
		line += w.schemaQ(lt.schema) + w.q(lt.name) + "." + w.dbmlFieldNames(l.fields) + " "
		line += relationOp(l.relation, rt.relation) + " "
		rtt := rt.fields[0].table
		line += w.schemaQ(rtt.schema) + w.q(rtt.name) + "." + w.dbmlFieldNames(rt.fields)
		var acts []string
		if r.onUpdate != nil && *r.onUpdate != "" {
			acts = append(acts, "update: "+strings.ToLower(*r.onUpdate))
		}
		if r.onDelete != nil && *r.onDelete != "" {
			acts = append(acts, "delete: "+strings.ToLower(*r.onDelete))
		}
		if r.inactive {
			acts = append(acts, "inactive")
		}
		if len(acts) > 0 {
			line += " [" + strings.Join(acts, ", ") + "]"
		}
		parts = append(parts, line+"\n")
	}
	return strings.Join(parts, "\n")
}

func (w *dbmlWriter) depEndpoint(ep model.DepEndpoint) string {
	s := ""
	if ep.SchemaName != nil && *ep.SchemaName != "" {
		s = w.q(*ep.SchemaName) + "."
	}
	s += w.q(ep.TableName)
	switch len(ep.FieldNames) {
	case 0:
		return s
	case 1:
		return s + "." + w.q(ep.FieldNames[0])
	}
	var fs []string
	for _, f := range ep.FieldNames {
		fs = append(fs, w.q(f))
	}
	return s + ".(" + strings.Join(fs, ", ") + ")"
}

func depCustom(v any) string {
	switch x := v.(type) {
	case string:
		return "'" + escapeString(x) + "'"
	case nil:
		return "null"
	}
	return jsString(v)
}

func (w *dbmlWriter) dbmlDeps(s *cSchema) string {
	var parts []string
	for _, d := range s.deps {
		hasNote := d.Note != nil && d.Note.Value != ""
		keys := make([]string, 0, len(d.Metadata))
		for k := range d.Metadata {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		if len(d.Edges) == 1 && !hasNote && len(keys) == 0 {
			parts = append(parts, "Dep: "+w.depEndpoint(d.Edges[0].Upstream)+" -> "+w.depEndpoint(d.Edges[0].Downstream)+"\n")
			continue
		}
		var b strings.Builder
		b.WriteString("Dep {\n")
		for _, e := range d.Edges {
			b.WriteString("  " + w.depEndpoint(e.Upstream) + " -> " + w.depEndpoint(e.Downstream) + "\n")
		}
		if hasNote || len(keys) > 0 {
			b.WriteString("\n")
			if hasNote {
				b.WriteString("  note: " + escapeNote(d.Note.Value) + "\n")
			}
			for _, k := range keys {
				b.WriteString("  " + k + ": " + depCustom(d.Metadata[k]) + "\n")
			}
		}
		b.WriteString("}\n")
		parts = append(parts, b.String())
	}
	return strings.Join(parts, "\n")
}

// FormatRecordValue mirrors @dbml/parse formatRecordValue.
func FormatRecordValue(v *model.RecordValue) string {
	if v == nil || v.Undefined || v.Value == nil {
		return "null"
	}
	if isBooleanType(v.Type) {
		if s, ok := v.Value.(string); ok && strings.ToLower(s) == "null" {
			return "null"
		}
	}
	raw := jsString(v.Value)
	if v.Type == "expression" {
		if s, ok := v.Value.(string); ok && strings.ToLower(strings.TrimSpace(s)) == "null" {
			return "null"
		}
		return "`" + raw + "`"
	}
	if isBooleanType(v.Type) {
		if b, ok := extractBoolean(v.Value); ok {
			if b {
				return "true"
			}
			return "false"
		}
		return "`" + raw + "`"
	}
	if isNumericType(v.Type) {
		if n, ok := extractNumeric(v.Value); ok {
			return jsNumber(n)
		}
		return "`" + raw + "`"
	}
	if isDateTimeType(v.Type) {
		if s, ok := v.Value.(string); ok {
			if d, ok := model.NormalizeDateTime(s); ok {
				return "'" + escapeString(d) + "'"
			}
		}
		return "`" + raw + "`"
	}
	return "'" + escapeString(raw) + "'"
}

func extractBoolean(v any) (bool, bool) {
	switch x := v.(type) {
	case bool:
		return x, true
	case float64:
		if x == 0 {
			return false, true
		}
		if x == 1 {
			return true, true
		}
	case string:
		switch strings.ToLower(x) {
		case "true", "yes", "y", "t", "1":
			return true, true
		case "false", "no", "n", "f", "0":
			return false, true
		}
	}
	return false, false
}

// extractNumeric mirrors tryExtractNumeric on a primitive (JS Number()).
func extractNumeric(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case bool:
		if x {
			return 1, true
		}
		return 0, true
	case string:
		return jsToNumber(x)
	}
	return 0, false
}

// jsToNumber mirrors JavaScript Number(string) (NaN -> false).
func jsToNumber(s string) (float64, bool) {
	t := strings.TrimSpace(s)
	if t == "" {
		return 0, true
	}
	lower := strings.ToLower(t)
	switch {
	case strings.HasPrefix(lower, "0x"), strings.HasPrefix(lower, "0o"), strings.HasPrefix(lower, "0b"):
		base := map[byte]int{'x': 16, 'o': 8, 'b': 2}[lower[1]]
		n, err := strconv.ParseUint(t[2:], base, 64)
		if err != nil {
			return 0, false
		}
		return float64(n), true
	}
	switch t {
	case "Infinity", "+Infinity":
		return 0, false
	case "-Infinity":
		return 0, false
	}
	if strings.ContainsAny(t, "_xXpP") || strings.HasPrefix(lower, "inf") || strings.HasPrefix(lower, "+inf") || strings.HasPrefix(lower, "-inf") || strings.Contains(lower, "nan") {
		return 0, false
	}
	f, err := strconv.ParseFloat(t, 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

func dbmlRecords(recs []*model.TableRecord) string {
	type group struct {
		schema *string
		table  string
		recs   []*model.TableRecord
	}
	var groups []*group
	idx := map[string]*group{}
	for _, r := range recs {
		k := str(r.SchemaName) + "." + r.TableName
		g := idx[k]
		if g == nil {
			g = &group{schema: r.SchemaName, table: r.TableName}
			idx[k] = g
			groups = append(groups, g)
		}
		g.recs = append(g.recs, r)
	}
	var parts []string
	for _, g := range groups {
		ref := addDoubleQuoteIfNeeded(g.table)
		if g.schema != nil && *g.schema != "" {
			ref = addDoubleQuoteIfNeeded(*g.schema) + "." + addDoubleQuoteIfNeeded(g.table)
		}
		var cols []string
		seen := map[string]bool{}
		for _, r := range g.recs {
			for _, c := range r.Columns {
				if !seen[c] {
					seen[c] = true
					cols = append(cols, c)
				}
			}
		}
		var qcols []string
		for _, c := range cols {
			qcols = append(qcols, addDoubleQuoteIfNeeded(c))
		}
		var rows []string
		example := false
		for _, r := range g.recs {
			if truthy(r.Example) {
				example = true
			}
			pos := map[string]int{}
			for i := len(r.Columns) - 1; i >= 0; i-- {
				pos[r.Columns[i]] = i
			}
			for _, row := range r.Values {
				var vals []string
				for _, c := range cols {
					i, ok := pos[c]
					if !ok || i >= len(row) {
						vals = append(vals, "null")
						continue
					}
					vals = append(vals, FormatRecordValue(row[i]))
				}
				rows = append(rows, "  "+strings.Join(vals, ", "))
			}
		}
		ex := ""
		if example {
			ex = " [example]"
		}
		parts = append(parts, "Records "+ref+"("+strings.Join(qcols, ", ")+")"+ex+" {\n"+strings.Join(rows, "\n")+"\n}\n")
	}
	return strings.Join(parts, "\n")
}

func (w *dbmlWriter) metadataBlock(header string, m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var lines []string
	for _, k := range keys {
		lines = append(lines, "  "+k+": "+escapeNote(m[k]))
	}
	return "Metadata " + header + " {\n" + strings.Join(lines, "\n") + "\n}\n"
}

func (w *dbmlWriter) dbmlMetadata(c *core) string {
	var blocks []string
	for _, s := range c.schemas {
		for _, t := range s.tables {
			prefix := w.schemaQ(t.schema)
			if t.src != nil && len(t.src.Metadata) > 0 {
				blocks = append(blocks, w.metadataBlock("Table "+prefix+w.q(t.name), t.src.Metadata))
			}
			for _, f := range t.fields {
				if len(f.metadata) > 0 {
					blocks = append(blocks, w.metadataBlock("Column "+prefix+w.q(t.name)+"."+w.q(f.name), f.metadata))
				}
			}
		}
		for _, g := range s.groups {
			if len(g.g.Metadata) > 0 {
				blocks = append(blocks, w.metadataBlock("TableGroup "+w.q(str(g.g.Name)), g.g.Metadata))
			}
		}
	}
	for _, n := range c.notes {
		if len(n.Metadata) > 0 {
			blocks = append(blocks, w.metadataBlock("Note "+n.Name, n.Metadata))
		}
	}
	return strings.Join(blocks, "\n")
}

// rawdq wraps s in double quotes without escaping, exactly like
// DbmlExporter (which does not escape names).
func rawdq(s string) string { return `"` + s + `"` }
