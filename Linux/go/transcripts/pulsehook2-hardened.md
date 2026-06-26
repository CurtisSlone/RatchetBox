# pulsehook2 - hardened rebuild (the review -> KB -> better-code loop)

The first `pulsehook` (see `pulsehook-build.md`) was clean but demo-grade; a code review found four real
production gaps. Each became a knowledge-base entry, and `pulsehook2` was composed from hardened specs
with those entries grounding generation - so the gaps are gone WITHOUT hand-writing the fixes. This is
the evidence loop: an observed weakness -> a KB entry -> the next build is born with the fix.

- Generated: 2026-06-26
- Command: `ratchet flow . compose --ws pulsehook2 "" ` (console: `/flow compose --ws pulsehook2`)
- KB added first: `kb/patterns/algo_production_http_server.md`, `kb/pitfalls/atomic-int64-alignment.md`;
  `add_unit` now also searches `pitfalls` + `guidelines`.
- Result: composed first try except `dispatcher.go` (one repair); `go build ./...` + `go test ./...` pass.

## The four findings, before -> after (model-generated both times)

### 1. Atomic counter alignment (pitfall: atomic-int64-alignment)
```go
// pulsehook (before): raw int64, misaligned-atomic risk on 386/arm
processed int64
atomic.AddInt64(&d.processed, 1)
return atomic.LoadInt64(&d.processed)

// pulsehook2 (after): typed atomic, always alignment-safe
processed atomic.Int64
d.processed.Add(1)
return d.processed.Load()
```

### 2. Request body cap (pattern: production-http-server)
```go
// before: unbounded read - a huge POST can OOM the process
body, err := io.ReadAll(r.Body)

// after: capped at 1 MiB, 413 on overflow
r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
body, err := io.ReadAll(r.Body)
if err != nil {
	http.Error(w, "Request Entity Too Large", http.StatusRequestEntityTooLarge)
	return
}
```

### 3 & 4. Server timeouts + graceful shutdown (pattern: production-http-server)
```go
// before (main.go): bare server, no timeouts, no shutdown - queued events lost on exit
log.Println("Listening on :8080")
if err := http.ListenAndServe(":8080", server.Routes()); err != nil {
	log.Fatal(err)
}

// after (main.go): timeouts + signal-driven graceful shutdown that drains the worker pool
srv := &http.Server{
	Addr:              ":8080",
	Handler:           server.Routes(),
	ReadHeaderTimeout: 5 * time.Second,
	ReadTimeout:       15 * time.Second,
	WriteTimeout:      15 * time.Second,
	IdleTimeout:       60 * time.Second,
}
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
defer stop()
go func() {
	log.Println("listening on", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}()
<-ctx.Done()
stop()
shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
if err := srv.Shutdown(shutCtx); err != nil {
	log.Println("shutdown:", err)
}
dispatcher.Stop() // drain queued events
```

## Takeaway

The specs named the production concerns; the KB supplied the exact idioms (the `production-http-server`
pattern and the `atomic-int64-alignment` pitfall, retrieved by `add_unit`'s per-unit search). Same
ratchet, same local model - the difference is the knowledge it was grounded on. That is how this ratchet
gets to production-grade by default: feed review findings back as KB entries.
