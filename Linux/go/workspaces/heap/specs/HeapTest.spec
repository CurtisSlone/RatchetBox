name: HeapTest
role: test
intent: Example tests and a fuzz property test for the heap, in one test file
api:
  - func TestHeap(t *testing.T)
  - func FuzzHeap(f *testing.F)
behavior:
  - TestHeap: basic push/pop example cases; popping yields non-decreasing order; popped multiset equals pushed multiset; edge cases (empty heap, single element)
  - FuzzHeap: fuzz a []byte, derive a sequence of ints (one per byte), push all, pop until empty; assert the popped sequence is non-decreasing and its multiset equals the pushed bytes
constraints: Uses testing package, standard library only; one _test.go file
