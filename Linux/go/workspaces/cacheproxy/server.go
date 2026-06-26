package main

import (
	"net/http"
)

type HTTPServer struct {
	proxy *Proxy
}

func NewHTTPServer(proxy *Proxy) *HTTPServer {
	return &HTTPServer{
		proxy: proxy,
	}
}

func (s *HTTPServer) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /get/{key}", s.handleGet)
	return mux
}

func (s *HTTPServer) handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	value, err := s.proxy.Get(key)
	if err != nil {
		http.Error(w, "Proxy error", http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}
