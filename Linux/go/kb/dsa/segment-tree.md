# Segment Tree (Range Query and Range Update with Lazy Propagation)

A segment tree is a binary tree built over an array where each node stores an aggregate (here, the sum) of a contiguous segment, with leaves holding single elements and the root covering the whole array. Use it for fast range queries (sum, min, max, gcd, etc.) combined with updates: a point or whole-range update plus a range query each run in O(log n), versus O(n) for a naive array. This implementation supports range-add updates using lazy propagation, where pending updates are deferred at internal nodes and pushed down only when needed, keeping range updates at O(log n) too. Building is O(n), queries and updates are O(log n), and the tree array uses O(n) space (allocated as 4n for safety). Keywords: segment tree range query range update range sum range min range max lazy propagation interval tree point update build query update aggregate logarithmic divide and conquer node leaf prefix.

## implementation

```go
package segmenttree

// SegmentTree supports range-sum queries and range-add updates in O(log n)
// using lazy propagation.
type SegmentTree struct {
	n    int
	tree []int // aggregate (sum) per node, 1-based node indexing
	lazy []int // pending range-add per node
}

// New builds a segment tree from the given values.
func New(values []int) *SegmentTree {
	n := len(values)
	st := &SegmentTree{
		n:    n,
		tree: make([]int, 4*n),
		lazy: make([]int, 4*n),
	}
	if n > 0 {
		st.build(values, 1, 0, n-1)
	}
	return st
}

func (s *SegmentTree) build(vals []int, node, lo, hi int) {
	if lo == hi {
		s.tree[node] = vals[lo]
		return
	}
	mid := (lo + hi) / 2
	s.build(vals, 2*node, lo, mid)
	s.build(vals, 2*node+1, mid+1, hi)
	s.tree[node] = s.tree[2*node] + s.tree[2*node+1]
}

// push applies and propagates a node's pending lazy value.
func (s *SegmentTree) push(node, lo, hi int) {
	if s.lazy[node] == 0 {
		return
	}
	s.tree[node] += (hi - lo + 1) * s.lazy[node]
	if lo != hi {
		s.lazy[2*node] += s.lazy[node]
		s.lazy[2*node+1] += s.lazy[node]
	}
	s.lazy[node] = 0
}

// Update adds delta to every element in [l, r] (0-based, inclusive).
func (s *SegmentTree) Update(l, r, delta int) {
	if s.n > 0 {
		s.update(1, 0, s.n-1, l, r, delta)
	}
}

func (s *SegmentTree) update(node, lo, hi, l, r, delta int) {
	s.push(node, lo, hi)
	if r < lo || hi < l {
		return // no overlap
	}
	if l <= lo && hi <= r {
		s.lazy[node] += delta
		s.push(node, lo, hi)
		return
	}
	mid := (lo + hi) / 2
	s.update(2*node, lo, mid, l, r, delta)
	s.update(2*node+1, mid+1, hi, l, r, delta)
	s.tree[node] = s.tree[2*node] + s.tree[2*node+1]
}

// Query returns the sum of elements in [l, r] (0-based, inclusive).
func (s *SegmentTree) Query(l, r int) int {
	if s.n == 0 {
		return 0
	}
	return s.query(1, 0, s.n-1, l, r)
}

func (s *SegmentTree) query(node, lo, hi, l, r int) int {
	s.push(node, lo, hi)
	if r < lo || hi < l {
		return 0 // no overlap
	}
	if l <= lo && hi <= r {
		return s.tree[node]
	}
	mid := (lo + hi) / 2
	return s.query(2*node, lo, mid, l, r) + s.query(2*node+1, mid+1, hi, l, r)
}
```

## usage / test

```go
package segmenttree

import "testing"

func bruteSum(a []int, l, r int) int {
	s := 0
	for i := l; i <= r; i++ {
		s += a[i]
	}
	return s
}

func TestSegmentTree(t *testing.T) {
	a := []int{1, 3, 5, 7, 9, 11, -2, 4}
	st := New(a)

	// Initial range sums match brute force.
	for l := 0; l < len(a); l++ {
		for r := l; r < len(a); r++ {
			if got, want := st.Query(l, r), bruteSum(a, l, r); got != want {
				t.Fatalf("Query(%d,%d)=%d want %d", l, r, got, want)
			}
		}
	}

	// Range update [2,5] += 4 must be reflected everywhere.
	st.Update(2, 5, 4)
	for i := 2; i <= 5; i++ {
		a[i] += 4
	}
	for l := 0; l < len(a); l++ {
		for r := l; r < len(a); r++ {
			if got, want := st.Query(l, r), bruteSum(a, l, r); got != want {
				t.Fatalf("after update Query(%d,%d)=%d want %d", l, r, got, want)
			}
		}
	}
}
```
