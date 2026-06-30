# encoding/pem (Go standard library)

Package pem implements the PEM data encoding, which originated in Privacy

Import path: encoding/pem   Toolchain: go1.26.4

package pem // import "encoding/pem"

Package pem implements the PEM data encoding, which originated in Privacy
Enhanced Mail. The most common use of PEM encoding today is in TLS keys and
certificates. See RFC 1421.

FUNCTIONS

func Encode(out io.Writer, b *Block) error
    Encode writes the PEM encoding of b to out.

func EncodeToMemory(b *Block) []byte
    EncodeToMemory returns the PEM encoding of b.

    If b has invalid headers and cannot be encoded, EncodeToMemory returns nil.
    If it is important to report details about this error case, use Encode
    instead.


TYPES

type Block struct {
	Type    string            // The type, taken from the preamble (i.e. "RSA PRIVATE KEY").
	Headers map[string]string // Optional headers.
	Bytes   []byte            // The decoded bytes of the contents. Typically a DER encoded ASN.1 structure.
}
    A Block represents a PEM encoded structure.

    The encoded form is:

        -----BEGIN Type-----
        Headers
        base64-encoded Bytes
        -----END Type-----

    where Block.Headers is a possibly empty sequence of Key: Value lines.

func Decode(data []byte) (p *Block, rest []byte)
    Decode will find the next PEM formatted block (certificate, private key etc)
    in the input. It returns that block and the remainder of the input. If no
    PEM data is found, p is nil and the whole of the input is returned in rest.
    Blocks must start at the beginning of a line and end at the end of a line.

## idiomatic usage

Decode a PEM-encoded block (such as a public key or certificate) into its DER bytes, and encode a `pem.Block` back to PEM text. Keywords: pem Decode Encode pem.Block Type Bytes Headers PEM certificate public key BEGIN END decode encode parse PEM block x509.

```go
import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func ExampleDecode() {
	var pubPEMData = []byte(`
-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAlRuRnThUjU8/prwYxbty
...
-----END PUBLIC KEY-----
and some more`)

	block, rest := pem.Decode(pubPEMData)
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatal("failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Got a %T, with remaining data: %q", pub, rest)
	// Output: Got a *rsa.PublicKey, with remaining data: "and some more"
}

func ExampleEncode() {
	block := &pem.Block{
		Type: "MESSAGE",
		Headers: map[string]string{
			"Animal": "Gopher",
		},
		Bytes: []byte("test"),
	}
	if err := pem.Encode(os.Stdout, block); err != nil {
		log.Fatal(err)
	}
	// Output:
	// -----BEGIN MESSAGE-----
	// Animal: Gopher
	//
	// dGVzdA==
	// -----END MESSAGE-----
}
```
