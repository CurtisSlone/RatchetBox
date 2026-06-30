# Binary Search

Binary search finds a target in a SORTED slice by repeatedly halving the search interval: compare the middle element to the target and discard the half that cannot contain it. It runs in O(log n) time and O(1) space and requires the input to be sorted. The two essential boundary variants are lower bound (index of the first element >= target, the leftmost insertion point) and upper bound (index of the first element > target); their difference gives the count of equal keys, and lower bound also tells you where to insert to keep the slice sorted. Always compute the midpoint as lo+(hi-lo)/2 to avoid integer overflow. Keywords: binary search sorted array logarithmic O(log n) divide and conquer bisect lower bound upper bound leftmost rightmost insertion point first occurrence last occurrence find index half interval search sort.Search slices.BinarySearch overflow safe midpoint

## implementation

```go
package search

import "cmp"

// BinarySearch returns the index of target in the sorted slice s, or -1 if it
// is not present. O(log n) time. If duplicates exist, any matching index may
// be returned; use LowerBound for the first match.
func BinarySearch[T cmp.Ordered](s []T, target T) int {
	lo, hi := 0, len(s)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2 // overflow-safe midpoint
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

// LowerBound returns the index of the first element in s that is >= target.
// This is the leftmost position where target could be inserted to keep s
// sorted; it returns len(s) if every element is < target.
func LowerBound[T cmp.Ordered](s []T, target T) int {
	lo, hi := 0, len(s) // half-open [lo, hi)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if s[mid] < target {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

// UpperBound returns the index of the first element in s that is > target.
// UpperBound(s, x) - LowerBound(s, x) is the number of elements equal to x.
func UpperBound[T cmp.Ordered](s []T, target T) int {
	lo, hi := 0, len(s)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if s[mid] <= target {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}
```

## usage / test

```go
package search

import "testing"

func TestBinarySearch(t *testing.T) {
	s := []int{1, 2, 2, 2, 5, 8, 13}
	if got := BinarySearch(s, 8); got != 5 {
		t.Errorf("BinarySearch(8) = %d, want 5", got)
	}
	if got := BinarySearch(s, 7); got != -1 {
		t.Errorf("BinarySearch(7) = %d, want -1", got)
	}
	// lower/upper bound around the run of 2s
	if lb := LowerBound(s, 2); lb != 1 {
		t.Errorf("LowerBound(2) = %d, want 1", lb)
	}
	if ub := UpperBound(s, 2); ub != 4 {
		t.Errorf("UpperBound(2) = %d, want 4", ub)
	}
	if cnt := UpperBound(s, 2) - LowerBound(s, 2); cnt != 3 {
		t.Errorf("count of 2 = %d, want 3", cnt)
	}
	// insertion point for a missing value
	if lb := LowerBound(s, 6); lb != 5 {
		t.Errorf("LowerBound(6) = %d, want 5", lb)
	}
}
```
