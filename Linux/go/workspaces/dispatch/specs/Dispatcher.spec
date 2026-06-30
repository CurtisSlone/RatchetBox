name: Dispatcher
role: component
intent: The orchestrator. It delivers jobs to an external destination through a pluggable Deliverer, recording every outcome in the append-only EventLog, gating calls with the circuit Breaker, and retrying failures per the RetryPolicy. It can process many jobs concurrently with a fixed worker pool.
api:
  - type Deliverer interface { Deliver(url string, payload []byte) error }
  - type Dispatcher struct { ... }
  - func NewDispatcher(log *EventLog, breaker *Breaker, policy RetryPolicy, deliverer Deliverer) *Dispatcher
  - func (d *Dispatcher) Dispatch(job Job) Job
  - func (d *Dispatcher) DispatchAll(jobs []Job, workers int) []Job
behavior:
  - "Deliverer is the injectable external call (so it is testable with a fake). Deliver returns nil on success, a non-nil error on failure."
  - "Dispatch delivers ONE job and returns the final job (with updated Attempts and State). Algorithm, attempt counting from 1:
      - If breaker.Allow() is false: do NOT call the deliverer; append a failed Event for this job (Attempt = current attempt) and return the job unchanged (still pending) - it was not delivered and is not dead.
      - Otherwise call deliverer.Deliver(job.URL, job.Payload):
          success -> breaker.Success(); append a delivered Event; set job.State = StateDelivered; set job.Attempts = attempt; return.
          failure -> breaker.Failure(); append a failed Event; set job.Attempts = attempt; if policy.ShouldRetry(attempt) is true, wait policy.Backoff(attempt) then try again (attempt+1); otherwise set job.State = StateDead and return."
  - "On enqueue: Dispatch should append one enqueued Event for the job before its first attempt."
  - "DispatchAll processes the given jobs concurrently using exactly `workers` goroutines (a fixed worker pool - feed jobs over a channel, collect results), and returns the final jobs (order need not match input). If workers < 1, treat it as 1. The EventLog and Breaker are safe for concurrent use, so the shared breaker and log are used directly."
  - "CONCURRENCY: DispatchAll must be race-free under go test -race. Use a sync.WaitGroup for the pool."
constraints: package main; standard library only; uses Job, EventLog, Breaker, RetryPolicy
