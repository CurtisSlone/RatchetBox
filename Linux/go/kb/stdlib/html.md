# html (Go standard library)

Package html provides functions for escaping and unescaping HTML text.

Import path: html   Toolchain: go1.26.4

package html // import "html"

Package html provides functions for escaping and unescaping HTML text.

FUNCTIONS

func EscapeString(s string) string
    EscapeString escapes special characters like "<" to become "&lt;".
    It escapes only five such characters: <, >, &, ' and ".
    UnescapeString(EscapeString(s)) == s always holds, but the converse isn't
    always true.

func UnescapeString(s string) string
    UnescapeString unescapes entities like "&lt;" to become "<".
    It unescapes a larger range of entities than EscapeString escapes.
    For example, "&aacute;" unescapes to "á", as does "&#225;" and "&#xE1;".
    UnescapeString(EscapeString(s)) == s always holds, but the converse isn't
    always true.

## idiomatic usage

Escape the five special HTML characters (`<`, `>`, `&`, `'`, `"`) in plain text for safe output, and reverse the process. Keywords: EscapeString UnescapeString html escape unescape sanitize encode decode entities special characters ampersand angle brackets.

```go
import (
	"fmt"
	"html"
)

func ExampleEscapeString() {
	const s = `"Fran & Freddie's Diner" <tasty@example.com>`
	fmt.Println(html.EscapeString(s))
	// Output: &#34;Fran &amp; Freddie&#39;s Diner&#34; &lt;tasty@example.com&gt;
}

func ExampleUnescapeString() {
	const s = `&quot;Fran &amp; Freddie&#39;s Diner&quot; &lt;tasty@example.com&gt;`
	fmt.Println(html.UnescapeString(s))
	// Output: "Fran & Freddie's Diner" <tasty@example.com>
}
```
