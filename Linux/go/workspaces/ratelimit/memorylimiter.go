package main

// file-kw: memory rate limiter in-process
import (
	"context"
	"sync"
	"time"
)

// kw: check client request allowance
type MemoryLimiter struct {
	limit    int64
	window   time.Duration
	requests map[string][]int64
	mutex    sync.Mutex
}

// kw: create new memory limiter
func NewMemoryLimiter(limit int, window time.Duration) *MemoryLimiter {
	return &MemoryLimiter{
		limit:    int64(limit),
		window:   window,
		requests: make(map[string][]int64),
	}
}

// kw: allow client request or deny
func (l *MemoryLimiter) Allow(ctx context.Context, clientID string) (bool, error) {
	now := time.Now().UnixNano()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	clientRequests := l.requests[clientID]
	// Drop old requests outside the window
	cutoff := now - l.window.Nanoseconds()
	var validRequests []int64
	for _, reqTime := range clientRequests {
		if reqTime > cutoff {
			validRequests = append(validRequests, reqTime)
		}
	}
	// Append current request
	validRequests = append(validRequests, now)
	allowed := int64(len(validRequests)) <= l.limit
	l.requests[clientID] = validRequests
	return allowed, nil
}
