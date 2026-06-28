#!/usr/bin/env bash
# tdd_red.sh - the RED oracle for the tdd flow. Writes the test file(s) (=== path.go === markers on
# stdin) into the workspace, then proves the test is MEANINGFUL: it must COMPILE against the stub
# (go vet) AND FAIL (go test exits non-zero because the stub bodies panic). A test that compiles AND
# PASSES against a stub tests nothing -> rejected. Exit 0 iff vet passes AND test fails (red achieved).
# Arg: proj (workspace). Marker payload on stdin (same === path.go === format as stage_files).
set -u
proj="${1:?usage: tdd_red <proj> (marker payload on stdin)}"
root="workspaces/$proj"; [ -d "$root" ] || root="$proj"
[ -d "$root" ] || { echo "no such workspace: $proj"; exit 1; }

written="$(ROOT="$root" python3 -c '
import sys, os, re
root=os.environ["ROOT"]; text=sys.stdin.read()
text="\n".join(l for l in text.splitlines() if not l.strip().startswith("```"))
marker=re.compile(r"^===\s*([A-Za-z0-9_./-]+\.go)\s*===\s*$")
cur=None; buf=[]; paths=[]
def flush():
    global cur, buf
    if cur is None: return
    name=os.path.normpath(cur)
    if name.startswith("..") or name.startswith("/"): return
    full=os.path.join(root, name)
    if os.path.dirname(full): os.makedirs(os.path.dirname(full), exist_ok=True)
    open(full,"w").write("\n".join(buf).rstrip()+"\n"); paths.append(name)
for line in text.splitlines():
    m=marker.match(line)
    if m: flush(); cur=m.group(1); buf=[]
    else: buf.append(line)
flush(); print(" ".join(paths))
')"
[ -z "$written" ] && { echo "no test files parsed (need === path.go === markers)"; exit 1; }
echo "wrote test: $written"
# Strip unused imports the model left in the test (goimports-lite) before the compile check.
for p in $written; do bash tools/prune_imports.sh "$root/$p" >/dev/null 2>&1 || true; done
cd "$root" || exit 1

# 1. must COMPILE (vet type-checks without requiring a linked binary)
if ! go vet ./... 2>/tmp/tdd_vet.$$; then
  echo "NOT RED (compile error). The test must be in the SAME package as the stub (do NOT import the"
  echo "module by path) and call its EXACT signatures. If the STUB's signatures are wrong for the"
  echo "behavior, the stub must change - not the test:"
  sed 's/^/  /' /tmp/tdd_vet.$$; rm -f /tmp/tdd_vet.$$; exit 1
fi
rm -f /tmp/tdd_vet.$$
# 2. must FAIL (the stub bodies panic -> a real test goes red)
if go test ./... >/tmp/tdd_test.$$ 2>&1; then
  echo "NOT RED (the test PASSES against the panic stub, so it asserts nothing real). Add concrete"
  echo "assertions on the returned values / invariants, not just 'does not panic':"
  tail -8 /tmp/tdd_test.$$ | sed 's/^/  /'; rm -f /tmp/tdd_test.$$; exit 1
fi
echo "RED: the test compiles and fails against the stub (as required)."
rm -f /tmp/tdd_test.$$; exit 0
