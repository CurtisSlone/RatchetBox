# Recipe: a bounded worker pool

Process a stream of jobs concurrently with a fixed number of workers draining a buffered channel - the
shape behind a low-latency producer (enqueue and return immediately; workers do the work). Generalizes
the pulsehook dispatcher.

- A buffered `chan Job` is the queue; `N` worker goroutines `for job := range queue { ... }`.
- A `sync.WaitGroup` tracks workers; `Add(1)` before launching, `defer Done()` inside.
- Non-blocking submit: `select { case queue <- j: ; default: /* full */ }` keeps the producer fast.
- Use `atomic.Int64` (typed) for any shared counter to stay alignment-safe.
- Stop by `close(queue)` then `wg.Wait()` so queued jobs drain.

```go
package main

import (
	"sync"
	"sync/atomic"
)

type Pool struct {
	jobs chan func()
	wg   sync.WaitGroup
	done atomic.Int64
}

func NewPool(workers, buffer int) *Pool {
	p := &Pool{jobs: make(chan func(), buffer)}
	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for job := range p.jobs {
				job()
				p.done.Add(1)
			}
		}()
	}
	return p
}

func (p *Pool) Submit(job func()) bool { // non-blocking; false if the queue is full
	select {
	case p.jobs <- job:
		return true
	default:
		return false
	}
}

func (p *Pool) Stop() { close(p.jobs); p.wg.Wait() }
func (p *Pool) Done() int64 { return p.done.Load() }
```

Verify with `go test -race` - the worker pool is exactly where the race detector earns its keep.
