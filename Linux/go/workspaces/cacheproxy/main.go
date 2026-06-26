package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	fetcher := NewMockFetcher(10 * time.Millisecond)
	cache := NewCache()
	proxy := NewProxy(cache, fetcher)
	httpServer := NewHTTPServer(proxy)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           httpServer.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return srv.ListenAndServe()
	})

	g.Go(func() error {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	})

	// Handle OS signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		ctx.Done()
	}()

	if err := g.Wait(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
