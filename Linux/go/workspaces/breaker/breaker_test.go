package main

import (
	"sync"
	"testing"
	"time"
)

func TestBreaker(t *testing.T) {
	t.Run("Closed start", func(t *testing.T) {
		b := NewBreaker(3, time.Minute)
		if got := b.State(); got != "closed" {
			t.Errorf("NewBreaker(3, time.Minute).State() = %q; want \"closed\"", got)
		}
	})

	t.Run("Trips open", func(t *testing.T) {
		b := NewBreaker(3, time.Minute)
		b.Failure()
		b.Failure()
		b.Failure()
		if got := b.State(); got != "open" {
			t.Errorf("after 3 failures, State() = %q; want \"open\"", got)
		}
		if got := b.Allow(); got {
			t.Errorf("after 3 failures, Allow() = %v; want false", got)
		}
	})

	t.Run("Cooldown -> half-open", func(t *testing.T) {
		b := NewBreaker(1, time.Millisecond)
		b.Failure()
		time.Sleep(5 * time.Millisecond)
		if got := b.Allow(); !got {
			t.Errorf("after cooldown, Allow() = %v; want true", got)
		}
		if got := b.State(); got != "half-open" {
			t.Errorf("after cooldown, State() = %q; want \"half-open\"", got)
		}
	})

	t.Run("Half-open success closes", func(t *testing.T) {
		b := NewBreaker(1, time.Millisecond)
		b.Failure()
		time.Sleep(5 * time.Millisecond)
		b.Success()
		if got := b.State(); got != "closed" {
			t.Errorf("after half-open success, State() = %q; want \"closed\"", got)
		}
	})

	t.Run("Half-open failure reopens", func(t *testing.T) {
		b := NewBreaker(1, time.Millisecond)
		b.Failure()
		time.Sleep(5 * time.Millisecond)
		b.Failure() // This should reopen it
		if got := b.State(); got != "open" {
			t.Errorf("after half-open failure, State() = %q; want \"open\"", got)
		}
	})

	t.Run("Concurrency", func(t *testing.T) {
		b := NewBreaker(3, time.Minute)
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				b.Allow()
				b.Success()
				b.Failure()
			}()
		}
		wg.Wait()
		// Just ensure no panic and valid state
		validStates := map[string]bool{"closed": true, "open": true, "half-open": true}
		if !validStates[b.State()] {
			t.Errorf("concurrent operations resulted in invalid state %q", b.State())
		}
	})
}

func FuzzBreaker(f *testing.F) {
	f.Add([]byte{2, 2, 2, 0, 1}) // seed
	f.Fuzz(func(t *testing.T, data []byte) {
		b := NewBreaker(3, time.Millisecond)
		for _, op := range data {
			switch op % 3 {
			case 0:
				b.Allow()
			case 1:
				b.Success()
			case 2:
				b.Failure()
			}
		}
		validStates := map[string]bool{"closed": true, "open": true, "half-open": true}
		if !validStates[b.State()] {
			t.Errorf("final state %q is not a valid breaker state", b.State())
		}
	})
}
