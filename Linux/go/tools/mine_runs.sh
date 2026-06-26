#!/usr/bin/env bash
# mine_runs.sh - the self-improvement engine. Scans runs/ for oracle FAILURES (the diagnostics that
# triggered repairs), clusters them into error classes, counts them, and flags which recurring classes
# are NOT yet covered by a kb/pitfalls entry - i.e. the best candidates for a new KB entry. This closes
# the evidence loop: what the oracle keeps rejecting becomes what the KB should teach. Run: /do mine_runs.
set -u
python3 - <<'PY'
import json, glob, re, os
# Error classes: a label + a regex matched against failed-oracle output, + keywords used to check
# whether kb/pitfalls already covers the class.
CLASSES = [
    ("unused-import",   r"imported and not used",                 ["unused", "import"]),
    ("unused-var",      r"declared (and )?not used",              ["unused", "variable"]),
    ("undefined-name",  r"undefined:",                            ["undefined"]),
    ("redeclared",      r"redeclared",                            ["redeclar"]),
    ("nil-map-write",   r"assignment to entry in nil map",        ["nil map"]),
    ("channel-deadlock",r"all goroutines are asleep|deadlock",    ["deadlock", "channel"]),
    ("data-race",       r"DATA RACE",                             ["race"]),
    ("printf-verb",     r"(Printf|Errorf|Sprintf).*(wrong type|no formatting directives|too few|too many)", ["printf", "format"]),
    ("type-mismatch",   r"cannot use .* as .* value|mismatched types", ["type"]),
    ("missing-return",  r"missing return",                        ["missing return"]),
    ("wrong-behavior",  r"--- FAIL|TEST FAILED",                  ["test"]),
    ("syntax-error",    r"syntax error|expected ",               ["syntax"]),
]
counts = {c[0]: 0 for c in CLASSES}
examples = {c[0]: "" for c in CLASSES}
runs_seen = 0
fails_seen = 0
for step in glob.glob("runs/*/step-*.json"):
    try:
        d = json.load(open(step))
    except Exception:
        continue
    if d.get("kind") != "action" or d.get("ok") is not False:
        continue
    fails_seen += 1
    out = d.get("output", "") or ""
    for label, pat, _ in CLASSES:
        m = re.search(pat, out)
        if m:
            counts[label] += 1
            if not examples[label]:
                # first matching line as the example
                for line in out.splitlines():
                    if re.search(pat, line):
                        examples[label] = line.strip()[:140]; break

# pitfalls coverage
pitfall_text = ""
for f in glob.glob("kb/pitfalls/*.md"):
    pitfall_text += open(f).read().lower() + "\n"
def covered(keywords):
    return all(k.lower() in pitfall_text for k in keywords) or any(k.lower() in pitfall_text for k in keywords)

ranked = sorted([c for c in CLASSES if counts[c[0]] > 0], key=lambda c: -counts[c[0]])
print("== mine_runs: recurring oracle failures across runs/ ==")
print("(scanned %d failed oracle steps)\n" % fails_seen)
if not ranked:
    print("no classified oracle failures found - nothing to mine yet.")
else:
    print("%-18s %6s  %-9s  %s" % ("class", "count", "covered?", "example"))
    print("-" * 90)
    candidates = []
    for label, pat, kw in ranked:
        cov = covered(kw)
        if not cov: candidates.append(label)
        print("%-18s %6d  %-9s  %s" % (label, counts[label], "yes" if cov else "NO", examples[label]))
    print()
    if candidates:
        print("KB CANDIDATES (recurring + not covered by kb/pitfalls): " + ", ".join(candidates))
        print("Draft one with:  ratchet flow . learn \"<the error class, e.g. %s>\"" % candidates[0])
    else:
        print("All recurring failure classes are already covered by kb/pitfalls. The loop is closed.")
PY
