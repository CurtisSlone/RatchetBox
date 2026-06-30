# Queue (FIFO)

A queue is a first-in-first-out (FIFO) container: you enqueue values at the back and dequeue them from the front, so elements leave in the same order they arrived. Use it for breadth-first search, task and job scheduling, producer/consumer buffering, rate limiting, and any pipeline where work is processed in arrival order. The two common implementations are a linked list (one node per element, O(1) enqueue/dequeue, no reallocation) and a ring/slice with head and tail indices that grows when full; both give amortized O(1) enqueue, dequeue, and peek, with O(n) space. Keywords: queue FIFO first-in-first-out enqueue dequeue offer poll peek front back rear head tail isEmpty size BFS breadth-first scheduling buffer producer-consumer linked-list-backed slice-backed ring amortized container shift push-back pop-front

## implementation

```go
package queue

// Queue is a FIFO queue backed by a growable circular slice buffer.
// The zero value is ready to use.
type Queue[T any] struct {
	data  []T
	head  int // index of the front element
	count int // number of stored elements
}

// New returns an empty queue; capacity is an optional hint.
func New[T any](capacity int) *Queue[T] {
	return &Queue[T]{data: make([]T, capacity)}
}

// Len reports the number of elements.
func (q *Queue[T]) Len() int { return q.count }

// IsEmpty reports whether the queue has no elements.
func (q *Queue[T]) IsEmpty() bool { return q.count == 0 }

// Enqueue adds v at the back in amortized O(1).
func (q *Queue[T]) Enqueue(v T) {
	if q.count == len(q.data) {
		q.grow()
	}
	tail := (q.head + q.count) % len(q.data)
	q.data[tail] = v
	q.count++
}

// Dequeue removes and returns the front value in O(1).
func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if q.count == 0 {
		return zero, false
	}
	v := q.data[q.head]
	q.data[q.head] = zero // release reference
	q.head = (q.head + 1) % len(q.data)
	q.count--
	return v, true
}

// Peek returns the front value without removing it.
func (q *Queue[T]) Peek() (T, bool) {
	var zero T
	if q.count == 0 {
		return zero, false
	}
	return q.data[q.head], true
}

// grow doubles capacity and re-linearizes the elements from head.
func (q *Queue[T]) grow() {
	newCap := len(q.data) * 2
	if newCap == 0 {
		newCap = 4
	}
	buf := make([]T, newCap)
	for i := 0; i < q.count; i++ {
		buf[i] = q.data[(q.head+i)%len(q.data)]
	}
	q.data = buf
	q.head = 0
}

// LinkedQueue is a FIFO queue backed by a singly linked list with
// head and tail pointers. The zero value is ready to use.
type LinkedQueue[T any] struct {
	head *qnode[T]
	tail *qnode[T]
	size int
}

type qnode[T any] struct {
	value T
	next  *qnode[T]
}

func (q *LinkedQueue[T]) Len() int      { return q.size }
func (q *LinkedQueue[T]) IsEmpty() bool { return q.size == 0 }

// Enqueue adds v at the tail in O(1).
func (q *LinkedQueue[T]) Enqueue(v T) {
	n := &qnode[T]{value: v}
	if q.tail == nil {
		q.head, q.tail = n, n
	} else {
		q.tail.next = n
		q.tail = n
	}
	q.size++
}

// Dequeue removes and returns the head value in O(1).
func (q *LinkedQueue[T]) Dequeue() (T, bool) {
	var zero T
	if q.head == nil {
		return zero, false
	}
	n := q.head
	q.head = n.next
	if q.head == nil {
		q.tail = nil
	}
	q.size--
	return n.value, true
}
```

## usage / test

```go
package queue

import "testing"

func TestQueueFIFO(t *testing.T) {
	q := New[int](2)
	for i := 1; i <= 5; i++ { // forces a grow past initial capacity
		q.Enqueue(i)
	}
	if q.Len() != 5 {
		t.Fatalf("len: got %d want 5", q.Len())
	}
	if v, _ := q.Peek(); v != 1 {
		t.Fatalf("peek: got %d want 1", v)
	}
	// Elements must come out in arrival order.
	for want := 1; want <= 5; want++ {
		got, ok := q.Dequeue()
		if !ok || got != want {
			t.Fatalf("dequeue: got %d want %d", got, want)
		}
	}
	if _, ok := q.Dequeue(); ok {
		t.Fatal("dequeue on empty must return false")
	}
}

func TestQueueWrapAround(t *testing.T) {
	q := New[int](4)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Dequeue()      // head advances
	q.Enqueue(3)
	q.Enqueue(4)
	q.Enqueue(5)     // wraps around the ring
	want := []int{2, 3, 4, 5}
	for _, w := range want {
		if got, _ := q.Dequeue(); got != w {
			t.Fatalf("got %d want %d", got, w)
		}
	}
}

func TestLinkedQueue(t *testing.T) {
	var q LinkedQueue[string]
	q.Enqueue("a")
	q.Enqueue("b")
	if v, _ := q.Dequeue(); v != "a" {
		t.Fatalf("got %q want a", v)
	}
	if v, _ := q.Dequeue(); v != "b" {
		t.Fatalf("got %q want b", v)
	}
}
```
