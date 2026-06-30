# Red-Black Tree (Self-Balancing Binary Search Tree)

A red-black tree is a self-balancing binary search tree where every node is colored red or black and a set of color rules keeps the longest root-to-leaf path at most twice the shortest, guaranteeing O(log n) search, insert, and delete in the worst case. The invariants are: the root is black, red nodes have black children (no two reds in a row), and every root-to-nil path passes through the same number of black nodes (the "black height"). After inserting (always as a red node) it restores the invariants with recolorings and rotations. Red-black trees do fewer rotations on updates than AVL trees, so they are the common choice for general-purpose ordered maps and sets (e.g. library map/set implementations). Space is O(n). Keywords: red-black tree red black rb tree self-balancing balanced binary search tree rotation recolor color invariant black height ordered map ordered set insert delete search guaranteed log n worst case rebalance sentinel nil leaf.

## implementation

```go
package rbtree

import "cmp"

type color bool

const (
	red   color = true
	black color = false
)

type node[T cmp.Ordered] struct {
	key                 T
	color               color
	left, right, parent *node[T]
}

// RBTree is a red-black self-balancing binary search tree.
// It uses a shared sentinel nil node colored black to simplify the logic.
type RBTree[T cmp.Ordered] struct {
	root *node[T]
	nul  *node[T] // sentinel for nil
	size int
}

func New[T cmp.Ordered]() *RBTree[T] {
	nul := &node[T]{color: black}
	return &RBTree[T]{root: nul, nul: nul}
}

func (t *RBTree[T]) Len() int { return t.size }

func (t *RBTree[T]) leftRotate(x *node[T]) {
	y := x.right
	x.right = y.left
	if y.left != t.nul {
		y.left.parent = x
	}
	y.parent = x.parent
	switch {
	case x.parent == t.nul:
		t.root = y
	case x == x.parent.left:
		x.parent.left = y
	default:
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func (t *RBTree[T]) rightRotate(x *node[T]) {
	y := x.left
	x.left = y.right
	if y.right != t.nul {
		y.right.parent = x
	}
	y.parent = x.parent
	switch {
	case x.parent == t.nul:
		t.root = y
	case x == x.parent.right:
		x.parent.right = y
	default:
		x.parent.left = y
	}
	y.right = x
	x.parent = y
}

// Insert adds key (ignoring duplicates) and restores the red-black invariants.
func (t *RBTree[T]) Insert(key T) {
	parent := t.nul
	cur := t.root
	for cur != t.nul {
		parent = cur
		switch {
		case key < cur.key:
			cur = cur.left
		case key > cur.key:
			cur = cur.right
		default:
			return // duplicate
		}
	}
	z := &node[T]{key: key, color: red, left: t.nul, right: t.nul, parent: parent}
	switch {
	case parent == t.nul:
		t.root = z
	case key < parent.key:
		parent.left = z
	default:
		parent.right = z
	}
	t.size++
	t.insertFixup(z)
}

func (t *RBTree[T]) insertFixup(z *node[T]) {
	for z.parent.color == red {
		if z.parent == z.parent.parent.left {
			uncle := z.parent.parent.right
			if uncle.color == red {
				z.parent.color = black
				uncle.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					z = z.parent
					t.leftRotate(z)
				}
				z.parent.color = black
				z.parent.parent.color = red
				t.rightRotate(z.parent.parent)
			}
		} else {
			uncle := z.parent.parent.left
			if uncle.color == red {
				z.parent.color = black
				uncle.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					t.rightRotate(z)
				}
				z.parent.color = black
				z.parent.parent.color = red
				t.leftRotate(z.parent.parent)
			}
		}
	}
	t.root.color = black
}

// Contains reports whether key is present.
func (t *RBTree[T]) Contains(key T) bool {
	cur := t.root
	for cur != t.nul {
		switch {
		case key < cur.key:
			cur = cur.left
		case key > cur.key:
			cur = cur.right
		default:
			return true
		}
	}
	return false
}

// InOrder returns all keys in ascending sorted order.
func (t *RBTree[T]) InOrder() []T {
	var out []T
	var walk func(n *node[T])
	walk = func(n *node[T]) {
		if n == t.nul {
			return
		}
		walk(n.left)
		out = append(out, n.key)
		walk(n.right)
	}
	walk(t.root)
	return out
}

// blackHeight returns the number of black nodes on any root-to-nil path,
// or -1 if the black-height invariant is violated. Used for validation.
func (t *RBTree[T]) blackHeight() int {
	var check func(n *node[T]) int
	check = func(n *node[T]) int {
		if n == t.nul {
			return 1
		}
		if n.color == red && (n.left.color == red || n.right.color == red) {
			return -1 // two reds in a row
		}
		lh := check(n.left)
		rh := check(n.right)
		if lh == -1 || rh == -1 || lh != rh {
			return -1
		}
		if n.color == black {
			return lh + 1
		}
		return lh
	}
	return check(t.root)
}
```

## usage / test

```go
package rbtree

import (
	"math/rand"
	"sort"
	"testing"
)

func TestRBTreeInvariants(t *testing.T) {
	tree := New[int]()
	seen := map[int]bool{}
	var want []int
	for i := 0; i < 2000; i++ {
		k := rand.Intn(5000)
		tree.Insert(k)
		if !seen[k] {
			seen[k] = true
			want = append(want, k)
		}
	}
	sort.Ints(want)

	// In-order traversal must be sorted with no duplicates.
	got := tree.InOrder()
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("at %d got %d want %d", i, got[i], want[i])
		}
	}
	if tree.root.color != black {
		t.Error("root is not black")
	}
	// Black-height invariant holds (and no red-red violations) => balanced.
	if bh := tree.blackHeight(); bh < 0 {
		t.Fatal("red-black invariants violated")
	}
}
```
