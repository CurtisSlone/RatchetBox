# Palindrome Check

A palindrome reads the same forward and backward. The check walks two indices inward from both ends, comparing characters until they meet. Decide up front whether the comparison is exact (byte/rune for byte) or "clean" (ignoring case, spaces, and punctuation) - the clean version is what people usually mean for phrases like "A man, a plan, a canal: Panama". Work on runes, not bytes, so multi-byte UTF-8 characters compare correctly. It runs in O(n) time and O(1) extra space (or O(n) to first build a filtered rune slice). Keywords: palindrome check is palindrome reverse string two pointer mirror symmetric same forwards backwards rune unicode case-insensitive alphanumeric clean phrase O(n) strings

## implementation

```go
package strs

import "unicode"

// IsPalindrome reports whether s reads the same forwards and backwards,
// comparing runes exactly. O(n) time, O(n) space to decode runes.
func IsPalindrome(s string) bool {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		if r[i] != r[j] {
			return false
		}
	}
	return true
}

// IsPalindromeClean reports whether s is a palindrome ignoring case and any
// non-alphanumeric characters (so "A man, a plan, a canal: Panama" qualifies).
func IsPalindromeClean(s string) bool {
	r := []rune(s)
	i, j := 0, len(r)-1
	for i < j {
		if !isAlnum(r[i]) {
			i++
			continue
		}
		if !isAlnum(r[j]) {
			j--
			continue
		}
		if unicode.ToLower(r[i]) != unicode.ToLower(r[j]) {
			return false
		}
		i++
		j--
	}
	return true
}

func isAlnum(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
```

## usage / test

```go
package strs

import "testing"

func TestIsPalindrome(t *testing.T) {
	exact := map[string]bool{"": true, "a": true, "aba": true, "abba": true, "abc": false, "ab": false}
	for s, want := range exact {
		if got := IsPalindrome(s); got != want {
			t.Errorf("IsPalindrome(%q) = %v, want %v", s, got, want)
		}
	}
	clean := map[string]bool{
		"A man, a plan, a canal: Panama": true,
		"No 'x' in Nixon":                true,
		"race a car":                     false,
	}
	for s, want := range clean {
		if got := IsPalindromeClean(s); got != want {
			t.Errorf("IsPalindromeClean(%q) = %v, want %v", s, got, want)
		}
	}
}
```
