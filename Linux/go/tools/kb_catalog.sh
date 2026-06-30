#!/usr/bin/env bash
# kb_catalog.sh - maintain kb/catalog.json: the high-level KB registry the ENGINE reads (one entry per
# knowledge library: name, path, doc-count, default, one-line subject). KBs are auto-discovered by
# walking kb/ for subdirs that contain a manifest.json, so the inventory + doc counts cannot drift from
# disk. catalog.json is the SOURCE OF TRUTH for the authored fields (name override e.g. idioms->"kb",
# default flag, summary): those are PRESERVED across regenerations by matching on path. Run after
# adding/removing a KB or editing a summary.
# Modes:
#   (default)  rewrite kb/catalog.json from disk (+ preserved authored fields)
#   --check    exit 1 if a KB on disk is missing from the catalog (or vice versa) or a summary is empty
#   --list     print compact "name: summary" lines (for injecting the catalog into a plan prompt)
set -u
mode="${1:-write}"
MODE="$mode" python3 - <<'PY'
import json, os, sys, glob

mode = os.environ.get("MODE", "write")

# Bootstrap subject lines (keyed by directory basename) - used ONLY for a brand-new KB with no summary
# in the catalog yet. After first write, kb/catalog.json is the source of truth (edit summaries there).
SEED = {
 "idioms":      "Hand-written Go idioms: buildable-code basics - errors, slices/maps, interfaces, zero values, defer.",
 "stdlib":      "Go standard library API reference (go doc): io, net/http, sync, context, encoding, and 170+ packages.",
 "patterns":    "Design & algorithm patterns: 23 GoF patterns plus authored algorithms (two-heap median, production HTTP server).",
 "deps":        "Third-party module API surface (go doc per added dependency).",
 "guidelines":  "Idiomatic Go style from Effective Go, Code Review Comments, and Go Proverbs.",
 "pitfalls":    "Builds-but-wrong traps: unused vars, nil-map writes, loop-var capture, channel deadlock, slice aliasing, RLock mutation, string(int).",
 "recipes":     "End-to-end app playbooks: unit breakdown plus the APIs/deps for a kind of program (e.g. an HTTP service).",
 "concurrency": "Concurrent Go: channel & goroutine patterns - fan-in/out, worker pool, select-timeout, context cancellation, bounded parallelism, semaphores, generators.",
 "cache":       "Concurrent in-memory cache building blocks: sharded store, reader-biased RWMutex, single-flight, timer-wheel TTL, W-TinyLFU, count-min sketch.",
 "dsa":         "Data structures & algorithms in Go: linked lists, stacks, queues, deques, ring buffers, sets, hash tables, heaps/priority queues, trees (BST/AVL/red-black/trie/segment/Fenwick), union-find, graphs (BFS/DFS/topological-sort/Dijkstra/MST), sorting, searching, number theory, hashing/checksums, encoding (RLE/base-N/varint) - worked Go implementations + tests.",
 "cqrs":        "Event-driven / CQRS building blocks: Message, Publisher/Subscriber, Router + middleware, command/event bus, processors, marshaler.",
}

CATALOG = "kb/catalog.json"

# Preserve authored fields from the existing catalog, keyed by path.
prev = {}
if os.path.exists(CATALOG):
    try:
        for e in json.load(open(CATALOG)).get("entries", []):
            prev[e["path"]] = e
    except Exception:
        pass

def doc_count(path):
    m = os.path.join(path, "manifest.json")
    if os.path.exists(m):
        try: return len(json.load(open(m)).get("entries", []))
        except Exception: pass
    return len([f for f in glob.glob(os.path.join(path, "*.md")) if os.path.basename(f).lower() != "readme.md"])

# Discover KB dirs: any kb/<x>/ that carries a manifest.json (catalog.json is a file, so it is excluded).
dirs = sorted(d for d in glob.glob("kb/*") if os.path.isfile(os.path.join(d, "manifest.json")))

entries, warns = [], []
for path in dirs:
    base = os.path.basename(path)
    p = prev.get(path, {})
    name = (p.get("name") or base).strip()
    summ = (p.get("summary") or SEED.get(base, "")).strip()
    if not summ: warns.append(f"empty summary for KB '{name}' ({path}) - add one in {CATALOG}")
    entries.append({
        "name": name, "path": path,
        "default": bool(p.get("default", False)),
        "docs": doc_count(path),
        "summary": summ,
    })

catalog = {
  "$comment": "High-level KB registry: one entry per knowledge library (name, path, doc-count, default, one-line subject). The ENGINE reads this as its knowledge registry (Knowledge() prefers it over ratchet.json knowledgeBases). KBs + doc-counts are auto-discovered from disk by tools/kb_catalog.sh; name/default/summary are authored here and preserved on regenerate.",
  "version": 1,
  "entries": entries,
}

if mode == "--list":
    for e in entries:
        print(f"{e['name']}: {e['summary']}")
    sys.exit(0)

if mode == "--check":
    bad = False
    if not os.path.exists(CATALOG):
        print(f"MISSING: {CATALOG} - run tools/kb_catalog.sh"); bad = True
    else:
        def struct(es): return [{k: e[k] for k in ("name", "path", "default", "docs")} for e in es]
        cur = json.load(open(CATALOG)).get("entries", [])
        if struct(cur) != struct(entries):
            print(f"STALE: {CATALOG} does not match the KB dirs on disk - run tools/kb_catalog.sh"); bad = True
    for w in warns: print("WARN:", w)
    if not bad and not warns: print("OK: catalog matches the KB dirs on disk; all summaries present")
    sys.exit(1 if (bad or warns) else 0)

with open(CATALOG, "w") as f:
    json.dump(catalog, f, indent=2); f.write("\n")
print(f"wrote {CATALOG} ({len(entries)} KBs, auto-discovered)")
for w in warns: print("WARN:", w)
PY
