# Go idioms for buildable library code

Guidance for emitting a single Go file that type-checks under `go build` as a library package.

## Package and imports
- Use a library package clause such as `package solution`. Do not write `func main` - a library
  package does not need an entry point, and it lets the oracle type-check pure functions.
- Import only what you use. Go treats an unused import as a compile error: `imported and not used`.
- Group standard-library imports in one `import ( ... )` block.

## Common standard-library packages
- `strings` - `strings.Builder`, `Split`, `Join`, `TrimSpace`, `Contains`, `ReplaceAll`, `HasPrefix`.
- `strconv` - `Atoi`, `Itoa`, `ParseFloat`, `FormatInt`.
- `sort` - `sort.Ints`, `sort.Strings`, `sort.Slice(s, func(i, j int) bool { ... })`.
- `errors` / `fmt` - `errors.New`, `fmt.Errorf("...: %w", err)` for wrapped errors.
- `unicode` - `unicode.Isletter/IsDigit/IsSpace`.

## Errors and zero values
- Return errors as values: `func f() (T, error)`. The zero value of T plus a non-nil error on failure.
- Check `if err != nil { return zero, err }`. Do not panic for ordinary failures.

## Unused variables
- A declared-but-unused local is a compile error. Use `_ = x` only as a deliberate discard, never to
  paper over dead code - prefer removing it.

## Slices and maps
- `make([]T, 0, n)` to preallocate; append with `s = append(s, v)`.
- `m := map[K]V{}`; comma-ok read: `v, ok := m[k]`.

## Concurrency (only if asked)
- A goroutine is `go f()`; synchronize with channels or `sync.WaitGroup`. Do not add concurrency
  unless the task calls for it.
