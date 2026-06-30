# Graph Representations (Adjacency List and Adjacency Matrix)

A graph is a set of vertices connected by edges, which may be directed or undirected and optionally weighted. The two standard in-memory representations are the adjacency list — each vertex stores a list (or map) of its neighbors — and the adjacency matrix — a V-by-V grid where cell [u][v] records whether (or how heavily) an edge connects u to v. Use an adjacency list for sparse graphs (few edges): it uses O(V + E) space and lets you iterate a vertex's neighbors in O(degree); it is the default for BFS, DFS, Dijkstra, and most traversal algorithms. Use an adjacency matrix for dense graphs or when you need O(1) "is there an edge u→v?" lookups: it uses O(V^2) space regardless of edge count. Adding an edge is O(1) in both; checking a specific edge is O(degree) in a list but O(1) in a matrix. Keywords: graph representation adjacency list adjacency matrix vertices vertex nodes edges directed undirected weighted unweighted sparse dense neighbors degree add edge has edge incidence matrix container map slice space complexity.

## implementation

```go
package graph

// AdjacencyList represents a graph as a map from each vertex to its
// neighbors and edge weights. Space O(V + E); good for sparse graphs.
type AdjacencyList struct {
	directed bool
	adj      map[int]map[int]int // adj[u][v] = weight of edge u->v
}

// NewAdjacencyList creates an empty graph. directed selects directed vs
// undirected semantics.
func NewAdjacencyList(directed bool) *AdjacencyList {
	return &AdjacencyList{directed: directed, adj: make(map[int]map[int]int)}
}

// AddVertex ensures v exists in the graph (with no edges).
func (g *AdjacencyList) AddVertex(v int) {
	if g.adj[v] == nil {
		g.adj[v] = make(map[int]int)
	}
}

// AddEdge adds an edge u->v with the given weight (use 1 for unweighted).
// For an undirected graph the reverse edge is added too.
func (g *AdjacencyList) AddEdge(u, v, weight int) {
	g.AddVertex(u)
	g.AddVertex(v)
	g.adj[u][v] = weight
	if !g.directed {
		g.adj[v][u] = weight
	}
}

// HasEdge reports whether an edge u->v exists, in O(1) amortized.
func (g *AdjacencyList) HasEdge(u, v int) bool {
	_, ok := g.adj[u][v]
	return ok
}

// Neighbors returns the neighbors of u and their edge weights.
func (g *AdjacencyList) Neighbors(u int) map[int]int { return g.adj[u] }

// AdjacencyMatrix represents a graph as a V-by-V weight grid.
// Space O(V^2); good for dense graphs and O(1) edge lookups.
type AdjacencyMatrix struct {
	directed bool
	n        int
	mat      [][]int  // mat[u][v] = weight, or 0 when no edge
	present  [][]bool // distinguishes a real 0-weight edge from "no edge"
}

// NewAdjacencyMatrix creates an n-vertex graph with no edges.
func NewAdjacencyMatrix(n int, directed bool) *AdjacencyMatrix {
	mat := make([][]int, n)
	present := make([][]bool, n)
	for i := range mat {
		mat[i] = make([]int, n)
		present[i] = make([]bool, n)
	}
	return &AdjacencyMatrix{directed: directed, n: n, mat: mat, present: present}
}

// AddEdge adds an edge u->v with the given weight.
func (g *AdjacencyMatrix) AddEdge(u, v, weight int) {
	g.mat[u][v] = weight
	g.present[u][v] = true
	if !g.directed {
		g.mat[v][u] = weight
		g.present[v][u] = true
	}
}

// HasEdge reports whether an edge u->v exists, in O(1).
func (g *AdjacencyMatrix) HasEdge(u, v int) bool { return g.present[u][v] }

// Weight returns the weight of edge u->v (only meaningful if HasEdge is true).
func (g *AdjacencyMatrix) Weight(u, v int) int { return g.mat[u][v] }
```

## usage / test

```go
package graph

import "testing"

func TestRepresentations(t *testing.T) {
	// Undirected list: edge 0-1 implies both directions.
	list := NewAdjacencyList(false)
	list.AddEdge(0, 1, 5)
	list.AddEdge(1, 2, 3)
	if !list.HasEdge(0, 1) || !list.HasEdge(1, 0) {
		t.Error("undirected edge should exist both ways")
	}
	if list.HasEdge(0, 2) {
		t.Error("0-2 edge should not exist")
	}
	if w := list.Neighbors(1)[2]; w != 3 {
		t.Errorf("weight 1->2 = %d, want 3", w)
	}

	// Directed matrix: edge 0->1 only.
	mat := NewAdjacencyMatrix(3, true)
	mat.AddEdge(0, 1, 7)
	if !mat.HasEdge(0, 1) || mat.HasEdge(1, 0) {
		t.Error("directed edge should be one-way")
	}
	if mat.Weight(0, 1) != 7 {
		t.Errorf("weight 0->1 = %d, want 7", mat.Weight(0, 1))
	}
}
```
