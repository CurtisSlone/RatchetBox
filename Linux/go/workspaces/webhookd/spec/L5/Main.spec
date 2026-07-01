name: Main
role: behavior
dependsOn: Server, Worker, Breaker, RetryPolicy, DeadLetter, Wal, Queue, Deliverer, Job
intent: The runnable entry point, now DURABLE, RESILIENT, and CORRECT-under-replay. It demonstrates two new guarantees: (1) a permanently-failing destination ends up in the dead-letter queue instead of being lost, and (2) a webhook replayed with the same Idempotency-Key is de-duplicated (only one delivery).
api:
  - func main()
behavior:
  - "Create a context with context.WithCancel(context.Background()); defer cancel()."
  - "Open the WAL: wal, err := OpenWal(filepath.Join(os.TempDir(), \"webhookd.wal\")); on error print and return. defer wal.Close()."
  - "RECOVERY: recovered, _ := wal.Replay(); q := NewQueue(); push every recovered job onto it. Print how many were recovered."
  - "DEAD destination (pinned): dead := httptest.NewServer(a handler that ALWAYS responds http.StatusInternalServerError). defer dead.Close(). Deliveries to it will exhaust retries and dead-letter."
  - "Wire the pipeline: deliverer := NewHTTPDeliverer(nil); breaker := NewBreaker(10, 100*time.Millisecond); policy := NewRetryPolicy(2*time.Millisecond, 2.0, 3); dlq := NewDeadLetter(); worker := NewWorker(q, deliverer, breaker, policy, dlq); go worker.Run(ctx); srv := NewServer(q, wal)."
  - "Stand up the ingest server: ingest := httptest.NewServer(srv.Handler()); defer ingest.Close()."
  - "DEAD-LETTER demo: POST one webhook to ingest for the dead destination ({\"url\": dead.URL, \"payload\": {\"x\":1}}). Expect 202."
  - "POLL with a DEADLINE LOOP (pinned - do NOT use a for/select with a break inside the select; break exits the select not the for, leaving the loop infinite and later code unreachable): `deadline := time.Now().Add(3 * time.Second); for dlq.Len() < 1 && time.Now().Before(deadline) { time.Sleep(10 * time.Millisecond) }`. Then print the dead-letter count (dlq.Len())."
  - "IDEMPOTENCY demo: send the SAME webhook TWICE with the same header \"Idempotency-Key\": \"dup-1\" (body {\"url\": dead.URL, \"payload\": {\"y\":2}}). CRITICAL (pinned): build a FRESH *http.Request with a NEW body reader for EACH send - do NOT reuse one *http.Request across two client.Do calls, because the first Do drains its Body and the second fails with \"ContentLength=N with Body length 0\". A helper `send := func() int { req, _ := http.NewRequest(\"POST\", ingest.URL+\"/webhook\", bytes.NewReader(payloadBytes)); req.Header.Set(\"Idempotency-Key\", \"dup-1\"); resp, err := http.DefaultClient.Do(req); if err != nil { panic(err) }; defer resp.Body.Close(); return resp.StatusCode }` called twice is the clean way. First status 202, second 200 (duplicate). Print both statuses."
  - "main must run to completion and exit normally (no panic). Keep total runtime under ~6 seconds."
constraints: package main; standard library only (context, net/http, net/http/httptest, encoding/json, sync/atomic, time, fmt, bytes, os, path/filepath); uses Server, Worker, Breaker, RetryPolicy, DeadLetter, Wal, Queue, Deliverer, Job
