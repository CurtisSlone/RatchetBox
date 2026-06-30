package main

// file-kw: main entry point http server redis limiter

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

// kw: parse flags configure
var (
	addr      = flag.String("addr", ":8080", "HTTP listen address")
	redisAddr = flag.String("redis", "localhost:6379", "Redis address")
	limit     = flag.Int("limit", 100, "requests per window")
	window    = flag.Duration("window", time.Minute, "window duration")
)

// kw: start http server with middleware
func main() {
	flag.Parse()

	client := redis.NewClient(&redis.Options{
		Addr: *redisAddr,
	})

	breaker := NewBreaker(5, 10*time.Second)
	limiter := NewRedisLimiter(client, breaker, *limit, *window)
	srv := NewServer(limiter)

	mux := srv.Handler()
	server := &http.Server{
		Addr:              *addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("listening on", *addr)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	stop()
	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutCtx); err != nil {
		log.Println("shutdown:", err)
	}
}
