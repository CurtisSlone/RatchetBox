# Counting Sort

Counting sort is a non-comparison integer sort: it counts how many times each key occurs, turns the counts into a prefix-sum of starting positions, then places each element directly into its sorted slot. Use it when keys are integers in a small known range k; it beats comparison sorts (which are bounded by O(n log n)) by running in O(n + k) time. It needs O(n + k) extra space and is stable when built with prefix sums and a right-to-left placement pass. It only works for integer (or integer-mappable) keys within a bounded range. Keywords: counting sort sorting non-comparison integer sort linear time O(n+k) stable prefix sum histogram bucket key range distribution sort radix building block

## implementation

```go
package sorting

// CountingSort returns a new slice with the non-negative integers of s sorted
// in non-decreasing order. Runs in O(n + max) time and space; stable.
func CountingSort(s []int) []int {
	if len(s) == 0 {
		return []int{}
	}
	max := s[0]
	for _, v := range s {
		if v > max {
			max = v
		}
	}
	// count[v] = number of occurrences of value v
	count := make([]int, max+1)
	for _, v := range s {
		count[v]++
	}
	// prefix sum: count[v] becomes the first output index for value v
	for v := 1; v <= max; v++ {
		count[v] += count[v-1]
	}
	out := make([]int, len(s))
	// place right-to-left so equal keys keep input order (stable)
	for i := len(s) - 1; i >= 0; i-- {
		v := s[i]
		count[v]--
		out[count[v]] = v
	}
	return out
}
```

## usage / test

```go
package sorting

import (
	"slices"
	"testing"
)

func TestCountingSort(t *testing.T) {
	cases := [][]int{
		{}, {0}, {4, 2, 2, 8, 3, 3, 1}, {5, 5, 5}, {0, 0, 9, 1, 9},
	}
	for _, in := range cases {
		got := CountingSort(in)
		want := slices.Clone(in)
		slices.Sort(want)
		if len(want) == 0 {
			want = []int{}
		}
		if !slices.Equal(got, want) {
			t.Errorf("CountingSort(%v) = %v, want %v", in, got, want)
		}
	}
}
```
