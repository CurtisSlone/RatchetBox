# Kruskal's Minimum Spanning Tree (MST)

Kruskal's algorithm finds a minimum spanning tree of a connected, undirected, weighted graph: a subset of edges that connects all vertices with the smallest possible total weight and no cycles. It sorts every edge by weight ascending, then greedily adds each edge if its two endpoints are not already connected, using a union-find (disjoint-set) structure to detect cycles in near-constant time. After processing it has chosen V-1 edges (for a connected graph). Use it for network design, clustering, and laying out minimum-cost connections; Kruskal is especially natural when edges are given as a flat list and the graph is sparse. Sorting dominates the cost: O(E log E) time (equivalently O(E log V)) and O(V + E) space. Keywords: kruskal minimum spanning tree mst greedy edges sort union find disjoint set cycle detection connect all vertices minimum cost weighted undirected graph network design clustering spanning forest least weight.

## implementation

```go
package graph

import "sort"

// WEdge is a weighted undirected edge between U and V.
type WEdge struct {
	U, V, Weight int
}

// dsu is a union-find with path compression and union by rank.
type dsu struct {
	parent, rank []int
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return &dsu{parent: p, rank: make([]int, n)}
}

func (d *dsu) find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *dsu) union(a, b int) bool {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return false
	}
	if d.rank[ra] < d.rank[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	if d.rank[ra] == d.rank[rb] {
		d.rank[ra]++
	}
	return true
}

// KruskalMST returns the edges of a minimum spanning tree of an n-vertex
// graph and the MST's total weight. Vertices are labelled 0..n-1.
func KruskalMST(n int, edges []WEdge) (mst []WEdge, total int) {
	sorted := make([]WEdge, len(edges))
	copy(sorted, edges)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Weight < sorted[j].Weight
	})

	d := newDSU(n)
	for _, e := range sorted {
		// Add the edge only if it joins two distinct components (no cycle).
		if d.union(e.U, e.V) {
			mst = append(mst, e)
			total += e.Weight
			if len(mst) == n-1 {
				break // spanning tree complete
			}
		}
	}
	return mst, total
}
```

## usage / test

```go
package graph

import "testing"

func TestKruskal(t *testing.T) {
	// Undirected weighted graph, 5 vertices.
	edges := []WEdge{
		{0, 1, 2},
		{0, 3, 6},
		{1, 2, 3},
		{1, 3, 8},
		{1, 4, 5},
		{2, 4, 7},
		{3, 4, 9},
	}
	mst, total := KruskalMST(5, edges)

	// A spanning tree of 5 vertices has exactly 4 edges.
	if len(mst) != 4 {
		t.Fatalf("MST has %d edges, want 4", len(mst))
	}
	// Known minimum total weight: 2 + 3 + 5 + 6 = 16.
	if total != 16 {
		t.Fatalf("MST weight = %d, want 16", total)
	}

	// The selected edges must connect all vertices (single component).
	d := newDSU(5)
	for _, e := range mst {
		d.union(e.U, e.V)
	}
	root := d.find(0)
	for v := 1; v < 5; v++ {
		if d.find(v) != root {
			t.Fatalf("vertex %d not connected in MST", v)
		}
	}
}
```
