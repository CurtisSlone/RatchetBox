# Factorial

The factorial of a non-negative integer n, written n!, is the product 1*2*...*n, with 0! defined as 1. Use it in combinatorics (permutations and combinations), probability, and series expansions. The iterative form runs in O(n) multiplications and O(1) space and is preferred over recursion to avoid stack growth. Note that n! overflows int64 around n=21; use math/big for larger inputs. Keywords: factorial n! product permutations combinations iterative recursive O(n) big integer combinatorics gamma overflow 0! equals 1

## implementation

```go
package mathx

import "math/big"

// Factorial returns n! for 0 <= n <= 20 using int64. It returns 1 for n == 0.
// It panics on negative input and overflows silently past 20 (use BigFactorial).
func Factorial(n int) int64 {
	if n < 0 {
		panic("Factorial: negative input")
	}
	result := int64(1)
	for i := int64(2); i <= int64(n); i++ {
		result *= i
	}
	return result
}

// BigFactorial returns n! as an arbitrary-precision integer for any n >= 0.
func BigFactorial(n int) *big.Int {
	if n < 0 {
		panic("BigFactorial: negative input")
	}
	result := big.NewInt(1)
	for i := int64(2); i <= int64(n); i++ {
		result.Mul(result, big.NewInt(i))
	}
	return result
}
```

## usage / test

```go
package mathx

import "testing"

func TestFactorial(t *testing.T) {
	cases := []struct {
		n    int
		want int64
	}{
		{0, 1}, {1, 1}, {5, 120}, {10, 3628800}, {20, 2432902008176640000},
	}
	for _, c := range cases {
		if got := Factorial(c.n); got != c.want {
			t.Errorf("Factorial(%d) = %d, want %d", c.n, got, c.want)
		}
		if got := BigFactorial(c.n).Int64(); got != c.want {
			t.Errorf("BigFactorial(%d) = %d, want %d", c.n, got, c.want)
		}
	}
	// big factorial beyond int64 range
	if BigFactorial(25).String() != "15511210043330985984000000" {
		t.Errorf("BigFactorial(25) wrong")
	}
}
```
