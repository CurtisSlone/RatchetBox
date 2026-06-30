name: Worker
role: component
dependsOn: Queue, Deliverer, Job
intent: The delivery loop - pull pending jobs off the queue and deliver each via the Deliverer, updating job state. In the walking skeleton it delivers exactly once with no retry; a later layer wraps this in a circuit breaker and exponential-backoff retry, and another adds a dead-letter path. Keeping the loop in its own unit is what makes those later changes a localized diff.
api:
  - type Worker struct { ... }
  - func NewWorker(q *Queue, d Deliverer) *Worker
  - func (w *Worker) RunOnce(ctx context.Context) bool
  - func (w *Worker) Run(ctx context.Context)
behavior:
  - "Worker holds the *Queue and the Deliverer. NewWorker stores both."
  - "RunOnce pops one job. If the queue is empty (Pop returns false) it returns false (nothing to do). Otherwise it increments j.Attempts, calls d.Deliver(ctx, j); on nil error sets j.State = StateDelivered, on error sets j.State = StateFailed. It returns true (a job was processed) regardless of delivery outcome."
  - "Run loops: while ctx.Err() == nil, call RunOnce; if RunOnce returned false (queue empty) sleep a short poll interval (10 * time.Millisecond) before looping again. Return when the context is cancelled."
  - "Run must be safe to start in its own goroutine."
constraints: package main; standard library only (context, time); uses Queue, Deliverer, Job
