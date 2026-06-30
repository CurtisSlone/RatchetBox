# Hash Table / Hash Map (with Separate-Chaining Collision Handling)

A hash table (hash map, dictionary, associative array) stores key/value pairs and uses a hash function to map each key to a bucket index, giving average O(1) insert, lookup, and delete. Collisions, where two keys hash to the same bucket, are resolved here by separate chaining: each bucket holds a linked list (or slice) of entries, and a matching key is found by scanning that short chain. The table resizes (rehashes all entries into a larger bucket array) when the load factor (entries / buckets) grows too high, which keeps chains short; worst-case operations degrade to O(n) if every key collides, and space is O(n + buckets). Go's built-in `map` is itself a highly optimized hash table; implement your own mainly to learn or to control hashing/collision policy. Keywords: hash table hash map dictionary associative array map key value bucket hash function collision separate chaining open addressing linear probing load factor rehash resize grow put get set delete remove contains keys values FNV maphash average O(1) buckets

## implementation

```go
package hashtable

import (
	"hash/maphash"
)

type entry[K comparable, V any] struct {
	key  K
	val  V
	next *entry[K, V]
}

// HashMap is a hash table with separate chaining and automatic growth.
type HashMap[K comparable, V any] struct {
	buckets []*entry[K, V]
	size    int
	seed    maphash.Seed
}

const maxLoadFactor = 0.75

// New returns an empty hash map.
func New[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{
		buckets: make([]*entry[K, V], 8),
		seed:    maphash.MakeSeed(),
	}
}

// Len reports the number of stored key/value pairs.
func (m *HashMap[K, V]) Len() int { return m.size }

// hash maps a key to a bucket index using the runtime's hasher.
func (m *HashMap[K, V]) hash(key K) int {
	var h maphash.Hash
	h.SetSeed(m.seed)
	h.WriteString(toString(key))
	return int(h.Sum64() % uint64(len(m.buckets)))
}

// Put inserts or updates the value for key in average O(1).
func (m *HashMap[K, V]) Put(key K, val V) {
	if float64(m.size+1)/float64(len(m.buckets)) > maxLoadFactor {
		m.resize()
	}
	i := m.hash(key)
	for e := m.buckets[i]; e != nil; e = e.next {
		if e.key == key {
			e.val = val // update existing key
			return
		}
	}
	m.buckets[i] = &entry[K, V]{key: key, val: val, next: m.buckets[i]}
	m.size++
}

// Get returns the value for key and whether it was found.
func (m *HashMap[K, V]) Get(key K) (V, bool) {
	var zero V
	for e := m.buckets[m.hash(key)]; e != nil; e = e.next {
		if e.key == key {
			return e.val, true
		}
	}
	return zero, false
}

// Delete removes key and reports whether it was present.
func (m *HashMap[K, V]) Delete(key K) bool {
	i := m.hash(key)
	var prev *entry[K, V]
	for e := m.buckets[i]; e != nil; e = e.next {
		if e.key == key {
			if prev == nil {
				m.buckets[i] = e.next
			} else {
				prev.next = e.next
			}
			m.size--
			return true
		}
		prev = e
	}
	return false
}

// Keys returns all stored keys in unspecified order.
func (m *HashMap[K, V]) Keys() []K {
	out := make([]K, 0, m.size)
	for _, head := range m.buckets {
		for e := head; e != nil; e = e.next {
			out = append(out, e.key)
		}
	}
	return out
}

// resize doubles the bucket count and rehashes every entry.
func (m *HashMap[K, V]) resize() {
	old := m.buckets
	m.buckets = make([]*entry[K, V], len(old)*2)
	for _, head := range old {
		for e := head; e != nil; {
			next := e.next
			i := m.hash(e.key)
			e.next = m.buckets[i]
			m.buckets[i] = e
			e = next
		}
	}
}

// toString renders a comparable key into bytes for hashing.
func toString[K comparable](key K) string {
	return string([]byte(sprint(key)))
}
```

```go
package hashtable

import "fmt"

// sprint is split out so toString stays generic-friendly.
func sprint[K comparable](key K) string { return fmt.Sprint(key) }
```

## usage / test

```go
package hashtable

import "testing"

func TestHashMap(t *testing.T) {
	m := New[string, int]()

	m.Put("alpha", 1)
	m.Put("beta", 2)
	m.Put("alpha", 10) // update, not insert

	if m.Len() != 2 {
		t.Fatalf("len: got %d want 2", m.Len())
	}
	if v, ok := m.Get("alpha"); !ok || v != 10 {
		t.Fatalf("get alpha: got %d %v want 10 true", v, ok)
	}
	if _, ok := m.Get("missing"); ok {
		t.Fatal("missing key should not be found")
	}

	if !m.Delete("beta") {
		t.Fatal("delete beta should succeed")
	}
	if _, ok := m.Get("beta"); ok {
		t.Fatal("beta should be gone")
	}
}

func TestHashMapGrowthKeepsEntries(t *testing.T) {
	m := New[int, int]()
	const n = 1000 // forces several rehashes past the load factor
	for i := 0; i < n; i++ {
		m.Put(i, i*i)
	}
	if m.Len() != n {
		t.Fatalf("len: got %d want %d", m.Len(), n)
	}
	for i := 0; i < n; i++ {
		if v, ok := m.Get(i); !ok || v != i*i {
			t.Fatalf("get %d: got %d %v", i, v, ok)
		}
	}
}
```
