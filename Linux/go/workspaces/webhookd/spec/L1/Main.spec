name: Main
role: behavior
dependsOn: Server, Worker, Wal, Queue, Deliverer, Job
intent: The runnable entry point, now DURABLE. It opens the write-ahead log, replays any jobs recovered from a previous run back onto the queue, then wires and demonstrates the pipeline end to end - an in-process destination receives a webhook that was ingested over HTTP, persisted to the WAL, and delivered by the worker.
api:
  - func main()
behavior:
  - "Create a context.Context with context.WithCancel(context.Background()); defer cancel()."
  - "Open the write-ahead log: wal, err := OpenWal(a temp path, e.g. filepath.Join(os.TempDir(), \"webhookd.wal\")); on error, print and return. defer wal.Close()."
  - "RECOVERY (pinned): recovered, _ := wal.Replay(); create q := NewQueue() and push every recovered job onto it (for _, j := range recovered { q.Push(j) }). Print how many jobs were recovered from the WAL."
  - "Stand up an in-process destination with net/http/httptest: dest := httptest.NewServer(an http.HandlerFunc that reads the body and records that it was hit - increment an int64 via sync/atomic). defer dest.Close()."
  - "Wire the pipeline: deliverer := NewHTTPDeliverer(nil); worker := NewWorker(q, deliverer); go worker.Run(ctx); srv := NewServer(q, wal)."
  - "Stand up the ingest server with httptest.NewServer(srv.Handler()); defer ingest.Close()."
  - "POST a webhook to the ingest server: build a JSON body {\"url\": dest.URL, \"payload\": {\"hello\":\"world\"}} and http.Post(ingest.URL+\"/webhook\", \"application/json\", the body). Expect 202."
  - "Poll for up to ~2 seconds (loop with short time.Sleep) until the destination's hit counter reaches 1, then print how many deliveries the destination received and the queue length. main must run to completion and exit normally (no panic); it should print a clear line like \"delivered N webhook(s)\"."
constraints: package main; standard library only (context, net/http, net/http/httptest, encoding/json, sync/atomic, time, fmt, bytes, os, path/filepath); uses Server, Worker, Wal, Queue, Deliverer, Job
