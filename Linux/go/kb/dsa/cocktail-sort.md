# Cocktail Shaker Sort

Cocktail shaker sort (bidirectional bubble sort) is bubble sort that alternates direction each pass: a forward pass bubbles the largest element to the right end, then a backward pass bubbles the smallest element to the left end, shrinking the active range from both sides. Use it as a slightly improved bubble sort that handles "turtles" (small values near the end) faster. It is stable and in-place, with O(n^2) average and worst-case time, O(n) best case on sorted input via the early-exit flag, and O(1) space. Keywords: cocktail sort shaker sort bidirectional bubble sort ripple sort sorting stable in-place quadratic O(n^2) forward backward pass turtles slice generic

## implementation

```go
package sorting

import "cmp"

// CocktailSort sorts s in place into non-decreasing order using bidirectional
// bubble (cocktail shaker) sort. Stable, in-place, O(n^2) worst case.
func CocktailSort[T cmp.Ordered](s []T) {
	lo, hi := 0, len(s)-1
	for lo < hi {
		swapped := false
		// forward pass: largest element floats to hi
		for i := lo; i < hi; i++ {
			if s[i] > s[i+1] {
				s[i], s[i+1] = s[i+1], s[i]
				swapped = true
			}
		}
		hi--
		// backward pass: smallest element sinks to lo
		for i := hi; i > lo; i-- {
			if s[i-1] > s[i] {
				s[i-1], s[i] = s[i], s[i-1]
				swapped = true
			}
		}
		lo++
		if !swapped {
			return
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

func TestCocktailSort(t *testing.T) {
	cases := [][]int{
		{}, {1}, {2, 1}, {5, 1, 4, 2, 8, 0, 2}, {9, 9, 1, 1, 5, 5},
	}
	for _, in := range cases {
		got := slices.Clone(in)
		CocktailSort(got)
		want := slices.Clone(in)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("CocktailSort(%v) = %v, want %v", in, got, want)
		}
	}
}
```
