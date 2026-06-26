package main

import (
	"time"
)

// Event represents the unit of work carried from the HTTP handler to the async workers
type Event struct {
	ID       string
	Body     []byte
	Received time.Time
}

// NewEvent creates a new Event with the given ID and body, setting Received to time.Now()
func NewEvent(id string, body []byte) Event {
	return Event{
		ID:       id,
		Body:     body,
		Received: time.Now(),
	}
}
