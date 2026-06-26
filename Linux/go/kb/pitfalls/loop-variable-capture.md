# Pitfall: capturing the loop variable in a goroutine/closure

Closures capture a variable by reference, not by value. A goroutine or closure that refers to the loop
variable may see its final value, not the value at the iteration it was created. (Go 1.22+ gives each
iteration a fresh loop variable, but write code that is correct regardless and clear to readers.)

- Pass the loop variable as an argument to the goroutine/closure, or rebind it inside the loop.
- This bug builds and vets clean; it shows up as wrong results under `go test` / the race detector.

```go
// RISKY - all goroutines may observe the same (final) v on older toolchains
for _, v := range items {
	go func() { process(v) }() // captures v by reference
}

// SAFE - pass as an argument (a fresh copy per call)
for _, v := range items {
	go func(v Item) { process(v) }(v)
}

// SAFE - rebind inside the loop
for _, v := range items {
	v := v // new variable each iteration
	go func() { process(v) }()
}
```
