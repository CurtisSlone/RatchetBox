package main

import (
	"sync"
	"time"
)

type Cache struct {
	mu    sync.RWMutex
	items map[string]item
}

type item struct {
	value  string
	expiry time.Time
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]item),
	}
}

func (c *Cache) Set(key string, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiry := time.Now().Add(ttl)
	// If ttl is zero or negative, expiry will be in the past
	// which means the entry is already expired
	c.items[key] = item{
		value:  value,
		expiry: expiry,
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return "", false
	}

	if time.Now().After(item.expiry) {
		// Entry expired, remove it
		// Note: This is a simplification. In a real cache,
		// we might want to remove expired entries lazily or
		// use a separate goroutine to clean up.
		// For this exercise, we'll remove it here to match
		// the behavior described in the test.
		delete(c.items, key)
		return "", false
	}

	return item.value, true
}
