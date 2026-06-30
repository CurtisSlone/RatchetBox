# tdd: the repair loop must feed the prior attempt (not just the verdict), + a three-domain replication

A correction to the earlier "the cache is at the small model's reliability frontier" finding, and a
replication of the assurance ladder across three different kinds of software.

## The bug: cold-start repair

The tdd repair loops were under-built. `tdd.test` (red loop) and `tdd.impl` (green/fuzz/harden loop) each
bound the oracle's VERDICT but NOT the failing artifact - so on every retry the model regenerated the
test/impl from scratch, aimed at the verdict, and reintroduced bugs it had already fixed (whack-a-mole).
This deviated from Ratchet's own established pattern: `edit_file.fix` binds `prev` (the previous attempt)
AND `errors`. The canonical repair gives the model its last code and the complaint; the tdd loops dropped
the code.

**Fix:** wire a self-binding into both looped nodes - `tdd.test` reads `tdd.test` (prev_test), `tdd.impl`
reads `tdd.impl` (prev_impl); empty on the first pass, the prior failing code on a retry. The self-edge is
lint-legal (it closes a feedback cycle). The prompts switch to REPAIR MODE: "start from your previous
attempt; make the SMALLEST change that fixes exactly what the oracle flagged; keep everything else."

## What it changed: the cache, re-run

Before (cold-start): the 3-type TTL cache RED-capped on trivial, varied compile errors - attributed to a
model "reliability frontier."

After (iterating on prev):
```text
red    -> PASS first try        (was: capped at 4 attempts)
green  -> PASS                  (compiling, race-clean 3-type sharded cache)
fuzz   -> FUZZ FAILURE in FuzzTTLCache on input 0xffffffffffffff9c (a NEGATIVE ttl)
         failing input saved to testdata/fuzz/... (replayable)
... impl cycle then capped trying to fix the negative-ttl bug while keeping green
```

The correction: the earlier "frontier" verdict was confounded by the cold-start repair *we* under-built.
With proper iteration the local model reaches GREEN, and the FUZZ rung catches the real edge bug (negative
ttl) the happy-path `-race` test passed over - the unexercised-path class the original cacheproxy shipped.
The genuine remaining frontier is fixing that fuzz-found bug while keeping green - the ideal escalate-on-
repair target (one impl node, send the failing impl + the fuzz repro, gate the result).

## Replication across domains - the ladder catches the edge bug

Same flow, different kinds of software, looking for the same result (the ladder catches what a naive impl
ships). Each caught it at a different rung:

| Domain | Axis | Where the ladder caught the bug |
|---|---|---|
| **TTL cache** | concurrency | **fuzz** - happy-path + `-race` passed; fuzz found a negative `ttl` |
| **run-length codec** | encoding round-trip | **green/example** - the red-verified test (spec named "digit characters") asserted `Encode("112233")` round-trips; naive count-prefix RLE is ambiguous on digits; green rejected all 4 impls |
| **min-heap** | ordering invariant | (see below) |

The RLE is the complement of the cache: the cache showed the FUZZ rung's value (an edge the example test
missed); the RLE showed the RED-VERIFIED EXAMPLE test's value (a test authored first from a spec that
names the edge catches it at green, before fuzz). Both: the unexercised-path bug a naive impl would ship is
rejected; nothing broken ships.

In both the cache and the RLE the local model's failure was IDENTICAL across attempts (negative-ttl
handling; the `Encode("112233")="212223"` ambiguity) - not whack-a-mole, a real design limit the small
model cannot iterate past. That is a clean, reproducible, model-can't-solve case with a concrete artifact
to escalate - exactly what escalate-on-repair (placement-router-spec S4) targets.

## The min-heap: fuzz-expressibility friction, then a real wall

The heap (a textbook structure) took four runs - not because heaps are hard, but because expressing a
SIGNED-INT-SEQUENCE fuzz under Go's primitive-only fuzz constraint is friction the small model could not
navigate. Each red-cap was a fixable, GENERAL prompt rule (all banked, all benefit every future tdd run):

1. fuzzed `[]int` (not allowed) -> rule: fuzz a `[]byte`, derive the sequence in the body.
2. spec over-split the test into `HeapTest` + `HeapFuzzTest` -> they collide (`TestHeap redeclared`) ->
   rule (spec flow): exactly ONE `role: test` unit, several `Test*`/`Fuzz*` funcs in one file.
3. seeded a negative byte in a `[]byte` literal -> rule: derive signed via `int(int8(b))`, seed 0-255 only.

With those fixed, the heap cleared RED first try and reached the IMPL phase - where it hit a GENUINE wall:
the model tried to delegate to `container/heap`, but its public `Pop() (int, bool)` cannot also satisfy
`heap.Interface.Pop() any` (a signature conflict). It produced the SAME error all 4 green attempts; green
rejected every one; failed clean. A real API-design confusion the small model can't iterate past - the
escalation target (a frontier model would wrap stdlib correctly or write a manual heap).

## Net: three domains, three axes, three walls, one pattern

| Domain | Axis | Ladder caught it at | Reproducible wall (the escalation target) |
|---|---|---|---|
| TTL cache | concurrency | **fuzz** (happy-path + `-race` passed) | negative `ttl` handling |
| run-length codec | encoding round-trip | **green/example** (red-verified, spec named the edge) | count-prefix encoding ambiguous on digits |
| min-heap | ordering / stdlib API | **green** | public `Pop() (int,bool)` vs `heap.Interface.Pop() any` |

Every domain: the assurance ladder catches the bug a naive impl would ship (or rejects code the model
can't get right); the failure is IDENTICAL across attempts (a real capability/design limit, not flaky
whack-a-mole); nothing broken ships; and each is a clean, reproducible, model-can't-solve case with a
concrete artifact to escalate. The `prev`-wiring fix is what let the model iterate at all (cache reached
green; all three reached the real wall instead of capping on trivia). Confirmed across concurrency,
encoding, and data-structure/stdlib axes.
