name: Main
role: behavior
intent: The runnable entry point. It reads configuration from flags, connects to Redis, builds a circuit-breaker-guarded Redis sliding-window limiter, and serves an HTTP endpoint protected by the rate-limit middleware.
api:
  - func main()
behavior:
  - "Define flags: -addr (HTTP listen address, default \":8080\"), -redis (Redis address, default \"localhost:6379\"), -limit (int requests per window, default 100), -window (time.Duration window, default time.Minute)."
  - "Connect: client := redis.NewClient(&redis.Options{Addr: *redisAddr}). Imports github.com/redis/go-redis/v9 as redis."
  - "Build: breaker := NewBreaker(5, 10*time.Second); limiter := NewRedisLimiter(client, breaker, *limit, *window); srv := NewServer(limiter)."
  - "Serve: build an http.Server with the srv.Handler() and sensible ReadTimeout/WriteTimeout, log the listen address, and ListenAndServe; log.Fatal on a non-ErrServerClosed error."
  - "main must compile and start without panicking; it is fine for it to block in ListenAndServe."
constraints: package main; uses Breaker, RedisLimiter, Server, github.com/redis/go-redis/v9; standard library (flag, log, net/http, time)
