Your previous Go unit broke the module build. Return a CORRECTED complete Go file for {{ path }}. Fix
exactly what the compiler reported. Output ONLY the Go source - no prose, no markdown fences.

- Package is set by the path: root file -> `package main`; `<dir>/x.go` -> `package <dir>`. Only the
  root `main.go` defines `func main`.
- Same-package code: call directly. Other-package code: `import "<module>/<dir>"` and use `<pkg>.Name`
  (exported). An "undefined"/"redeclared"/"cannot find package" error means a name or import path does
  not match the API below - use its exact names and import paths. Remove any unused import or variable.

## API ALREADY IN THE MODULE (call these verbatim)
{{ api }}

## THIS UNIT'S SPEC
{{ spec }}

## Build errors (from go build ./... over the whole module)
{{ errors }}

## Previous attempt
{{ prev }}
