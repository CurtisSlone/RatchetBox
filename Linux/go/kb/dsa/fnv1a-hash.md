# FNV-1a Hash

FNV-1a is a fast, simple non-cryptographic hash for byte strings: start from a fixed offset basis, then for each byte XOR it into the hash and multiply by a fixed prime. Use it for hash tables, checksums, bloom filters, and bucket selection where speed and good distribution matter but cryptographic security does NOT. It runs in O(n) time and O(1) space. Do not use it for passwords, signatures, or anything adversarial; it is not collision-resistant against attackers. Keywords: FNV-1a FNV hash non-cryptographic fast hash function 32-bit 64-bit offset basis prime XOR multiply byte hash table bucket bloom filter checksum string hash distribution

## implementation

```go
package hashx

// FNV-1a 32-bit constants.
const (
	fnvOffset32 uint32 = 2166136261
	fnvPrime32  uint32 = 16777619
)

// FNV1a32 returns the 32-bit FNV-1a hash of data. For each byte it XORs the
// byte into the hash then multiplies by the FNV prime. O(n) time, O(1) space.
func FNV1a32(data []byte) uint32 {
	hash := fnvOffset32
	for _, b := range data {
		hash ^= uint32(b)
		hash *= fnvPrime32
	}
	return hash
}

// FNV-1a 64-bit constants.
const (
	fnvOffset64 uint64 = 14695981039346656037
	fnvPrime64  uint64 = 1099511628211
)

// FNV1a64 returns the 64-bit FNV-1a hash of data.
func FNV1a64(data []byte) uint64 {
	hash := fnvOffset64
	for _, b := range data {
		hash ^= uint64(b)
		hash *= fnvPrime64
	}
	return hash
}
```

## usage / test

```go
package hashx

import (
	"hash/fnv"
	"testing"
)

func TestFNV1a(t *testing.T) {
	inputs := [][]byte{nil, []byte(""), []byte("a"), []byte("hello"), []byte("The quick brown fox")}
	for _, in := range inputs {
		// cross-check against the standard library implementation
		h32 := fnv.New32a()
		h32.Write(in)
		if got := FNV1a32(in); got != h32.Sum32() {
			t.Errorf("FNV1a32(%q) = %d, want %d", in, got, h32.Sum32())
		}
		h64 := fnv.New64a()
		h64.Write(in)
		if got := FNV1a64(in); got != h64.Sum64() {
			t.Errorf("FNV1a64(%q) = %d, want %d", in, got, h64.Sum64())
		}
	}
	// same input must hash to the same value (determinism)
	if FNV1a32([]byte("x")) != FNV1a32([]byte("x")) {
		t.Error("hash not deterministic")
	}
}
```
