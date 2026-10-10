// Usage: node import.mjs <repoRoot> <dirContainingNodeModules>
//
// Regenerates testdata/import/<format>/<set>/<name>.dbml (or .err when
// @dbml/core throws) with @dbml/core importer.import(src, format) for
// format in mysql, postgres, mssql, schemarb. Inputs (vendored verbatim
// from holistics/dbml packages/dbml-core/__tests__/examples):
//   importer/<f>_importer/input/*.in.sql  -> set "importer"
//   parser/<f>-parse/input/*.in.{sql,rb}  -> set "parser"
//   testdata/import/inputs/<f>/*.{sql,rb} -> set "hand" (hand-written)
//   testdata/export/<f>/<set>/<name>.sql  -> set "exported" (@dbml/core's own SQL export, deduplicated)
import { readFileSync, readdirSync, writeFileSync, mkdirSync, rmSync, existsSync } from 'node:fs';
import path from 'node:path';
import { createRequire } from 'node:module';

const [,, root, nm] = process.argv;
const require = createRequire(nm + '/package.json');
const core = require('@dbml/core');

const formats = ['mysql', 'postgres', 'mssql', 'schemarb'];
const outRoot = path.join(root, 'testdata/import');
const hc = path.join(root, 'testdata/holistics/dbml-core');

let n = 0;
for (const f of formats) {
  rmSync(path.join(outRoot, f), { recursive: true, force: true });
  const sources = [
    ['importer', path.join(hc, 'importer', `${f}_importer`, 'input')],
    ['parser', path.join(hc, 'parser', `${f}-parse`, 'input')],
    ['hand', path.join(outRoot, 'inputs', f)],
  ];
  // round trip: every SQL file @dbml/core exported (testdata/export/<f>/<set>/<name>.sql),
  // deduplicated by content, imported back -> set "exported" (<set>__<name>).
  const exp = path.join(root, 'testdata/export', f);
  if (existsSync(exp)) {
    const seen = new Set();
    const out = path.join(outRoot, f, 'exported');
    for (const set of readdirSync(exp).sort()) {
      for (const file of readdirSync(path.join(exp, set)).sort()) {
        if (!file.endsWith('.sql')) continue;
        const src = readFileSync(path.join(exp, set, file), 'utf-8');
        if (seen.has(src) || src.trim() === '') continue;
        seen.add(src);
        mkdirSync(out, { recursive: true });
        n++;
        const name = `${set}__${file.replace(/\.sql$/, '')}`;
        try {
          writeFileSync(path.join(out, name + '.dbml'), core.importer.import(src, f));
        } catch (e) {
          const msg = e && e.diags ? e.diags.map((x) => x.message).join('\n') : String(e && e.message || e);
          writeFileSync(path.join(out, name + '.err'), msg + '\n');
        }
      }
    }
  }
  for (const [set, dir] of sources) {
    if (!existsSync(dir)) continue;
    for (const file of readdirSync(dir).sort()) {
      if (!/\.(sql|rb)$/.test(file)) continue;
      const name = file.replace(/(\.in)?\.(sql|rb)$/, '');
      const src = readFileSync(path.join(dir, file), 'utf-8');
      const out = path.join(outRoot, f, set);
      mkdirSync(out, { recursive: true });
      n++;
      try {
        writeFileSync(path.join(out, name + '.dbml'), core.importer.import(src, f));
      } catch (e) {
        const msg = e && e.diags ? e.diags.map((x) => x.message).join('\n') : String(e && e.message || e);
        writeFileSync(path.join(out, name + '.err'), msg + '\n');
      }
    }
  }
}
console.log(`imported ${n} inputs`);
