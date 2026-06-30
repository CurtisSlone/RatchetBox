# Priority Queue / Binary Heap (Array-Backed)

A binary heap is a complete binary tree stored compactly in an array where every parent satisfies the heap order property (in a min-heap each parent is <= its children, in a max-heap each parent is >= its children), so the minimum or maximum is always at index 0. It is the standard way to implement a priority queue: peek at the best element in O(1), and push (sift-up) or pop the best (sift-down) in O(log n); building a heap from n items with heapify is O(n). For index i the parent is (i-1)/2 and the children are 2i+1 and 2i+2, so no pointers are needed and space is O(n). Use it for Dijkstra/Prim/A* frontiers, scheduling by priority, k-largest/k-smallest selection, and merging sorted streams. Keywords: heap priority queue min-heap max-heap binary heap push pop peek top sift-up sift-down bubble-up bubble-down heapify build-heap insert extract-min extract-max delete decrease-key complete-binary-tree array-backed parent child index logarithmic dijkstra prim scheduling container/heap heap.Interface heap.Init heap.Push heap.Pop heap.Fix

## implementation

```go
package heapq

import "cmp"

// Heap is an array-backed binary heap. less defines the ordering:
// pass cmp.Less for a min-heap, or a reversed comparator for a max-heap.
type Heap[T any] struct {
	data []T
	less func(a, b T) bool
}

// New returns an empty heap ordered by less.
func New[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{less: less}
}

// NewMin returns a min-heap over an ordered type.
func NewMin[T cmp.Ordered]() *Heap[T] {
	return &Heap[T]{less: func(a, b T) bool { return a < b }}
}

// NewMax returns a max-heap over an ordered type.
func NewMax[T cmp.Ordered]() *Heap[T] {
	return &Heap[T]{less: func(a, b T) bool { return a > b }}
}

// FromSlice builds a heap from items in O(n) using bottom-up heapify.
func FromSlice[T any](items []T, less func(a, b T) bool) *Heap[T] {
	h := &Heap[T]{data: items, less: less}
	for i := len(h.data)/2 - 1; i >= 0; i-- {
		h.down(i)
	}
	return h
}

// Len reports the number of elements.
func (h *Heap[T]) Len() int { return len(h.data) }

// Peek returns the best element (min or max) without removing it, O(1).
func (h *Heap[T]) Peek() (T, bool) {
	var zero T
	if len(h.data) == 0 {
		return zero, false
	}
	return h.data[0], true
}

// Push inserts v and restores heap order in O(log n).
func (h *Heap[T]) Push(v T) {
	h.data = append(h.data, v)
	h.up(len(h.data) - 1)
}

// Pop removes and returns the best element in O(log n).
func (h *Heap[T]) Pop() (T, bool) {
	var zero T
	n := len(h.data)
	if n == 0 {
		return zero, false
	}
	top := h.data[0]
	h.data[0] = h.data[n-1]
	h.data[n-1] = zero
	h.data = h.data[:n-1]
	if len(h.data) > 0 {
		h.down(0)
	}
	return top, true
}

// up sifts the element at i toward the root.
func (h *Heap[T]) up(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if !h.less(h.data[i], h.data[parent]) {
			break
		}
		h.data[i], h.data[parent] = h.data[parent], h.data[i]
		i = parent
	}
}

// down sifts the element at i toward the leaves.
func (h *Heap[T]) down(i int) {
	n := len(h.data)
	for {
		left := 2*i + 1
		if left >= n {
			break
		}
		best := left
		if right := left + 1; right < n && h.less(h.data[right], h.data[left]) {
			best = right
		}
		if !h.less(h.data[best], h.data[i]) {
			break
		}
		h.data[i], h.data[best] = h.data[best], h.data[i]
		i = best
	}
}
```

### wrapping the standard library `container/heap`

The idiomatic Go approach is to define an inner unexported type that implements
`heap.Interface` (Len/Less/Swap plus pointer-receiver Push/Pop), then expose a
clean typed wrapper so callers never touch `interface{}` or the `heap` package
directly.

```go
package taskpq

import "container/heap"

// Task is the concrete element stored in the priority queue.
type Task struct {
	Name     string
	Priority int // smaller value = higher priority
}

// taskHeap is the inner type implementing heap.Interface.
type taskHeap []Task

func (h taskHeap) Len() int            { return len(h) }
func (h taskHeap) Less(i, j int) bool  { return h[i].Priority < h[j].Priority }
func (h taskHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *taskHeap) Push(x any)         { *h = append(*h, x.(Task)) }
func (h *taskHeap) Pop() any {
	old := *h
	n := len(old)
	v := old[n-1]
	*h = old[:n-1]
	return v
}

// PriorityQueue is the public typed wrapper around container/heap.
type PriorityQueue struct {
	h taskHeap
}

// NewPriorityQueue returns an empty priority queue.
func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{}
	heap.Init(&pq.h)
	return pq
}

// Len reports the number of tasks.
func (pq *PriorityQueue) Len() int { return pq.h.Len() }

// Push adds a task in O(log n).
func (pq *PriorityQueue) Push(t Task) { heap.Push(&pq.h, t) }

// Pop removes and returns the highest-priority task in O(log n).
func (pq *PriorityQueue) Pop() (Task, bool) {
	if pq.h.Len() == 0 {
		return Task{}, false
	}
	return heap.Pop(&pq.h).(Task), true
}

// Peek returns the highest-priority task without removing it, O(1).
func (pq *PriorityQueue) Peek() (Task, bool) {
	if pq.h.Len() == 0 {
		return Task{}, false
	}
	return pq.h[0], true
}
```

## usage / test

```go
package heapq

import (
	"sort"
	"testing"
)

func TestMinHeapPopsAscending(t *testing.T) {
	in := []int{5, 1, 9, 3, 7, 2, 8, 4, 6, 0}
	h := FromSlice(append([]int(nil), in...), func(a, b int) bool { return a < b })

	var out []int
	for h.Len() > 0 {
		v, _ := h.Pop()
		out = append(out, v)
	}

	// A min-heap must yield non-decreasing order.
	want := append([]int(nil), in...)
	sort.Ints(want)
	for i := range want {
		if out[i] != want[i] {
			t.Fatalf("at %d: got %d want %d (full %v)", i, out[i], want[i], out)
		}
	}
	for i := 1; i < len(out); i++ {
		if out[i-1] > out[i] {
			t.Fatalf("not non-decreasing at %d: %v", i, out)
		}
	}
}

func TestMaxHeapPeek(t *testing.T) {
	h := NewMax[int]()
	for _, v := range []int{3, 1, 4, 1, 5, 9, 2, 6} {
		h.Push(v)
	}
	if top, _ := h.Peek(); top != 9 {
		t.Fatalf("max-heap peek: got %d want 9", top)
	}
	if v, _ := h.Pop(); v != 9 {
		t.Fatalf("max-heap pop: got %d want 9", v)
	}
}
```
