# Binary Search Tree (BST: Insert, Search, Delete, Min, Max)

A binary search tree is an ordered binary tree where every node's left subtree holds only smaller keys and its right subtree only larger keys. Use it for a dynamic ordered set or map that supports fast lookup, insertion, deletion, minimum, maximum, predecessor/successor, and in-order iteration that yields keys in sorted order. Search, insert, and delete take O(h) time where h is the height: O(log n) on a balanced tree but O(n) in the worst case if keys are inserted in sorted order (for guaranteed O(log n), use an AVL or red-black tree). Space is O(n). Keywords: bst binary search tree ordered set ordered map dictionary insert add put search find lookup contains has delete remove erase min minimum max maximum successor predecessor inorder sorted ordered traversal balanced height log n comparison key.

## implementation

```go
package bst

import "cmp"

// node is a single BST node holding an ordered key.
type node[T cmp.Ordered] struct {
	key         T
	left, right *node[T]
}

// BST is a binary search tree of ordered keys.
type BST[T cmp.Ordered] struct {
	root *node[T]
	size int
}

// New returns an empty BST.
func New[T cmp.Ordered]() *BST[T] { return &BST[T]{} }

// Len returns the number of keys stored.
func (t *BST[T]) Len() int { return t.size }

// Insert adds key if it is not already present. Duplicates are ignored.
func (t *BST[T]) Insert(key T) {
	var inserted bool
	t.root, inserted = insert(t.root, key)
	if inserted {
		t.size++
	}
}

func insert[T cmp.Ordered](n *node[T], key T) (*node[T], bool) {
	if n == nil {
		return &node[T]{key: key}, true
	}
	var ins bool
	switch {
	case key < n.key:
		n.left, ins = insert(n.left, key)
	case key > n.key:
		n.right, ins = insert(n.right, key)
	default:
		ins = false // already present
	}
	return n, ins
}

// Contains reports whether key is in the tree.
func (t *BST[T]) Contains(key T) bool {
	n := t.root
	for n != nil {
		switch {
		case key < n.key:
			n = n.left
		case key > n.key:
			n = n.right
		default:
			return true
		}
	}
	return false
}

// Min returns the smallest key and true, or the zero value and false if empty.
func (t *BST[T]) Min() (T, bool) {
	var zero T
	if t.root == nil {
		return zero, false
	}
	return minNode(t.root).key, true
}

func minNode[T cmp.Ordered](n *node[T]) *node[T] {
	for n.left != nil {
		n = n.left
	}
	return n
}

// Max returns the largest key and true, or the zero value and false if empty.
func (t *BST[T]) Max() (T, bool) {
	var zero T
	if t.root == nil {
		return zero, false
	}
	n := t.root
	for n.right != nil {
		n = n.right
	}
	return n.key, true
}

// Delete removes key if present and reports whether it was found.
func (t *BST[T]) Delete(key T) bool {
	var removed bool
	t.root, removed = remove(t.root, key)
	if removed {
		t.size--
	}
	return removed
}

func remove[T cmp.Ordered](n *node[T], key T) (*node[T], bool) {
	if n == nil {
		return nil, false
	}
	var removed bool
	switch {
	case key < n.key:
		n.left, removed = remove(n.left, key)
	case key > n.key:
		n.right, removed = remove(n.right, key)
	default:
		removed = true
		// Cases: zero or one child -> splice it out.
		if n.left == nil {
			return n.right, true
		}
		if n.right == nil {
			return n.left, true
		}
		// Two children: replace key with in-order successor (min of right
		// subtree), then delete that successor from the right subtree.
		succ := minNode(n.right)
		n.key = succ.key
		n.right, _ = remove(n.right, succ.key)
	}
	return n, removed
}

// InOrder returns all keys in ascending sorted order.
func (t *BST[T]) InOrder() []T {
	var out []T
	var walk func(n *node[T])
	walk = func(n *node[T]) {
		if n == nil {
			return
		}
		walk(n.left)
		out = append(out, n.key)
		walk(n.right)
	}
	walk(t.root)
	return out
}
```

## usage / test

```go
package bst

import (
	"reflect"
	"sort"
	"testing"
)

func TestBST(t *testing.T) {
	tree := New[int]()
	keys := []int{5, 3, 8, 1, 4, 7, 9, 2, 6}
	for _, k := range keys {
		tree.Insert(k)
	}
	tree.Insert(5) // duplicate ignored

	// In-order traversal of a BST must be sorted.
	want := append([]int(nil), keys...)
	sort.Ints(want)
	if got := tree.InOrder(); !reflect.DeepEqual(got, want) {
		t.Fatalf("InOrder = %v, want %v", got, want)
	}
	if tree.Len() != len(want) {
		t.Fatalf("Len = %d, want %d", tree.Len(), len(want))
	}

	if mn, _ := tree.Min(); mn != 1 {
		t.Errorf("Min = %d, want 1", mn)
	}
	if mx, _ := tree.Max(); mx != 9 {
		t.Errorf("Max = %d, want 9", mx)
	}
	if !tree.Contains(7) || tree.Contains(42) {
		t.Errorf("Contains wrong")
	}

	// Delete a node with two children and confirm order stays sorted.
	if !tree.Delete(5) {
		t.Fatal("Delete(5) returned false")
	}
	if tree.Contains(5) {
		t.Error("5 still present after delete")
	}
	got := tree.InOrder()
	if !sort.IntsAreSorted(got) {
		t.Errorf("after delete not sorted: %v", got)
	}
}
```
