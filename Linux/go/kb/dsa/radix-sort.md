# Radix Sort

Radix sort is a non-comparison integer sort that sorts numbers digit by digit, from least significant digit (LSD) to most significant, using a stable counting sort as the per-digit subroutine. Because each pass is stable, processing all digits leaves the whole slice sorted. Use it for large arrays of fixed-width integers or strings where it can beat O(n log n) comparison sorts. With base b and d digits it runs in O(d * (n + b)) time and O(n + b) space. This LSD implementation handles non-negative integers in base 256 (one byte per pass). Keywords: radix sort LSD MSD sorting non-comparison integer sort digit by digit linear stable counting sort base bucket O(d*(n+b)) bytewise fixed width keys

## implementation

```go
package sorting

// RadixSort returns a new slice with the non-negative integers of s sorted in
// non-decreasing order using LSD radix sort in base 256. O(d*(n+256)) time.
func RadixSort(s []int) []int {
	out := make([]int, len(s))
	copy(out, s)
	if len(out) < 2 {
		return out
	}
	max := out[0]
	for _, v := range out {
		if v > max {
			max = v
		}
	}
	// process one byte (8-bit digit) per pass, least significant first
	for shift := uint(0); max>>shift > 0; shift += 8 {
		out = countingByByte(out, shift)
	}
	return out
}

// countingByByte stably sorts s by the byte at the given bit shift.
func countingByByte(s []int, shift uint) []int {
	const radix = 256
	count := make([]int, radix)
	for _, v := range s {
		count[(v>>shift)&0xFF]++
	}
	for i := 1; i < radix; i++ {
		count[i] += count[i-1]
	}
	out := make([]int, len(s))
	for i := len(s) - 1; i >= 0; i-- {
		d := (s[i] >> shift) & 0xFF
		count[d]--
		out[count[d]] = s[i]
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

func TestRadixSort(t *testing.T) {
	cases := [][]int{
		{}, {1}, {170, 45, 75, 90, 802, 24, 2, 66}, {255, 256, 257, 1, 0}, {7, 7, 7},
	}
	for _, in := range cases {
		got := RadixSort(in)
		want := slices.Clone(in)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("RadixSort(%v) = %v, want %v", in, got, want)
		}
	}
}
```
