# Topological Sort (Kahn's Algorithm and DFS)

A topological sort of a directed acyclic graph (DAG) is a linear ordering of its vertices such that for every directed edge u→v, u appears before v. Use it to order tasks with dependencies: build systems, course prerequisites, package/install ordering, spreadsheet recalculation, and job scheduling. Two standard methods exist. Kahn's algorithm repeatedly removes vertices with in-degree zero (no remaining dependencies), pushing each onto the order and decrementing its neighbors' in-degrees; if not all vertices come out, the graph has a cycle and no topological order exists. The DFS method appends each vertex after visiting all its descendants, then reverses the result. Both run in O(V + E) time and O(V) space. Keywords: topological sort topological ordering toposort dag directed acyclic graph dependency order kahn in-degree indegree zero in degree queue dfs postorder reverse cycle detection prerequisite scheduling build order ordering linearization.

## implementation

```go
package graph

// TopoSortKahn returns a topological ordering of an n-vertex DAG (vertices
// labelled 0..n-1) using Kahn's algorithm. edges[i] = {u, v} is a directed
// edge u->v (v depends on u). The second return value is false if the graph
// has a cycle, in which case no topological order exists.
func TopoSortKahn(n int, edges [][2]int) ([]int, bool) {
	adj := make([][]int, n)
	inDegree := make([]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		inDegree[v]++
	}

	// Start with every vertex that has no incoming edges.
	queue := make([]int, 0, n)
	for v := 0; v < n; v++ {
		if inDegree[v] == 0 {
			queue = append(queue, v)
		}
	}

	order := make([]int, 0, n)
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, u)
		for _, v := range adj[u] {
			inDegree[v]--
			if inDegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	if len(order) != n {
		return nil, false // cycle: some vertices never reached in-degree 0
	}
	return order, true
}

// TopoSortDFS returns a topological ordering using depth-first search.
// It reports false if a cycle is detected.
func TopoSortDFS(n int, edges [][2]int) ([]int, bool) {
	adj := make([][]int, n)
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
	}

	const (
		white = 0 // unvisited
		gray  = 1 // on the current DFS path
		black = 2 // fully processed
	)
	state := make([]int, n)
	order := make([]int, 0, n)
	cyclic := false

	var visit func(u int)
	visit = func(u int) {
		state[u] = gray
		for _, v := range adj[u] {
			if state[v] == gray {
				cyclic = true // back edge to a vertex on the path
				return
			}
			if state[v] == white {
				visit(v)
			}
		}
		state[u] = black
		order = append(order, u) // finished: prepend by reversing later
	}

	for v := 0; v < n; v++ {
		if state[v] == white {
			visit(v)
			if cyclic {
				return nil, false
			}
		}
	}
	// order holds vertices in reverse topological order; reverse it.
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}
	return order, true
}
```

## usage / test

```go
package graph

import "testing"

// validTopo checks that every edge u->v has u positioned before v.
func validTopo(order []int, edges [][2]int) bool {
	pos := make(map[int]int, len(order))
	for i, v := range order {
		pos[v] = i
	}
	for _, e := range edges {
		if pos[e[0]] >= pos[e[1]] {
			return false
		}
	}
	return true
}

func TestTopoSort(t *testing.T) {
	// 0 -> 1, 0 -> 2, 1 -> 3, 2 -> 3, 3 -> 4
	edges := [][2]int{{0, 1}, {0, 2}, {1, 3}, {2, 3}, {3, 4}}

	for _, fn := range []struct {
		name string
		f    func(int, [][2]int) ([]int, bool)
	}{
		{"kahn", TopoSortKahn},
		{"dfs", TopoSortDFS},
	} {
		order, ok := fn.f(5, edges)
		if !ok {
			t.Fatalf("%s: expected a valid DAG order", fn.name)
		}
		if len(order) != 5 || !validTopo(order, edges) {
			t.Fatalf("%s: invalid topological order %v", fn.name, order)
		}
	}

	// Add a cycle 4 -> 0; both algorithms must report failure.
	cyclic := append(edges, [2]int{4, 0})
	if _, ok := TopoSortKahn(5, cyclic); ok {
		t.Error("kahn: should detect cycle")
	}
	if _, ok := TopoSortDFS(5, cyclic); ok {
		t.Error("dfs: should detect cycle")
	}
}
```
