# Binary Tree and Traversals (Inorder, Preorder, Postorder, Level-Order)

A binary tree is a hierarchical structure where each node has at most two children, conventionally called left and right. Use it to model hierarchies, expression trees, and as the basis for search trees and heaps. The four standard traversals visit every node exactly once: preorder (node, left, right), inorder (left, node, right), postorder (left, right, node), and level-order / breadth-first (top to bottom, left to right using a queue). All four run in O(n) time; recursive traversals use O(h) stack space where h is the tree height (O(n) worst case for a skewed tree, O(log n) when balanced), and level-order uses O(w) queue space where w is the maximum width. Keywords: binary tree tree node left right child children traversal traverse walk visit inorder in-order preorder pre-order postorder post-order depth-first dfs level-order levelorder breadth-first bfs queue recursion stack height depth root leaf expression-tree hierarchy.

## implementation

```go
package binarytree

// TreeNode is a node in a binary tree holding an ordered value.
type TreeNode[T any] struct {
	Val   T
	Left  *TreeNode[T]
	Right *TreeNode[T]
}

// Insert adds a child manually; helper for building trees in tests.
func NewNode[T any](v T) *TreeNode[T] { return &TreeNode[T]{Val: v} }

// Preorder visits node, then left subtree, then right subtree.
func Preorder[T any](root *TreeNode[T]) []T {
	var out []T
	var walk func(n *TreeNode[T])
	walk = func(n *TreeNode[T]) {
		if n == nil {
			return
		}
		out = append(out, n.Val)
		walk(n.Left)
		walk(n.Right)
	}
	walk(root)
	return out
}

// Inorder visits left subtree, then node, then right subtree.
// For a binary search tree this yields values in sorted order.
func Inorder[T any](root *TreeNode[T]) []T {
	var out []T
	var walk func(n *TreeNode[T])
	walk = func(n *TreeNode[T]) {
		if n == nil {
			return
		}
		walk(n.Left)
		out = append(out, n.Val)
		walk(n.Right)
	}
	walk(root)
	return out
}

// Postorder visits left subtree, then right subtree, then node.
func Postorder[T any](root *TreeNode[T]) []T {
	var out []T
	var walk func(n *TreeNode[T])
	walk = func(n *TreeNode[T]) {
		if n == nil {
			return
		}
		walk(n.Left)
		walk(n.Right)
		out = append(out, n.Val)
	}
	walk(root)
	return out
}

// LevelOrder visits nodes breadth-first, top to bottom, left to right,
// using a FIFO queue.
func LevelOrder[T any](root *TreeNode[T]) []T {
	if root == nil {
		return nil
	}
	var out []T
	queue := []*TreeNode[T]{root}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		out = append(out, n.Val)
		if n.Left != nil {
			queue = append(queue, n.Left)
		}
		if n.Right != nil {
			queue = append(queue, n.Right)
		}
	}
	return out
}

// Height returns the number of edges on the longest root-to-leaf path.
// An empty tree has height -1, a single node has height 0.
func Height[T any](root *TreeNode[T]) int {
	if root == nil {
		return -1
	}
	l := Height(root.Left)
	r := Height(root.Right)
	if l > r {
		return l + 1
	}
	return r + 1
}
```

## usage / test

```go
package binarytree

import (
	"reflect"
	"testing"
)

// buildTree constructs:
//          4
//        /   \
//       2     6
//      / \   / \
//     1   3 5   7
func buildTree() *TreeNode[int] {
	n := func(v int, l, r *TreeNode[int]) *TreeNode[int] {
		return &TreeNode[int]{Val: v, Left: l, Right: r}
	}
	return n(4,
		n(2, n(1, nil, nil), n(3, nil, nil)),
		n(6, n(5, nil, nil), n(7, nil, nil)),
	)
}

func TestTraversals(t *testing.T) {
	root := buildTree()
	cases := []struct {
		name string
		got  []int
		want []int
	}{
		{"preorder", Preorder(root), []int{4, 2, 1, 3, 6, 5, 7}},
		{"inorder", Inorder(root), []int{1, 2, 3, 4, 5, 6, 7}}, // sorted for a BST
		{"postorder", Postorder(root), []int{1, 3, 2, 5, 7, 6, 4}},
		{"levelorder", LevelOrder(root), []int{4, 2, 6, 1, 3, 5, 7}},
	}
	for _, c := range cases {
		if !reflect.DeepEqual(c.got, c.want) {
			t.Errorf("%s = %v, want %v", c.name, c.got, c.want)
		}
	}
	if h := Height(root); h != 2 {
		t.Errorf("height = %d, want 2", h)
	}
}
```
