package main

import (
	"sync"
	"testing"
	"time"
)

func TestProxy(t *testing.T) {
	// Herd test: a couple hundred goroutines request the SAME key at once
	t.Run("herd", func(t *testing.T) {
		cache := NewCache()
		fetcher := NewMockFetcher(10 * time.Millisecond)
		proxy := NewProxy(cache, fetcher)

		const numGoroutines = 200
		var wg sync.WaitGroup
		results := make([]string, numGoroutines)
		errors := make([]error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				results[i], errors[i] = proxy.Get("key")
			}(i)
		}

		wg.Wait()

		// Check all results are the same value
		if len(results) == 0 {
			t.Fatal("no results")
		}
		firstValue := results[0]
		for i, err := range errors {
			if err != nil {
				t.Errorf("goroutine %d got error: %v", i, err)
			}
			if results[i] != firstValue {
				t.Errorf("goroutine %d got value %q, want %q", i, results[i], firstValue)
			}
		}

		// Check the origin was called exactly once
		if fetcher.Count() != 1 {
			t.Errorf("origin called %d times, want 1", fetcher.Count())
		}
	})

	// Cache-hit test: requesting a warm key a second time does not call the origin again
	t.Run("cache-hit", func(t *testing.T) {
		cache := NewCache()
		fetcher := NewMockFetcher(0)
		proxy := NewProxy(cache, fetcher)

		// First request
		value1, err1 := proxy.Get("key")
		if err1 != nil {
			t.Fatalf("first request failed: %v", err1)
		}

		// Second request
		value2, err2 := proxy.Get("key")
		if err2 != nil {
			t.Fatalf("second request failed: %v", err2)
		}

		// Values should be equal
		if value1 != value2 {
			t.Errorf("first and second requests returned different values: %q vs %q", value1, value2)
		}

		// Origin should only be called once
		if fetcher.Count() != 1 {
			t.Errorf("origin called %d times, want 1", fetcher.Count())
		}
	})

	// Distinct-keys test: concurrent requests across a handful of different keys
	t.Run("distinct-keys", func(t *testing.T) {
		cache := NewCache()
		fetcher := NewMockFetcher(5 * time.Millisecond)
		proxy := NewProxy(cache, fetcher)

		const numKeys = 5
		const numGoroutines = 100
		var wg sync.WaitGroup
		results := make([]string, numGoroutines)
		errors := make([]error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				key := "key" + string(rune('0'+i%numKeys))
				results[i], errors[i] = proxy.Get(key)
			}(i)
		}

		wg.Wait()

		// Check all results are the expected value
		for i, err := range errors {
			if err != nil {
				t.Errorf("goroutine %d got error: %v", i, err)
			}
		}

		// Check the origin was called exactly once per unique key
		if fetcher.Count() != numKeys {
			t.Errorf("origin called %d times, want %d", fetcher.Count(), numKeys)
		}
	})
}
