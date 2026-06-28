#!/usr/bin/env bash
# prune_imports.sh - remove UNUSED imports from Go files (a goimports-lite, since goimports is absent in
# this environment). Go-specific. Heuristic, and safe: an unused import is always removable.
#   - blank (`_ "x"`) and dot (`. "x"`) imports are KEPT (side-effects / dotted scope).
#   - an aliased import keeps if its alias is used as `alias.`; a plain import keeps if the package name
#     (the path's last element) is used as `name.` anywhere outside the import block.
# Args: one or more .go files and/or directories (dirs are walked for *.go). Edits in place.
set -u
[ $# -ge 1 ] || { echo "usage: prune_imports <file.go|dir> ..."; exit 0; }
python3 - "$@" <<'PY'
import sys, os, re, glob

def gofiles(args):
    out = []
    for a in args:
        if os.path.isdir(a):
            out += glob.glob(os.path.join(a, "**", "*.go"), recursive=True)
        elif a.endswith(".go") and os.path.isfile(a):
            out.append(a)
    return sorted(set(out))

BLOCK = re.compile(r'(?ms)^import\s*\(\n(.*?)\n\)[ \t]*$')
SINGLE = re.compile(r'(?m)^import[ \t]+((?:[A-Za-z0-9_.]+[ \t]+)?(?:_[ \t]+|\.[ \t]+)?"[^"]+")[ \t]*$')

def spec_name(spec):
    spec = spec.strip()
    m = re.match(r'^([A-Za-z0-9_.]+)[ \t]+"([^"]+)"$', spec)   # alias / _ / . + path
    if m:
        return m.group(1), m.group(2)
    m = re.match(r'^"([^"]+)"$', spec)                          # plain path
    if m:
        return m.group(1).split("/")[-1], m.group(1)
    return None, None

def used(name, body):
    if name in ("_", "."):
        return True
    return re.search(r'\b' + re.escape(name) + r'[ \t]*\.', body) is not None

def prune(fn):
    src = open(fn).read()
    orig = src
    mb = BLOCK.search(src)
    if mb:
        body = src[:mb.start()] + src[mb.end():]   # everything outside the block
        kept = []
        for ln in mb.group(1).split("\n"):
            s = ln.strip()
            if not s or s.startswith("//"):
                kept.append(ln); continue
            name, _ = spec_name(s)
            if name is None or used(name, body):
                kept.append(ln)
        nonblank = [l for l in kept if l.strip()]
        if not nonblank:
            src = src[:mb.start()] + src[mb.end():]
        else:
            src = src[:mb.start()] + "import (\n" + "\n".join(kept) + "\n)" + src[mb.end():]
    for m in list(SINGLE.finditer(src)):
        name, _ = spec_name(m.group(1))
        body = src[:m.start()] + src[m.end():]
        if name and not used(name, body):
            src = src[:m.start()] + src[m.end():]
    if src != orig:
        open(fn, "w").write(src)
        print("pruned imports:", fn)

for fn in gofiles(sys.argv[1:]):
    prune(fn)
PY
