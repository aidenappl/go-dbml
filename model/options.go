package model

// Options bounds the resources a compile may use. Zero fields take the
// corresponding default, so Options{} behaves like DefaultOptions().
type Options struct {
	// MaxBytes is the largest source file parsed, in bytes. Larger files are
	// skipped with an INPUT_TOO_LARGE error.
	MaxBytes int
	// MaxFiles bounds the files a project may load through imports.
	MaxFiles int
	// MaxDiagnostics caps the returned diagnostics; the rest are summarised
	// by one trailing TOO_MANY_DIAGNOSTICS entry.
	MaxDiagnostics int
	// MaxDepth bounds syntactic nesting (see parser.DefaultMaxDepth).
	MaxDepth int
	// MaxRecordRows bounds the rows of all Records blocks together; the
	// excess is dropped with a RECORDS_TOO_LARGE error. Each cell costs
	// roughly 1 KB through lexing, parsing and the model, so keep this low
	// where memory is scarce (the WASM build defaults to 2,000).
	MaxRecordRows int
	// MaxRecordBytes bounds the source bytes of all Records blocks
	// together; a block that would exceed it is skipped with a
	// RECORDS_TOO_LARGE error.
	MaxRecordBytes int
	// MaxTokens bounds the non-trivia tokens of one file. The effective
	// limit is also at most len(source)/2 + 65,536: hand-written DBML has
	// well under one token per 2 bytes, while inputs such as "a,a,a,..."
	// (one token per byte) exist only to exhaust memory. Past it the rest
	// of the file is not parsed and one TOO_MANY_TOKENS error is reported.
	MaxTokens int
	// MaxElements bounds the elements of one block, list, tuple or record
	// row (columns of a table, values of an enum, columns of an index,
	// cells of a row, ...); the rest is skipped with one TOO_MANY_ELEMENTS
	// error.
	MaxElements int
	// Cancel aborts the compile when closed (e.g. ctx.Done()). It is polled
	// by the lexer, parser, validator, binder and interpreter; a canceled
	// compile returns only a CANCELED diagnostic (Database and Partial
	// nil).
	Cancel <-chan struct{}
}

// Default limits: generous for real schemas, safe for untrusted input.
const (
	DefaultMaxBytes       = 16 << 20 // 16 MiB per file
	DefaultMaxFiles       = 256
	DefaultMaxDiagnostics = 1000
	DefaultMaxDepth       = 256
	DefaultMaxRecordRows  = 10000
	DefaultMaxRecordBytes = 4 << 20 // 4 MiB
	DefaultMaxTokens      = 2_000_000
	DefaultMaxElements    = 10_000
)

// DefaultOptions returns the default limits.
func DefaultOptions() Options {
	return Options{MaxBytes: DefaultMaxBytes, MaxFiles: DefaultMaxFiles, MaxDiagnostics: DefaultMaxDiagnostics, MaxDepth: DefaultMaxDepth,
		MaxRecordRows: DefaultMaxRecordRows, MaxRecordBytes: DefaultMaxRecordBytes, MaxTokens: DefaultMaxTokens, MaxElements: DefaultMaxElements}
}

func (o Options) withDefaults() Options {
	d := DefaultOptions()
	if o.MaxBytes <= 0 {
		o.MaxBytes = d.MaxBytes
	}
	if o.MaxFiles <= 0 {
		o.MaxFiles = d.MaxFiles
	}
	if o.MaxDiagnostics <= 0 {
		o.MaxDiagnostics = d.MaxDiagnostics
	}
	if o.MaxDepth <= 0 {
		o.MaxDepth = d.MaxDepth
	}
	if o.MaxRecordRows <= 0 {
		o.MaxRecordRows = d.MaxRecordRows
	}
	if o.MaxRecordBytes <= 0 {
		o.MaxRecordBytes = d.MaxRecordBytes
	}
	if o.MaxTokens <= 0 {
		o.MaxTokens = d.MaxTokens
	}
	if o.MaxElements <= 0 {
		o.MaxElements = d.MaxElements
	}
	return o
}

// tokenLimit is the effective MaxTokens for a source of n bytes.
func (o Options) tokenLimit(n int) int {
	if d := n/2 + 65536; d < o.MaxTokens {
		return d
	}
	return o.MaxTokens
}
