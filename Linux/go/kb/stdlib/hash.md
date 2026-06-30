# hash (Go standard library)

Package hash provides interfaces for hash functions.

Import path: hash   Toolchain: go1.26.4

package hash // import "hash"

Package hash provides interfaces for hash functions.

TYPES

type Cloner interface {
	Hash
	Clone() (Cloner, error)
}
    A Cloner is a hash function whose state can be cloned, returning a value
    with equivalent and independent state.

    All Hash implementations in the standard library implement this interface,
    unless GOFIPS140=v1.0.0 is set.

    If a hash can only determine at runtime if it can be cloned (e.g.
    if it wraps another hash), Clone may return an error wrapping
    errors.ErrUnsupported. Otherwise, Clone must always return a nil error.

type Hash interface {
	// Write (via the embedded io.Writer interface) adds more data to the running hash.
	// It never returns an error.
	io.Writer

	// Sum appends the current hash to b and returns the resulting slice.
	// It does not change the underlying hash state.
	Sum(b []byte) []byte

	// Reset resets the Hash to its initial state.
	Reset()

	// Size returns the number of bytes Sum will return.
	Size() int

	// BlockSize returns the hash's underlying block size.
	// The Write method must be able to accept any amount
	// of data, but it may operate more efficiently if all writes
	// are a multiple of the block size.
	BlockSize() int
}
    Hash is the common interface implemented by all hash functions.

    Hash implementations in the standard library (e.g. hash/crc32
    and crypto/sha256) implement the encoding.BinaryMarshaler,
    encoding.BinaryAppender, encoding.BinaryUnmarshaler and Cloner interfaces.
    Marshaling a hash implementation allows its internal state to be saved and
    used for additional processing later, without having to re-write the data
    previously written to the hash. The hash state may contain portions of the
    input in its original form, which users are expected to handle for any
    possible security implications.

    Compatibility: Any future changes to hash or crypto packages will endeavor
    to maintain compatibility with state encoded using previous versions.
    That is, any released versions of the packages should be able to decode
    data written with any previously released version, subject to issues such
    as security fixes. See the Go compatibility document for background:
    https://golang.org/doc/go1compat

type Hash32 interface {
	Hash
	Sum32() uint32
}
    Hash32 is the common interface implemented by all 32-bit hash functions.

type Hash64 interface {
	Hash
	Sum64() uint64
}
    Hash64 is the common interface implemented by all 64-bit hash functions.

type XOF interface {
	// Write absorbs more data into the XOF's state. It panics if called
	// after Read.
	io.Writer

	// Read reads more output from the XOF. It may return io.EOF if there
	// is a limit to the XOF output length.
	io.Reader

	// Reset resets the XOF to its initial state.
	Reset()

	// BlockSize returns the XOF's underlying block size.
	// The Write method must be able to accept any amount
	// of data, but it may operate more efficiently if all writes
	// are a multiple of the block size.
	BlockSize() int
}
    XOF (extendable output function) is a hash function with arbitrary or
    unlimited output length.

## idiomatic usage

Idiomatic usage of `hash` drawn from the package's own runnable examples. Keywords: hash hash usage example idiomatic how to use  binary Marshaler.

```go
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding"
	"fmt"
	"log"
)

func main() {
	const (
		input1 = "The tunneling gopher digs downwards, "
		input2 = "unaware of what he will find."
	)

	first := sha256.New()
	first.Write([]byte(input1))

	marshaler, ok := first.(encoding.BinaryMarshaler)
	if !ok {
		log.Fatal("first does not implement encoding.BinaryMarshaler")
	}
	state, err := marshaler.MarshalBinary()
	if err != nil {
		log.Fatal("unable to marshal hash:", err)
	}

	second := sha256.New()

	unmarshaler, ok := second.(encoding.BinaryUnmarshaler)
	if !ok {
		log.Fatal("second does not implement encoding.BinaryUnmarshaler")
	}
	if err := unmarshaler.UnmarshalBinary(state); err != nil {
		log.Fatal("unable to unmarshal hash:", err)
	}

	first.Write([]byte(input2))
	second.Write([]byte(input2))

	fmt.Printf("%x\n", first.Sum(nil))
	fmt.Println(bytes.Equal(first.Sum(nil), second.Sum(nil)))
}

// Output:
// 57d51a066f3a39942649cd9a76c77e97ceab246756ff3888659e6aa5a07f4a52
// true
```
