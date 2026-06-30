# sort (Go standard library)

Package sort provides primitives for sorting slices and user-defined

Import path: sort   Toolchain: go1.26.4

package sort // import "sort"

Package sort provides primitives for sorting slices and user-defined
collections.

FUNCTIONS

func Find(n int, cmp func(int) int) (i int, found bool)
    Find uses binary search to find and return the smallest index i in [0,
    n) at which cmp(i) <= 0. If there is no such index i, Find returns i = n.
    The found result is true if i < n and cmp(i) == 0. Find calls cmp(i) only
    for i in the range [0, n).

    To permit binary search, Find requires that cmp(i) > 0 for a leading prefix
    of the range, cmp(i) == 0 in the middle, and cmp(i) < 0 for the final suffix
    of the range. (Each subrange could be empty.) The usual way to establish
    this condition is to interpret cmp(i) as a comparison of a desired target
    value t against entry i in an underlying indexed data structure x, returning
    <0, 0, and >0 when t < x[i], t == x[i], and t > x[i], respectively.

    For example, to look for a particular string in a sorted, random-access list
    of strings:

        i, found := sort.Find(x.Len(), func(i int) int {
            return strings.Compare(target, x.At(i))
        })
        if found {
            fmt.Printf("found %s at entry %d\n", target, i)
        } else {
            fmt.Printf("%s not found, would insert at %d", target, i)
        }

func Float64s(x []float64)
    Float64s sorts a slice of float64s in increasing order. Not-a-number (NaN)
    values are ordered before other values.

    Note: as of Go 1.22, this function simply calls slices.Sort.

func Float64sAreSorted(x []float64) bool
    Float64sAreSorted reports whether the slice x is sorted in increasing order,
    with not-a-number (NaN) values before any other values.

    Note: as of Go 1.22, this function simply calls slices.IsSorted.

func Ints(x []int)
    Ints sorts a slice of ints in increasing order.

    Note: as of Go 1.22, this function simply calls slices.Sort.

func IntsAreSorted(x []int) bool
    IntsAreSorted reports whether the slice x is sorted in increasing order.

    Note: as of Go 1.22, this function simply calls slices.IsSorted.

func IsSorted(data Interface) bool
    IsSorted reports whether data is sorted.

    Note: in many situations, the newer slices.IsSortedFunc function is more
    ergonomic and runs faster.

func Search(n int, f func(int) bool) int
    Search uses binary search to find and return the smallest index i in [0,
    n) at which f(i) is true, assuming that on the range [0, n), f(i) == true
    implies f(i+1) == true. That is, Search requires that f is false for some
    (possibly empty) prefix of the input range [0, n) and then true for the
    (possibly empty) remainder; Search returns the first true index. If there is
    no such index, Search returns n. (Note that the "not found" return value is
    not -1 as in, for instance, strings.Index.) Search calls f(i) only for i in
    the range [0, n).

    A common use of Search is to find the index i for a value x in a sorted,
    indexable data structure such as an array or slice. In this case,
    the argument f, typically a closure, captures the value to be searched for,
    and how the data structure is indexed and ordered.

    For instance, given a slice data sorted in ascending order, the call
    Search(len(data), func(i int) bool { return data[i] >= 23 }) returns the
    smallest index i such that data[i] >= 23. If the caller wants to find
    whether 23 is in the slice, it must test data[i] == 23 separately.

    Searching data sorted in descending order would use the <= operator instead
    of the >= operator.

    To complete the example above, the following code tries to find the value x
    in an integer slice data sorted in ascending order:

        x := 23
        i := sort.Search(len(data), func(i int) bool { return data[i] >= x })
        if i < len(data) && data[i] == x {
        	// x is present at data[i]
        } else {
        	// x is not present in data,
        	// but i is the index where it would be inserted.
        }

    As a more whimsical example, this program guesses your number:

        func GuessingGame() {
        	var s string
        	fmt.Printf("Pick an integer from 0 to 100.\n")
        	answer := sort.Search(100, func(i int) bool {
        		fmt.Printf("Is your number <= %d? ", i)
        		fmt.Scanf("%s", &s)
        		return s != "" && s[0] == 'y'
        	})
        	fmt.Printf("Your number is %d.\n", answer)
        }

func SearchFloat64s(a []float64, x float64) int
    SearchFloat64s searches for x in a sorted slice of float64s and returns the
    index as specified by Search. The return value is the index to insert x if x
    is not present (it could be len(a)). The slice must be sorted in ascending
    order.

func SearchInts(a []int, x int) int
    SearchInts searches for x in a sorted slice of ints and returns the index as
    specified by Search. The return value is the index to insert x if x is not
    present (it could be len(a)). The slice must be sorted in ascending order.

func SearchStrings(a []string, x string) int
    SearchStrings searches for x in a sorted slice of strings and returns the
    index as specified by Search. The return value is the index to insert x if x
    is not present (it could be len(a)). The slice must be sorted in ascending
    order.

func Slice(x any, less func(i, j int) bool)
    Slice sorts the slice x given the provided less function. It panics if x is
    not a slice.

    The sort is not guaranteed to be stable: equal elements may be reversed from
    their original order. For a stable sort, use SliceStable.

    The less function must satisfy the same requirements as the Interface type's
    Less method.

    Note: in many situations, the newer slices.SortFunc function is more
    ergonomic and runs faster.

func SliceIsSorted(x any, less func(i, j int) bool) bool
    SliceIsSorted reports whether the slice x is sorted according to the
    provided less function. It panics if x is not a slice.

    Note: in many situations, the newer slices.IsSortedFunc function is more
    ergonomic and runs faster.

func SliceStable(x any, less func(i, j int) bool)
    SliceStable sorts the slice x using the provided less function, keeping
    equal elements in their original order. It panics if x is not a slice.

    The less function must satisfy the same requirements as the Interface type's
    Less method.

    Note: in many situations, the newer slices.SortStableFunc function is more
    ergonomic and runs faster.

func Sort(data Interface)
    Sort sorts data in ascending order as determined by the Less method. It
    makes one call to data.Len to determine n and O(n*log(n)) calls to data.Less
    and data.Swap. The sort is not guaranteed to be stable.

    Note: in many situations, the newer slices.SortFunc function is more
    ergonomic and runs faster.

func Stable(data Interface)
    Stable sorts data in ascending order as determined by the Less method,
    while keeping the original order of equal elements.

    It makes one call to data.Len to determine n, O(n*log(n)) calls to data.Less
    and O(n*log(n)*log(n)) calls to data.Swap.

    Note: in many situations, the newer slices.SortStableFunc function is more
    ergonomic and runs faster.

func Strings(x []string)
    Strings sorts a slice of strings in increasing order.

    Note: as of Go 1.22, this function simply calls slices.Sort.

func StringsAreSorted(x []string) bool
    StringsAreSorted reports whether the slice x is sorted in increasing order.

    Note: as of Go 1.22, this function simply calls slices.IsSorted.


TYPES

type Float64Slice []float64
    Float64Slice implements Interface for a []float64, sorting in increasing
    order, with not-a-number (NaN) values ordered before other values.

func (x Float64Slice) Len() int

func (x Float64Slice) Less(i, j int) bool
    Less reports whether x[i] should be ordered before x[j], as required by
    the sort Interface. Note that floating-point comparison by itself is
    not a transitive relation: it does not report a consistent ordering for
    not-a-number (NaN) values. This implementation of Less places NaN values
    before any others, by using:

        x[i] < x[j] || (math.IsNaN(x[i]) && !math.IsNaN(x[j]))

func (p Float64Slice) Search(x float64) int
    Search returns the result of applying SearchFloat64s to the receiver and x.

func (x Float64Slice) Sort()
    Sort is a convenience method: x.Sort() calls Sort(x).

func (x Float64Slice) Swap(i, j int)

type IntSlice []int
    IntSlice attaches the methods of Interface to []int, sorting in increasing
    order.

func (x IntSlice) Len() int

func (x IntSlice) Less(i, j int) bool

func (p IntSlice) Search(x int) int
    Search returns the result of applying SearchInts to the receiver and x.

func (x IntSlice) Sort()
    Sort is a convenience method: x.Sort() calls Sort(x).

func (x IntSlice) Swap(i, j int)

type Interface interface {
	// Len is the number of elements in the collection.
	Len() int

	// Less reports whether the element with index i
	// must sort before the element with index j.
	//
	// If both Less(i, j) and Less(j, i) are false,
	// then the elements at index i and j are considered equal.
	// Sort may place equal elements in any order in the final result,
	// while Stable preserves the original input order of equal elements.
	//
	// Less must describe a [Strict Weak Ordering]. For example:
	//  - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
	//  - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
	//
	// Note that floating-point comparison (the < operator on float32 or float64 values)
	// is not a strict weak ordering when not-a-number (NaN) values are involved.
	// See Float64Slice.Less for a correct implementation for floating-point values.
	//
	// [Strict Weak Ordering]: https://en.wikipedia.org/wiki/Weak_ordering#Strict_weak_orderings
	Less(i, j int) bool

	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
}
    An implementation of Interface can be sorted by the routines in this
    package. The methods refer to elements of the underlying collection by
    integer index.

func Reverse(data Interface) Interface
    Reverse returns the reverse order for data.

type StringSlice []string
    StringSlice attaches the methods of Interface to []string, sorting in
    increasing order.

func (x StringSlice) Len() int

func (x StringSlice) Less(i, j int) bool

func (p StringSlice) Search(x string) int
    Search returns the result of applying SearchStrings to the receiver and x.

func (x StringSlice) Sort()
    Sort is a convenience method: x.Sort() calls Sort(x).

func (x StringSlice) Swap(i, j int)

## idiomatic usage

Idiomatic usage of `sort` drawn from the package's own runnable examples. Keywords: sort sort usage example idiomatic how to use basic Find Float64s.

```go
package main

import (
	"fmt"
	"sort"
)

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s: %d", p.Name, p.Age)
}

// ByAge implements sort.Interface for []Person based on
// the Age field.
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

func main() {
	people := []Person{
		{"Bob", 31},
		{"John", 42},
		{"Michael", 17},
		{"Jenny", 26},
	}

	fmt.Println(people)
	// There are two ways to sort a slice. First, one can define
	// a set of methods for the slice type, as with ByAge, and
	// call sort.Sort. In this first example we use that technique.
	sort.Sort(ByAge(people))
	fmt.Println(people)

	// The other way is to use sort.Slice with a custom Less
	// function, which can be provided as a closure. In this
	// case no methods are needed. (And if they exist, they
	// are ignored.) Here we re-sort in reverse order: compare
	// the closure with ByAge.Less.
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age > people[j].Age
	})
	fmt.Println(people)

}

// Output:
// [Bob: 31 John: 42 Michael: 17 Jenny: 26]
// [Michael: 17 Jenny: 26 Bob: 31 John: 42]
// [John: 42 Bob: 31 Jenny: 26 Michael: 17]
```

```go
package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	a := []string{"apple", "banana", "lemon", "mango", "pear", "strawberry"}

	for _, x := range []string{"banana", "orange"} {
		i, found := sort.Find(len(a), func(i int) int {
			return strings.Compare(x, a[i])
		})
		if found {
			fmt.Printf("found %s at index %d\n", x, i)
		} else {
			fmt.Printf("%s not found, would insert at %d\n", x, i)
		}
	}

}

// Output:
// found banana at index 1
// orange not found, would insert at 4
```

```go
package main

import (
	"fmt"
	"math"
	"sort"
)

func main() {
	s := []float64{5.2, -1.3, 0.7, -3.8, 2.6} // unsorted
	sort.Float64s(s)
	fmt.Println(s)

	s = []float64{math.Inf(1), math.NaN(), math.Inf(-1), 0.0} // unsorted
	sort.Float64s(s)
	fmt.Println(s)

}

// Output:
// [-3.8 -1.3 0.7 2.6 5.2]
// [NaN -Inf 0 +Inf]
```

## key idioms (curated)

There are three idiomatic ways to sort, newest first. For a slice, prefer the generic `slices` package
(Go 1.21+): `slices.Sort(s)` for ordered element types, or `slices.SortFunc(s, cmp)` where `cmp(a, b)`
returns a negative/zero/positive int (use `cmp.Compare` or `strings.Compare`, or subtract for ints). For
ad-hoc sorting without a custom type, `sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })`. To make
your OWN collection sortable, implement `sort.Interface` (Len, Less, Swap) and call `sort.Sort`. The
`slices` package also has `Contains`, `Index`, `BinarySearch`, `Min`, `Max`, `Reverse`, and `Equal` -
reach for these before hand-writing loops. Keywords: sort slices sort.Slice sort.SortFunc sort.Sort
sort.Interface Len Less Swap slices.Sort slices.SortFunc slices.Contains slices.Index slices.BinarySearch
slices.Reverse cmp.Compare comparator order ascending descending stable SliceStable by field key.

```go
package main

import (
	"cmp"
	"fmt"
	"slices"
	"sort"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	nums := []int{3, 1, 2}
	slices.Sort(nums) // [1 2 3] - simplest for ordered types
	fmt.Println(nums, slices.Contains(nums, 2))

	people := []Person{{"Bo", 30}, {"Al", 30}, {"Cy", 25}}

	// Sort by Age, then Name, with a comparator (Go 1.21+).
	slices.SortFunc(people, func(a, b Person) int {
		if d := cmp.Compare(a.Age, b.Age); d != 0 {
			return d
		}
		return cmp.Compare(a.Name, b.Name)
	})
	fmt.Println(people)

	// Older ad-hoc style, still common:
	sort.Slice(people, func(i, j int) bool { return people[i].Name < people[j].Name })
	fmt.Println(people)
}
```
