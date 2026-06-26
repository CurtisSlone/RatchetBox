name: Server
role: component
intent: the HTTP layer - accept a webhook POST, hand it to the Dispatcher, and reply immediately
api:
  - type Server struct holding a *Dispatcher (field name: disp)
  - func NewServer(d *Dispatcher) *Server
  - method (*Server) Webhook(w http.ResponseWriter, r *http.Request)   // the handler:
      * if r.Method is not POST, http.Error with status 405 and return
      * read the body with io.ReadAll(r.Body); on error, 400 and return
      * build an Event with NewEvent (any short id is fine, e.g. the time or a counter)
      * if s.disp.Enqueue(event) is false (queue full), reply 503 Service Unavailable and return
      * otherwise reply 202 Accepted immediately (do NOT process inline) and write a tiny body
  - method (*Server) Routes() http.Handler   // a *http.ServeMux with HandleFunc("/webhook", s.Webhook)
behavior:
  - the handler must return right after enqueuing - the actual work happens in the Dispatcher workers,
    which is what keeps request latency low
constraints: standard library only (net/http, io); package main; no func main in this file; uses the
  existing Dispatcher and Event API verbatim
