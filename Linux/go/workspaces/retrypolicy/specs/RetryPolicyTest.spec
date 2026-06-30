name: RetryPolicyTest
role: test
intent: Prove the exponential-backoff policy: exact backoff values, the retry boundary, the edge cases, and the monotonic-growth property.
api:
  - func TestRetryPolicy(t *testing.T)
behavior:
  - "Exact backoff: p := NewRetryPolicy(time.Millisecond, 2.0, 5); p.Backoff(1) == time.Millisecond; p.Backoff(2) == 2*time.Millisecond; p.Backoff(3) == 4*time.Millisecond."
  - "Attempt < 1 is treated as 1: p.Backoff(0) == time.Millisecond and p.Backoff(-3) == time.Millisecond."
  - "Retry boundary: with MaxAttempts == 3, p.ShouldRetry(1) and p.ShouldRetry(2) are true, p.ShouldRetry(3) is false."
  - "Edge: NewRetryPolicy(time.Millisecond, 2.0, 0).ShouldRetry(1) is false (MaxAttempts <= 0 never retries); never panics."
  - "PROPERTY (fuzz): FuzzBackoff(f) fuzzes an int attempt; with Factor >= 1 the backoff is monotonic non-decreasing - assert Backoff(attempt+1) >= Backoff(attempt) for attempt clamped to a sane range (e.g. 1..30 to avoid overflow), and Backoff(attempt) >= Base for all attempt. Seed with f.Add(1) and f.Add(10)."
constraints: package main; standard library only
