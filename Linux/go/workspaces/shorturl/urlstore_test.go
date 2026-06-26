package main

import (
	"sync"
	"testing"
)

func TestURLStore(t *testing.T) {
	store := NewURLStore()

	// Test storing and retrieving a URL
	code := "abc123"
	url := "https://example.com"

	if err := store.Put(code, url); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	retrievedURL, exists := store.Get(code)
	if !exists {
		t.Fatalf("Get failed: code %s not found", code)
	}

	if retrievedURL != url {
		t.Errorf("Expected %s, got %s", url, retrievedURL)
	}

	// Test Get returns false for non-existent codes
	_, exists = store.Get("nonexistent")
	if exists {
		t.Errorf("Get should return false for non-existent code")
	}

	// Test Put overwrites existing entries
	newURL := "https://newexample.com"
	if err := store.Put(code, newURL); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	retrievedURL, exists = store.Get(code)
	if !exists {
		t.Fatalf("Get failed: code %s not found", code)
	}

	if retrievedURL != newURL {
		t.Errorf("Expected %s, got %s", newURL, retrievedURL)
	}

	// Test concurrent access is safe
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			code := "concurrent" + string(rune(i))
			testURL := "https://test" + string(rune(i)) + ".com"
			store.Put(code, testURL)
			_, _ = store.Get(code)
		}(i)
	}
	wg.Wait()
}
