name: Server
role: component
intent: HTTP layer - accept a webhook POST, cap the body, enqueue, reply immediately
api:
  - type Server struct holding a *Dispatcher (field disp)
  - func NewServer(d *Dispatcher) *Server
  - method (*Server) Webhook(w http.ResponseWriter, r *http.Request):
      * if method != POST -> 405 and return
      * cap the body: r.Body = http.MaxBytesReader(w, r.Body, 1<<20) BEFORE reading, so a huge POST
        cannot exhaust memory; read with io.ReadAll; on error -> 413 Request Entity Too Large and return
      * Enqueue an Event; if the queue is full -> 503 and return; else 202 Accepted immediately
  - method (*Server) Routes() http.Handler   // *http.ServeMux, HandleFunc("/webhook", s.Webhook)
constraints: standard library (net/http, io); package main; uses Dispatcher + Event verbatim
