package main

// file-kw: main runnable entry point wires whole dispatcher together demonstrates end including deliberately flaky destination

import (
	"fmt"
	"sync"
	"time"
)

// DemoDeliverer implements the Deliverer interface with a flaky delivery
// that fails the first two attempts and then succeeds.
// kw: demo deliverer runnable entry point
type DemoDeliverer struct {
	counter int
	mutex   sync.Mutex
}

// kw: deliver demo deliverer url payload runnable entry point
func (d *DemoDeliverer) Deliver(url string, payload []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.counter++
	if d.counter <= 2 {
		return fmt.Errorf("delivery failed on attempt %d", d.counter)
	}
	return nil
}

// kw: main runnable entry point
func main() {
	log := NewEventLog()
	breaker := NewBreaker(3, 50*time.Millisecond)
	policy := NewRetryPolicy(1*time.Millisecond, 2.0, 5)
	deliverer := &DemoDeliverer{}
	d := NewDispatcher(log, breaker, policy, deliverer)

	jobs := []Job{
		NewJob("job1", "http://example.com/1", []byte("payload1")),
		NewJob("job2", "http://example.com/2", []byte("payload2")),
		NewJob("job3", "http://example.com/3", []byte("payload3")),
	}

	resultJobs := d.DispatchAll(jobs, 4)

	for _, job := range resultJobs {
		fmt.Printf("Job ID: %s, State: %s\n", job.ID, job.State)
	}

	fmt.Printf("Number of events in log: %d\n", len(log.Events()))
}
