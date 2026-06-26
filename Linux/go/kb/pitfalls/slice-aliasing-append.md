# Pitfall: slice aliasing and append surprises

A slice is a view (pointer, len, cap) over a backing array. Slices that share a backing array alias
each other's data, and `append` mutates the backing array in place when there is spare capacity. Builds
clean; manifests as data corruption under `go test`.

- `append` may modify a shared backing array (within cap) or allocate a new one (beyond cap) - do not
  assume either. Always use the returned slice: `s = append(s, v)`.
- A sub-slice shares storage with its parent; writing through one is visible through the other. Copy
  (`copy`, or a full three-index slice) when you need independence.
- Re-slicing keeps the original capacity, so a later append can overwrite "dropped" elements.

```go
// Aliasing: b shares a's backing array
a := []int{1, 2, 3, 4}
b := a[1:3]   // b == [2 3], shares storage with a
b[0] = 99     // a is now [1 99 3 4]

// append may overwrite shared storage when cap allows:
base := make([]int, 3, 5) // len 3, cap 5
x := append(base, 1)      // writes into base's backing array (within cap)
y := append(base, 2)      // overwrites the same slot -> x[3] is now 2, not 1

// Copy for an independent slice:
indep := make([]int, len(a))
copy(indep, a)
```
