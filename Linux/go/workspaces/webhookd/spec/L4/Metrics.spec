name: Metrics
role: component
dependsOn: ""
intent: Operational counters for the dispatcher, exposed over HTTP so the service is observable in production. Four monotonic counters - accepted, delivered, failed, dead-lettered - incremented from many goroutines, plus a /metrics-style text handler and a snapshot for printing.
api:
  - type Metrics struct { ... }
  - func NewMetrics() *Metrics
  - func (m *Metrics) IncAccepted()
  - func (m *Metrics) IncDelivered()
  - func (m *Metrics) IncFailed()
  - func (m *Metrics) IncDeadLettered()
  - func (m *Metrics) Snapshot() map[string]int64
  - func (m *Metrics) Handler() http.Handler
behavior:
  - "Hold four counters as atomic int64 values (sync/atomic; either four atomic.Int64 fields or four int64 fields incremented with atomic.AddInt64). They are incremented concurrently by the worker and the HTTP handlers, so all access MUST be atomic."
  - "Each Inc* method atomically adds 1 to its counter."
  - "Snapshot atomically Loads all four and returns a map[string]int64 with keys exactly: \"accepted\", \"delivered\", \"failed\", \"dead_lettered\"."
  - "Handler returns an http.HandlerFunc that writes the counters in Prometheus text exposition format - one line per counter like \"webhookd_accepted N\\n\", \"webhookd_delivered N\\n\", \"webhookd_failed N\\n\", \"webhookd_dead_lettered N\\n\" - using the snapshot values. Content-Type text/plain."
constraints: package main; standard library only (net/http, sync/atomic, fmt); no dependencies on other units
