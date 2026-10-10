// Portions of this file are derived from holistics/dbml packages/dbml-parse/src/core/types/schemaJson/types.ts
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

// Package model holds go-dbml's resolved DBML model and the semantic
// analysis (validation, binding and interpretation) that produces it.
//
// The JSON encoding of Database matches the per-file "Database" schema JSON
// emitted by @dbml/parse (the format @dbml/core's Parser consumes and its
// "json" import format accepts). Field names, nesting and the presence of
// optional fields follow the reference implementation so the Holistics
// fixture suite can be compared directly. Intentional deviations are listed
// in AGENTS.md (e.g. Token.Filepath is a plain string).
package model

import (
	"encoding/json"
	"sort"

	"github.com/aidenappl/go-dbml/diag"
)

// JSONSchemaVersion identifies the JSON shape of Database. It tracks the
// @dbml/parse database JSON (10.3.1) and is bumped whenever go-dbml changes
// the shape, so stored model_json can be re-compiled when it is stale.
// Version 2 (go-dbml Task 0.3 round 3) adds Column.generated; tokens that
// are absent are omitted instead of null.
const JSONSchemaVersion = 2

// Position is a 1-based source position; Offset/Column are UTF-16 based.
type Position = diag.Pos

// Token is the source span of a model element.
type Token struct {
	Start    Position `json:"start"`
	End      Position `json:"end"`
	Filepath string   `json:"filepath"`
}

// NoteValue is a note attached to an element.
type NoteValue struct {
	Value string `json:"value"`
	Token *Token `json:"token,omitempty"`
}

// Database is the interpreted content of one DBML file (plus anything it imports).
type Database struct {
	Schemas       []Schema        `json:"schemas"`
	Tables        []*Table        `json:"tables"`
	Notes         []*StickyNote   `json:"notes"`
	Refs          []*Ref          `json:"refs"`
	Deps          []*Dep          `json:"deps"`
	Enums         []*Enum         `json:"enums"`
	TableGroups   []*TableGroup   `json:"tableGroups"`
	Aliases       []*Alias        `json:"aliases"`
	Project       *Project        `json:"project,omitempty"`
	TablePartials []*TablePartial `json:"tablePartials"`
	Records       []*TableRecord  `json:"records"`
	Externals     Externals       `json:"externals"`
	DiagramViews  []*DiagramView  `json:"diagramViews"`
	Token         *Token          `json:"token,omitempty"`
}

// Schema is a placeholder for the per-file "schemas" list, which
// @dbml/parse always emits empty (schemas are implied by element names).
type Schema struct{}

// Externals lists elements imported from other files.
type Externals struct {
	Tables        []*ElementRef `json:"tables"`
	Enums         []*ElementRef `json:"enums"`
	TableGroups   []*ElementRef `json:"tableGroups"`
	TablePartials []*ElementRef `json:"tablePartials"`
	Notes         []*ElementRef `json:"notes"`
}

// ElementRef identifies an imported element.
type ElementRef struct {
	Name         string        `json:"name"`
	SchemaName   *string       `json:"schemaName"`
	Filepath     string        `json:"filepath"`
	VisibleNames []VisibleName `json:"visibleNames"`
}

// VisibleName is a local name under which an imported element is reachable.
type VisibleName struct {
	SchemaName *string `json:"schemaName"`
	Name       string  `json:"name"`
}

// Table is a table definition.
type Table struct {
	Name        string            `json:"name"`
	SchemaName  *string           `json:"schemaName"`
	Alias       *string           `json:"alias"`
	Fields      []*Column         `json:"fields"`
	Checks      []*Check          `json:"checks"`
	Partials    []*PartialInject  `json:"partials"`
	Token       *Token            `json:"token,omitempty"`
	Indexes     []*Index          `json:"indexes"`
	HeaderColor *string           `json:"headerColor,omitempty"`
	Note        *NoteValue        `json:"note,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	// Line is the 1-based start line, set only by StripTokens (compact JSON).
	Line int `json:"line,omitempty"`
}

// ColumnType is a column's data type.
type ColumnType struct {
	SchemaName *string `json:"schemaName"`
	TypeName   string  `json:"type_name"`
	Args       *string `json:"args"`
}

// Default is a column default value. Value is a float64 for numbers and a
// string otherwise.
type Default struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}

// Column is a table or table-partial column.
type Column struct {
	Name       string            `json:"name"`
	Type       ColumnType        `json:"type"`
	Token      *Token            `json:"token,omitempty"`
	InlineRefs []*InlineRef      `json:"inline_refs"`
	InlineDeps []*InlineDep      `json:"inline_deps"`
	Checks     []*Check          `json:"checks"`
	PK         *bool             `json:"pk,omitempty"`
	DBDefault  *Default          `json:"dbdefault,omitempty"`
	Increment  *bool             `json:"increment,omitempty"`
	Unique     *bool             `json:"unique,omitempty"`
	NotNull    *bool             `json:"not_null,omitempty"`
	Note       *NoteValue        `json:"note,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	// Generated is set when the column carries the go-dbml custom
	// properties [generated: '<expr>'] (and optionally
	// generated_stored: 'true'): a generated / computed column. Additive
	// (go-dbml only; @dbml/parse keeps them as plain metadata).
	Generated *GeneratedColumn `json:"generated,omitempty"`
	// Line is the 1-based start line, set only by StripTokens (compact JSON).
	Line int `json:"line,omitempty"`
}

// Custom property keys describing a generated (computed) column.
const (
	MetaGenerated       = "generated"
	MetaGeneratedStored = "generated_stored"
)

// GeneratedColumn is a generated / computed column: Expression is SQL
// (standard SQL in DBML; exporters validate it per dialect), Stored
// selects STORED / PERSISTED over VIRTUAL.
type GeneratedColumn struct {
	Expression string `json:"expression"`
	Stored     bool   `json:"stored"`
}

// GeneratedFromMetadata derives a column's Generated value from its
// custom properties (nil when there is no non-empty "generated" key).
func GeneratedFromMetadata(m map[string]string) *GeneratedColumn {
	e := m[MetaGenerated]
	if e == "" {
		return nil
	}
	st := m[MetaGeneratedStored]
	return &GeneratedColumn{Expression: e, Stored: st == "true" || st == "stored" || st == "1"}
}

// IndexColumn is one column or expression of an index.
type IndexColumn struct {
	Value string `json:"value"`
	Type  string `json:"type"` // "column" | "expression"
	Token *Token `json:"token,omitempty"`
}

// Index is a table index.
type Index struct {
	Columns []*IndexColumn `json:"columns"`
	Token   *Token         `json:"token,omitempty"`
	Unique  *bool          `json:"unique,omitempty"`
	PK      *bool          `json:"pk,omitempty"`
	Name    *string        `json:"name,omitempty"`
	Note    *NoteValue     `json:"note,omitempty"`
	Type    *string        `json:"type,omitempty"`
}

// Check is a check constraint.
type Check struct {
	Token      *Token  `json:"token,omitempty"`
	Expression string  `json:"expression"`
	Name       *string `json:"name,omitempty"`
}

// InlineRef is a `[ref: > t.c]` column setting.
type InlineRef struct {
	SchemaName *string  `json:"schemaName"`
	TableName  string   `json:"tableName"`
	FieldNames []string `json:"fieldNames"`
	Relation   string   `json:"relation"`
	Token      *Token   `json:"token,omitempty"`
}

// Ref is a relationship.
type Ref struct {
	SchemaName *string        `json:"schemaName"`
	Name       *string        `json:"name"`
	Endpoints  [2]RefEndpoint `json:"endpoints"`
	Color      *string        `json:"color,omitempty"`
	OnDelete   *string        `json:"onDelete,omitempty"`
	OnUpdate   *string        `json:"onUpdate,omitempty"`
	Inactive   *bool          `json:"inactive,omitempty"`
	Token      *Token         `json:"token,omitempty"`
	// Line is the 1-based start line, set only by StripTokens (compact JSON).
	Line int `json:"line,omitempty"`
}

// RefEndpoint is one side of a relationship. Relation is a cardinality:
// "1", "*", "0..1" or "0..*".
type RefEndpoint struct {
	SchemaName *string  `json:"schemaName"`
	TableName  string   `json:"tableName"`
	FieldNames []string `json:"fieldNames"`
	Relation   string   `json:"relation"`
	Token      *Token   `json:"token,omitempty"`
}

// Dep is a data-lineage dependency.
type Dep struct {
	SchemaName *string        `json:"schemaName"`
	Name       *string        `json:"name"`
	Edges      []*DepEdge     `json:"edges"`
	Color      *string        `json:"color,omitempty"`
	Note       *NoteValue     `json:"note,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	Token      *Token         `json:"token,omitempty"`
}

// DepEdge is one upstream -> downstream edge.
type DepEdge struct {
	Upstream   DepEndpoint `json:"upstream"`
	Downstream DepEndpoint `json:"downstream"`
	Token      *Token      `json:"token,omitempty"`
}

// DepEndpoint is a table (FieldNames empty) or column(s).
type DepEndpoint struct {
	SchemaName *string  `json:"schemaName"`
	TableName  string   `json:"tableName"`
	FieldNames []string `json:"fieldNames"`
	Token      *Token   `json:"token,omitempty"`
}

// InlineDep is a `[dep: -> t.c]` column setting.
type InlineDep struct {
	SchemaName *string  `json:"schemaName"`
	TableName  string   `json:"tableName"`
	FieldNames []string `json:"fieldNames"`
	Direction  string   `json:"direction"`
	Token      *Token   `json:"token,omitempty"`
}

// Enum is an enum definition.
type Enum struct {
	Name       string       `json:"name"`
	SchemaName *string      `json:"schemaName"`
	Token      *Token       `json:"token,omitempty"`
	Values     []*EnumValue `json:"values"`
	// Line is the 1-based start line, set only by StripTokens (compact JSON).
	Line int `json:"line,omitempty"`
}

// EnumValue is one enum member.
type EnumValue struct {
	Name  string     `json:"name"`
	Token *Token     `json:"token,omitempty"`
	Note  *NoteValue `json:"note,omitempty"`
}

// TableGroup groups tables.
type TableGroup struct {
	Name       *string           `json:"name"`
	SchemaName *string           `json:"schemaName"`
	Tables     []*TableGroupItem `json:"tables"`
	Token      *Token            `json:"token,omitempty"`
	Color      *string           `json:"color,omitempty"`
	Note       *NoteValue        `json:"note,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	// Line is the 1-based start line, set only by StripTokens (compact JSON).
	Line int `json:"line,omitempty"`
}

// TableGroupItem names a grouped table.
type TableGroupItem struct {
	Name       string  `json:"name"`
	SchemaName *string `json:"schemaName"`
}

// Alias records `Table x as y`.
type Alias struct {
	Name  string     `json:"name"`
	Kind  string     `json:"kind"`
	Value AliasValue `json:"value"`
}

// AliasValue is the aliased table.
type AliasValue struct {
	TableName  string  `json:"tableName"`
	SchemaName *string `json:"schemaName"`
}

// TablePartial is a reusable set of fields, indexes and checks.
type TablePartial struct {
	Name        string     `json:"name"`
	Fields      []*Column  `json:"fields"`
	Token       *Token     `json:"token,omitempty"`
	Indexes     []*Index   `json:"indexes"`
	Checks      []*Check   `json:"checks"`
	HeaderColor *string    `json:"headerColor,omitempty"`
	Note        *NoteValue `json:"note,omitempty"`
	// Line is the 1-based start line, set only by StripTokens (compact JSON).
	Line int `json:"line,omitempty"`
}

// PartialInject is a `~partial` injection inside a table.
type PartialInject struct {
	Name  string `json:"name"`
	Order int    `json:"order"`
	Token *Token `json:"token,omitempty"`
}

// StickyNote is a top-level `Note name { ... }`.
type StickyNote struct {
	Name     string            `json:"name"`
	Content  string            `json:"content"`
	Token    *Token            `json:"token,omitempty"`
	Color    *string           `json:"color,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	// Line is the 1-based start line, set only by StripTokens (compact JSON).
	Line int `json:"line,omitempty"`
}

// RecordValue is one cell of a records block. Value is nil, a float64, a
// bool or a string; Type is the inferred value type (e.g. "integer",
// "string", "bool", "expression").
type RecordValue struct {
	Value any    `json:"value"`
	Type  string `json:"type"`
	Token *Token `json:"token,omitempty"`
	// Undefined marks a value @dbml/parse leaves undefined (e.g. an
	// `Enum.member` reference in a records row); "value" is then omitted
	// from the JSON.
	Undefined bool `json:"-"`
}

// UnmarshalJSON marks a value without a "value" key as Undefined.
func (v *RecordValue) UnmarshalJSON(b []byte) error {
	type plain RecordValue
	if err := json.Unmarshal(b, (*plain)(v)); err != nil {
		return err
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	_, has := raw["value"]
	v.Undefined = !has
	return nil
}

// MarshalJSON omits "value" for undefined values.
func (v RecordValue) MarshalJSON() ([]byte, error) {
	if v.Undefined {
		return json.Marshal(struct {
			Type  string `json:"type"`
			Token *Token `json:"token,omitempty"`
		}{v.Type, v.Token})
	}
	type plain RecordValue
	return json.Marshal(plain(v))
}

// TableRecord is a records block.
type TableRecord struct {
	SchemaName *string          `json:"schemaName"`
	TableName  string           `json:"tableName"`
	Columns    []string         `json:"columns"`
	Values     [][]*RecordValue `json:"values"`
	Example    *bool            `json:"example,omitempty"`
	Token      *Token           `json:"token,omitempty"`
}

// FilterConfig is the visible-entity filter of a DiagramView. For each
// category a nil slice means "hide all", an empty non-nil slice means "show
// all" and a non-empty slice lists the visible items.
type FilterConfig struct {
	Tables      []DiagramViewTable `json:"tables"`
	StickyNotes []NamedItem        `json:"stickyNotes"`
	TableGroups []NamedItem        `json:"tableGroups"`
	Schemas     []NamedItem        `json:"schemas"`
}

// DiagramViewTable names a table in a view.
type DiagramViewTable struct {
	Name       string `json:"name"`
	SchemaName string `json:"schemaName"`
}

// NamedItem is a `{name}` object.
type NamedItem struct {
	Name string `json:"name"`
}

// DiagramView is a named view over the diagram.
type DiagramView struct {
	Name            string       `json:"name"`
	SchemaName      *string      `json:"schemaName"`
	VisibleEntities FilterConfig `json:"visibleEntities"`
	Token           *Token       `json:"token,omitempty"`
}

// Project is the Project block. Extra holds free-form string fields such as
// database_type.
type Project struct {
	Name          *string           `json:"name"`
	Tables        []*Table          `json:"tables"`
	Refs          []*Ref            `json:"refs"`
	Enums         []*Enum           `json:"enums"`
	TableGroups   []*TableGroup     `json:"tableGroups"`
	TablePartials []*TablePartial   `json:"tablePartials"`
	Note          *NoteValue        `json:"note,omitempty"`
	Token         *Token            `json:"token,omitempty"`
	Extra         map[string]string `json:"-"`
}

// MarshalJSON flattens Extra into the object like the reference implementation.
func (p *Project) MarshalJSON() ([]byte, error) {
	type plain Project
	base, err := json.Marshal((*plain)(p))
	if err != nil || len(p.Extra) == 0 {
		return base, err
	}
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(p.Extra))
	for k := range p.Extra {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if _, reserved := m[k]; reserved {
			continue
		}
		v, _ := json.Marshal(p.Extra[k])
		m[k] = v
	}
	return json.Marshal(m)
}

// UnmarshalJSON restores Extra from keys that are not standard fields.
func (p *Project) UnmarshalJSON(b []byte) error {
	type plain Project
	if err := json.Unmarshal(b, (*plain)(p)); err != nil {
		return err
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	known := map[string]bool{"name": true, "tables": true, "refs": true, "enums": true, "tableGroups": true,
		"tablePartials": true, "note": true, "token": true}
	for k, v := range raw {
		if known[k] {
			continue
		}
		var str string
		if json.Unmarshal(v, &str) == nil {
			if p.Extra == nil {
				p.Extra = map[string]string{}
			}
			p.Extra[k] = str
		}
	}
	return nil
}

// MarshalJSON emits nil slices as null and empty slices as [] (the tri-state).
func (f FilterConfig) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"tables":      nilOrSlice(f.Tables),
		"stickyNotes": nilOrSlice(f.StickyNotes),
		"tableGroups": nilOrSlice(f.TableGroups),
		"schemas":     nilOrSlice(f.Schemas),
	}
	return json.Marshal(m)
}

func nilOrSlice[T any](s []T) any {
	if s == nil {
		return nil
	}
	return s
}

func strPtr(s string) *string { return &s }

func boolPtr(b bool) *bool { return &b }
