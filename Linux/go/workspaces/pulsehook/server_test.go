package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebhookAcceptsAndProcesses(t *testing.T) {
	dispatcher := NewDispatcher(2, 16)
	dispatcher.Start()
	server := NewServer(dispatcher)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/webhook", strings.NewReader("hello"))

	server.Webhook(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("Expected status code %d, got %d", http.StatusAccepted, rec.Code)
	}

	dispatcher.Stop()
	if dispatcher.Processed() != 1 {
		t.Errorf("Expected processed count 1, got %d", dispatcher.Processed())
	}
}

func TestWebhookRejectsGet(t *testing.T) {
	dispatcher := NewDispatcher(2, 16)
	dispatcher.Start()
	server := NewServer(dispatcher)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/webhook", nil)

	server.Webhook(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}

	dispatcher.Stop()
}
