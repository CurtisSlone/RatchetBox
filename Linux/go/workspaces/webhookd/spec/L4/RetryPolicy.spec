name: RetryPolicy
role: component
intent: A deterministic exponential-backoff retry policy. Given an attempt number it returns how long to wait before the next try, and whether another try is allowed. Deterministic (no jitter) so it is exactly testable.
api:
  - type RetryPolicy struct { Base time.Duration; Factor float64; MaxAttempts int }
  - func NewRetryPolicy(base time.Duration, factor float64, maxAttempts int) RetryPolicy
  - func (p RetryPolicy) Backoff(attempt int) time.Duration
  - func (p RetryPolicy) ShouldRetry(attempt int) bool
behavior:
  - "Attempt numbers are 1-based: attempt 1 is the first try."
  - "Backoff(attempt) = Base * Factor^(attempt-1), as a time.Duration. So Backoff(1) == Base, Backoff(2) == Base*Factor, Backoff(3) == Base*Factor^2. For attempt < 1, treat it as 1 (return Base). No jitter - the result is exact and deterministic."
  - "ShouldRetry(attempt) reports whether another attempt is allowed: it returns true iff attempt < MaxAttempts. So with MaxAttempts == 3, ShouldRetry(1) and ShouldRetry(2) are true and ShouldRetry(3) is false (3 attempts have been made)."
  - "Edge: MaxAttempts <= 0 means never retry - ShouldRetry returns false for every attempt. Never panic."
constraints: package main; standard library only
