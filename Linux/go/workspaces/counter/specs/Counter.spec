name: Counter
role: component
intent: A thread-safe integer counter that is correct under heavy concurrent use. Many goroutines may call Inc and Add at the same time; every increment must be counted exactly once with no lost updates or data races.
api:
  - func NewCounter() *Counter
  - func (c *Counter) Inc()
  - func (c *Counter) Add(n int64)
  - func (c *Counter) Value() int64
behavior:
  - "NewCounter returns a counter whose Value is 0."
  - "Inc increases the counter by exactly 1. Add(n) increases it by exactly n (n may be negative)."
  - "Value returns the current total."
  - "CONCURRENCY INVARIANT: if G goroutines each call Inc K times and then all finish, Value equals G*K exactly. The type must be safe for concurrent use by multiple goroutines (no data race under go test -race) - use a sync.Mutex or sync/atomic."
constraints: package: main
