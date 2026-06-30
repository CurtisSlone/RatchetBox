# crypto/md5 (Go standard library)

Package md5 implements the MD5 hash algorithm as defined in RFC 1321.

Import path: crypto/md5   Toolchain: go1.26.4

package md5 // import "crypto/md5"

Package md5 implements the MD5 hash algorithm as defined in RFC 1321.

MD5 is cryptographically broken and should not be used for secure applications.

CONSTANTS

const BlockSize = 64
    The blocksize of MD5 in bytes.

const Size = 16
    The size of an MD5 checksum in bytes.


FUNCTIONS

func New() hash.Hash
    New returns a new hash.Hash computing the MD5 checksum. The Hash
    also implements encoding.BinaryMarshaler, encoding.BinaryAppender and
    encoding.BinaryUnmarshaler to marshal and unmarshal the internal state of
    the hash.

func Sum(data []byte) [Size]byte
    Sum returns the MD5 checksum of the data.

## idiomatic usage

Compute an MD5 checksum either in one shot with Sum or by streaming data through a hash.Hash via New (e.g. hashing a file with io.Copy). Keywords: New Sum md5.New md5.Sum hash checksum digest io.WriteString io.Copy h.Sum file hash MD5 message digest fingerprint.

```go
import (
	"crypto/md5"
	"fmt"
	"io"
)

// One-shot checksum of a byte slice.
func ExampleSum() {
	data := []byte("These pretzels are making me thirsty.")
	fmt.Printf("%x", md5.Sum(data))
	// Output: b0804ec967f48520697662a204f5fe72
}

// Streaming a hash with New.
func ExampleNew() {
	h := md5.New()
	io.WriteString(h, "The fog is getting thicker!")
	io.WriteString(h, "And Leon's getting laaarger!")
	fmt.Printf("%x", h.Sum(nil))
	// Output: e2c569be17396eca2a2e3c11578123ed
}
```
