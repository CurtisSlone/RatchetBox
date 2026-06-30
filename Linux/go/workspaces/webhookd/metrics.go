package main

// file-kw: metrics operational counters dispatcher exposed over http service observable production four monotonic accepted delivered

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// Metrics holds operational counters for the dispatcher.
// kw: metrics operational counters dispatcher
type Metrics struct {
	accepted     atomic.Int64
	delivered    atomic.Int64
	failed       atomic.Int64
	deadLettered atomic.Int64
}

// NewMetrics creates a new Metrics instance.
// kw: metrics operational counters dispatcher
func NewMetrics() *Metrics {
	return &Metrics{}
}

// IncAccepted increments the accepted counter.
// kw: inc accepted metrics operational counters dispatcher
func (m *Metrics) IncAccepted() {
	m.accepted.Add(1)
}

// IncDelivered increments the delivered counter.
// kw: inc delivered metrics operational counters dispatcher
func (m *Metrics) IncDelivered() {
	m.delivered.Add(1)
}

// IncFailed increments the failed counter.
// kw: inc failed metrics operational counters dispatcher
func (m *Metrics) IncFailed() {
	m.failed.Add(1)
}

// IncDeadLettered increments the dead-lettered counter.
// kw: inc dead lettered metrics operational counters dispatcher
func (m *Metrics) IncDeadLettered() {
	m.deadLettered.Add(1)
}

// Snapshot returns a copy of all counters.
// kw: snapshot metrics operational counters dispatcher
func (m *Metrics) Snapshot() map[string]int64 {
	return map[string]int64{
		"accepted":      m.accepted.Load(),
		"delivered":     m.delivered.Load(),
		"failed":        m.failed.Load(),
		"dead_lettered": m.deadLettered.Load(),
	}
}

// Handler returns an HTTP handler that serves metrics in Prometheus text format.
// kw: handler metrics http operational counters dispatcher
func (m *Metrics) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		snap := m.Snapshot()
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "webhookd_accepted %d\n", snap["accepted"])
		fmt.Fprintf(w, "webhookd_delivered %d\n", snap["delivered"])
		fmt.Fprintf(w, "webhookd_failed %d\n", snap["failed"])
		fmt.Fprintf(w, "webhookd_dead_lettered %d\n", snap["dead_lettered"])
	})
}
