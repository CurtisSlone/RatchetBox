# Prim's Minimum Spanning Tree (MST)

Prim's algorithm finds a minimum spanning tree of a connected, undirected, weighted graph by growing a single tree outward from a starting vertex. It keeps a min-priority queue of edges crossing from the in-tree set to the rest; at each step it pulls the cheapest crossing edge whose far endpoint is not yet in the tree, adds that vertex, and pushes its newly exposed edges. It finishes when all vertices are included. Use it for the same problems as Kruskal (network design, minimum-cost connections); Prim tends to be the natural choice when the graph is dense or given as an adjacency list and you want to grow the tree from a root. With a binary heap it runs in O(E log V) time and O(V + E) space. Keywords: prim prim's minimum spanning tree mst greedy grow tree priority queue min heap crossing edge cut frontier weighted undirected graph network design dense graph adjacency list root least weight connect vertices container/heap.

## implementation

```go
package graph

import "container/heap"

// PEdge is a weighted undirected edge used by Prim's algorithm.
type PEdge struct {
	From, To, Weight int
}

type primHeap []PEdge

func (h primHeap) Len() int           { return len(h) }
func (h primHeap) Less(i, j int) bool { return h[i].Weight < h[j].Weight }
func (h primHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *primHeap) Push(x any)        { *h = append(*h, x.(PEdge)) }
func (h *primHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// PrimMST returns the edges of a minimum spanning tree of an n-vertex graph
// and its total weight, growing the tree from `start`. adj[u] lists every
// edge incident to u (both directions present for an undirected graph).
func PrimMST(n int, adj map[int][]PEdge, start int) (mst []PEdge, total int) {
	inTree := make([]bool, n)
	h := &primHeap{}
	heap.Init(h)

	add := func(u int) {
		inTree[u] = true
		for _, e := range adj[u] {
			if !inTree[e.To] {
				heap.Push(h, e)
			}
		}
	}
	add(start)

	for h.Len() > 0 && len(mst) < n-1 {
		e := heap.Pop(h).(PEdge)
		if inTree[e.To] {
			continue // far endpoint already in the tree; skip to avoid a cycle
		}
		mst = append(mst, e)
		total += e.Weight
		add(e.To)
	}
	return mst, total
}
```

## usage / test

```go
package graph

import "testing"

// buildUndirected turns a list of {u,v,w} into an adjacency list with both
// directions present.
func buildUndirected(triples [][3]int) map[int][]PEdge {
	adj := map[int][]PEdge{}
	for _, t := range triples {
		u, v, w := t[0], t[1], t[2]
		adj[u] = append(adj[u], PEdge{From: u, To: v, Weight: w})
		adj[v] = append(adj[v], PEdge{From: v, To: u, Weight: w})
	}
	return adj
}

func TestPrim(t *testing.T) {
	adj := buildUndirected([][3]int{
		{0, 1, 2},
		{0, 3, 6},
		{1, 2, 3},
		{1, 3, 8},
		{1, 4, 5},
		{2, 4, 7},
		{3, 4, 9},
	})
	mst, total := PrimMST(5, adj, 0)

	// 5 vertices -> 4 edges, same minimum weight as Kruskal: 16.
	if len(mst) != 4 {
		t.Fatalf("MST has %d edges, want 4", len(mst))
	}
	if total != 16 {
		t.Fatalf("MST weight = %d, want 16", total)
	}

	// Every vertex must be incident to at least one MST edge (all connected).
	seen := map[int]bool{0: true}
	for _, e := range mst {
		seen[e.From] = true
		seen[e.To] = true
	}
	for v := 0; v < 5; v++ {
		if !seen[v] {
			t.Fatalf("vertex %d missing from MST", v)
		}
	}
}
```
