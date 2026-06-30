name: RedisLimiter
role: component
intent: A Redis-backed sliding-window rate limiter (the real service integration). It keeps one Redis sorted set per client whose members are request timestamps; a request is allowed iff the number of timestamps still inside the window is at most the limit. A circuit breaker guards the Redis calls so a Redis outage fails OPEN (allows traffic) instead of cascading.
api:
  - type RedisLimiter struct { ... }
  - func NewRedisLimiter(client *redis.Client, breaker *Breaker, limit int, window time.Duration) *RedisLimiter
  - func (l *RedisLimiter) Allow(ctx context.Context, clientID string) (bool, error)
behavior:
  - "Imports github.com/redis/go-redis/v9 as `redis`. The constructor stores the *redis.Client, the *Breaker, the integer limit (max requests per window), and the window duration."
  - "SLIDING WINDOW (pinned, using a sorted set per client): the key is \"ratelimit:\" + clientID. On each Allow, with now := time.Now().UnixNano() and a per-instance atomic sequence counter for uniqueness:
      1. if !breaker.Allow() -> FAIL OPEN: return (true, nil) without touching Redis (Redis is considered unhealthy).
      2. run a redis TxPipeline (l.client.TxPipeline()):
           pipe.ZRemRangeByScore(ctx, key, \"0\", strconv.FormatInt(now-int64(l.window), 10))   // drop entries older than the window
           pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: fmt.Sprintf(\"%d-%d\", now, seq)})  // add this request (unique member)
           count := pipe.ZCard(ctx, key)
           pipe.Expire(ctx, key, l.window)
           _, err := pipe.Exec(ctx)
      3. on err != nil -> breaker.Failure(); FAIL OPEN: return (true, err).
      4. on success -> breaker.Success(); allowed := count.Val() <= int64(l.limit); return (allowed, nil)."
  - "The window-start score is now - window-in-nanoseconds (window is a time.Duration, which is int64 nanoseconds, so int64(l.window) is the nanosecond span)."
  - "Members must be UNIQUE even within the same nanosecond - use an atomic counter (sync/atomic) appended to the timestamp - or a request in the same nanosecond would overwrite a prior one and undercount."
  - "Allow must be safe for concurrent use (the *redis.Client is concurrency-safe; guard the sequence counter atomically)."
constraints: package main; uses github.com/redis/go-redis/v9, Breaker; standard library plus that one dependency
