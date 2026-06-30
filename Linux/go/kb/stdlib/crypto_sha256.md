# crypto/sha256 (Go standard library)

Package sha256 implements the SHA224 and SHA256 hash algorithms as defined in

Import path: crypto/sha256   Toolchain: go1.26.4

package sha256 // import "crypto/sha256"

Package sha256 implements the SHA224 and SHA256 hash algorithms as defined in
FIPS 180-4.

CONSTANTS

const BlockSize = 64
    The blocksize of SHA256 and SHA224 in bytes.

const Size = 32
    The size of a SHA256 checksum in bytes.

const Size224 = 28
    The size of a SHA224 checksum in bytes.


FUNCTIONS

func New() hash.Hash
    New returns a new hash.Hash computing the SHA256 checksum. The Hash
    also implements encoding.BinaryMarshaler, encoding.BinaryAppender and
    encoding.BinaryUnmarshaler to marshal and unmarshal the internal state of
    the hash.

func New224() hash.Hash
    New224 returns a new hash.Hash computing the SHA224 checksum. The Hash
    also implements encoding.BinaryMarshaler, encoding.BinaryAppender and
    encoding.BinaryUnmarshaler to marshal and unmarshal the internal state of
    the hash.

func Sum224(data []byte) [Size224]byte
    Sum224 returns the SHA224 checksum of the data.

func Sum256(data []byte) [Size]byte
    Sum256 returns the SHA256 checksum of the data.

## idiomatic usage

Compute a SHA-256 digest either in one shot with sha256.Sum256, or incrementally via sha256.New and streaming data with io.Copy. Keywords: sha256 Sum256 New hash digest checksum fingerprint hash.Hash Write io.Copy streaming file hashing SHA-256.

```go
// One-shot checksum.
sum := sha256.Sum256([]byte("hello world\n"))
fmt.Printf("%x", sum)
// Output: a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447

// Streaming a file into the hash.
f, err := os.Open("file.txt")
if err != nil {
	log.Fatal(err)
}
defer f.Close()
h := sha256.New()
if _, err := io.Copy(h, f); err != nil {
	log.Fatal(err)
}
fmt.Printf("%x", h.Sum(nil))
```
