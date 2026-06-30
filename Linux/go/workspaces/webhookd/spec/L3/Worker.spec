name: Worker
role: component
dependsOn: Queue, Deliverer, Breaker, RetryPolicy, DeadLetter, Job
intent: The delivery loop, now resilient AND with a dead-letter path. It delivers through a circuit breaker plus exponential-backoff retry; a job that exhausts its retries or is rejected by an open breaker is no longer silently left failed - it is parked in the dead-letter queue so nothing is lost.
api:
  - type Worker struct { ... }
  - func NewWorker(q *Queue, d Deliverer, breaker *Breaker, policy RetryPolicy, dlq *DeadLetter) *Worker
  - func (w *Worker) RunOnce(ctx context.Context) bool
  - func (w *Worker) Run(ctx context.Context)
behavior:
  - "Worker holds the *Queue, the Deliverer, a *Breaker, a RetryPolicy, and a *DeadLetter. NewWorker stores all five."
  - "RunOnce pops one job. If the queue is empty (Pop returns false) it returns false. Otherwise it runs the RESILIENT DELIVERY loop below and returns true when done."
  - "RESILIENT DELIVERY (pinned): loop attempts until delivered or give up:
      1. j.Attempts++ (1-based).
      2. if !w.breaker.Allow() -> set j.State = StateFailed; w.dlq.Add(j); stop (break).
      3. err := w.d.Deliver(ctx, j).
      4. on nil err -> w.breaker.Success(); j.State = StateDelivered; stop (break) (do NOT dead-letter a delivered job).
      5. on err -> w.breaker.Failure(). If !w.policy.ShouldRetry(j.Attempts) -> set j.State = StateFailed; w.dlq.Add(j); stop (break). Otherwise sleep w.policy.Backoff(j.Attempts) (respecting ctx) and loop to retry."
  - "So a job lands in the dead-letter queue EXACTLY when it gives up undelivered (open breaker, or retries exhausted) - never when it was delivered."
  - "Run loops: while ctx.Err() == nil, call RunOnce; if RunOnce returned false (queue empty) sleep 10 * time.Millisecond before looping. Return when the context is cancelled. Run must be safe to start in its own goroutine."
constraints: package main; standard library only (context, time); uses Queue, Deliverer, Breaker, RetryPolicy, DeadLetter, Job
