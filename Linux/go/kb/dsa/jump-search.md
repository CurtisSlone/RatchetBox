# Jump Search

Jump search works on a SORTED slice by jumping ahead in fixed blocks of size sqrt(n) until it passes the target, then doing a linear scan backward through the last block. Use it when binary search is awkward (for example on data where jumping forward is cheap but jumping backward is costly, like certain tape or linked structures), since it only ever steps backward within one block. Optimal block size is sqrt(n), giving O(sqrt(n)) time and O(1) space; the input must be sorted. Keywords: jump search block search sorted array square root sqrt(n) O(sqrt n) skip search block linear scan sorted divide blocks find index

## implementation

```go
package search

import (
	"cmp"
	"math"
)

// JumpSearch returns the index of target in the sorted slice s, or -1 if it is
// not present. It jumps in blocks of size sqrt(n) then scans the final block.
// O(sqrt(n)) time, O(1) space.
func JumpSearch[T cmp.Ordered](s []T, target T) int {
	n := len(s)
	if n == 0 {
		return -1
	}
	step := int(math.Sqrt(float64(n)))
	if step < 1 {
		step = 1
	}
	// find the block whose right edge is >= target
	prev := 0
	for next := step; ; next += step {
		if next >= n {
			next = n
		}
		if s[next-1] >= target {
			// linear scan within [prev, next)
			for i := prev; i < next; i++ {
				if s[i] == target {
					return i
				}
			}
			return -1
		}
		prev = next
		if next == n {
			return -1
		}
	}
}
```

## usage / test

```go
package search

import "testing"

func TestJumpSearch(t *testing.T) {
	s := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55}
	for i, v := range s {
		if got := JumpSearch(s, v); s[got] != v {
			t.Errorf("JumpSearch(%d) = %d (value %d), want a match at value %d", v, got, s[got], v)
		}
		_ = i
	}
	if got := JumpSearch(s, 7); got != -1 {
		t.Errorf("JumpSearch(7) = %d, want -1", got)
	}
	if got := JumpSearch([]int{}, 1); got != -1 {
		t.Errorf("JumpSearch(empty) = %d, want -1", got)
	}
}
```
