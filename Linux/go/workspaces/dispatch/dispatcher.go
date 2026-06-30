package main

// file-kw: dispatcher orchestrator delivers jobs external destination through pluggable deliverer recording outcome event log gating

import (
	"sync"
	"time"
)

// kw: deliverer orchestrator delivers jobs
type Deliverer interface {
	Deliver(url string, payload []byte) error
}

// kw: dispatcher orchestrator delivers jobs
type Dispatcher struct {
	log     *EventLog
	breaker *Breaker
	policy  RetryPolicy
	deliver Deliverer
}

// kw: dispatcher log event breaker policy retry deliverer orchestrator delivers
func NewDispatcher(log *EventLog, breaker *Breaker, policy RetryPolicy, deliverer Deliverer) *Dispatcher {
	return &Dispatcher{
		log:     log,
		breaker: breaker,
		policy:  policy,
		deliver: deliverer,
	}
}

// kw: dispatch dispatcher job orchestrator delivers jobs
func (d *Dispatcher) Dispatch(job Job) Job {
	attempt := job.Attempts + 1
	d.log.Append(Event{
		JobID:   job.ID,
		Type:    EventEnqueued,
		Attempt: attempt,
		At:      time.Now(),
	})

	if !d.breaker.Allow() {
		d.log.Append(Event{
			JobID:   job.ID,
			Type:    EventFailed,
			Attempt: attempt,
			At:      time.Now(),
		})
		return job
	}

	for {
		err := d.deliver.Deliver(job.URL, job.Payload)
		if err == nil {
			d.breaker.Success()
			d.log.Append(Event{
				JobID:   job.ID,
				Type:    EventDelivered,
				Attempt: attempt,
				At:      time.Now(),
			})
			job.State = StateDelivered
			job.Attempts = attempt
			return job
		}

		d.breaker.Failure()
		d.log.Append(Event{
			JobID:   job.ID,
			Type:    EventFailed,
			Attempt: attempt,
			At:      time.Now(),
		})

		if !d.policy.ShouldRetry(attempt) {
			job.State = StateDead
			job.Attempts = attempt
			return job
		}

		time.Sleep(d.policy.Backoff(attempt))
		attempt++
	}
}

// kw: dispatch dispatcher jobs job workers orchestrator delivers
func (d *Dispatcher) DispatchAll(jobs []Job, workers int) []Job {
	if workers < 1 {
		workers = 1
	}

	jobChan := make(chan Job, len(jobs))
	for _, job := range jobs {
		jobChan <- job
	}
	close(jobChan)

	var wg sync.WaitGroup
	resultChan := make(chan Job, len(jobs))

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobChan {
				result := d.Dispatch(job)
				resultChan <- result
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var results []Job
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}
