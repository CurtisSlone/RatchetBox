Your previous Go failed `go vet` or `go test`. Return a CORRECTED pair of files. Fix exactly what the
diagnostics report - the failure may be in the implementation OR in the test.

Emit EXACTLY two marker-separated files, same as before:

// === solution.go ===
package solution
// the corrected implementation (library package, no `func main`; no unused imports/vars).

// === solution_test.go ===
package solution
// the corrected test (import "testing"; real assertions; not a trivial always-pass test).

- A `TEST FAILED` message means the code ran but produced the wrong result: fix the logic so the
  asserted behavior holds (do not weaken the test to make it pass).
- A `VET FAILED` or compile message means it did not build cleanly: fix the reported construct.
- Output ONLY the two marker lines and their Go source - no prose, no markdown fences.

## Task
{{ task }}

## Diagnostics
{{ errors }}

## Previous attempt
{{ prev }}
