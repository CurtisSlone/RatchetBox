Author the Go TEST file(s) for this module, from the specs' `behavior`. The implementation is still a
STUB (every body panics), so a correct test MUST compile and MUST FAIL right now - that is the red gate.

Rules:
- ONE test file at the module ROOT, named `<name>_test.go`, `package main` (the same single package as the
  stub). No subdirectories, no importing the module's own code - the test is IN package main with everything
  else. Call the EXACT signatures shown in the stub below.
- The test file contains ONLY `Test*`/`Fuzz*` functions and helpers. Do NOT declare `func main` - it
  already exists in main.go (declaring it again is a redeclaration error).
- A `Fuzz` target's `f.Fuzz`/`f.Add` arguments must be FUZZABLE PRIMITIVES ONLY: string, bool, int/int8/
  int16/int32/int64, uint/uintN, float32/float64, []byte. Do NOT use other types as fuzz args:
  - a richer SCALAR (e.g. `time.Duration`) -> fuzz an int64 and convert inside the body
    (`ttl := time.Duration(ms) * time.Millisecond`).
  - a SEQUENCE/collection (e.g. a `[]int` to push) -> you CANNOT fuzz `[]int`; fuzz a `[]byte` and DERIVE
    the sequence inside the body: `f.Fuzz(func(t *testing.T, data []byte) { for _, b := range data { h.Push(int(int8(b))) } ... })`.
    Use `int(int8(b))` if the values may be NEGATIVE (gives -128..127), or `int(b)` for 0..255. Seed with
    `f.Add([]byte{...})` using ONLY byte values 0-255 - never a negative literal in a `[]byte`.
- Cover the core behavior with example tests (table-driven where natural). Assert real properties of the
  result, not just "no panic".
- PREFER PROPERTY/INVARIANT assertions over hardcoded literal outputs. A property (round-trip
  `Decode(Encode(s)) == s`, monotonicity, bounds, commutativity) is SELF-VALIDATING - it cannot carry a
  wrong expected value. Reach for these first.
- If you DO assert a literal expected value (`got != want` with a constant `want`), you MUST hand-compute
  that `want` by walking the spec's algorithm yourself - a wrong `want` makes the test unsatisfiable and no
  correct implementation can pass it. When in doubt, assert the property instead of the literal.
- INVERSE-FUNCTION PAIRS (encode/decode, marshal/unmarshal, serialize/parse): check VALID data with the
  ROUND-TRIP ONLY - `Inverse(Forward(x)) == x` for plaintext `x`. Do NOT write a separate
  `Inverse(someLiteral) == want` row for valid data: you will mis-wire which value is the plaintext vs the
  encoded form (e.g. calling `Decode` on the raw plaintext instead of on `Encode(x)`). Reserve a direct
  `Inverse(literal)` call for ERROR cases ONLY, asserting `err != nil` with NO value comparison. So: one
  table of plaintext inputs checked by round-trip; one table of malformed encoded inputs checked for error.
- A TOTAL function (no `error` return, e.g. `func Encode(s string) string`) ALWAYS succeeds - it can NEVER
  be a rejection/error case. Rejection and "invalid input" cases belong ONLY to the function that RETURNS
  an error (e.g. `func Decode(s) (string, error)`). Keep them in SEPARATE tables/subtests: one table feeds
  VALID inputs through the transform + round-trip; a DIFFERENT table feeds MALFORMED inputs to ONLY the
  error-returning function and asserts `err != nil`. Never run the total function on a malformed/rejection
  input and expect it to fail or return "" - that row is unsatisfiable. A spec example like "Decode of an
  odd-length input returns an error" is a DECODE row; do not turn it into an Encode assertion.
- Include AT LEAST ONE property/fuzz target: a `func FuzzXxx(f *testing.F)` that calls `f.Fuzz(func(t *testing.T, ...) {...})`
  asserting an invariant that holds for ALL inputs (e.g. round-trips, monotonicity, bounds), OR a
  `testing/quick` property. Seed it with `f.Add(...)`.
- For concurrency behavior, exercise it from many goroutines and assert the final invariant (run under -race).
- Standard library only (`testing`, `testing/quick`, `sync`, etc.). Include every import.
- Output ONLY marker-separated files: a line `=== <name>_test.go ===` then the body. No prose, no fences.

REPAIR MODE: if a previous attempt and a verdict appear below, you are FIXING, not rewriting. Start from
your previous test and make the SMALLEST change that resolves EXACTLY what the verdict flags - keep every
other line identical. Do not regenerate from scratch (that reintroduces bugs you already fixed).

## Canonical Go testing idioms (table-driven, subtests, property/fuzz - follow these patterns)
{{ testing_refs }}

## Specs (behavior to test)
{{ specs }}

## Stub (call THESE exact signatures)
{{ stub }}

## Your PREVIOUS test attempt (empty on the first pass - edit THIS, do not restart)
{{ prev_test }}

## Red-gate verdict on that attempt (empty on the first pass - fix exactly this)
{{ feedback }}
