name: App
role: behavior
intent: a small runnable demo that shows the herd collapsing to a single origin call
behavior:
  - in main, build a MockFetcher (with a little latency), a Cache, and a Proxy wired together
  - launch many goroutines (around 100) that all request the SAME key concurrently, and wait for them
    all to finish (a sync.WaitGroup)
  - print the value once and the number of times the origin was actually called - it should be 1
constraints: standard library (sync, fmt); uses the existing Proxy, Cache, and MockFetcher; this file
  (main.go) is the only one with func main
