# Two-heap streaming median (Go algorithm pattern)

Maintain the running median of a stream with two heaps: a max-heap for the lower half and a min-heap
for the upper half, using `container/heap`. Authored pattern (not from the GoF set); targets the common
bug where the size invariant is wrong and the median is read from the wrong heap.

THE INVARIANT (get this exactly right): after each insert, `len(lower) == len(upper)` or
`len(lower) == len(upper)+1`. Rebalance whenever the difference is OUTSIDE that range - i.e. move one
element when `len(lower) > len(upper)+1` or when `len(upper) > len(lower)`. A `+1` tolerance on BOTH
sides (e.g. `len(lower) > len(upper)+1` paired with `len(upper) > len(lower)+1`) lets the halves
diverge and produces wrong medians.

- `lower` is a MAX-heap (its root is the largest of the small half); `upper` is a MIN-heap.
- Insert: if `lower` is empty or `x <= lower.top()`, push to `lower`, else push to `upper`; then rebalance.
- Median: if `len(lower) > len(upper)`, it is `lower.top()`; otherwise `(lower.top()+upper.top())/2`.

```go
package solution

import "container/heap"

// intHeap is a heap of ints; max==true makes it a max-heap, else a min-heap.
type intHeap struct {
	data []int
	max  bool
}

func (h intHeap) Len() int { return len(h.data) }
func (h intHeap) Less(i, j int) bool {
	if h.max {
		return h.data[i] > h.data[j]
	}
	return h.data[i] < h.data[j]
}
func (h intHeap) Swap(i, j int)      { h.data[i], h.data[j] = h.data[j], h.data[i] }
func (h *intHeap) Push(x any)        { h.data = append(h.data, x.(int)) }
func (h *intHeap) Pop() any {
	old := h.data
	n := len(old)
	v := old[n-1]
	h.data = old[:n-1]
	return v
}
func (h *intHeap) top() int { return h.data[0] }

// RunningMedian returns the median of nums[0..i] for each i.
func RunningMedian(nums []int) []float64 {
	lower := &intHeap{max: true}  // small half (max-heap)
	upper := &intHeap{max: false} // large half (min-heap)
	out := make([]float64, 0, len(nums))

	for _, x := range nums {
		if lower.Len() == 0 || x <= lower.top() {
			heap.Push(lower, x)
		} else {
			heap.Push(upper, x)
		}
		// Rebalance to keep len(lower) == len(upper) or len(lower) == len(upper)+1.
		if lower.Len() > upper.Len()+1 {
			heap.Push(upper, heap.Pop(lower))
		} else if upper.Len() > lower.Len() {
			heap.Push(lower, heap.Pop(upper))
		}

		if lower.Len() > upper.Len() {
			out = append(out, float64(lower.top()))
		} else {
			out = append(out, float64(lower.top()+upper.top())/2)
		}
	}
	return out
}
```

Verified: `RunningMedian([]int{5, 4, 3, 2, 1})` returns `[5 4.5 4 3.5 3]`.
