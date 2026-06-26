package main

import (
	"errors"
	"sync/atomic"
)

// URLShortener combines store and encoder to provide URL shortening functionality
type URLShortener struct {
	store   *URLStore
	encoder *Base62Encoder
	counter atomic.Int64
}

// NewURLShortener creates a new URLShortener instance
func NewURLShortener() *URLShortener {
	return &URLShortener{
		store:   NewURLStore(),
		encoder: NewBase62Encoder(),
	}
}

// Shorten returns a unique short code for the given URL
func (s *URLShortener) Shorten(url string) (string, error) {
	if url == "" {
		return "", errors.New("url cannot be empty")
	}

	id := s.counter.Add(1)
	code := s.encoder.Encode(id)
	err := s.store.Put(code, url)
	if err != nil {
		return "", err
	}
	return code, nil
}

// Expand returns the original URL for a valid short code
func (s *URLShortener) Expand(code string) (string, bool) {
	url, exists := s.store.Get(code)
	return url, exists
}
