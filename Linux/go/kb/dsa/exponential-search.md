# Exponential Search

Exponential search (also called galloping or doubling search) finds a range to search in a SORTED slice by doubling an index (1, 2, 4, 8, ...) until it passes the target, then runs binary search within the last doubled interval. Use it when the array is unbounded or very large and the target is expected near the front, since it touches only O(log i) elements to bracket a target at index i. Total time is O(log i) where i is the target's position (O(log n) worst case), with O(1) space; the input must be sorted. Keywords: exponential search galloping search doubling search unbounded search sorted array O(log n) O(log i) bracket range then binary search find index infinite array

## implementation

```go
package search

import "cmp"

// ExponentialSearch returns the index of target in the sorted slice s, or -1 if
// it is not present. It doubles a bound until it passes target, then binary
// searches the bracketed range. O(log i) time where i is target's position.
func ExponentialSearch[T cmp.Ordered](s []T, target T) int {
	n := len(s)
	if n == 0 {
		return -1
	}
	if s[0] == target {
		return 0
	}
	// double the bound until s[bound] >= target or we run off the end
	bound := 1
	for bound < n && s[bound] < target {
		bound *= 2
	}
	lo := bound / 2
	hi := bound
	if hi > n-1 {
		hi = n - 1
	}
	// binary search within [lo, hi]
	for lo <= hi {
		mid := lo + (hi-lo)/2
		switch {
		case s[mid] == target:
			return mid
		case s[mid] < target:
			lo = mid + 1
		default:
			hi = mid - 1
		}
	}
	return -1
}
```

## usage / test

```go
package search

import "testing"

func TestExponentialSearch(t *testing.T) {
	s := []int{2, 3, 4, 10, 40, 50, 60, 70, 80, 90, 100}
	cases := map[int]bool{2: true, 10: true, 100: true, 5: false, 99: false}
	for target, present := range cases {
		got := ExponentialSearch(s, target)
		if present && (got < 0 || s[got] != target) {
			t.Errorf("ExponentialSearch(%d) = %d, want a match", target, got)
		}
		if !present && got != -1 {
			t.Errorf("ExponentialSearch(%d) = %d, want -1", target, got)
		}
	}
}
```
