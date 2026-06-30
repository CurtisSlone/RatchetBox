# Set (Hash Set of Unique Elements)

A set is an unordered collection of unique elements that answers membership questions and supports the algebraic operations union, intersection, and difference. In Go it is idiomatically built on a `map[T]struct{}`, where the empty struct value costs zero bytes and the key carries the data. Use a set to deduplicate items, test "have I seen this?" in a loop (visited sets in graph traversal), or compute relationships between two collections; add, remove, and contains are average O(1), iteration is O(n), and the binary operations are O(n+m). Keywords: set hash set unique deduplicate dedup membership contains has add insert remove delete union intersection difference symmetric-difference subset superset disjoint map[T]struct{} visited seen comparable cardinality size elements collection

## implementation

```go
package set

// Set is an unordered collection of unique comparable elements.
type Set[T comparable] map[T]struct{}

// New returns a set seeded with the given elements.
func New[T comparable](elems ...T) Set[T] {
	s := make(Set[T], len(elems))
	s.Add(elems...)
	return s
}

// Add inserts elements; duplicates are ignored. Average O(1) each.
func (s Set[T]) Add(elems ...T) {
	for _, e := range elems {
		s[e] = struct{}{}
	}
}

// Remove deletes an element if present.
func (s Set[T]) Remove(e T) { delete(s, e) }

// Contains reports membership in average O(1).
func (s Set[T]) Contains(e T) bool {
	_, ok := s[e]
	return ok
}

// Len reports the number of elements (cardinality).
func (s Set[T]) Len() int { return len(s) }

// Elements returns the members in unspecified order.
func (s Set[T]) Elements() []T {
	out := make([]T, 0, len(s))
	for e := range s {
		out = append(out, e)
	}
	return out
}

// Union returns the elements present in either set.
func (s Set[T]) Union(other Set[T]) Set[T] {
	out := make(Set[T], len(s)+len(other))
	for e := range s {
		out[e] = struct{}{}
	}
	for e := range other {
		out[e] = struct{}{}
	}
	return out
}

// Intersect returns the elements present in both sets.
func (s Set[T]) Intersect(other Set[T]) Set[T] {
	// Iterate the smaller set for efficiency.
	small, large := s, other
	if len(other) < len(s) {
		small, large = other, s
	}
	out := make(Set[T])
	for e := range small {
		if _, ok := large[e]; ok {
			out[e] = struct{}{}
		}
	}
	return out
}

// Difference returns the elements in s that are not in other.
func (s Set[T]) Difference(other Set[T]) Set[T] {
	out := make(Set[T])
	for e := range s {
		if _, ok := other[e]; !ok {
			out[e] = struct{}{}
		}
	}
	return out
}

// IsSubset reports whether every element of s is in other.
func (s Set[T]) IsSubset(other Set[T]) bool {
	if len(s) > len(other) {
		return false
	}
	for e := range s {
		if _, ok := other[e]; !ok {
			return false
		}
	}
	return true
}
```

## usage / test

```go
package set

import (
	"sort"
	"testing"
)

func sorted(s Set[int]) []int {
	out := s.Elements()
	sort.Ints(out)
	return out
}

func eq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSetOps(t *testing.T) {
	a := New(1, 2, 3, 3) // duplicate ignored
	b := New(2, 3, 4)

	if a.Len() != 3 {
		t.Fatalf("dedup len: got %d want 3", a.Len())
	}
	if !a.Contains(2) || a.Contains(9) {
		t.Fatal("membership wrong")
	}
	if got := sorted(a.Union(b)); !eq(got, []int{1, 2, 3, 4}) {
		t.Fatalf("union: got %v", got)
	}
	if got := sorted(a.Intersect(b)); !eq(got, []int{2, 3}) {
		t.Fatalf("intersect: got %v", got)
	}
	if got := sorted(a.Difference(b)); !eq(got, []int{1}) {
		t.Fatalf("difference: got %v", got)
	}
	if !New(2, 3).IsSubset(b) {
		t.Fatal("expected {2,3} subset of {2,3,4}")
	}
}
```
