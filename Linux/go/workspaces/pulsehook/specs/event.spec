name: Event
role: data
intent: the unit of work carried from the HTTP handler to the async workers
api:
  - type Event struct with fields: ID string; Body []byte; Received time.Time
  - func NewEvent(id string, body []byte) Event   // sets Received to time.Now()
constraints: standard library only; package main; no func main in this file
