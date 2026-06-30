# Luhn Algorithm (Mod 10 Checksum)

The Luhn algorithm is a simple checksum that detects single-digit errors and most adjacent transpositions in identification numbers such as credit cards, IMEIs, and national IDs. Walking from the rightmost digit leftward, double every second digit; if a doubled value exceeds 9 subtract 9 (equivalently sum its two digits), then sum all digits. The number is valid if that total is divisible by 10. It runs in O(n) time and O(1) space. It is an error-detecting check digit, NOT a cryptographic or security mechanism. Keywords: luhn algorithm mod 10 modulus 10 checksum check digit credit card validation IMEI double every second digit error detection validate number identification check digit verify card

## implementation

```go
package checksum

// Luhn reports whether the digit string s passes the Luhn (mod 10) check.
// It processes digits right-to-left, doubling every second one. Returns false
// if s is empty or contains a non-digit. O(n) time, O(1) space.
func Luhn(s string) bool {
	sum := 0
	double := false // the rightmost digit is not doubled
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		if c < '0' || c > '9' {
			return false
		}
		d := int(c - '0')
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	return s != "" && sum%10 == 0
}

// LuhnCheckDigit returns the digit (0-9) that, appended to the digit string s,
// makes the whole number pass the Luhn check.
func LuhnCheckDigit(s string) int {
	sum := 0
	double := true // the appended digit shifts parity: next-to-last gets doubled
	for i := len(s) - 1; i >= 0; i-- {
		d := int(s[i] - '0')
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	return (10 - sum%10) % 10
}
```

## usage / test

```go
package checksum

import (
	"strconv"
	"testing"
)

func TestLuhn(t *testing.T) {
	valid := []string{"4539148803436467", "79927398713", "0"}
	for _, s := range valid {
		if !Luhn(s) {
			t.Errorf("Luhn(%q) = false, want true", s)
		}
	}
	invalid := []string{"4539148803436468", "79927398710", "", "12a4"}
	for _, s := range invalid {
		if Luhn(s) {
			t.Errorf("Luhn(%q) = true, want false", s)
		}
	}
	// appending the computed check digit must yield a valid number
	base := "7992739871"
	full := base + strconv.Itoa(LuhnCheckDigit(base))
	if !Luhn(full) {
		t.Errorf("computed check digit produced invalid number %q", full)
	}
}
```
