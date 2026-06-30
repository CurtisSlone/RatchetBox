# Stack (LIFO, Array-Backed and Linked-List-Backed)

A stack is a last-in-first-out (LIFO) container that supports pushing a value on top, popping the most recently pushed value off the top, and peeking at the top without removing it. Use it for function-call and recursion unwinding, expression parsing and bracket matching, undo history, depth-first search, and backtracking; reach for it whenever the most recently added item is the next one you need. Push, pop, and peek are all amortized O(1); a slice-backed stack stores values contiguously and grows by doubling, while a linked-list-backed stack allocates a node per element and never reallocates. Keywords: stack LIFO last-in-first-out push pop peek top isEmpty size array-backed slice-backed linked-list-backed call-stack backtracking DFS undo bracket-matching parentheses balanced container drop pop-front

## implementation

```go
package stack

// Stack is a slice-backed LIFO stack. The zero value is ready to use.
type Stack[T any] struct {
	data []T
}

// New returns an empty stack; capacity is an optional hint.
func New[T any](capacity int) *Stack[T] {
	return &Stack[T]{data: make([]T, 0, capacity)}
}

// Len reports the number of elements.
func (s *Stack[T]) Len() int { return len(s.data) }

// IsEmpty reports whether the stack has no elements.
func (s *Stack[T]) IsEmpty() bool { return len(s.data) == 0 }

// Push adds v on top in amortized O(1).
func (s *Stack[T]) Push(v T) { s.data = append(s.data, v) }

// Pop removes and returns the top value in O(1).
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	n := len(s.data)
	if n == 0 {
		return zero, false
	}
	v := s.data[n-1]
	s.data[n-1] = zero // release reference for the GC
	s.data = s.data[:n-1]
	return v, true
}

// Peek returns the top value without removing it.
func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	n := len(s.data)
	if n == 0 {
		return zero, false
	}
	return s.data[n-1], true
}

// LinkedStack is a linked-list-backed LIFO stack: no reallocation,
// one allocation per element. The zero value is ready to use.
type LinkedStack[T any] struct {
	top  *snode[T]
	size int
}

type snode[T any] struct {
	value T
	next  *snode[T]
}

// Len reports the number of elements.
func (s *LinkedStack[T]) Len() int { return s.size }

// IsEmpty reports whether the stack has no elements.
func (s *LinkedStack[T]) IsEmpty() bool { return s.size == 0 }

// Push adds v on top in O(1).
func (s *LinkedStack[T]) Push(v T) {
	s.top = &snode[T]{value: v, next: s.top}
	s.size++
}

// Pop removes and returns the top value in O(1).
func (s *LinkedStack[T]) Pop() (T, bool) {
	var zero T
	if s.top == nil {
		return zero, false
	}
	n := s.top
	s.top = n.next
	s.size--
	return n.value, true
}

// Peek returns the top value without removing it.
func (s *LinkedStack[T]) Peek() (T, bool) {
	var zero T
	if s.top == nil {
		return zero, false
	}
	return s.top.value, true
}
```

## usage / test

```go
package stack

import "testing"

func TestStackLIFO(t *testing.T) {
	s := New[int](0)
	if !s.IsEmpty() {
		t.Fatal("new stack should be empty")
	}
	for _, v := range []int{1, 2, 3} {
		s.Push(v)
	}
	if top, _ := s.Peek(); top != 3 {
		t.Fatalf("peek: got %d want 3", top)
	}
	// Items must come out in reverse insertion order: 3, 2, 1.
	for _, want := range []int{3, 2, 1} {
		got, ok := s.Pop()
		if !ok || got != want {
			t.Fatalf("pop: got %d want %d", got, want)
		}
	}
	if _, ok := s.Pop(); ok {
		t.Fatal("pop on empty must return false")
	}
}

// BalancedParens demonstrates a classic stack use: bracket matching.
func BalancedParens(s string) bool {
	st := New[rune](len(s))
	pairs := map[rune]rune{')': '(', ']': '[', '}': '{'}
	for _, r := range s {
		switch r {
		case '(', '[', '{':
			st.Push(r)
		case ')', ']', '}':
			top, ok := st.Pop()
			if !ok || top != pairs[r] {
				return false
			}
		}
	}
	return st.IsEmpty()
}

func TestBalancedParens(t *testing.T) {
	cases := map[string]bool{"(a[b]{c})": true, "(]": false, "((": false, "": true}
	for in, want := range cases {
		if got := BalancedParens(in); got != want {
			t.Fatalf("BalancedParens(%q): got %v want %v", in, got, want)
		}
	}
}

func TestLinkedStack(t *testing.T) {
	var s LinkedStack[string]
	s.Push("a")
	s.Push("b")
	if v, _ := s.Pop(); v != "b" {
		t.Fatalf("got %q want b", v)
	}
	if v, _ := s.Pop(); v != "a" {
		t.Fatalf("got %q want a", v)
	}
}
```
