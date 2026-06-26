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
