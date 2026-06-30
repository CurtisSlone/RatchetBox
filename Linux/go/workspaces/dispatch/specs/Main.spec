name: Main
role: behavior
intent: The runnable entry point that wires the whole dispatcher together and demonstrates it end to end - including a deliberately flaky destination so the circuit breaker and retries actually fire.
api:
  - func main()
behavior:
  - "Define a demo Deliverer (a struct implementing Deliverer) whose Deliver fails the first two calls (returns an error) and then succeeds, so retries and the breaker are exercised. Guard its internal counter with a mutex (it is called concurrently)."
  - "Wire the components: log := NewEventLog(); breaker := NewBreaker(3, 50*time.Millisecond); policy := NewRetryPolicy(1*time.Millisecond, 2.0, 5); d := NewDispatcher(log, breaker, policy, deliverer)."
  - "Build a few Jobs with NewJob (distinct ids and urls), call d.DispatchAll(jobs, 4), then print each returned job's ID and State, and finally print the number of events in the log (len(log.Events()))."
  - "main must run to completion without panicking and exit normally."
constraints: package main; standard library only; uses Job, EventLog, Breaker, RetryPolicy, Dispatcher
