package main

import (
	"sync"
)

// URLStore stores mappings from short codes to long URLs
type URLStore struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewURLStore creates and returns a new URLStore
func NewURLStore() *URLStore {
	return &URLStore{
		data: make(map[string]string),
	}
}

// Put stores a mapping from code to url
// If code already exists, it overwrites the previous value
func (s *URLStore) Put(code string, url string) error {
	s.mu.Lock()
	s.data[code] = url
	s.mu.Unlock()
	return nil
}

// Get retrieves the url for a given code
// Returns false if code does not exist
func (s *URLStore) Get(code string) (string, bool) {
	s.mu.RLock()
	url, exists := s.data[code]
	s.mu.RUnlock()
	return url, exists
}
