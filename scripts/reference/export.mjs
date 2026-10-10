// Usage: node export.mjs <repoRoot> <dirContainingNodeModules>
//
// Regenerates testdata/export/<dialect>/<set>/<name>.sql (or .err when
// @dbml/core throws while building its model) with
// @dbml/core exporter.export(src, dialect) for dialect in mysql, postgres,
// mssql. Inputs: every single-file testdata/reference/<set>/*.dbml that
// compiles without errors, plus the vendored Holistics exporter fixtures
// testdata/holistics/dbml-core/exporter/<d>_exporter/input/*.in.dbml
// (set "holistics", named <d>_exporter__<name>).
import { readFileSync, readdirSync, writeFileSync, mkdirSync, rmSync, existsSync } from 'node:fs';
import path from 'node:path';
import { createRequire } from 'node:module';

const [,, root, nm] = process.argv;
const require = createRequire(nm + '/package.json');
const core = require('@dbml/core');
const P = require('@dbml/parse');
const { Compiler, MemoryProjectLayout, DEFAULT_ENTRY } = P;

const dialects = ['mysql', 'postgres', 'mssql', 'dbml'];
const outRoot = path.join(root, 'testdata/export');
for (const d of dialects) rmSync(path.join(outRoot, d), { recursive: true, force: true });

function compiles(src) {
  const layout = new MemoryProjectLayout();
  layout.setSource(DEFAULT_ENTRY, src);
  const c = new Compiler(layout);
  return c.interpretFile(DEFAULT_ENTRY).getErrors().length === 0;
}

const inputs = []; // [set, name, path]
const refDir = path.join(root, 'testdata/reference');
for (const set of readdirSync(refDir).sort()) {
  if (set === 'multi') continue;
  for (const f of readdirSync(path.join(refDir, set)).sort()) {
    if (f.endsWith('.dbml')) inputs.push([set, f.replace(/\.dbml$/, ''), path.join(refDir, set, f)]);
  }
}
const hx = path.join(root, 'testdata/holistics/dbml-core/exporter');
for (const d of readdirSync(hx).sort()) {
  const dir = path.join(hx, d, 'input');
  if (!existsSync(dir)) continue;
  for (const f of readdirSync(dir).sort()) {
    if (f.endsWith('.in.dbml')) inputs.push(['holistics', `${d}__${f.replace(/\.in\.dbml$/, '')}`, path.join(dir, f)]);
  }
}

let n = 0;
for (const [set, name, p] of inputs) {
  const src = readFileSync(p, 'utf-8');
  if (!compiles(src)) continue;
  n++;
  for (const d of dialects) {
    const dir = path.join(outRoot, d, set);
    mkdirSync(dir, { recursive: true });
    try {
      writeFileSync(path.join(dir, name + '.sql'), core.exporter.export(src, d));
    } catch (e) {
      const msg = e && e.diags ? e.diags.map((x) => x.message).join('\n') : String(e && e.message || e);
      writeFileSync(path.join(dir, name + '.err'), msg + '\n');
    }
  }
}
console.log(`exported ${n} inputs x ${dialects.length} dialects`);
