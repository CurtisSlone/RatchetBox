name: HTTPServer
role: behavior
intent: HTTP API with POST /shorten and GET /{code} redirect
api:
  - func NewHTTPServer(shortener *URLShortener) *HTTPServer
  - func (s *HTTPServer) Start(addr string) error
  - func (s *HTTPServer) Shutdown() error
behavior:
  - POST /shorten should accept a JSON body with "url" field and return a short code
  - GET /{code} should redirect to the original URL with 302 status
  - Should handle request size limits and timeouts
  - Should gracefully shutdown on SIGINT/SIGTERM
constraints: Use http.Server with ReadHeaderTimeout, ReadTimeout, WriteTimeout, IdleTimeout; use http.MaxBytesReader; package: main
