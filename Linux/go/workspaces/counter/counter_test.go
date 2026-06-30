package main

import (
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		c := NewCounter()
		if got := c.Value(); got != 0 {
			t.Errorf("New counter Value() = %d; want 0", got)
		}
		c.Inc()
		if got := c.Value(); got != 1 {
			t.Errorf("After Inc Value() = %d; want 1", got)
		}
		c.Add(5)
		if got := c.Value(); got != 6 {
			t.Errorf("After Add(5) Value() = %d; want 6", got)
		}
		c.Add(-2)
		if got := c.Value(); got != 4 {
			t.Errorf("After Add(-2) Value() = %d; want 4", got)
		}
	})

	t.Run("Concurrency", func(t *testing.T) {
		const goroutines = 100
		const increments = 1000
		c := NewCounter()
		var wg sync.WaitGroup
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < increments; j++ {
					c.Inc()
				}
			}()
		}
		wg.Wait()
		if got := c.Value(); got != goroutines*increments {
			t.Errorf("After %d goroutines each calling Inc %d times, Value() = %d; want %d", goroutines, increments, got, goroutines*increments)
		}
	})
}

func FuzzCounterAdd(f *testing.F) {
	f.Add([]byte{0})
	f.Add([]byte{1})
	f.Add([]byte{255}) // -1 in signed byte
	f.Add([]byte{100})
	f.Add([]byte{200}) // -56 in signed byte

	f.Fuzz(func(t *testing.T, data []byte) {
		c := NewCounter()
		expectedSum := int64(0)
		for _, b := range data {
			delta := int64(int8(b))
			c.Add(delta)
			expectedSum += delta
		}
		if got := c.Value(); got != expectedSum {
			t.Errorf("After adding deltas from data, Value() = %d; want %d", got, expectedSum)
		}
	})
}
