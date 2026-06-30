# Doubly Linked List

A doubly linked list is a linear sequence of nodes where each node carries a value plus two pointers, one to the next node and one to the previous node, with the list holding head and tail pointers. Use it when you need O(1) insertion and removal at both ends or at a known node (it backs LRU caches, the standard library `container/list`, and deques), or when you must walk the sequence in either direction; avoid it when you need random indexed access (O(n)) or want minimal per-node memory. Push/pop at front or back and unlinking a known node are O(1), search and indexing are O(n), and space is O(n) with two pointers of overhead per node. Keywords: doubly linked list double linked list two-way list node next prev previous pointer head tail push-front push-back pop-front pop-back prepend append insert-before insert-after unlink delete remove reverse bidirectional traverse forward backward container/list sentinel LRU deque

## implementation

```go
package dll

// Node is an element of a doubly linked list.
type Node[T any] struct {
	Value T
	next  *Node[T]
	prev  *Node[T]
}

// Next returns the following node, or nil at the tail.
func (n *Node[T]) Next() *Node[T] { return n.next }

// Prev returns the preceding node, or nil at the head.
func (n *Node[T]) Prev() *Node[T] { return n.prev }

// List is a doubly linked list with head and tail pointers.
type List[T any] struct {
	head *Node[T]
	tail *Node[T]
	size int
}

// New returns an empty list.
func New[T any]() *List[T] { return &List[T]{} }

// Len reports the element count in O(1).
func (l *List[T]) Len() int { return l.size }

// Front and Back return the boundary nodes (nil if empty).
func (l *List[T]) Front() *Node[T] { return l.head }
func (l *List[T]) Back() *Node[T]  { return l.tail }

// PushFront inserts a value at the head and returns its node, O(1).
func (l *List[T]) PushFront(v T) *Node[T] {
	n := &Node[T]{Value: v, next: l.head}
	if l.head != nil {
		l.head.prev = n
	} else {
		l.tail = n
	}
	l.head = n
	l.size++
	return n
}

// PushBack appends a value at the tail and returns its node, O(1).
func (l *List[T]) PushBack(v T) *Node[T] {
	n := &Node[T]{Value: v, prev: l.tail}
	if l.tail != nil {
		l.tail.next = n
	} else {
		l.head = n
	}
	l.tail = n
	l.size++
	return n
}

// PopFront removes and returns the head value in O(1).
func (l *List[T]) PopFront() (T, bool) {
	if l.head == nil {
		var zero T
		return zero, false
	}
	n := l.head
	l.Remove(n)
	return n.Value, true
}

// PopBack removes and returns the tail value in O(1).
func (l *List[T]) PopBack() (T, bool) {
	if l.tail == nil {
		var zero T
		return zero, false
	}
	n := l.tail
	l.Remove(n)
	return n.Value, true
}

// Remove unlinks a node that belongs to this list in O(1).
func (l *List[T]) Remove(n *Node[T]) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		l.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		l.tail = n.prev
	}
	n.next, n.prev = nil, nil
	l.size--
}

// InsertAfter inserts v after node at and returns the new node, O(1).
func (l *List[T]) InsertAfter(v T, at *Node[T]) *Node[T] {
	n := &Node[T]{Value: v, prev: at, next: at.next}
	if at.next != nil {
		at.next.prev = n
	} else {
		l.tail = n
	}
	at.next = n
	l.size++
	return n
}

// Slice returns values from head to tail.
func (l *List[T]) Slice() []T {
	out := make([]T, 0, l.size)
	for cur := l.head; cur != nil; cur = cur.next {
		out = append(out, cur.Value)
	}
	return out
}

// SliceReverse returns values from tail to head.
func (l *List[T]) SliceReverse() []T {
	out := make([]T, 0, l.size)
	for cur := l.tail; cur != nil; cur = cur.prev {
		out = append(out, cur.Value)
	}
	return out
}
```

## usage / test

```go
package dll

import (
	"reflect"
	"testing"
)

func TestDoublyLinkedList(t *testing.T) {
	l := New[int]()
	l.PushBack(2)
	mid := l.PushBack(3)
	l.PushFront(1)        // 1 2 3
	l.InsertAfter(99, mid) // 1 2 3 99

	if got := l.Slice(); !reflect.DeepEqual(got, []int{1, 2, 3, 99}) {
		t.Fatalf("forward: got %v", got)
	}
	if got := l.SliceReverse(); !reflect.DeepEqual(got, []int{99, 3, 2, 1}) {
		t.Fatalf("reverse: got %v", got)
	}

	l.Remove(mid) // 1 2 99
	if got := l.Slice(); !reflect.DeepEqual(got, []int{1, 2, 99}) {
		t.Fatalf("after remove: got %v", got)
	}

	if v, ok := l.PopFront(); !ok || v != 1 {
		t.Fatalf("pop front: got %v %v", v, ok)
	}
	if v, ok := l.PopBack(); !ok || v != 99 {
		t.Fatalf("pop back: got %v %v", v, ok)
	}
	if l.Len() != 1 {
		t.Fatalf("len: got %d want 1", l.Len())
	}
}
```
