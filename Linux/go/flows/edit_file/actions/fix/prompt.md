Your previous edit of {{ path }} failed `go vet` or `go test`. Return the CORRECTED complete file. Fix
exactly what the diagnostics report. Output ONLY the Go source - no prose, no markdown fences.

- Return the WHOLE file. Keep `package main` and everything not affected by the request.
- An "undefined"/"redeclared" error means a name does not match the rest of the module - use the real
  names; do not import or redeclare same-package names. Remove any unused import or variable. A TEST
  FAILED means behavior is wrong - fix it, do not weaken any test.

## Requested change
{{ request }}

## Current (original) contents of {{ path }}
{{ current }}

## Diagnostics (go vet / go test over the whole module)
{{ errors }}

## Previous attempt
{{ prev }}
