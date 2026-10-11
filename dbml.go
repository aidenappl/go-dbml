// Package dbml is the entry point of go-dbml, a pure-Go implementation of
// DBML (Database Markup Language, https://dbml.dbdiagram.io).
//
// The common case is a single call:
//
//	db, diags := dbml.Parse(src)
//	if diags.HasErrors() { ... }   // db is nil when there are errors
//
// The returned model.Database serialises (encoding/json) to the same shape
// as the per-file database JSON produced by @dbml/parse / @dbml/core. Every
// diagnostic carries a 1-based line/column range in UTF-16 units, so
// diags.Markers() can be handed straight to Monaco's setModelMarkers.
//
// Multi-file projects using `use` / `reuse` imports are compiled with
// ParseProject.
package dbml

import (
	"context"
	"io"
	"io/fs"
	"path"
	"strings"

	"github.com/aidenappl/go-dbml/diag"
	"github.com/aidenappl/go-dbml/model"
	"github.com/aidenappl/go-dbml/parser"
)

// Version is the go-dbml compiler version. Store it (with
// model.JSONSchemaVersion) alongside persisted model JSON.
const Version = "0.1.1"

// DefaultFilepath is the virtual path used for single-source parsing.
const DefaultFilepath = "/main.dbml"

// Options bounds resource usage (file size, file count, diagnostics,
// nesting depth). The zero value means DefaultOptions().
type Options = model.Options

// DefaultOptions returns the default limits: 16 MiB per file, 256 files,
// 1000 diagnostics, nesting depth 256, 10,000 record rows / 4 MiB of
// Records source.
func DefaultOptions() Options { return model.DefaultOptions() }

// Parse compiles a single DBML source with default options. The database
// is nil if any error-severity diagnostic was produced; warnings and infos
// never block it. Use ParseFile for the best-effort Result.Partial model.
func Parse(src string) (*model.Database, diag.List) {
	r := model.Compile(src, DefaultFilepath)
	return r.Database, r.Diagnostics
}

// ParseFile compiles src under an explicit (slash-separated, absolute) file
// path, which is recorded in every token and diagnostic.
func ParseFile(filepath, src string) model.Result {
	return model.Compile(src, filepath)
}

// ParseFileWith is ParseFile with explicit limits.
func ParseFileWith(filepath, src string, opts Options) model.Result {
	return model.CompileWithOptions(src, filepath, opts)
}

// ParseProject compiles entry and the files it imports. files maps
// ParseFileContext is ParseFileWith that stops when ctx is done: the
// result then carries only a CANCELED diagnostic (Database and Partial
// nil). The lexer, parser, validator, binder and interpreter poll ctx.
func ParseFileContext(ctx context.Context, filepath, src string, opts Options) model.Result {
	opts.Cancel = ctx.Done()
	if ctx.Err() != nil {
		return model.Result{Diagnostics: diag.List{diag.New(diag.CANCELED, "Compilation was canceled", filepath,
			diag.Range{Start: diag.Pos{Line: 1, Column: 1}, End: diag.Pos{Line: 1, Column: 1}})}}
	}
	return ParseFileWith(filepath, src, opts)
}

// absolute slash-separated paths (e.g. "/schema/users.dbml") to sources.
func ParseProject(entry string, files map[string]string) model.Result {
	return ParseProjectWith(entry, files, Options{})
}

// ParseProjectWith is ParseProject with explicit limits.
func ParseProjectWith(entry string, files map[string]string, opts Options) model.Result {
	return model.CompileProjectWithOptions(entry, func(p string) (string, bool) {
		s, ok := files[p]
		return s, ok
	}, opts)
}

// ParseFS compiles entry (a path inside fsys, without a leading slash) and
// the files it imports from fsys.
//
// Security: imports are resolved inside fsys, but os.DirFS follows symbolic
// links out of its root. For untrusted directories pass a root-confined
// file system instead, e.g.
//
//	root, err := os.OpenRoot(dir) // Go 1.24+
//	if err != nil { ... }
//	defer root.Close()
//	res := dbml.ParseFS(root.FS(), "main.dbml")
func ParseFS(fsys fs.FS, entry string) model.Result {
	return ParseFSWith(fsys, entry, Options{})
}

// ParseFSWith is ParseFS with explicit limits. Files larger than
// opts.MaxBytes are not read into memory.
func ParseFSWith(fsys fs.FS, entry string, opts Options) model.Result {
	opts = withDefaults(opts)
	entry = "/" + strings.TrimPrefix(path.Clean(entry), "/")
	return model.CompileProjectWithOptions(entry, func(p string) (string, bool) {
		f, err := fsys.Open(strings.TrimPrefix(p, "/"))
		if err != nil {
			return "", false
		}
		defer f.Close()
		// Read at most MaxBytes+1: enough for the compiler to report
		// INPUT_TOO_LARGE without loading an arbitrarily large file.
		b, err := io.ReadAll(io.LimitReader(f, int64(opts.MaxBytes)+1))
		if err != nil {
			return "", false
		}
		return string(b), true
	}, opts)
}

func withDefaults(o Options) Options {
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

// ParseSyntax only lexes and parses src, returning the concrete syntax tree
// and syntax diagnostics (no semantic analysis).
func ParseSyntax(src string) parser.Result {
	return parser.Parse(src, DefaultFilepath)
}
