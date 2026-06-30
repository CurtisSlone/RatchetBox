package main

// file-kw: job data structure webhook delivery

// kw: job state constants
const (
	StatePending   = "pending"
	StateDelivered = "delivered"
	StateFailed    = "failed"
)

// kw: job data record
type Job struct {
	ID       string
	URL      string
	Payload  []byte
	State    string
	Attempts int
}

// kw: create new job
func NewJob(id, url string, payload []byte) *Job {
	return &Job{
		ID:       id,
		URL:      url,
		Payload:  payload,
		State:    StatePending,
		Attempts: 0,
	}
}
