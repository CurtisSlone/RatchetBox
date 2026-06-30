# Union-Find (Disjoint Set Union, DSU) with Path Compression and Union by Rank

Union-find, also called disjoint-set union (DSU), maintains a collection of disjoint sets and answers "are these two elements in the same set?" while merging sets on the fly. Each set is a tree whose root is its representative; `Find` returns the root and `Union` links two roots. Use it for dynamic connectivity, detecting cycles while building a graph, grouping connected components, Kruskal's minimum spanning tree, and equivalence-class problems. With both optimizations — path compression (flattening the tree during `Find`) and union by rank (attaching the shorter tree under the taller) — each operation runs in nearly O(1), specifically O(α(n)) amortized where α is the inverse Ackermann function. Space is O(n). Keywords: union find disjoint set dsu union-find disjoint-set-union connected components dynamic connectivity merge sets representative root find union path compression union by rank inverse ackermann cycle detection equivalence kruskal grouping same set.

## implementation

```go
package unionfind

// UnionFind maintains disjoint sets over elements 0..n-1 with path
// compression and union by rank for near-constant-time operations.
type UnionFind struct {
	parent []int
	rank   []int // upper bound on tree height
	count  int   // number of disjoint sets
}

// New creates n singleton sets {0}, {1}, ..., {n-1}.
func New(n int) *UnionFind {
	parent := make([]int, n)
	rank := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	return &UnionFind{parent: parent, rank: rank, count: n}
}

// Find returns the representative (root) of x's set, compressing the path.
func (u *UnionFind) Find(x int) int {
	for u.parent[x] != x {
		u.parent[x] = u.parent[u.parent[x]] // path halving
		x = u.parent[x]
	}
	return x
}

// Union merges the sets containing a and b. It reports whether a merge
// actually happened (false if they were already in the same set).
func (u *UnionFind) Union(a, b int) bool {
	ra, rb := u.Find(a), u.Find(b)
	if ra == rb {
		return false
	}
	// Attach the smaller-rank tree under the larger-rank root.
	if u.rank[ra] < u.rank[rb] {
		ra, rb = rb, ra
	}
	u.parent[rb] = ra
	if u.rank[ra] == u.rank[rb] {
		u.rank[ra]++
	}
	u.count--
	return true
}

// Connected reports whether a and b are in the same set.
func (u *UnionFind) Connected(a, b int) bool {
	return u.Find(a) == u.Find(b)
}

// Count returns the current number of disjoint sets.
func (u *UnionFind) Count() int { return u.count }
```

## usage / test

```go
package unionfind

import "testing"

func TestUnionFind(t *testing.T) {
	uf := New(10)
	if uf.Count() != 10 {
		t.Fatalf("Count = %d, want 10", uf.Count())
	}

	uf.Union(0, 1)
	uf.Union(1, 2) // {0,1,2}
	uf.Union(3, 4) // {3,4}
	uf.Union(5, 6)
	uf.Union(6, 7) // {5,6,7}

	if !uf.Connected(0, 2) {
		t.Error("0 and 2 should be connected")
	}
	if uf.Connected(2, 3) {
		t.Error("2 and 3 should not be connected")
	}
	// Merging an already-connected pair must be a no-op.
	if uf.Union(0, 2) {
		t.Error("Union of already-connected pair returned true")
	}
	// Sets: {0,1,2} {3,4} {5,6,7} {8} {9} => 5 components.
	if uf.Count() != 5 {
		t.Fatalf("Count = %d, want 5", uf.Count())
	}

	uf.Union(2, 4) // join {0,1,2} and {3,4}
	if !uf.Connected(0, 3) {
		t.Error("0 and 3 should now be connected")
	}
	if uf.Count() != 4 {
		t.Fatalf("Count = %d, want 4", uf.Count())
	}
}
```
