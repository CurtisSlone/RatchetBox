package main

// file-kw: job unit work dispatcher delivers webhook destination url payload attempt count delivery state plain

// kw: job state unit work dispatcher
type JobState string

const (
	StatePending   JobState = "pending"
	StateDelivered JobState = "delivered"
	StateDead      JobState = "dead"
)

// kw: job unit work dispatcher
type Job struct {
	ID       string
	URL      string
	Payload  []byte
	Attempts int
	State    JobState
}

// kw: job url payload unit work dispatcher
func NewJob(id, url string, payload []byte) Job {
	return Job{
		ID:       id,
		URL:      url,
		Payload:  payload,
		Attempts: 0,
		State:    StatePending,
	}
}
