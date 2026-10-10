// Package importer converts SQL DDL (MySQL/MariaDB, PostgreSQL, SQL
// Server) and Rails schema.rb into a DBML model.
//
// The SQL importers mirror @dbml/core 10.3.1's importer.import(src, format)
// (the ANTLR AST generators in packages/dbml-core/src/parse/ANTLR) over a
// hand-written streaming tokenizer and statement parser: only CREATE TABLE
// / CREATE INDEX / ALTER TABLE ... ADD constraint / COMMENT ON / CREATE
// TYPE ... AS ENUM / INSERT statements are read, everything else is
// skipped.
//
// Two modes:
//   - ModeClean (the default, for product use): SQL string literals and
//     quoted identifiers are unescaped, quoted identifiers inside CHECK /
//     DEFAULT / index expressions are rewritten so the DBML compiles,
//     records for unknown tables and duplicate ref names are dropped, and
//     schema.rb keeps column options (null/default/limit/precision/scale,
//     comments, primary keys, indexes, implicit id type).
//   - ModeParity: byte-identical to @dbml/core (CoreDBML), quirks
//     included; used by the golden suite.
//
// INSERT statements become Records only with ImportOptions.IncludeRecords,
// capped at MaxRecordRows (rows beyond the cap are dropped and
// Result.RecordsTruncated is set); skipped statements are never tokenised
// into memory.
//
// Match rates against @dbml/core goldens are reported in AGENTS.md.
package importer

import (
	"errors"
	"fmt"
	"math"
	"strings"

	dbml "github.com/aidenappl/go-dbml"
	"github.com/aidenappl/go-dbml/diag"
	"github.com/aidenappl/go-dbml/emit"
	"github.com/aidenappl/go-dbml/export"
	"github.com/aidenappl/go-dbml/model"
)

// Format is an import source format.
type Format string

// Supported formats.
const (
	MySQL    Format = "mysql" // also MariaDB
	Postgres Format = "postgres"
	MSSQL    Format = "mssql"
	SchemaRB Format = "schemarb"
)

// Formats lists the supported formats.
var Formats = []Format{MySQL, Postgres, MSSQL, SchemaRB}

// Mode selects clean (usable) or parity (@dbml/core-identical) output.
type Mode int

const (
	// ModeClean unescapes literals/identifiers and fixes output that would
	// not compile. The default.
	ModeClean Mode = iota
	// ModeParity reproduces @dbml/core 10.3.1's importer.import exactly.
	ModeParity
)

// MaxInputBytes bounds the input size.
const MaxInputBytes = 16 << 20

// ErrInputTooLarge is returned for inputs larger than MaxInputBytes.
var ErrInputTooLarge = errors.New("import: input exceeds 16 MiB")

// ErrUnsupported is wrapped by errors for input the reference importer
// rejects (go-dbml mirrors the failure rather than guessing).
var ErrUnsupported = errors.New("unsupported SQL")

func errUnsupported(what string) error { return fmt.Errorf("import: %w: %s", ErrUnsupported, what) }

// InternalError reports a recovered panic inside the importer (a bug).
type InternalError struct{ Value any }

func (e *InternalError) Error() string { return fmt.Sprintf("import: internal error: %v", e.Value) }

// ImportOptions tunes ImportWith. The zero value is clean mode without
// records and default compile limits.
type ImportOptions struct {
	Mode Mode
	// IncludeRecords turns INSERT ... VALUES into Records blocks.
	IncludeRecords bool
	// MaxRecordRows caps imported rows (0 = model.DefaultMaxRecordRows).
	MaxRecordRows int
	// Model bounds compiling the generated DBML (zero = defaults).
	Model model.Options
	// OmitProject leaves out the `Project { database_type: '…' }` block
	// clean mode adds to record the source dialect (export.SQLWith's
	// TranslateTypes reads it). Set it when the DBML is merged into a
	// document that has its own Project.
	OmitProject bool
	// MaxInputBytes lowers the input cap for this import (0 = MaxInputBytes,
	// 16 MiB; larger values are clamped to it). Callers can set it per
	// format, e.g. smaller for schema.rb than for SQL dumps.
	MaxInputBytes int
}

// Result is an import result.
type Result struct {
	// Database is the compiled model, nil when the generated DBML does not
	// compile (then Partial is go-dbml's best-effort model).
	Database *model.Database
	Partial  *model.Database
	// DBML is the generated text: canonical go-dbml formatting (emit) in
	// clean mode when it compiles, @dbml/core's layout otherwise.
	DBML        string
	Diagnostics diag.List
	// RecordsTruncated is set when rows beyond MaxRecordRows were dropped.
	RecordsTruncated bool
	// Warnings lists clean-mode notes: expressions dropped because they
	// could not be normalised, and string literals whose backslashes keep
	// the source dialect's escaping.
	Warnings []string
}

// CompileError reports that the DBML generated for an import does not
// compile with go-dbml (common in parity mode: @dbml/core's own output
// emits records for tables never created, unescaped quotes, duplicate ref
// names).
type CompileError struct {
	DBML        string
	Diagnostics diag.List
}

func (e *CompileError) Error() string {
	return fmt.Sprintf("import: generated DBML does not compile: %v", e.Diagnostics)
}

// ParseFormat maps a name (mysql, mariadb, postgres, postgresql, pg,
// mssql, sqlserver, schemarb, rails) to a Format.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "mysql", "mariadb":
		return MySQL, nil
	case "postgres", "postgresql", "pg":
		return Postgres, nil
	case "mssql", "sqlserver", "tsql":
		return MSSQL, nil
	case "schemarb", "schema.rb", "rails":
		return SchemaRB, nil
	}
	return "", fmt.Errorf("unsupported import format %q", s)
}

// config is the per-import state shared with the dialect importers.
type config struct {
	clean     bool
	records   bool
	maxRows   int
	rows      int
	truncated bool
	maxInput  int
}

// skipStatement drops INSERT statements before they are stored when
// records are off or the row cap is reached.
func (c *config) skipStatement(t token) bool {
	if !t.is("INSERT") {
		return false
	}
	if c.records && c.rows >= c.maxRows {
		c.truncated = true
		return true
	}
	return !c.records
}

// maxStatementTokens bounds the tokens stored for one statement when
// records are on (a single INSERT can hold a whole dump).
func (c *config) maxStatementTokens() int {
	if !c.records || c.maxRows >= math.MaxInt32/512 {
		return 0
	}
	return c.maxRows*256 + 4096
}

// addRow reserves one record row; false once the cap is reached.
func (c *config) addRow() bool {
	if c.rows >= c.maxRows {
		c.truncated = true
		return false
	}
	c.rows++
	return true
}

func rawImport(src string, f Format, cfg *config) (raw *rawDB, err error) {
	limit := MaxInputBytes
	if cfg.maxInput > 0 && cfg.maxInput < limit {
		limit = cfg.maxInput
	}
	if len(src) > limit {
		return nil, ErrInputTooLarge
	}
	defer func() {
		if r := recover(); r != nil {
			raw, err = nil, &InternalError{Value: r}
		}
	}()
	switch f {
	case MySQL:
		return importMySQL(src, cfg)
	case Postgres:
		return importPostgres(src, cfg)
	case MSSQL:
		return importMSSQL(src, cfg)
	case SchemaRB:
		return importSchemaRB(src, cfg)
	}
	return nil, fmt.Errorf("unsupported import format %q", f)
}

// CoreDBML returns exactly what @dbml/core 10.3.1's importer.import(src,
// format) returns (parity mode, records included, no row cap).
func CoreDBML(src string, f Format) (string, error) {
	cfg := &config{records: true, maxRows: math.MaxInt}
	raw, err := rawImport(src, f, cfg)
	if err != nil {
		return "", err
	}
	return export.CoreDBML(raw.toModel(), export.CoreDBMLOptions{OverrideCardinality: true})
}

// ImportWith imports src. Errors are input errors (ErrInputTooLarge,
// ErrUnsupported, *InternalError); a generated DBML that does not compile
// is reported through Result.Database == nil and Result.Diagnostics, not
// as an error.
func ImportWith(src string, f Format, opts ImportOptions) (*Result, error) {
	cfg := &config{clean: opts.Mode == ModeClean, records: opts.IncludeRecords, maxRows: opts.MaxRecordRows, maxInput: opts.MaxInputBytes}
	if cfg.maxRows <= 0 {
		cfg.maxRows = model.DefaultMaxRecordRows
	}
	// never raise a stricter compiler limit
	mo := opts.Model
	if mo.MaxRecordRows > 0 && mo.MaxRecordRows < cfg.maxRows {
		cfg.maxRows = mo.MaxRecordRows
	}
	if mo.MaxRecordRows <= 0 {
		mo.MaxRecordRows = cfg.maxRows
	}
	raw, err := rawImport(src, f, cfg)
	if err != nil {
		return nil, err
	}
	db := raw.toModel()
	var warns []string
	if cfg.clean {
		warns = normalizeModel(db, sourceDialect(f))
		if f == MSSQL {
			unqualifyDefaultSchema(db, "dbo")
		}
		cleanModel(db)
	}
	text, err := export.CoreDBML(db, export.CoreDBMLOptions{OverrideCardinality: true, Escape: cfg.clean})
	if err != nil {
		return nil, err
	}
	if cfg.clean && !opts.OmitProject {
		text = "Project {\n  database_type: '" + databaseType(f, src) + "'\n}\n\n" + text
	}
	r := dbml.ParseFileWith(dbml.DefaultFilepath, text, mo)
	res := &Result{Database: r.Database, Partial: r.Partial, DBML: text, Diagnostics: r.Diagnostics, RecordsTruncated: cfg.truncated, Warnings: warns}
	if cfg.clean && r.Database != nil {
		if pretty, err := emit.DBML(r.Database); err == nil {
			res.DBML = pretty
		}
	}
	return res, nil
}

// Import imports src in clean mode without records and returns the
// compiled model.
//
// Contract: on success the model is complete. If the generated DBML does
// not compile, Import returns go-dbml's best-effort model (Result.Partial,
// possibly with dangling references) together with a *CompileError
// carrying the text and diagnostics — callers must check the error before
// trusting the model. Other errors (ErrInputTooLarge, ErrUnsupported,
// *InternalError) return a nil model.
func Import(src string, f Format) (*model.Database, error) {
	r, err := ImportWith(src, f, ImportOptions{})
	if err != nil {
		return nil, err
	}
	if r.Database == nil {
		return r.Partial, &CompileError{DBML: r.DBML, Diagnostics: r.Diagnostics}
	}
	return r.Database, nil
}

// sourceDialect is the SQL dialect of f's expressions.
func sourceDialect(f Format) dialect {
	switch f {
	case MySQL:
		return dMySQL
	case Postgres:
		return dPostgres
	case MSSQL:
		return dMSSQL
	}
	return dRails
}

// databaseType is the Project database_type recording f's dialect.
func databaseType(f Format, src string) string {
	switch f {
	case MySQL:
		return export.DatabaseType(export.MySQL)
	case Postgres:
		return export.DatabaseType(export.Postgres)
	case MSSQL:
		return export.DatabaseType(export.MSSQL)
	}
	if rbIsMySQL(src) {
		return export.DatabaseType(export.MySQL)
	}
	return export.DatabaseType(export.Postgres)
}
