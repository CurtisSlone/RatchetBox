# Circular Linked List

A circular linked list is a linked list whose last node points back to the first node instead of to nil, so iteration can loop forever and there is no natural end; it is usually implemented with a single tail pointer because tail.Next is the head, giving O(1) access to both ends. Use it for round-robin scheduling, repeating playlists or buffers, and any problem (like the Josephus problem) where you advance cyclically through elements; it can be singly or doubly linked. Insert at front or back is O(1) with a tail pointer, deletion of a known node is O(1) for a doubly circular list or O(n) to find a value in a singly circular list, traversal of all n elements is O(n), and space is O(n). Keywords: circular linked list circularly linked ring list round-robin cyclic loop tail-points-to-head josephus next prev rotate advance push-front push-back insert delete traverse no-nil-terminator playlist scheduler

## implementation

```go
package circular

// node is one element; in a non-empty list next never points to nil.
type node[T any] struct {
	value T
	next  *node[T]
}

// List is a singly circular linked list tracked by its tail node.
// tail.next is always the head, giving O(1) front and back access.
type List[T any] struct {
	tail *node[T]
	size int
}

// New returns an empty circular list.
func New[T any]() *List[T] { return &List[T]{} }

// Len reports the element count in O(1).
func (l *List[T]) Len() int { return l.size }

// IsEmpty reports whether the list has no elements.
func (l *List[T]) IsEmpty() bool { return l.size == 0 }

// PushBack appends a value after the tail in O(1).
func (l *List[T]) PushBack(v T) {
	n := &node[T]{value: v}
	if l.tail == nil {
		n.next = n // single node points to itself
	} else {
		n.next = l.tail.next // new node points to head
		l.tail.next = n
	}
	l.tail = n
	l.size++
}

// PushFront inserts a value before the head in O(1).
func (l *List[T]) PushFront(v T) {
	oldTail := l.tail
	l.PushBack(v)
	// PushBack appended after the old tail and made the new node the tail.
	// Keep the old tail as the tail so the new node becomes the head.
	if oldTail != nil {
		l.tail = oldTail
	}
}

// PopFront removes and returns the head value in O(1).
func (l *List[T]) PopFront() (T, bool) {
	var zero T
	if l.tail == nil {
		return zero, false
	}
	head := l.tail.next
	if head == l.tail { // single element
		l.tail = nil
	} else {
		l.tail.next = head.next
	}
	l.size--
	return head.value, true
}

// Do calls fn once per element starting at the head, in order.
func (l *List[T]) Do(fn func(T)) {
	if l.tail == nil {
		return
	}
	cur := l.tail.next
	for i := 0; i < l.size; i++ {
		fn(cur.value)
		cur = cur.next
	}
}

// Slice returns the values starting from the head.
func (l *List[T]) Slice() []T {
	out := make([]T, 0, l.size)
	l.Do(func(v T) { out = append(out, v) })
	return out
}

// Josephus eliminates every k-th element cyclically and returns the
// value of the survivor, demonstrating circular traversal.
func Josephus[T any](l *List[T], k int) (T, bool) {
	var zero T
	if l.tail == nil || k < 1 {
		return zero, false
	}
	for l.size > 1 {
		for i := 1; i < k; i++ {
			l.tail = l.tail.next // advance the "previous" pointer
		}
		// l.tail.next is the k-th node from the previous survivor; drop it.
		l.tail.next = l.tail.next.next
		l.size--
	}
	return l.tail.value, true
}
```

## usage / test

```go
package circular

import (
	"reflect"
	"testing"
)

func TestCircularList(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.PushFront(0) // 0 1 2 3

	if got := l.Slice(); !reflect.DeepEqual(got, []int{0, 1, 2, 3}) {
		t.Fatalf("order: got %v", got)
	}

	// Tail wraps back to head: walking 2*size elements stays bounded.
	count := 0
	cur := l.tail
	for i := 0; i < 2*l.Len(); i++ {
		cur = cur.next
		count++
	}
	if count != 2*l.Len() {
		t.Fatalf("circular walk count: got %d", count)
	}

	if v, ok := l.PopFront(); !ok || v != 0 {
		t.Fatalf("pop front: got %v %v", v, ok)
	}
}

func TestJosephus(t *testing.T) {
	l := New[int]()
	for i := 1; i <= 5; i++ {
		l.PushBack(i)
	}
	// Classic n=5, k=2 -> survivor is 3.
	if v, ok := Josephus(l, 2); !ok || v != 3 {
		t.Fatalf("josephus survivor: got %v %v want 3", v, ok)
	}
}
```
