Author the Go TEST file(s) for this module, from the specs' `behavior`. The implementation is still a
STUB (every body panics), so a correct test MUST compile and MUST FAIL right now - that is the red gate.

Rules:
- ONE test file at the module ROOT, named `<name>_test.go`, `package main` (the same single package as the
  stub). No subdirectories, no importing the module's own code - the test is IN package main with everything
  else. Call the EXACT signatures shown in the stub below.
- The test file contains ONLY `Test*`/`Fuzz*` functions and helpers. Do NOT declare `func main` - it
  already exists in main.go (declaring it again is a redeclaration error).
- A `Fuzz` target's `f.Fuzz`/`f.Add` arguments must be FUZZABLE PRIMITIVES ONLY: string, bool, int/int8/
  int16/int32/int64, uint/uintN, float32/float64, []byte. Do NOT use other types (e.g. `time.Duration`) as
  fuzz args - instead fuzz an int64 (e.g. milliseconds) and convert inside the body:
  `ttl := time.Duration(ms) * time.Millisecond`.
- Cover the core behavior with example tests (table-driven where natural). Assert real properties of the
  result, not just "no panic".
- Include AT LEAST ONE property/fuzz target: a `func FuzzXxx(f *testing.F)` that calls `f.Fuzz(func(t *testing.T, ...) {...})`
  asserting an invariant that holds for ALL inputs (e.g. round-trips, monotonicity, bounds), OR a
  `testing/quick` property. Seed it with `f.Add(...)`.
- For concurrency behavior, exercise it from many goroutines and assert the final invariant (run under -race).
- Standard library only (`testing`, `testing/quick`, `sync`, etc.). Include every import.
- Output ONLY marker-separated files: a line `=== <name>_test.go ===` then the body. No prose, no fences.

## Specs (behavior to test)
{{ specs }}

## Stub (call THESE exact signatures)
{{ stub }}

## Red-gate feedback on the previous attempt (empty on the first pass)
{{ feedback }}
