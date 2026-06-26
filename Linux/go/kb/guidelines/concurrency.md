# Concurrency idioms

Goroutines and channels (Effective Go; proverbs). Don't communicate by sharing memory; share memory by
communicating. Concurrency is not parallelism.

- `go f()` starts a goroutine. Synchronize with channels or `sync.WaitGroup`, not sleeps.
- `make(chan T)` is unbuffered (send blocks until receive); `make(chan T, n)` buffers n.
- Close a channel from the SENDER when no more values will be sent; `for v := range ch` drains until
  close. Receiving from a closed channel yields the zero value with ok=false.
- A buffered channel can act as a semaphore to bound concurrency.
- Protect shared mutable state with `sync.Mutex` when communication doesn't fit.

```go
// Wait for N goroutines with a WaitGroup.
var wg sync.WaitGroup
for _, item := range items {
	wg.Add(1)
	go func(it Item) { // pass the loop var as an arg (avoid capture bugs)
		defer wg.Done()
		process(it)
	}(item)
}
wg.Wait()

// Fan-in results over a channel.
results := make(chan int)
for _, n := range nums {
	go func(n int) { results <- n * n }(n)
}
sum := 0
for range nums {
	sum += <-results
}

// Bounded concurrency with a buffered channel as a semaphore.
sem := make(chan struct{}, maxConcurrent)
for _, r := range reqs {
	sem <- struct{}{}
	go func(r Req) { defer func() { <-sem }(); handle(r) }(r)
}
```
