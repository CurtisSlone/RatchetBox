name: Breaker
role: component
intent: A circuit breaker that stops calls to a failing destination so failures do not cascade. It is a three-state machine (closed, open, half-open) driven by recorded successes and failures. Safe for concurrent use.
api:
  - type Breaker struct { ... }
  - func NewBreaker(maxFailures int, cooldown time.Duration) *Breaker
  - func (b *Breaker) Allow() bool
  - func (b *Breaker) Success()
  - func (b *Breaker) Failure()
  - func (b *Breaker) State() string
behavior:
  - "STATES are exactly: \"closed\", \"open\", \"half-open\". A new breaker starts \"closed\" with a failure count of 0."
  - "TIME TRANSITION (applied at the START of Allow, Success, AND Failure - not only Allow): if the breaker is open and at least cooldown has elapsed since it opened, it transitions to half-open FIRST. This makes the open->half-open transition observable no matter which method is called after the cooldown - the state machine does not depend on Allow being called. (A private helper called at the top of each method is the natural way.)"
  - "Allow reports whether a call may proceed, after applying the time transition above:
      closed    -> returns true.
      open      -> returns false (cooldown not yet elapsed; had it elapsed, the transition above already moved it to half-open and Allow returns true).
      half-open -> returns true."
  - "Failure drives the state (after the time transition):
      closed    -> increment the consecutive-failure count; if it reaches maxFailures, transition to open and record the open time.
      half-open -> transition straight back to open and record the open time.
      Recording a failure always resets the half-open probe."
  - "Success drives the state (after the time transition):
      closed    -> reset the consecutive-failure count to 0.
      half-open -> transition to closed and reset the failure count to 0.
      So Failure then (cooldown elapses) then Success closes the breaker, because Success applies the time transition to half-open first, then closes."
  - "State returns the current stored state string. It does NOT itself advance time (callers observe a transition after the next Allow/Success/Failure)."
  - "CONCURRENCY: all methods may be called from many goroutines; guard all state with a sync.Mutex. Use time.Now / time.Since for the cooldown. Simplification (state explicitly): half-open admits any caller (no single-probe gating); the first Success closes it, the first Failure reopens it."
constraints: package main; standard library only
