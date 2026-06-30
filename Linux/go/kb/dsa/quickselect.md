# Quickselect (Select K / Kth Smallest)

Quickselect finds the k-th smallest element (0-indexed) of an unsorted slice without fully sorting it. It uses quicksort's partition step but recurses into only the one side that contains the target rank, so it runs in O(n) average time (O(n^2) worst case with bad pivots, avoidable with a random pivot). Use it to find medians, percentiles, or "top-k" thresholds faster than a full O(n log n) sort. It works in place, rearranging the input, and uses O(1) extra space. Keywords: quickselect select k kth smallest kth largest order statistic median percentile partition Hoare Lomuto top-k selection O(n) average nth_element introselect rank find kth random pivot

## implementation

```go
package search

import (
	"cmp"
	"math/rand"
)

// QuickSelect returns the k-th smallest element (0-indexed) of s. It rearranges
// s in place and runs in O(n) average time. It panics if k is out of range.
func QuickSelect[T cmp.Ordered](s []T, k int) T {
	if k < 0 || k >= len(s) {
		panic("QuickSelect: k out of range")
	}
	lo, hi := 0, len(s)-1
	for {
		if lo == hi {
			return s[lo]
		}
		p := partitionSelect(s, lo, hi)
		switch {
		case k == p:
			return s[p]
		case k < p:
			hi = p - 1
		default:
			lo = p + 1
		}
	}
}

// partitionSelect partitions s[lo..hi] around a random pivot (Lomuto) and
// returns the pivot's final index.
func partitionSelect[T cmp.Ordered](s []T, lo, hi int) int {
	r := lo + rand.Intn(hi-lo+1)
	s[r], s[hi] = s[hi], s[r] // move random pivot to the end
	pivot := s[hi]
	i := lo
	for j := lo; j < hi; j++ {
		if s[j] < pivot {
			s[i], s[j] = s[j], s[i]
			i++
		}
	}
	s[i], s[hi] = s[hi], s[i]
	return i
}
```

## usage / test

```go
package search

import (
	"slices"
	"testing"
)

func TestQuickSelect(t *testing.T) {
	for _, in := range [][]int{
		{3, 1, 2}, {7, 7, 7}, {9, 1, 8, 2, 7, 3, 6, 4, 5},
	} {
		sorted := slices.Clone(in)
		slices.Sort(sorted)
		for k := 0; k < len(in); k++ {
			work := slices.Clone(in)
			if got := QuickSelect(work, k); got != sorted[k] {
				t.Errorf("QuickSelect(%v, %d) = %d, want %d", in, k, got, sorted[k])
			}
		}
	}
	// median of an odd-length slice
	work := []int{5, 2, 8, 1, 9}
	if got := QuickSelect(work, 2); got != 5 {
		t.Errorf("median = %d, want 5", got)
	}
}
```
