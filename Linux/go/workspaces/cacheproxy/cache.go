package main

import (
	"sync"
	"time"
)

const shardCount = 32

// Cache is a concurrency-safe in-memory store of cached values, keyed by string
type Cache struct {
	shards [shardCount]struct {
		mu    sync.RWMutex
		store map[string]cacheEntry
	}
}

// cacheEntry holds a cached value and its expiration time
type cacheEntry struct {
	value     string
	expiresAt time.Time
}

// NewCache creates and returns a new Cache
func NewCache() *Cache {
	c := &Cache{}
	for i := range c.shards {
		c.shards[i].store = make(map[string]cacheEntry)
	}
	return c
}

// Get looks up a value by key, reporting whether it was present
func (c *Cache) Get(key string) (string, bool) {
	shardIndex := fnv1aHash(key) % shardCount
	shard := &c.shards[shardIndex]

	shard.mu.RLock()
	defer shard.mu.RUnlock()

	entry, exists := shard.store[key]
	if !exists {
		return "", false
	}

	// Check if the entry has expired
	if time.Now().After(entry.expiresAt) {
		delete(shard.store, key)
		return "", false
	}

	return entry.value, true
}

// Put stores a value under a key
func (c *Cache) Put(key string, value string, ttl time.Duration) {
	shardIndex := fnv1aHash(key) % shardCount
	shard := &c.shards[shardIndex]

	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.store[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

// fnv1aHash implements FNV-1a hash algorithm
func fnv1aHash(s string) uint32 {
	var hash uint32 = 2166136261 // FNV offset basis
	for i := 0; i < len(s); i++ {
		hash ^= uint32(s[i])
		hash *= 16777619 // FNV prime
	}
	return hash
}
