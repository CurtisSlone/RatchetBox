# encoding/hex (Go standard library)

Package hex implements hexadecimal encoding and decoding.

Import path: encoding/hex   Toolchain: go1.26.4

package hex // import "encoding/hex"

Package hex implements hexadecimal encoding and decoding.

VARIABLES

var ErrLength = errors.New("encoding/hex: odd length hex string")
    ErrLength reports an attempt to decode an odd-length input using Decode or
    DecodeString. The stream-based Decoder returns io.ErrUnexpectedEOF instead
    of ErrLength.


FUNCTIONS

func AppendDecode(dst, src []byte) ([]byte, error)
    AppendDecode appends the hexadecimally decoded src to dst and returns the
    extended buffer. If the input is malformed, it returns the partially decoded
    src and an error.

func AppendEncode(dst, src []byte) []byte
    AppendEncode appends the hexadecimally encoded src to dst and returns the
    extended buffer.

func Decode(dst, src []byte) (int, error)
    Decode decodes src into DecodedLen(len(src)) bytes, returning the actual
    number of bytes written to dst.

    Decode expects that src contains only hexadecimal characters and that src
    has even length. If the input is malformed, Decode returns the number of
    bytes decoded before the error.

func DecodeString(s string) ([]byte, error)
    DecodeString returns the bytes represented by the hexadecimal string s.

    DecodeString expects that src contains only hexadecimal characters and that
    src has even length. If the input is malformed, DecodeString returns the
    bytes decoded before the error.

func DecodedLen(x int) int
    DecodedLen returns the length of a decoding of x source bytes. Specifically,
    it returns x / 2.

func Dump(data []byte) string
    Dump returns a string that contains a hex dump of the given data. The format
    of the hex dump matches the output of `hexdump -C` on the command line.

func Dumper(w io.Writer) io.WriteCloser
    Dumper returns a io.WriteCloser that writes a hex dump of all written data
    to w. The format of the dump matches the output of `hexdump -C` on the
    command line.

func Encode(dst, src []byte) int
    Encode encodes src into EncodedLen(len(src)) bytes of dst. As a convenience,
    it returns the number of bytes written to dst, but this value is always
    EncodedLen(len(src)). Encode implements hexadecimal encoding.

func EncodeToString(src []byte) string
    EncodeToString returns the hexadecimal encoding of src.

func EncodedLen(n int) int
    EncodedLen returns the length of an encoding of n source bytes.
    Specifically, it returns n * 2.

func NewDecoder(r io.Reader) io.Reader
    NewDecoder returns an io.Reader that decodes hexadecimal characters from r.
    NewDecoder expects that r contain only an even number of hexadecimal
    characters.

func NewEncoder(w io.Writer) io.Writer
    NewEncoder returns an io.Writer that writes lowercase hexadecimal characters
    to w.


TYPES

type InvalidByteError byte
    InvalidByteError values describe errors resulting from an invalid byte in a
    hex string.

func (e InvalidByteError) Error() string

## idiomatic usage

Encode bytes to a hexadecimal string and decode hex back to bytes, or produce a hexdump of binary data. Keywords: hex encode decode EncodeToString DecodeString Encode Decode EncodedLen DecodedLen Dump Dumper hexadecimal hexdump bytes to hex string.

```go
import (
	"encoding/hex"
	"fmt"
	"log"
)

func ExampleEncodeToString() {
	src := []byte("Hello")
	encodedStr := hex.EncodeToString(src)
	fmt.Printf("%s\n", encodedStr)
	// Output:
	// 48656c6c6f
}

func ExampleDecodeString() {
	const s = "48656c6c6f20476f7068657221"
	decoded, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", decoded)
	// Output:
	// Hello Gopher!
}

func ExampleDump() {
	content := []byte("Go is an open source programming language.")
	fmt.Printf("%s", hex.Dump(content))
	// Output:
	// 00000000  47 6f 20 69 73 20 61 6e  20 6f 70 65 6e 20 73 6f  |Go is an open so|
	// 00000010  75 72 63 65 20 70 72 6f  67 72 61 6d 6d 69 6e 67  |urce programming|
	// 00000020  20 6c 61 6e 67 75 61 67  65 2e                    | language.|
}
```
