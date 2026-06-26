package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhookAcceptsAndProcesses(t *testing.T) {
	dispatcher := NewDispatcher(2, 16)
	dispatcher.Start()
	server := NewServer(dispatcher)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/webhook", nil)
	request.Body = http.NoBody

	server.Webhook(recorder, request)
	if recorder.Code != 202 {
		t.Errorf("Expected status 202, got %d", recorder.Code)
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

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/webhook", nil)

	server.Webhook(recorder, request)
	if recorder.Code != 405 {
		t.Errorf("Expected status 405, got %d", recorder.Code)
	}

	dispatcher.Stop()
}
