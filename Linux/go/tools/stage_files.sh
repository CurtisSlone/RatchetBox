#!/usr/bin/env bash
# stage_files.sh - the TRANSACTIONAL multi-file oracle for coedit. Reads marker-separated Go files on
# stdin (=== path ===), snapshots every target, writes them all, then verifies the WHOLE module with
# `go vet ./...` + `go test -race ./...`. Keeps the change only if everything is clean; otherwise rolls
# back EVERY touched file (restoring prior contents, deleting files that were new). This makes a
# cross-cutting, multi-file change atomic - the thing single-file edit_file cannot do. Arg: proj; the
# marker payload on stdin.
set -u
proj="${1:?usage: stage_files <proj> (marker payload on stdin)}"
root="$proj"; [ -d "$root" ] || root="workspaces/$proj"
[ -d "$root" ] || { echo "no such workspace: $proj"; exit 1; }
payload="$(cat)"
snap="$(mktemp -d)"; trap 'rm -rf "$snap"' EXIT

written="$(printf '%s' "$payload" | ROOT="$root" SNAP="$snap" python3 -c '
import sys, os, re
root=os.environ["ROOT"]; snap=os.environ["SNAP"]
text=sys.stdin.read()
text="\n".join(l for l in text.splitlines() if not l.strip().startswith("```"))
marker=re.compile(r"^===\s*([A-Za-z0-9_./-]+\.go)\s*===\s*$")
cur=None; buf=[]; paths=[]
def flush():
    global cur, buf
    if cur is None: return
    name=os.path.normpath(cur)
    if name.startswith("..") or name.startswith("/"): return
    full=os.path.join(root, name)
    if os.path.exists(full):
        d=os.path.join(snap, os.path.dirname(name))
        if d: os.makedirs(d, exist_ok=True)
        open(os.path.join(snap, name), "w").write(open(full).read())
    if os.path.dirname(full): os.makedirs(os.path.dirname(full), exist_ok=True)
    open(full, "w").write("\n".join(buf).rstrip()+"\n")
    paths.append(name)
for line in text.splitlines():
    m=marker.match(line)
    if m: flush(); cur=m.group(1); buf=[]
    else: buf.append(line)
flush()
print(" ".join(paths))
')"
[ -z "$written" ] && { echo "no files parsed from the payload (need === path === markers)"; exit 1; }
echo "staged: $written"

# Strip unused imports the model left behind (goimports-lite) from the just-staged files before verifying.
for p in $written; do bash tools/prune_imports.sh "$root/$p" >/dev/null 2>&1 || true; done

RACE=""
if [ "$(go env CGO_ENABLED 2>/dev/null)" != "0" ] && { command -v gcc >/dev/null 2>&1 || command -v clang >/dev/null 2>&1 || command -v cc >/dev/null 2>&1; }; then RACE="-race"; fi
v="$(cd "$root" && GOFLAGS=-mod=mod go vet ./... 2>&1)"; vs=$?
ts=0; to=""
if [ "$vs" -eq 0 ]; then to="$(cd "$root" && GOFLAGS=-mod=mod go test $RACE ./... 2>&1)"; ts=$?; fi

if [ "$vs" -eq 0 ] && [ "$ts" -eq 0 ]; then
  echo "OK: ${written} staged; module vets and tests pass${RACE:+ (-race)}"
  exit 0
fi
echo "VERIFY FAILED - rolling back all ${written}:"
[ "$vs" -ne 0 ] && printf '%s\n' "$v" || printf '%s\n' "$to"
for p in $written; do
  if [ -f "$snap/$p" ]; then cp "$snap/$p" "$root/$p"; else rm -f "$root/$p"; fi
done
echo "rolled back (workspace restored)."
exit 1
