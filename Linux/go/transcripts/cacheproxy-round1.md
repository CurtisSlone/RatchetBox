# cacheproxy - Round 1 build transcript

Building the Round 1 concurrent memory-cached proxy from `workspaces/cacheproxy/design/Round1.md`,
one ratchet command at a time. Raw input/output of each step.

## Step 1 - add the dependency

```text
$ /do add_dep cacheproxy golang.org/x/sync/singleflight
go: upgraded go 1.21 => 1.25.0
go: added golang.org/x/sync v0.21.0
ingested golang.org/x/sync/singleflight -> kb/deps/golang_org_x_sync_singleflight.md
reindexed kb/deps
OK: golang.org/x/sync/singleflight added to workspace cacheproxy and grounded in the deps KB
[exit 0]
```

## Step 2 - compose (FAILED at the final gate)

```text
$ ratchet flow . compose --ws cacheproxy ""
  compose.read -> plan -> worklist -> fan (foreach 5 units):
    fetcher.go     (1 repair, then ok)
    cache.go       (first try)
    proxy.go       (first try)
    main.go        (first try)
    proxy_test.go  (first try)
  compose.build (module_check: go build ./... + go test -race ./...) -> FAILED
  compose.fail

MODULE TEST FAILED (-race):
--- FAIL: TestProxy/herd (0.01s)
    proxy_test.go:37: goroutine 0 got value "value-for-key-call-1", want "value"   (x200, all same)
--- FAIL: TestProxy/cache-hit (0.00s)
    proxy_test.go:59: first request got "value-for-key-call-1", want "value"
--- FAIL: TestProxy/distinct-keys (0.01s)
    proxy_test.go:107: goroutine N got "value-for-key0-call-1" / "value-for-key1-call-5" / ... want "value"
FAIL	cacheproxy	0.028s
[exit code 1]
```

Diagnosis (input level, no code reviewed): NOT a data race and NOT a build error - a value-format
mismatch between two units. The `MockFetcher` returns values like `value-for-<key>-call-<n>`, but the
test asserts `want "value"`. The fetcher.spec and proxy_test.spec under-specified the mock's return
value, so the units disagreed. The race detector found no race; the test oracle caught the mismatch.

## Step 2b - edit_file (targeted fix, instead of re-composing)

The failure was in one unit's assertions, so we fed the error + fix into `edit_file` rather than
re-running compose on all 5 units.

```text
$ ratchet flow . edit_file --ws cacheproxy "proxy_test.go <assert the property (same value per key +
    call counts), not the mock's literal string>"
  edit_file.read -> readproj -> generate -> check (go vet + go test -race) -> log -> done
  logged edit to proxy_test.go
  [exit 0]   # first try, no repair: the whole module passes `go test -race`
```

Result: the module now passes `go vet` + `go test -race ./...`. The targeted edit fixed the one
offending file; the other four units were untouched.

## Step 3 - harden (the full gate)

```text
$ ratchet flow . harden --ws cacheproxy ""
  harden.check -> harden.done
  == go_quality: workspaces/cacheproxy ==
  ok   (gofmt)
  ok   (go vet)
  ok   (go build)
  ok   (go test -race)
  ok   (staticcheck)
  ok   (govulncheck (known CVEs))     # over golang.org/x/sync/singleflight - no known CVEs
  --   gosec absent (skipped)
  PRODUCTION-CLEAN: all available gates passed.
  [exit 0]
```

## Step 4 - run (the demo)

```text
$ ratchet flow . run --ws cacheproxy ""
  run.exec -> run.done
  built workspaces/cacheproxy; running for up to 5s...
  --- program output ---
  Value: value-for-key-call-1
  Origin calls: 1            # 100 concurrent same-key requests collapsed to ONE origin call
  ----------------------
  exited with status 0.
```

## Round 1 - Definition of Done: met

| DoD item | result |
|---|---|
| composes + `go build ./...` + `go test -race ./...` pass | yes (step 2 + targeted fix 2b) |
| `harden` PRODUCTION-CLEAN (incl. govulncheck on singleflight) | yes (step 3) |
| no data race under `-race`; single-flight collapse (1 call per key) | yes (race detector clean; distinct-keys one call/key) |
| `run` prints calls = 1 for the concurrent same-key burst | yes (step 4: Origin calls: 1) |

Round 1 is complete: the concurrent cache-proxy core, behavior-verified under the race detector and the
full quality/security gate. Deferred to later rounds (via edit_file): HTTP server, TTL, eviction/bound,
sharding, benchmarks, context.

---

# cacheproxy - Round 2 (edits on the live codebase)

## Step R2.1a - add_file server.go

Attempt 1 FAILED - the model wrote a `func main` into server.go (the word "server" cues an entry point),
colliding with main.go:

```text
$ ratchet flow . add_file --ws cacheproxy "server.go an HTTP front-end ..."
  ... generate -> check FAIL -> fix -> recheck FAIL -> add_file.fail
  VET FAILED after staging server.go:
  vet: ./server.go:39:6: main redeclared in this block
  [exit code 1]
```

Retry with the constraint fed into the request ("a LIBRARY file; do NOT declare func main") - first try:

```text
$ ratchet flow . add_file --ws cacheproxy "server.go a LIBRARY file (NOT an entry; no func main) ..."
  readproj -> api -> generate -> check (go vet + go test -race) -> register -> done
  registered server.go
  [exit 0]
```

### Step R2.1a (aside) - self-improvement on the live failure

The `func main redeclared` failure above became a KB entry, via the ratchet's own loop:

```text
$ ratchet flow . learn "redeclared top-level name ... func main redeclared ..."
  learn.generate -> learn.write -> learn.done
  wrote kb/pitfalls/redeclared-top-level-name.md ; reindexed (pitfalls: 8)
```

Also wired `pitfalls` into `add_file` and `edit_file` generate (they had only stdlib/patterns/deps), so
the entry now grounds those flows - the failure is unlikely to recur.

## Step R2.1b - edit_file main.go (demo -> production http.Server)

```text
$ ratchet flow . edit_file --ws cacheproxy "main.go ... production *http.Server + graceful shutdown"
  read -> readproj -> generate -> check (go vet + go test -race) -> log -> done
  logged edit to main.go
  [exit 0]   # first try, no repair
```

Increment 1 (HTTP server) is in: server.go (handler + Routes) + main.go (production server). The module
still passes `go vet` + `go test -race`.

## Step R2.1c - harden gate caught a staticcheck issue, fixed via edit_file

```text
$ ratchet flow . harden --ws cacheproxy ""        # gate
  ok (gofmt) (go vet) (go build) (go test -race)
  FAIL (staticcheck): main.go:44:33: os.Kill cannot be trapped (SA1016)
  NOT CLEAN  [exit 1]
```

staticcheck is advisory in the per-unit oracle (so the 1b edit passed) but GATES in harden - caught the
SIGKILL-can't-be-trapped bug. Fixed by feeding the finding into edit_file:

```text
$ ratchet flow . edit_file --ws cacheproxy "main.go fix SA1016: os.Kill -> syscall.SIGTERM"
  ... -> done   (first try)
$ ratchet flow . harden --ws cacheproxy ""
  ok (staticcheck) ; PRODUCTION-CLEAN
```

## Step R2.1d - prove it serves (built binary + real curl)

```text
$ go build -o /tmp/cacheproxy . && /tmp/cacheproxy &     # listening on :8080
$ curl localhost:8080/get/alpha   -> 200  value-for-alpha-call-2   in 0.0106s   (miss -> fetch)
$ curl localhost:8080/get/alpha   -> 200  value-for-alpha-call-2   in 0.0017s   (cache HIT, ~6x faster, same value)
$ curl localhost:8080/get/beta    -> 200  value-for-beta-call-3    in 0.0108s   (new key -> fetch)
```

Increment 1 done: the proxy serves over HTTP; a warm key is served from cache (no refetch, much faster).

## Step R2.2 - edit_file cache.go (single mutex -> 32 FNV shards), contract preserved

```text
$ ratchet flow . edit_file --ws cacheproxy "cache.go shard into 32 FNV shards; keep Get/Set signatures"
  read -> readproj -> generate -> check (go vet + go test -race) -> log -> done
  logged edit to cache.go
  [exit 0]   # first try, no repair
```

The advanced edit: rewrote the cache internals to per-shard RWMutex locking while keeping the exact
Get/Set contract - so proxy.go and the test were untouched, and the EXISTING `-race` concurrency test
passed as-is. The contract held and the refactor is race-clean.

## Step R2.2b - harden (sharded cache)

```text
$ ratchet flow . harden --ws cacheproxy ""
  ok (gofmt) (go vet) (go build) (go test -race) (staticcheck) (govulncheck)
  PRODUCTION-CLEAN  [exit 0]
```

## Step R2.3 - add_file cache_bench_test.go + run benchmarks

```text
$ ratchet flow . add_file --ws cacheproxy "cache_bench_test.go b.RunParallel benchmarks ..."
  ... generate -> check FAIL -> fix -> recheck -> register -> done   (1 repair)
  registered cache_bench_test.go

$ go test -bench . -benchmem -run='^$' -cpu=8        # (manual: no bench flow yet)
  BenchmarkCacheSet-8    9507300    122.3 ns/op    36 B/op    1 allocs/op
  BenchmarkCacheGet-8   46355934     26.11 ns/op    0 B/op    0 allocs/op
  BenchmarkProxyGet-8   42683137     27.68 ns/op    5 B/op    1 allocs/op
```

Concurrent reads are ~26 ns with zero allocations under 8-way parallelism - the sharded RWMutex read
path. (A `bench` ratchet flow/tool is a natural Round 3 addition; the run above was a manual go test.)

## Round 2 increments 1-3: complete
- HTTP server (additive + main edit), hardened, served via curl (cache hit ~6x faster).
- Sharded cache (32 FNV shards), contract preserved, existing -race test passed, hardened clean.
- Benchmarks (b.RunParallel) reporting throughput/allocs.
Remaining: the context-propagation probe (the cross-cutting, multi-file edit experiment).

## Step R2.probe - context propagation (the cross-cutting edit) -> FINDINGS

```text
$ ratchet flow . edit_file --ws cacheproxy "fetcher.go add context.Context to Fetch"
  read -> readproj -> generate -> check FAIL -> fix -> recheck FAIL -> edit_file.fail
  VET FAILED after staging fetcher.go:
  vet: ./proxy.go:28:36: not enough arguments in call to p.fetcher.Fetch
       have (string) want (context.Context, string)
  [exit 1]
```

Two findings (the point of the probe):
1. CROSS-CUTTING edits don't fit single-file edit_file: the breakage is in proxy.go, which a fetcher.go
   edit cannot touch; the whole-module oracle correctly refuses. Needs a multi-file / transactional edit
   flow (read N files + one request -> regenerate all -> verify once).
2. edit_file has NO ROLLBACK: it staged the new fetcher.go and left it on failure, so the module is now
   broken (fetcher.go has ctx, proxy.go does not). edit_file should snapshot/restore the target on
   failure, like the bump tool does.

State after probe: workspace cacheproxy does not build (fetcher.go updated, callers not).
