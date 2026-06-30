name: MemoryLimiter
role: component
intent: An in-process sliding-window rate limiter with the same policy as RedisLimiter, backed by a map of per-client timestamp slices. It needs no Redis, so it is the limiter the tests use (and a usable fallback when Redis is absent).
api:
  - type MemoryLimiter struct { ... }
  - func NewMemoryLimiter(limit int, window time.Duration) *MemoryLimiter
  - func (l *MemoryLimiter) Allow(ctx context.Context, clientID string) (bool, error)
behavior:
  - "Stores, per clientID, the timestamps (time.Time or int64 nanos) of recent requests, guarded by a sync.Mutex."
  - "On Allow with now := time.Now(): drop this client's timestamps older than now.Add(-window); append now; allowed := len(remaining) <= limit; return (allowed, nil). It never returns an error (in-memory, nothing to fail)."
  - "SAME POLICY as RedisLimiter: a request is allowed iff the count of timestamps within the last `window` (inclusive of the new one) is <= limit. So with limit=2, the 1st and 2nd requests in a window are allowed and the 3rd is denied, until the window slides past the older ones."
  - "ctx is accepted to satisfy the Limiter interface but is not otherwise needed. Allow must be safe for concurrent use."
constraints: package main; standard library only
