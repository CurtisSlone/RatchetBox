#!/usr/bin/env bash
# write_specs.sh - validate marker-separated .spec drafts (via spec_check) and, if ALL are well-formed,
# write each into <workspace>/specs/ (bootstrapping the workspace's go.mod + specs/ if absent, so a
# `spec` run alone produces a composable workspace). Atomic: writes nothing if any spec is invalid.
# Args: proj (workspace name or path); the marker-separated specs on stdin.
set -u
proj="${1:?usage: write_specs <proj> (specs on stdin)}"
payload="$(cat)"

# Oracle first: reject the whole batch if any spec is malformed (-> the flow repairs).
v="$(printf '%s' "$payload" | bash tools/spec_check.sh 2>&1)"
if [ $? -ne 0 ]; then echo "SPEC INVALID (nothing written):"; printf '%s\n' "$v"; exit 1; fi

root="$proj"; [ -d "$root" ] || root="workspaces/$proj"
mod="$(basename "$root")"
mkdir -p "$root/specs"
[ -f "$root/go.mod" ]    || printf 'module %s\n\ngo 1.21\n' "$mod" > "$root/go.mod"
[ -f "$root/PROJECT.md" ] || printf '# %s\n\nA composed Go module.\n\n## Units\n' "$mod" > "$root/PROJECT.md"

printf '%s' "$payload" | ROOT="$root" python3 -c '
import sys, re, os
root = os.environ["ROOT"]
text = sys.stdin.read()
text = "\n".join(l for l in text.splitlines() if not l.strip().startswith("```"))
marker = re.compile(r"^===\s*([A-Za-z0-9_./-]+\.spec)\s*===\s*$")
cur, buf, written = None, [], []
def flush():
    global cur, buf
    if cur is not None:
        name = os.path.basename(cur)   # no path escape
        open(os.path.join(root, "specs", name), "w").write("\n".join(buf).strip() + "\n")
        written.append(name)
for line in text.splitlines():
    m = marker.match(line)
    if m: flush(); cur, buf = m.group(1), []
    else: buf.append(line)
flush()
print("wrote %d spec(s) to %s/specs: %s" % (len(written), root, ", ".join(written)))
'
echo "OK: specs written; review them, then: ratchet flow . compose --ws $mod \"\""
exit 0
