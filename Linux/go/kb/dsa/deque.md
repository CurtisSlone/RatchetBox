# Deque (Double-Ended Queue)

A deque (double-ended queue, pronounced "deck") is a sequence that supports adding and removing elements at both the front and the back, so it generalizes both the stack and the queue. Use it for sliding-window algorithms (e.g. the monotonic-deque sliding-window maximum), work-stealing schedulers, undo/redo buffers, and any case where you push and pop from either end. Backed by a growable circular slice with head and tail indices, PushFront, PushBack, PopFront, PopBack, and peeking at either end are all amortized O(1), with O(n) space; a doubly linked list gives the same O(1) bounds without reallocation. Keywords: deque double-ended queue dequeue head-tail push-front push-back pop-front pop-back append-left append-right peek-front peek-back front back both-ends ring circular buffer sliding-window monotonic stack queue generalization amortized O(1) container

## implementation

```go
package deque

// Deque is a double-ended queue backed by a growable circular slice.
// The zero value is ready to use.
type Deque[T any] struct {
	data  []T
	head  int // index of the front element
	count int // number of stored elements
}

// New returns an empty deque; capacity is an optional hint.
func New[T any](capacity int) *Deque[T] {
	return &Deque[T]{data: make([]T, capacity)}
}

// Len reports the number of elements.
func (d *Deque[T]) Len() int { return d.count }

// IsEmpty reports whether the deque has no elements.
func (d *Deque[T]) IsEmpty() bool { return d.count == 0 }

// PushBack appends v at the back in amortized O(1).
func (d *Deque[T]) PushBack(v T) {
	if d.count == len(d.data) {
		d.grow()
	}
	tail := (d.head + d.count) % len(d.data)
	d.data[tail] = v
	d.count++
}

// PushFront prepends v at the front in amortized O(1).
func (d *Deque[T]) PushFront(v T) {
	if d.count == len(d.data) {
		d.grow()
	}
	d.head = (d.head - 1 + len(d.data)) % len(d.data)
	d.data[d.head] = v
	d.count++
}

// PopFront removes and returns the front value in O(1).
func (d *Deque[T]) PopFront() (T, bool) {
	var zero T
	if d.count == 0 {
		return zero, false
	}
	v := d.data[d.head]
	d.data[d.head] = zero
	d.head = (d.head + 1) % len(d.data)
	d.count--
	return v, true
}

// PopBack removes and returns the back value in O(1).
func (d *Deque[T]) PopBack() (T, bool) {
	var zero T
	if d.count == 0 {
		return zero, false
	}
	tail := (d.head + d.count - 1) % len(d.data)
	v := d.data[tail]
	d.data[tail] = zero
	d.count--
	return v, true
}

// Front returns the front value without removing it.
func (d *Deque[T]) Front() (T, bool) {
	var zero T
	if d.count == 0 {
		return zero, false
	}
	return d.data[d.head], true
}

// Back returns the back value without removing it.
func (d *Deque[T]) Back() (T, bool) {
	var zero T
	if d.count == 0 {
		return zero, false
	}
	return d.data[(d.head+d.count-1)%len(d.data)], true
}

func (d *Deque[T]) grow() {
	newCap := len(d.data) * 2
	if newCap == 0 {
		newCap = 4
	}
	buf := make([]T, newCap)
	for i := 0; i < d.count; i++ {
		buf[i] = d.data[(d.head+i)%len(d.data)]
	}
	d.data = buf
	d.head = 0
}
```

## usage / test

```go
package deque

import "testing"

func TestDequeBothEnds(t *testing.T) {
	d := New[int](1) // tiny capacity forces growth
	d.PushBack(2)
	d.PushBack(3)
	d.PushFront(1)
	d.PushFront(0) // 0 1 2 3

	if v, _ := d.Front(); v != 0 {
		t.Fatalf("front: got %d want 0", v)
	}
	if v, _ := d.Back(); v != 3 {
		t.Fatalf("back: got %d want 3", v)
	}

	// Drain from the front: 0 1 2 3.
	for want := 0; want <= 3; want++ {
		got, ok := d.PopFront()
		if !ok || got != want {
			t.Fatalf("pop front: got %d want %d", got, want)
		}
	}
	if !d.IsEmpty() {
		t.Fatal("deque should be empty")
	}
}

// SlidingWindowMax uses a monotonic deque of indices to compute the
// maximum of every window of width k in O(n).
func SlidingWindowMax(nums []int, k int) []int {
	if k <= 0 || len(nums) == 0 {
		return nil
	}
	dq := New[int](k) // holds indices, values decreasing front->back
	var out []int
	for i, v := range nums {
		if front, ok := dq.Front(); ok && front <= i-k {
			dq.PopFront() // drop indices that fell out of the window
		}
		for {
			back, ok := dq.Back()
			if !ok || nums[back] >= v {
				break
			}
			dq.PopBack()
		}
		dq.PushBack(i)
		if i >= k-1 {
			front, _ := dq.Front()
			out = append(out, nums[front])
		}
	}
	return out
}

func TestSlidingWindowMax(t *testing.T) {
	got := SlidingWindowMax([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3)
	want := []int{3, 3, 5, 5, 6, 7}
	if len(got) != len(want) {
		t.Fatalf("len: got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("at %d: got %d want %d", i, got[i], want[i])
		}
	}
}
```
