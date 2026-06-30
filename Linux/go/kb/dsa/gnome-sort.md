# Gnome Sort

Gnome sort (stupid sort) walks an index forward; whenever the current element is smaller than its left neighbor it swaps them and steps back one position, otherwise it steps forward. The effect is the same as insertion sort but expressed with a single moving index and no nested loop. Use it only for its simplicity on tiny or nearly-sorted inputs. It is stable and in-place, with O(n^2) average and worst-case time, O(n) on already-sorted input, and O(1) space. Keywords: gnome sort stupid sort sorting stable in-place quadratic O(n^2) single loop swap back insertion sort variant slice generic

## implementation

```go
package sorting

import "cmp"

// GnomeSort sorts s in place into non-decreasing order using gnome sort.
// Stable, in-place, O(n^2) worst case, O(n) on sorted input.
func GnomeSort[T cmp.Ordered](s []T) {
	i := 0
	for i < len(s) {
		if i == 0 || s[i-1] <= s[i] {
			i++ // in order: step forward
		} else {
			s[i-1], s[i] = s[i], s[i-1]
			i-- // out of order: swap and step back
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

func TestGnomeSort(t *testing.T) {
	cases := [][]int{
		{}, {1}, {2, 1}, {34, 2, 10, -9, 0, 7}, {4, 4, 3, 3, 2, 1, 1},
	}
	for _, in := range cases {
		got := slices.Clone(in)
		GnomeSort(got)
		want := slices.Clone(in)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("GnomeSort(%v) = %v, want %v", in, got, want)
		}
	}
}
```
