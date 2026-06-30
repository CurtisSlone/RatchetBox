package main

// file-kw: queue job fifo concurrency
import "sync"

// kw: append job queue concurrent
type Queue struct {
	jobs []*Job
	mu   sync.Mutex
}

// kw: create new queue
func NewQueue() *Queue {
	return &Queue{
		jobs: make([]*Job, 0),
	}
}

// kw: push job to queue
func (q *Queue) Push(j *Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.jobs = append(q.jobs, j)
}

// kw: pop job from queue
func (q *Queue) Pop() (*Job, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.jobs) == 0 {
		return nil, false
	}
	j := q.jobs[0]
	q.jobs = q.jobs[1:]
	return j, true
}

// kw: get queue length
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.jobs)
}
