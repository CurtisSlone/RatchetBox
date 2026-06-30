name: Server
role: component
dependsOn: Queue, Job
intent: The HTTP ingest surface - accept webhook delivery requests and enqueue them. A client POSTs JSON describing where to deliver and what; the server creates a Job, enqueues it, and returns 202 Accepted immediately (delivery is asynchronous, done by the worker). This is the front door of the service.
api:
  - type Server struct { ... }
  - func NewServer(q *Queue) *Server
  - func (s *Server) Handler() http.Handler
behavior:
  - "Server holds the *Queue plus an atomic int64 sequence counter (sync/atomic) used to mint unique job IDs. NewServer stores the queue."
  - "Handler returns an *http.ServeMux with one route: POST /webhook."
  - "The /webhook handler: reject non-POST with http.StatusMethodNotAllowed. Decode the JSON body into a struct { URL string `json:\"url\"`; Payload json.RawMessage `json:\"payload\"` }. On a decode error or empty URL, respond http.StatusBadRequest."
  - "Mint an id with atomic.AddInt64 on the counter, formatted as fmt.Sprintf(\"job-%d\", n). Create j := NewJob(id, body.URL, []byte(body.Payload)), q.Push(j), then write http.StatusAccepted (202) and the job id as the response body."
constraints: package main; standard library only (net/http, encoding/json, sync/atomic, fmt); uses Queue, Job
