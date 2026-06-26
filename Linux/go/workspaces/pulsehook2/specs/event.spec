name: Event
role: data
intent: the unit of work from handler to workers
api:
  - type Event struct { ID string; Body []byte; Received time.Time }
  - func NewEvent(id string, body []byte) Event   // Received = time.Now()
constraints: standard library only; package main; no func main here
