# Base-N Integer Conversion

Base-N conversion turns a non-negative integer into its string representation in an arbitrary base (radix) and parses it back. To encode, repeatedly take the remainder modulo the base to get the least-significant digit and integer-divide by the base until zero, then reverse the collected digits. To decode, accumulate value = value*base + digit across the string left to right. Use it for compact IDs, base-36/base-62 short codes, and custom number formats. Both directions are O(d) in the number of digits. The pair must ROUND-TRIP: decode(encode(n)) == n for any n and valid base. Keywords: base-n conversion radix base conversion encode decode integer to string parse base-2 base-16 base-36 base-62 hexadecimal binary octal change base digits round trip strconv FormatInt ParseInt

## implementation

```go
package encoding

import (
	"errors"
	"strings"
)

const digits = "0123456789abcdefghijklmnopqrstuvwxyz" // supports bases 2..36

// ToBase encodes a non-negative integer n as a string in the given base
// (2..36). O(log_base(n)) time.
func ToBase(n, base int) (string, error) {
	if base < 2 || base > 36 {
		return "", errors.New("base must be in [2, 36]")
	}
	if n < 0 {
		return "", errors.New("n must be non-negative")
	}
	if n == 0 {
		return "0", nil
	}
	var sb strings.Builder
	for n > 0 {
		sb.WriteByte(digits[n%base]) // least-significant digit first
		n /= base
	}
	return reverse(sb.String()), nil
}

// FromBase decodes a string in the given base (2..36) back to an integer.
func FromBase(s string, base int) (int, error) {
	if base < 2 || base > 36 {
		return 0, errors.New("base must be in [2, 36]")
	}
	if s == "" {
		return 0, errors.New("empty string")
	}
	n := 0
	for i := 0; i < len(s); i++ {
		d := strings.IndexByte(digits, lower(s[i]))
		if d < 0 || d >= base {
			return 0, errors.New("invalid digit for base")
		}
		n = n*base + d
	}
	return n, nil
}

func reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func lower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}
```

## usage / test

```go
package encoding

import "testing"

func TestBaseNRoundTrip(t *testing.T) {
	values := []int{0, 1, 7, 10, 15, 16, 255, 1000, 123456789}
	bases := []int{2, 8, 10, 16, 36}
	for _, v := range values {
		for _, b := range bases {
			s, err := ToBase(v, b)
			if err != nil {
				t.Fatalf("ToBase(%d, %d): %v", v, b, err)
			}
			got, err := FromBase(s, b)
			if err != nil {
				t.Fatalf("FromBase(%q, %d): %v", s, b, err)
			}
			if got != v {
				t.Errorf("round trip base %d: %d -> %q -> %d", b, v, s, got)
			}
		}
	}
	// known encodings
	if s, _ := ToBase(255, 16); s != "ff" {
		t.Errorf("ToBase(255,16) = %q, want ff", s)
	}
	if s, _ := ToBase(10, 2); s != "1010" {
		t.Errorf("ToBase(10,2) = %q, want 1010", s)
	}
}
```
