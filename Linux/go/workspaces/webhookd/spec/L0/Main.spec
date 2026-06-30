name: Main
role: behavior
dependsOn: Server, Worker, Queue, Deliverer, Job
intent: The runnable entry point that wires the walking skeleton together and demonstrates it end to end - an in-process destination receives a webhook that was ingested over HTTP and delivered by the worker. Proves the ingest -> queue -> worker -> deliver path works.
api:
  - func main()
behavior:
  - "Create a context.Context with context.WithCancel(context.Background()); defer cancel()."
  - "Stand up an in-process destination with net/http/httptest: dest := httptest.NewServer(an http.HandlerFunc that reads the body and records that it was hit - increment an int64 via sync/atomic). defer dest.Close(). This is where the webhook will be delivered."
  - "Wire the pipeline: q := NewQueue(); deliverer := NewHTTPDeliverer(nil); worker := NewWorker(q, deliverer); go worker.Run(ctx); srv := NewServer(q)."
  - "Stand up the ingest server with httptest.NewServer(srv.Handler()); defer ingest.Close()."
  - "POST a webhook to the ingest server: build a JSON body {\"url\": dest.URL, \"payload\": {\"hello\":\"world\"}} and http.Post(ingest.URL+\"/webhook\", \"application/json\", the body). Expect 202."
  - "Poll for up to ~2 seconds (loop with short time.Sleep) until the destination's hit counter reaches 1, then print how many deliveries the destination received and the queue length. main must run to completion and exit normally (no panic); it should print a clear line like \"delivered N webhook(s)\"."
constraints: package main; standard library only (context, net/http, net/http/httptest, encoding/json, sync/atomic, time, fmt, bytes); uses Server, Worker, Queue, Deliverer, Job
