package main

// file-kw: deadletter dead letter queue holding pen jobs could delivered after exhausting retries because breaker

import (
	"sync"
)

// DeadLetter is a concurrency-safe dead-letter queue for jobs that could not be delivered.
// kw: dead letter queue
type DeadLetter struct {
	jobs []*Job
	mu   sync.Mutex
}

// NewDeadLetter creates a new dead-letter queue.
// kw: dead letter queue
func NewDeadLetter() *DeadLetter {
	return &DeadLetter{}
}

// Add appends a job to the dead-letter queue.
// kw: add dead letter job queue
func (d *DeadLetter) Add(j *Job) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.jobs = append(d.jobs, j)
}

// Jobs returns a copy of the internal slice of jobs.
// kw: jobs dead letter job queue
func (d *DeadLetter) Jobs() []*Job {
	d.mu.Lock()
	defer d.mu.Unlock()
	cpy := make([]*Job, len(d.jobs))
	copy(cpy, d.jobs)
	return cpy
}

// Len returns the number of jobs in the dead-letter queue.
// kw: dead letter queue
func (d *DeadLetter) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.jobs)
}
