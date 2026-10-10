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

// Options tunes FromRealm.
type Options struct {
	// DefaultSchema is the schema whose tables are emitted without a schema
	// prefix. Empty means: "public", or the only schema of the realm.
	DefaultSchema string
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
	if r == nil {
		return nil, fmt.Errorf("import: nil realm")
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
	enums := map[string]bool{}
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
					if en.Name == "" {
						en.Name = t.Name + "_" + c.Name + "_enum"
					}
					key := str(en.SchemaName) + "." + en.Name
					if !enums[key] {
						enums[key] = true
						for _, v := range et.Values {
							en.Values = append(en.Values, &model.EnumValue{Name: v})
						}
						db.Enums = append(db.Enums, en)
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
				mt.Indexes = append(mt.Indexes, indexOf(ix))
			}
			for _, ck := range t.Checks() {
				mc := &model.Check{Expression: ck.Expr}
				if ck.Name != "" {
					mc.Name = strp(ck.Name)
				}
				mt.Checks = append(mt.Checks, mc)
			}
			db.Tables = append(db.Tables, mt)
		}
	}
	for _, s := range r.Schemas {
		for _, t := range s.Tables {
			for _, fk := range t.ForeignKeys {
				if fk.RefTable == nil || len(fk.Columns) == 0 {
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
				ref.Endpoints[0] = model.RefEndpoint{TableName: fk.RefTable.Name, SchemaName: schemaOf(fk.RefTable.Schema), FieldNames: refCols, Relation: rel}
				ref.Endpoints[1] = model.RefEndpoint{TableName: t.Name, SchemaName: schemaOf(s), FieldNames: cols, Relation: "0..*"}
				db.Refs = append(db.Refs, ref)
			}
		}
	}
	text, err := emit.DBML(db)
	if err != nil {
		return nil, err
	}
	res := dbml.ParseFile(dbml.DefaultFilepath, text)
	if res.Database == nil {
		return res.Partial, &importer.CompileError{DBML: text, Diagnostics: res.Diagnostics}
	}
	return res.Database, nil
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
