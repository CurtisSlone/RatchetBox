# Interface conventions

How idiomatic Go defines and uses interfaces (Effective Go; Code Review Comments; proverbs). The bigger
the interface, the weaker the abstraction; `interface{}` says nothing.

- Define interfaces in the package that USES them (the consumer), not the implementing package.
- Implementations should return concrete types, not interfaces, so callers gain methods/fields and the
  package can add methods without breaking users.
- Keep interfaces small (often one method, named "-er"); compose larger ones by embedding.
- Accept interfaces, return structs.

```go
// Implementor returns a concrete type.
package producer

type Thinger struct{ /* ... */ }

func (t Thinger) Thing() bool { return true }

func NewThinger() Thinger { return Thinger{} }

// Consumer declares the small interface IT needs.
package consumer

type Thinger interface{ Thing() bool }

func Foo(t Thinger) string { /* ... */ return "" }

// Compose interfaces by embedding.
type Reader interface{ Read(p []byte) (n int, err error) }
type Writer interface{ Write(p []byte) (n int, err error) }

type ReadWriter interface {
	Reader
	Writer
}

// Compile-time check that *RawMessage implements json.Marshaler.
var _ json.Marshaler = (*RawMessage)(nil)
```
