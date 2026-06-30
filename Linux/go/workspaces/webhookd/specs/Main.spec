name: Main
role: behavior
dependsOn: Server, Worker, Breaker, RetryPolicy, DeadLetter, Metrics, Wal, Queue, Deliverer, Job
intent: The runnable entry point, now the full production shape - DURABLE, RESILIENT, CORRECT, and OBSERVABLE with GRACEFUL DRAIN. It runs the dead-letter and idempotency demos, then drains the queue cleanly before shutting down and prints the final metrics (also scrapeable at /metrics).
api:
  - func main()
behavior:
  - "Create a context with context.WithCancel(context.Background()); defer cancel()."
  - "Open the WAL: wal, err := OpenWal(filepath.Join(os.TempDir(), \"webhookd.wal\")); on error print and return. defer wal.Close()."
  - "RECOVERY: recovered, _ := wal.Replay(); q := NewQueue(); push every recovered job onto it. Print how many were recovered."
  - "DEAD destination (pinned): dead := httptest.NewServer(a handler that ALWAYS responds http.StatusInternalServerError). defer dead.Close()."
  - "Wire the pipeline: metrics := NewMetrics(); deliverer := NewHTTPDeliverer(nil); breaker := NewBreaker(10, 100*time.Millisecond); policy := NewRetryPolicy(2*time.Millisecond, 2.0, 3); dlq := NewDeadLetter(); worker := NewWorker(q, deliverer, breaker, policy, dlq, metrics); go worker.Run(ctx); srv := NewServer(q, wal, metrics)."
  - "Stand up the ingest server: ingest := httptest.NewServer(srv.Handler()); defer ingest.Close()."
  - "DEAD-LETTER demo: POST one webhook to ingest for the dead destination ({\"url\": dead.URL, \"payload\": {\"x\":1}}). Expect 202."
  - "POLL with a DEADLINE LOOP (pinned - never use a for/select with a break inside the select; break exits the select not the for, leaving the loop infinite and later code unreachable): `deadline := time.Now().Add(3 * time.Second); for dlq.Len() < 1 && time.Now().Before(deadline) { time.Sleep(10 * time.Millisecond) }`. Then print the dead-letter count (dlq.Len())."
  - "IDEMPOTENCY demo: send the SAME webhook TWICE with the same header \"Idempotency-Key\": \"dup-1\" (body {\"url\": dead.URL, \"payload\": {\"y\":2}}). CRITICAL (pinned): build a FRESH *http.Request with a NEW body reader for EACH of the two sends - do NOT reuse one *http.Request across two client.Do calls, because the first Do drains its Body and the second fails with \"ContentLength=N with Body length 0\". A helper like `send := func() int { req, _ := http.NewRequest(\"POST\", ingest.URL+\"/webhook\", bytes.NewReader(payloadBytes)); req.Header.Set(\"Idempotency-Key\", \"dup-1\"); resp, err := http.DefaultClient.Do(req); if err != nil { panic(err) }; defer resp.Body.Close(); return resp.StatusCode }` called twice is the clean way. First status 202, second 200 (duplicate). Print both statuses."
  - "GRACEFUL DRAIN (pinned): after the demos, wait for the queue to drain with a DEADLINE LOOP (same idiom, no for/select-break): `dd := time.Now().Add(3 * time.Second); for q.Len() > 0 && time.Now().Before(dd) { time.Sleep(10 * time.Millisecond) }`. THEN cancel the context to stop the worker. This ensures no in-flight job is abandoned at shutdown."
  - "OBSERVABILITY: GET ingest.URL+\"/metrics\" and print the body (the Prometheus-format counters). Also print metrics.Snapshot(). main must run to completion and exit normally (no panic); keep total runtime under ~8 seconds."
constraints: package main; standard library only (context, net/http, net/http/httptest, encoding/json, io, sync/atomic, time, fmt, bytes, os, path/filepath); uses Server, Worker, Breaker, RetryPolicy, DeadLetter, Metrics, Wal, Queue, Deliverer, Job
