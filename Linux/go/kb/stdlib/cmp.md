# cmp (Go standard library)

Package cmp provides types and functions related to comparing ordered values.

Import path: cmp   Toolchain: go1.26.4

package cmp // import "cmp"

Package cmp provides types and functions related to comparing ordered values.

FUNCTIONS

func Compare[T Ordered](x, y T) int
    Compare returns

        -1 if x is less than y,
         0 if x equals y,
        +1 if x is greater than y.

    For floating-point types, a NaN is considered less than any non-NaN,
    a NaN is considered equal to a NaN, and -0.0 is equal to 0.0.

func Less[T Ordered](x, y T) bool
    Less reports whether x is less than y. For floating-point types, a NaN is
    considered less than any non-NaN, and -0.0 is not less than (is equal to)
    0.0.

func Or[T comparable](vals ...T) T
    Or returns the first of its arguments that is not equal to the zero value.
    If no argument is non-zero, it returns the zero value.


TYPES

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}
    Ordered is a constraint that permits any ordered type: any type that
    supports the operators < <= >= >. If future releases of Go add new ordered
    types, this constraint will be modified to include them.

    Note that floating-point types may contain NaN ("not-a-number") values. An
    operator such as == or < will always report false when comparing a NaN value
    with any other value, NaN or not. See the Compare function for a consistent
    way to compare NaN values.

## idiomatic usage

Use cmp.Compare/cmp.Less for ordering of ordered types (handy as slices.SortFunc comparators), and cmp.Or to pick the first non-zero value, which chains tie-breaking comparisons or supplies defaults. Keywords: cmp Compare Less Or three-way comparison ordered generics default fallback first non-zero tie-break sort comparator SortFunc multi-key sort NaN.

```go
import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

// Compare returns -1, 0, or +1; Less returns a bool.
func basics() {
	fmt.Println(cmp.Compare(1, 2)) // -1
	fmt.Println(cmp.Less("a", "aa")) // true
}

// cmp.Or returns its first non-zero argument (e.g. defaults).
func defaults(userInput string) string {
	return cmp.Or(userInput, "default")
}

// cmp.Or chains comparisons for multi-key sorting.
func sortOrders(orders []Order) {
	slices.SortFunc(orders, func(a, b Order) int {
		return cmp.Or(
			strings.Compare(a.Customer, b.Customer),
			strings.Compare(a.Product, b.Product),
			cmp.Compare(b.Price, a.Price), // higher price first
		)
	})
}

type Order struct {
	Product, Customer string
	Price             float64
}
```
