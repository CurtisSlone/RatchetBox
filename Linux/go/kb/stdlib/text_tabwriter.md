# text/tabwriter (Go standard library)

Package tabwriter implements a write filter (tabwriter.Writer) that translates

Import path: text/tabwriter   Toolchain: go1.26.4

package tabwriter // import "text/tabwriter"

Package tabwriter implements a write filter (tabwriter.Writer) that translates
tabbed columns in input into properly aligned text.

The package is using the Elastic Tabstops algorithm described at
http://nickgravgaard.com/elastictabstops/index.html.

The text/tabwriter package is frozen and is not accepting new features.

CONSTANTS

const (
	// Ignore html tags and treat entities (starting with '&'
	// and ending in ';') as single characters (width = 1).
	FilterHTML uint = 1 << iota

	// Strip Escape characters bracketing escaped text segments
	// instead of passing them through unchanged with the text.
	StripEscape

	// Force right-alignment of cell content.
	// Default is left-alignment.
	AlignRight

	// Handle empty columns as if they were not present in
	// the input in the first place.
	DiscardEmptyColumns

	// Always use tabs for indentation columns (i.e., padding of
	// leading empty cells on the left) independent of padchar.
	TabIndent

	// Print a vertical bar ('|') between columns (after formatting).
	// Discarded columns appear as zero-width columns ("||").
	Debug
)
    Formatting can be controlled with these flags.

const Escape = '\xff'
    To escape a text segment, bracket it with Escape characters. For instance,
    the tab in this string "Ignore this tab: \xff\t\xff" does not terminate
    a cell and constitutes a single character of width one for formatting
    purposes.

    The value 0xff was chosen because it cannot appear in a valid UTF-8
    sequence.


TYPES

type Writer struct {
	// Has unexported fields.
}
    A Writer is a filter that inserts padding around tab-delimited columns in
    its input to align them in the output.

    The Writer treats incoming bytes as UTF-8-encoded text consisting of cells
    terminated by horizontal ('\t') or vertical ('\v') tabs, and newline ('\n')
    or formfeed ('\f') characters; both newline and formfeed act as line breaks.

    Tab-terminated cells in contiguous lines constitute a column. The Writer
    inserts padding as needed to make all cells in a column have the same width,
    effectively aligning the columns. It assumes that all characters have
    the same width, except for tabs for which a tabwidth must be specified.
    Column cells must be tab-terminated, not tab-separated: non-tab terminated
    trailing text at the end of a line forms a cell but that cell is not part
    of an aligned column. For instance, in this example (where | stands for a
    horizontal tab):

        aaaa|bbb|d
        aa  |b  |dd
        a   |
        aa  |cccc|eee

    the b and c are in distinct columns (the b column is not contiguous all the
    way). The d and e are not in a column at all (there's no terminating tab,
    nor would the column be contiguous).

    The Writer assumes that all Unicode code points have the same width;
    this may not be true in some fonts or if the string contains combining
    characters.

    If DiscardEmptyColumns is set, empty columns that are terminated entirely
    by vertical (or "soft") tabs are discarded. Columns terminated by horizontal
    (or "hard") tabs are not affected by this flag.

    If a Writer is configured to filter HTML, HTML tags and entities are passed
    through. The widths of tags and entities are assumed to be zero (tags) and
    one (entities) for formatting purposes.

    A segment of text may be escaped by bracketing it with Escape characters.
    The tabwriter passes escaped text segments through unchanged. In particular,
    it does not interpret any tabs or line breaks within the segment. If the
    StripEscape flag is set, the Escape characters are stripped from the output;
    otherwise they are passed through as well. For the purpose of formatting,
    the width of the escaped text is always computed excluding the Escape
    characters.

    The formfeed character acts like a newline but it also terminates all
    columns in the current line (effectively calling Writer.Flush). Tab-
    terminated cells in the next line start new columns. Unless found inside an
    HTML tag or inside an escaped text segment, formfeed characters appear as
    newlines in the output.

    The Writer must buffer input internally, because proper spacing of one line
    may depend on the cells in future lines. Clients must call Flush when done
    calling Writer.Write.

func NewWriter(output io.Writer, minwidth, tabwidth, padding int, padchar byte, flags uint) *Writer
    NewWriter allocates and initializes a new Writer. The parameters are the
    same as for the Init function.

func (b *Writer) Flush() error
    Flush should be called after the last call to Writer.Write to ensure that
    any data buffered in the Writer is written to output. Any incomplete escape
    sequence at the end is considered complete for formatting purposes.

func (b *Writer) Init(output io.Writer, minwidth, tabwidth, padding int, padchar byte, flags uint) *Writer
    A Writer must be initialized with a call to Init. The first parameter
    (output) specifies the filter output. The remaining parameters control the
    formatting:

        minwidth	minimal cell width including any padding
        tabwidth	width of tab characters (equivalent number of spaces)
        padding		padding added to a cell before computing its width
        padchar		ASCII char used for padding
        		if padchar == '\t', the Writer will assume that the
        		width of a '\t' in the formatted output is tabwidth,
        		and cells are left-aligned independent of align_left
        		(for correct-looking results, tabwidth must correspond
        		to the tab width in the viewer displaying the result)
        flags		formatting control

func (b *Writer) Write(buf []byte) (n int, err error)
    Write writes buf to the writer b. The only errors returned are ones
    encountered while writing to the underlying output stream.

## idiomatic usage

Idiomatic usage of `text/tabwriter` drawn from the package's own runnable examples. Keywords: text/tabwriter tabwriter usage example idiomatic how to use Writer Init  elastic  trailing Tab.

```go
package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func main() {
	w := new(tabwriter.Writer)

	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "a\tb\tc\td\t.")
	fmt.Fprintln(w, "123\t12345\t1234567\t123456789\t.")
	fmt.Fprintln(w)
	w.Flush()

	// Format right-aligned in space-separated columns of minimal width 5
	// and at least one blank of padding (so wider column entries do not
	// touch each other).
	w.Init(os.Stdout, 5, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "a\tb\tc\td\t.")
	fmt.Fprintln(w, "123\t12345\t1234567\t123456789\t.")
	fmt.Fprintln(w)
	w.Flush()

}

// Output:
// a	b	c	d		.
// 123	12345	1234567	123456789	.
// 
//     a     b       c         d.
//   123 12345 1234567 123456789.
```

```go
package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func main() {
	// Observe how the b's and the d's, despite appearing in the
	// second cell of each line, belong to different columns.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "a\tb\tc")
	fmt.Fprintln(w, "aa\tbb\tcc")
	fmt.Fprintln(w, "aaa\t") // trailing tab
	fmt.Fprintln(w, "aaaa\tdddd\teeee")
	w.Flush()

}

// Output:
// ....a|..b|c
// ...aa|.bb|cc
// ..aaa|
// .aaaa|.dddd|eeee
```

```go
package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func main() {
	// Observe that the third line has no trailing tab,
	// so its final cell is not part of an aligned column.
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, '-', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "a\tb\taligned\t")
	fmt.Fprintln(w, "aa\tbb\taligned\t")
	fmt.Fprintln(w, "aaa\tbbb\tunaligned") // no trailing tab
	fmt.Fprintln(w, "aaaa\tbbbb\taligned\t")
	w.Flush()

}

// Output:
// ------a|------b|---aligned|
// -----aa|-----bb|---aligned|
// ----aaa|----bbb|unaligned
// ---aaaa|---bbbb|---aligned|
```
