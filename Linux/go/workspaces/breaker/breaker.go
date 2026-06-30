package main

import (
	"sync"
	"time"
)

type Breaker struct {
	maxFailures  int
	cooldown     time.Duration
	state        string
	failureCount int
	lastOpenTime time.Time
	mutex        sync.Mutex
}

func NewBreaker(maxFailures int, cooldown time.Duration) *Breaker {
	return &Breaker{
		maxFailures:  maxFailures,
		cooldown:     cooldown,
		state:        "closed",
		failureCount: 0,
	}
}

func (b *Breaker) transitionTime() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.state == "open" {
		if time.Since(b.lastOpenTime) >= b.cooldown {
			b.state = "half-open"
		}
	}
}

func (b *Breaker) Allow() bool {
	b.transitionTime()
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.state == "closed" || b.state == "half-open"
}

func (b *Breaker) Success() {
	b.transitionTime()
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.state == "closed" {
		b.failureCount = 0
	} else if b.state == "half-open" {
		b.state = "closed"
		b.failureCount = 0
	}
}

func (b *Breaker) Failure() {
	b.transitionTime()
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.state == "closed" {
		b.failureCount++
		if b.failureCount >= b.maxFailures {
			b.state = "open"
			b.lastOpenTime = time.Now()
		}
	} else if b.state == "half-open" {
		b.state = "open"
		b.lastOpenTime = time.Now()
		b.failureCount = 0
	}
}

func (b *Breaker) State() string {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.state
}
