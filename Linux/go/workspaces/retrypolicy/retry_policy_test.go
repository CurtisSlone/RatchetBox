package main

import (
	"testing"
	"time"
)

func TestRetryPolicy(t *testing.T) {
	t.Run("ExactBackoff", func(t *testing.T) {
		p := NewRetryPolicy(time.Millisecond, 2.0, 5)
		if got := p.Backoff(1); got != time.Millisecond {
			t.Errorf("Backoff(1) = %v, want %v", got, time.Millisecond)
		}
		if got := p.Backoff(2); got != 2*time.Millisecond {
			t.Errorf("Backoff(2) = %v, want %v", got, 2*time.Millisecond)
		}
		if got := p.Backoff(3); got != 4*time.Millisecond {
			t.Errorf("Backoff(3) = %v, want %v", got, 4*time.Millisecond)
		}
	})

	t.Run("AttemptLessThanOne", func(t *testing.T) {
		p := NewRetryPolicy(time.Millisecond, 2.0, 5)
		if got := p.Backoff(0); got != time.Millisecond {
			t.Errorf("Backoff(0) = %v, want %v", got, time.Millisecond)
		}
		if got := p.Backoff(-3); got != time.Millisecond {
			t.Errorf("Backoff(-3) = %v, want %v", got, time.Millisecond)
		}
	})

	t.Run("RetryBoundary", func(t *testing.T) {
		p := NewRetryPolicy(time.Millisecond, 2.0, 3)
		if !p.ShouldRetry(1) {
			t.Error("ShouldRetry(1) should be true")
		}
		if !p.ShouldRetry(2) {
			t.Error("ShouldRetry(2) should be true")
		}
		if p.ShouldRetry(3) {
			t.Error("ShouldRetry(3) should be false")
		}
	})

	t.Run("EdgeMaxAttemptsZero", func(t *testing.T) {
		p := NewRetryPolicy(time.Millisecond, 2.0, 0)
		if p.ShouldRetry(1) {
			t.Error("ShouldRetry(1) should be false when MaxAttempts <= 0")
		}
	})

	t.Run("EdgeMaxAttemptsNegative", func(t *testing.T) {
		p := NewRetryPolicy(time.Millisecond, 2.0, -1)
		if p.ShouldRetry(1) {
			t.Error("ShouldRetry(1) should be false when MaxAttempts <= 0")
		}
	})
}

func FuzzBackoff(f *testing.F) {
	f.Add(1)
	f.Add(10)

	f.Fuzz(func(t *testing.T, attempt int) {
		// Clamp attempt to avoid overflow
		if attempt < 1 {
			attempt = 1
		}
		if attempt > 30 {
			attempt = 30
		}

		p := NewRetryPolicy(time.Millisecond, 2.0, 10)
		backoff := p.Backoff(attempt)
		base := p.Base

		// Backoff should be at least base
		if backoff < base {
			t.Errorf("Backoff(%d) = %v, want >= %v", attempt, backoff, base)
		}

		// For factor >= 1, backoff should be monotonic non-decreasing
		if attempt > 1 {
			prevBackoff := p.Backoff(attempt - 1)
			if backoff < prevBackoff {
				t.Errorf("Backoff(%d) = %v, want >= %v (previous)", attempt, backoff, prevBackoff)
			}
		}
	})
}
