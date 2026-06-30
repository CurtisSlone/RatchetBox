name: Queue
role: component
dependsOn: Job
intent: An in-memory, concurrency-safe FIFO queue of pending webhook jobs. The HTTP server pushes accepted jobs onto it; a worker pops them off to deliver. This is the simplest possible queue for the walking skeleton - a later layer makes durability the write-ahead log's job, not the queue's.
api:
  - type Queue struct { ... }
  - func NewQueue() *Queue
  - func (q *Queue) Push(j *Job)
  - func (q *Queue) Pop() (*Job, bool)
  - func (q *Queue) Len() int
behavior:
  - "Back the queue with a slice []*Job guarded by a sync.Mutex (every method locks). It is accessed concurrently by the HTTP handler and the worker, so it MUST be safe for concurrent use."
  - "Push appends the job to the tail."
  - "Pop removes and returns the head job and true; if the queue is empty it returns (nil, false). FIFO order."
  - "Len returns the current number of queued jobs (lock, read len, unlock)."
constraints: package main; standard library only (sync); uses Job
