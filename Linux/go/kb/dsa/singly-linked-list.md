# Singly Linked List

A singly linked list is a linear sequence of nodes where each node holds a value and a single pointer to the next node; the list keeps a head pointer (and often a tail pointer and length counter). Use it when you need O(1) insertion or removal at the front (or at a known node), a stack/queue backbone, or when you cannot afford the contiguous reallocation that a slice/array needs; avoid it when you need random indexed access, which is O(n). Inserting at the head is O(1), appending at the tail is O(1) with a tail pointer (O(n) without), searching or indexing is O(n), and space is O(n) plus one pointer of overhead per node. Keywords: singly linked list single linked list one-way list node next pointer head tail prepend append push-front insert-at-head insert-at-tail add-to-head add-to-end delete remove-by-value remove-by-index traverse iterate search reverse length size container/list sequence chain

## implementation

```go
package linkedlist

// Node is a single element of a singly linked list.
type Node[T any] struct {
	Value T
	Next  *Node[T]
}

// List is a singly linked list with head and tail pointers and a size counter.
type List[T any] struct {
	head *Node[T]
	tail *Node[T]
	size int
}

// New returns an empty list.
func New[T any]() *List[T] { return &List[T]{} }

// Len reports the number of elements in O(1).
func (l *List[T]) Len() int { return l.size }

// IsEmpty reports whether the list has no elements.
func (l *List[T]) IsEmpty() bool { return l.size == 0 }

// PushFront inserts a value at the head in O(1).
func (l *List[T]) PushFront(v T) {
	n := &Node[T]{Value: v, Next: l.head}
	l.head = n
	if l.tail == nil {
		l.tail = n
	}
	l.size++
}

// PushBack appends a value at the tail in O(1).
func (l *List[T]) PushBack(v T) {
	n := &Node[T]{Value: v}
	if l.tail == nil {
		l.head, l.tail = n, n
	} else {
		l.tail.Next = n
		l.tail = n
	}
	l.size++
}

// PopFront removes and returns the head value in O(1).
func (l *List[T]) PopFront() (T, bool) {
	var zero T
	if l.head == nil {
		return zero, false
	}
	n := l.head
	l.head = n.Next
	if l.head == nil {
		l.tail = nil
	}
	n.Next = nil
	l.size--
	return n.Value, true
}

// Front returns the head value without removing it.
func (l *List[T]) Front() (T, bool) {
	var zero T
	if l.head == nil {
		return zero, false
	}
	return l.head.Value, true
}

// Remove deletes the first node whose value satisfies eq in O(n).
func (l *List[T]) Remove(target T, eq func(a, b T) bool) bool {
	var prev *Node[T]
	for cur := l.head; cur != nil; cur = cur.Next {
		if eq(cur.Value, target) {
			if prev == nil {
				l.head = cur.Next
			} else {
				prev.Next = cur.Next
			}
			if cur == l.tail {
				l.tail = prev
			}
			cur.Next = nil
			l.size--
			return true
		}
		prev = cur
	}
	return false
}

// Reverse reverses the list in place in O(n) time and O(1) extra space.
func (l *List[T]) Reverse() {
	var prev *Node[T]
	l.tail = l.head
	cur := l.head
	for cur != nil {
		next := cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}
	l.head = prev
}

// Slice returns the values from head to tail.
func (l *List[T]) Slice() []T {
	out := make([]T, 0, l.size)
	for cur := l.head; cur != nil; cur = cur.Next {
		out = append(out, cur.Value)
	}
	return out
}
```

## usage / test

```go
package linkedlist

import (
	"reflect"
	"testing"
)

func TestSinglyLinkedList(t *testing.T) {
	l := New[int]()
	if !l.IsEmpty() {
		t.Fatal("new list should be empty")
	}

	l.PushBack(1)
	l.PushBack(2)
	l.PushFront(0) // 0 1 2

	if got := l.Slice(); !reflect.DeepEqual(got, []int{0, 1, 2}) {
		t.Fatalf("order: got %v", got)
	}
	if l.Len() != 3 {
		t.Fatalf("len: got %d want 3", l.Len())
	}

	eq := func(a, b int) bool { return a == b }
	if !l.Remove(1, eq) {
		t.Fatal("expected to remove 1")
	}
	if got := l.Slice(); !reflect.DeepEqual(got, []int{0, 2}) {
		t.Fatalf("after remove: got %v", got)
	}

	l.Reverse() // 2 0
	if got := l.Slice(); !reflect.DeepEqual(got, []int{2, 0}) {
		t.Fatalf("after reverse: got %v", got)
	}

	if v, ok := l.PopFront(); !ok || v != 2 {
		t.Fatalf("pop front: got %v %v", v, ok)
	}
	if v, ok := l.Front(); !ok || v != 0 {
		t.Fatalf("front: got %v %v", v, ok)
	}
}
```
