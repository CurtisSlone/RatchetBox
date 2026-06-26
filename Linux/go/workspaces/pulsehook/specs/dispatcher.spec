name: Dispatcher
role: component
intent: a worker pool that processes Events asynchronously off a buffered channel, so the HTTP handler
  can accept a webhook and return immediately (low latency) instead of doing the work inline
api:
  - type Dispatcher struct holding a buffered chan Event (the queue), a sync.WaitGroup, a worker count,
    and an int64 processed counter updated with sync/atomic
  - func NewDispatcher(workers, buffer int) *Dispatcher   // makes the channel with capacity buffer
  - method (*Dispatcher) Start()   // launches `workers` goroutines, each draining the queue with
    `for e := range d.queue { ... atomic.AddInt64(&d.processed, 1) }`
  - method (*Dispatcher) Enqueue(e Event) bool   // NON-BLOCKING: select { case d.queue <- e: return true; default: return false }
    so a full queue never blocks the caller (this is the low-latency guarantee)
  - method (*Dispatcher) Stop()   // close(d.queue) then d.wg.Wait() to drain workers cleanly
  - method (*Dispatcher) Processed() int64   // atomic.LoadInt64 of the processed counter
behavior:
  - each worker goroutine must call wg.Done() when the range loop ends (defer wg.Done() at the top)
  - Start must call wg.Add(1) per worker before launching it
constraints: standard library only (sync, sync/atomic); package main; no func main in this file;
  pointer receivers (it holds a sync.WaitGroup and must not be copied)
