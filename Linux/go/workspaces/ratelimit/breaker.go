package main

// file-kw: circuit breaker state machine concurrent
import (
	"sync"
	"time"
)

// kw: create breaker struct
type Breaker struct {
	maxFailures  int
	cooldown     time.Duration
	state        string
	failureCount int
	lastFailure  time.Time
	mutex        sync.Mutex
}

// kw: create new breaker
func NewBreaker(maxFailures int, cooldown time.Duration) *Breaker {
	return &Breaker{
		maxFailures:  maxFailures,
		cooldown:     cooldown,
		state:        "closed",
		failureCount: 0,
		lastFailure:  time.Time{},
	}
}

// kw: apply time transition
func (b *Breaker) applyTimeTransition() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.state == "open" && time.Since(b.lastFailure) >= b.cooldown {
		b.state = "half-open"
	}
}

// kw: check if call allowed
func (b *Breaker) Allow() bool {
	b.applyTimeTransition()

	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.state == "closed" || b.state == "half-open"
}

// kw: record success
func (b *Breaker) Success() {
	b.applyTimeTransition()

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.state == "closed" {
		b.failureCount = 0
	} else if b.state == "half-open" {
		b.state = "closed"
		b.failureCount = 0
	}
}

// kw: record failure
func (b *Breaker) Failure() {
	b.applyTimeTransition()

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.state == "closed" {
		b.failureCount++
		if b.failureCount >= b.maxFailures {
			b.state = "open"
			b.lastFailure = time.Now()
		}
	} else if b.state == "half-open" {
		b.state = "open"
		b.lastFailure = time.Now()
		b.failureCount = 0
	}
}

// kw: get current state
func (b *Breaker) State() string {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.state
}
