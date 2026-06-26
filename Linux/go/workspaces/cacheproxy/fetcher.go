package main

import (
	"context"
	"sync"
	"time"
)

// Fetcher interface defines the contract for fetching values by key
type Fetcher interface {
	Fetch(ctx context.Context, key string) ([]byte, error)
}

// MockFetcher implements the Fetcher interface for testing and demo purposes
type MockFetcher struct {
	mu      sync.Mutex
	count   int
	latency time.Duration
}

// NewMockFetcher creates a new MockFetcher with optional artificial latency
func NewMockFetcher(latency time.Duration) *MockFetcher {
	return &MockFetcher{
		latency: latency,
	}
}

// Fetch returns a deterministic value derived from the key and increments the call counter
func (mf *MockFetcher) Fetch(ctx context.Context, key string) ([]byte, error) {
	// Simulate artificial latency
	if mf.latency > 0 {
		time.Sleep(mf.latency)
	}

	mf.mu.Lock()
	mf.count++
	count := mf.count
	mf.mu.Unlock()

	// Return deterministic value based on key and call count
	return []byte("value-for-" + key + "-call-" + string(rune(count+'0'))), nil
}

// Count returns the number of times Fetch has been called
func (mf *MockFetcher) Count() int {
	mf.mu.Lock()
	defer mf.mu.Unlock()
	return mf.count
}
