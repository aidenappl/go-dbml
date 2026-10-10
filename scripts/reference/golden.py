#!/usr/bin/env python3
"""Convert harness output (<out>/<set>/<name>.json) into committed goldens.

usage: golden.py <corpus_dir> <harness_out_dir> <testdata/reference/<set>>

Copies every <name>.dbml (or multi-file project directory) from corpus_dir
into the target set and writes <name>.golden.json (minified, sorted keys).
"""
import json, os, shutil, sys

corpus, out, dst = sys.argv[1:4]
os.makedirs(dst, exist_ok=True)
for name in sorted(os.listdir(corpus)):
    src = os.path.join(corpus, name)
    if os.path.isdir(src):
        shutil.copytree(src, os.path.join(dst, name), dirs_exist_ok=True)
        base = name
    elif name.endswith('.dbml'):
        shutil.copy(src, os.path.join(dst, name))
        base = name[:-5]
    else:
        continue
    with open(os.path.join(out, base + '.json')) as f:
        data = json.load(f)
    with open(os.path.join(dst, base + '.golden.json'), 'w') as f:
        json.dump(data, f, separators=(',', ':'), sort_keys=True)
