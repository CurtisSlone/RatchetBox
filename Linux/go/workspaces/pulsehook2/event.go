package main

import (
	"time"
)

type Event struct {
	ID       string
	Body     []byte
	Received time.Time
}

func NewEvent(id string, body []byte) Event {
	return Event{
		ID:       id,
		Body:     body,
		Received: time.Now(),
	}
}
