# Pitfall: writing to a nil map

A nil map reads fine (returns zero values) but PANICS on write at runtime - `go build` accepts it; the
panic only shows when the code runs (so `go test` catches it, `go build` does not).

- A map's zero value is nil. You must initialize it with `make` or a literal before assigning keys.
- A struct field of map type is nil until you set it; initialize in the constructor.

```go
// WRONG - panic: assignment to entry in nil map
var m map[string]int
m["a"] = 1 // runtime panic

// RIGHT - initialize first
m := make(map[string]int) // or map[string]int{}
m["a"] = 1

// RIGHT - initialize a map field in the constructor
type Counter struct{ counts map[string]int }

func NewCounter() *Counter { return &Counter{counts: make(map[string]int)} }
func (c *Counter) Inc(k string) { c.counts[k]++ } // safe: counts is non-nil
```
