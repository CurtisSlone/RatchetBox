# Control structures

Idiomatic if/for/switch (Effective Go). No parentheses around conditions; opening brace on the same
line (semicolon insertion requires it).

- `if`/`switch`/`for` take an optional init statement: `if err := f(); err != nil { ... }`.
- One `for` keyword covers C-style, while-style, and infinite loops; `range` iterates slices, maps,
  strings (by rune), and channels.
- `switch` needs no `break`; cases don't fall through (use `fallthrough` explicitly). A bare `switch`
  (switch on true) is an idiomatic if/else-if chain. Cases may be comma-separated.
- A type switch dispatches on dynamic type.

```go
if err := file.Chmod(0664); err != nil {
	log.Print(err)
	return err
}

for i := 0; i < n; i++ { /* C-style */ }
for cond { /* while */ }
for { /* infinite */ }
for k, v := range m { use(k, v) }     // range; use _ to drop key or value
for _, r := range "héllo" { use(r) }  // ranges runes, not bytes

switch { // switch on true: an if/else-if chain
case x < 0:
	return "neg"
case x == 0:
	return "zero"
default:
	return "pos"
}

switch v := x.(type) { // type switch
case int:
	return v * 2
case string:
	return len(v)
default:
	return 0
}
```
