# crypto/sha1 (Go standard library)

Package sha1 implements the SHA-1 hash algorithm as defined in RFC 3174.

Import path: crypto/sha1   Toolchain: go1.26.4

package sha1 // import "crypto/sha1"

Package sha1 implements the SHA-1 hash algorithm as defined in RFC 3174.

SHA-1 is cryptographically broken and should not be used for secure
applications.

CONSTANTS

const BlockSize = 64
    The blocksize of SHA-1 in bytes.

const Size = 20
    The size of a SHA-1 checksum in bytes.


FUNCTIONS

func New() hash.Hash
    New returns a new hash.Hash computing the SHA1 checksum. The Hash
    also implements encoding.BinaryMarshaler, encoding.BinaryAppender and
    encoding.BinaryUnmarshaler to marshal and unmarshal the internal state of
    the hash.

func Sum(data []byte) [Size]byte
    Sum returns the SHA-1 checksum of the data.

## idiomatic usage

Compute a SHA-1 digest either in one shot with sha1.Sum, or incrementally via sha1.New and streaming data with io.Copy/io.WriteString. Keywords: sha1 New Sum hash digest checksum fingerprint hash.Hash Write io.Copy streaming file hashing SHA-1.

```go
// One-shot checksum.
data := []byte("This page intentionally left blank.")
fmt.Printf("% x", sha1.Sum(data))
// Output: af 06 49 23 bb f2 30 15 96 aa c4 c2 73 ba 32 17 8e bc 4a 96

// Streaming a file into the hash.
f, err := os.Open("file.txt")
if err != nil {
	log.Fatal(err)
}
defer f.Close()
h := sha1.New()
if _, err := io.Copy(h, f); err != nil {
	log.Fatal(err)
}
fmt.Printf("% x", h.Sum(nil))
```
