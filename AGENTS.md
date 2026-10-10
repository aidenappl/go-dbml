# AGENTS.md — go-dbml

> **⚠️ Golden rule — keep this file current:** any change to the public API, the JSON shape
> of `model.Database`, the package layout, the conformance suites, the resource limits or the
> commands below MUST update this AGENTS.md (and README.md where relevant) in the SAME change.

---

## 1. What this repo is

`go-dbml` is a pure-Go, zero-cgo implementation of **DBML** (Database Markup Language,
<https://dbml.dbdiagram.io>): a hand-written lexer, a recursive-descent/Pratt parser producing a
concrete syntax tree, and a semantic layer (validation, binding, interpretation) producing a
resolved model — plus a DBML pretty-printer, SQL DDL exporters, SQL/schema.rb importers, an Atlas
(live-database) importer and a WASM entry point for the browser. It is the canonical DBML
compiler for the `dbdiagram.appleby.cloud` project (plan:
`~/.claude/plans/twinkly-weaving-muffin.md`, Recommendation 1 / Tasks 0.2 + 0.3).

It **owns**: parsing DBML v2 (incl. `TablePartial`, `records`, `Dep`, `DiagramView`, `Metadata`
blocks, custom properties, `?` optional refs, multi-file `use`/`reuse`), diagnostics with
Monaco-ready ranges, a JSON model in the shape of `@dbml/parse`'s per-file database JSON, and
resource limits that make it safe on untrusted input.

It also owns (Task 0.3): canonical DBML emitting (`emit/`), DDL export byte-compatible with
`@dbml/core` 10.3.1 for MySQL/MariaDB, PostgreSQL and SQL Server (`export/`), import of
MySQL/MariaDB, PostgreSQL and SQL Server DDL and Rails `schema.rb` mirroring `@dbml/core`'s
`importer.import` (`importer/`), Atlas `schema.Realm` → model for live-DB inspection
(`importer/atlasrealm`, native only), and the browser API (`wasm/`, `goDbml` global).

It does **not** own: the editor's language service (that stays `@dbml/parse` in the browser),
Oracle/Snowflake import/export, diffing, layout, or any service/deployment. It is a library;
the only `main` is the WASM entry (`wasm/`).

Licensing: Apache-2.0 (`LICENSE`). Large parts are ported from Holistics' Apache-2.0
`@dbml/parse` and `@dbml/core`, and `importer/pluralize.go` from the MIT `pluralize` 8.0.0; see
`NOTICE` and the header of each ported file.

## 2. Stack & dependencies

- Go **1.25** (`go.mod`: `go 1.25.0`), module `github.com/aidenappl/go-dbml`. The go directive
  must stay ≤ 1.25 (Alpine `golang:1.25` builders).
- **Standard library only, with one exception:** `ariga.io/atlas` **v1.2.2** (Apache-2.0), used
  only by `importer/atlasrealm` (`//go:build !js`) for `sql/schema` (itself stdlib-only); every
  other package, `importer` included, is stdlib-only. Pinned
  at v1.2.2 because v1.2.3+/v1.3.0 require `go 1.26.4`. Do not add other modules.
- Evaluated and **rejected** for the importers (measured Oct 9 2026):
  - `github.com/pingcap/tidb/pkg/parser` (pure Go, builds for js/wasm) — alone 25.4 MB raw /
    4.2 MB gzipped as WASM (over the 3 MB budget by itself), and lossy for parity: it
    normalises types (`BOOL`→`tinyint(1)`, `INT`→`int(11)`, `SERIAL`→`bigint(20) UNSIGNED`,
    `REAL`→`double`), loses keyword case and keeps no original expression text, so defaults and
    checks cannot match `@dbml/core` (which renders ANTLR `getText`).
  - `github.com/wasilibs/go-pgquery` (wasm-embedded libpg_query via wazero, no cgo) — works
    natively (10.4 MB binary) and builds for js, but is 15.0 MB raw / 3.1 MB gzipped alone, and
    being the real Postgres grammar it diverges from the reference's ANTLR grammar.
  - `pganalyze/pg_query_go` — cgo, not allowed.
  The SQL importers are therefore hand-written (shared tokenizer + per-dialect statement
  parsers mirroring the reference's ANTLR AST generators); they are what the WASM build ships.
- Reference implementations: `holistics/dbml` at commit
  `34108bb9346df5f22d5224549c106ff47ca931c8`, published as `@dbml/parse` **10.3.1** (parser
  goldens) and `@dbml/core` **10.3.1** (export/import goldens). **Ground truth is the
  published npm package**, not the repo source: they differ in places (e.g. CREATE VIEW is
  silently dropped by the published importer).
- Node is only needed to *regenerate* goldens (`scripts/`), never to test.

## 3. Project structure

```
go-dbml/
  dbml.go                          # public API: Parse, ParseFile(With), ParseFileContext, ParseProject(With),
                                   #   ParseFS(With), ParseSyntax, Options/DefaultOptions, Version
  dbml_test.go                     # API tests + FuzzParse
  conformance_test.go              # Holistics interpreter (+nan) snapshots: DBML -> database JSON (gate ≥95%)
                                   #   + infinite_loops snapshot
  diagnostics_conformance_test.go  # Holistics lexer/parser/validator/binder snapshots (gate ≥95%)
  reference_test.go                # differential suite vs @dbml/parse 10.3.1 goldens (gate ≥95% + no regressions)
  limits_test.go                   # depth (10^6 nesting), perf (20k columns, 4k tables), Options, Partial
  dense_limits_test.go             # 2 MiB token/element-dense inputs: TOO_MANY_* diagnostics, bounded allocation
  cancel_test.go                   # ParseFileContext cancellation latency
  roundtrip_test.go                # Database -> JSON -> Database -> JSON stability
  race_{on,off}_test.go            # timing scale under -race
  diag/                            # Diagnostic, Code (numbering == @dbml/parse CompileErrorCode; 9000+ = go-dbml),
                                   #   Monaco markers
  lexer/                           # tokens (trivia kept as ranges, re-lexed on demand); int32 positions in UTF-16
                                   #   units (JS/Monaco semantics); MaxTokens, Cancel
  ast/                             # syntax tree node types, Element interface, factory (spans), walk helpers
  parser/                          # panic-mode recursive descent + Pratt expressions; never aborts; depth-,
                                   #   element- and records-limited (skips without building nodes)
  model/                           # semantic layer + resolved model
    types.go                       #   JSON model (Database, Table, Column, Ref, Dep, ...), JSONSchemaVersion
    options.go                     #   Options / DefaultOptions (resource limits, Cancel)
    compact.go                     #   StripTokens (compact JSON with `line` fields)
    compiler.go                    #   pipeline, recover, files/imports, symbols & scopes, memoised columns, Validate()
    syntax.go                      #   AST predicates/extractors (port of utils/expression + utils/validate)
    settings.go                    #   setting-list aggregation + per-kind setting validation
    validate.go                    #   per-element structural validation (local modules)
    bind.go                        #   reference resolution + binding/duplicate/circular-ref/import errors
    interpret.go                   #   AST -> model (tables, partials, enums, groups, notes, refs, deps, views, metadata)
    records.go                     #   records value typing + record infos (type, PK, UNIQUE, FK)
    refinfo.go                     #   INVALID_REF_RELATIONSHIP infos (operator vs nullability/uniqueness)
    partials.go                    #   Database.ResolveTable: @dbml/core partial-injection merge
  emit/                            # DBML pretty-printer: DBML(db) / DBMLWith(db, Options) -> (string, error);
    dbml.go validate.go            #   Ident, String; ErrNilDatabase, ErrInvalidModel, ErrUnrepresentable, InternalError
  export/                          # @dbml/core-compatible exporters (ports, Apache-2.0 headers)
    core.go                        #   port of @dbml/core model_structure (schema order, partial injection, enum
                                   #   binding, ref validation, injected refs) with map lookups (linear on 20k tables)
    export.go                      #   SQL(db, Dialect), Dialects, ParseDialect, ErrModel, InternalError, quoting helpers
    fragment.go                    #   raw-fragment validation per target dialect: ErrUnsafeExpression, UnsafeExpressionError
    translate.go                   #   SQLWith(db, Dialect, Options{TranslateTypes, SourceDialect}) -> *Result{SQL, Warnings};
                                   #   SourceDialectOf / DatabaseType (Project database_type); cross-dialect type translation
    reserved.go                    #   MySQL reserved words: IsMySQLReserved, "key" -> `key` in MySQL fragments
    mysql.go postgres.go mssql.go  #   MysqlExporter.js / PostgresExporter.js / SqlServerExporter.js (+ identifier/literal escaping)
    dbml.go                        #   CoreDBML(db, CoreDBMLOptions): DbmlExporter.ts + adjustCardinalitiesFromColumns (+ Escape mode)
    integration_test.go            #   //go:build integration: DDL executed on MariaDB 11.4, Postgres 16 (+ scs=off) and
                                   #   SQL Server 2022 (Docker): adversarial names, smuggling, clean round trips, DBML corpus
    testdata/smuggle/              #   cross-dialect smuggling inputs (re-verifier, Oct 2026)
    testdata/roundtrip/<format>/   #   real-world-shaped dumps for the round-trip test (re-verifier corpus)
  importer/                        # SQL / schema.rb -> model (stdlib only)
    importer.go                    #   Import, ImportWith(src, Format, ImportOptions) -> *Result, CoreDBML (parity text),
                                   #   Mode (ModeClean default, ModeParity), ErrInputTooLarge, ErrUnsupported,
                                   #   InternalError, CompileError, Formats, ParseFormat, MaxInputBytes
    lexer.go cursor.go             #   streaming tokenizer (MySQL DELIMITER, T-SQL GO, $$ quotes, N''/E''/b'' strings),
                                   #   statement splitter (drops skipped INSERTs unstored), cursor with paren-match table
    mysql.go postgres.go mssql.go  #   ports of MySQLASTGen / PostgresASTGen / MssqlASTGen semantics (hand-written parsers)
    clean.go normalize.go          #   clean mode: unescaping, expression normalisation (Result.Warnings), cleanModel,
                                   #   generated columns, dbo unqualification
    schemarb_clean.go              #   clean schema.rb (Rails types → SQL types, enums, virtual columns, references FKs)
    schemarb.go pluralize.go       #   parity schema.rb (PEG grammar port) + pluralize 8.0.0 (MIT)
    atlasrealm/                    #   FromRealm(*schema.Realm, Options) — //go:build !js; the only ariga.io/atlas user
  wasm/                            # browser entry: goDbml.{parse,exportSQL,importSQL,emit,version}, typed JSON responses
    api.go main_js.go main_other.go#   handlers (natively tested), js wiring (js && wasm), non-js stub main
  scripts/                         # dev-only golden regeneration (see scripts/README.md)
    reference/{ref,multi,export,import}.mjs, golden.py
  testdata/holistics/              # vendored Holistics fixtures (verbatim, NOTICE + LICENSE);
                                   #   dbml-core/{exporter,importer,parser}/ = @dbml/core SQL in/out fixtures
  testdata/reference/<set>/        # <name>.dbml (or project dir) + <name>.golden.json from @dbml/parse 10.3.1
                                   #   (hand4 = exporter-focused hand-written cases)
  testdata/export/<dialect>/<set>/ # <name>.sql|.err from @dbml/core 10.3.1 exporter.export (dialect = mysql,
                                   #   postgres, mssql, dbml); inputs are testdata/reference + Holistics exporter fixtures
  testdata/import/<format>/<set>/  # <name>.dbml|.err from @dbml/core 10.3.1 importer.import; sets: importer, parser
                                   #   (Holistics fixtures), hand (testdata/import/inputs/<format>/), exported
                                   #   (@dbml/core's own SQL export from testdata/export, deduplicated)
  testdata/fuzz/FuzzParse/         # fuzz regression corpus (keep: past crashers); emit/ and importer/testdata/fuzz likewise
  LICENSE, NOTICE, Devfile.yaml
```

## 4. Running, building & testing

No prerequisites beyond Go 1.25 (no DB, no network, no Node). Only `dev integration` needs
Docker (pulls `mariadb:11.4`, `postgres:16` and `mcr.microsoft.com/mssql/server:2022-latest`,
which runs as linux/amd64 — emulated on Apple silicon; ~2 min in total).

| Command | What it does |
|---|---|
| `dev build` | `go build ./...` + `GOOS=js GOARCH=wasm go build ./...` |
| `dev test` | `go test ./...` (unit, conformance, reference, limits, round trip, export/import goldens, wasm API) |
| `dev run` | all golden suites verbose: pass rates and failing cases |
| `dev fuzz` | 30 s each: `FuzzParse`, `FuzzEmitRoundTrip`, `FuzzTranslate`, `FuzzImport{MySQL,Postgres,MSSQL,SchemaRB}` |
| `dev integration` | `go test -tags integration -run Against ./export`: executes exported DDL on MariaDB 11.4, PostgreSQL 16 (also `standard_conforming_strings=off`) and SQL Server 2022 containers — adversarial names, smuggling (sentinel must survive), clean-import round trips (≥ 90% gate for MySQL/PostgreSQL), DBML corpus |
| `dev wasm` | builds `dist/go-dbml.wasm` (`-ldflags='-s -w'`) + copies `wasm_exec.js`, prints raw/gzip size |
| `dev check` | `gofmt -w -s . && go vet ./... && go test ./...` |

Benchmarks: `go test -run '^$' -bench . . ./importer` (`BenchmarkWideTable`, `BenchmarkManyTables`,
`BenchmarkPostgresDump` — the 8.7k-line Holistics pg_dump fixture imports in ~19 ms).

Browser use: load `wasm_exec.js` (from `$(go env GOROOT)/lib/wasm/`), optionally define
`globalThis.goDbmlReady = () => …`, instantiate `go-dbml.wasm` and `go.run(instance)`; then call
`goDbml.parse(src, optionsJSON?)` etc. Every function takes strings and returns a JSON string
(contract in §6, types in `wasm/api.go`).

## 5. How code is written here

- Flat package-per-concern layout (no `cmd/`, `internal/`, `pkg/`). Standard `testing` only
  (no testify), table-driven `t.Run` tests.
- **Port faithfully, don't reinvent.** When behaviour is unclear, read the Holistics source
  (`packages/dbml-parse/src/core/...`) and mirror it — quirks included (e.g. a records value
  written `Enum.member` is an *invalid* enum value in the reference). The suites compare exact
  token spans and diagnostic start positions, so "close enough" fails. Comments name the JS
  function being mirrored; ported files carry a provenance header (keep it when editing).
- Positions: lexer positions are 0-based with `Offset`/`Column` in **UTF-16 units**; model and
  diagnostics use 1-based line/column (Monaco). Use `lexer.ToPos`; never emit byte offsets.
- The parser never aborts: methods return `(node, *pErr)` where `pErr` carries a partial node
  and the context that should synchronise (JS `PartialParsingError`).
- **Every recursion must stay bounded.** The parser enforces `MaxDepth` (default 256) at every
  recursive entry (`expressionBP`, block/list/tuple) and on left-associative operator chains;
  too-deep input is skipped *iteratively* with one `NESTING_TOO_DEEP` error, so later walkers
  only ever see bounded trees. New recursive code must rely on that bound.
- **No quadratic passes.** Merged columns are memoised per (table, scope) with a name index;
  scopes index visible symbols by symbol; Metadata blocks are indexed once per compile. Keep
  lookups map-based (`limits_test.go` enforces 20k columns / 4k tables in < 1 s).
- Any input must produce a result without panicking — `FuzzParse` enforces it, and
  `CompileProjectWithOptions` converts a panic into an `INTERNAL_ERROR` diagnostic as a last
  line of defence (a bug if it ever fires). Add crashers to `testdata/fuzz/FuzzParse/`.
- **Exporters/importers are ports of the published `@dbml/core` 10.3.1 JS**, quirks included
  (raw unescaped defaults, `int IDENTITY(1,1)` as a T-SQL type name, `doubleprecision`, unary
  minus dropped from Postgres INSERT values, N'..' kept verbatim, …). Port the visitor's
  *semantics*, not ANTLR: text a visitor renders with `ctx.getText()` is `cursor.getText`
  (token texts concatenated, no whitespace) and `getOriginalText` is `cursor.origText`. When
  the reference throws (syntax error, TypeError), mirror it with an `ErrUnsupported` /
  `ErrModel` error instead of guessing — but only for cases the goldens show.
- Public entry points recover panics into typed errors (`export.SQL`/`CoreDBML` →
  `*export.InternalError`, `importer.Import`/`ImportWith`/`CoreDBML` → `*importer.InternalError`,
  `emit.DBML` → `*emit.InternalError`, every WASM function); the fuzz targets call the internals without the recover so panics still fail.
- `emit` is the canonical formatter (minimal quoting, source order); `export.CoreDBML` is only the
  `@dbml/core` layout used to mirror `importer.import`. Never "fix" one to look like the other.
- Importers: parity behaviour lives behind `cfg.clean == false`; clean-mode relaxations and
  unescaping must never change parity output (the parity goldens are the check). Use the
  `im.str` / `im.expr` / `im.ident` helpers instead of reading token text directly.
- New escapers go through `bq`/`dq`/`sb` (identifiers) and `myLit`/`pgLit`/`stdLit`/`nLit`
  (literals); never concatenate a user-controlled name or string into DDL without one, and
  never emit a raw SQL fragment (type, expression) without `safeExpr`/`safeType`.
- Clean-mode importers keep raw expression text (`origText`); `normalizeModel` rewrites it
  once for the whole model — do not unescape expressions in the per-dialect parsers.
- Memory: a lexer token is ~104 bytes and an AST node ~104+; do not add per-token or
  per-node fields or allocations without measuring (`dense_limits_test.go`).

## 6. Domain & architecture

Pipeline (per file, then across imports): **lex → parse → validate → bind → interpret**.
Any error from lex/parse/validate/bind suppresses interpreter diagnostics and makes
`Result.Database` nil (as `@dbml/parse`); warnings/infos never block. `Result.Partial` is
**always** interpreted (best effort, for editors) and equals `Database` when there are no errors.

- **Limits** (`model.Options`, zero value = defaults): `MaxBytes` 16 MiB per file
  (`INPUT_TOO_LARGE`, `ParseFS` reads at most MaxBytes+1), `MaxFiles` 256 (`TOO_MANY_FILES`),
  `MaxDiagnostics` 1000 (rest summarised by one `TOO_MANY_DIAGNOSTICS`; lexer/parser stop
  recording early), `MaxDepth` 256 (`NESTING_TOO_DEEP`), plus `MaxRecordRows`,
  `MaxRecordBytes`, `MaxTokens`, `MaxElements` and `Cancel` (table under *Resource limits and
  memory* below). go-dbml-only codes are 9000+ (9000–9008).
- **Symbols** (`model/compiler.go`): top-level `Table`, `Enum`, `TablePartial`, `TableGroup`,
  sticky `Note`, `DiagramView`. Schema defaults to `public` (reported as `schemaName: null`).
  Table aliases are program-level table names. Columns of a table = own columns + injected
  partial columns (own wins; among partials the last injection wins).
- **Refs**: short/long/inline, composite `t.(a, b)`, cross-schema `s.t.c`, actions, `color`,
  `inactive`, all 16 operators (`< > - <>` with `?` on either side → cardinalities
  `1 * 0..1 0..*`). Inline refs appear both in `column.inline_refs` and `db.refs`.
  `INVALID_REF_RELATIONSHIP` infos flag operators that disagree with nullability/uniqueness.
- **Records**: values typed per column type family (integer/real/bool/string/date/time/
  datetime/enum/expression); `INVALID_RECORDS_FIELD` infos for type mismatches, NOT NULL,
  string length, numeric precision/scale, and PK (missing/null/duplicate), UNIQUE and FK
  violations across tables.
- **Deps**: short form, block form (note/color/custom keys), inline on table header and column.
- **TablePartial**: the JSON keeps partials separate (`table.partials[]` with `order`) exactly
  like `@dbml/parse`; `Database.ResolveTable` applies the documented merge precedence the way
  `@dbml/core` does (local > last-injected partial; injected indexes/checks placed before the table's own, as
  `@dbml/core` does; injected
  check names prefixed `<table>.`; injected inline refs returned as `InjectedRefs`).
- **Metadata blocks / custom properties**: inline `[k: 'v']` + `Metadata <Kind> <name> {}`;
  precedence inline < imported-file blocks < current-file blocks (later wins).
- **Multi-file**: `use`/`reuse` with `*` or selective `{ table|enum|tablepartial|note|schema|tablegroup name [as alias] }`.
  `use` is not transitive, `reuse` re-exports. Paths are relative, `.dbml` optional. Missing
  files/elements produce the reference's NONEXISTENT_MODULE / BINDING_ERROR diagnostics; an
  unresolved selective import still declares its name (as the reference's `UseSymbol` does).
  Imported elements are interpreted into the entry database and listed in `externals`;
  refs/deps/records from imported files are included when all their tables are visible.
- **ParseFS security**: `os.DirFS` follows symlinks out of its root. Callers handling
  untrusted directories must pass `os.OpenRoot(dir)`'s `FS()` (documented on `ParseFS`).
- **Versioning**: `dbml.Version` (compiler) and `model.JSONSchemaVersion` (JSON shape) — store
  both with persisted `model_json`.
- **JSON round trip**: `Database` → JSON → `Database` → JSON is byte-stable (`Project.Extra`
  and undefined record values are preserved by custom (Un)MarshalJSON).

### Emit, export, import, WASM (Task 0.3)

- **`emit.DBML(db) (string, error)`** writes canonical DBML: 2-space indent, identifiers quoted
  only when needed (keywords like `note`/`table` always quoted), column settings in a fixed
  order (`pk, increment, not null|null, unique, default, ref, dep, check, note, custom`),
  multi-line notes as indented `'''` blocks when DBML's note normalisation leaves them
  unchanged (otherwise an escaped single-line string). Elements come out in source order when
  the model has tokens and in kind order otherwise. Inline refs / column deps are written
  inline and their `db.refs` / `db.deps` twins skipped; table-header deps go back on the
  header; records become top-level `Records t(cols)` blocks; `externals` are not emitted; a
  `Column.Generated` without its custom properties gets them written. It never panics:
  `ErrNilDatabase`, `ErrInvalidModel` (nil elements, e.g. hand-written JSON),
  `ErrUnrepresentable` (an expression containing a backtick — DBML cannot escape it, so it is
  refused rather than altered), `*InternalError` (recovered panic). Numeric defaults that are
  not numbers are written as strings.
- **`export.SQL(db, dialect)`** rebuilds `@dbml/core`'s normalized model (`export/core.go`):
  schema creation order (enums, tables, notes, refs, deps, groups, then partial-injected refs
  into `public`), partial injection with injected indexes/checks *before* the table's own,
  the shared-partial-type enum binding quirk, ref validation errors (`ErrModel`). Lookups and
  ref de-duplication are map-based (20k tables + 20k refs: 35–80 ms per dialect). `<>` refs
  become `<a>_<b>` junction tables, `[example]` records are skipped.
  - **Names and strings are escaped:** every identifier is quoted with the dialect's quote
    doubled (`` ` `` → ``` `` ```, `"` → `""`, `]` → `]]`); every string literal is escaped:
    MySQL `\` → `\\` and `'` → `''`; PostgreSQL `'` → `''`, and a value containing a
    backslash is written `E'…'` with the backslashes doubled, so it reads the same with
    `standard_conforming_strings` off (`pgLit`); SQL Server `'` → `''`, with
    `sp_addextendedproperty` values and level names as `N'…'` (`nLit`). This covers
    table/column/index/check/ref/schema names, enum values, string defaults, notes/comments,
    extended properties and record strings. Referential actions and index types must be
    plain keywords (else `ErrModel`). A non-numeric record value in a numeric column, a
    non-numeric number default, or a non-hex binary value is quoted (the reference prints it
    raw).
  - **Raw fragments are validated, never trusted** (`export/fragment.go`): column types,
    expression defaults, CHECK expressions, expression index parts, expression record values,
    `CAST(… AS type)` types and generated-column expressions are SQL text and are emitted as
    written (like `@dbml/core`) — but only after being lexed with the **target** dialect's
    rules, in every lexical mode the server may be in (MySQL with and without
    `NO_BACKSLASH_ESCAPES` / `ANSI_QUOTES`, PostgreSQL with `standard_conforming_strings` on
    and off, SQL Server; `"…"` as identifier and as string). A fragment must be one balanced
    expression in all of them: no `;`, no comment (`--`, `/*`, MySQL `#`), no unterminated
    string / quoted identifier / `[name]` / `$tag$` quote, balanced parentheses, no backslash
    outside a string literal, no `GO` line (sqlcmd batch separator), no `DELIMITER` line, no
    NUL byte, no ODBC `{…}` escape (SQL Server). Anything else fails the whole export with
    `*UnsafeExpressionError` (wraps `ErrUnsafeExpression`; names the fragment, where it was
    and why). Nothing is rewritten to make a fragment pass: a string such as `'x\''` (MySQL
    escaping) is accepted for SQL Server (one string there) but refused for MySQL/PostgreSQL
    where its end depends on the mode. The only rewrite is for MySQL: `"w"` where `w` is a
    MySQL reserved word becomes `` `w` `` (a column such as `key` reaches DBML as `"key"`, which
    MySQL would read as a string); other `"…"` tokens are left alone. Semantics of a fragment across dialects stay the author's business
    (a MySQL-only function exported to PostgreSQL fails to execute; it cannot inject).
  - **Generated / computed columns** are model fields (`Column.Generated{Expression,
    Stored}`, from the custom properties `[generated: '<expr>', generated_stored: 'true']`):
    MySQL `type GENERATED ALWAYS AS (expr) STORED|VIRTUAL`, PostgreSQL
    `type GENERATED ALWAYS AS (expr) STORED` (PostgreSQL < 18 has no VIRTUAL), SQL Server
    `[c] AS (expr) [PERSISTED]` (no type; NOT NULL only when persisted). Defaults and
    increment are dropped for them.
  - **Other go-dbml-only export fixes:** MySQL `FULLTEXT`/`SPATIAL` index types become
    `CREATE FULLTEXT INDEX` (not `USING`); a MySQL expression index part `col(n)` naming a
    column of the table is written as a prefix key part `` `col`(n) ``; PostgreSQL
    parameterised / multi-word / aliased built-in types (`VARCHAR(20)`,
    `character varying(255)`, `timestamp(6) with time zone`, `double precision`,
    `bit varying(8)`, `TIMESTAMPTZ`, `INT8`, …) are not quoted as one identifier; SQL Server
    index column lists use `[x]` (the reference's `"x"` is a string under
    `QUOTED_IDENTIFIER OFF`, sqlcmd's default); `int unsigned`-style types count as numeric
    for record values; a MySQL foreign key from a table in a named schema (a MySQL database)
    to a default-schema table is a prepared statement that qualifies the referenced table
    with `DATABASE()` (`SET @dbml_stmt = CONCAT(…); PREPARE …; EXECUTE …`) — the reference
    leaves it unqualified, so MySQL looks for it in the child's database; PostgreSQL
    `interval day to second`-style types are not quoted.
  - `export/integration_test.go` (Docker) proves it on real servers — MariaDB 11.4,
    PostgreSQL 16 (also with `standard_conforming_strings=off`) and SQL Server 2022 (sqlcmd,
    `QUOTED_IDENTIFIER` off for the attack tests): 17 adversarial name/string models execute
    with exactly the intended tables; 19 smuggling inputs (the re-verifier's MySQL / PG /
    T-SQL / schema.rb dumps under `testdata/smuggle` imported in clean mode, plus 14
    hand-written DBML fragment attacks: `);DROP`, `$$`, `[`, `GO` lines, types, index parts,
    generated columns, record expressions, `e'\''`) exported to every dialect are refused or
    fail harmlessly, the sentinel table surviving everywhere (mariadb 14 refused / 2 ran / 3
    failed to execute; postgres 12/4/3; postgres scs=off 12/4/3; mssql 11/4/4).
- **Cross-dialect type translation** (`export.SQLWith(db, d, Options{TranslateTypes: true})`,
  `export/translate.go`): off by default (`SQL` and the goldens are untranslated); needs the
  source dialect — `Options.SourceDialect` or the model's `Project { database_type: '…' }`
  (`SourceDialectOf`; clean imports record it) — and does nothing when source = target.
  It maps types through a neutral kind (bool, int1–int8 with unsigned widening outside
  MySQL, numeric, float4/8, money, char/varchar/text with `(MAX)`, binary/varbinary/blob,
  date/time/timetz/timestamp/timestamptz (precision ≤ 6), interval, year, json/jsonb, uuid,
  xml, bit/varbit, inet/macaddr, set, arrays, geometric types; serial/`nextval` →
  increment; MySQL `tinyint(1)` and SQL Server `bit` → bool; text/blob used in a key →
  bounded `varchar(255)` / `nvarchar(450)` / …); translates defaults and CHECK / generated /
  index expressions (PostgreSQL `::casts` stripped, `now()`/`getdate()`/
  `CURRENT_TIMESTAMP`, `curdate()`/`CURRENT_DATE`, uuid functions, `ISNULL`/`IFNULL` →
  `COALESCE`, `LEN`/`CHAR_LENGTH`/`length`, `~` ↔ `REGEXP`, `interval '1 year'` ↔
  `INTERVAL 1 YEAR`, T-SQL `'a' + b` → `||`, `"x"` → `` `x` `` for MySQL, upper-case
  column names quoted for PostgreSQL, `col = 1` → `col = true` on PostgreSQL booleans,
  boolean / bit / numeric / zero-date / tz-suffixed defaults). A MySQL `tinyint(1)` /
  `bit(1)` / `BOOLEAN` becomes a boolean only when its default is absent or exactly
  0/1/true/false and every CHECK, index expression and generated column uses it only as a
  boolean (`= TRUE|FALSE|1|0`, `<>`, `IS [NOT] TRUE|FALSE|NULL`, `NOT col`, alone between
  `AND`/`OR`); a numeric use (BETWEEN, ranges, arithmetic, other integers) keeps it numeric
  (`smallint`) with a warning — a default that is not clearly boolean is never turned into
  `false`. It drops (with a warning)
  what cannot be translated safely — expressions calling functions outside a portable list,
  PostgreSQL-only operators, column references in defaults, non-deterministic MySQL CHECKs,
  non-immutable PostgreSQL generated columns, functional index parts for MariaDB / SQL
  Server, indexes on unindexable target types, foreign keys to non-unique columns
  (PostgreSQL / SQL Server; MySQL gets an index instead), `ON … SET NULL` on NOT NULL
  columns and cascading self / multi-path references for SQL Server; renames index names
  that would clash (PostgreSQL: per schema); aligns a foreign key column's type with the
  referenced column; creates MySQL schemas `IF NOT EXISTS`. Unknown types are kept as
  written with a warning. `Result.Warnings` lists every decision, each naming
  `table.column` (or the table / index): kept or unknown types, numeric `tinyint(1)`, and
  every lossy or semantic mapping — unsigned widened or dropped, SET → text, YEAR →
  smallint, json/jsonb/xml/interval/inet/macaddr/arrays/uuid stored as another type, money,
  dropped time zones, fractional-seconds or numeric precision changes, key columns bounded
  to `varchar(255)` / `nvarchar(450)`, long char/varchar becoming text, bit strings,
  geometric types, enums becoming a CHECK (SQL Server), dropped index methods, removed
  PostgreSQL casts, dropped defaults / checks / generated expressions / index parts /
  foreign keys / referential actions, renamed indexes and retyped foreign key columns
  (`TestTranslateWarnings`, `TestTranslateTinyintBool`, `TestTranslateBoolInChecks` with the
  Holistics `mysql_importer/checks.in.sql` fixture).
  - **Measured** (real servers): the web team's `check-sql-export` corpus (go-dbml import
    fixtures imported in clean mode via the WASM `importSQL`, exported via `exportSQL` with
    default options): round trip **37/37**, cross-dialect MySQL ↔ PostgreSQL (incl. SQL
    Server and schema.rb sources) **137/147 (93.2%)** — was 37/137; the 10 failures are
    source types that do not exist (`orders_status_enum` ×3, `intdsfsd`,
    `key_column_name_type`, a domain `Name`), two cases where the harness re-creates a
    server-wide MySQL database another case already filled (`Table 'users' already exists`),
    a schema.rb CHECK comparing a uuid with 0, and one MySQL generated column PostgreSQL
    deems not immutable.
    `TestCleanRoundTripAgainstRealDatabases` exports the round-trip corpus to every other
    dialect: **97/98 (99%)** execute (gate 85%; the failure is `x1_postgres` → SQL Server, a
    generated column `("Name" IS NOT NULL) * 2` T-SQL cannot express).
- **`export.CoreDBML(db, opts)`** is `DbmlExporter` (quotes everything, escapes nothing); with
  `OverrideCardinality` it applies `adjustCardinalitiesFromColumns`; `Escape` makes the output
  always lex (used by clean imports).
- **Importers** (`importer`): `ImportWith(src, format, ImportOptions{Mode, IncludeRecords,
  MaxRecordRows, Model})` returns `*Result{Database, Partial, DBML, Diagnostics,
  RecordsTruncated, Warnings}`; `Import(src, format)` is clean mode without records and
  returns the model, or `Result.Partial` + `*CompileError` when the generated DBML does not
  compile (callers must check the error before trusting the model); `CoreDBML(src, format)`
  is the exact `@dbml/core` text (parity, records, no cap). The row cap is
  `min(MaxRecordRows (default 10,000), Model.MaxRecordRows)` — a stricter compiler limit is
  never raised.
  - **ModeClean (default):** SQL literals and quoted identifiers are unescaped (MySQL `''`,
    `\'` and backslash escapes, `N'..'`, `_charset'..'`; Postgres `''`, `E'..'`, `$$..$$`, `""`
    identifiers; T-SQL `N'..'`, `[a]]b]`), comments keep their real text, the reference's
    syntax-error mirrors are relaxed (non-MySQL types, `INT(11)` in Postgres, multi-word T-SQL
    types, `CHECK CONSTRAINT ALL`, positional `sp_addextendedproperty`, expression index keys
    in T-SQL, …), and `cleanModel` drops records/refs/index columns pointing at unknown tables
    or columns, duplicate refs and duplicate ref names.
    - **Expressions are normalised** (`importer/normalize.go`): every kept expression (default,
      CHECK, index part, record expression, generated column) is re-lexed with the *source*
      dialect: string literals become `'…'` with `''` (MySQL `"…"` strings and escapes,
      PostgreSQL `E'…'` and `$tag$…$tag$`), `` `x` `` / `[x]` become `x` (or `"x"`), `"x"` is
      kept. A value with a backslash keeps the source convention in its mode-independent form
      (MySQL `'\\'` never `\'`, PostgreSQL `E'\\'`, T-SQL verbatim) and adds a
      `Result.Warnings` entry; from schema.rb (dialect unknown) such an expression is dropped.
      An expression with a comment, `;`, unterminated quote or unbalanced parentheses is
      dropped with a warning. The exporters still validate every fragment for the target.
    - **The source dialect is recorded:** clean mode starts the DBML with
      `Project { database_type: 'MySQL' | 'PostgreSQL' | 'SQL Server' }` (schema.rb: MySQL
      when the schema shows MySQL options, else PostgreSQL), which `TranslateTypes` and the
      WASM `exportSQL` read. `ImportOptions.OmitProject` (WASM `omitProject`) leaves it out —
      use it when merging the DBML into a document that has its own Project.
    - **Fidelity:** MySQL reserved-word identifiers (`` `key` ``) stay quoted (`"key"`, see
      the MySQL fragment rewrite above); MySQL keeps `unsigned`/`zerofill` in the type, `FULLTEXT`/`SPATIAL` as
      the index type, prefix key parts as `col(n)` expressions, generated columns
      (`GENERATED ALWAYS AS (…) STORED|VIRTUAL|PERSISTENT`), and applies `ALTER TABLE … ADD
      [COLUMN] … [FIRST|AFTER]`, `ADD (…)`, `ADD INDEX|KEY|UNIQUE|FULLTEXT|SPATIAL`,
      `MODIFY`, `CHANGE` (renames follow into indexes, refs, checks and generated
      expressions), `RENAME COLUMN`, `DROP COLUMN|INDEX|KEY|FOREIGN KEY|PRIMARY KEY`. `DESC` in index keys is dropped (DBML has no
      direction). PostgreSQL keeps multi-word types (`double precision`, `bit varying(n)`,
      `character varying(n)`, `timestamp(n) with time zone`, `time with time zone`) and folds
      unquoted type names to lower case; `GENERATED ALWAYS AS (…) STORED` becomes a generated
      column; `uuid_generate_v4()` (uuid-ossp) becomes the built-in `gen_random_uuid()`
      (also from schema.rb). T-SQL `int IDENTITY(1,1)` becomes `int` + increment, `c AS (expr) [PERSISTED]` a
      computed column (type `sql_variant`), and the `dbo` schema is the default (no `dbo.`
      prefix, so an export does not `CREATE SCHEMA dbo`).
    - **schema.rb** gets a separate clean parser keeping null/default/limit/precision/scale/
      comment, table comments, `id: false` / `id: :uuid` (+ `default:`) / `primary_key:`, an
      implicit id typed `bigint` (Rails ≥ 5.1, from `Schema[x.y]` or the version date) or
      `integer`, `t.timestamps`, `t.references` / `belongs_to` (with the index Rails creates
      and `foreign_key: true` / `{ to_table:, on_delete: }` refs), `t.index`, `add_index`,
      `add_foreign_key` with any options (`:nullify` → `set null`), `create_enum` +
      `t.enum enum_type:`, `t.virtual … as:, stored:` (generated column), check constraints,
      and it skips unknown lines instead of dropping the table. A foreign key column typed
      unlike the key it references (a `bigint` `t.references` to a `uuid` id) takes the key's
      type. Rails abstract types become the
      adapter's SQL types: PostgreSQL by default (`string` → `varchar[(n)]`, `datetime` →
      `timestamp[(p)]`, `float` → `double precision`, `binary` → `bytea`, `integer limit:` →
      `smallint`/`integer`/`bigint`, …), MySQL when the schema shows MySQL table options
      (`charset:`, `collation:`, `options: "ENGINE…"`, `unsigned:`) (`varchar(255)`,
      `datetime(p)`, `tinyint(1)`, `tinyint unsigned`, `mediumtext`, `char(36)` for uuid, …).
  - **Same-dialect round trip** (clean import, records on → export the same dialect →
    execute on a fresh database; `TestCleanRoundTripAgainstRealDatabases`, corpus
    `export/testdata/roundtrip`: the re-verifier's 45 real-world-shaped dumps + 4 mixed
    stress files): **mysql 15/16 (94%), postgres 17/17, mssql 8/8, schema.rb 7/8**. The two
    failures are inputs that are invalid themselves (`m15`: `ON DELETE SET NULL` on a NOT NULL
    column — MariaDB rejects the source too; `r04b`: a foreign key to a non-unique column).
    The PostgreSQL target has `citext`/`uuid-ossp`/`hstore`/`pgcrypto` installed (the dumps
    use them); SQL Server runs with `QUOTED_IDENTIFIER ON` (as SSMS). 8 of the 31 SQL sources
    fail when run raw on a fresh database (roles, `USE Shop`, missing schemas or sequences) —
    the round trip works for all of them. The 12 compiling DBML corpus models exported to
    each dialect execute 6/12 (mysql), 8/12 (postgres), 4/12 (mssql): failures are the
    models' own (PostgreSQL-only types/functions, refs to non-unique columns, a duplicate
    primary key), identical in `@dbml/core`.
  - **ModeParity:** byte-identical to `importer.import`, including its failures
    (`ErrUnsupported`).
  - **Records:** INSERTs are skipped unless `IncludeRecords`; skipped statements are never
    stored (the splitter streams), so a 16 MB INSERT dump imports with ~0 MB allocated
    (~0.25 s) in every dialect — T-SQL included: a GO batch is further split at `;` and
    before `CREATE`/`INSERT`/`EXEC` at paren depth 0 (routine/view batches are never split
    nor stored), so an SSMS-style dump of one INSERT per line streams (was 970 MB). With
    records, rows are capped at the row cap and stored tokens per statement at
    `rows*256+4096`; `RecordsTruncated` reports the cap (16 MB dump, records on: ~50 MB
    allocated, 0.3 s).
  - Inputs above 16 MiB → `ErrInputTooLarge`; panics → `*InternalError`. 20k-table dump:
    ~0.8 s (`Import`, incl. compiling the DBML). Deep nesting (200k parens / unary chains) is
    linear and non-recursive in every dialect (a precomputed paren-match table; T-SQL
    expressions peel parens/unary operators iteratively, beyond 256 levels the raw text is
    kept). schema.rb arguments (hashes, arrays, `-> { }` lambdas) are parsed in one pass with
    an explicit stack — no recursion, no re-scan; the per-line regexps and the MySQL-option
    detection are plain scans. 200k-deep hashes / arrays / lambdas and a 1 MiB
    `{{{…}}}` line import in 11–66 ms natively (clean) and < 1 ms (parity); in WASM
    (2 MiB allowed) 180–450 ms, no stack error, instance alive (`TestSchemaRBDeepNesting`).
  - **Input cap:** `MaxInputBytes` (16 MiB) by default; `ImportOptions.MaxInputBytes` lowers
    it per call (callers choose it per format). WASM `importSQL` uses `maxBytes` (1 MiB) or
    a per-format `maxInputBytes` object, e.g. `{"maxInputBytes": {"schemarb": 262144}}`.
- **`atlasrealm.FromRealm(realm, opts)`** (native only, own package so `importer` stays
  stdlib-only) maps an Atlas realm: raw column types (split into name+args), NOT NULL,
  defaults, auto-increment (`AutoIncrement`/`Identity`/`SerialType` by type name), comments,
  PKs, indexes (non-btree type kept), checks, enums, FKs → refs. The default schema (`public`
  or the realm's only schema) is unqualified.

#### Resource limits and memory

`model.Options` (all additive; zero = default):

| Option | Default | Enforced by | Diagnostic |
|---|---|---|---|
| `MaxBytes` | 16 MiB | compiler | `INPUT_TOO_LARGE` (9002) |
| `MaxFiles` | 256 | compiler | `TOO_MANY_FILES` (9003) |
| `MaxDiagnostics` | 1000 | everywhere | `TOO_MANY_DIAGNOSTICS` (9001) |
| `MaxDepth` | 256 | parser | `NESTING_TOO_DEEP` (9000) |
| `MaxRecordRows` | 10,000 | parser (skips later Records bodies) + interpreter | `RECORDS_TOO_LARGE` (9005) |
| `MaxRecordBytes` | 4 MiB | parser (skips the body before building nodes) + interpreter | `RECORDS_TOO_LARGE` |
| `MaxTokens` | 2,000,000 non-trivia tokens, and at most `len(src)/2 + 65,536` | lexer stops | `TOO_MANY_TOKENS` (9006) |
| `MaxElements` | 10,000 per block / list / tuple / row (columns, enum values, index columns, cells, rows) | parser skips the rest | `TOO_MANY_ELEMENTS` (9007) |
| `Cancel` | nil | lexer, parser, validator, binder, interpreter poll it | `CANCELED` (9008), `Database` and `Partial` nil |

`dbml.ParseFileContext(ctx, path, src, opts)` wires `ctx.Done()` to `Cancel`; a cancel
returns within a few ms (4 MiB schema: cancel at 10/33/66% of the 1.4 s compile returned
after 10/1/1 ms). The density guard (`len/2`) is what bounds adversarial input: hand-written
DBML has well under one significant token per 2 bytes, `a,a,a,…` has one per byte.

Memory work in this round: trivia (whitespace/comments) is no longer materialised as tokens —
a token stores its trivia ranges and `Token.LeadingTrivia()/TrailingTrivia()` re-lex them on
demand (lexer conformance unchanged); `lexer.Position` fields and the AST `NodeBase` are
int32/uint8-packed (token 192 → 104 bytes, node 136 → 104); settings lists, `Children` walks
and empty metadata maps no longer allocate per element. Measured peak Go heap, native, 2 MiB
inputs at default options (re-verifier shapes): one table of 157k columns 85 MB, enum of 220k
values 62 MB, index of 1M columns 196 MB, one row of 1M cells 196 MB, 524k one-cell rows
108 MB, 84k one-row Records blocks 340 MB, a *valid* schema of 30k small tables 413 MB (its
result holds 252 MB: 207 MB syntax tree + tokens, 45 MB model). Large valid schemas are the
remaining cost (~200 bytes of peak heap per source byte); size `MaxBytes` for the service
accordingly (dbdiagram-api: 2 MiB → ≤ ~420 MB per concurrent compile).

#### WASM JSON contract (`wasm/api.go`)

`goDbml.parse(src, optionsJSON?)`, `goDbml.exportSQL(src, dialect, optionsJSON?)`,
`goDbml.importSQL(sql, format, optionsJSON?)`, `goDbml.emit(modelJSON)`, `goDbml.version`.
Every function returns a JSON string with `ok` (bool); when `ok` is false `error` (string) is
set. `diagnostics` (array of Monaco `IMarkerData`) is always present on parse/exportSQL/
importSQL. Shapes:

| Call | ok = true | ok = false |
|---|---|---|
| parse | `database`, `diagnostics`, `version`, `schemaVersion` — with `{"tokens": false}` the model has no `token` objects and tables, columns, refs, enums, table groups, partials and sticky notes carry `line` (1-based start line, for click-to-source); the JSON is ~2.4x smaller | `error`, `partial` (compact too with `tokens: false`), `diagnostics`, … |
| exportSQL | `sql`, `diagnostics`, `translated?` + `sourceDialect?` (types were translated), `warnings?` | `error`, `diagnostics` (an unsafe fragment: `error` names it) |
| importSQL | `dbml`, `database`, `diagnostics`, `recordsTruncated?`, `warnings?` | `error`, `diagnostics`, and when the import ran but the DBML does not compile also `dbml`, `partial`, `warnings?` |
| emit | `dbml` | `error` |

`optionsJSON` (all optional): `maxBytes` (default **1 MiB**), `maxRecordRows` (2,000),
`maxRecordBytes` (1 MiB), `maxDiagnostics` (500), `maxDepth` (256), `maxTokens`,
`maxElements` (native defaults), `mode` (`clean`|`parity`, importSQL), `includeRecords`
(importSQL), `omitProject` (importSQL), `maxInputBytes` (importSQL: per-format input cap,
`{"mysql": n, "schemarb": n, …}`), `tokens` (parse / importSQL: `false` = compact
model, see the table), `translateTypes` (exportSQL; default **on** when the model records a
source dialect — `Project database_type`, e.g. any clean `importSQL` result — that differs
from the target; `false` exports untranslated), `sourceDialect` (exportSQL: overrides the
recorded one). Values above the native defaults are clamped. Every key added in this round is
additive.

The WASM build sets `GOGC=50` and the handlers drop the syntax tree and collect garbage
before encoding a large response. Linear memory (never shrinks) at the defaults: any 2 MiB
probe is refused up front (13 MB); a valid 1 MiB schema of 15k tables peaks at ~300 MB
(2.4 s), 34k one-column tables 260 MB. With `maxBytes: 2 MiB`: dense adversarial inputs
34–210 MB (1M-column index 210 MB), 84k one-row Records blocks 374 MB, a valid 2 MiB schema
of 30k tables ~590 MB (5 s) — hence the 1 MiB default.

**Browser worker contract:** run the WASM in a dedicated Web Worker; wrap every call in
`try/catch` — if anything throws (a Go runtime fatal error cannot be recovered and leaves the
instance dead), terminate the worker and create a fresh instance; enforce a per-call timeout
(terminate + recreate on expiry: `Cancel` cannot fire inside a synchronous WASM call); never
reuse an instance after an exception. The `dbdiagram-web` worker implements this.

#### JSON model changes (`model.JSONSchemaVersion` = 2)

Version 2 adds `column.generated` (`{expression, stored}`, only for generated columns) and
omits a nil `token` instead of writing `"token": null` (every `@dbml/parse` golden is
unchanged: their tokens are present). `line` fields appear only after
`model.StripTokens(db)`, which removes every token (in place, idempotent) for a compact
encoding (a 2,000-table schema's JSON: 4.6 MB → 1.9 MB).

### Measured conformance (this commit)

| Suite | Result |
|---|---|
| Holistics interpreter + nan snapshots (model JSON incl. every token, diagnostics) | 45/45 (100%), infos 45/45 |
| Holistics lexer / parser / validator / binder snapshots | 13/13, 17/17, 21/21, 25/25 |
| Holistics infinite_loops snapshot | terminates |
| Differential vs `@dbml/parse` 10.3.1 (`testdata/reference`, 751 inputs: fixtures, core, examples, docs, hand-written, all 16 operators, DiagramView, partials, Dep, multi-file, export-focused hand4) | 751/751 (100%) |
| `emit` round trip: parse → emit → parse model-identical (tokens ignored) over every compiling input in `testdata/reference` + `testdata/holistics` (multi-file excluded) | 660/660 (100%) |
| `emit` idempotence: emit(parse(emit(x))) == emit(x) | 660/660 (100%) |
| `export.SQL` vs `@dbml/core` 10.3.1 `exporter.export` (582 inputs; `.err` = both throw), modulo the systematic documented deviations (SQL Server `N'…'` extended properties + `[x]` index columns; PostgreSQL unquoted built-in types; MySQL cross-database foreign keys via `DATABASE()`) | mysql 578/582 (99.3%), postgres 577/582 (99.1%), mssql 577/582 (99.1%) — 13 known (`knownExportFailures`): quoted non-numeric record values (7), `int(--1)` type refused as unsafe (6), PostgreSQL `E'…'` backslash value (1). Byte-exact without the normaliser: mysql 558, postgres 536, mssql 467 |
| `export.CoreDBML` vs `exporter.export(src, "dbml")` | 579/582 (99.5%) — 3 known: Metadata key order |
| `importer.CoreDBML` vs `@dbml/core` 10.3.1 `importer.import` (Holistics importer+parser fixtures, hand-written, and @dbml/core's own SQL exports re-imported) | mysql 419/420 (99.8%), postgres 412/421 (97.9%), mssql 385/395 (97.5%), schemarb 8/8 (100%) |
| Clean mode (`ImportWith`, records on) over the same 1,244 inputs: imports and compiles | mysql 420/420, postgres 421/421, mssql 395/395, schemarb 8/8 (100%, gated) |
| Import corpus composition | mysql 38 Holistics + 8 hand + 374 exported; postgres 35 + 3 + 383; mssql 21 + 2 + 372; schemarb 5 + 3 |

"Exact" (parser suites) means: same diagnostics (level, code, start line/column) **and** the same
database JSON after dropping `null`/`[]`/`{}` values, with array order compared strictly. Export
and import goldens compare text exactly after trimming trailing whitespace per line; a golden
`.err` passes when go-dbml also returns an error (messages are not compared). Gates: export
≥95% per dialect, import ≥90% per format, plus no regressions outside `knownExportFailures` /
`knownImportFailures`.

### Remaining deviations from @dbml/parse (not covered by any golden)

- Diagnostic **messages** follow the reference wording where ported but are not compared; only
  level, code and start position are. Diagnostic *end* positions are not compared.
- `token.filepath` is a plain string (the reference serialises a `Filepath` object whose
  `absolute` string we emit).
- Datetime normalisation in records uses Go layouts instead of luxon; exotic ISO-8601 forms
  (week dates, ordinal dates, fractional hours) may be rejected where luxon accepts them.
- DiagramView validator *warnings* (e.g. unknown names) are not emitted; binding errors are.
- Tables/enums/groups/partials nested inside `Project` are not interpreted into
  `project.tables` etc. (the reference rejects them as an invalid context, too).
- `@dbml/core`'s `normalize()` / `export()` shapes are not produced (only `ResolveTable`).
- Multi-file edge cases beyond the `multi` goldens (deep `reuse` chains of schema imports,
  aliasing a whole schema) follow the documented semantics but are unverified against the
  reference.

### Known differences in export/import (Task 0.3)

- Exporters quote identifiers/literals containing quote characters or backslashes where the
  reference emits them raw, quote non-numeric record values in numeric columns / non-hex
  binary values (7 goldens), refuse unsafe raw fragments such as the type `int(--1)` (6
  goldens: the reference's output comments out the rest of the line), write PostgreSQL
  backslash values as `E'…'` (1 golden), SQL Server extended properties as `N'…'` and index
  columns as `[x]` (58 and 76 goldens; 104 SQL Server goldens touched), PostgreSQL built-in
  types like `VARCHAR(20)` / `character varying(255)` unquoted (39 goldens), and add generated columns, FULLTEXT /
  SPATIAL index kinds and MySQL prefix key parts (no golden), and qualify a MySQL foreign
  key from a named schema to a default-schema table with `DATABASE()` (14 goldens). `export_test.go` compares
  modulo the systematic ones (`documentedDeviations`) and lists the rest in
  `knownExportFailures`; it logs the byte-exact count too.
- `export.CoreDBML`: model metadata is a Go map, so keys inside generated `Metadata` blocks are
  sorted; `@dbml/core` keeps source order (3 goldens). Same for Dep custom keys.
- Parity importers are **lenient** where the reference grammar rejects expression syntax
  foreign to the dialect inside CHECK/DEFAULT/index expressions or reserved words used as names,
  and MariaDB's `CREATE OR REPLACE TABLE` (20 goldens, `knownImportFailures`). ANTLR error
  messages are not mirrored; malformed DDL may import partially where the reference throws.
- The MySQL importer drops `[ ] { } ? \` like the reference lexer, which reproduces texts like
  `` check: `status] IN (…)` `` in parity mode.
- `FromRealm` has no reference; it is tested against a hand-built realm only. Charset,
  collation, generated columns and partitions are ignored.
- Lexer trivia are re-lexed on demand (`Token.LeadingTrivia()` is a method, no longer a
  field); `lexer.Position` / `ast.NodeBase` fields are int32. These packages are internal to
  go-dbml's consumers (dbdiagram-api uses `dbml`, `model`, `diag`, whose APIs only grew).
- Compiles with more than `MaxElements` elements in one block (e.g. a 20k-column table) or
  more significant tokens than `MaxTokens` fail with `TOO_MANY_ELEMENTS` / `TOO_MANY_TOKENS`
  where `@dbml/parse` would compile them; raise the options when that is wanted.

## 7. Ecosystem & related repos

- Consumers (planned): `dbdiagram-api` (compile on save → `model_json` via
  `ParseFileContext` with its request budget, use `Result.Database`;
  store `dbml.Version` + `model.JSONSchemaVersion`, re-compile stored rows whose version is
  older than 2; `export.SQL` for SQL downloads;
  `importer.Import` for repo sync of `schema.sql`/`schema.rb`), `dbdiagram-cli` (Atlas
  inspectors → `atlasrealm.FromRealm` → `emit.DBML` → push), `dbdiagram-web` via the WASM build
  in a Web Worker (`goDbml.parse` with `{"tokens": false}` → `database`/`partial` for live previews, `importSQL`/`exportSQL` for
  the import/export UI). Browser editor intelligence still comes from `@dbml/parse`; on
  disagreement the server-side go-dbml result wins.
- No dependency on Forta/Keyring/Monitor — it is a pure library.

## 8. Operations

Not deployed. Released by git tag (`vX.Y.Z`) and consumed via the Go module proxy. Not yet
pushed to GitHub (local repo only at time of writing). Safe defaults mean an API can compile
user-submitted DBML directly; still apply an HTTP body limit and a request timeout upstream.
The WASM artifact (`dev wasm` → `dist/go-dbml.wasm` + `wasm_exec.js`; `dist/` is git-ignored)
is copied into `dbdiagram-web` and must be served with the `wasm_exec.js` of the same Go
version.

## 9. Rules & guardrails

- Stdlib only, except `ariga.io/atlas` in `importer/atlasrealm` (`//go:build !js`). Never add
  cgo or another module; never let `go.mod`'s go directive exceed 1.25 (check a dependency's
  go directive before bumping it — atlas ≥ v1.2.3 needs 1.26.4). Keep
  `GOOS=js GOARCH=wasm go build ./...` green and the WASM ≤ 3 MB gzipped (`dev wasm` prints it).
- Never edit `testdata/export/**` or `testdata/import/<format>/**` goldens by hand; regenerate
  them with `scripts/reference/export.mjs` / `import.mjs`. Hand-written importer inputs live in
  `testdata/import/inputs/<format>/`.
- Export/import behaviour must follow the published `@dbml/core` 10.3.1 goldens; a deliberate
  deviation goes into `knownExportFailures` / `knownImportFailures` with a reason and into §6.
- Never edit files under `testdata/holistics/` or `*.golden.json` by hand. Refresh Holistics
  fixtures by re-copying (update the SHA in both NOTICE files); regenerate goldens with
  `scripts/` (see `scripts/README.md`).
- Never loosen a suite's normalisation to make a case pass; fix the implementation or add the
  case to `knownReferenceFailures` with a reason and document it in §6. (`documentedDeviations`
  in `export_test.go` encodes deliberate, documented output changes only.)
- Every raw SQL fragment an exporter emits goes through `safeExpr`/`safeType`; any new
  fragment kind needs a case in `fragment_test.go` and in the integration smuggling test.
- Do not change `diag.Code` numbering — codes < 9000 must equal `CompileErrorCode`.
- Keep the provenance header on ported files; new ports get one too.
- Commit, push, merge, tag and release **only after asking Aiden and getting an explicit yes**
  for that specific action (see `~/CLAUDE.md` *Ask-first actions*). **Never** amend, rebase,
  force-push, or run destructive git ops.

## 10. Verification

All must pass before work is done:

```sh
gofmt -w -s . && gofmt -l .           # prints nothing
go vet ./...
go build ./...
GOOS=js GOARCH=wasm go build ./...
go vet ./... && GOOS=js GOARCH=wasm go vet ./...
go test -race -count=1 ./...          # conformance (≥95%), reference (≥95% + no regressions),
                                      # depth 10^6, perf (<1 s), limits, JSON round trip,
                                      # emit round trip + idempotence, export (≥95%/dialect),
                                      # import (≥90%/format), clean import (100%), escaping,
                                      # 20k tables, 16 MB dump, 200k nesting, atlas, wasm API
go test -count=1 -tags integration -run Against ./export   # Docker: sentinel survives every smuggling
                                      # input; adversarial DDL exact; round trips >= 90% (MySQL/PG)
go test -run='^$' -fuzz=FuzzParse -fuzztime=60s .
go test -run='^$' -fuzz=FuzzEmitRoundTrip -fuzztime=30s ./emit
go test -run='^$' -fuzz=FuzzTranslate -fuzztime=30s ./export
for f in FuzzImportMySQL FuzzImportPostgres FuzzImportMSSQL FuzzImportSchemaRB; do
  go test -run='^$' -fuzz=$f -fuzztime=30s ./importer; done
dev wasm                              # ≤ 3 MB gzipped
```

`limits_test.go`'s gates are wall-clock based (scaled under `-race`: 15 s perf, 40 s for the
10^6-deep inputs); on an overloaded machine re-run the root package alone before treating a
timing failure as a regression.

## 11. Keeping this file updated

Any change to structure, public API, JSON shape, limits, supported syntax, commands,
conformance suites or deviations updates this file (and README.md if setup/commands/role
changed) in the same change — not a follow-up.
