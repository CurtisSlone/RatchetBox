# cacheproxy

A composed Go module. Every unit is a file in `package main` at the module root.

## Units
- `fetcher.go` (composed unit)
- `cache.go` (composed unit)
- `proxy.go` (composed unit)
- `main.go` (composed unit)
- `proxy_test.go` (composed unit)
- edited `proxy_test.go`: the tests compare each result to the literal string "value", but MockFetcher returns a different value such as value-for-key-call-1, so every assertion fails. Fix the tests to NOT depend on the mock's exact string: for the herd and distinct-keys tests, assert that all goroutines requesting the same key received the SAME value, equal to one reference Get of that key; for cache-hit, assert the two requests return equal values. KEEP all the call-count checks unchanged (herd: Calls()==1; distinct-keys: Calls()==number of unique keys; cache-hit: Calls()==1 after the second request). Do not weaken the concurrency or call-count assertions, and do not add sleeps.
- `server.go` (added file)
- edited `main.go`: replace the concurrent demo in func main with a real server: build a MockFetcher, a Cache, a Proxy, and an HTTPServer (NewHTTPServer); construct an *http.Server with Addr :8080, Handler set to the HTTPServer's Routes(), and ReadHeaderTimeout, ReadTimeout, WriteTimeout, IdleTimeout all set; run srv.ListenAndServe in a goroutine (ignore http.ErrServerClosed); use signal.NotifyContext on SIGINT/SIGTERM and on ctx.Done call srv.Shutdown with a timeout context. Keep func main only in this file.
- edited `main.go`: staticcheck SA1016 flags that os.Kill cannot be trapped by a signal handler. In the signal.NotifyContext call, replace os.Kill with syscall.SIGTERM (keep os.Interrupt for SIGINT). Ensure syscall is imported and remove any now-unused import. Change nothing else.
- edited `cache.go`: refactor the internals to reduce lock contention by sharding: instead of one sync.RWMutex and one map, use a fixed array of 32 shards, each shard its own sync.RWMutex plus its own map. Route a key to a shard with an FNV-1a hash of the key modulo 32. Get takes a read lock on only that key's shard; Set takes a write lock on only that key's shard. CRITICAL: keep the EXACT same exported method signatures - NewCache() *Cache, (*Cache) Get(key) ([]byte, bool), (*Cache) Set(key, value) - so no other file changes. Initialize every shard's map in NewCache.
- `cache_bench_test.go` (added file)
