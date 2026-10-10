import { readFileSync, readdirSync, statSync, writeFileSync, mkdirSync } from 'node:fs';
import path from 'node:path';
import { createRequire } from 'node:module';
const [,, corpus, out, nm] = process.argv;
const require = createRequire(nm + '/package.json');
const { Compiler, MemoryProjectLayout, Filepath, CompileErrorCode } = require('@dbml/parse');
mkdirSync(out, { recursive: true });
const walk = (d, b = '') => readdirSync(d).flatMap(f => statSync(path.join(d, f)).isDirectory() ? walk(path.join(d, f), b + '/' + f) : [[b + '/' + f, readFileSync(path.join(d, f), 'utf-8')]]);
const dg = (level, d) => ({ level, code: CompileErrorCode[d.code], file: d.filepath?.absolute, start: [d.nodeOrToken.startPos.line, d.nodeOrToken.startPos.column], end: [d.nodeOrToken.endPos.line, d.nodeOrToken.endPos.column], msg: d.diagnostic });
for (const p of readdirSync(corpus).sort()) {
  const files = Object.fromEntries(walk(path.join(corpus, p)));
  const c = new Compiler(new MemoryProjectLayout(files));
  const r = c.interpretFile(Filepath.from('/main.dbml'));
  const errs = r.getErrors();
  const res = { database: errs.length ? null : r.getValue(), diags: [...errs.map(e => dg('error', e)), ...r.getWarnings().map(e => dg('warning', e)), ...(r.getInfos?.() ?? []).map(e => dg('info', e))] };
  writeFileSync(path.join(out, p + '.json'), JSON.stringify(res, (k, v) => v instanceof Filepath ? v.absolute : v, 2));
}
