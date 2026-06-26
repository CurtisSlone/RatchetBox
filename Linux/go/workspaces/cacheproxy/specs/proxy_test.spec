name: ProxyTest
role: test
intent: prove, under the race detector, that the proxy is concurrency-safe and collapses the herd
behavior:
  - a Go test file in package main
  - herd test: a couple hundred goroutines request the SAME key at once through one Proxy backed by a
    MockFetcher with latency; afterward every result is the expected value and the origin was called
    exactly once
  - cache-hit test: requesting a warm key a second time does not call the origin again
  - distinct-keys test: concurrent requests across a handful of different keys call the origin once per
    unique key
  - lean on `go test -race` to prove there is no data race; do not hide races behind added sleeps
constraints: standard library (sync, testing, bytes); uses the existing Proxy, Cache, and MockFetcher
  exactly as they are defined; package main
