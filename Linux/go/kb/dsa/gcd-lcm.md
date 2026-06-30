# Greatest Common Divisor (GCD) and Least Common Multiple (LCM)

The greatest common divisor of two integers is the largest integer dividing both; the Euclidean algorithm computes it by repeatedly replacing (a, b) with (b, a mod b) until b is zero. The least common multiple is the smallest positive integer divisible by both and equals a/gcd(a,b)*b (divide before multiplying to avoid overflow). Use GCD to reduce fractions and LCM to find common periods or denominators. The Euclidean algorithm runs in O(log(min(a,b))) time and O(1) space. Keywords: gcd greatest common divisor lcm least common multiple euclidean algorithm euclid modulo coprime reduce fraction common denominator divisor multiple O(log n) number theory

## implementation

```go
package mathx

// GCD returns the greatest common divisor of a and b using the Euclidean
// algorithm. The result is non-negative. GCD(0, 0) is 0.
func GCD(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM returns the least common multiple of a and b. It divides before
// multiplying to reduce the chance of overflow. LCM(x, 0) is 0.
func LCM(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	g := GCD(a, b)
	res := a / g * b
	if res < 0 {
		res = -res
	}
	return res
}
```

## usage / test

```go
package mathx

import "testing"

func TestGCDLCM(t *testing.T) {
	gcdCases := []struct{ a, b, want int }{
		{12, 18, 6}, {17, 5, 1}, {0, 9, 9}, {-12, 18, 6}, {48, 36, 12},
	}
	for _, c := range gcdCases {
		if got := GCD(c.a, c.b); got != c.want {
			t.Errorf("GCD(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
	lcmCases := []struct{ a, b, want int }{
		{4, 6, 12}, {21, 6, 42}, {5, 0, 0}, {7, 7, 7},
	}
	for _, c := range lcmCases {
		if got := LCM(c.a, c.b); got != c.want {
			t.Errorf("LCM(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}
```
