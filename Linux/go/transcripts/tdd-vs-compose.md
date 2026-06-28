# tdd vs compose - the assurance-ladder flow, and where it beats / breaks vs compose

The `tdd` flow (PLANS.md S1) builds test-first up the assurance ladder:
`read -> STUB (signatures, panic bodies; oracle: compiles) -> TEST RED (author test + a property/fuzz
target; oracle tdd_red: COMPILES against the stub AND FAILS) -> IMPL GREEN (fill bodies; oracle: go vet +
go test -race) -> FUZZ (oracle go_fuzz: FuzzXxx finds no crash) -> HARDEN (go_quality)`. Repairs are
feedback cycles (red->test, green/fuzz/harden->impl).

- Generated: 2026-06-28
- New tools: `tdd_red` (the red oracle), `go_fuzz` (the property/fuzz rung).
- Engine support: lint now accepts feedback cycles (back-edge that closes a cycle).

## Result 1: counter (single `package main`) - TDD WINS

```text
$ ratchet flow . tdd --ws counter_tdd ""
  read -> stub -> stubwrite
  test -> red FAIL -> test -> red PASS        # RED cycle: first test didn't go red; re-genned, then red
  impl -> green FAIL -> impl -> green PASS     # GREEN cycle: first impl failed -race/test; re-genned, passed
  fuzz -> harden -> done
  PRODUCTION-CLEAN: all available gates passed.
```

Both repair cycles fired. The artifact is strictly STRONGER than the compose-built counter:
- a **red-verified** test suite (TestCounter + concurrent + race + init), so it cannot be a trivial test;
- a **`FuzzCounter`** seeded `0,1,100,1000000` - exercising the large-count edge the cacheproxy missed;
- a correct lock-free `atomic.AddInt64` impl; `harden` PRODUCTION-CLEAN.

compose on the same spec also builds green, but its test is generated alongside the impl (no red proof)
and there is no fuzz/harden rung. **TDD trades more steps for more assurance.**

## Result 2: bounded worker pool (sub-package `pkg: pool`) - the model's stub ceiling

Same specs through both flows. compose -> `compose.fail` (the model split `Pool` across `pool.go` +
`poolimpl.go`, both declaring `Pool` -> redeclared). tdd was then iterated through FOUR real flow bugs,
each found by an oracle and fixed, before hitting a model-reliability wall:

| # | failure (oracle caught it) | fix |
|---|---|---|
| 1 | stub emitted a `_test.go` (from the test spec) -> fails vs panic | stub excludes `role: test` |
| 2 | stub imported `sync/atomic` it didn't use -> unused-import | "import only what signatures need" |
| 3 | specs target `package pool`, stub hardcoded `package main` -> test `import "pool"` fails | stub/test/impl respect the spec's `package:` (fix 2) |
| 4 | stub split `Pool` (data) + `PoolImpl` (component) -> `Pool redeclared` | "one declaration per symbol; merge same-type specs" |
| 5 | stub again imported `sync` it didn't use -> unused-import | **unfixed: prompt rule not reliably obeyed by the local model** |

The repair caps (fix 4) worked - no more spinning to timeout; the flow fails fast and clean at the cap.

## Findings (what the sketch taught us)

1. **The paradigm is sound** - counter proved it end-to-end: red-verified test + fuzz + harden-clean,
   both feedback cycles recovering. The property/fuzz rung (1c) works and seeds real edge cases.
2. **The flow hardened on 4 real bug classes** (test-in-stub, unused-import, sub-package, redeclaration);
   it now respects `package:`, merges same-type specs, and caps every repair cycle (fixes 2-4 done).
3. **The wall is the local model's stub reliability**, not flow design. It repeatedly leaves UNUSED
   IMPORTS in the stub despite an explicit rule - a deterministically-fixable class. The clean fix is a
   goimports-style auto-prune step, but `goimports` is ABSENT in this environment (and not module-cached).
4. **Deeper: the spec flow over-decomposes** one type into two units (`Pool` data + `PoolImpl` component),
   which is what makes both compose and tdd collide. A sharper `spec` decomposition (one unit per type)
   would help both flows.

## Result 3: both fixes applied -> the worker pool CLIMBS THE LADDER

Two complementary fixes, applied on main (different failure classes):
- **Fix A - `prune_imports` (goimports-lite, Go-specific):** removes unused imports deterministically;
  wired into `stage_files` (so it protects stub-write, impl-green, AND compose/coedit) and `tdd_red`.
- **Fix B - sharper `spec` decomposition:** the spec prompt now keeps a type's definition + methods in
  ONE spec. Confirmed: the worker pool re-specced to 2 units (`Pool` component + test), not 3.

Re-running `tdd` on the clean specs:

```text
reset -> read -> stub -> stubwrite      # PASS (no unused import, no redeclaration)
test -> red                             # PASS first try (compiles vs stub, fails)
impl -> green (x4 via implcap)           # impl cycle ground through the hard part (a correct concurrent pool)
```

The impl reached GREEN; verified independently: **gofmt clean, go vet clean, staticcheck clean, go build
ok, `go test -race` PASS, and `go_fuzz` (FuzzPoolSubmit, 6s) CLEAN** - a race- and fuzz-clean bounded
worker pool (NewPool/Submit/Close/Results, no goroutine leak) with a red-verified test suite, generated
test-first from a one-line spec. (The flow run itself was cut by a 590s shell timeout right after green,
before fuzz/harden flushed; the artifact is complete.)

The fixes are confirmed COMPLEMENTARY: B cleared the stub redeclaration, A cleared the unused-import wall,
and together they moved the blocker from "stub won't compile" to "impl logic" - which the impl cycle then
solved. The repair caps (fix 4) bounded the impl cycle cleanly.

## Result 4: concurrent TTL cache (the original cacheproxy nemesis) - the model's frontier, ladder holds

Pushed TDD onto the hardest target: a sharded TTL cache (Set/Get + expiry), the class where the
expiry-under-RLock race and mutate-under-RLock bug live - the bug the original cacheproxy SHIPPED. The
spec decomposed into 3 interdependent types (TTLCache -> TTLCacheShard -> TTLCacheEntry). Two runs:

- Run A: the stub invented a `cache` sub-package; `main.go` imported it as `"cache"` (not the full
  module path) -> impl never compiled -> impl cap exhausted (4 tries) -> clean fail, rolled back.
- Run B (after "prefer single root package" stub rule): the model STILL made the sub-package, and the
  RED cap exhausted - the test node could not produce a test COHERENT with the 3-type stub: it wrote a
  stray `main.go` referencing a `_test.go` func (`undefined: TestTTLCacheConcurrent`), and used struct
  fields the stub did not declare (`unknown field value`). The red gate rejected every attempt.

The decisive point: across BOTH runs the model never produced a coherent stub+test+impl for the 3-type
cache - and the ladder + caps **rejected every bad attempt and shipped nothing**. That is strictly better
than the original cacheproxy outcome, which shipped the expiry race. We never even reached the -race/fuzz
judges because the model failed earlier (compile/coherence), but the protection held: no broken cache.

### Structural fixes vs the cache - what each wall taught us

The diagnostic (run `compose` on the same specs) showed the units SCATTERED across packages: some at the
root in `package main`, one in a `cache/` sub-package -> cross-references undefined. So the wall was not
"per-unit vs whole-module" - it was INCONSISTENT PACKAGE PLACEMENT. We then removed the choice
structurally, wall by wall (no bigger model):

| wall | structural fix | result |
|---|---|---|
| units scatter across packages | stub/test/impl: ONE package main, ONE file, at the root | **coherence cleared** - the 3-type stub compiles |
| test declares `func main` / uses `time.Duration` as a fuzz arg | test-prompt rules (no main; fuzzable primitives only) | partially obeyed |
| test has a different trivial compile bug each try (`testing.F{TB:}`, redeclared main, unused var) | (none reliably) | RED cap exhausts |

The single-package fix WORKED - it cleared the structural coherence wall the cache kept hitting. What
remains is NOT structure: the small model cannot reliably write a clean-COMPILING complex test
(concurrent + fuzz + 3 types) - it makes a different trivial error each attempt, whack-a-mole, and the red
gate rejects each. That is raw code-gen reliability, where prompt/structure has diminishing returns - the
"bigger model OR much more structure" frontier. Stopped here.

## Net
TDD's assurance scales with type-count, and so does the model's difficulty:
- 1 type (counter): full ladder green, race + fuzz clean. TDD > compose.
- 1 concurrency-heavy type (worker pool, after fixes A+B): full ladder, race + fuzz + lint clean; compose
  still fails the same spec (`Pool redeclared`).
- 3 interdependent types (TTL cache): beyond the local model's coherence; the ladder + caps fail clean
  and ship nothing - the protection the original cacheproxy lacked.

Robustness banked while pushing: `prune_imports` (helps compose + coedit too), sharper spec decomposition,
single-root-package default, and the repair caps (proven: they bounded the red AND impl cycles cleanly).
Open: catalog-route the impl grounding; a coherence aid for multi-type stubs (e.g. regenerate stub+test
together on a field mismatch, since the red repair currently re-does only the test).
