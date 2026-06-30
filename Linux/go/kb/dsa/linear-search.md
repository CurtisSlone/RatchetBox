# Linear Search

Linear search scans a slice element by element from the start and returns the index of the first match (or -1 if absent). Use it when the data is unsorted, small, or only searched a few times, since it needs no preprocessing. It runs in O(n) time and O(1) space and works on any equality-comparable element type. For sorted data prefer binary search; for repeated lookups prefer a map. Keywords: linear search sequential search find index scan O(n) unsorted brute force contains lookup first match slice generic comparable

## implementation

```go
package search

// Linear returns the index of the first element in s equal to target, or -1 if
// target is not present. O(n) time, O(1) space; works on unsorted data.
func Linear[T comparable](s []T, target T) int {
	for i, v := range s {
		if v == target {
			return i
		}
	}
	return -1
}
```

## usage / test

```go
package search

import "testing"

func TestLinear(t *testing.T) {
	s := []int{5, 3, 8, 1, 9, 2}
	cases := []struct {
		target int
		want   int
	}{
		{5, 0}, {2, 5}, {8, 2}, {7, -1},
	}
	for _, c := range cases {
		if got := Linear(s, c.target); got != c.want {
			t.Errorf("Linear(%v, %d) = %d, want %d", s, c.target, got, c.want)
		}
	}
	// works on strings too
	if got := Linear([]string{"a", "b", "c"}, "c"); got != 2 {
		t.Errorf("Linear strings = %d, want 2", got)
	}
}
```
