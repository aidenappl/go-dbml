// Usage: node ref.mjs <corpusDir> <outDir> <dirContainingNodeModules>
import { readFileSync, readdirSync, writeFileSync, mkdirSync } from 'node:fs';
import path from 'node:path';
import { createRequire } from 'node:module';
const require = createRequire(process.argv[4] + '/package.json');
const P = require('@dbml/parse');
const { Compiler, MemoryProjectLayout, DEFAULT_ENTRY, Filepath, CompileErrorCode } = P;
const [,, corpus, out] = process.argv;
mkdirSync(out, { recursive: true });
function diag(level, d) {
  const n = d.nodeOrToken;
  return { level, code: CompileErrorCode[d.code],
    start: n ? [n.startPos.line, n.startPos.column] : null,
    end: n ? [n.endPos.line, n.endPos.column] : null, msg: d.diagnostic };
}
for (const f of readdirSync(corpus).filter(f => f.endsWith('.dbml')).sort()) {
  const src = readFileSync(path.join(corpus, f), 'utf-8');
  let res;
  try {
    const layout = new MemoryProjectLayout();
    layout.setSource(DEFAULT_ENTRY, src);
    const compiler = new Compiler(layout);
    const report = compiler.interpretFile(DEFAULT_ENTRY);
    const errors = report.getErrors();
    const db = errors.length > 0 ? null : report.getValue();
    res = { database: db ?? null,
      diags: [...errors.map(e => diag('error', e)), ...report.getWarnings().map(e => diag('warning', e)), ...(report.getInfos?.() ?? []).map(e => diag('info', e))] };
  } catch (e) { res = { crash: String(e && e.stack || e) }; }
  const json = JSON.stringify(res, (k, v) => (v instanceof Filepath ? v.absolute : (typeof v === 'bigint' ? Number(v) : v)), 2);
  writeFileSync(path.join(out, f.replace(/\.dbml$/, '.json')), json);
}
