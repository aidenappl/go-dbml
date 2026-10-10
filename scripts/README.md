# scripts (dev-only)

Nothing here runs during `go test`; tests only read the committed goldens.

## reference/ — regenerate `testdata/reference` goldens

`testdata/reference/<set>/` pairs DBML inputs with the output of the reference
implementation, `@dbml/parse` **10.3.1**. `reference_test.go` compares go-dbml
against them (diagnostic level/code/start position + the database JSON with
null/empty values stripped).

To add or refresh a set (requires Node 20+, run from this directory):

```sh
cd scripts/reference
npm install                       # installs @dbml/parse@10.3.1 locally (not committed)
# single-file corpus: a directory of *.dbml files
node ref.mjs   <corpus_dir> /tmp/out_<set>   "$PWD/node_modules/.."
# multi-file corpus: a directory of project dirs, each with /main.dbml
node multi.mjs <corpus_dir> /tmp/out_<set>   "$PWD/node_modules/.."
python3 golden.py <corpus_dir> /tmp/out_<set> ../../testdata/reference/<set>
cd ../.. && go test -run ReferenceSuite -v .
```

The third argument to `ref.mjs`/`multi.mjs` is the directory whose
`node_modules` holds `@dbml/parse`.

When bumping the reference version, regenerate **every** set, review the
diff, and update the version in `reference_test.go`, `AGENTS.md` and
`NOTICE`. Never hand-edit a golden file.

Sets: `fixtures` (Holistics interpreter snapshots), `core` (DBML from
@dbml/core fixtures), `examples` (snippets extracted from @dbml/parse example
tests), `docs` (snippets from the DBML docs), `hand*` (hand-written edge
cases), `ops` (all 16 relationship operators, DiagramView, partials, Dep),
`multi` (multi-file projects), `hand4` (export-focused: schemas, enums,
partial type binding, records, defaults, indexes, junction-name collisions).

## reference/export.mjs — regenerate `testdata/export` goldens

Runs `@dbml/core` **10.3.1** `exporter.export(src, dialect)` for `mysql`,
`postgres`, `mssql` and `dbml` over every single-file input in
`testdata/reference/<set>/` that compiles, plus the vendored Holistics
exporter fixtures (`testdata/holistics/dbml-core/exporter/*/input`, set
`holistics`). Writes `testdata/export/<dialect>/<set>/<name>.sql`, or `.err`
when @dbml/core throws while building its model. `export/export_test.go`
compares them (trailing whitespace trimmed).

```sh
cd scripts/reference && npm install     # @dbml/core@10.3.1 + @dbml/parse@10.3.1
node export.mjs "$PWD/../.." "$PWD"     # <repoRoot> <dir whose node_modules has @dbml/core>
cd ../.. && go test -run TestExportGoldens -v ./export
```

## reference/import.mjs — regenerate `testdata/import` goldens

Runs `@dbml/core` 10.3.1 `importer.import(src, format)` for `mysql`,
`postgres`, `mssql` and `schemarb` over the vendored Holistics importer and
parser fixtures (sets `importer`, `parser`), the hand-written inputs in
`testdata/import/inputs/<format>/` (set `hand`) and — run `export.mjs` first —
every distinct SQL file in `testdata/export/<format>/` (set `exported`).
Writes `testdata/import/<format>/<set>/<name>.dbml`, or `.err` when @dbml/core
throws. `importer/importer_test.go` compares them.

```sh
node import.mjs "$PWD/../.." "$PWD"
cd ../.. && go test -run TestImportGoldens -v ./importer
```

Ground truth is the **published** npm package, not the repo source (they
differ, e.g. CREATE VIEW handling). Never hand-edit an export/import golden.
