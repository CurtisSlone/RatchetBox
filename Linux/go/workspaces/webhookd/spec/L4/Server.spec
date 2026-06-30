name: Server
role: component
dependsOn: Queue, Wal, Metrics, Job
intent: The HTTP ingest surface, now DURABLE, IDEMPOTENT, and OBSERVABLE. Same accept/idempotency/durability behavior, plus it counts accepted webhooks and exposes the metrics endpoint so the whole service is scrapeable from one server.
api:
  - type Server struct { ... }
  - func NewServer(q *Queue, wal *Wal, metrics *Metrics) *Server
  - func (s *Server) Handler() http.Handler
behavior:
  - "Server holds the *Queue, the *Wal, the *Metrics, an atomic int64 sequence counter (sync/atomic) for unique job IDs, AND an idempotency map[string]string guarded by a sync.Mutex. NewServer stores the queue, wal, and metrics and initialises the map."
  - "Handler returns an *http.ServeMux with TWO routes: POST /webhook (the ingest handler below) and GET /metrics (serve s.metrics.Handler() - mux.Handle(\"GET /metrics\", s.metrics.Handler()))."
  - "The /webhook handler: reject non-POST with http.StatusMethodNotAllowed. Decode the JSON body into a struct { URL string `json:\"url\"`; Payload json.RawMessage `json:\"payload\"` }. On a decode error or empty URL, respond http.StatusBadRequest."
  - "IDEMPOTENCY (pinned): read key := r.Header.Get(\"Idempotency-Key\"). If key != \"\" and already present in the map, write http.StatusOK (200) with the EXISTING job id and return (no append, no enqueue, no accepted-count)."
  - "Otherwise mint an id with atomic.AddInt64, formatted fmt.Sprintf(\"job-%d\", n). Create j := NewJob(id, body.URL, []byte(body.Payload))."
  - "DURABILITY ORDER (pinned): call s.wal.Append(j) FIRST. If it errors, respond http.StatusInternalServerError and return. On success: if key != \"\" record key->id in the map; q.Push(j); s.metrics.IncAccepted(); then write http.StatusAccepted (202) with the job id."
constraints: package main; standard library only (net/http, encoding/json, sync, sync/atomic, fmt); uses Queue, Wal, Metrics, Job
