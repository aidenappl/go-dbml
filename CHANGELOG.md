# Changelog

## 0.1.1

`importer/atlasrealm` fixes found by dbdiagram-cli against real MariaDB 11.4 and PostgreSQL 16
inspections:

- **Expressions are normalised per dialect.** CHECKs (table and column level), expression
  defaults, generated columns and expression index parts go through the clean importer's
  normaliser (new `importer.NormalizeExpr`). MySQL/MariaDB report them with backtick
  identifiers, which made `emit.DBML` fail with `ErrUnrepresentable` for any schema with a
  CHECK. Expressions that cannot be kept are left out with a warning.
- **Inline MySQL/MariaDB ENUMs are named `<table>_<column>_enum`** (`_2`, `_3`, … on a clash)
  instead of all collapsing into one enum called `enum`. Named PostgreSQL enum types are kept
  and shared.
- **Foreign keys outside the inspected schemas are left out with a warning** instead of
  producing DBML that does not compile.
- New `atlasrealm.Import(realm, opts) (*Result, error)` returns the model, the DBML text and
  the warnings; `FromRealm` keeps its signature. New `Options.Dialect` (detected when
  empty); the dialect is recorded as `Project { database_type }`.
- Generated columns (`schema.GeneratedExpr`) are imported as `Column.Generated`.
- `dbml.Version` is `0.1.1`.
- Tests: performance gates compare two input sizes in process CPU time (load-robust);
  absolute budgets moved behind `DBML_PERF_STRICT=1` (`dev perf`). The Docker integration
  harness probes MariaDB / PostgreSQL over TCP so it no longer races their init-time server.

## 0.1.0

Initial release: DBML compiler (`dbml`, `model`, `diag`), canonical emitter (`emit`),
`@dbml/core`-compatible exporters with fragment validation and type translation (`export`),
MySQL / PostgreSQL / SQL Server / schema.rb importers (`importer`), Atlas realm converter
(`importer/atlasrealm`) and the WASM build (`wasm`). See AGENTS.md §6.
