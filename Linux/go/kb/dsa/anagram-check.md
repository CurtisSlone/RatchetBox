# Anagram Check

Two strings are anagrams if one is a rearrangement of the other, i.e. they contain exactly the same multiset of characters. The fastest check counts each character's frequency in the first string (incrementing a map) and decrements while scanning the second; they are anagrams if and only if all counts return to zero and the lengths match. Use it for word games, grouping, and duplicate-detection. The counting method is O(n) time and O(k) space for k distinct characters; an alternative is to sort both strings (O(n log n)) and compare. Decide whether to normalize case or strip spaces first. Keywords: anagram check is anagram permutation rearrange same letters character frequency count map histogram sort compare multiset O(n) strings runes case-insensitive valid anagram

## implementation

```go
package strs

// IsAnagram reports whether a and b are anagrams: same multiset of runes.
// O(n) time, O(k) space for k distinct runes. Comparison is exact (no case
// folding or space stripping; normalize the inputs first if you need that).
func IsAnagram(a, b string) bool {
	ra, rb := []rune(a), []rune(b)
	if len(ra) != len(rb) {
		return false
	}
	counts := make(map[rune]int, len(ra))
	for _, r := range ra {
		counts[r]++
	}
	for _, r := range rb {
		counts[r]--
		if counts[r] < 0 {
			return false // b has a rune a lacks
		}
	}
	// lengths matched and no count went negative, so all counts are zero
	return true
}
```

## usage / test

```go
package strs

import "testing"

func TestIsAnagram(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"listen", "silent", true},
		{"anagram", "nagaram", true},
		{"rat", "car", false},
		{"", "", true},
		{"ab", "a", false},
		{"aabb", "abab", true},
		{"héllo", "olléh", true}, // multi-byte runes
	}
	for _, c := range cases {
		if got := IsAnagram(c.a, c.b); got != c.want {
			t.Errorf("IsAnagram(%q, %q) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}
```
