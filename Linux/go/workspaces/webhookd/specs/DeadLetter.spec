name: DeadLetter
role: component
dependsOn: Job
intent: A dead-letter queue - the holding pen for jobs that could not be delivered after exhausting retries (or because the breaker was open). Instead of silently dropping a permanently-failed webhook, the worker parks it here so it can be inspected or replayed. Concurrency-safe.
api:
  - type DeadLetter struct { ... }
  - func NewDeadLetter() *DeadLetter
  - func (d *DeadLetter) Add(j *Job)
  - func (d *DeadLetter) Jobs() []*Job
  - func (d *DeadLetter) Len() int
behavior:
  - "Back it with a slice []*Job guarded by a sync.Mutex (the worker adds from its goroutine while main reads). Every method locks."
  - "Add appends the job to the slice."
  - "Jobs returns a COPY of the slice (allocate a new []*Job, copy, return) so callers cannot mutate the internal slice under the lock."
  - "Len returns the current count."
constraints: package main; standard library only (sync); uses Job
