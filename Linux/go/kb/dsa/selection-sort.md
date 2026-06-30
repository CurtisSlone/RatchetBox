# Selection Sort

Selection sort divides the slice into a sorted prefix and an unsorted suffix; on each pass it scans the suffix for the minimum element and swaps it into the next sorted position. Use it when the number of swaps must be minimized (it performs at most n-1 swaps), but otherwise it is slow. It is in-place and NOT stable in its simple swap form, with O(n^2) time in every case (best, average, and worst) and O(1) extra space. Keywords: selection sort sorting comparison sort in-place minimum element find min swap quadratic O(n^2) unstable fewest swaps slice generic

## implementation

```go
package sorting

import "cmp"

// SelectionSort sorts s in place into non-decreasing order using selection sort.
// It always runs in O(n^2) time and performs at most n-1 swaps.
func SelectionSort[T cmp.Ordered](s []T) {
	for i := 0; i < len(s); i++ {
		minIdx := i
		for j := i + 1; j < len(s); j++ {
			if s[j] < s[minIdx] {
				minIdx = j
			}
		}
		if minIdx != i {
			s[i], s[minIdx] = s[minIdx], s[i]
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

func TestSelectionSort(t *testing.T) {
	cases := [][]int{
		{}, {1}, {3, 2, 1}, {64, 25, 12, 22, 11}, {0, -1, -1, 2, 2},
	}
	for _, in := range cases {
		got := slices.Clone(in)
		SelectionSort(got)
		want := slices.Clone(in)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("SelectionSort(%v) = %v, want %v", in, got, want)
		}
	}
}
```
