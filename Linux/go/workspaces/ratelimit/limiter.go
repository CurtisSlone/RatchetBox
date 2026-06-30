package main

// file-kw: limiter rate abstraction decides whether client request right now under sliding window policy code

import (
	"context"
	"sync"
	"time"
)

// kw: limiter rate abstraction
type Limiter interface {
	Allow(ctx context.Context, client string) (bool, error)
}

// kw: sliding window limiter rate abstraction
type slidingWindowLimiter struct {
	windowSize  time.Duration
	maxRequests int64
	requests    map[string][]time.Time
	mutex       sync.RWMutex
}

// kw: sliding window limiter size time duration max requests rate
func NewSlidingWindowLimiter(windowSize time.Duration, maxRequests int64) *slidingWindowLimiter {
	return &slidingWindowLimiter{
		windowSize:  windowSize,
		maxRequests: maxRequests,
		requests:    make(map[string][]time.Time),
	}
}

// kw: allow sliding window limiter ctx context client rate abstraction
func (l *slidingWindowLimiter) Allow(ctx context.Context, client string) (bool, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now()
	requests := l.requests[client]

	// Remove requests outside the window
	for len(requests) > 0 && now.Sub(requests[0]) >= l.windowSize {
		requests = requests[1:]
	}

	// Update the requests slice
	l.requests[client] = requests

	// Check if we're under the limit
	if int64(len(requests)) < l.maxRequests {
		l.requests[client] = append(requests, now)
		return true, nil
	}

	return false, nil
}

// kw: count sliding window limiter client rate abstraction
func (l *slidingWindowLimiter) Count(client string) int64 {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	now := time.Now()
	requests := l.requests[client]

	// Remove requests outside the window
	for len(requests) > 0 && now.Sub(requests[0]) >= l.windowSize {
		requests = requests[1:]
	}

	return int64(len(requests))
}
