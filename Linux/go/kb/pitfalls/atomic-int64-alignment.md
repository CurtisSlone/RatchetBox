# Pitfall: 64-bit atomic field alignment

`sync/atomic`'s 64-bit ops (`AddInt64`, `LoadInt64`, ...) require their operand to be 64-bit aligned.
The language only guarantees alignment for the FIRST word of a struct/allocated value, so an `int64`
that is not the first field can be MISaligned on 32-bit platforms (GOARCH=386, arm) and the atomic op
panics ("unaligned 64-bit atomic operation") at runtime. Builds clean on amd64; only bites on 32-bit.

- Best fix (Go 1.19+): use the typed atomics `atomic.Int64` / `atomic.Uint64`. They carry their own
  alignment guarantee and a clearer API - no `&` and no field-ordering rules.
- If you must use a raw `int64` with the `atomic.AddInt64` functions, make it the FIRST field of the
  struct (or otherwise ensure 8-byte alignment).

```go
// RISKY - processed is not the first field; can panic on 386/arm
type Dispatcher struct {
	queue     chan int
	wg        sync.WaitGroup
	processed int64 // misaligned on 32-bit
}

func (d *Dispatcher) inc() { atomic.AddInt64(&d.processed, 1) }

// PREFERRED - typed atomic, always safe regardless of field order (Go 1.19+)
type Dispatcher struct {
	queue     chan int
	wg        sync.WaitGroup
	processed atomic.Int64
}

func (d *Dispatcher) inc() { d.processed.Add(1) }
func (d *Dispatcher) count() int64 { return d.processed.Load() }
```
