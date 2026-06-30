# Comb Sort

Comb sort improves bubble sort by comparing elements a shrinking "gap" apart instead of always adjacent, which quickly eliminates small values stranded near the end (so-called "turtles"). Start with gap = n, divide it by the shrink factor 1.3 each pass, and finish with gap 1 (ordinary bubble passes) until no swaps occur. Use it as a simple in-place sort that is much faster than bubble sort in practice. Average time is around O(n^2 / 2^p) trending to O(n log n), worst case O(n^2), space O(1); it is NOT stable. Keywords: comb sort combsort sorting bubble sort improvement shrink factor 1.3 gap turtles in-place unstable O(n^2) O(n log n) swap slice generic

## implementation

```go
package sorting

import "cmp"

// CombSort sorts s in place into non-decreasing order using comb sort with the
// standard shrink factor of 1.3. In-place, not stable.
func CombSort[T cmp.Ordered](s []T) {
	n := len(s)
	gap := n
	swapped := true
	for gap > 1 || swapped {
		// shrink the gap by 1.3 (integer division), never below 1
		gap = gap * 10 / 13
		if gap < 1 {
			gap = 1
		}
		swapped = false
		for i := 0; i+gap < n; i++ {
			if s[i] > s[i+gap] {
				s[i], s[i+gap] = s[i+gap], s[i]
				swapped = true
			}
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

func TestCombSort(t *testing.T) {
	cases := [][]int{
		{}, {1}, {2, 1}, {8, 4, 1, 56, 3, -44, 23, -6, 28, 0}, {3, 3, 2, 2, 1, 1},
	}
	for _, in := range cases {
		got := slices.Clone(in)
		CombSort(got)
		want := slices.Clone(in)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("CombSort(%v) = %v, want %v", in, got, want)
		}
	}
}
```
