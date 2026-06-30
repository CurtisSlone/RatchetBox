package main

// file-kw: http deliverer interface http client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

// kw: deliver job http post
type Deliverer interface {
	Deliver(ctx context.Context, j *Job) error
}

// kw: http deliverer struct
type HTTPDeliverer struct {
	client *http.Client
}

// kw: create http deliverer
func NewHTTPDeliverer(client *http.Client) *HTTPDeliverer {
	if client == nil {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	return &HTTPDeliverer{client: client}
}

// kw: deliver job over http
func (d *HTTPDeliverer) Deliver(ctx context.Context, j *Job) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, j.URL, bytes.NewReader(j.Payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("delivery failed with status %d", resp.StatusCode)
}
