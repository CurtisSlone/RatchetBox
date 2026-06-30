# Breadth-First Search (BFS) on a Graph

Breadth-first search explores a graph level by level: it visits the start vertex, then all its neighbors, then their unvisited neighbors, and so on, using a FIFO queue to manage the frontier and a visited set to avoid revisiting. Use it to find the shortest path in an unweighted graph (fewest edges), to compute distances from a source, to test reachability, and to discover connected components. Because it expands outward in rings, the first time BFS reaches a vertex it has done so by a minimum number of edges. Time is O(V + E) and space is O(V) for the queue and visited set. Keywords: bfs breadth first search graph traversal shortest path unweighted queue fifo frontier level order distance reachability connected components visit explore neighbors adjacency layer wave hops fewest edges.

## implementation

```go
package graph

// BFS returns the order in which vertices are visited starting from `start`,
// exploring the graph breadth-first using a FIFO queue.
// adj is an adjacency list: adj[u] lists u's neighbors.
func BFS(adj map[int][]int, start int) []int {
	visited := map[int]bool{start: true}
	queue := []int{start}
	var order []int
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, u)
		for _, v := range adj[u] {
			if !visited[v] {
				visited[v] = true
				queue = append(queue, v)
			}
		}
	}
	return order
}

// ShortestPath returns the shortest path (fewest edges) from start to target
// in an unweighted graph, or nil if target is unreachable. The path includes
// both endpoints.
func ShortestPath(adj map[int][]int, start, target int) []int {
	if start == target {
		return []int{start}
	}
	visited := map[int]bool{start: true}
	prev := map[int]int{} // prev[v] = vertex we reached v from
	queue := []int{start}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, v := range adj[u] {
			if visited[v] {
				continue
			}
			visited[v] = true
			prev[v] = u
			if v == target {
				return reconstruct(prev, start, target)
			}
			queue = append(queue, v)
		}
	}
	return nil
}

func reconstruct(prev map[int]int, start, target int) []int {
	path := []int{target}
	for path[len(path)-1] != start {
		path = append(path, prev[path[len(path)-1]])
	}
	// reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}
```

## usage / test

```go
package graph

import (
	"reflect"
	"testing"
)

func TestBFSShortestPath(t *testing.T) {
	// Undirected graph:
	// 0-1, 0-2, 1-3, 2-3, 3-4, 4-5
	adj := map[int][]int{
		0: {1, 2},
		1: {0, 3},
		2: {0, 3},
		3: {1, 2, 4},
		4: {3, 5},
		5: {4},
	}

	// BFS from 0 visits closer vertices first.
	order := BFS(adj, 0)
	if order[0] != 0 {
		t.Fatalf("BFS should start at 0, got %v", order)
	}
	if len(order) != 6 {
		t.Fatalf("BFS should visit all 6 vertices, got %v", order)
	}

	// Shortest unweighted path 0 -> 5 has length 4 edges (5 vertices):
	// 0 -> 1 -> 3 -> 4 -> 5  (or via 2). It must have 5 nodes.
	path := ShortestPath(adj, 0, 5)
	if len(path) != 5 || path[0] != 0 || path[len(path)-1] != 5 {
		t.Fatalf("path 0->5 = %v, want length-5 path", path)
	}

	// Unreachable target returns nil.
	adj[6] = []int{} // isolated vertex
	if p := ShortestPath(adj, 0, 6); p != nil {
		t.Fatalf("path to isolated vertex should be nil, got %v", p)
	}

	// Same-vertex path is just that vertex.
	if got := ShortestPath(adj, 3, 3); !reflect.DeepEqual(got, []int{3}) {
		t.Fatalf("self path = %v, want [3]", got)
	}
}
```
