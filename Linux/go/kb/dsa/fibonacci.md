# Fibonacci Numbers (Iterative and Matrix)

The Fibonacci sequence is defined by F(0)=0, F(1)=1, and F(n)=F(n-1)+F(n-2). The iterative method keeps the last two values and adds them n times in O(n) time and O(1) space; this is the right default. The matrix-power method raises [[1,1],[1,0]] to the n-th power using fast exponentiation by squaring, computing F(n) in O(log n) time, which matters for very large n. Avoid the naive recursion, which is exponential O(phi^n). Note that F(n) overflows int64 around n=92. Keywords: fibonacci sequence fib golden ratio iterative dynamic programming matrix exponentiation fast doubling O(n) O(log n) F(n) F(n-1) recurrence number sequence

## implementation

```go
package mathx

// Fibonacci returns F(n) for n >= 0 iteratively in O(n) time, O(1) space.
func Fibonacci(n int) int {
	if n < 2 {
		return n
	}
	prev, curr := 0, 1
	for i := 2; i <= n; i++ {
		prev, curr = curr, prev+curr
	}
	return curr
}

// FibonacciMatrix returns F(n) for n >= 0 in O(log n) time by raising the
// matrix [[1,1],[1,0]] to the n-th power via exponentiation by squaring.
func FibonacciMatrix(n int) int {
	if n < 2 {
		return n
	}
	result := [2][2]int{{1, 0}, {0, 1}} // identity
	base := [2][2]int{{1, 1}, {1, 0}}
	for n > 0 {
		if n&1 == 1 {
			result = mul2x2(result, base)
		}
		base = mul2x2(base, base)
		n >>= 1
	}
	return result[0][1] // F(n)
}

func mul2x2(a, b [2][2]int) [2][2]int {
	return [2][2]int{
		{a[0][0]*b[0][0] + a[0][1]*b[1][0], a[0][0]*b[0][1] + a[0][1]*b[1][1]},
		{a[1][0]*b[0][0] + a[1][1]*b[1][0], a[1][0]*b[0][1] + a[1][1]*b[1][1]},
	}
}
```

## usage / test

```go
package mathx

import "testing"

func TestFibonacci(t *testing.T) {
	want := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89}
	for n, w := range want {
		if got := Fibonacci(n); got != w {
			t.Errorf("Fibonacci(%d) = %d, want %d", n, got, w)
		}
		if got := FibonacciMatrix(n); got != w {
			t.Errorf("FibonacciMatrix(%d) = %d, want %d", n, got, w)
		}
	}
	// both methods must agree for larger n
	for n := 0; n < 90; n++ {
		if Fibonacci(n) != FibonacciMatrix(n) {
			t.Errorf("methods disagree at n=%d", n)
		}
	}
}
```
