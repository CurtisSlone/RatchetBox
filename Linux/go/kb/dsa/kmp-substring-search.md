# KMP Substring Search

The Knuth-Morris-Pratt algorithm finds all occurrences of a pattern in a text in O(n + m) time without ever backtracking in the text. It first builds a "failure" / longest-proper-prefix-suffix (LPS) table for the pattern, recording for each position the length of the longest prefix that is also a suffix; on a mismatch it uses that table to slide the pattern forward instead of rechecking characters. Use it for repeated or large-text substring search where naive O(n*m) scanning is too slow. Time is O(n + m), space is O(m) for the table. Keywords: KMP Knuth-Morris-Pratt substring search string matching pattern search find substring index LPS failure function prefix suffix table strings.Index O(n+m) text pattern occurrences

## implementation

```go
package strs

// KMPSearch returns the starting indices of every (possibly overlapping)
// occurrence of pattern in text. O(len(text) + len(pattern)) time.
func KMPSearch(text, pattern string) []int {
	if pattern == "" {
		// empty pattern matches at every position, including the end
		idx := make([]int, len(text)+1)
		for i := range idx {
			idx[i] = i
		}
		return idx
	}
	lps := buildLPS(pattern)
	var matches []int
	j := 0 // number of pattern chars currently matched
	for i := 0; i < len(text); i++ {
		for j > 0 && text[i] != pattern[j] {
			j = lps[j-1] // fall back using the table, no text backtrack
		}
		if text[i] == pattern[j] {
			j++
		}
		if j == len(pattern) {
			matches = append(matches, i-j+1)
			j = lps[j-1]
		}
	}
	return matches
}

// buildLPS computes the longest-proper-prefix-suffix table for pattern.
func buildLPS(pattern string) []int {
	lps := make([]int, len(pattern))
	length := 0
	for i := 1; i < len(pattern); {
		if pattern[i] == pattern[length] {
			length++
			lps[i] = length
			i++
		} else if length > 0 {
			length = lps[length-1]
		} else {
			lps[i] = 0
			i++
		}
	}
	return lps
}
```

## usage / test

```go
package strs

import (
	"slices"
	"strings"
	"testing"
)

func TestKMPSearch(t *testing.T) {
	cases := []struct {
		text, pattern string
		want          []int
	}{
		{"abxabcabcaby", "abcaby", []int{6}},
		{"aaaaa", "aa", []int{0, 1, 2, 3}},
		{"abcabcabc", "abc", []int{0, 3, 6}},
		{"hello", "z", nil},
	}
	for _, c := range cases {
		if got := KMPSearch(c.text, c.pattern); !slices.Equal(got, c.want) {
			t.Errorf("KMPSearch(%q, %q) = %v, want %v", c.text, c.pattern, got, c.want)
		}
		// first match must agree with strings.Index
		got := KMPSearch(c.text, c.pattern)
		first := -1
		if len(got) > 0 {
			first = got[0]
		}
		if first != strings.Index(c.text, c.pattern) {
			t.Errorf("first index mismatch for %q in %q", c.pattern, c.text)
		}
	}
}
```
