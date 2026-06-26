#!/usr/bin/env bash
# plan_units.sh - turn a compose build-plan JSON (on stdin) into the foreach worklist: one line per unit
# in dependency order, "<targetpath> <specfile>". Mapping:
#   - entry unit (role behavior/gui) -> main.go (always at the module root, package main)
#   - a unit with a `pkg` (subdirectory/package) -> <pkg>/<name>.go (package <pkg>)
#   - a test unit (role test, or name containing "test") -> [<pkg>/]<base>_test.go
#   - otherwise -> <name>.go at the module root (package main)
# An empty/"main" pkg keeps the unit at the root (flat single-package model). Sub-packages let the
# module have a real directory layout; the generated files import each other via "<module>/<pkg>".
# NOTE: uses `python3 -c` (not a heredoc) so the piped plan JSON stays on stdin.
set -u
python3 -c '
import sys, json, re
try:
    plan = json.load(sys.stdin)
except Exception as e:
    sys.stderr.write("plan_units: invalid plan JSON: %s\n" % e); sys.exit(1)
units = plan.get("units") or []
if not units:
    sys.stderr.write("plan_units: no units in plan\n"); sys.exit(1)
def slug(s):
    s = (s or "unit").lower().replace("-", "_").replace(" ", "_")
    s = re.sub(r"[^a-z0-9_]", "", s)
    return s or "unit"
def pkgdir(u):
    p = (u.get("pkg") or "").strip().strip("/").lower()
    p = re.sub(r"[^a-z0-9_/]", "", p)
    return "" if p in ("", "main") else p
seen_entry = False
for u in units:
    name = u.get("name", "unit")
    role = (u.get("role") or "").lower()
    spec = u.get("spec", "")
    s = slug(name)
    pkg = pkgdir(u)
    pre = (pkg + "/") if pkg else ""
    if role in ("behavior", "gui") and not seen_entry:
        target, seen_entry = "main.go", True   # entry is always root package main
    elif role == "test" or "test" in s:
        base = re.sub(r"_*test$", "", s) or "unit"
        target = pre + base + "_test.go"
    else:
        target = pre + s + ".go"
    print(target, spec)
'
