# math/bits (Go standard library)

Package bits implements bit counting and manipulation functions for the

Import path: math/bits   Toolchain: go1.26.4

package bits // import "math/bits"

Package bits implements bit counting and manipulation functions for the
predeclared unsigned integer types.

Functions in this package may be implemented directly by the compiler,
for better performance. For those functions the code in this package will
not be used. Which functions are implemented by the compiler depends on the
architecture and the Go release.

CONSTANTS

const UintSize = uintSize
    UintSize is the size of a uint in bits.


FUNCTIONS

func Add(x, y, carry uint) (sum, carryOut uint)
    Add returns the sum with carry of x, y and carry: sum = x + y + carry.
    The carry input must be 0 or 1; otherwise the behavior is undefined.
    The carryOut output is guaranteed to be 0 or 1.

    This function's execution time does not depend on the inputs.

func Add32(x, y, carry uint32) (sum, carryOut uint32)
    Add32 returns the sum with carry of x, y and carry: sum = x + y + carry.
    The carry input must be 0 or 1; otherwise the behavior is undefined.
    The carryOut output is guaranteed to be 0 or 1.

    This function's execution time does not depend on the inputs.

func Add64(x, y, carry uint64) (sum, carryOut uint64)
    Add64 returns the sum with carry of x, y and carry: sum = x + y + carry.
    The carry input must be 0 or 1; otherwise the behavior is undefined.
    The carryOut output is guaranteed to be 0 or 1.

    This function's execution time does not depend on the inputs.

func Div(hi, lo, y uint) (quo, rem uint)
    Div returns the quotient and remainder of (hi, lo) divided by y: quo = (hi,
    lo)/y, rem = (hi, lo)%y with the dividend bits' upper half in parameter hi
    and the lower half in parameter lo. Div panics for y == 0 (division by zero)
    or y <= hi (quotient overflow).

func Div32(hi, lo, y uint32) (quo, rem uint32)
    Div32 returns the quotient and remainder of (hi, lo) divided by y: quo =
    (hi, lo)/y, rem = (hi, lo)%y with the dividend bits' upper half in parameter
    hi and the lower half in parameter lo. Div32 panics for y == 0 (division by
    zero) or y <= hi (quotient overflow).

func Div64(hi, lo, y uint64) (quo, rem uint64)
    Div64 returns the quotient and remainder of (hi, lo) divided by y: quo =
    (hi, lo)/y, rem = (hi, lo)%y with the dividend bits' upper half in parameter
    hi and the lower half in parameter lo. Div64 panics for y == 0 (division by
    zero) or y <= hi (quotient overflow).

func LeadingZeros(x uint) int
    LeadingZeros returns the number of leading zero bits in x; the result is
    UintSize for x == 0.

func LeadingZeros16(x uint16) int
    LeadingZeros16 returns the number of leading zero bits in x; the result is
    16 for x == 0.

func LeadingZeros32(x uint32) int
    LeadingZeros32 returns the number of leading zero bits in x; the result is
    32 for x == 0.

func LeadingZeros64(x uint64) int
    LeadingZeros64 returns the number of leading zero bits in x; the result is
    64 for x == 0.

func LeadingZeros8(x uint8) int
    LeadingZeros8 returns the number of leading zero bits in x; the result is 8
    for x == 0.

func Len(x uint) int
    Len returns the minimum number of bits required to represent x; the result
    is 0 for x == 0.

func Len16(x uint16) (n int)
    Len16 returns the minimum number of bits required to represent x; the result
    is 0 for x == 0.

func Len32(x uint32) (n int)
    Len32 returns the minimum number of bits required to represent x; the result
    is 0 for x == 0.

func Len64(x uint64) (n int)
    Len64 returns the minimum number of bits required to represent x; the result
    is 0 for x == 0.

func Len8(x uint8) int
    Len8 returns the minimum number of bits required to represent x; the result
    is 0 for x == 0.

func Mul(x, y uint) (hi, lo uint)
    Mul returns the full-width product of x and y: (hi, lo) = x * y with the
    product bits' upper half returned in hi and the lower half returned in lo.

    This function's execution time does not depend on the inputs.

func Mul32(x, y uint32) (hi, lo uint32)
    Mul32 returns the 64-bit product of x and y: (hi, lo) = x * y with the
    product bits' upper half returned in hi and the lower half returned in lo.

    This function's execution time does not depend on the inputs.

func Mul64(x, y uint64) (hi, lo uint64)
    Mul64 returns the 128-bit product of x and y: (hi, lo) = x * y with the
    product bits' upper half returned in hi and the lower half returned in lo.

    This function's execution time does not depend on the inputs.

func OnesCount(x uint) int
    OnesCount returns the number of one bits ("population count") in x.

func OnesCount16(x uint16) int
    OnesCount16 returns the number of one bits ("population count") in x.

func OnesCount32(x uint32) int
    OnesCount32 returns the number of one bits ("population count") in x.

func OnesCount64(x uint64) int
    OnesCount64 returns the number of one bits ("population count") in x.

func OnesCount8(x uint8) int
    OnesCount8 returns the number of one bits ("population count") in x.

func Rem(hi, lo, y uint) uint
    Rem returns the remainder of (hi, lo) divided by y. Rem panics for y == 0
    (division by zero) but, unlike Div, it doesn't panic on a quotient overflow.

func Rem32(hi, lo, y uint32) uint32
    Rem32 returns the remainder of (hi, lo) divided by y. Rem32 panics for y
    == 0 (division by zero) but, unlike Div32, it doesn't panic on a quotient
    overflow.

func Rem64(hi, lo, y uint64) uint64
    Rem64 returns the remainder of (hi, lo) divided by y. Rem64 panics for y
    == 0 (division by zero) but, unlike Div64, it doesn't panic on a quotient
    overflow.

func Reverse(x uint) uint
    Reverse returns the value of x with its bits in reversed order.

func Reverse16(x uint16) uint16
    Reverse16 returns the value of x with its bits in reversed order.

func Reverse32(x uint32) uint32
    Reverse32 returns the value of x with its bits in reversed order.

func Reverse64(x uint64) uint64
    Reverse64 returns the value of x with its bits in reversed order.

func Reverse8(x uint8) uint8
    Reverse8 returns the value of x with its bits in reversed order.

func ReverseBytes(x uint) uint
    ReverseBytes returns the value of x with its bytes in reversed order.

    This function's execution time does not depend on the inputs.

func ReverseBytes16(x uint16) uint16
    ReverseBytes16 returns the value of x with its bytes in reversed order.

    This function's execution time does not depend on the inputs.

func ReverseBytes32(x uint32) uint32
    ReverseBytes32 returns the value of x with its bytes in reversed order.

    This function's execution time does not depend on the inputs.

func ReverseBytes64(x uint64) uint64
    ReverseBytes64 returns the value of x with its bytes in reversed order.

    This function's execution time does not depend on the inputs.

func RotateLeft(x uint, k int) uint
    RotateLeft returns the value of x rotated left by (k mod UintSize) bits.
    To rotate x right by k bits, call RotateLeft(x, -k).

    This function's execution time does not depend on the inputs.

func RotateLeft16(x uint16, k int) uint16
    RotateLeft16 returns the value of x rotated left by (k mod 16) bits.
    To rotate x right by k bits, call RotateLeft16(x, -k).

    This function's execution time does not depend on the inputs.

func RotateLeft32(x uint32, k int) uint32
    RotateLeft32 returns the value of x rotated left by (k mod 32) bits.
    To rotate x right by k bits, call RotateLeft32(x, -k).

    This function's execution time does not depend on the inputs.

func RotateLeft64(x uint64, k int) uint64
    RotateLeft64 returns the value of x rotated left by (k mod 64) bits.
    To rotate x right by k bits, call RotateLeft64(x, -k).

    This function's execution time does not depend on the inputs.

func RotateLeft8(x uint8, k int) uint8
    RotateLeft8 returns the value of x rotated left by (k mod 8) bits. To rotate
    x right by k bits, call RotateLeft8(x, -k).

    This function's execution time does not depend on the inputs.

func Sub(x, y, borrow uint) (diff, borrowOut uint)
    Sub returns the difference of x, y and borrow: diff = x - y - borrow.
    The borrow input must be 0 or 1; otherwise the behavior is undefined.
    The borrowOut output is guaranteed to be 0 or 1.

    This function's execution time does not depend on the inputs.

func Sub32(x, y, borrow uint32) (diff, borrowOut uint32)
    Sub32 returns the difference of x, y and borrow, diff = x - y - borrow.
    The borrow input must be 0 or 1; otherwise the behavior is undefined.
    The borrowOut output is guaranteed to be 0 or 1.

    This function's execution time does not depend on the inputs.

func Sub64(x, y, borrow uint64) (diff, borrowOut uint64)
    Sub64 returns the difference of x, y and borrow: diff = x - y - borrow.
    The borrow input must be 0 or 1; otherwise the behavior is undefined.
    The borrowOut output is guaranteed to be 0 or 1.

    This function's execution time does not depend on the inputs.

func TrailingZeros(x uint) int
    TrailingZeros returns the number of trailing zero bits in x; the result is
    UintSize for x == 0.

func TrailingZeros16(x uint16) int
    TrailingZeros16 returns the number of trailing zero bits in x; the result is
    16 for x == 0.

func TrailingZeros32(x uint32) int
    TrailingZeros32 returns the number of trailing zero bits in x; the result is
    32 for x == 0.

func TrailingZeros64(x uint64) int
    TrailingZeros64 returns the number of trailing zero bits in x; the result is
    64 for x == 0.

func TrailingZeros8(x uint8) int
    TrailingZeros8 returns the number of trailing zero bits in x; the result is
    8 for x == 0.

## idiomatic usage

Count set/leading/trailing bits, rotate and reverse bits, and do full-width multi-word add/sub/mul/div arithmetic with carry/borrow. Keywords: bits.OnesCount bits.LeadingZeros bits.TrailingZeros bits.Len bits.RotateLeft bits.Reverse bits.ReverseBytes bits.Add64 bits.Sub64 bits.Mul64 bits.Div64 popcount bit counting carry borrow wide arithmetic.

```go
// Bit counting and manipulation.
fmt.Printf("OnesCount(%b) = %d\n", 14, bits.OnesCount(14))          // popcount -> 3
fmt.Printf("LeadingZeros8(%08b) = %d\n", 1, bits.LeadingZeros8(1))  // 7
fmt.Printf("TrailingZeros8(%08b) = %d\n", 14, bits.TrailingZeros8(14)) // 1
fmt.Printf("Len8(%08b) = %d\n", 8, bits.Len8(8))                    // 4
fmt.Printf("%08b\n", bits.RotateLeft8(15, 2))                       // 00111100
fmt.Printf("%08b\n", bits.Reverse8(19))                             // 11001000
```

```go
// Full-width addition with carry across two words.
n1 := []uint64{33, 12}
n2 := []uint64{21, 23}
d1, carry := bits.Add64(n1[1], n2[1], 0)
d0, _ := bits.Add64(n1[0], n2[0], carry)
fmt.Printf("%v + %v = %v\n", n1, n2, []uint64{d0, d1})

// 64x64 -> 128-bit multiply, and 128/64 divide.
hi, lo := bits.Mul64(0x8000000000000000, 2)
fmt.Println(hi, lo) // 1 0
quo, rem := bits.Div64(0, 6, 3)
fmt.Println(quo, rem) // 2 0
```
