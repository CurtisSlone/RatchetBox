# crypto/des (Go standard library)

Package des implements the Data Encryption Standard (DES) and the Triple Data

Import path: crypto/des   Toolchain: go1.26.4

package des // import "crypto/des"

Package des implements the Data Encryption Standard (DES) and the Triple Data
Encryption Algorithm (TDEA) as defined in U.S. Federal Information Processing
Standards Publication 46-3.

DES is cryptographically broken and should not be used for secure applications.

CONSTANTS

const BlockSize = 8
    The DES block size in bytes.


FUNCTIONS

func NewCipher(key []byte) (cipher.Block, error)
    NewCipher creates and returns a new cipher.Block.

func NewTripleDESCipher(key []byte) (cipher.Block, error)
    NewTripleDESCipher creates and returns a new cipher.Block.


TYPES

type KeySizeError int

func (k KeySizeError) Error() string

## idiomatic usage

Construct a Triple DES (3DES) block cipher with des.NewTripleDESCipher (24-byte key; build an EDE2 key by repeating the first 8 bytes), then use it as a cipher.Block with crypto/cipher modes. Keywords: des.NewTripleDESCipher des.NewCipher Triple DES 3DES EDE2 EDE3 cipher.Block key 24 bytes encrypt decrypt KeySizeError.

```go
import "crypto/des"

// EDE2: duplicate the first 8 bytes of a 16-byte key to make a 24-byte key.
ede2Key := []byte("example key 1234")
var tripleDESKey []byte
tripleDESKey = append(tripleDESKey, ede2Key[:16]...)
tripleDESKey = append(tripleDESKey, ede2Key[:8]...)

block, err := des.NewTripleDESCipher(tripleDESKey)
if err != nil {
	panic(err)
}
_ = block // use as a cipher.Block with crypto/cipher modes (e.g. CBC, CTR)
```
