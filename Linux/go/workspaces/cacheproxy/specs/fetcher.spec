name: Fetcher
role: interface
intent: the upstream origin the proxy calls on a cache miss; an interface so it can be swapped or mocked
behavior:
  - a Fetcher fetches the value for a key, returning the value and an error
  - a MockFetcher implements it for the demo and tests: it returns a deterministic value derived from
    the key, counts how many times it has actually been called (the count must be readable safely while
    other goroutines are fetching), and can be given a small artificial latency to widen the race window
constraints: standard library only; the call counter must be concurrency-safe; package main; no func main here
