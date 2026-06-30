package main

// file-kw: retrypolicy deterministic exponential backoff retry policy attempt number how long wait before next try

import (
	"math"
	"time"
)

// kw: retry policy deterministic exponential backoff
type RetryPolicy struct {
	Base        time.Duration
	Factor      float64
	MaxAttempts int
}

// kw: retry policy base time duration factor max attempts deterministic
func NewRetryPolicy(base time.Duration, factor float64, maxAttempts int) RetryPolicy {
	return RetryPolicy{
		Base:        base,
		Factor:      factor,
		MaxAttempts: maxAttempts,
	}
}

// kw: backoff retry policy attempt time duration deterministic exponential
func (p RetryPolicy) Backoff(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}

	if p.MaxAttempts <= 0 {
		return 0
	}

	exp := math.Pow(p.Factor, float64(attempt-1))
	duration := float64(p.Base) * exp
	return time.Duration(duration)
}

// kw: should retry policy attempt deterministic exponential backoff
func (p RetryPolicy) ShouldRetry(attempt int) bool {
	if p.MaxAttempts <= 0 {
		return false
	}
	return attempt < p.MaxAttempts
}
