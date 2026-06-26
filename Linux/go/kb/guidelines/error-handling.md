# Error handling idioms

How idiomatic Go reports and handles errors (Effective Go; Code Review Comments). Errors are values.

- Return errors as the last result value; check and handle them, do not panic for ordinary failures.
- Error strings are not capitalized and do not end with punctuation (they are often wrapped).
- Indent the error-handling branch and keep the happy path at minimal indentation; return early.
- Wrap with `%w` to preserve the chain; inspect with `errors.Is` / `errors.As`.
- Prefer multiple return values over in-band sentinels (e.g. -1, "", nil) for "not found".

```go
// Error strings: lower-case, no trailing punctuation.
return fmt.Errorf("parse config %q: %w", path, err)   // not "Parse config..." and not "...: %w."

// Indent error flow; keep the normal path un-indented; return early instead of else.
f, err := os.Open(name)
if err != nil {
	return err
}
defer f.Close()
codeUsing(f)

// In-band errors: prefer (value, ok) or (value, error) over a sentinel return.
func Lookup(key string) (value string, ok bool)   // good
// func Lookup(key string) string                 // avoid: caller can't tell "missing" from ""

// A custom error type carries structured context and still satisfies error.
type PathError struct {
	Op   string
	Path string
	Err  error
}

func (e *PathError) Error() string { return e.Op + " " + e.Path + ": " + e.Err.Error() }

// Inspect a wrapped error by type.
var perr *PathError
if errors.As(err, &perr) && perr.Err == syscall.ENOSPC {
	deleteTempFiles()
}
```
