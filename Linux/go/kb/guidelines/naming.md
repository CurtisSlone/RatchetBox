# Naming conventions

Idiomatic Go names (Effective Go; Code Review Comments). Names carry meaning; keep them MixedCaps.

- Use `MixedCaps` / `mixedCaps`, never underscores, for multiword names. Exported = leading capital.
- Package names: short, lower-case, single word; the package name qualifies its contents, so don't
  repeat it in identifiers (`chubby.File`, not `chubby.ChubbyFile`). Avoid `util`/`common`/`misc`.
- Getters omit `Get`: `obj.Owner()` / `obj.SetOwner(u)`.
- One-method interfaces are named method+"er": `Reader`, `Writer`, `Stringer`.
- Initialisms keep one case: `URL`, `ID`, `HTTP` -> `ServeHTTP`, `urlPony`, `appID` (not `ServeHttp`).
- Receiver names: short and consistent (one or two letters), never `this`/`self`.
- Local variable names: short where scope is small (`i`, `c`, `r`); more descriptive the farther the
  use is from the declaration.

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}

owner := obj.Owner()           // getter has no Get prefix
if owner != user {
	obj.SetOwner(user)
}

func (c *Client) Do(req *Request) (*Response, error) { /* ... */ }   // receiver c, not this/self

func ServeHTTP(w ResponseWriter, r *Request) {}   // initialism HTTP stays upper-case
```
