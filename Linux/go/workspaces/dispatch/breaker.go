package main

// file-kw: breaker circuit stops calls failing destination failures cascade state machine closed open half driven

import (
	"sync"
	"time"
)

// kw: breaker circuit stops
type Breaker struct {
	maxFailures int
	cooldown    time.Duration

	state           string
	failureCount    int
	lastFailureTime time.Time
	mutex           sync.Mutex
}

// kw: breaker max failures cooldown time duration circuit stops
func NewBreaker(maxFailures int, cooldown time.Duration) *Breaker {
	return &Breaker{
		maxFailures: maxFailures,
		cooldown:    cooldown,
		state:       "closed",
	}
}

// kw: allow breaker circuit stops
func (b *Breaker) Allow() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.state == "open" {
		if time.Since(b.lastFailureTime) >= b.cooldown {
			b.state = "half-open"
			return true
		}
		return false
	}

	return true
}

// kw: success breaker circuit stops
func (b *Breaker) Success() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.state == "closed" {
		b.failureCount = 0
	} else if b.state == "half-open" {
		b.state = "closed"
		b.failureCount = 0
	}
}

// kw: failure breaker circuit stops
func (b *Breaker) Failure() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.failureCount++

	if b.state == "closed" {
		if b.failureCount >= b.maxFailures {
			b.state = "open"
			b.lastFailureTime = time.Now()
		}
	} else if b.state == "half-open" {
		b.state = "open"
		b.lastFailureTime = time.Now()
		b.failureCount = 0
	}
}

// kw: state breaker circuit stops
func (b *Breaker) State() string {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.state
}
