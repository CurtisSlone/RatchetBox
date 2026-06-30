# Ternary Search

Ternary search comes in two flavors. On a SORTED array it splits the range into three parts with two midpoints and discards two-thirds each step to find a target; this is O(log3 n) comparisons but does more work per step than binary search, so binary search is usually preferred. Its more important use is finding the extremum (minimum or maximum) of a UNIMODAL function on a real interval: pick two interior points, drop the third that cannot contain the peak, and shrink the interval until it is smaller than a tolerance. The function version runs in O(log((b-a)/eps)) evaluations and O(1) space. Keywords: ternary search unimodal function maximum minimum peak finding optimization golden section sorted array O(log3 n) two midpoints divide into thirds interval extremum convex concave epsilon tolerance

## implementation

```go
package search

import "cmp"

// TernarySearch returns the index of target in the sorted slice s, or -1 if it
// is not present, by discarding two of three partitions each step. O(log n).
func TernarySearch[T cmp.Ordered](s []T, target T) int {
	lo, hi := 0, len(s)-1
	for lo <= hi {
		third := (hi - lo) / 3
		m1 := lo + third
		m2 := hi - third
		switch {
		case s[m1] == target:
			return m1
		case s[m2] == target:
			return m2
		case target < s[m1]:
			hi = m1 - 1
		case target > s[m2]:
			lo = m2 + 1
		default:
			lo, hi = m1+1, m2-1
		}
	}
	return -1
}

// TernaryMin returns the x in [a, b] minimizing the unimodal function f, to
// within tolerance eps. For a maximum, negate f or flip the comparison.
func TernaryMin(a, b, eps float64, f func(float64) float64) float64 {
	for b-a > eps {
		m1 := a + (b-a)/3
		m2 := b - (b-a)/3
		if f(m1) < f(m2) {
			b = m2
		} else {
			a = m1
		}
	}
	return (a + b) / 2
}
```

## usage / test

```go
package search

import (
	"math"
	"testing"
)

func TestTernarySearch(t *testing.T) {
	s := []int{1, 3, 5, 7, 9, 11, 13}
	if got := TernarySearch(s, 9); got != 4 {
		t.Errorf("TernarySearch(9) = %d, want 4", got)
	}
	if got := TernarySearch(s, 8); got != -1 {
		t.Errorf("TernarySearch(8) = %d, want -1", got)
	}

	// minimum of (x-2)^2 is at x = 2
	x := TernaryMin(-10, 10, 1e-7, func(x float64) float64 {
		return (x - 2) * (x - 2)
	})
	if math.Abs(x-2) > 1e-3 {
		t.Errorf("TernaryMin found x = %v, want ~2", x)
	}
}
```
