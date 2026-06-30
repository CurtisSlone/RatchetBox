name: Server
role: component
dependsOn: Queue, Wal, Job
intent: The HTTP ingest surface, now made DURABLE. It accepts a webhook, appends it to the write-ahead log FIRST, and only then enqueues it - so an accepted (202) webhook survives a crash. A client POSTs JSON describing where to deliver and what; the server creates a Job, persists it, enqueues it, and returns 202 Accepted.
api:
  - type Server struct { ... }
  - func NewServer(q *Queue, wal *Wal) *Server
  - func (s *Server) Handler() http.Handler
behavior:
  - "Server holds the *Queue, the *Wal, plus an atomic int64 sequence counter (sync/atomic) used to mint unique job IDs. NewServer stores the queue and the wal."
  - "Handler returns an *http.ServeMux with one route: POST /webhook."
  - "The /webhook handler: reject non-POST with http.StatusMethodNotAllowed. Decode the JSON body into a struct { URL string `json:\"url\"`; Payload json.RawMessage `json:\"payload\"` }. On a decode error or empty URL, respond http.StatusBadRequest."
  - "Mint an id with atomic.AddInt64 on the counter, formatted as fmt.Sprintf(\"job-%d\", n). Create j := NewJob(id, body.URL, []byte(body.Payload))."
  - "DURABILITY ORDER (pinned): call s.wal.Append(j) FIRST. If it returns an error, respond http.StatusInternalServerError (the webhook is NOT accepted because it could not be made durable) and return. Only on a successful append do q.Push(j) and then write http.StatusAccepted (202) with the job id as the response body."
constraints: package main; standard library only (net/http, encoding/json, sync/atomic, fmt); uses Queue, Wal, Job
