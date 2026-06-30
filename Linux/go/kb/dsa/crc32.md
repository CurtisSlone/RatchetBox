# CRC (Cyclic Redundancy Check)

A cyclic redundancy check treats the data as a large binary polynomial and computes the remainder when divided (mod 2) by a fixed generator polynomial; that remainder is the checksum. CRCs catch burst transmission and storage errors far better than a simple sum, which is why they appear in Ethernet, ZIP, PNG, and disk formats. A 256-entry lookup table makes it run in O(n) time with O(1) working space after table setup. Below is a table-driven CRC-32 using the standard IEEE (reflected) polynomial. It is error-detecting, NOT cryptographic. Keywords: CRC cyclic redundancy check CRC32 CRC-32 IEEE checksum polynomial division lookup table reflected error detection Ethernet zip png frame check sequence hash data integrity

## implementation

```go
package checksum

// crc32IEEETable holds the precomputed reflected CRC-32 (IEEE) lookup table.
var crc32IEEETable = makeCRC32Table(0xEDB88320)

// makeCRC32Table builds the 256-entry table for a reflected CRC-32 with the
// given (already-reflected) polynomial, e.g. 0xEDB88320 for IEEE.
func makeCRC32Table(poly uint32) [256]uint32 {
	var t [256]uint32
	for i := uint32(0); i < 256; i++ {
		crc := i
		for j := 0; j < 8; j++ {
			if crc&1 == 1 {
				crc = (crc >> 1) ^ poly
			} else {
				crc >>= 1
			}
		}
		t[i] = crc
	}
	return t
}

// CRC32 returns the CRC-32 (IEEE) checksum of data. O(n) time, O(1) space.
func CRC32(data []byte) uint32 {
	crc := ^uint32(0) // start with all ones
	for _, b := range data {
		crc = crc32IEEETable[byte(crc)^b] ^ (crc >> 8)
	}
	return ^crc // final XOR with all ones
}
```

## usage / test

```go
package checksum

import (
	"hash/crc32"
	"testing"
)

func TestCRC32(t *testing.T) {
	inputs := [][]byte{nil, []byte(""), []byte("a"), []byte("123456789"), []byte("The quick brown fox")}
	for _, in := range inputs {
		// cross-check against the standard library IEEE CRC-32
		want := crc32.ChecksumIEEE(in)
		if got := CRC32(in); got != want {
			t.Errorf("CRC32(%q) = %08x, want %08x", in, got, want)
		}
	}
	// the canonical "123456789" CRC-32/IEEE check value is 0xCBF43926
	if got := CRC32([]byte("123456789")); got != 0xCBF43926 {
		t.Errorf("CRC32 check value = %08x, want CBF43926", got)
	}
}
```
