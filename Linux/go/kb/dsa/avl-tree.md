# AVL Tree (Self-Balancing Binary Search Tree)

An AVL tree is a binary search tree that keeps itself height-balanced: for every node the heights of its left and right subtrees differ by at most one. After each insert or delete it restores balance with rotations (left, right, left-right, right-left), guaranteeing O(log n) search, insert, and delete in the worst case — unlike a plain BST which can degrade to O(n). Each node stores its subtree height (or a balance factor) so imbalance can be detected and fixed on the way back up the recursion. Use it when you need a strictly balanced ordered set/map with reliably fast operations and relatively more lookups than updates (AVL is more rigidly balanced than a red-black tree, giving slightly faster lookups but more rotations on writes). Space is O(n). Keywords: avl tree self-balancing balanced binary search tree rotation rotate left right left-right right-left balance factor height insert delete search ordered set map guaranteed log n worst case height-balanced adelson-velsky landis rebalance.

## implementation

```go
package avl

import "cmp"

type node[T cmp.Ordered] struct {
	key         T
	height      int
	left, right *node[T]
}

// AVL is a self-balancing binary search tree.
type AVL[T cmp.Ordered] struct {
	root *node[T]
	size int
}

func New[T cmp.Ordered]() *AVL[T] { return &AVL[T]{} }

func (t *AVL[T]) Len() int { return t.size }

func height[T cmp.Ordered](n *node[T]) int {
	if n == nil {
		return 0
	}
	return n.height
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func update[T cmp.Ordered](n *node[T]) {
	n.height = 1 + max(height(n.left), height(n.right))
}

// balanceFactor > 1 means left-heavy, < -1 means right-heavy.
func balanceFactor[T cmp.Ordered](n *node[T]) int {
	if n == nil {
		return 0
	}
	return height(n.left) - height(n.right)
}

func rotateRight[T cmp.Ordered](y *node[T]) *node[T] {
	x := y.left
	y.left = x.right
	x.right = y
	update(y)
	update(x)
	return x
}

func rotateLeft[T cmp.Ordered](x *node[T]) *node[T] {
	y := x.right
	x.right = y.left
	y.left = x
	update(x)
	update(y)
	return y
}

func rebalance[T cmp.Ordered](n *node[T]) *node[T] {
	update(n)
	bf := balanceFactor(n)
	if bf > 1 { // left-heavy
		if balanceFactor(n.left) < 0 {
			n.left = rotateLeft(n.left) // left-right case
		}
		return rotateRight(n)
	}
	if bf < -1 { // right-heavy
		if balanceFactor(n.right) > 0 {
			n.right = rotateRight(n.right) // right-left case
		}
		return rotateLeft(n)
	}
	return n
}

// Insert adds key, ignoring duplicates, and rebalances.
func (t *AVL[T]) Insert(key T) {
	var added bool
	t.root, added = insert(t.root, key)
	if added {
		t.size++
	}
}

func insert[T cmp.Ordered](n *node[T], key T) (*node[T], bool) {
	if n == nil {
		return &node[T]{key: key, height: 1}, true
	}
	var added bool
	switch {
	case key < n.key:
		n.left, added = insert(n.left, key)
	case key > n.key:
		n.right, added = insert(n.right, key)
	default:
		return n, false
	}
	return rebalance(n), added
}

// Contains reports whether key is present.
func (t *AVL[T]) Contains(key T) bool {
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

// Delete removes key if present and rebalances.
func (t *AVL[T]) Delete(key T) bool {
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
		if n.left == nil {
			return n.right, true
		}
		if n.right == nil {
			return n.left, true
		}
		succ := n.right
		for succ.left != nil {
			succ = succ.left
		}
		n.key = succ.key
		n.right, _ = remove(n.right, succ.key)
	}
	return rebalance(n), removed
}

// InOrder returns all keys in ascending sorted order.
func (t *AVL[T]) InOrder() []T {
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

// Height returns the height of the whole tree (0 for empty).
func (t *AVL[T]) Height() int { return height(t.root) }
```

## usage / test

```go
package avl

import (
	"sort"
	"testing"
)

func TestAVLBalancedAndSorted(t *testing.T) {
	tree := New[int]()
	// Inserting already-sorted keys would make a plain BST degrade to a list;
	// the AVL tree must stay shallow.
	const n = 1023
	for i := 1; i <= n; i++ {
		tree.Insert(i)
	}
	if tree.Len() != n {
		t.Fatalf("Len = %d, want %d", tree.Len(), n)
	}
	// For n=1023 a perfectly balanced tree has height 10; AVL stays close.
	if h := tree.Height(); h > 14 {
		t.Fatalf("height %d too large, tree is not balanced", h)
	}
	if got := tree.InOrder(); !sort.IntsAreSorted(got) {
		t.Fatal("InOrder not sorted")
	}
	// Delete half the keys; tree must remain sorted and balanced.
	for i := 1; i <= n; i += 2 {
		if !tree.Delete(i) {
			t.Fatalf("Delete(%d) returned false", i)
		}
	}
	if !sort.IntsAreSorted(tree.InOrder()) {
		t.Fatal("not sorted after deletes")
	}
	if h := tree.Height(); h > 14 {
		t.Fatalf("height %d too large after deletes", h)
	}
}
```
