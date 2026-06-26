You are making ONE coordinated change across SEVERAL Go files at once. Below are the current files, each
under a `=== path ===` marker. Apply the requested change CONSISTENTLY across all of them so the whole
module still compiles and its tests pass.

Output rules:
- Return EVERY file, each under its exact `=== path ===` marker, as the COMPLETE updated file (not a
  diff). Preserve everything not affected by the change (package clause, imports, other funcs/types).
- Keep all call sites in sync: if you change a function's signature, update every caller in these files.
- Add/remove imports to match. Only main.go may have func main.
- Output ONLY the marker blocks - no prose, no code fences.

## Requested change
{{ request }}

## Current files (apply the change across all of them)
{{ current }}

## Reference (may be empty - use only if relevant)
### Go standard library
{{ stdlib_refs }}
### Pitfalls
{{ pitfall_refs }}
