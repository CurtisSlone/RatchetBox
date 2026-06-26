#!/usr/bin/env bash
# spec_check.sh - the spec oracle: validate that marker-separated .spec drafts are well-formed. Reads on
# stdin one or more specs, each introduced by a line "=== <name>.spec ===". A spec is well-formed if it
# has a `name:` line, a `role:` line whose value is one of data/interface/component/behavior/gui/test
# (or empty), and at least one of `intent:`/`behavior:`/`api:`. Exit 0 iff EVERY spec is well-formed;
# otherwise print which spec failed and why, and exit 1. No side effects (validation only). Run via
# /do spec_check (specs on stdin), or used as the oracle inside write_specs and the `spec` flow.
set -u
python3 -c '
import sys, re
text = sys.stdin.read()
text = "\n".join(l for l in text.splitlines() if not l.strip().startswith("```"))
marker = re.compile(r"^===\s*([A-Za-z0-9_./-]+\.spec)\s*===\s*$")
specs, cur, buf = [], None, []
for line in text.splitlines():
    m = marker.match(line)
    if m:
        if cur is not None: specs.append((cur, "\n".join(buf)))
        cur, buf = m.group(1), []
    else:
        buf.append(line)
if cur is not None: specs.append((cur, "\n".join(buf)))
if not specs:
    specs = [("unit.spec", text)]   # no markers: treat the whole input as one spec
roles = {"data","interface","component","behavior","gui","test",""}
def has(body, key): return re.search(r"(?mi)^\s*"+key+r"\s*:", body) is not None
bad = []
for name, body in specs:
    errs = []
    if not has(body, "name"): errs.append("missing name:")
    rm = re.search(r"(?mi)^\s*role\s*:\s*(\S*)", body)
    if not rm: errs.append("missing role:")
    elif rm.group(1).lower() not in roles:
        errs.append("role %r not one of data/interface/component/behavior/gui/test" % rm.group(1))
    if not (has(body,"intent") or has(body,"behavior") or has(body,"api")):
        errs.append("needs at least one of intent:/behavior:/api:")
    if errs: bad.append((name, errs))
if bad:
    for name, errs in bad: print("INVALID %s: %s" % (name, "; ".join(errs)))
    sys.exit(1)
print("OK: %d spec(s) well-formed (%s)" % (len(specs), ", ".join(n for n,_ in specs)))
'
