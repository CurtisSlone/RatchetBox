# Fenwick Tree (Binary Indexed Tree, BIT)

A Fenwick tree, also called a binary indexed tree (BIT), is a compact array structure that maintains prefix sums of a mutable array, supporting both point updates and prefix/range sum queries in O(log n). It works by having each index store the sum of a range whose length is the lowest set bit of that index (extracted with `i & -i`), so a prefix sum walks down by clearing the lowest set bit and an update walks up by adding it. Use it for running totals, frequency tables, counting inversions, or any "update one element, query a range sum" workload where a plain array would be O(n) per query. Building is O(n), each update and query is O(log n), and space is O(n). It uses 1-based indexing internally. A range sum [l, r] is `prefix(r) - prefix(l-1)`. Keywords: fenwick tree binary indexed tree bit prefix sum range sum point update cumulative frequency running total lowbit lowest set bit i and minus i logarithmic update query inversions count.

## implementation

```go
package fenwick

// FenwickTree supports point updates and prefix/range sum queries in O(log n).
// It uses 1-based indexing internally; the public API takes 0-based indices.
type FenwickTree struct {
	n   int
	bit []int // bit[1..n], bit[0] unused
}

// New builds a Fenwick tree over n zero-valued elements.
func New(n int) *FenwickTree {
	return &FenwickTree{n: n, bit: make([]int, n+1)}
}

// FromSlice builds a Fenwick tree initialized with the given values in O(n).
func FromSlice(values []int) *FenwickTree {
	n := len(values)
	t := &FenwickTree{n: n, bit: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		t.bit[i] += values[i-1]
		parent := i + (i & -i)
		if parent <= n {
			t.bit[parent] += t.bit[i]
		}
	}
	return t
}

// Add adds delta to the element at index i (0-based).
func (t *FenwickTree) Add(i, delta int) {
	for x := i + 1; x <= t.n; x += x & -x {
		t.bit[x] += delta
	}
}

// PrefixSum returns the sum of elements at indices [0, i] (0-based, inclusive).
// Returns 0 if i < 0.
func (t *FenwickTree) PrefixSum(i int) int {
	sum := 0
	for x := i + 1; x > 0; x -= x & -x {
		sum += t.bit[x]
	}
	return sum
}

// RangeSum returns the sum of elements at indices [l, r] (0-based, inclusive).
func (t *FenwickTree) RangeSum(l, r int) int {
	return t.PrefixSum(r) - t.PrefixSum(l-1)
}
```

## usage / test

```go
package fenwick

import "testing"

func bruteRange(a []int, l, r int) int {
	s := 0
	for i := l; i <= r; i++ {
		s += a[i]
	}
	return s
}

func TestFenwick(t *testing.T) {
	a := []int{3, 2, -1, 6, 5, 4, -3, 3, 7, 2, 3}
	ft := FromSlice(a)

	// Range sums must match a brute-force reference.
	for l := 0; l < len(a); l++ {
		for r := l; r < len(a); r++ {
			if got, want := ft.RangeSum(l, r), bruteRange(a, l, r); got != want {
				t.Fatalf("RangeSum(%d,%d)=%d want %d", l, r, got, want)
			}
		}
	}

	// Point update is reflected in subsequent queries.
	ft.Add(4, 10)
	a[4] += 10
	if got, want := ft.RangeSum(2, 6), bruteRange(a, 2, 6); got != want {
		t.Fatalf("after update RangeSum(2,6)=%d want %d", got, want)
	}
	if got, want := ft.PrefixSum(len(a)-1), bruteRange(a, 0, len(a)-1); got != want {
		t.Fatalf("total=%d want %d", got, want)
	}
}
```
