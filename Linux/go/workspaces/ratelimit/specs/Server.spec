name: Server
role: component
intent: The HTTP layer. It wraps any handler with rate-limit middleware that consults a Limiter, keyed by the client's IP, and rejects over-limit requests with 429 Too Many Requests.
api:
  - type Server struct { ... }
  - func NewServer(limiter Limiter) *Server
  - func (s *Server) Middleware(next http.Handler) http.Handler
  - func (s *Server) Handler() http.Handler
behavior:
  - "NewServer stores the Limiter (the interface, so either implementation works)."
  - "Middleware returns an http.Handler that, for each request, derives the client key from the request (use r.RemoteAddr - take the host part via net.SplitHostPort, falling back to the raw RemoteAddr), calls limiter.Allow(r.Context(), clientKey), and:
      - allowed true  -> call next.ServeHTTP(w, r).
      - allowed false -> write http.StatusTooManyRequests (429) with a short \"rate limit exceeded\" body and a Retry-After: 1 header.
      - on a non-nil error from Allow -> the limiter already decided allowed (fail-open), so honor the bool: if allowed, proceed; the middleware does not itself turn an error into a rejection."
  - "Handler returns a ready-to-serve handler: a simple OK endpoint (200, body \"ok\") wrapped in Middleware, so the server is runnable on its own."
  - "Middleware must be safe for concurrent use (it only reads the limiter, which is concurrency-safe)."
constraints: package main; uses Limiter; standard library (net/http, net)
