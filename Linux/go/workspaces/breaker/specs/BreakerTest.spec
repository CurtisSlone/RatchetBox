name: BreakerTest
role: test
intent: Prove the circuit breaker's state machine: the exact transitions, the cooldown timing, safety under any operation sequence, and race-freedom under concurrent use.
api:
  - func TestBreaker(t *testing.T)
behavior:
  - "Closed start: NewBreaker(3, time.Minute).State() == \"closed\"."
  - "Trips open: on a fresh breaker, call Failure() three times (maxFailures); State() is then \"open\", and Allow() returns false (cooldown has not elapsed)."
  - "Cooldown -> half-open: NewBreaker(1, time.Millisecond); call Failure() (now open); time.Sleep(5*time.Millisecond); Allow() returns true and State() is then \"half-open\"."
  - "Half-open success closes: from the half-open state above, call Success(); State() == \"closed\"."
  - "Half-open failure reopens: reach half-open again, call Failure(); State() == \"open\"."
  - "INVARIANT (fuzz): FuzzBreaker(f) drives a random sequence of operations derived from a []byte input - for each byte, 0 => Allow(), 1 => Success(), 2 => Failure() (use b%3). After the whole sequence, assert State() is always one of exactly {\"closed\",\"open\",\"half-open\"} and that no call panicked. Seed with f.Add([]byte{2,2,2,0,1}). Use a small cooldown so time transitions are reachable."
  - "CONCURRENCY: a test that launches many goroutines each calling Allow/Success/Failure on one shared breaker, joined with a sync.WaitGroup, asserting it does not panic and State() is still a valid value afterward. Meant to run under -race."
constraints: package main; standard library only
