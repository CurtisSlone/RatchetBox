# Production HTTP server (Go pattern)

A network-facing `net/http` server should not use bare `http.ListenAndServe`. Use an explicit
`*http.Server` with timeouts, cap request bodies, and shut down gracefully so in-flight work drains.
Authored pattern; targets the gaps a plain handler leaves open (slowloris, unbounded bodies, lost work
on exit).

Three things every production server needs:
- TIMEOUTS on the `http.Server` (`ReadHeaderTimeout` at minimum; plus `ReadTimeout`/`WriteTimeout`/
  `IdleTimeout`) - a bare server has none and is vulnerable to slow-client (slowloris) attacks.
- BODY LIMITS on untrusted input: `http.MaxBytesReader(w, r.Body, max)` before reading, so a huge POST
  cannot exhaust memory.
- GRACEFUL SHUTDOWN: catch SIGINT/SIGTERM with `signal.NotifyContext`, call `srv.Shutdown(ctx)` to stop
  accepting and drain active requests, then stop any worker pools.

```go
package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func handle(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // cap body at 1 MiB
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "request too large", http.StatusRequestEntityTooLarge)
		return
	}
	_ = body
	w.WriteHeader(http.StatusAccepted)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", handle)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Listen for SIGINT/SIGTERM; ctx is cancelled on signal.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("listening on", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done() // wait for the signal
	stop()
	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutCtx); err != nil { // drain in-flight requests
		log.Println("shutdown:", err)
	}
	// stop worker pools here (e.g. dispatcher.Stop()) so queued work finishes.
}
```
