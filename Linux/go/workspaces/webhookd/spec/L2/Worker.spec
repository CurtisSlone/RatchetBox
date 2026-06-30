name: Worker
role: component
dependsOn: Queue, Deliverer, Breaker, RetryPolicy, Job
intent: The delivery loop, now RESILIENT. It pulls pending jobs off the queue and delivers each through a circuit breaker plus exponential-backoff retry, so a flaky or briefly-down destination is retried instead of dropped, and a destination that keeps failing trips the breaker so we stop hammering it. Updating job state as before.
api:
  - type Worker struct { ... }
  - func NewWorker(q *Queue, d Deliverer, breaker *Breaker, policy RetryPolicy) *Worker
  - func (w *Worker) RunOnce(ctx context.Context) bool
  - func (w *Worker) Run(ctx context.Context)
behavior:
  - "Worker holds the *Queue, the Deliverer, a *Breaker, and a RetryPolicy. NewWorker stores all four."
  - "RunOnce pops one job. If the queue is empty (Pop returns false) it returns false. Otherwise it runs the RESILIENT DELIVERY loop below and returns true when done (a job was processed)."
  - "RESILIENT DELIVERY (pinned): loop attempts until delivered or give up:
      1. j.Attempts++ (1-based attempt counter on the job).
      2. if !w.breaker.Allow() -> the destination is considered unhealthy: set j.State = StateFailed and stop (break).
      3. err := w.d.Deliver(ctx, j).
      4. on nil err -> w.breaker.Success(); j.State = StateDelivered; stop (break).
      5. on err -> w.breaker.Failure(). If !w.policy.ShouldRetry(j.Attempts) -> set j.State = StateFailed and stop (break). Otherwise sleep w.policy.Backoff(j.Attempts) (respecting ctx: if ctx is cancelled during the wait, stop) and loop to retry."
  - "Run loops: while ctx.Err() == nil, call RunOnce; if RunOnce returned false (queue empty) sleep a short poll interval (10 * time.Millisecond) before looping again. Return when the context is cancelled. Run must be safe to start in its own goroutine."
constraints: package main; standard library only (context, time); uses Queue, Deliverer, Breaker, RetryPolicy, Job
