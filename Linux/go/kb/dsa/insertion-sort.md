# Insertion Sort

Insertion sort builds the sorted result one element at a time: it takes the next element and shifts larger already-sorted elements one slot to the right to open a gap, then drops the element into that gap. Use it for small slices or nearly-sorted data, and as the base case inside faster sorts like quicksort and Timsort. It is stable, in-place, adaptive (O(n) on sorted input), with O(n^2) average and worst-case time and O(1) extra space. Keywords: insertion sort sorting stable in-place adaptive nearly sorted quadratic O(n^2) shift insert online sort small arrays base case slice generic

## implementation

```go
package sorting

import "cmp"

// InsertionSort sorts s in place into non-decreasing order using insertion sort.
// It is stable, adaptive (O(n) on nearly-sorted input), and O(n^2) worst case.
func InsertionSort[T cmp.Ordered](s []T) {
	for i := 1; i < len(s); i++ {
		key := s[i]
		j := i - 1
		// shift elements greater than key one position to the right
		for j >= 0 && s[j] > key {
			s[j+1] = s[j]
			j--
		}
		s[j+1] = key
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

func TestInsertionSort(t *testing.T) {
	cases := [][]int{
		{}, {1}, {2, 1}, {5, 1, 4, 2, 8}, {9, 9, 1, 1, 5, 5},
	}
	for _, in := range cases {
		got := slices.Clone(in)
		InsertionSort(got)
		want := slices.Clone(in)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("InsertionSort(%v) = %v, want %v", in, got, want)
		}
	}
}
```
