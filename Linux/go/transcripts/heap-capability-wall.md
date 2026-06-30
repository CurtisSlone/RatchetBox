# Transcript: Heap — a genuine capability wall (local model, FAIL)

**Plain-language summary:** We asked the small local model to build a min-heap (a priority queue) with a
simple public API: `Push(int)`, `Pop() (int, bool)`, `Len()`. It failed — the same way, four times in a
row. This is the one case in the sweep where the workflow could *not* be unblocked by a better spec,
better prompt, or mechanical fix. It is a real reasoning limit of the small model, and therefore the
justified place to escalate to a stronger model. The value here is that the workflow **localizes the wall
precisely** instead of flailing.

---

## 1. The input task (the spec)

```
name: Heap   role: component
api:
  - func New() *Heap
  - func (h *Heap) Push(x int)
  - func (h *Heap) Pop() (int, bool)   // smallest element; ok=false when empty
  - func (h *Heap) Len() int
behavior:
  - Pop returns the smallest pushed value; popped order is non-decreasing.
  - The popped multiset equals the pushed multiset.
```

**Input reasoning:** the spec is unambiguous (so this is *not* the RLE ambiguity problem) and the grounding
included the canonical `container/heap` wrapper pattern AND a from-scratch array heap. So the model had a
clear contract and the right knowledge in front of it. What's left to test is pure *capability*.

## 2. What the model produced (output) — the same mistake every attempt

The test it wrote is fine — it asserts the real properties:

```go
h := New(); h.Push(3); h.Push(1); h.Push(2)
got, ok := h.Pop()              // want 1, true
// ... and: popped order non-decreasing; popped multiset == pushed multiset
```

The implementation is where it breaks. Final attempt (identical shape to all four):

```go
import "container/heap"
type Heap struct{ ints []int }

func New() *Heap { h := &Heap{}; heap.Init(h); return h }      // <-- compile error originates here

func (h *Heap) Push(x int)        { heap.Push(h, x) }
func (h *Heap) Pop() (int, bool)  {                            // the PUBLIC API the test calls
    if h.Len() == 0 { return 0, false }
    return heap.Pop(h).(int), true
}
func (h *Heap) Len() int          { return len(h.ints) }
func (h *Heap) Less(i, j int) bool{ return h.ints[i] < h.ints[j] }
func (h *Heap) Swap(i, j int)     { h.ints[i], h.ints[j] = h.ints[j], h.ints[i] }
func (h *Heap) Push2(x interface{}) { h.ints = append(h.ints, x.(int)) }   // tried to satisfy the interface
func (h *Heap) Pop2() interface{}   { /* remove last */ }                  // ...but with the WRONG name
```

## 3. The oracle's verdict — identical all four attempts

```
vet: ./heap.go:11:12: cannot use h (*Heap) as heap.Interface value in heap.Init:
     *Heap does not implement heap.Interface (wrong type for method Pop)
         have Pop() (int, bool)
         want Pop() any
```

## 4. What the failure actually is (plain language)

`container/heap` is a *helper*: you give it a type that implements its 5-method interface
(`Len, Less, Swap, Push(any), Pop() any`) and it runs the heap algorithm for you. The trap: that interface
requires a method **named `Pop` that returns `any`** — but our public API requires a method **named `Pop`
that returns `(int, bool)`**. One type cannot have two methods with the same name and different signatures.

A senior engineer resolves this instantly with the **two-type wrapper**: an unexported inner type
(`intHeap []int`) implements `heap.Interface` with `Pop() any`, and the public `Heap` *wraps* it, exposing
`Pop() (int, bool)` that calls `heap.Pop(&inner).(int)`. (Or: skip `container/heap` entirely and write the
~15-line sift-up/sift-down array heap by hand — also in the grounding.)

The small model can see it needs a `Pop() any` for the interface — that's why it invented `Pop2`. But it
**cannot perform the restructuring**: it can't split into two types, and naming the second method `Pop2`
doesn't satisfy the interface (which demands the literal name `Pop`). It is stuck between the two
requirements and produces the same non-compiling shape every attempt.

## 5. Why a better spec / prompt / mechanical fix can't save this one

- The spec is already unambiguous → spec-pinning (which cleared RLE) does nothing here.
- The grounding already contains both correct solutions (the wrapper and the from-scratch heap) → adding
  knowledge does nothing; we proved the entry reaches the prompt and is ignored.
- It is a compile-time type error, not formatting or imports → no deterministic pipeline fix applies.

The repair loop fed the exact `vet` error back four times; the model could not act on it. **This is the
definition of a capability wall**, and exactly the case the "escalate-on-repair" design is for: after N
identical failures on one node, hand *just this node* to a stronger model, which applies the wrapper and
the gate passes.

## 6. Bottom line

Heap is the lone genuine **frontier** case from the sweep. Counter and cache (spec/test/mechanical-class
problems) cleared on the local model; heap (a structural-reasoning problem) did not, and won't, regardless
of spec or prompt work. The workflow's contribution is that it **names the wall precisely** — same error,
same line, four times — so escalation is targeted and cheap rather than a guess.
