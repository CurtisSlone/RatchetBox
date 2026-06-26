package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HTTPServer struct {
	shortener *URLShortener
	server    *http.Server
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Code string `json:"code"`
}

func NewHTTPServer(shortener *URLShortener) *HTTPServer {
	mux := http.NewServeMux()
	s := &HTTPServer{
		shortener: shortener,
		server: &http.Server{
			Addr:              ":8080",
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}

	mux.HandleFunc("POST /shorten", s.handleShorten)
	mux.HandleFunc("GET /{code}", s.handleExpand)

	return s
}

func (s *HTTPServer) Start(addr string) error {
	s.server.Addr = addr
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) handleShorten(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB limit
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "request too large", http.StatusRequestEntityTooLarge)
		return
	}

	var req ShortenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "url field is required", http.StatusBadRequest)
		return
	}

	code, err := s.shortener.Shorten(req.URL)
	if err != nil {
		http.Error(w, "failed to shorten URL", http.StatusInternalServerError)
		return
	}

	resp := ShortenResponse{Code: code}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *HTTPServer) handleExpand(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		http.Error(w, "code is required", http.StatusBadRequest)
		return
	}

	url, found := s.shortener.Expand(code)
	if !found {
		http.Error(w, "code not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func main() {
	_ = NewURLStore()
	encoder := NewBase62Encoder()
	shortener := NewURLShortener()
	server := NewHTTPServer(shortener)

	// Listen for SIGINT/SIGTERM; ctx is cancelled on signal.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("listening on", server.server.Addr)
		if err := server.Start(server.server.Addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done() // wait for the signal
	stop()
	if err := server.Shutdown(); err != nil {
		log.Println("shutdown:", err)
	}
}
