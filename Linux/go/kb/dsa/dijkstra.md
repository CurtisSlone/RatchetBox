# Dijkstra's Shortest Path Algorithm

Dijkstra's algorithm finds the shortest paths from a single source vertex to all other vertices in a graph with non-negative edge weights. It grows a set of finalized vertices: repeatedly it picks the unfinalized vertex with the smallest tentative distance (using a min-priority queue / min-heap), finalizes it, and relaxes its outgoing edges (updating a neighbor's distance if going through the current vertex is shorter). Use it for road networks, network routing, and any weighted shortest-path problem where weights are non-negative (for negative weights use Bellman-Ford). With a binary heap it runs in O((V + E) log V) time and O(V) space. Keywords: dijkstra shortest path single source weighted graph non-negative edge weights priority queue min heap relaxation relax distance tentative finalize visited greedy routing road network container/heap path cost spt shortest path tree.

## implementation

```go
package graph

import "container/heap"

// Edge is a weighted directed edge to vertex To.
type Edge struct {
	To     int
	Weight int
}

// item is a vertex with its current best distance, for the priority queue.
type item struct {
	vertex int
	dist   int
}

type minHeap []item

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].dist < h[j].dist }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x any)         { *h = append(*h, x.(item)) }
func (h *minHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Dijkstra computes the shortest distance from src to every vertex in an
// n-vertex graph with non-negative weights. adj[u] lists u's outgoing edges.
// Unreachable vertices get distance -1.
func Dijkstra(n int, adj map[int][]Edge, src int) []int {
	const inf = int(^uint(0) >> 1)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = inf
	}
	dist[src] = 0

	pq := &minHeap{{vertex: src, dist: 0}}
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(item)
		if cur.dist > dist[cur.vertex] {
			continue // stale entry, a better distance was already found
		}
		for _, e := range adj[cur.vertex] {
			if nd := cur.dist + e.Weight; nd < dist[e.To] {
				dist[e.To] = nd
				heap.Push(pq, item{vertex: e.To, dist: nd})
			}
		}
	}

	for i := range dist {
		if dist[i] == inf {
			dist[i] = -1
		}
	}
	return dist
}
```

## usage / test

```go
package graph

import (
	"reflect"
	"testing"
)

func TestDijkstra(t *testing.T) {
	// Directed weighted graph (6 vertices):
	// 0->1 (7), 0->2 (9), 0->5 (14), 1->2 (10), 1->3 (15),
	// 2->3 (11), 2->5 (2), 3->4 (6), 5->4 (9)
	adj := map[int][]Edge{
		0: {{1, 7}, {2, 9}, {5, 14}},
		1: {{2, 10}, {3, 15}},
		2: {{3, 11}, {5, 2}},
		3: {{4, 6}},
		5: {{4, 9}},
	}
	got := Dijkstra(6, adj, 0)
	// Known shortest distances from vertex 0.
	want := []int{0, 7, 9, 20, 20, 11}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Dijkstra = %v, want %v", got, want)
	}

	// Unreachable vertex reports -1.
	adj2 := map[int][]Edge{0: {{1, 5}}}
	d := Dijkstra(3, adj2, 0)
	if d[2] != -1 {
		t.Errorf("unreachable vertex distance = %d, want -1", d[2])
	}
}
```
