# Merge Sort

Merge sort is a divide-and-conquer sort that recursively splits the slice in half, sorts each half, then merges the two sorted halves into one sorted run. Use it when you need guaranteed O(n log n) worst-case performance and a stable sort, or when sorting linked lists or external data too large for memory. It runs in O(n log n) time in all cases and uses O(n) extra space for the merge buffer. The merge step is the heart of the algorithm: compare the fronts of the two sorted halves and copy the smaller out each time. Keywords: merge sort sorting divide and conquer stable O(n log n) merge two sorted halves recursive top-down external sort linked list guaranteed worst case slice generic

## implementation

```go
package sorting

import "cmp"

// MergeSort returns a new sorted slice in non-decreasing order. It is stable
// and runs in O(n log n) time in all cases using O(n) auxiliary space.
func MergeSort[T cmp.Ordered](s []T) []T {
	if len(s) <= 1 {
		out := make([]T, len(s))
		copy(out, s)
		return out
	}
	mid := len(s) / 2
	left := MergeSort(s[:mid])
	right := MergeSort(s[mid:])
	return merge(left, right)
}

// merge combines two sorted slices into one sorted slice, preserving order of
// equal elements (left wins ties) so the overall sort is stable.
func merge[T cmp.Ordered](left, right []T) []T {
	out := make([]T, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			out = append(out, left[i])
			i++
		} else {
			out = append(out, right[j])
			j++
		}
	}
	out = append(out, left[i:]...)
	out = append(out, right[j:]...)
	return out
}
```

## usage / test

```go
package sorting

import (
	"slices"
	"testing"
)

func TestMergeSort(t *testing.T) {
	cases := [][]int{
		{}, {1}, {2, 1}, {38, 27, 43, 3, 9, 82, 10}, {5, 5, 4, 4, 3, 3},
	}
	for _, in := range cases {
		got := MergeSort(in)
		want := slices.Clone(in)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("MergeSort(%v) = %v, want %v", in, got, want)
		}
	}
}
```
