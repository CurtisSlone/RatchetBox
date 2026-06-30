# Interpolation Search

Interpolation search improves on binary search for SORTED, UNIFORMLY DISTRIBUTED numeric data by estimating the target's position with linear interpolation between the low and high values, instead of always probing the middle. On uniform data it runs in O(log log n) time; on skewed or clustered data it degrades to O(n). It needs O(1) space and requires sorted numeric keys. Guard against division by zero when low and high values are equal, and clamp the probe index to the current range. Keywords: interpolation search sorted array uniform distribution numeric keys O(log log n) probe position estimate linear interpolation guess index phone book search find index

## implementation

```go
package search

// InterpolationSearch returns the index of target in the sorted, ideally
// uniformly distributed int slice s, or -1 if absent. O(log log n) on uniform
// data, O(n) worst case.
func InterpolationSearch(s []int, target int) int {
	lo, hi := 0, len(s)-1
	for lo <= hi && target >= s[lo] && target <= s[hi] {
		if s[lo] == s[hi] {
			// flat range: only equal values remain
			if s[lo] == target {
				return lo
			}
			return -1
		}
		// estimate position by linear interpolation, then clamp into range
		pos := lo + int(int64(target-s[lo])*int64(hi-lo)/int64(s[hi]-s[lo]))
		if pos < lo {
			pos = lo
		} else if pos > hi {
			pos = hi
		}
		switch {
		case s[pos] == target:
			return pos
		case s[pos] < target:
			lo = pos + 1
		default:
			hi = pos - 1
		}
	}
	return -1
}
```

## usage / test

```go
package search

import "testing"

func TestInterpolationSearch(t *testing.T) {
	s := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	for i, v := range s {
		if got := InterpolationSearch(s, v); got != i {
			t.Errorf("InterpolationSearch(%d) = %d, want %d", v, got, i)
		}
	}
	if got := InterpolationSearch(s, 55); got != -1 {
		t.Errorf("InterpolationSearch(55) = %d, want -1", got)
	}
	if got := InterpolationSearch([]int{}, 1); got != -1 {
		t.Errorf("InterpolationSearch(empty) = %d, want -1", got)
	}
}
```
