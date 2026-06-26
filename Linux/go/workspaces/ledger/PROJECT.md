# ledger

A composed Go module. Every unit is a file in `package main` at the module root.

## Units
- `counter.go` (composed unit)
- `main.go` (composed unit)
- `counter_test.go` (composed unit)
- `double.go` (added file)
- edited `double.go`: change Double to saturate at math.MaxInt instead of overflowing when n is very large
