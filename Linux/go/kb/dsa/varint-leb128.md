# Varint / LEB128 Encoding

A variable-length integer (varint, LEB128) stores an unsigned integer in as few bytes as possible: it emits 7 bits of the value per byte, little-endian, and sets the high bit (0x80) of every byte except the last as a "more bytes follow" continuation flag. Small numbers take one byte, large ones up to ten for a 64-bit value. Use it for compact binary formats and wire protocols (Protocol Buffers, DWARF, WebAssembly). Signed values use zig-zag mapping so that small-magnitude negatives also stay small. Encoding and decoding are O(bytes) ~ O(log value). Keywords: varint LEB128 variable length integer encoding base-128 continuation bit MSB protobuf wire format zig-zag signed unsigned uvarint compact integer encode decode 7 bits per byte little-endian binary.PutUvarint

## implementation

```go
package encoding

import "errors"

// AppendUvarint appends the LEB128 encoding of x to dst and returns the result.
// Each byte carries 7 value bits; the high bit marks "more bytes follow".
func AppendUvarint(dst []byte, x uint64) []byte {
	for x >= 0x80 {
		dst = append(dst, byte(x)|0x80) // low 7 bits + continuation flag
		x >>= 7
	}
	return append(dst, byte(x)) // final byte, high bit clear
}

// Uvarint decodes one LEB128 unsigned integer from buf, returning the value and
// the number of bytes consumed, or an error if the input is truncated/overlong.
func Uvarint(buf []byte) (uint64, int, error) {
	var x uint64
	var shift uint
	for i, b := range buf {
		if i == 10 { // a 64-bit value needs at most 10 bytes
			return 0, 0, errors.New("varint overflows 64 bits")
		}
		x |= uint64(b&0x7F) << shift
		if b&0x80 == 0 { // last byte
			return x, i + 1, nil
		}
		shift += 7
	}
	return 0, 0, errors.New("truncated varint")
}

// AppendVarint encodes a signed integer using zig-zag mapping so small
// magnitudes (including negatives) stay short, then LEB128.
func AppendVarint(dst []byte, x int64) []byte {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux // zig-zag: 0,-1,1,-2,2,... -> 0,1,2,3,4,...
	}
	return AppendUvarint(dst, ux)
}

// Varint decodes a zig-zag LEB128 signed integer.
func Varint(buf []byte) (int64, int, error) {
	ux, n, err := Uvarint(buf)
	if err != nil {
		return 0, 0, err
	}
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, n, nil
}
```

## usage / test

```go
package encoding

import (
	"encoding/binary"
	"testing"
)

func TestVarintRoundTrip(t *testing.T) {
	unsigned := []uint64{0, 1, 127, 128, 300, 16384, 1<<32 - 1, 1<<64 - 1}
	for _, v := range unsigned {
		enc := AppendUvarint(nil, v)
		// cross-check encoding against the standard library
		ref := make([]byte, binary.MaxVarintLen64)
		ref = ref[:binary.PutUvarint(ref, v)]
		if string(enc) != string(ref) {
			t.Errorf("AppendUvarint(%d) = % x, want % x", v, enc, ref)
		}
		got, n, err := Uvarint(enc)
		if err != nil || got != v || n != len(enc) {
			t.Errorf("Uvarint round trip %d: got %d n=%d err=%v", v, got, n, err)
		}
	}

	signed := []int64{0, -1, 1, -300, 300, -1 << 40, 1 << 40}
	for _, v := range signed {
		enc := AppendVarint(nil, v)
		got, n, err := Varint(enc)
		if err != nil || got != v || n != len(enc) {
			t.Errorf("Varint round trip %d: got %d n=%d err=%v", v, got, n, err)
		}
	}
}
```
