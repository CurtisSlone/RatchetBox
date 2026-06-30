# Heap Sort

Heap sort builds a binary max-heap over the slice in place, then repeatedly swaps the maximum (root) to the end of the unsorted region and sifts the new root down to restore the heap, growing a sorted suffix. Use it when you need guaranteed O(n log n) time with O(1) extra space and do not need stability. It always runs in O(n log n) (build heap is O(n), each of n extractions is O(log n)) and sorts entirely in place; it is NOT stable. Keywords: heapsort heap sort sorting binary heap max-heap sift down heapify in-place O(n log n) priority queue extract max guaranteed worst case unstable slice generic

## implementation

```go
package sorting

import "cmp"

// HeapSort sorts s in place into non-decreasing order using an in-place binary
// max-heap. O(n log n) time in all cases, O(1) extra space, not stable.
func HeapSort[T cmp.Ordered](s []T) {
	n := len(s)
	// build a max-heap bottom-up: O(n)
	for i := n/2 - 1; i >= 0; i-- {
		siftDown(s, i, n)
	}
	// repeatedly move the max to the end and shrink the heap
	for end := n - 1; end > 0; end-- {
		s[0], s[end] = s[end], s[0]
		siftDown(s, 0, end)
	}
}

// siftDown restores the max-heap property for the subtree rooted at i within
// s[0:size] by sinking s[i] down to its correct level.
func siftDown[T cmp.Ordered](s []T, i, size int) {
	for {
		largest := i
		left, right := 2*i+1, 2*i+2
		if left < size && s[left] > s[largest] {
			largest = left
		}
		if right < size && s[right] > s[largest] {
			largest = right
		}
		if largest == i {
			return
		}
		s[i], s[largest] = s[largest], s[i]
		i = largest
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

func TestHeapSort(t *testing.T) {
	cases := [][]int{
		{}, {1}, {3, 1, 2}, {12, 11, 13, 5, 6, 7}, {2, 2, 1, 1, 3, 3},
	}
	for _, in := range cases {
		got := slices.Clone(in)
		HeapSort(got)
		want := slices.Clone(in)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("HeapSort(%v) = %v, want %v", in, got, want)
		}
	}
}
```
