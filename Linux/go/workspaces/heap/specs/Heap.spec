name: Heap
role: component
intent: Implements a min-heap of integers using the standard library's heap package
api:
  - func New() *Heap
  - func (h *Heap) Push(x int)
  - func (h *Heap) Pop() (int, bool)
  - func (h *Heap) Len() int
behavior:
  - The heap maintains the min-heap property where the smallest element is at the root
  - Push adds elements to the heap while maintaining the heap property
  - Pop removes and returns the minimum element, returning false if the heap is empty
  - Len returns the number of elements in the heap
constraints: Uses container/heap package, standard library only, no external dependencies
