Your previous Go file failed `go vet` or `go test`. Return a CORRECTED complete file for {{ path }}. Fix
exactly what the diagnostics report. Output ONLY the Go source - no prose, no markdown fences.

- Keep `package main`. Other files are in the same package: call their names directly, do not import or
  redeclare them. Only `main.go` defines `func main`.
- A "redeclared"/"undefined" error usually means a name does not match the existing API below - use the
  API's exact names. Remove any unused import or variable. A TEST FAILED means the logic is wrong - fix
  it, do not weaken any test.

## API ALREADY IN THE MODULE (call these verbatim)
{{ api }}

## Request
{{ request }}

## Diagnostics (go vet / go test over the whole module)
{{ errors }}

## Previous attempt
{{ prev }}
