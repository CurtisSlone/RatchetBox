# Transcript: RLE — pinning an unambiguous spec (local model)

## Hypothesis

The earlier RLE wall was diagnosed not as a model-capability limit but as an **underspecified spec**.
The original spec said "compress runs of the same byte ... handling digits" without ever pinning the
encoding scheme, so run-length encoding of digit-containing strings is ambiguous ("31" = three `1`s, or
the literal bytes `3`,`1`?). Test and impl each picked a *different* disambiguation and disagreed exactly
on those cases. Prediction: **pinning the scheme should let the LOCAL model clear it — no frontier needed.**

## Setup

- New workspace `workspaces/rle2`, grounded test phase (test node now retrieves `stdlib/testing` and is
  nudged toward self-validating property assertions).
- Spec pinned to an **unambiguous byte-pair scheme**: each maximal run of a byte `b` of length `n`
  (1..255) encodes to exactly two bytes `byte(n), b`; runs > 255 split into ≤255 chunks; `Decode` reads
  `(count, data)` pairs and errors on odd length or a zero count. Digits are ordinary data bytes and can
  never be confused with counts. Concrete examples given in the spec (`Encode("aaa") == {3,'a'}`,
  `Encode("12") == {1,'1',1,'2'}`).

## What happened along the way (infra + flow hardening)

1. **WSL→Ollama networking broke mid-session.** `localhost:11434` stopped forwarding to the Windows-hosted
   Ollama (cause of an earlier `exit 2` crash — not a flow bug). The server was still reachable via the
   WSL gateway IP. Fix: a small Python TCP forwarder `localhost:11434 -> 172.18.160.1:11434`, leaving the
   git-tracked `ratchet.json` (`ollama_url: localhost`) untouched.
2. **Stub fumble, caught clean.** On one run the stub phase emitted the test functions
   (`TestRunLengthCodec`, `FuzzRoundTrip` using `*testing.T/F`) into `main.go` without importing
   `testing` → `vet: undefined: testing`. The stub gate failed the compile and **rolled back** (the
   "ship nothing broken" guarantee held). Hardened the stub prompt to explicitly forbid
   `Test*`/`Fuzz*`/`Benchmark*` and any `*testing.T/F/B` signature. Next run got past stub cleanly.

## Result

**The spec-pin worked for its target.** With the scheme pinned:

- The **test followed the byte-pair scheme exactly** — `wantEnc: string([]byte{3,'a'})`, digits as
  `string([]byte{1,'1',1,'2'})`, the 256-run split as `string([]byte{255,'a',1,'a'})`. No more textual
  "3a" ambiguity; the digit-disambiguation disagreement is gone.
- The **implementation is fully correct**: `Encode` is greedy run-length with the 255-cap split; `Decode`
  reads `(count,data)` pairs and returns errors on odd length and zero count. The local model wrote a
  correct codec once the scheme was unambiguous.

**One residual blocker — a structural test-table bug (not the impl).** The model built a single table-row
schema that runs `enc := Encode(tt.input)` on *every* row, including two rows that are really
*Decode-rejection* cases:

```
{ name: "odd length input", input: {3,'a',2}, wantEnc: "", wantErr: true }
{ name: "zero count",       input: {0,'a'},   wantEnc: "", wantErr: true }
```

`Encode` never errors, and `Encode("\x03a\x02")` is correctly `{1,3,1,'a',1,2}` (three distinct bytes,
each a run of 1), not `""`. So `enc != tt.wantEnc` fails — the row is **unsatisfiable**. Identical failure
on all four impl attempts (the impl can't satisfy a wrong test). The `{3,'a',2}` value is the spec's
example of an input **`Decode`** should reject; the model misapplied it to `Encode`.

## Lessons

1. **Spec-pin beats frontier for the ambiguity class — validated.** The local model produced a correct
   implementation the moment the scheme stopped being ambiguous. The original "RLE wall" was never a
   capability limit; it was a specification gap. Cost: a few lines of spec, no escalation.
2. **The red gate proves a test *fails against the stub*, not that its assertions are *satisfiable*.** A
   structurally-flawed table (running `Encode` on inputs meant only for `Decode` rejection, with a wrong
   `wantEnc`) passes red and then makes green unreachable for any correct impl. This is the same
   wrong-`want` failure mode, now narrowed to a single mis-structured table — and it remains the weak link.
3. **Test-authoring fix (local-model territory, not frontier):** separate "encode a plaintext → check the
   encoding + round-trip" rows from "decode a malformed byte string → expect an error" rows. Never run
   `Encode` on a decode-rejection input; `Encode` is total (never errors). Rely on the round-trip property
   for the happy path and a dedicated decode-error sub-test for rejection.

## Bottom line

RLE went the whole way: **scheme-ambiguity wall → (spec-pin) correct impl → (inverse-pair test rule)
correct test → (gofmt + func-main pipeline fixes) clean DONE** — the first wall cleared end-to-end on the
LOCAL model with **zero frontier escalation**.

### How the residual blockers resolved (each routed to a fix tier; none to the frontier)

1. **Test mis-wiring (3 distinct variants across reruns):** wrong literal `want` → running `Encode` on
   decode-rejection rows → calling `Decode` on the raw plaintext instead of `Encode(x)`. Each prompt tweak
   relocated the bug, which was itself the signal: hand-wiring a bidirectional-codec test is the local
   model's edge. The fix that removed the *root* (not whack-a-mole): an INVERSE-PAIR rule in the test
   prompt — validity is checked by the ROUND-TRIP only (`Inverse(Forward(x)) == x`); direct
   `Inverse(literal)` is reserved for error cases (`err != nil`, no value compare). After this the test was
   correct and **green passed on the first impl attempt**.
2. **Harden failed on two MECHANICAL issues**, not impl logic (green/fuzz/-race/vet/staticcheck/govulncheck
   all passed): the code wasn't `gofmt`-clean, and a pure library in `package main` (single-package rule)
   has no `func main`, so `go build` (which links a command) failed. Both fixed deterministically in the
   staging pipeline — the `prune_imports` pattern extended: `gofmt -w` the staged files, and synthesize a
   trivial `func main() {}` when a `package main` module declares none.

Final run: `reset → read → stub → stubwrite → test → red → impl → green → fuzz → harden → DONE` in one
pass. Independent check of the workspace: `gofmt -l` clean, `go vet` clean, `go build` OK, `go test -race`
ok.

The routing thesis held: every blocker was a spec fix, a prompt fix, or a deterministic mechanical fix, and
the workflow *localized* each one. The frontier stays reserved for genuine capability walls — heap remains
the lone such case (reconciling the public `Pop() (int,bool)` with `container/heap`'s `Pop() any`, which
grounding, spec-pinning, and test-grounding all cannot touch).
