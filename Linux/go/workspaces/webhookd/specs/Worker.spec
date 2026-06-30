name: Worker
role: component
dependsOn: Queue, Deliverer, Breaker, RetryPolicy, DeadLetter, Metrics, Job
intent: The delivery loop, now resilient, dead-lettering, AND instrumented. Same delivery logic as before, but every terminal outcome bumps an operational counter (delivered / failed / dead-lettered) so the service is observable.
api:
  - type Worker struct { ... }
  - func NewWorker(q *Queue, d Deliverer, breaker *Breaker, policy RetryPolicy, dlq *DeadLetter, metrics *Metrics) *Worker
  - func (w *Worker) RunOnce(ctx context.Context) bool
  - func (w *Worker) Run(ctx context.Context)
behavior:
  - "Worker holds the *Queue, the Deliverer, a *Breaker, a RetryPolicy, a *DeadLetter, and a *Metrics. NewWorker stores all six."
  - "RunOnce pops one job. If the queue is empty (Pop returns false) it returns false. Otherwise it runs the RESILIENT DELIVERY loop below and returns true when done."
  - "RESILIENT DELIVERY (pinned): loop attempts until delivered or give up:
      1. j.Attempts++ (1-based).
      2. if !w.breaker.Allow() -> set j.State = StateFailed; w.dlq.Add(j); w.metrics.IncFailed(); w.metrics.IncDeadLettered(); stop (break).
      3. err := w.d.Deliver(ctx, j).
      4. on nil err -> w.breaker.Success(); j.State = StateDelivered; w.metrics.IncDelivered(); stop (break).
      5. on err -> w.breaker.Failure(). If !w.policy.ShouldRetry(j.Attempts) -> set j.State = StateFailed; w.dlq.Add(j); w.metrics.IncFailed(); w.metrics.IncDeadLettered(); stop (break). Otherwise sleep w.policy.Backoff(j.Attempts) (respecting ctx) and loop to retry."
  - "So exactly one terminal counter is bumped per job: IncDelivered on success, or IncFailed+IncDeadLettered when it gives up undelivered."
  - "Run loops: while ctx.Err() == nil, call RunOnce; if RunOnce returned false (queue empty) sleep 10 * time.Millisecond before looping. Return when the context is cancelled. Run must be safe to start in its own goroutine."
constraints: package main; standard library only (context, time); uses Queue, Deliverer, Breaker, RetryPolicy, DeadLetter, Metrics, Job
