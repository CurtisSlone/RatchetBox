# Pitfall: channel deadlocks and leaks

Sends/receives on an unbuffered channel block until the other side is ready. Mismatched send/receive
counts deadlock (panic "all goroutines are asleep") or leak goroutines. Builds clean; fails at runtime.

- An unbuffered send blocks until a receiver is ready - never send on an unbuffered channel from the
  same goroutine that must also receive it.
- Receive exactly as many values as are sent; or close the channel from the sender and `range` it.
- Close from the SENDER, once; never send on a closed channel (panics). Receiving from a closed channel
  returns the zero value with ok=false.

```go
// WRONG - deadlock: unbuffered send with no concurrent receiver
ch := make(chan int)
ch <- 1     // blocks forever; fatal error: all goroutines are asleep
fmt.Println(<-ch)

// RIGHT - receive from another goroutine
ch := make(chan int)
go func() { ch <- 1 }()
fmt.Println(<-ch)

// RIGHT - sender closes, receiver ranges until close
out := make(chan int)
go func() {
	defer close(out)
	for i := 0; i < 3; i++ {
		out <- i
	}
}()
for v := range out { // ends when out is closed
	use(v)
}
```
