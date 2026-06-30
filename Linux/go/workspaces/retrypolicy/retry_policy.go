package main

import (
	"math"
	"time"
)

type RetryPolicy struct {
	Base        time.Duration
	Factor      float64
	MaxAttempts int
}

func NewRetryPolicy(base time.Duration, factor float64, maxAttempts int) RetryPolicy {
	return RetryPolicy{
		Base:        base,
		Factor:      factor,
		MaxAttempts: maxAttempts,
	}
}

func (p RetryPolicy) Backoff(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}

	// Calculate factor^attempt-1
	power := math.Pow(p.Factor, float64(attempt-1))

	// Convert to time.Duration
	backoff := float64(p.Base) * power
	return time.Duration(backoff)
}

func (p RetryPolicy) ShouldRetry(attempt int) bool {
	return attempt < p.MaxAttempts && p.MaxAttempts > 0
}
