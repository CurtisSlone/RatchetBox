# Slices and maps idioms

Working with slices and maps (Effective Go; Code Review Comments). Slices and maps are references to an
underlying structure.

- Declare an empty slice as `var t []T` (a nil slice), not `t := []T{}`; nil appends and ranges fine.
- Grow with `append`; reassign the result (`s = append(s, v)`). Preallocate with `make([]T, 0, n)`.
- Read a map with the comma-ok idiom to distinguish "absent" from the zero value.
- `make` initializes slices, maps, and channels; `new(T)` returns `*T` to zeroed storage (a `new([]int)`
  is a pointer to a nil slice, almost never what you want - use `make`).
- A nil map is read-only; writing to it panics. Initialize before writing.

```go
var names []string                 // nil slice; preferred over []string{}
names = append(names, "ann", "joe") // reassign the append result
nums := make([]int, 0, 100)        // preallocated, length 0 capacity 100

m := map[string]int{}              // ready to write (make(map[string]int) is equivalent)
if v, ok := m["k"]; ok {           // comma-ok: ok distinguishes missing from zero
	use(v)
}
delete(m, "k")

// new vs make:
p := new([]int)                    // *[]int, *p == nil  (rarely useful)
s := make([]int, 100)              // []int with 100 zeroed elements (idiomatic)
```
