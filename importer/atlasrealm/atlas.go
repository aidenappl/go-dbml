//go:build !js

// Package atlasrealm converts an inspected Atlas realm
// (ariga.io/atlas/sql/schema) into a compiled DBML model. It is the only
// go-dbml package with a third-party dependency and is excluded from
// GOOS=js builds; the importer package itself stays stdlib-only.
package atlasrealm

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"ariga.io/atlas/sql/schema"

	dbml "github.com/aidenappl/go-dbml"
	"github.com/aidenappl/go-dbml/emit"
	"github.com/aidenappl/go-dbml/importer"
	"github.com/aidenappl/go-dbml/model"
)

// Options tunes FromRealm / Import.
type Options struct {
	// DefaultSchema is the schema whose tables are emitted without a schema
	// prefix. Empty means: "public", or the only schema of the realm.
	DefaultSchema string
	// Dialect is the inspected database's dialect (importer.MySQL for MySQL
	// and MariaDB, importer.Postgres, importer.MSSQL). It decides how raw
	// expressions are normalised and is recorded as the Project
	// database_type. Empty: detected from Atlas's dialect attributes
	// (sql/mysql, sql/postgres) or, failing that, from backtick-quoted
	// expressions (MySQL).
	Dialect importer.Format
}

// Result is Import's output.
type Result struct {
	// Database is the compiled model.
	Database *model.Database
	// DBML is the canonical DBML text of Database.
	DBML string
	// Warnings lists what was left out: foreign keys to tables outside the
	// inspected schemas, expressions that cannot be kept, dropped parts.
	Warnings []string
}

// FromRealm converts an inspected Atlas realm (ariga.io/atlas/sql/schema,
// e.g. from mysql/postgres/sqlite Driver.InspectRealm) into a compiled DBML
// model: tables, columns (raw database types, NOT NULL, defaults,
// auto-increment/identity/serial, comments), primary keys, indexes,
// checks, enums and foreign keys. Foreign keys become refs the way the SQL
// importers produce them (the FK side "many and optional", the referenced
// side required when every FK column is NOT NULL). Views are skipped.
//
// FromRealm is native-only (excluded from GOOS=js builds).
func FromRealm(r *schema.Realm, opts Options) (*model.Database, error) {
	res, err := Import(r, opts)
	if res == nil {
		return nil, err
	}
	return res.Database, err
}

// Import is FromRealm returning the DBML text and the warnings too.
//
//   - Raw expressions (CHECKs, expression defaults, generated columns,
//     expression index parts) are normalised per dialect like the clean
//     SQL importer's (importer.NormalizeExpr): MySQL/MariaDB report them
//     with `backtick` identifiers, which DBML cannot hold.
//   - Inline enums (MySQL/MariaDB report every one as "enum") are named
//     <table>_<column>_enum, with a numeric suffix on a clash; named enum
//     types (PostgreSQL) keep their name and are shared.
//   - A foreign key to a table outside the inspected schemas (or to
//     columns it does not have) is dropped with a warning, so the result
//     always compiles.
func Import(r *schema.Realm, opts Options) (*Result, error) {
	if r == nil {
		return nil, fmt.Errorf("import: nil realm")
	}
	dialect := opts.Dialect
	if dialect == "" {
		dialect = detectDialect(r)
	}
	var warns []string
	warn := func(format string, a ...any) { warns = append(warns, fmt.Sprintf(format, a...)) }
	norm := func(e, where string) (string, bool) {
		f := dialect
		if f == "" {
			f = importer.Postgres // standard SQL quoting
		}
		out, w, ok := importer.NormalizeExpr(e, f)
		if !ok {
			warn("%s: expression %q left out (%s)", where, e, w)
			return "", false
		}
		if w != "" {
			warn("%s: %s", where, w)
		}
		return out, true
	}
	def := opts.DefaultSchema
	if def == "" {
		def = "public"
		if len(r.Schemas) == 1 {
			def = r.Schemas[0].Name
		}
	}
	schemaOf := func(s *schema.Schema) *string {
		if s == nil || s.Name == "" || s.Name == def {
			return nil
		}
		return strp(s.Name)
	}
	db := &model.Database{}
	enums := map[string]string{} // schema.name -> values (to detect clashes)
	for _, s := range r.Schemas {
		for _, t := range s.Tables {
			mt := &model.Table{Name: t.Name, SchemaName: schemaOf(s)}
			if c := comment(t.Attrs); c != "" {
				mt.Note = &model.NoteValue{Value: c}
			}
			pkCols := map[string]bool{}
			if t.PrimaryKey != nil && len(t.PrimaryKey.Parts) == 1 && t.PrimaryKey.Parts[0].C != nil {
				pkCols[t.PrimaryKey.Parts[0].C.Name] = true
			}
			for _, c := range t.Columns {
				col := &model.Column{Name: c.Name}
				if et, ok := enumOf(c); ok {
					en := &model.Enum{Name: et.T, SchemaName: mt.SchemaName}
					if et.Schema != nil {
						en.SchemaName = schemaOf(et.Schema)
					}
					values := strings.Join(et.Values, "\x00")
					inline := en.Name == "" || strings.EqualFold(en.Name, "enum") || dialect == importer.MySQL
					if inline {
						// MySQL/MariaDB inline ENUM: one enum per column
						base := t.Name + "_" + c.Name + "_enum"
						en.Name = base
						for n := 2; enums[str(en.SchemaName)+"."+en.Name] != ""; n++ {
							en.Name = base + "_" + strconv.Itoa(n)
						}
					}
					key := str(en.SchemaName) + "." + en.Name
					switch prev := enums[key]; {
					case prev == "":
						enums[key] = values + "\x01"
						for _, v := range et.Values {
							en.Values = append(en.Values, &model.EnumValue{Name: v})
						}
						db.Enums = append(db.Enums, en)
					case prev != values+"\x01":
						warn("column %s.%s: enum type %s seen with different values; the first definition is kept", t.Name, c.Name, en.Name)
					}
					col.Type = model.ColumnType{TypeName: en.Name, SchemaName: en.SchemaName}
				} else {
					col.Type = columnType(typeString(c.Type))
				}
				if pkCols[c.Name] {
					col.PK = boolp(true)
				}
				if c.Type != nil && !c.Type.Null && !pkCols[c.Name] {
					col.NotNull = boolp(true)
				}
				if isIncrement(c) {
					col.Increment = boolp(true)
				}
				col.DBDefault = defaultOf(c.Default)
				where := "column " + t.Name + "." + c.Name
				if d := col.DBDefault; d != nil && d.Type == "expression" {
					if e, ok := norm(fmt.Sprint(d.Value), "default of "+where); ok {
						d.Value = e
					} else {
						col.DBDefault = nil
					}
				}
				for _, a := range c.Attrs {
					if ck, ok := a.(*schema.Check); ok && ck.Expr != "" {
						if e, ok := norm(ck.Expr, "check on "+where); ok {
							mc := &model.Check{Expression: e}
							if ck.Name != "" {
								mc.Name = strp(ck.Name)
							}
							col.Checks = append(col.Checks, mc)
						}
					}
					if g, ok := a.(*schema.GeneratedExpr); ok && g.Expr != "" {
						if e, ok := norm(g.Expr, "generated "+where); ok {
							col.Metadata = map[string]string{model.MetaGenerated: e}
							if st := strings.ToUpper(g.Type); st == "STORED" || st == "PERSISTENT" {
								col.Metadata[model.MetaGeneratedStored] = "true"
							}
							col.DBDefault = nil
						}
					}
				}
				if cm := comment(c.Attrs); cm != "" {
					col.Note = &model.NoteValue{Value: cm}
				}
				mt.Fields = append(mt.Fields, col)
			}
			if t.PrimaryKey != nil && len(t.PrimaryKey.Parts) > 1 {
				ix := indexOf(t.PrimaryKey)
				ix.PK, ix.Unique, ix.Name = boolp(true), nil, nil
				mt.Indexes = append(mt.Indexes, ix)
			}
			for _, ix := range t.Indexes {
				mi := indexOf(ix)
				cols := mi.Columns[:0]
				for _, ic := range mi.Columns {
					if ic.Type == "expression" {
						e, ok := norm(ic.Value, "index "+str(mi.Name)+" on "+t.Name)
						if !ok {
							continue
						}
						ic.Value = e
					}
					cols = append(cols, ic)
				}
				mi.Columns = cols
				if len(cols) == 0 {
					warn("index %s on %s left out (no usable part)", str(mi.Name), t.Name)
					continue
				}
				mt.Indexes = append(mt.Indexes, mi)
			}
			for _, ck := range t.Checks() {
				e, ok := norm(ck.Expr, "check "+ck.Name+" on "+t.Name)
				if !ok {
					continue
				}
				mc := &model.Check{Expression: e}
				if ck.Name != "" {
					mc.Name = strp(ck.Name)
				}
				mt.Checks = append(mt.Checks, mc)
			}
			db.Tables = append(db.Tables, mt)
		}
	}
	inRealm := map[*schema.Table]bool{}
	byName := map[string]*schema.Table{}
	for _, s := range r.Schemas {
		for _, t := range s.Tables {
			inRealm[t] = true
			byName[s.Name+"."+t.Name] = t
		}
	}
	hasCols := func(t *schema.Table, cs []*schema.Column) bool {
		for _, c := range cs {
			found := false
			for _, tc := range t.Columns {
				found = found || tc.Name == c.Name
			}
			if !found {
				return false
			}
		}
		return true
	}
	for _, s := range r.Schemas {
		for _, t := range s.Tables {
			for _, fk := range t.ForeignKeys {
				if fk.RefTable == nil || len(fk.Columns) == 0 {
					continue
				}
				refSchema := s.Name
				if fk.RefTable.Schema != nil && fk.RefTable.Schema.Name != "" {
					refSchema = fk.RefTable.Schema.Name
				}
				target := fk.RefTable
				if !inRealm[target] {
					target = byName[refSchema+"."+fk.RefTable.Name]
				}
				if target == nil {
					warn("foreign key %s on %s.%s left out: it references %s.%s, outside the inspected schemas", fk.Symbol, s.Name, t.Name, refSchema, fk.RefTable.Name)
					continue
				}
				if !hasCols(target, fk.RefColumns) || len(fk.RefColumns) != len(fk.Columns) {
					warn("foreign key %s on %s.%s left out: the referenced columns of %s.%s were not inspected", fk.Symbol, s.Name, t.Name, refSchema, fk.RefTable.Name)
					continue
				}
				required := true
				var cols, refCols []string
				for _, c := range fk.Columns {
					cols = append(cols, c.Name)
					if c.Type == nil || c.Type.Null {
						required = false
					}
				}
				for _, c := range fk.RefColumns {
					refCols = append(refCols, c.Name)
				}
				rel := "0..1"
				if required {
					rel = "1"
				}
				ref := &model.Ref{}
				if fk.Symbol != "" {
					ref.Name = strp(fk.Symbol)
				}
				if fk.OnDelete != "" && fk.OnDelete != schema.NoAction {
					ref.OnDelete = strp(strings.ToLower(string(fk.OnDelete)))
				}
				if fk.OnUpdate != "" && fk.OnUpdate != schema.NoAction {
					ref.OnUpdate = strp(strings.ToLower(string(fk.OnUpdate)))
				}
				ref.Endpoints[0] = model.RefEndpoint{TableName: target.Name, SchemaName: schemaOf(target.Schema), FieldNames: refCols, Relation: rel}
				ref.Endpoints[1] = model.RefEndpoint{TableName: t.Name, SchemaName: schemaOf(s), FieldNames: cols, Relation: "0..*"}
				db.Refs = append(db.Refs, ref)
			}
		}
	}
	if dt := databaseType(dialect); dt != "" {
		db.Project = &model.Project{Extra: map[string]string{"database_type": dt}}
	}
	text, err := emit.DBML(db)
	if err != nil {
		return &Result{Warnings: warns}, err
	}
	res := dbml.ParseFile(dbml.DefaultFilepath, text)
	out := &Result{Database: res.Database, DBML: text, Warnings: warns}
	if res.Database == nil {
		out.Database = res.Partial
		return out, &importer.CompileError{DBML: text, Diagnostics: res.Diagnostics}
	}
	return out, nil
}

func databaseType(f importer.Format) string {
	switch f {
	case importer.MySQL:
		return "MySQL"
	case importer.Postgres:
		return "PostgreSQL"
	case importer.MSSQL:
		return "SQL Server"
	}
	return ""
}

// detectDialect guesses the dialect from Atlas's dialect attribute types
// (ariga.io/atlas/sql/mysql, …/postgres) and, failing that, from
// backtick-quoted expressions (MySQL / MariaDB).
func detectDialect(r *schema.Realm) importer.Format {
	backtick := false
	var found importer.Format
	pkg := func(v any) {
		if v == nil || found != "" {
			return
		}
		t := reflect.TypeOf(v)
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		switch p := t.PkgPath(); {
		case strings.HasSuffix(p, "/sql/mysql"):
			found = importer.MySQL
		case strings.HasSuffix(p, "/sql/postgres"):
			found = importer.Postgres
		}
	}
	for _, s := range r.Schemas {
		for _, a := range s.Attrs {
			pkg(a)
		}
		for _, t := range s.Tables {
			for _, a := range t.Attrs {
				pkg(a)
				if ck, ok := a.(*schema.Check); ok && strings.Contains(ck.Expr, "`") {
					backtick = true
				}
			}
			for _, c := range t.Columns {
				for _, a := range c.Attrs {
					pkg(a)
					if g, ok := a.(*schema.GeneratedExpr); ok && strings.Contains(g.Expr, "`") {
						backtick = true
					}
				}
				if c.Type != nil {
					pkg(c.Type.Type)
				}
				if x, ok := c.Default.(*schema.RawExpr); ok && strings.Contains(x.X, "`") {
					backtick = true
				}
			}
		}
	}
	if found == "" && backtick {
		found = importer.MySQL
	}
	return found
}

func str(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

var typeArgsRe = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_ ]*)\(([A-Za-z0-9_, ]*)\)$`)

// columnType splits "varchar(255)" into the name + args shape DBML
// parsing produces (TypeName keeps the suffix, Args the arguments).
func columnType(raw string) model.ColumnType {
	m := typeArgsRe.FindStringSubmatch(raw)
	if m == nil || strings.Contains(m[1], " ") {
		return model.ColumnType{TypeName: raw}
	}
	parts := strings.Split(m[2], ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	args := strings.Join(parts, ",")
	return model.ColumnType{TypeName: m[1] + "(" + args + ")", Args: &args}
}

func comment(attrs []schema.Attr) string {
	for _, a := range attrs {
		if c, ok := a.(*schema.Comment); ok {
			return c.Text
		}
	}
	return ""
}

func enumOf(c *schema.Column) (*schema.EnumType, bool) {
	if c.Type == nil {
		return nil, false
	}
	e, ok := c.Type.Type.(*schema.EnumType)
	return e, ok
}

// typeName returns the T field of an Atlas type (every type has one).
func typeName(t schema.Type) string {
	v := reflect.ValueOf(t)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		if f := v.FieldByName("T"); f.IsValid() && f.Kind() == reflect.String {
			return f.String()
		}
	}
	return ""
}

func typeString(ct *schema.ColumnType) string {
	if ct == nil {
		return "unknown"
	}
	if ct.Raw != "" {
		return ct.Raw
	}
	switch t := ct.Type.(type) {
	case *schema.StringType:
		if t.Size > 0 {
			return t.T + "(" + strconv.Itoa(t.Size) + ")"
		}
		return t.T
	case *schema.BinaryType:
		if t.Size != nil {
			return t.T + "(" + strconv.Itoa(*t.Size) + ")"
		}
		return t.T
	case *schema.DecimalType:
		if t.Precision > 0 {
			return t.T + "(" + strconv.Itoa(t.Precision) + "," + strconv.Itoa(t.Scale) + ")"
		}
		return t.T
	case *schema.TimeType:
		if t.Precision != nil {
			return t.T + "(" + strconv.Itoa(*t.Precision) + ")"
		}
		return t.T
	}
	if n := typeName(ct.Type); n != "" {
		return n
	}
	return "unknown"
}

// isIncrement detects dialect attributes without importing the dialect
// packages: mysql.AutoIncrement, postgres.Identity, postgres.SerialType.
func isIncrement(c *schema.Column) bool {
	for _, a := range c.Attrs {
		n := reflect.TypeOf(a).String()
		if strings.HasSuffix(n, ".AutoIncrement") || strings.HasSuffix(n, ".Identity") {
			return true
		}
	}
	if c.Type != nil && c.Type.Type != nil && strings.HasSuffix(reflect.TypeOf(c.Type.Type).String(), ".SerialType") {
		return true
	}
	return false
}

func defaultOf(x schema.Expr) *model.Default {
	if x == nil {
		return nil
	}
	switch v := schema.UnderlyingExpr(x).(type) {
	case *schema.Literal:
		s := v.V
		low := strings.ToLower(s)
		switch {
		case low == "null":
			return nil
		case len(s) >= 2 && (s[0] == '\'' || s[0] == '"') && s[len(s)-1] == s[0]:
			inner := s[1 : len(s)-1]
			return &model.Default{Type: "string", Value: strings.ReplaceAll(inner, string(s[0])+string(s[0]), string(s[0]))}
		case low == "true" || low == "false":
			return &model.Default{Type: "boolean", Value: low}
		}
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return &model.Default{Type: "number", Value: f}
		}
		return &model.Default{Type: "expression", Value: s}
	case *schema.RawExpr:
		return &model.Default{Type: "expression", Value: v.X}
	}
	return nil
}

func indexOf(ix *schema.Index) *model.Index {
	mi := &model.Index{}
	if ix.Name != "" {
		mi.Name = strp(ix.Name)
	}
	if ix.Unique {
		mi.Unique = boolp(true)
	}
	for _, a := range ix.Attrs {
		if strings.HasSuffix(reflect.TypeOf(a).String(), ".IndexType") {
			if t := attrT(a); t != "" && !strings.EqualFold(t, "btree") {
				mi.Type = strp(strings.ToLower(t))
			}
		}
	}
	for _, p := range ix.Parts {
		switch {
		case p.C != nil:
			mi.Columns = append(mi.Columns, &model.IndexColumn{Type: "column", Value: p.C.Name})
		case p.X != nil:
			if raw, ok := p.X.(*schema.RawExpr); ok {
				mi.Columns = append(mi.Columns, &model.IndexColumn{Type: "expression", Value: raw.X})
			}
		}
	}
	return mi
}

func attrT(a schema.Attr) string {
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		if f := v.FieldByName("T"); f.IsValid() && f.Kind() == reflect.String {
			return f.String()
		}
	}
	return ""
}

func strp(s string) *string { return &s }

func boolp(b bool) *bool { return &b }
