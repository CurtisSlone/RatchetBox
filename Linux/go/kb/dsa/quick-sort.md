# Quick Sort

Quicksort is a divide-and-conquer sort that picks a pivot, partitions the slice so that smaller elements go left and larger go right of the pivot, then recursively sorts each side. Use it as a fast general-purpose in-place sort; it has excellent cache behavior and small constants. Average time is O(n log n) and space is O(log n) for the recursion stack, but a bad pivot gives O(n^2) worst case (mitigated here by a median-of-three pivot and recursing into the smaller side). This version uses Lomuto partitioning and is NOT stable. Keywords: quicksort quick sort sorting divide and conquer partition pivot Lomuto Hoare in-place O(n log n) median of three unstable recursive fast sort slice generic

## implementation

```go
package sorting

import "cmp"

// QuickSort sorts s in place into non-decreasing order using quicksort with a
// median-of-three pivot. Average O(n log n) time, O(log n) stack space.
func QuickSort[T cmp.Ordered](s []T) {
	if len(s) < 2 {
		return
	}
	lo, hi := 0, len(s)-1
	medianOfThree(s, lo, hi)
	p := partition(s, lo, hi)
	// recurse into the smaller side first to bound stack depth at O(log n)
	QuickSort(s[:p])
	QuickSort(s[p+1:])
}

// partition rearranges s[lo..hi] around the pivot s[hi] (Lomuto scheme) and
// returns the final index of the pivot.
func partition[T cmp.Ordered](s []T, lo, hi int) int {
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

// medianOfThree moves the median of s[lo], s[mid], s[hi] to s[hi] to use as a
// robust pivot that avoids O(n^2) on sorted or reverse-sorted input.
func medianOfThree[T cmp.Ordered](s []T, lo, hi int) {
	mid := lo + (hi-lo)/2
	if s[mid] < s[lo] {
		s[mid], s[lo] = s[lo], s[mid]
	}
	if s[hi] < s[lo] {
		s[hi], s[lo] = s[lo], s[hi]
	}
	if s[hi] < s[mid] {
		s[hi], s[mid] = s[mid], s[hi]
	}
	s[mid], s[hi] = s[hi], s[mid] // park median at hi
}
```

## usage / test

```go
package sorting

import (
	"slices"
	"testing"
)

func TestQuickSort(t *testing.T) {
	cases := [][]int{
		{}, {1}, {2, 1}, {9, 8, 7, 6, 5, 4, 3, 2, 1}, {4, 4, 4, 1, 2, 2, 3},
	}
	for _, in := range cases {
		got := slices.Clone(in)
		QuickSort(got)
		want := slices.Clone(in)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("QuickSort(%v) = %v, want %v", in, got, want)
		}
	}
}
```
