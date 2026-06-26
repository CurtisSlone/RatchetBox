name: Proxy
role: component
intent: the caching proxy - serve from cache, and on a miss fetch once from the origin and cache it
behavior:
  - a Proxy is built from a Cache and a Fetcher
  - getting a key returns the cached value when present; on a miss it fetches from the origin, stores
    the result in the cache, and returns it
  - collapse the thundering herd: when many goroutines miss the SAME key at the same moment, only one
    fetch should run and all the callers share its result - use golang.org/x/sync/singleflight for this
constraints: uses the existing Cache and Fetcher; golang.org/x/sync/singleflight for stampede
  protection; pointer receivers; package main; no func main here
