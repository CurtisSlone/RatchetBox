# Comments and doc conventions

Documentation conventions (Effective Go; Code Review Comments). Documentation is for users; doc comments
become the package's godoc.

- A doc comment is a full sentence that begins with the name being declared.
- Every exported (capitalized) name should have a doc comment.
- The package comment immediately precedes `package` with no blank line; one file states it, starting
  "Package x ...".
- Comment WHY, not WHAT the code obviously does. Clear is better than clever.

```go
// Package tarutil provides helpers for reading and writing tar archives.
package tarutil

// Request represents a request to run a command.
type Request struct {
	Args []string
}

// Encode writes the JSON encoding of req to w.
func Encode(w io.Writer, req *Request) error {
	// ...
	return nil
}
```
