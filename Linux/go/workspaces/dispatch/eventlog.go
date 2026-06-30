package main

// file-kw: eventlog event log write ahead sourcing recording delivery order state never mutated place ordered

import (
	"sync"
	"time"
)

// kw: event log write
type EventType string

const (
	EventEnqueued  EventType = "enqueued"
	EventDelivered EventType = "delivered"
	EventFailed    EventType = "failed"
)

// kw: event log write
type Event struct {
	JobID   string
	Type    EventType
	Attempt int
	At      time.Time
}

// kw: event log write
type EventLog struct {
	events []Event
	mutex  sync.Mutex
}

// kw: event log write
func NewEventLog() *EventLog {
	return &EventLog{
		events: make([]Event, 0),
	}
}

// kw: event log write
func (l *EventLog) Append(e Event) {
	l.mutex.Lock()
	l.events = append(l.events, e)
	l.mutex.Unlock()
}

// kw: events event log write
func (l *EventLog) Events() []Event {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Return a copy to prevent callers from mutating the backing storage
	eventsCopy := make([]Event, len(l.events))
	copy(eventsCopy, l.events)
	return eventsCopy
}

// kw: replay event log job state write
func (l *EventLog) Replay() map[string]JobState {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Fold the events into current job states
	jobStates := make(map[string]JobState)

	for _, event := range l.events {
		// Only the last event for a job determines its state
		switch event.Type {
		case EventEnqueued:
			jobStates[event.JobID] = StatePending
		case EventDelivered:
			jobStates[event.JobID] = StateDelivered
		case EventFailed:
			// If there's no delivered event after this failed event, it's still pending
			// But we don't track that here - we just keep the last event's state
			// If a job has a failed event and no delivered event, it's still pending
			// unless we have a convention that failed events mean dead
			// For now, we'll treat it as pending unless there's a delivered event
			if jobStates[event.JobID] != StateDelivered {
				jobStates[event.JobID] = StatePending
			}
		}
	}

	// Convert the map to the expected return type
	result := make(map[string]JobState)
	for jobID, state := range jobStates {
		result[jobID] = state
	}

	return result
}
