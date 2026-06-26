package main

import (
	"testing"
	"time"
)

func BenchmarkCacheSet(b *testing.B) {
	cache := NewCache()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := "key-" + string(rune(time.Now().UnixNano()%1000000))
			cache.Put(key, "value", time.Hour)
		}
	})
}

func BenchmarkCacheGet(b *testing.B) {
	cache := NewCache()
	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		key := "key-" + string(rune(i))
		cache.Put(key, "value", time.Hour)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := "key-" + string(rune(time.Now().UnixNano()%1000))
			cache.Get(key)
		}
	})
}

func BenchmarkProxyGet(b *testing.B) {
	fetcher := NewMockFetcher(0)
	cache := NewCache()
	proxy := NewProxy(cache, fetcher)
	// Pre-populate cache with some keys
	for i := 0; i < 100; i++ {
		key := "key-" + string(rune(i))
		value, _ := fetcher.Fetch(key)
		cache.Put(key, value, time.Hour)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := "key-" + string(rune(time.Now().UnixNano()%100))
			proxy.Get(key)
		}
	})
}
