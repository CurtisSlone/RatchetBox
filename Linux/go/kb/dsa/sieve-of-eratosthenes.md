# Sieve of Eratosthenes and Primality Testing

The Sieve of Eratosthenes lists all primes up to n by marking the multiples of each prime as composite, starting at the prime's square. Use it when you need many primes or repeated primality checks up to a known bound; it runs in O(n log log n) time and O(n) space. For a single ad-hoc test, trial division up to sqrt(n) is simpler at O(sqrt n) time and O(1) space. Both are included below. Keywords: sieve of eratosthenes primes prime numbers primality test is prime trial division composite mark multiples O(n log log n) sqrt number theory prime generation crossing out

## implementation

```go
package mathx

// Sieve returns all prime numbers <= n using the Sieve of Eratosthenes.
// O(n log log n) time, O(n) space.
func Sieve(n int) []int {
	if n < 2 {
		return nil
	}
	composite := make([]bool, n+1) // composite[i] true means i is not prime
	for p := 2; p*p <= n; p++ {
		if !composite[p] {
			for m := p * p; m <= n; m += p {
				composite[m] = true
			}
		}
	}
	var primes []int
	for i := 2; i <= n; i++ {
		if !composite[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

// IsPrime reports whether n is prime using trial division up to sqrt(n).
// O(sqrt(n)) time, O(1) space.
func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	if n%3 == 0 {
		return n == 3
	}
	// check candidates of the form 6k +/- 1
	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}
```

## usage / test

```go
package mathx

import (
	"slices"
	"testing"
)

func TestSieveAndIsPrime(t *testing.T) {
	got := Sieve(30)
	want := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	if !slices.Equal(got, want) {
		t.Errorf("Sieve(30) = %v, want %v", got, want)
	}
	if Sieve(1) != nil {
		t.Errorf("Sieve(1) should be empty")
	}

	primeChecks := map[int]bool{0: false, 1: false, 2: true, 3: true, 4: false, 17: true, 18: false, 97: true, 100: false}
	for n, want := range primeChecks {
		if got := IsPrime(n); got != want {
			t.Errorf("IsPrime(%d) = %v, want %v", n, got, want)
		}
	}
}
```
