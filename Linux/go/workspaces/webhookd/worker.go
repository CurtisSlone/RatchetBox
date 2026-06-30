package main

// file-kw: worker delivery loop now resilient dead lettering instrumented same logic before but terminal outcome

import (
	"context"
	"time"
)

// Worker is a resilient delivery loop that pulls jobs from a queue and delivers them.
// kw: worker delivery loop now
type Worker struct {
	queue   *Queue
	d       Deliverer
	breaker *Breaker
	policy  RetryPolicy
	dlq     *DeadLetter
	metrics *Metrics
}

// NewWorker creates a new worker with the given components.
// kw: worker queue deliverer breaker policy retry dlq dead letter
func NewWorker(q *Queue, d Deliverer, breaker *Breaker, policy RetryPolicy, dlq *DeadLetter, metrics *Metrics) *Worker {
	return &Worker{
		queue:   q,
		d:       d,
		breaker: breaker,
		policy:  policy,
		dlq:     dlq,
		metrics: metrics,
	}
}

// RunOnce pops one job from the queue and processes it with resilient delivery.
// kw: run once worker ctx context delivery loop now
func (w *Worker) RunOnce(ctx context.Context) bool {
	j, ok := w.queue.Pop()
	if !ok {
		return false
	}

	// RESILIENT DELIVERY LOOP
	for {
		j.Attempts++
		if !w.breaker.Allow() {
			j.State = StateFailed
			w.dlq.Add(j)
			w.metrics.IncFailed()
			w.metrics.IncDeadLettered()
			break
		}

		err := w.d.Deliver(ctx, j)
		if err == nil {
			w.breaker.Success()
			j.State = StateDelivered
			w.metrics.IncDelivered()
			break
		}

		w.breaker.Failure()
		if !w.policy.ShouldRetry(j.Attempts) {
			j.State = StateFailed
			w.dlq.Add(j)
			w.metrics.IncFailed()
			w.metrics.IncDeadLettered()
			break
		}

		// Sleep with context cancellation support
		select {
		case <-time.After(w.policy.Backoff(j.Attempts)):
		case <-ctx.Done():
			return false
		}
	}

	return true
}

// Run runs the worker loop until the context is cancelled.
// kw: run worker ctx context delivery loop now
func (w *Worker) Run(ctx context.Context) {
	for ctx.Err() == nil {
		if !w.RunOnce(ctx) {
			time.Sleep(10 * time.Millisecond)
		}
	}
}
