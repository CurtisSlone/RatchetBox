package main

// file-kw: http server rate limiter middleware

import (
	"net"
	"net/http"
)

// kw: http server rate limit middleware
type Server struct {
	limiter Limiter
}

// kw: create new http server with rate limiter
func NewServer(limiter Limiter) *Server {
	return &Server{
		limiter: limiter,
	}
}

// kw: http middleware for rate limiting
func (s *Server) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientKey, err := getClientKey(r)
		if err != nil {
			// If we can't determine the client key, allow the request to proceed
			// to avoid blocking legitimate requests due to misconfiguration
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := s.limiter.Allow(r.Context(), clientKey)
		if err != nil {
			// Fail-open: if there's an error from the limiter, allow the request
			// to proceed to avoid blocking legitimate requests due to system issues
			next.ServeHTTP(w, r)
			return
		}

		if !allowed {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// kw: http server handler with default endpoint
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	return s.Middleware(mux)
}

// kw: extract client key from http request
func getClientKey(r *http.Request) (string, error) {
	remoteAddr := r.RemoteAddr
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// If SplitHostPort fails, use the raw RemoteAddr
		// This handles cases where there's no port (e.g., unix sockets)
		return remoteAddr, nil
	}
	return host, nil
}
