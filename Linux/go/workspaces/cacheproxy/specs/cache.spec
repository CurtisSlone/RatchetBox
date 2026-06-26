name: Cache
role: data
intent: a concurrency-safe in-memory store of cached values, keyed by string
behavior:
  - look up a value by key, reporting whether it was present
  - store a value under a key
  - many goroutines read at once while writes are serialized: guard the map with a sync.RWMutex (a read
    lock for lookups, a write lock for stores) - a plain map would panic under concurrent use
constraints: standard library only; pointer receivers (it holds a mutex and must not be copied);
  initialize the map before first use; package main; no func main here
