name: Main
role: behavior
dependsOn: Server, Worker, Breaker, RetryPolicy, Wal, Queue, Deliverer, Job
intent: The runnable entry point, now DURABLE and RESILIENT. It opens the write-ahead log, replays recovered jobs, wires the pipeline with a circuit breaker and retry policy, and demonstrates end to end that a FLAKY destination (failing the first couple of attempts) is still delivered to thanks to retries.
api:
  - func main()
behavior:
  - "Create a context.Context with context.WithCancel(context.Background()); defer cancel()."
  - "Open the write-ahead log: wal, err := OpenWal(filepath.Join(os.TempDir(), \"webhookd.wal\")); on error, print and return. defer wal.Close()."
  - "RECOVERY: recovered, _ := wal.Replay(); q := NewQueue(); push every recovered job onto it. Print how many jobs were recovered."
  - "FLAKY destination (pinned, to exercise retries): stand up dest := httptest.NewServer(handler) where the handler fails the FIRST 2 calls (respond http.StatusInternalServerError) and succeeds afterward (200), counting successful deliveries with an atomic int64. Guard the call counter with sync/atomic. defer dest.Close()."
  - "Wire the resilient pipeline: deliverer := NewHTTPDeliverer(nil); breaker := NewBreaker(5, 100*time.Millisecond); policy := NewRetryPolicy(5*time.Millisecond, 2.0, 6); worker := NewWorker(q, deliverer, breaker, policy); go worker.Run(ctx); srv := NewServer(q, wal)."
  - "Stand up the ingest server with httptest.NewServer(srv.Handler()); defer ingest.Close()."
  - "POST a webhook to the ingest server: JSON body {\"url\": dest.URL, \"payload\": {\"hello\":\"world\"}} via http.Post(ingest.URL+\"/webhook\", \"application/json\", body). Expect 202."
  - "POLL with a DEADLINE LOOP (pinned - do NOT use a for/select with a break inside the select, because break exits the select not the for, leaving the loop infinite and the code after it unreachable): `deadline := time.Now().Add(3 * time.Second); for atomic.LoadInt64(&successes) < 1 && time.Now().Before(deadline) { time.Sleep(10 * time.Millisecond) }`. After the loop, print the number of successful deliveries and the breaker state (breaker.State()) with a clear line like \"delivered N webhook(s) after retries, breaker=closed\". main must run to completion and exit normally (no panic)."
constraints: package main; standard library only (context, net/http, net/http/httptest, encoding/json, sync/atomic, time, fmt, bytes, os, path/filepath); uses Server, Worker, Breaker, RetryPolicy, Wal, Queue, Deliverer, Job
