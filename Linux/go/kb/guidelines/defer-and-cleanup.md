# Defer and cleanup

Using `defer` for cleanup (Effective Go). A deferred call runs when the surrounding function returns.

- Defer the cleanup right after acquiring the resource, so it is never forgotten on any return path.
- Deferred calls run in LIFO order. Arguments are evaluated when `defer` executes, not when it runs.
- A deferred closure can read and modify named return values (useful for wrapping errors).
- Beware deferring inside a loop: the calls pile up until the function returns, not each iteration -
  for per-iteration cleanup, wrap the body in its own function.

```go
func Contents(name string) (string, error) {
	f, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer f.Close() // runs on every return below
	// ... read f ...
	return string(b), nil
}

// LIFO + argument capture at defer time:
for i := 0; i < 3; i++ {
	defer fmt.Print(i) // prints 2 1 0 when the function returns
}

// Per-iteration cleanup: defer inside a helper, not the loop body.
for _, name := range names {
	func() {
		f, err := os.Open(name)
		if err != nil {
			return
		}
		defer f.Close() // closes at the end of THIS iteration
		use(f)
	}()
}
```
