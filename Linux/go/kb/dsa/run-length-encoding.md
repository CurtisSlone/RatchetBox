# Run-Length Encoding (RLE), Unambiguous

Run-length encoding compresses data by replacing each run of identical symbols with a (count, symbol) pair, which is effective when the input has long repeats. The NAIVE scheme of writing the decimal count as a prefix ("3A") is AMBIGUOUS for two reasons: (1) if the data itself contains digits, "3A" could be a run of three 'A's or the literal text "3A"; and (2) a run longer than 9 like "12W" is indistinguishable from "1" followed by "2W". The fix shown here is to length-PREFIX every run with a fixed-width byte count and emit the raw symbol byte, so the decoder always knows exactly how many bytes a token spans regardless of the data. Runs longer than 255 are split into multiple tokens. Encoding and decoding are both O(n) time. Keywords: run length encoding RLE run-length decode compression count prefix ambiguous unambiguous length prefix escape digits long runs greater than 9 round trip repeat symbol token byte compress decompress

## implementation

```go
package compression

// RLE uses a fixed-format token of two bytes: [count][symbol], where count is
// 1..255. This is unambiguous for ANY input -- including bytes that look like
// ASCII digits and runs longer than 9 -- because the decoder reads a fixed
// 1-byte count and then exactly one symbol byte, never parsing decimal text.
// Runs longer than 255 are emitted as several consecutive tokens.

// RLEncode compresses data into a sequence of [count][symbol] byte pairs.
func RLEncode(data []byte) []byte {
	var out []byte
	for i := 0; i < len(data); {
		sym := data[i]
		run := 1
		for i+run < len(data) && data[i+run] == sym && run < 255 {
			run++
		}
		out = append(out, byte(run), sym)
		i += run
	}
	return out
}

// RLEdecode reverses RLEncode. It expects an even-length stream of
// [count][symbol] pairs and returns the original bytes.
func RLEdecode(data []byte) []byte {
	var out []byte
	for i := 0; i+1 < len(data); i += 2 {
		count := int(data[i])
		sym := data[i+1]
		for j := 0; j < count; j++ {
			out = append(out, sym)
		}
	}
	return out
}
```

## usage / test

```go
package compression

import (
	"bytes"
	"strings"
	"testing"
)

func TestRLERoundTrip(t *testing.T) {
	cases := [][]byte{
		nil,
		[]byte(""),
		[]byte("A"),
		[]byte("AAAAA"),
		[]byte("WWWWWWWWWWWWBWWWWWWWWWWWWBBB"),
		[]byte("3A"),                          // input contains digits
		[]byte("12W"),                         // ambiguous under naive count-prefix RLE
		[]byte("aaaaaaaaaa1111111111"),        // digits in long runs
		[]byte(strings.Repeat("Z", 1000)),     // run far longer than 255
		[]byte{0, 0, 0, 255, 255, 1, 2, 3, 3}, // arbitrary bytes
	}
	for _, in := range cases {
		got := RLEdecode(RLEncode(in))
		if len(in) == 0 && len(got) == 0 {
			continue
		}
		if !bytes.Equal(got, in) {
			t.Errorf("round trip failed for %q: got %q", in, got)
		}
	}
}
```
