Write a COMPLETE Go file for ONE unit of a larger system. Output ONLY the Go source - no prose, no
markdown fences.

KEYWORD TAGS (so the code can be searched by intent - do this in EVERY file):
- Directly after the `package` line, emit `// file-kw: <8-14 lowercase keywords summarizing this file's purpose>`.
- Above EACH top-level func and type, emit `// kw: <4-8 lowercase keywords for what THIS symbol does>` -
  its action plus domain (e.g. `// kw: append event log write-ahead durable`). Be concise and specific;
  these replace prose comments and let later steps find the exact code.

The file you are writing is: {{ path }}

PACKAGE (derive from the path):
- A file at the module ROOT (no slash in the path, e.g. `main.go`, `store.go`) is `package main`.
- A file at `<dir>/<name>.go` is `package <dir>` (the last path element of the directory).
- Only the root `main.go` defines `func main`. No other file may define `func main`.

CALLING OTHER CODE:
- SAME package (same directory): call its types/functions directly - do NOT import or redeclare them.
- ANOTHER package (different directory): add `import "<module>/<dir>"` (the module path is in go.mod
  in "THE MODULE SO FAR" below) and call its EXPORTED (capitalized) names as `<pkg>.Name`. Never import
  your own package; never create an import cycle.
- Use the EXISTING names in "API ALREADY IN THE MODULE" VERBATIM, with the exact import path each entry
  shows. Do not invent or rename them. If the spec disagrees with the API, the API wins.

Include every import you use; use every import and variable (Go rejects unused ones). Standard library
only unless the spec says otherwise.

## API ALREADY IN THE MODULE (authoritative names + how to import each)
{{ api }}

## THIS UNIT'S SPEC (what to build)
{{ spec }}

## THE MODULE SO FAR (PROJECT.md, go.mod with the module path, file tree)
{{ project }}

## Reference (retrieved for this spec; any section may be empty - use only if relevant)
### Go standard library
{{ stdlib_refs }}
### Design / algorithm pattern
{{ pattern_refs }}
### Pitfalls to avoid (builds-but-wrong)
{{ pitfall_refs }}
### Idiomatic style
{{ guideline_refs }}
### Recipe (playbook for this kind of unit - follow its API shape)
{{ recipe_refs }}
### Concurrency patterns (channels/select/fan-in/worker-pool/context; use for concurrent code)
{{ concurrency_refs }}
### Concurrent cache building blocks (theine: sharded store, reader-biased mutex, single-flight, timer-wheel TTL, TinyLFU)
{{ cache_refs }}
### Third-party module (already in this workspace's go.mod)
{{ dep_refs }}
