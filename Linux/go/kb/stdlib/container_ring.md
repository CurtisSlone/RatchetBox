# container/ring (Go standard library)

Package ring implements operations on circular lists.

Import path: container/ring   Toolchain: go1.26.4

package ring // import "container/ring"

Package ring implements operations on circular lists.

TYPES

type Ring struct {
	Value any // for use by client; untouched by this library
	// Has unexported fields.
}
    A Ring is an element of a circular list, or ring. Rings do not have a
    beginning or end; a pointer to any ring element serves as reference to the
    entire ring. Empty rings are represented as nil Ring pointers. The zero
    value for a Ring is a one-element ring with a nil Value.

func New(n int) *Ring
    New creates a ring of n elements.

func (r *Ring) Do(f func(any))
    Do calls function f on each element of the ring, in forward order. The
    behavior of Do is undefined if f changes *r.

func (r *Ring) Len() int
    Len computes the number of elements in ring r. It executes in time
    proportional to the number of elements.

func (r *Ring) Link(s *Ring) *Ring
    Link connects ring r with ring s such that r.Next() becomes s and returns
    the original value for r.Next(). r must not be empty.

    If r and s point to the same ring, linking them removes the elements between
    r and s from the ring. The removed elements form a subring and the result
    is a reference to that subring (if no elements were removed, the result is
    still the original value for r.Next(), and not nil).

    If r and s point to different rings, linking them creates a single ring
    with the elements of s inserted after r. The result points to the element
    following the last element of s after insertion.

func (r *Ring) Move(n int) *Ring
    Move moves n % r.Len() elements backward (n < 0) or forward (n >= 0) in the
    ring and returns that ring element. r must not be empty.

func (r *Ring) Next() *Ring
    Next returns the next ring element. r must not be empty.

func (r *Ring) Prev() *Ring
    Prev returns the previous ring element. r must not be empty.

func (r *Ring) Unlink(n int) *Ring
    Unlink removes n % r.Len() elements from the ring r, starting at r.Next().
    If n % r.Len() == 0, r remains unchanged. The result is the removed subring.
    r must not be empty.

## idiomatic usage

Create a circular list with ring.New(n), set each element's Value while advancing with Next, and traverse the whole ring with Do (or Next/Prev/Move). Keywords: ring.New Len Next Prev Move Do Link Unlink Value circular list ring buffer traverse iterate rotate.

```go
import (
	"container/ring"
	"fmt"
)

r := ring.New(5)
n := r.Len()

// Fill the ring with values.
for i := 0; i < n; i++ {
	r.Value = i
	r = r.Next()
}

// Visit every element.
r.Do(func(p any) {
	fmt.Println(p.(int)) // 0 1 2 3 4
})
```
