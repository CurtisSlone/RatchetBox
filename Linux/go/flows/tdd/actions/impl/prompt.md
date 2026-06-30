Implement the module: replace the stub `panic("TODO")` bodies with real code so the existing test PASSES
(go vet + go test -race + the fuzz target). Do NOT change the signatures the test calls.

Rules:
- Output the SAME files as the stub (same `=== <path>.go ===` paths, all at the module ROOT, all
  `package main`), now with full implementations. ONE package, no subdirectories, no imports of your own
  code. Keep every exported signature identical; declare each type/function EXACTLY once.
- The test is the contract - satisfy exactly what it asserts, including the property/fuzz invariant and any
  concurrent (-race) expectations. Use the reference material for correct concurrency + to avoid traps.
- Standard library only. Include every import; remove unused ones.
- KEYWORD TAGS: keep a `// file-kw:` line after `package` and a `// kw:` line above EACH top-level func/type
  (4-8 lowercase keywords for what that symbol does - its action + domain). These make the code searchable
  by intent; never drop them.
- Output ONLY marker-separated files. No prose, no code fences.

REPAIR MODE: if a previous implementation and oracle feedback appear below, you are FIXING, not rewriting.
Start from your previous implementation and make the SMALLEST change that resolves EXACTLY what the oracle
flagged - keep everything else identical. Do not regenerate from scratch (that reintroduces fixed bugs).

## Specs
{{ specs }}

## The test you must satisfy
{{ test }}

## Stub (the signatures to keep)
{{ stub }}

## Your PREVIOUS implementation attempt (empty on the first pass - edit THIS, do not restart)
{{ prev_impl }}

## Reference (data structures & algorithms / concurrency / cache / pitfalls / stdlib)
### Data structures & algorithms (use this canonical implementation if it fits the task)
{{ dsa_refs }}
### Other
{{ conc_refs }}
{{ cache_refs }}
{{ pitfall_refs }}
{{ stdlib_refs }}

## Oracle feedback to fix (empty unless a prior attempt failed)
### go test / build (green)
{{ green_fb }}
### fuzz
{{ fuzz_fb }}
### harden (vet/staticcheck/govulncheck)
{{ harden_fb }}
