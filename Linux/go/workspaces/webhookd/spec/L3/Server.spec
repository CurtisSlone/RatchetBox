name: Server
role: component
dependsOn: Queue, Wal, Job
intent: The HTTP ingest surface, now DURABLE and IDEMPOTENT. A client may send an Idempotency-Key header; if the same key is replayed (e.g. a client retry after a network blip), the server returns the original job id WITHOUT enqueuing a duplicate delivery. Otherwise it persists to the WAL, enqueues, and returns 202 as before.
api:
  - type Server struct { ... }
  - func NewServer(q *Queue, wal *Wal) *Server
  - func (s *Server) Handler() http.Handler
behavior:
  - "Server holds the *Queue, the *Wal, an atomic int64 sequence counter (sync/atomic) for unique job IDs, AND an idempotency map[string]string (idempotency-key -> job id) guarded by a sync.Mutex. NewServer stores the queue and wal and initialises the map."
  - "Handler returns an *http.ServeMux with one route: POST /webhook."
  - "The /webhook handler: reject non-POST with http.StatusMethodNotAllowed. Decode the JSON body into a struct { URL string `json:\"url\"`; Payload json.RawMessage `json:\"payload\"` }. On a decode error or empty URL, respond http.StatusBadRequest."
  - "IDEMPOTENCY (pinned): read key := r.Header.Get(\"Idempotency-Key\"). If key != \"\", lock the mutex and check the map: if the key is already present, write http.StatusOK (200) with the EXISTING job id and return (no WAL append, no enqueue - this is a duplicate). Keep the lock only as needed."
  - "Otherwise mint an id with atomic.AddInt64, formatted fmt.Sprintf(\"job-%d\", n). Create j := NewJob(id, body.URL, []byte(body.Payload))."
  - "DURABILITY ORDER (pinned): call s.wal.Append(j) FIRST. If it errors, respond http.StatusInternalServerError and return. On success: if key != \"\" record key->id in the idempotency map (under the mutex); then q.Push(j) and write http.StatusAccepted (202) with the job id as the body."
constraints: package main; standard library only (net/http, encoding/json, sync, sync/atomic, fmt); uses Queue, Wal, Job
