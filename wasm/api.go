// Command wasm is go-dbml's browser entry point (GOOS=js GOARCH=wasm). It
// installs a global `goDbml` object; every function takes strings and
// returns a JSON string and never throws (a panic inside a js.FuncOf
// callback would kill the Go instance, so every call recovers).
//
//	goDbml.parse(src, optionsJSON?)              -> ParseResponse
//	goDbml.exportSQL(src, dialect, optionsJSON?) -> ExportResponse
//	goDbml.importSQL(sql, format, optionsJSON?)  -> ImportResponse
//	goDbml.emit(modelJSON)                       -> EmitResponse
//	goDbml.version                               -> "<go-dbml version>"
//
// Every response has "ok" and, when ok is false, "error"; "diagnostics"
// (Monaco IMarkerData, 1-based, UTF-16 columns) is always an array where
// present. "partial" (the best-effort model) appears only when ok is false.
// The JSON contract is documented in AGENTS.md §6; the handlers live in this
// file (no build tag) so they are unit-tested natively.
package main

import (
	"encoding/json"
	"fmt"
	"runtime"

	dbml "github.com/aidenappl/go-dbml"
	"github.com/aidenappl/go-dbml/diag"
	"github.com/aidenappl/go-dbml/emit"
	"github.com/aidenappl/go-dbml/export"
	"github.com/aidenappl/go-dbml/importer"
	"github.com/aidenappl/go-dbml/model"
)

// Browser defaults: lower than the native ones because the instance's
// linear memory never shrinks. A 1 MiB schema of small tables (the
// worst legitimate shape) peaks at ~300 MB of linear memory; 2 MiB at
// ~600 MB (measured, AGENTS.md §6).
const (
	browserMaxBytes       = 1 << 20 // 1 MiB of DBML/SQL (~300 MB linear memory worst case)
	browserMaxRecordRows  = 2000
	browserMaxRecordBytes = 1 << 20
	browserMaxDiagnostics = 500
	maxModelJSON          = 8 << 20
)

// Options is the optional options JSON (all fields optional; zero = browser
// default). Values above the native defaults are clamped to them.
type Options struct {
	MaxBytes       int    `json:"maxBytes,omitempty"`
	MaxRecordRows  int    `json:"maxRecordRows,omitempty"`
	MaxRecordBytes int    `json:"maxRecordBytes,omitempty"`
	MaxDiagnostics int    `json:"maxDiagnostics,omitempty"`
	MaxDepth       int    `json:"maxDepth,omitempty"`
	MaxTokens      int    `json:"maxTokens,omitempty"`
	MaxElements    int    `json:"maxElements,omitempty"`
	Mode           string `json:"mode,omitempty"` // importSQL: "clean" (default) | "parity"
	IncludeRecords bool   `json:"includeRecords,omitempty"`
	// Tokens=false returns the compact model (parse, importSQL): no source
	// tokens; tables, columns, refs, enums, groups, partials and notes
	// carry "line" instead.
	Tokens *bool `json:"tokens,omitempty"`
	// TranslateTypes (exportSQL): translate types/defaults from the
	// dialect recorded in the model's Project database_type to the target.
	// Default: on when a different source dialect is recorded.
	TranslateTypes *bool `json:"translateTypes,omitempty"`
	// SourceDialect (exportSQL) overrides the recorded source dialect.
	SourceDialect string `json:"sourceDialect,omitempty"`
	// OmitProject (importSQL): no `Project { database_type }` block in the
	// generated DBML.
	OmitProject bool `json:"omitProject,omitempty"`
	// MaxInputBytes (importSQL) caps the SQL / schema.rb input per call,
	// e.g. per format: {"mysql": 1048576, "schemarb": 262144}; a format
	// missing from the object uses maxBytes (default 1 MiB).
	MaxInputBytes map[string]int `json:"maxInputBytes,omitempty"`
}

func (o Options) compact() bool { return o.Tokens != nil && !*o.Tokens }

// ParseResponse answers parse.
type ParseResponse struct {
	OK            bool                `json:"ok"`
	Error         string              `json:"error,omitempty"`
	Database      *model.Database     `json:"database,omitempty"`
	Partial       *model.Database     `json:"partial,omitempty"`
	Diagnostics   []diag.MonacoMarker `json:"diagnostics"`
	Version       string              `json:"version"`
	SchemaVersion int                 `json:"schemaVersion"`
}

// ExportResponse answers exportSQL.
type ExportResponse struct {
	OK          bool                `json:"ok"`
	Error       string              `json:"error,omitempty"`
	SQL         string              `json:"sql,omitempty"`
	Diagnostics []diag.MonacoMarker `json:"diagnostics"`
	// Translated reports that types were translated from SourceDialect.
	Translated    bool   `json:"translated,omitempty"`
	SourceDialect string `json:"sourceDialect,omitempty"`
	// Warnings are translation notes (types kept as written, dropped
	// fragments).
	Warnings []string `json:"warnings,omitempty"`
}

// ImportResponse answers importSQL. dbml is set whenever the import ran
// (also when ok is false because the generated DBML does not compile).
type ImportResponse struct {
	OK               bool                `json:"ok"`
	Error            string              `json:"error,omitempty"`
	DBML             string              `json:"dbml,omitempty"`
	Database         *model.Database     `json:"database,omitempty"`
	Partial          *model.Database     `json:"partial,omitempty"`
	Diagnostics      []diag.MonacoMarker `json:"diagnostics"`
	RecordsTruncated bool                `json:"recordsTruncated,omitempty"`
	// Warnings are clean-mode import notes (dropped or dialect-specific
	// expressions).
	Warnings []string `json:"warnings,omitempty"`
}

// EmitResponse answers emit.
type EmitResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
	DBML  string `json:"dbml,omitempty"`
}

func markers(l diag.List) []diag.MonacoMarker {
	m := l.Markers()
	if m == nil {
		m = []diag.MonacoMarker{}
	}
	return m
}

func clamp(v, def, max int) int {
	if v <= 0 {
		return def
	}
	if v > max {
		return max
	}
	return v
}

// parseOptions decodes optionsJSON ("" = defaults).
func parseOptions(js string) (Options, model.Options, error) {
	var o Options
	if js != "" {
		if len(js) > 1<<16 {
			return o, model.Options{}, fmt.Errorf("options JSON too large")
		}
		if err := json.Unmarshal([]byte(js), &o); err != nil {
			return o, model.Options{}, fmt.Errorf("invalid options JSON: %v", err)
		}
	}
	d := model.DefaultOptions()
	mo := model.Options{
		MaxBytes:       clamp(o.MaxBytes, browserMaxBytes, d.MaxBytes),
		MaxRecordRows:  clamp(o.MaxRecordRows, browserMaxRecordRows, d.MaxRecordRows),
		MaxRecordBytes: clamp(o.MaxRecordBytes, browserMaxRecordBytes, d.MaxRecordBytes),
		MaxDiagnostics: clamp(o.MaxDiagnostics, browserMaxDiagnostics, d.MaxDiagnostics),
		MaxDepth:       clamp(o.MaxDepth, d.MaxDepth, d.MaxDepth),
		MaxTokens:      clamp(o.MaxTokens, d.MaxTokens, d.MaxTokens),
		MaxElements:    clamp(o.MaxElements, d.MaxElements, d.MaxElements),
		MaxFiles:       1,
	}
	return o, mo, nil
}

// guard runs f, converting a panic into an {"ok":false,"error":...} reply.
func guard(f func() any) (out string) {
	defer func() {
		if r := recover(); r != nil {
			out = encode(struct {
				OK    bool   `json:"ok"`
				Error string `json:"error"`
			}{false, fmt.Sprintf("internal error: %v", r)})
		}
	}()
	return encode(f())
}

func encode(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		b, _ = json.Marshal(struct {
			OK    bool   `json:"ok"`
			Error string `json:"error"`
		}{false, "encode: " + err.Error()})
	}
	return string(b)
}

func apiParse(src, optionsJSON string) string {
	return guard(func() any {
		o, mo, err := parseOptions(optionsJSON)
		if err != nil {
			return ParseResponse{Error: err.Error(), Diagnostics: []diag.MonacoMarker{}, Version: dbml.Version, SchemaVersion: model.JSONSchemaVersion}
		}
		r := dbml.ParseFileWith(dbml.DefaultFilepath, src, mo)
		r.Program = nil // the syntax tree is most of the memory: release it before encoding
		resp := ParseResponse{OK: r.Database != nil, Database: r.Database, Diagnostics: markers(r.Diagnostics),
			Version: dbml.Version, SchemaVersion: model.JSONSchemaVersion}
		if !resp.OK {
			resp.Error = "DBML has errors"
			resp.Partial = r.Partial
		}
		if o.compact() {
			model.StripTokens(resp.Database)
			model.StripTokens(resp.Partial)
		}
		releaseIfLarge(len(src))
		return resp
	})
}

func apiExportSQL(src, dialect, optionsJSON string) string {
	return guard(func() any {
		resp := ExportResponse{Diagnostics: []diag.MonacoMarker{}}
		d, err := export.ParseDialect(dialect)
		if err != nil {
			resp.Error = err.Error()
			return resp
		}
		o, mo, err := parseOptions(optionsJSON)
		if err != nil {
			resp.Error = err.Error()
			return resp
		}
		r := dbml.ParseFileWith(dbml.DefaultFilepath, src, mo)
		r.Program = nil
		resp.Diagnostics = markers(r.Diagnostics)
		if r.Database == nil {
			resp.Error = "DBML has errors"
			return resp
		}
		releaseIfLarge(len(src))
		eo := export.Options{}
		if o.SourceDialect != "" {
			if eo.SourceDialect, err = export.ParseDialect(o.SourceDialect); err != nil {
				resp.Error = err.Error()
				return resp
			}
		}
		src := eo.SourceDialect
		if src == "" {
			src = export.SourceDialectOf(r.Database)
		}
		eo.TranslateTypes = src != "" && src != d
		if o.TranslateTypes != nil {
			eo.TranslateTypes = *o.TranslateTypes && src != "" && src != d
		}
		res, err := export.SQLWith(r.Database, d, eo)
		if err != nil {
			resp.Error = err.Error()
			return resp
		}
		resp.OK, resp.SQL, resp.Warnings = true, res.SQL, res.Warnings
		if eo.TranslateTypes {
			resp.Translated, resp.SourceDialect = true, string(src)
		}
		return resp
	})
}

func apiImportSQL(src, format, optionsJSON string) string {
	return guard(func() any {
		resp := ImportResponse{Diagnostics: []diag.MonacoMarker{}}
		f, err := importer.ParseFormat(format)
		if err != nil {
			resp.Error = err.Error()
			return resp
		}
		o, mo, err := parseOptions(optionsJSON)
		if err != nil {
			resp.Error = err.Error()
			return resp
		}
		limit := mo.MaxBytes
		if v, ok := o.MaxInputBytes[string(f)]; ok && v > 0 {
			limit = min(v, importer.MaxInputBytes)
		}
		if len(src) > limit {
			resp.Error = fmt.Sprintf("input exceeds %d bytes", limit)
			return resp
		}
		opts := importer.ImportOptions{IncludeRecords: o.IncludeRecords, MaxRecordRows: mo.MaxRecordRows, Model: mo, OmitProject: o.OmitProject, MaxInputBytes: limit}
		opts.Model.MaxBytes = importer.MaxInputBytes // the generated DBML may be larger than the SQL
		if o.Mode == "parity" {
			opts.Mode = importer.ModeParity
		} else if o.Mode != "" && o.Mode != "clean" {
			resp.Error = "unknown mode " + o.Mode
			return resp
		}
		r, err := importer.ImportWith(src, f, opts)
		if err != nil {
			resp.Error = err.Error()
			return resp
		}
		resp.DBML, resp.Diagnostics, resp.RecordsTruncated, resp.Warnings = r.DBML, markers(r.Diagnostics), r.RecordsTruncated, r.Warnings
		if o.compact() {
			model.StripTokens(r.Database)
			model.StripTokens(r.Partial)
		}
		releaseIfLarge(len(src))
		if r.Database == nil {
			resp.Error = "generated DBML does not compile"
			resp.Partial = r.Partial
			return resp
		}
		resp.OK, resp.Database = true, r.Database
		return resp
	})
}

// releaseIfLarge collects the compile's garbage before the response is
// encoded, so the (never shrinking) WASM heap does not grow to hold both.
func releaseIfLarge(n int) {
	if n > 256<<10 {
		runtime.GC()
	}
}

func apiEmit(modelJSON string) string {
	return guard(func() any {
		if len(modelJSON) > maxModelJSON {
			return EmitResponse{Error: "model JSON too large"}
		}
		var db model.Database
		if err := json.Unmarshal([]byte(modelJSON), &db); err != nil {
			return EmitResponse{Error: "invalid model JSON: " + err.Error()}
		}
		out, err := emit.DBML(&db)
		if err != nil {
			return EmitResponse{Error: err.Error()}
		}
		return EmitResponse{OK: true, DBML: out}
	})
}
