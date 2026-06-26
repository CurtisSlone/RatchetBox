# pulsehook - build transcript (specs -> generated, verified, running)

A low-latency webhook receiver built by the go ratchet from `.spec` files only. The five spec files were the ONLY hand-written input; the local model (qwen3-coder via Ollama) generated every line of Go through the `compose` and `add_file` flows, and the deterministic oracles verified each step.

- Generated: 2026-06-26
- Build: `ratchet flow . compose --ws pulsehook ""`  (console: `/flow compose --ws pulsehook`)
- Test added later: `ratchet flow . add_file --ws pulsehook "server_test.go ..."`
- Oracles: per-unit `stage_build` (go vet over the whole module), final `module_check` (go build ./... + go test ./...).

## Input specs (the only hand-written artifacts)

### `app.spec`

```yaml
name: App
role: behavior
intent: the program entry point - wire the dispatcher and server together and serve
behavior:
  - in func main: create a Dispatcher with NewDispatcher(4, 1024) and call Start()
  - create a Server with NewServer(dispatcher)
  - log that it is listening, then http.ListenAndServe(":8080", server.Routes())
  - if ListenAndServe returns an error, log.Fatal it
constraints: standard library only (net/http, log); package main; this file (main.go) is the ONLY file
  with func main; uses the existing Dispatcher and Server API verbatim
```

### `dispatcher.spec`

```yaml
name: Dispatcher
role: component
intent: a worker pool that processes Events asynchronously off a buffered channel, so the HTTP handler
  can accept a webhook and return immediately (low latency) instead of doing the work inline
api:
  - type Dispatcher struct holding a buffered chan Event (the queue), a sync.WaitGroup, a worker count,
    and an int64 processed counter updated with sync/atomic
  - func NewDispatcher(workers, buffer int) *Dispatcher   // makes the channel with capacity buffer
  - method (*Dispatcher) Start()   // launches `workers` goroutines, each draining the queue with
    `for e := range d.queue { ... atomic.AddInt64(&d.processed, 1) }`
  - method (*Dispatcher) Enqueue(e Event) bool   // NON-BLOCKING: select { case d.queue <- e: return true; default: return false }
    so a full queue never blocks the caller (this is the low-latency guarantee)
  - method (*Dispatcher) Stop()   // close(d.queue) then d.wg.Wait() to drain workers cleanly
  - method (*Dispatcher) Processed() int64   // atomic.LoadInt64 of the processed counter
behavior:
  - each worker goroutine must call wg.Done() when the range loop ends (defer wg.Done() at the top)
  - Start must call wg.Add(1) per worker before launching it
constraints: standard library only (sync, sync/atomic); package main; no func main in this file;
  pointer receivers (it holds a sync.WaitGroup and must not be copied)
```

### `event.spec`

```yaml
name: Event
role: data
intent: the unit of work carried from the HTTP handler to the async workers
api:
  - type Event struct with fields: ID string; Body []byte; Received time.Time
  - func NewEvent(id string, body []byte) Event   // sets Received to time.Now()
constraints: standard library only; package main; no func main in this file
```

### `server.spec`

```yaml
name: Server
role: component
intent: the HTTP layer - accept a webhook POST, hand it to the Dispatcher, and reply immediately
api:
  - type Server struct holding a *Dispatcher (field name: disp)
  - func NewServer(d *Dispatcher) *Server
  - method (*Server) Webhook(w http.ResponseWriter, r *http.Request)   // the handler:
      * if r.Method is not POST, http.Error with status 405 and return
      * read the body with io.ReadAll(r.Body); on error, 400 and return
      * build an Event with NewEvent (any short id is fine, e.g. the time or a counter)
      * if s.disp.Enqueue(event) is false (queue full), reply 503 Service Unavailable and return
      * otherwise reply 202 Accepted immediately (do NOT process inline) and write a tiny body
  - method (*Server) Routes() http.Handler   // a *http.ServeMux with HandleFunc("/webhook", s.Webhook)
behavior:
  - the handler must return right after enqueuing - the actual work happens in the Dispatcher workers,
    which is what keeps request latency low
constraints: standard library only (net/http, io); package main; no func main in this file; uses the
  existing Dispatcher and Event API verbatim
```

### `server_test.spec`

```yaml
name: ServerTest
role: test
intent: prove the webhook returns immediately with 202 and the event is processed asynchronously
behavior:
  - a Go test file (package main) with func TestWebhookAcceptsAndProcesses(t *testing.T)
  - build a Dispatcher with NewDispatcher(2, 16) and call Start()
  - build a Server with NewServer(dispatcher)
  - use net/http/httptest: make an httptest.NewRecorder() and an httptest.NewRequest("POST", "/webhook",
    strings.NewReader("hello")) and call server.Webhook(rec, req) directly
  - assert rec.Code == http.StatusAccepted (202)
  - call dispatcher.Stop() (this drains the workers), then assert dispatcher.Processed() == 1
  - also add func TestWebhookRejectsGet(t *testing.T): a GET request must yield 405
constraints: standard library only (net/http, net/http/httptest, strings, testing); package main; uses
  the existing Dispatcher and Server API verbatim
```

---

## Compose: plan -> worklist -> per-unit generation -> whole-module verify

### Build plan (model-inferred from the specs)

```json
{"contracts":[{"name":"Event","normalizedFrom":"","type":"struct"},{"name":"Dispatcher","normalizedFrom":"","type":"struct"},{"name":"Server","normalizedFrom":"","type":"struct"}],"units":[{"dependsOn":[],"name":"Event","role":"data","spec":"event.spec"},{"dependsOn":["Event"],"name":"Dispatcher","role":"component","spec":"dispatcher.spec"},{"dependsOn":["Dispatcher","Event"],"name":"Server","role":"component","spec":"server.spec"},{"dependsOn":["Dispatcher","Server"],"name":"App","role":"behavior","spec":"app.spec"}]}
```

### Worklist (dependency order)

```text
event.go event.spec
dispatcher.go dispatcher.spec
server.go server.spec
main.go app.spec
```

> NOTE: this first compose plan DROPPED the test spec - the plan schema's `role` enum had no `test`, so the model could not classify `server_test.spec` and omitted it (`module_check` then reported `[no test files]`). Fixed after: added `test` to the role enum + a plan-prompt rule that test specs become `_test.go` units ordered last. The test was added to this app via `add_file` (below), which also exercises the lifecycle flow on a real multi-file project.

### Per-unit generation (`add_unit` sub-runs)

#### `event.go event.spec`  (first try)

Generated:

```go
package main

import (
	"time"
)

// Event represents the unit of work carried from the HTTP handler to the async workers
type Event struct {
	ID        string
	Body      []byte
	Received  time.Time
}

// NewEvent creates a new Event with the given ID and body, setting Received to time.Now()
func NewEvent(id string, body []byte) Event {
	return Event{
		ID:       id,
		Body:     body,
		Received: time.Now(),
	}
}
```
```

stage_build (whole-module vet): pass

```text
OK: staged event.go; whole module type-checks (go vet clean)
```

#### `dispatcher.go dispatcher.spec`  (one repair)

Generated:

```go
package main

import (
	"sync"
	"sync/atomic"
)

type Dispatcher struct {
	queue    chan Event
	wg       sync.WaitGroup
	workers  int
	processed int64
}

func NewDispatcher(workers, buffer int) *Dispatcher {
	return &Dispatcher{
		queue:   make(chan Event, buffer),
		workers: workers,
	}
}

func (d *Dispatcher) Start() {
	for i := 0; i < d.workers; i++ {
		d.wg.Add(1)
		go func() {
			defer d.wg.Done()
			for e := range d.queue {
				atomic.AddInt64(&d.processed, 1)
			}
		}()
	}
}

func (d *Dispatcher) Enqueue(e Event) bool {
	select {
	case d.queue <- e:
		return true
	default:
		return false
	}
}

func (d *Dispatcher) Stop() {
	close(d.queue)
	d.wg.Wait()
}

func (d *Dispatcher) Processed() int64 {
	return atomic.LoadInt64(&d.processed)
}
```
```

stage_build (whole-module vet): FAIL

```text
BUILD FAILED after staging dispatcher.go:
# pulsehook
# [pulsehook]
vet: ./dispatcher.go:27:8: declared and not used: e
[exit code 1]
```

#### `server.go server.spec`  (first try)

Generated:

```go
package main

import (
	"io"
	"net/http"
	"time"
)

type Server struct {
	disp *Dispatcher
}

func NewServer(d *Dispatcher) *Server {
	return &Server{disp: d}
}

func (s *Server) Webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	event := NewEvent(time.Now().String(), body)

	if !s.disp.Enqueue(event) {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("accepted"))
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", s.Webhook)
	return mux
}
```
```

stage_build (whole-module vet): pass

```text
OK: staged server.go; whole module type-checks (go vet clean)
```

#### `main.go app.spec`  (one repair)

Generated:

```go
package main

import (
	"http"
	"log"
)

func main() {
	dispatcher := NewDispatcher(4, 1024)
	dispatcher.Start()
	server := NewServer(dispatcher)
	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", server.Routes()); err != nil {
		log.Fatal(err)
	}
}
```

stage_build (whole-module vet): FAIL

```text
BUILD FAILED after staging main.go:
main.go:4:2: package http is not in std (/usr/local/go/src/http)
[exit code 1]
```

### Final whole-module oracle (module_check)

```text
OK: module builds and tests pass with go1.26.4
?   	pulsehook	[no test files]
```

---

## Adding the behavior test (`add_file`, first try)

Generated `server_test.go`:

```go
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
```

stage_check (go vet + go test): pass

```text
OK: staged server_test.go; module vets and tests pass
ok  	pulsehook	0.002s
```

---

## Verification (run by hand - the run flow is not built yet)

Whole-module tests, verbose:

```text
$ go test -v ./...
=== RUN   TestWebhookAcceptsAndProcesses
--- PASS: TestWebhookAcceptsAndProcesses (0.00s)
=== RUN   TestWebhookRejectsGet
--- PASS: TestWebhookRejectsGet (0.00s)
PASS
ok  	pulsehook	0.002s
```

Live server (built binary, real curl):

```text
$ go build -o /tmp/pulsehook . && /tmp/pulsehook &      # listening on :8080
$ curl -X POST localhost:8080/webhook -d event-1   ->  HTTP 202  in 0.000878s
$ curl -X POST localhost:8080/webhook -d event-2   ->  HTTP 202  in 0.000437s
$ curl -X POST localhost:8080/webhook -d event-3   ->  HTTP 202  in 0.000405s
$ curl     localhost:8080/webhook                  ->  HTTP 405
time-to-first-byte on a POST: 0.000355s
```

The handler enqueues to a buffered channel and returns 202 immediately (non-blocking select), while a 4-worker pool drains the queue asynchronously - sub-millisecond request latency regardless of processing cost. That is the design the specs asked for and the model produced.
