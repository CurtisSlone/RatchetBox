package main

import (
	"io"
	"net/http"
)

// Server handles HTTP webhook requests
type Server struct {
	disp *Dispatcher
}

// NewServer creates a new Server with the given dispatcher
func NewServer(d *Dispatcher) *Server {
	return &Server{disp: d}
}

// Webhook handles incoming webhook POST requests
func (s *Server) Webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request Entity Too Large", http.StatusRequestEntityTooLarge)
		return
	}

	event := NewEvent("", body)
	if !s.disp.Enqueue(event) {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// Routes returns the HTTP handler for the server
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", s.Webhook)
	return mux
}
