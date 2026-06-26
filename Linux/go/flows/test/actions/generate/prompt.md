Write a COMPLETE Go implementation for the task below AND a Go test that proves it behaves correctly.
Both files compile and the tests must pass under `go vet` + `go test`.

Emit EXACTLY two files, each introduced by its marker line, in this order:

// === solution.go ===
package solution
// the implementation. A library package - do NOT write `func main`.
// Include every import you use; use every import and variable (Go rejects unused ones).

// === solution_test.go ===
package solution
// import "testing" and write at least one `func TestXxx(t *testing.T)` that exercises the
// implementation with real inputs and asserts the expected outputs with t.Errorf/t.Fatalf.
// Cover the normal case and at least one edge case. Do NOT write a trivial always-pass test.

Rules:
- Both files use `package solution` (same package, so the test can call unexported names).
- Standard library only unless the task says otherwise.
- Output ONLY the two marker lines and their Go source - no prose, no markdown fences.

## Task
{{ task }}

## Reference material (only the libraries the task needed were retrieved; sections may be empty)

### Go standard library
{{ stdlib_refs }}

### Design / algorithm pattern
{{ pattern_refs }}

### Idiomatic style (Effective Go / Code Review Comments)
{{ guideline_refs }}

### Pitfalls to avoid (builds-but-wrong)
{{ pitfall_refs }}

### Go idioms
{{ refs }}

Use the reference above when it is relevant; ignore an empty section. Prefer the exact signatures and
identifiers shown over recalled ones.
