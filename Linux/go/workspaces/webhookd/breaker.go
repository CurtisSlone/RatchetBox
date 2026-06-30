package main

// file-kw: breaker circuit stops calls failing destination failures cascade state machine closed open half driven

import (
	"sync"
	"time"
)

// Breaker is a circuit breaker that stops calls to a failing destination
// so failures do not cascade. It is a three-state machine (closed, open, half-open)
// driven by recorded successes and failures. Safe for concurrent use.
// kw: breaker circuit stops
type Breaker struct {
	maxFailures int
	cooldown    time.Duration

	state          string
	failureCount   int
	lastOpenedTime time.Time
	mu             sync.Mutex
}

// NewBreaker creates a new breaker with the given maximum failures and cooldown.
// kw: breaker max failures cooldown time duration circuit stops
func NewBreaker(maxFailures int, cooldown time.Duration) *Breaker {
	return &Breaker{
		maxFailures:    maxFailures,
		cooldown:       cooldown,
		state:          "closed",
		failureCount:   0,
		lastOpenedTime: time.Time{},
	}
}

// Allow reports whether a call may proceed, after applying the time transition.
// kw: allow breaker circuit stops
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.applyTimeTransition()
	switch b.state {
	case "closed":
		return true
	case "open":
		return false
	case "half-open":
		return true
	default:
		return false
	}
}

// Success records a successful call.
// kw: success breaker circuit stops
func (b *Breaker) Success() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.applyTimeTransition()
	switch b.state {
	case "closed":
		b.failureCount = 0
	case "half-open":
		b.state = "closed"
		b.failureCount = 0
	}
}

// Failure records a failed call.
// kw: failure breaker circuit stops
func (b *Breaker) Failure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.applyTimeTransition()
	switch b.state {
	case "closed":
		b.failureCount++
		if b.failureCount >= b.maxFailures {
			b.state = "open"
			b.lastOpenedTime = time.Now()
		}
	case "half-open":
		b.state = "open"
		b.lastOpenedTime = time.Now()
	}
}

// State returns the current stored state string.
// kw: state breaker circuit stops
func (b *Breaker) State() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// applyTimeTransition applies the time transition if needed.
// kw: apply time transition breaker circuit stops
func (b *Breaker) applyTimeTransition() {
	if b.state == "open" && time.Since(b.lastOpenedTime) >= b.cooldown {
		b.state = "half-open"
	}
}
