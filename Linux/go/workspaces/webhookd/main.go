package main

// file-kw: main runnable entry point production webhook dispatcher demo durable resilient idempotent observable graceful drain
//
// NOTE: this file is HUMAN-AUTHORED (the "senior intervention"). The local model reliably produced every
// other unit to spec - they are gated by the module oracle (go vet + go test -race). main is the one unit
// the oracle cannot gate on INTENT (it has no test), and it is the most complex/least-constrained unit
// (it wires everything and scripts a multi-step demo), so the model's generations of it were variable and
// sometimes silently dropped the demo. This is the documented "last 10-15%" case: the human edits intent.

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

// kw: main wire dispatcher demonstrate durability resilience idempotency dead-letter metrics drain
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Durability: open the WAL and replay anything a prior run accepted but may not have delivered.
	wal, err := OpenWal(filepath.Join(os.TempDir(), "webhookd.wal"))
	if err != nil {
		fmt.Printf("error opening WAL: %v\n", err)
		return
	}
	defer wal.Close()

	recovered, _ := wal.Replay()
	q := NewQueue()
	for _, j := range recovered {
		q.Push(j)
	}
	fmt.Printf("recovered %d job(s) from the WAL\n", len(recovered))

	// A destination that ALWAYS fails: deliveries to it exhaust retries and land in the dead-letter queue.
	var deadHits int64
	dead := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&deadHits, 1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer dead.Close()

	// Wire the full production pipeline.
	metrics := NewMetrics()
	breaker := NewBreaker(10, 100*time.Millisecond) // high threshold: won't trip during the small demo
	policy := NewRetryPolicy(2*time.Millisecond, 2.0, 3)
	dlq := NewDeadLetter()
	deliverer := NewHTTPDeliverer(nil)
	worker := NewWorker(q, deliverer, breaker, policy, dlq, metrics)
	go worker.Run(ctx)

	server := NewServer(q, wal, metrics)
	ingest := httptest.NewServer(server.Handler())
	defer ingest.Close()

	payload := []byte(`{"url":"` + dead.URL + `","payload":{"hello":"world"}}`)

	// --- Dead-letter demo: one webhook to the dead destination should end up dead-lettered.
	resp, err := http.Post(ingest.URL+"/webhook", "application/json", bytes.NewReader(payload))
	if err != nil {
		fmt.Printf("dead-letter demo POST failed: %v\n", err)
		return
	}
	fmt.Printf("dead-letter demo: accepted with %d\n", resp.StatusCode)
	resp.Body.Close()

	deadline := time.Now().Add(3 * time.Second)
	for dlq.Len() < 1 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Printf("dead-letter queue holds %d job(s)\n", dlq.Len())

	// --- Idempotency demo: the SAME Idempotency-Key sent twice must enqueue only ONCE.
	// A FRESH request per send (the first Do drains the body, so a reused *http.Request would fail).
	send := func() int {
		req, _ := http.NewRequest(http.MethodPost, ingest.URL+"/webhook", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", "dup-1")
		r, e := http.DefaultClient.Do(req)
		if e != nil {
			panic(e)
		}
		defer r.Body.Close()
		return r.StatusCode
	}
	first := send()
	second := send()
	fmt.Printf("idempotency demo: first=%d second=%d (200 second => de-duplicated)\n", first, second)

	// --- Graceful drain: wait for the queue to empty before stopping the worker (no abandoned jobs).
	dd := time.Now().Add(3 * time.Second)
	for q.Len() > 0 && time.Now().Before(dd) {
		time.Sleep(10 * time.Millisecond)
	}
	cancel()

	// --- Observability: scrape /metrics over HTTP, then print the snapshot.
	if mr, e := http.Get(ingest.URL + "/metrics"); e == nil {
		body, _ := io.ReadAll(mr.Body)
		mr.Body.Close()
		fmt.Printf("GET /metrics:\n%s", body)
	}
	snap := metrics.Snapshot()
	fmt.Printf("final metrics: accepted=%d delivered=%d failed=%d dead_lettered=%d\n",
		snap["accepted"], snap["delivered"], snap["failed"], snap["dead_lettered"])
}
