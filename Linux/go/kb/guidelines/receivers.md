# Value vs pointer receivers

Choosing a method receiver type (Code Review Comments). Be consistent: don't mix value and pointer
receivers on the same type.

- Use a POINTER receiver if the method mutates the receiver.
- Use a pointer if the struct is large (cheaper than copying) or contains a `sync.Mutex` or other
  field that must not be copied.
- Use a VALUE receiver for small, immutable types (basic-type wrappers, small structs, `time.Time`).
- If some methods need a pointer, give ALL methods on that type a pointer receiver for consistency.
- When in doubt, use a pointer receiver.

```go
type Counter struct {
	mu sync.Mutex // must not be copied -> pointer receivers
	n  int
}

func (c *Counter) Add(n int) { c.mu.Lock(); c.n += n; c.mu.Unlock() } // mutates -> pointer
func (c *Counter) Value() int { c.mu.Lock(); defer c.mu.Unlock(); return c.n }

type Point struct{ X, Y float64 } // small, immutable -> value receiver is fine
func (p Point) Add(q Point) Point { return Point{p.X + q.X, p.Y + q.Y} }
```
