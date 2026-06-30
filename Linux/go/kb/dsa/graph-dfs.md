# Depth-First Search (DFS) on a Graph

Depth-first search explores a graph by going as deep as possible along each branch before backtracking, using either recursion (the call stack) or an explicit stack, plus a visited set to avoid cycles. Use it for reachability, enumerating connected components, detecting cycles, generating topological orders, finding bridges/articulation points, path-finding when any path will do, and exhaustive search/backtracking. Unlike BFS it does not find shortest paths in general, but it uses less memory on wide graphs. Time is O(V + E) and space is O(V) for the visited set plus recursion/stack depth. Keywords: dfs depth first search graph traversal recursion stack backtracking visit explore neighbors adjacency reachability connected components cycle detection preorder postorder path finding deep branch backtrack.

## implementation

```go
package graph

// DFS returns the order vertices are first visited starting from `start`,
// exploring depth-first using recursion.
func DFS(adj map[int][]int, start int) []int {
	visited := map[int]bool{}
	var order []int
	var visit func(u int)
	visit = func(u int) {
		visited[u] = true
		order = append(order, u)
		for _, v := range adj[u] {
			if !visited[v] {
				visit(v)
			}
		}
	}
	visit(start)
	return order
}

// DFSIterative returns the visit order using an explicit stack instead of
// recursion (avoids stack overflow on very deep graphs).
func DFSIterative(adj map[int][]int, start int) []int {
	visited := map[int]bool{}
	var order []int
	stack := []int{start}
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if visited[u] {
			continue
		}
		visited[u] = true
		order = append(order, u)
		// Push neighbors; they will be popped in reverse insertion order.
		for _, v := range adj[u] {
			if !visited[v] {
				stack = append(stack, v)
			}
		}
	}
	return order
}

// ConnectedComponents returns the groups of mutually reachable vertices in an
// undirected graph. `vertices` lists every vertex (so isolated ones count).
func ConnectedComponents(adj map[int][]int, vertices []int) [][]int {
	visited := map[int]bool{}
	var comps [][]int
	var visit func(u int, comp *[]int)
	visit = func(u int, comp *[]int) {
		visited[u] = true
		*comp = append(*comp, u)
		for _, v := range adj[u] {
			if !visited[v] {
				visit(v, comp)
			}
		}
	}
	for _, s := range vertices {
		if !visited[s] {
			var comp []int
			visit(s, &comp)
			comps = append(comps, comp)
		}
	}
	return comps
}
```

## usage / test

```go
package graph

import (
	"sort"
	"testing"
)

func TestDFS(t *testing.T) {
	adj := map[int][]int{
		0: {1, 2},
		1: {0, 3},
		2: {0},
		3: {1},
		// component 2: 4-5
		4: {5},
		5: {4},
	}

	// DFS from 0 reaches exactly its component {0,1,2,3}.
	order := DFS(adj, 0)
	if len(order) != 4 || order[0] != 0 {
		t.Fatalf("DFS(0) = %v, want 4 vertices starting at 0", order)
	}

	// Iterative DFS reaches the same set of vertices.
	it := DFSIterative(adj, 0)
	sort.Ints(it)
	if len(it) != 4 {
		t.Fatalf("DFSIterative(0) reached %v, want 4 vertices", it)
	}

	// Two connected components overall.
	comps := ConnectedComponents(adj, []int{0, 1, 2, 3, 4, 5})
	if len(comps) != 2 {
		t.Fatalf("got %d components, want 2: %v", len(comps), comps)
	}
}
```
