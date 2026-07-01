#!/usr/bin/env bash
# spec_diff.sh - the layered-spec increment deriver. Given two FULL spec snapshots
# (workspaces/<proj>/spec/<from> and .../spec/<to>), compute which units the next evolve pass must
# (re)generate: every spec ADDED in <to> or whose content CHANGED vs <from>. Emit them as a
# foreach-ready worklist ("<targetpath> <specfile>") in dependency order, so add_unit regenerates
# only the delta against the already-built module. Unchanged units stay on disk (the workspace is the
# state). Removed units are NOT emitted here (deletion is handled by evolve.prune via spec_removed.sh).
#
# Determinism is the point: the changeset is DERIVED from the snapshots, never hand-authored, so it can
# never drift from the specs. Args: proj, layers ("<from> <to>", e.g. "L0 L1").
#
# Ordering: topological by an optional `dependsOn:` spec line, tie-broken by role rank
# (data<interface<component<behavior<gui<test) then name - so a changed unit always generates AFTER any
# (added or pre-existing) unit it depends on, and the entry behavior maps to main.go last.
set -u
proj="${1:?usage: spec_diff <proj> \"<from> <to>\" [regen|removed]}"
layers="${2:?usage: spec_diff <proj> \"<from> <to>\" [regen|removed]}"
mode="${3:-regen}"    # regen (default) = added+changed worklist; removed = target paths of dropped units
root="workspaces/$proj"; [ -d "$root" ] || root="$proj"
[ -d "$root" ] || { echo "no such workspace: $proj" >&2; exit 1; }

python3 -c '
import sys, os, re
root, layers, mode = sys.argv[1], sys.argv[2].split(), sys.argv[3]
if len(layers) != 2:
    sys.stderr.write("spec_diff: layers must be \"<from> <to>\" (e.g. \"L0 L1\")\n"); sys.exit(1)
frm, to = layers
fromdir = os.path.join(root, "spec", frm)
todir   = os.path.join(root, "spec", to)
if not os.path.isdir(todir):
    sys.stderr.write("spec_diff: no such layer dir: %s\n" % todir); sys.exit(1)

RANK = {"data":0, "interface":1, "component":2, "behavior":3, "gui":4, "test":5}

def slug(s):
    s = (s or "unit").lower().replace("-", "_").replace(" ", "_")
    s = re.sub(r"[^a-z0-9_]", "", s)
    return s or "unit"

def parse(path):
    name = os.path.splitext(os.path.basename(path))[0]
    role, pkg, deps = "component", "", []
    for line in open(path, encoding="utf-8", errors="replace"):
        m = re.match(r"\s*(name|role|package|dependsOn|depends_on)\s*:\s*(.*)", line)
        if not m: continue
        k, v = m.group(1), m.group(2).strip()
        if   k == "name" and v: name = v
        elif k == "role": role = v.lower() or "component"
        elif k == "package": pkg = v.strip().strip("/").lower()
        elif k in ("dependsOn", "depends_on"):
            deps = [d for d in re.split(r"[,\s]+", v) if d]
    if pkg in ("main",): pkg = ""
    pkg = re.sub(r"[^a-z0-9_/]", "", pkg)
    return {"name": name, "role": role, "pkg": pkg, "deps": [slug(d) for d in deps],
            "slug": slug(name), "spec": os.path.basename(path)}

def units(d):
    out = {}
    if os.path.isdir(d):
        for f in sorted(os.listdir(d)):
            if f.endswith(".spec"):
                u = parse(os.path.join(d, f)); out[u["spec"]] = u
    return out

tou  = units(todir)
frou = units(fromdir)

# changeset = added (in to, not in from) + changed (in both, content differs)
def differs(spec):
    a, b = os.path.join(fromdir, spec), os.path.join(todir, spec)
    if not os.path.exists(a): return True
    return open(a, "rb").read() != open(b, "rb").read()
changed = {s for s in tou if differs(s)}

# topological order over the FULL to-set (so deps that are unchanged still order before dependents),
# Kahn with a stable (rank, name) tiebreak; cycle-safe (falls back to rank order for any leftover).
by_slug = {u["slug"]: u for u in tou.values()}
indeg = {u["spec"]: 0 for u in tou.values()}
edges = {u["spec"]: [] for u in tou.values()}          # dep.spec -> [dependent.spec]
for u in tou.values():
    for d in u["deps"]:
        dep = by_slug.get(d)
        if dep and dep["spec"] != u["spec"]:
            edges[dep["spec"]].append(u["spec"]); indeg[u["spec"]] += 1
def keyf(spec):
    u = tou[spec]; return (RANK.get(u["role"], 2), u["slug"])
ready = sorted([s for s in tou if indeg[s] == 0], key=keyf)
order, placed = [], set()
while ready:
    s = ready.pop(0); order.append(s); placed.add(s)
    for nxt in edges[s]:
        indeg[nxt] -= 1
        if indeg[nxt] == 0:
            ready.append(nxt); ready.sort(key=keyf)
for s in sorted([s for s in tou if s not in placed], key=keyf):  # cycle fallback
    order.append(s)

def target(u, seen):
    pre = (u["pkg"] + "/") if u["pkg"] else ""
    if u["role"] in ("behavior", "gui") and not seen[0]:
        seen[0] = True; return "main.go"
    if u["role"] == "test" or "test" in u["slug"]:
        base = re.sub(r"_*test$", "", u["slug"]) or "unit"
        return pre + base + "_test.go"
    return pre + u["slug"] + ".go"

seen = [False]
for spec in order:                      # assign paths over the FULL ordered set (entry = first behavior)
    u = tou[spec]; u["target"] = target(u, seen)

if mode == "removed":
    # units present in <from> but absent in <to> -> their target .go files must be deleted. Map each
    # standalone (a removed entry unit would have been main.go; usually a removed unit is a component).
    rseen = [False]
    for spec in sorted(frou, key=lambda s: (RANK.get(frou[s]["role"], 2), frou[s]["slug"])):
        if spec not in tou:
            print(target(frou[spec], rseen), spec)
else:
    for spec in order:                  # emit ONLY the changed/added units, in dependency order
        if spec in changed:
            print(tou[spec]["target"], spec)
' "$root" "$layers" "$mode"
