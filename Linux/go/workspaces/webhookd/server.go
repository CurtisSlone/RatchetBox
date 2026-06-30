package main

// file-kw: server http ingest surface now durable idempotent observable same accept idempotency durability behavior plus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
)

// Server is the HTTP ingest surface, now durable and idempotent.
// kw: server http ingest surface
type Server struct {
	queue         *Queue
	wal           *Wal
	metrics       *Metrics
	counter       atomic.Int64
	idempotency   map[string]string
	idempotencyMu sync.Mutex
}

// NewServer creates a new server with the given queue, wal, and metrics.
// kw: server queue wal metrics http ingest surface
func NewServer(q *Queue, wal *Wal, metrics *Metrics) *Server {
	return &Server{
		queue:       q,
		wal:         wal,
		metrics:     metrics,
		idempotency: make(map[string]string),
	}
}

// Handler returns an HTTP handler for the server.
// kw: handler server http ingest surface
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /webhook", s.handleWebhook)
	mux.Handle("GET /metrics", s.metrics.Handler())
	return mux
}

// handleWebhook handles incoming webhook requests.
// kw: handle webhook server http response writer request ingest surface
func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB limit

	var body struct {
		URL     string          `json:"url"`
		Payload json.RawMessage `json:"payload"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if body.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	key := r.Header.Get("Idempotency-Key")
	if key != "" {
		s.idempotencyMu.Lock()
		if existingID, exists := s.idempotency[key]; exists {
			s.idempotencyMu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"id": existingID})
			return
		}
		s.idempotencyMu.Unlock()
	}

	id := fmt.Sprintf("job-%d", s.counter.Add(1))
	j := NewJob(id, body.URL, []byte(body.Payload))

	if err := s.wal.Append(j); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if key != "" {
		s.idempotencyMu.Lock()
		s.idempotency[key] = id
		s.idempotencyMu.Unlock()
	}

	s.queue.Push(j)
	s.metrics.IncAccepted()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}
