# Trie (Prefix Tree)

A trie, also called a prefix tree or digital tree, is a tree where each edge is labeled with a character and each root-to-node path spells a prefix; nodes flagged as end-of-word mark complete keys. Use it for fast string lookup, autocomplete, prefix queries ("does any stored word start with this prefix?"), spell-checkers, and IP routing tables — anything keyed by sequences of symbols. Insert, search, and prefix checks run in O(L) time where L is the length of the key, independent of how many words are stored; space is O(total characters across all keys), with each node holding a small map of child characters. Keywords: trie prefix tree digital tree radix retrieval autocomplete autosuggest typeahead prefix search startswith starts-with string dictionary word insert add search find contains has delete remove children character rune map spell-check completion.

## implementation

```go
package trie

// node is a single trie node. children maps the next rune to a child node.
type node struct {
	children map[rune]*node
	isWord   bool // true if a stored key ends here
}

func newNode() *node {
	return &node{children: make(map[rune]*node)}
}

// Trie is a prefix tree of string keys.
type Trie struct {
	root *node
	size int
}

func New() *Trie { return &Trie{root: newNode()} }

// Len returns the number of distinct words stored.
func (t *Trie) Len() int { return t.size }

// Insert adds word to the trie.
func (t *Trie) Insert(word string) {
	cur := t.root
	for _, r := range word {
		next, ok := cur.children[r]
		if !ok {
			next = newNode()
			cur.children[r] = next
		}
		cur = next
	}
	if !cur.isWord {
		cur.isWord = true
		t.size++
	}
}

// find walks to the node at the end of s, or nil if the path does not exist.
func (t *Trie) find(s string) *node {
	cur := t.root
	for _, r := range s {
		next, ok := cur.children[r]
		if !ok {
			return nil
		}
		cur = next
	}
	return cur
}

// Contains reports whether word is a stored key (not just a prefix).
func (t *Trie) Contains(word string) bool {
	n := t.find(word)
	return n != nil && n.isWord
}

// StartsWith reports whether any stored word has the given prefix.
func (t *Trie) StartsWith(prefix string) bool {
	return t.find(prefix) != nil
}

// WordsWithPrefix returns all stored words that begin with prefix, in
// lexicographic order is NOT guaranteed (map iteration); callers sort if needed.
func (t *Trie) WordsWithPrefix(prefix string) []string {
	start := t.find(prefix)
	if start == nil {
		return nil
	}
	var out []string
	var dfs func(n *node, path []rune)
	dfs = func(n *node, path []rune) {
		if n.isWord {
			out = append(out, prefix+string(path))
		}
		for r, child := range n.children {
			dfs(child, append(path, r))
		}
	}
	dfs(start, nil)
	return out
}

// Delete removes word if present and reports whether it was found.
func (t *Trie) Delete(word string) bool {
	n := t.find(word)
	if n == nil || !n.isWord {
		return false
	}
	n.isWord = false
	t.size--
	return true
}
```

## usage / test

```go
package trie

import (
	"sort"
	"testing"
)

func TestTrie(t *testing.T) {
	tr := New()
	words := []string{"cat", "car", "card", "dog", "do"}
	for _, w := range words {
		tr.Insert(w)
	}
	tr.Insert("cat") // duplicate ignored

	if tr.Len() != len(words) {
		t.Fatalf("Len = %d, want %d", tr.Len(), len(words))
	}
	// Exact membership vs prefix.
	if !tr.Contains("car") || tr.Contains("ca") {
		t.Error("Contains wrong: 'car' should be a word, 'ca' should not")
	}
	if !tr.StartsWith("ca") || tr.StartsWith("z") {
		t.Error("StartsWith wrong")
	}

	got := tr.WordsWithPrefix("ca")
	sort.Strings(got)
	want := []string{"car", "card", "cat"}
	if len(got) != len(want) {
		t.Fatalf("prefix 'ca' = %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("prefix 'ca' = %v, want %v", got, want)
		}
	}

	if !tr.Delete("car") || tr.Contains("car") {
		t.Error("Delete('car') failed")
	}
	// 'card' must still be reachable after deleting the shorter 'car'.
	if !tr.Contains("card") {
		t.Error("'card' lost after deleting 'car'")
	}
}
```
