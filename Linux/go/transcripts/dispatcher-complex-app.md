# Transcript: Durable webhook dispatcher — a complex multi-unit app (local model)

**Plain-language summary:** We asked the small local model (qwen3-coder 30B-A3B — a mixture-of-experts model
with ~3B *active* parameters) to build a real, complex application: a durable webhook/job dispatcher that
combines two production patterns — an append-only **event log / write-ahead log** (event sourcing) and a
**circuit breaker + retry** for the flaky external calls. It is six interdependent units, with concurrency.
The model built a coherent, building, *running* app — after one deterministic tooling fix — and then the
assurance ladder, applied to the units in isolation, caught a design smell the running app had hidden. No
frontier model was used at any point.

---

## The app (6 units, one `package main`, stdlib only)

| Unit | Pattern it carries |
|---|---|
| `Job` | the work item + a 3-value state enum |
| `EventLog` | append-only WAL: `Append`, `Events`, `Replay` (fold to current state) |
| `Breaker` | circuit-breaker state machine: closed / open / half-open |
| `RetryPolicy` | deterministic exponential backoff |
| `Dispatcher` | worker pool: breaker-gated delivery + retry, every outcome logged |
| `Main` | wires it together with a deliberately flaky deliverer |

Specs were written "sharp" per the methodology: pinned invariants, named edge cases, exact signatures.

## Round 1 — compose: a cross-unit coherence wall, fixed at the tooling layer

The `compose` flow builds units in dependency order, each checked against the real code already built. All
six units generated, but the whole-module build failed on **two cross-unit reference errors**:

```
dispatcher.go: undefined: EventTypeEnqueued
   (EventLog actually defines the constant EventEnqueued — the Dispatcher GUESSED the name)
main.go: *DemoDeliverer does not implement Deliverer
   have Deliver(Job) error      want Deliver(string, []byte) error
   (Main GUESSED the interface signature)
```

**Root cause — an incomplete binding, not a model limit.** Composition's whole point is to feed a new unit
the *real* API of the units already built so it calls them verbatim. But the tool that produces that API,
`module_api.sh`, was a one-liner: `grep -E '^(type |func )'`. That emits `type Deliverer interface {` (the
opening line) but **not the method signatures inside it**, and it skips **every `const`**. So the Dispatcher
and Main never saw `EventEnqueued` exists or what `Deliver` looks like — they had to guess, and guessed
wrong.

**Fix (deterministic, ~8 lines):** make `module_api` emit constants and the *full* body of type/interface
blocks. Before vs after, for the same module:

```
# before:  type Deliverer interface {        # after:  type Deliverer interface {
                                              #             Deliver(url string, payload []byte) error
                                              #         }
# (constants absent entirely)                 #         EventEnqueued  EventType = "enqueued"
```

Re-running compose with the richer binding: **all six units cohere, the module builds and vets clean, and
the app runs end-to-end** — 3 jobs delivered through the flaky deliverer, 8 events in the log. This fix
benefits every future compose run; it is residual #1 ("bind the real contract") closed once, in a tool.

## Round 2 — the assurance ladder on the units (what compose can't tell you)

Compose verifies *coherence + build + run*. It does **not** prove correctness — there were no tests. So we
put the two self-contained units through the full TDD ladder (red-gated test → impl → green → fuzz →
`-race` → harden).

**`RetryPolicy` → DONE, first attempt.** Pure deterministic math (`Backoff(attempt) = Base*Factor^(attempt-1)`,
the `ShouldRetry` boundary, the `MaxAttempts<=0` edge). Red-gated examples + a monotonicity fuzz property,
all green, race-clean, production-clean. Nothing to report — which is the point: well-specified pure logic
is squarely in the model's wheelhouse.

**`Breaker` → FAILED first, and the failure was the gem.** The red-gated test failed identically four times
on one subtest:

```
breaker_test.go: Half-open success closes: State() = "open"; want "closed"
```

The test set up the breaker as `Failure() → sleep past cooldown → Success()` and expected `closed`. But in
my spec, the `open → half-open` transition fired **only as a side effect of `Allow()`** — so without an
`Allow()` call, the breaker was still `open`, and `Success()` (which only acts on closed/half-open) did
nothing.

**Why compose hid it:** the running `Dispatcher` *always* calls `Allow()` before acting, so the lazy
transition always fired in practice. The happy-path integration masked a fragile design; testing the state
machine *in isolation* exposed it. **This is exactly what the assurance ladder is for.**

**The fix was a design change (residual #2), not a model problem.** A state machine whose transition only
fires inside one query method is fragile. I revised the spec so the cooldown transition is applied at the
*start of every method* — the model implemented it as a private helper:

```go
func (b *Breaker) transitionTime() {
    if b.state == "open" && time.Since(b.lastOpenTime) >= b.cooldown {
        b.state = "half-open"
    }
}
// ...called at the top of Allow(), Success(), and Failure()
```

Re-run with the corrected design: **DONE** — red gate, fuzz (random op-sequence → `State()` always valid,
5s clean), `-race`, vet, staticcheck, govulncheck all pass.

## Lessons

1. **Multi-unit coherence is a binding problem, not a capability ceiling.** Six interdependent units cohered
   once the binding surfaced the real constants and interface signatures. The wall localized to one
   ~8-line tool, fixed deterministically — the `prune_imports`/`gofmt` pattern again.
2. **Compose proves coherence + build + run; the ladder proves correctness — and they disagree.** The
   Breaker *ran fine* in the app and *failed* the isolated ladder. A passing integration can hide a unit
   that is wrong (or fragile) on a path the integration never stresses. If correctness matters, harden the
   units, do not trust "it builds and runs".
3. **The ladder is a design critic, not just a bug finder.** It did not find a coding bug in the Breaker —
   it found a *fragile design* (transition coupled to one query method). The fix was to the spec/design, and
   then the model implemented the better design cleanly.
4. **Every residual tier showed up in one app, none needed a bigger model:** a tooling-binding fix
   (coherence), a clean first pass (RetryPolicy), and a spec/design fix surfaced by the ladder (Breaker).
   The 30B-A3B local model built the whole thing.

## Status & the remaining gap

`Job`, `EventLog`, `Breaker`, `RetryPolicy`, `Dispatcher`, `Main` cohere, build, vet, and run; `Breaker` and
`RetryPolicy` are additionally proven through the full ladder. The open assurance gap is the **`Dispatcher`
itself** — it cannot be hardened in isolation (it needs Job/EventLog/Breaker/RetryPolicy present), so the
next step is a `role: test` spec added to the compose set (or a dedicated multi-unit test workspace) to put
the worker-pool concurrency under `-race` with a fake deliverer. That is where a cache-style latent
concurrency bug would surface, if there is one.
