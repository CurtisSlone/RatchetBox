package main

// file-kw: retrypolicy deterministic exponential backoff retry policy attempt number how long wait before next try

import (
	"time"
)

// RetryPolicy is a deterministic exponential-backoff retry policy.
// kw: retry policy deterministic exponential backoff
type RetryPolicy struct {
	Base        time.Duration
	Factor      float64
	MaxAttempts int
}

// NewRetryPolicy creates a new retry policy with the given base duration,
// factor, and maximum number of attempts.
// kw: retry policy base time duration factor max attempts deterministic
func NewRetryPolicy(base time.Duration, factor float64, maxAttempts int) RetryPolicy {
	return RetryPolicy{
		Base:        base,
		Factor:      factor,
		MaxAttempts: maxAttempts,
	}
}

// Backoff returns the backoff duration for the given attempt number.
// kw: backoff retry policy attempt time duration deterministic exponential
func (p RetryPolicy) Backoff(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	factor := 1.0
	for i := 1; i < attempt; i++ {
		factor *= p.Factor
	}
	return time.Duration(factor * float64(p.Base))
}

// ShouldRetry reports whether another attempt is allowed.
// kw: should retry policy attempt deterministic exponential backoff
func (p RetryPolicy) ShouldRetry(attempt int) bool {
	if p.MaxAttempts <= 0 {
		return false
	}
	return attempt < p.MaxAttempts
}
