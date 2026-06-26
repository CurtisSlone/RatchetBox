package main

import (
	"sync"
	"sync/atomic"
)

type Dispatcher struct {
	queue     chan Event
	wg        sync.WaitGroup
	workers   int
	processed atomic.Int64
}

func NewDispatcher(workers, buffer int) *Dispatcher {
	return &Dispatcher{
		queue:   make(chan Event, buffer),
		workers: workers,
	}
}

func (d *Dispatcher) Start() {
	for i := 0; i < d.workers; i++ {
		d.wg.Add(1)
		go func() {
			defer d.wg.Done()
			for range d.queue {
				d.processed.Add(1)
			}
		}()
	}
}

func (d *Dispatcher) Enqueue(e Event) bool {
	select {
	case d.queue <- e:
		return true
	default:
		return false
	}
}

func (d *Dispatcher) Stop() {
	close(d.queue)
	d.wg.Wait()
}

func (d *Dispatcher) Processed() int64 {
	return d.processed.Load()
}
