# image/png (Go standard library)

Package png implements a PNG image decoder and encoder.

Import path: image/png   Toolchain: go1.26.4

package png // import "image/png"

Package png implements a PNG image decoder and encoder.

The PNG specification is at https://www.w3.org/TR/PNG/.

FUNCTIONS

func Decode(r io.Reader) (image.Image, error)
    Decode reads a PNG image from r and returns it as an image.Image. The type
    of Image returned depends on the PNG contents.

func DecodeConfig(r io.Reader) (image.Config, error)
    DecodeConfig returns the color model and dimensions of a PNG image without
    decoding the entire image.

func Encode(w io.Writer, m image.Image) error
    Encode writes the Image m to w in PNG format. Any Image may be encoded,
    but images that are not image.NRGBA might be encoded lossily.


TYPES

type CompressionLevel int
    CompressionLevel indicates the compression level.

const (
	DefaultCompression CompressionLevel = 0
	NoCompression      CompressionLevel = -1
	BestSpeed          CompressionLevel = -2
	BestCompression    CompressionLevel = -3
)
type Encoder struct {
	CompressionLevel CompressionLevel

	// BufferPool optionally specifies a buffer pool to get temporary
	// EncoderBuffers when encoding an image.
	BufferPool EncoderBufferPool
}
    Encoder configures encoding PNG images.

func (enc *Encoder) Encode(w io.Writer, m image.Image) error
    Encode writes the Image m to w in PNG format.

type EncoderBuffer encoder
    EncoderBuffer holds the buffers used for encoding PNG images.

type EncoderBufferPool interface {
	Get() *EncoderBuffer
	Put(*EncoderBuffer)
}
    EncoderBufferPool is an interface for getting and returning temporary
    instances of the EncoderBuffer struct. This can be used to reuse buffers
    when encoding multiple images.

type FormatError string
    A FormatError reports that the input is not a valid PNG.

func (e FormatError) Error() string

type UnsupportedError string
    An UnsupportedError reports that the input uses a valid but unimplemented
    PNG feature.

func (e UnsupportedError) Error() string

## idiomatic usage

Encode an in-memory image to a PNG file, or decode a PNG stream into an `image.Image`. Keywords: png Encode Decode write read save load PNG image NewNRGBA color os.Create Bounds At pixels.

```go
import (
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
)

func ExampleEncode() {
	const width, height = 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8((x + y) & 255),
				G: uint8((x + y) << 1 & 255),
				B: uint8((x + y) << 2 & 255),
				A: 255,
			})
		}
	}

	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func decode(r io.Reader) {
	img, err := png.Decode(r)
	if err != nil {
		log.Fatal(err)
	}
	_ = img.Bounds()
}
```
