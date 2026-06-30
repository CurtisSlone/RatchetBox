name: EventLog
role: component
intent: An append-only event log (write-ahead log / event sourcing) recording every delivery event in order. State is never mutated in place - it is the ordered sequence of events, and current job states are derived by folding (replaying) the log. Safe for concurrent use.
api:
  - type EventType string
  - const EventEnqueued EventType = "enqueued"
  - const EventDelivered EventType = "delivered"
  - const EventFailed EventType = "failed"
  - type Event struct { JobID string; Type EventType; Attempt int; At time.Time }
  - type EventLog struct { ... }
  - func NewEventLog() *EventLog
  - func (l *EventLog) Append(e Event)
  - func (l *EventLog) Events() []Event
  - func (l *EventLog) Replay() map[string]JobState
behavior:
  - "APPEND-ONLY: Append adds an event to the end of the log. There is NO update or delete. The log only grows; existing events are never changed."
  - "Events returns all appended events in append order (a copy, so callers cannot mutate the log's backing storage)."
  - "Replay FOLDS the log into the current state per job id: an enqueued event -> pending; a failed event -> pending (still retryable unless a later event supersedes); a delivered event -> delivered. The LAST event for a job id determines its state; if the last event is failed AND a separate convention marks exhaustion, treat a job with no delivered event after its final failed event as pending. Keep the fold simple: delivered if any delivered event exists for the id, else pending."
  - "CONCURRENCY: Append, Events, and Replay may be called from many goroutines at once; guard the backing slice with a sync.Mutex so there is no data race under go test -race."
  - "An empty log: Events returns an empty slice and Replay returns an empty map."
constraints: package main; standard library only; uses Job (JobState)
