# Bubble Sort

Bubble sort repeatedly steps through a slice, compares each adjacent pair, and swaps them if they are out of order, so larger elements "bubble" toward the end on each pass. Use it only for teaching, tiny inputs, or nearly-sorted data where its early-exit makes it run in O(n); otherwise it is one of the slowest sorts. It is a stable, in-place comparison sort with O(n^2) average and worst-case time and O(1) extra space. The optimized version stops as soon as a full pass makes no swaps. Keywords: bubble sort sorting comparison sort adjacent swap stable in-place quadratic O(n^2) sinking sort exchange sort ascending order swap pass early exit slice generic

## implementation

```go
package sorting

import "cmp"

// BubbleSort sorts s in place into non-decreasing order using bubble sort.
// It is stable and runs in O(n) on already-sorted input thanks to the
// early-exit swapped flag, and O(n^2) in the average and worst case.
func BubbleSort[T cmp.Ordered](s []T) {
	for end := len(s); end > 1; end-- {
		swapped := false
		for i := 0; i < end-1; i++ {
			if s[i] > s[i+1] {
				s[i], s[i+1] = s[i+1], s[i]
				swapped = true
			}
		}
		if !swapped {
			return // no swaps in a full pass means the slice is sorted
		}
	}
}
```

## usage / test

```go
package sorting

import (
	"slices"
	"testing"
)

func TestBubbleSort(t *testing.T) {
	cases := [][]int{
		{},
		{1},
		{5, 4, 3, 2, 1},
		{3, 1, 2, 3, 1},
		{-2, 0, -5, 7, 7, 3},
	}
	for _, in := range cases {
		got := slices.Clone(in)
		BubbleSort(got)
		want := slices.Clone(in)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("BubbleSort(%v) = %v, want %v", in, got, want)
		}
		// assert the key property: non-decreasing order
		if !slices.IsSorted(got) {
			t.Errorf("result %v is not sorted", got)
		}
	}
}
```
