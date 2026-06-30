name: CounterTest
role: test
intent: Test the counter's arithmetic and, most importantly, its correctness under concurrency.
api:
  - func TestCounter(t *testing.T)
behavior:
  - "Basic: a new counter has Value 0; after Inc it is 1; after Add(5) it is 6; after Add(-2) it is 4."
  - "CONCURRENCY: launch many goroutines (e.g. 100) that each call Inc a fixed number of times (e.g. 1000), wait for all of them with a sync.WaitGroup, then assert Value equals goroutines*increments exactly. This test is meant to be run under -race."
  - "Include a fuzz target FuzzCounterAdd(f) that adds a sequence of int64 deltas (derive them from a []byte input) and asserts Value equals the plain sum of those deltas."
constraints: package: main
