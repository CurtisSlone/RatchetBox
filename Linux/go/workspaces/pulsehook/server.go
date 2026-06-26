package main

import (
	"io"
	"net/http"
	"time"
)

type Server struct {
	disp *Dispatcher
}

func NewServer(d *Dispatcher) *Server {
	return &Server{disp: d}
}

func (s *Server) Webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	event := NewEvent(time.Now().String(), body)

	if !s.disp.Enqueue(event) {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("accepted"))
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", s.Webhook)
	return mux
}
