package main

// file-kw: redislimiter redis backed sliding window rate limiter real service integration keeps sorted set client

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

// kw: redis limiter backed sliding
type RedisLimiter struct {
	client  *redis.Client
	breaker *Breaker
	limit   int64
	window  time.Duration
	seq     *atomic.Int64
	mutex   sync.Mutex
}

// kw: redis limiter client breaker limit window time duration backed
func NewRedisLimiter(client *redis.Client, breaker *Breaker, limit int, window time.Duration) *RedisLimiter {
	return &RedisLimiter{
		client:  client,
		breaker: breaker,
		limit:   int64(limit),
		window:  window,
		seq:     &atomic.Int64{},
	}
}

// kw: allow redis limiter ctx context client backed sliding
func (l *RedisLimiter) Allow(ctx context.Context, clientID string) (bool, error) {
	now := time.Now().UnixNano()
	key := "ratelimit:" + clientID

	if !l.breaker.Allow() {
		return true, nil // FAIL OPEN
	}

	seq := l.seq.Add(1)
	pipe := l.client.TxPipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(now-int64(l.window), 10))
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d-%d", now, seq)})
	count := pipe.ZCard(ctx, key)
	pipe.Expire(ctx, key, l.window)
	_, err := pipe.Exec(ctx)

	if err != nil {
		l.breaker.Failure()
		return true, err // FAIL OPEN
	}

	l.breaker.Success()
	allowed := count.Val() <= l.limit
	return allowed, nil
}
