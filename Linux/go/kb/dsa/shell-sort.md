# Shell Sort

Shell sort is a generalization of insertion sort that first sorts elements far apart (a "gap"), then progressively reduces the gap to 1, so that out-of-place elements move long distances early and the final gap-1 pass has little work left. Use it as a simple in-place sort that beats plain insertion sort on medium inputs without recursion. Time depends on the gap sequence: O(n^2) worst case for the simple halving sequence, around O(n^(3/2)) to O(n log^2 n) for better sequences; space is O(1). It is in-place and NOT stable. Keywords: shell sort shellsort sorting diminishing increment gap sequence insertion sort in-place unstable O(n^2) O(n log^2 n) gapped insertion slice generic

## implementation

```go
package sorting

import "cmp"

// ShellSort sorts s in place into non-decreasing order using shell sort with
// the classic halving gap sequence. In-place, not stable.
func ShellSort[T cmp.Ordered](s []T) {
	n := len(s)
	for gap := n / 2; gap > 0; gap /= 2 {
		// gapped insertion sort: each element jumps in steps of `gap`
		for i := gap; i < n; i++ {
			key := s[i]
			j := i
			for j >= gap && s[j-gap] > key {
				s[j] = s[j-gap]
				j -= gap
			}
			s[j] = key
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

func TestShellSort(t *testing.T) {
	cases := [][]int{
		{}, {1}, {2, 1}, {23, 12, 1, 8, 34, 54, 2, 3}, {5, 5, 4, 1, 1, 3},
	}
	for _, in := range cases {
		got := slices.Clone(in)
		ShellSort(got)
		want := slices.Clone(in)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("ShellSort(%v) = %v, want %v", in, got, want)
		}
	}
}
```
