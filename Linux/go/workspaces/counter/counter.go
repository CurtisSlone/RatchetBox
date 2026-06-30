package main

import (
	"sync/atomic"
)

type Counter struct {
	value int64
}

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) Inc() {
	atomic.AddInt64(&c.value, 1)
}

func (c *Counter) Add(n int64) {
	atomic.AddInt64(&c.value, n)
}

func (c *Counter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}
