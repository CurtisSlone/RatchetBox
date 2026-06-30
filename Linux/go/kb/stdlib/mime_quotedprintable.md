# mime/quotedprintable (Go standard library)

Package quotedprintable implements quoted-printable encoding as specified by RFC

Import path: mime/quotedprintable   Toolchain: go1.26.4

package quotedprintable // import "mime/quotedprintable"

Package quotedprintable implements quoted-printable encoding as specified by RFC
2045.

TYPES

type Reader struct {
	// Has unexported fields.
}
    Reader is a quoted-printable decoder.

func NewReader(r io.Reader) *Reader
    NewReader returns a quoted-printable reader, decoding from r.

func (r *Reader) Read(p []byte) (n int, err error)
    Read reads and decodes quoted-printable data from the underlying reader.

type Writer struct {
	// Binary mode treats the writer's input as pure binary and processes end of
	// line bytes as binary data.
	Binary bool

	// Has unexported fields.
}
    A Writer is a quoted-printable writer that implements io.WriteCloser.

func NewWriter(w io.Writer) *Writer
    NewWriter returns a new Writer that writes to w.

func (w *Writer) Close() error
    Close closes the Writer, flushing any unwritten data to the underlying
    io.Writer, but does not close the underlying io.Writer.

func (w *Writer) Write(p []byte) (n int, err error)
    Write encodes p using quoted-printable encoding and writes it to the
    underlying io.Writer. It limits line length to 76 characters. The encoded
    bytes are not necessarily flushed until the Writer is closed.

## idiomatic usage

Decode and encode quoted-printable content (used in MIME email bodies) via NewReader and NewWriter. Keywords: quotedprintable.NewReader quotedprintable.NewWriter Write Close decode encode quoted-printable QP MIME email transfer encoding =3D.

```go
import (
	"fmt"
	"io"
	"mime/quotedprintable"
	"os"
	"strings"
)

func main() {
	// Decode quoted-printable input.
	r := quotedprintable.NewReader(strings.NewReader("=48=65=6C=6C=6F=2C=20=47=6F=70=68=65=72=73=21"))
	b, _ := io.ReadAll(r)
	fmt.Printf("%s\n", b) // Hello, Gophers!

	// Encode to quoted-printable.
	w := quotedprintable.NewWriter(os.Stdout)
	w.Write([]byte("These symbols will be escaped: = \t"))
	w.Close()
	// These symbols will be escaped: =3D =09
}
```
