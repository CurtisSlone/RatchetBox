# Transcript: Counter — a thread-safe counter (local model, DONE)

**Plain-language summary:** We asked the small local model (qwen3-coder, 7B, running on the machine) to
build a counter that many goroutines can increment at the same time without losing counts. It passed every
gate on the first implementation attempt — including the Go race detector. This is the kind of task the
hardened workflow now handles end-to-end with no help from a bigger model.

---

## 1. The input task (the spec) — and why it's written this way

The workflow reads `.spec` files, not code. The component spec we gave it:

```
name: Counter   role: component
api:
  - func NewCounter() *Counter
  - func (c *Counter) Inc()
  - func (c *Counter) Add(n int64)
  - func (c *Counter) Value() int64
behavior:
  - NewCounter returns a counter whose Value is 0.
  - Inc increases by exactly 1. Add(n) increases by exactly n (n may be negative).
  - CONCURRENCY INVARIANT: if G goroutines each call Inc K times and then all finish,
    Value equals G*K exactly. Safe for concurrent use (no data race under go test -race).
```

**Input reasoning:** the spec states the *invariant* (G×K exactly, no race) rather than dictating the
mechanism. That leaves the model free to pick `sync.Mutex` or `sync/atomic`, but pins the property the
oracle will actually check. The lesson from earlier walls — *pin the contract, not the prose* — is applied
here.

## 2. What the model was fed at the implementation step (input)

The implementation prompt was **16,004 characters**. Most of it is retrieved knowledge: the prompt carried
the data-structures grounding, the concurrency knowledge base, and explicitly the strings `sync.Mutex` and
`atomic`. So the model didn't have to *know* Go concurrency from memory — the relevant idioms were placed
in front of it.

## 3. What the model produced (output)

**Stub (signatures only, every body panics) — the type-check rung:**

```go
package main
type Counter struct{ /* empty; impl will add sync.Mutex or atomic */ }
func NewCounter() *Counter   { panic("TODO") }
func (c *Counter) Inc()      { panic("TODO") }
func (c *Counter) Add(n int64){ panic("TODO") }
func (c *Counter) Value() int64 { panic("TODO") }
```

**Test (written before the implementation exists) — note the concurrency case and the fuzz target:**

```go
t.Run("Concurrency", func(t *testing.T) {
    const goroutines, increments = 100, 1000
    c := NewCounter()
    var wg sync.WaitGroup
    wg.Add(goroutines)
    for i := 0; i < goroutines; i++ {
        go func() { defer wg.Done(); for j := 0; j < increments; j++ { c.Inc() } }()
    }
    wg.Wait()
    if got := c.Value(); got != goroutines*increments { t.Errorf(...) }
})

func FuzzCounterAdd(f *testing.F) {
    f.Fuzz(func(t *testing.T, data []byte) {
        c := NewCounter(); expected := int64(0)
        for _, b := range data { d := int64(int8(b)); c.Add(d); expected += d }
        if got := c.Value(); got != expected { t.Errorf(...) }
    })
}
```

The model used the exact idioms from the grounding: a `sync.WaitGroup` to wait for all goroutines, and the
`[]byte` → `int8` trick to fuzz a sequence of signed deltas (you cannot fuzz a `[]int64` directly).

**Implementation (the green rung):**

```go
import ("sync"; "sync/atomic")
type Counter struct{ value int64 }
func NewCounter() *Counter   { return &Counter{} }
func (c *Counter) Inc()       { atomic.AddInt64(&c.value, 1) }
func (c *Counter) Add(n int64){ atomic.AddInt64(&c.value, n) }
func (c *Counter) Value() int64 { return atomic.LoadInt64(&c.value) }
```

It chose lock-free atomics — correct and race-free. (It imported `sync` but never used it; the pipeline's
deterministic `prune_imports` step silently dropped the unused import before the compile, so the model's
small sloppiness never reached a gate.)

## 4. The oracle's verdicts (what the gates said back)

```
green  : OK: counter.go staged; module vets and tests pass (-race)
fuzz   : all targets clean (5s each)
harden : PRODUCTION-CLEAN — gofmt, go vet, go build, go test -race, staticcheck, govulncheck all pass
```

Flow path: `reset → read → stub → stubwrite → test → red → (red-cap retry) → test → red → impl → green →
fuzz → harden → DONE`, 14 steps.

## 5. Why it worked

- The spec pinned the **invariant** (G×K, no race), so the test could assert it directly.
- The grounding put the concurrency idioms (`WaitGroup`, atomics, the `[]byte`-fuzz trick) in front of the
  model, so it didn't rely on memory.
- The deterministic pipeline (`prune_imports`, `gofmt`) absorbed the mechanical slips.

**Bottom line:** a concurrency-correct counter, race-clean, first implementation attempt, **local model
only**. No frontier escalation.
