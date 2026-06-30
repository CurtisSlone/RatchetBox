# Ring Buffer (Circular Buffer, Fixed Capacity)

A ring buffer (circular buffer) is a fixed-capacity FIFO backed by a single array whose head and tail indices wrap around modulo the length, so no element is ever shifted and no memory is reallocated. Use it for bounded streaming buffers, audio/IO pipelines, the most-recent-N log or telemetry window, and producer/consumer queues where a bounded memory footprint matters; when full it either rejects writes or overwrites the oldest element (an overwriting ring is a sliding window of the last N items). Push, pop, and peek are O(1) with O(capacity) space; reads and writes only touch one slot, which makes it cache-friendly and lock-free-friendly. Keywords: ring buffer circular buffer cyclic buffer fixed capacity bounded FIFO head tail read write index wrap-around modulo overwrite oldest full empty isFull isEmpty push pop peek streaming audio IO producer consumer last-n window O(1) container

## implementation

```go
package ringbuffer

// RingBuffer is a fixed-capacity circular FIFO buffer.
type RingBuffer[T any] struct {
	data     []T
	head     int // index of the oldest element
	count    int // number of stored elements
	overwrite bool
}

// New returns an empty ring buffer of the given capacity. When
// overwrite is true, writing to a full buffer evicts the oldest
// element; otherwise Push reports false when the buffer is full.
func New[T any](capacity int, overwrite bool) *RingBuffer[T] {
	if capacity < 1 {
		capacity = 1
	}
	return &RingBuffer[T]{data: make([]T, capacity), overwrite: overwrite}
}

// Cap reports the fixed capacity.
func (r *RingBuffer[T]) Cap() int { return len(r.data) }

// Len reports the number of stored elements.
func (r *RingBuffer[T]) Len() int { return r.count }

// IsEmpty and IsFull report the boundary states.
func (r *RingBuffer[T]) IsEmpty() bool { return r.count == 0 }
func (r *RingBuffer[T]) IsFull() bool  { return r.count == len(r.data) }

// Push writes v at the tail in O(1). When full and overwrite is
// disabled it returns false; when overwrite is enabled it evicts the
// oldest element and returns true.
func (r *RingBuffer[T]) Push(v T) bool {
	if r.IsFull() {
		if !r.overwrite {
			return false
		}
		r.head = (r.head + 1) % len(r.data) // drop the oldest
		r.count--
	}
	tail := (r.head + r.count) % len(r.data)
	r.data[tail] = v
	r.count++
	return true
}

// Pop removes and returns the oldest element in O(1).
func (r *RingBuffer[T]) Pop() (T, bool) {
	var zero T
	if r.count == 0 {
		return zero, false
	}
	v := r.data[r.head]
	r.data[r.head] = zero
	r.head = (r.head + 1) % len(r.data)
	r.count--
	return v, true
}

// Peek returns the oldest element without removing it.
func (r *RingBuffer[T]) Peek() (T, bool) {
	var zero T
	if r.count == 0 {
		return zero, false
	}
	return r.data[r.head], true
}

// Slice returns the stored elements oldest-first.
func (r *RingBuffer[T]) Slice() []T {
	out := make([]T, r.count)
	for i := 0; i < r.count; i++ {
		out[i] = r.data[(r.head+i)%len(r.data)]
	}
	return out
}
```

## usage / test

```go
package ringbuffer

import (
	"reflect"
	"testing"
)

func TestRingBufferFIFOAndBounds(t *testing.T) {
	r := New[int](3, false)
	if !r.Push(1) || !r.Push(2) || !r.Push(3) {
		t.Fatal("pushes within capacity must succeed")
	}
	if r.Push(4) { // buffer is full, overwrite disabled
		t.Fatal("push past capacity must fail when overwrite is off")
	}
	if r.Cap() != 3 || r.Len() != 3 {
		t.Fatalf("cap/len: got %d/%d want 3/3", r.Cap(), r.Len())
	}
	// FIFO order is preserved.
	for want := 1; want <= 3; want++ {
		got, _ := r.Pop()
		if got != want {
			t.Fatalf("pop: got %d want %d", got, want)
		}
	}
}

func TestRingBufferOverwrite(t *testing.T) {
	r := New[int](3, true)
	for i := 1; i <= 5; i++ {
		r.Push(i) // 4 and 5 overwrite 1 and 2
	}
	// Only the last 3 values survive, oldest-first.
	if got := r.Slice(); !reflect.DeepEqual(got, []int{3, 4, 5}) {
		t.Fatalf("overwrite window: got %v want [3 4 5]", got)
	}
	if r.Len() != r.Cap() {
		t.Fatalf("length must stay bounded at capacity")
	}
}
```
