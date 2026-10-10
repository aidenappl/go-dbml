# go-dbml

A pure-Go DBML toolkit: compiler, pretty-printer, SQL exporters, SQL/schema.rb importers and a WASM build.

> **appleby.cloud** · Go library · `github.com/aidenappl/go-dbml`

---

## Overview

`go-dbml` parses [DBML](https://dbml.dbdiagram.io) v2 — tables, columns, indexes, checks,
enums, refs (all operators incl. `?` optional and `<>`), `TablePartial` injection, records,
`Dep` lineage, `TableGroup`, sticky notes, `DiagramView`, custom properties and `Metadata`
blocks, `Project`, and multi-file `use`/`reuse` imports — into a resolved model whose JSON
follows `@dbml/parse`'s per-file database shape. On top of the model it:

- **emits** canonical DBML (`emit`) — parse → emit → parse is model-identical and emitting is
  idempotent on every compiling test input;
- **exports** SQL DDL for MySQL/MariaDB, PostgreSQL and SQL Server (`export`), matching
  `@dbml/core` 10.3.1 on 99% of 582 golden inputs (junction tables for `<>`, enums, indexes,
  composite keys, ref actions, comments, checks, generated columns, records → INSERT), except
  where go-dbml deliberately writes safer or valid SQL (AGENTS.md §6); with
  `export.SQLWith(…, Options{TranslateTypes: true})` it translates types, defaults and
  expressions from the model's source dialect (93% of a cross-dialect corpus executes);
- **imports** MySQL/MariaDB, PostgreSQL and SQL Server DDL and Rails `schema.rb` (`importer`),
  matching `@dbml/core`'s `importer.import` on 97.5–100% of goldens per format (1,244 in total), plus an
  Atlas `schema.Realm` → model converter for live-database inspection (`importer/atlasrealm`,
  native builds only). Imports default to a *clean* mode (unescaped strings/names, normalised
  expressions, real comments, generated columns, Rails types mapped to SQL types,
  always-compiling DBML; INSERT data skipped unless asked for); a clean import exported back to
  the same dialect executes on a real database for 94–100% of a real-world dump corpus;
- runs **in the browser** as WASM (`wasm`, ~2.0 MB gzipped): `goDbml.parse`, `exportSQL`,
  `importSQL`, `emit`, with lower browser limits (1 MiB input, 2,000 record rows) and an
  optional compact model (`{"tokens": false}`).

It has no cgo, is safe on untrusted input (size, file, diagnostic, nesting, token, element and
records limits; cancellation via `dbml.ParseFileContext`; importers are linear and
input-capped), and depends only on the standard library plus `ariga.io/atlas` for the
native-only Atlas converter.

**Exported DDL:** names and string literals are quoted/escaped per dialect. Column types,
expression defaults, CHECK expressions, expression index parts and generated-column
expressions are SQL and are emitted as written, but each is first lexed with the target
dialect's rules (in every server mode) and refused with `export.ErrUnsafeExpression` unless it
is one balanced expression — no `;`, comments, unterminated quotes or batch separators.
Integration tests run cross-dialect smuggling attempts against MariaDB, PostgreSQL and SQL
Server; a sentinel table always survives. See AGENTS.md §6 for exactly what "matches" covers and the
known deviations.

```go
import (
    "github.com/aidenappl/go-dbml"
    "github.com/aidenappl/go-dbml/emit"
    "github.com/aidenappl/go-dbml/export"
    "github.com/aidenappl/go-dbml/importer"
)

db, diags := dbml.Parse(src)
if diags.HasErrors() {
    markers := diags.Markers() // monaco.editor.IMarkerData[]
    ...
}
pretty, err := emit.DBML(db)                   // canonical DBML
ddl, err := export.SQL(db, export.Postgres)    // or export.MySQL, export.MSSQL
fromSQL, err := importer.Import(sql, importer.MySQL) // *model.Database (clean mode, no records)
res, err := importer.ImportWith(sql, importer.Postgres, importer.ImportOptions{
    IncludeRecords: true, MaxRecordRows: 1000,
}) // res.Database, res.DBML, res.Diagnostics, res.RecordsTruncated, res.Warnings

// Servers: bound the work by a deadline.
ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
defer cancel()
r := dbml.ParseFileContext(ctx, "/schema.dbml", src, dbml.Options{MaxBytes: 2 << 20})

// Editors: a best-effort model is always available, even with errors.
preview := r.Partial
```

## Role in the Trailblaze ecosystem

It is the canonical DBML compiler for **dbdiagram.appleby.cloud** (appleby.cloud, not
Trailblaze). `dbdiagram-api` compiles on save, exports SQL and imports synced `schema.sql` /
`schema.rb` with it; `dbdiagram-cli` turns Atlas inspections into DBML through it; and
`dbdiagram-web` runs it as WASM in a Web Worker.

## Tech stack

- Go 1.25, standard library (+ `ariga.io/atlas` v1.2.2 for `atlasrealm.FromRealm`, native only)
- Hand-written recursive-descent + Pratt parser ported from Holistics' `@dbml/parse`
  (Apache-2.0); exporters and importer semantics ported from `@dbml/core` (Apache-2.0)
- Hand-written SQL tokenizer/statement parsers (no ANTLR, no cgo, no tidb/pg_query — see
  AGENTS.md §2 for why)

## Getting started

### Prerequisites

- Go 1.25+ (Node 20+ only to regenerate goldens)

### Setup

```sh
go get github.com/aidenappl/go-dbml
```

| Function | Use |
|---|---|
| `dbml.Parse(src)` | single source → `(*model.Database, diag.List)` |
| `dbml.ParseFile(path, src)` / `ParseFileWith` | `model.Result` with `Database`, `Partial`, `Diagnostics`, `Program` |
| `dbml.ParseFileContext(ctx, path, src, opts)` | `ParseFileWith` that stops with a `CANCELED` diagnostic when `ctx` is done |
| `dbml.ParseProject(entry, files)` / `ParseProjectWith` | multi-file project from a `map[path]source` |
| `dbml.ParseFS(fsys, entry)` / `ParseFSWith` | multi-file project from an `fs.FS` (use `os.OpenRoot(dir).FS()` for untrusted dirs) |
| `dbml.ParseSyntax(src)` | syntax tree + syntax errors only |
| `model.Validate(src, path)` | syntax + structural diagnostics without binding |
| `dbml.Options` / `DefaultOptions()` | limits: 16 MiB/file, 256 files, 1000 diagnostics, depth 256, 10,000 record rows, 4 MiB of Records, 2M tokens (and ≤ len/2), 10,000 elements per block; `Cancel` channel |
| `model.StripTokens(db)` | compact model for JSON (tokens removed, `line` kept) |
| `dbml.Version`, `model.JSONSchemaVersion` | store with persisted model JSON |
| `emit.DBML(db)` / `DBMLWith(db, emit.Options)` | canonical DBML text, `(string, error)`; never panics |
| `export.SQL(db, dialect)`, `export.ParseDialect(s)` | DDL for `mysql` (MariaDB), `postgres`, `mssql`; `ErrUnsafeExpression` / `*UnsafeExpressionError` for a refused raw fragment |
| `export.SQLWith(db, dialect, export.Options{TranslateTypes, SourceDialect})` | DDL translated from the source dialect (`Project database_type`, recorded by clean imports) → `*export.Result{SQL, Warnings}` |
| `export.CoreDBML(db, export.CoreDBMLOptions)` | `@dbml/core`'s DbmlExporter layout |
| `importer.Import(src, format)` | SQL/schema.rb → model (clean mode, records skipped); `*importer.CompileError` carries a partial model; formats `mysql`, `postgres`, `mssql`, `schemarb` |
| `importer.ImportWith(src, format, importer.ImportOptions)` | `Mode` (`ModeClean` / `ModeParity`), `IncludeRecords`, `MaxRecordRows`, `MaxInputBytes`, `OmitProject`, compiler `Model` options → `*importer.Result` (incl. `Warnings`) |
| `importer.CoreDBML(src, format)` | `@dbml/core`-identical DBML text (parity) |
| `atlasrealm.FromRealm(realm, atlasrealm.Options)` | Atlas `schema.Realm` → model (not in WASM) |

Browser: build with `dev wasm`, load `dist/wasm_exec.js` + `dist/go-dbml.wasm`, run it, then
call `goDbml.parse(src, opts?)`, `goDbml.exportSQL(src, "postgres", opts?)`,
`goDbml.importSQL(sql, "mysql", opts?)`, `goDbml.emit(modelJSON)` — each returns a JSON string
(`{ok, …}`; diagnostics are Monaco markers; `opts` is an optional JSON string such as
`{"maxRecordRows":500,"includeRecords":true,"tokens":false}`; `exportSQL` translates types
from the recorded source dialect unless `{"translateTypes":false}`). Run it in a Web Worker, and terminate and
recreate the worker after any thrown exception or timeout (contract in AGENTS.md §6). Define `globalThis.goDbmlReady` to be notified when the API is installed.

## Development

| Command | What it does |
|---|---|
| `dev build` | Build all packages, natively and for `GOOS=js GOARCH=wasm` |
| `dev test` | Unit, conformance, reference, limits, round-trip, export/import golden and WASM API tests |
| `dev run` | All golden suites with pass rates |
| `dev fuzz` | 30 s fuzz each: parser, emitter round trip, four importers |
| `dev integration` | Execute exported DDL on MariaDB 11.4, PostgreSQL 16 and SQL Server 2022 containers: adversarial names, smuggling attempts, clean-import round trips (Docker) |
| `dev wasm` | Build `dist/go-dbml.wasm` + `wasm_exec.js`, print raw/gzip size |
| `dev check` | gofmt + vet (native and js) + test |
| `dev tidy` | `go mod tidy` without toolchain upgrades |

## Project structure

```
dbml.go         public parse API
diag/           diagnostics, codes (== @dbml/parse), Monaco markers
lexer/          tokens + trivia, UTF-16 positions
ast/            syntax tree
parser/         error-recovering, depth-limited parser
model/          validation, binding, interpretation, JSON model, ResolveTable, Options
emit/           canonical DBML pretty-printer
export/         @dbml/core-compatible SQL (MySQL, Postgres, MSSQL) and DBML exporters, raw-fragment validation
importer/       MySQL / Postgres / MSSQL / schema.rb importers (clean + parity modes)
  atlasrealm/   Atlas realm converter (native only)
wasm/           browser entry point (goDbml global)
scripts/        dev-only golden regeneration (Node)
testdata/       Holistics fixtures, @dbml/parse + @dbml/core goldens, fuzz corpus
```

## Deployment

Library — released by git tag and fetched through the Go module proxy. The WASM file is a
build artifact copied into `dbdiagram-web`. Nothing to deploy.

## Contributing & further reading

- [AGENTS.md](AGENTS.md) — architecture, conventions, conformance, deviations, verification
- [scripts/README.md](scripts/README.md) — regenerating goldens
- [LICENSE](LICENSE) (Apache-2.0) and [NOTICE](NOTICE) — portions derived from
  [holistics/dbml](https://github.com/holistics/dbml) `packages/dbml-parse` and
  `packages/dbml-core`, and from `pluralize` (MIT)
- Related: `dbdiagram-api`, `dbdiagram-web`, `dbdiagram-cli` (planned)
