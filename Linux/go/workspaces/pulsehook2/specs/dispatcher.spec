name: Dispatcher
role: component
intent: a worker pool draining a buffered channel so the handler returns immediately (low latency)
api:
  - type Dispatcher struct with: a buffered chan Event; a sync.WaitGroup; a worker count; and a
    processed counter that is an atomic.Int64 (use the typed atomic to avoid 64-bit alignment issues)
  - func NewDispatcher(workers, buffer int) *Dispatcher
  - method (*Dispatcher) Start()   // launch `workers` goroutines, each: defer wg.Done(); for range queue { d.processed.Add(1) }
  - method (*Dispatcher) Enqueue(e Event) bool   // NON-BLOCKING: select { case queue <- e: true; default: false }
  - method (*Dispatcher) Stop()   // close(queue); wg.Wait()
  - method (*Dispatcher) Processed() int64   // return d.processed.Load()
constraints: standard library (sync, sync/atomic); package main; pointer receivers; use atomic.Int64 (NOT a raw int64 field)
