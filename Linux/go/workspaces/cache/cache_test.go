package main

import (
	"sync"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	t.Run("Hit", func(t *testing.T) {
		c := NewCache()
		c.Set("k", "v", time.Minute)
		got, ok := c.Get("k")
		if got != "v" || !ok {
			t.Errorf("Get(k) = (%q, %v); want (\"v\", true)", got, ok)
		}
	})

	t.Run("Miss", func(t *testing.T) {
		c := NewCache()
		_, ok := c.Get("k")
		if ok {
			t.Errorf("Get(k) = (_, true); want (_, false)")
		}
	})

	t.Run("Expiry", func(t *testing.T) {
		c := NewCache()
		c.Set("k", "v", time.Millisecond)
		time.Sleep(5 * time.Millisecond)
		got, ok := c.Get("k")
		if got != "" || ok {
			t.Errorf("Get(k) = (%q, %v); want (\"\", false)", got, ok)
		}
	})

	t.Run("Edge", func(t *testing.T) {
		c := NewCache()
		c.Set("k", "v", 0)
		got, ok := c.Get("k")
		if got != "" || ok {
			t.Errorf("Get(k) = (%q, %v); want (\"\", false)", got, ok)
		}

		c.Set("k", "v", -time.Second)
		got, ok = c.Get("k")
		if got != "" || ok {
			t.Errorf("Get(k) = (%q, %v); want (\"\", false)", got, ok)
		}
	})

	t.Run("Replace", func(t *testing.T) {
		c := NewCache()
		c.Set("k", "a", time.Minute)
		c.Set("k", "b", time.Minute)
		got, ok := c.Get("k")
		if got != "b" || !ok {
			t.Errorf("Get(k) = (%q, %v); want (\"b\", true)", got, ok)
		}
	})

	t.Run("Concurrency", func(t *testing.T) {
		c := NewCache()
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(2)
			go func() {
				c.Set("k", "v", time.Minute)
				wg.Done()
			}()
			go func() {
				c.Get("k")
				wg.Done()
			}()
		}
		wg.Wait()
	})
}

func FuzzCache(f *testing.F) {
	f.Add("key", "value", int64(1000))
	f.Fuzz(func(t *testing.T, key, value string, ttlMs int64) {
		c := NewCache()
		ttl := time.Duration(ttlMs) * time.Millisecond
		c.Set(key, value, ttl)
		got, ok := c.Get(key)
		if ttl <= 0 {
			if got != "" || ok {
				t.Errorf("For ttl <= 0, Get should return (\"\", false)")
			}
		} else {
			if got != value || !ok {
				t.Errorf("Round-trip failed: got %q, ok %v; want %q, true", got, ok, value)
			}
		}
	})
}
