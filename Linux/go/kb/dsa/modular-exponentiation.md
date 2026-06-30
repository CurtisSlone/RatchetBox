# Modular Exponentiation

Modular exponentiation computes (base^exp) mod m efficiently using exponentiation by squaring: it scans the exponent's binary digits, squaring a running base each step and multiplying it into the result whenever the current bit is 1, reducing modulo m throughout to keep numbers small. Use it for cryptography (RSA, Diffie-Hellman), hashing, and any large-power-mod computation where computing base^exp directly would overflow. It runs in O(log exp) multiplications and O(1) space. Reduce the base modulo m first and use a wide accumulator (or math/big) to avoid overflow in the multiply. Keywords: modular exponentiation modpow power mod exponentiation by squaring binary exponentiation fast power RSA cryptography (base^exp) mod m O(log n) square and multiply pow modulus

## implementation

```go
package mathx

// ModPow returns (base^exp) mod m using exponentiation by squaring.
// O(log exp) multiplications. Requires exp >= 0 and m > 0. Uses uint64
// intermediates so it is safe for moduli up to about 2^32.
func ModPow(base, exp, m uint64) uint64 {
	if m == 1 {
		return 0
	}
	result := uint64(1)
	base %= m
	for exp > 0 {
		if exp&1 == 1 { // current low bit set: multiply base into result
			result = (result * base) % m
		}
		exp >>= 1
		base = (base * base) % m // square the base for the next bit
	}
	return result
}
```

## usage / test

```go
package mathx

import (
	"math/big"
	"testing"
)

func TestModPow(t *testing.T) {
	cases := []struct{ base, exp, m uint64 }{
		{2, 10, 1000}, {3, 0, 7}, {7, 256, 13}, {10, 9, 6}, {123, 45, 67},
	}
	for _, c := range cases {
		got := ModPow(c.base, c.exp, c.m)
		// cross-check against math/big
		want := new(big.Int).Exp(
			new(big.Int).SetUint64(c.base),
			new(big.Int).SetUint64(c.exp),
			new(big.Int).SetUint64(c.m),
		).Uint64()
		if got != want {
			t.Errorf("ModPow(%d, %d, %d) = %d, want %d", c.base, c.exp, c.m, got, want)
		}
	}
	if got := ModPow(5, 3, 1); got != 0 {
		t.Errorf("ModPow mod 1 = %d, want 0", got)
	}
}
```
