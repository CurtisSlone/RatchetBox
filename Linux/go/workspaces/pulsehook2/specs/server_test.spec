name: ServerTest
role: test
intent: prove the webhook accepts (202) and processes asynchronously
behavior:
  - package main test with func TestWebhookAcceptsAndProcesses(t *testing.T): NewDispatcher(2,16);
    Start(); NewServer; httptest.NewRecorder + httptest.NewRequest POST /webhook body "hello";
    call server.Webhook; assert 202; dispatcher.Stop(); assert Processed()==1
  - func TestWebhookRejectsGet: GET -> 405
constraints: standard library (net/http, net/http/httptest, strings, testing); package main
