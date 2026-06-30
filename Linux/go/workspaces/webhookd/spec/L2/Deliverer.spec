name: Deliverer
role: interface
dependsOn: Job
intent: The delivery abstraction - how a job's payload actually reaches its destination. An interface so the worker can be tested against a fake, and the real HTTP implementation can be swapped. The real one POSTs the payload over HTTP; 2xx means delivered.
api:
  - type Deliverer interface { Deliver(ctx context.Context, j *Job) error }
  - type HTTPDeliverer struct { ... }
  - func NewHTTPDeliverer(client *http.Client) *HTTPDeliverer
  - func (d *HTTPDeliverer) Deliver(ctx context.Context, j *Job) error
behavior:
  - "Deliverer is the interface: Deliver(ctx context.Context, j *Job) error. nil error = delivered, non-nil = a delivery failure (retryable by later layers)."
  - "HTTPDeliverer holds an *http.Client. NewHTTPDeliverer stores it; if the passed client is nil, default to &http.Client{Timeout: 5 * time.Second}."
  - "Deliver builds an http.NewRequestWithContext(ctx, http.MethodPost, j.URL, bytes.NewReader(j.Payload)) with Content-Type application/json, sends it via the client, and ALWAYS closes resp.Body (defer)."
  - "A 2xx status (resp.StatusCode >= 200 && < 300) returns nil. Any non-2xx returns fmt.Errorf with the status code. A transport error (client.Do returns err) is returned as-is."
constraints: package main; standard library only (net/http, bytes, context, fmt, time); uses Job
