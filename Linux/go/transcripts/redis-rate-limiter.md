# Transcript: Redis-backed rate limiter — a complex app with a real external service (local model)

**Plain-language summary:** We asked the small local model (qwen3-coder 30B-A3B, ~3B active) to build a
real microservice: a **sliding-window rate limiter backed by Redis**, exposed over HTTP, with the circuit
breaker + retry resilience from our earlier work wrapping the Redis calls. Six interdependent units, a
genuine third-party dependency (`go-redis`), and a live Redis to talk to. The model built it; it works
end-to-end against real Redis. The only blockers were two deterministic tooling bugs, both fixed.

---

## The app (6 units, package main, + the go-redis dependency)

| Unit | Role |
|---|---|
| `Limiter` | interface: `Allow(ctx, client) (bool, error)` |
| `RedisLimiter` | go-redis sorted-set sliding window + a fail-open circuit breaker |
| `MemoryLimiter` | in-memory sliding window (tests run without Redis) |
| `Breaker` | circuit breaker (reused from the dispatcher) |
| `Server` | net/http rate-limit middleware -> 429 when over limit |
| `Main` | flags, connect Redis, wire, serve |

## Setup — bringing the real service in

- **Redis** ran in Docker (`docker run redis:7-alpine`).
- **The dependency** was added with the ratchet's own `add_dep`, which runs `go get` + `go mod tidy` and
  ingests `go doc -all github.com/redis/go-redis/v9` (~1.1 MB) into the `deps` KB so compose can ground
  generation on the real API surface.
- **Testability principle:** Redis sits behind the `Limiter` interface, so `go test` uses `MemoryLimiter`
  (no Redis), and the real `RedisLimiter` is exercised end-to-end. Same dependency-injection move as the
  dispatcher's `Deliverer`.

## The headline result — spec-pinning carries a THIRD-PARTY API, not just stdlib

The `RedisLimiter` spec pinned the exact go-redis calls. The local model produced correct, idiomatic
go-redis code from it:

```go
func (l *RedisLimiter) Allow(ctx context.Context, clientID string) (bool, error) {
	now := time.Now().UnixNano()
	key := "ratelimit:" + clientID

	if !l.breaker.Allow() {
		return true, nil // FAIL OPEN
	}

	seq := l.seq.Add(1)
	pipe := l.client.TxPipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(now-int64(l.window), 10))
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d-%d", now, seq)})
	count := pipe.ZCard(ctx, key)
	pipe.Expire(ctx, key, l.window)
	_, err := pipe.Exec(ctx)

	if err != nil {
		l.breaker.Failure()
		return true, err // FAIL OPEN
	}

	l.breaker.Success()
	allowed := count.Val() <= l.limit
	return allowed, nil
}
```

A correct sliding window (drop-old, add-this, count, expire), a unique member via an atomic counter so two
requests in the same nanosecond don't collide, and a breaker that fails OPEN on a Redis outage. The
`deps`-KB grounding was thin (2500 chars of a 1.1 MB doc - just the package head), so the real grounding
was the **pinned spec** plus the model's training. The lesson from the stdlib work holds for dependencies:
pin the contract and the model fills it in.

## What actually broke — two POST-PROCESSING bugs (not context, not capability)

The first compose run failed the whole-module build. Be precise about *where*: a failure can live in three
distinct places, and these two were neither of the first two.

- **Context assembly** - what the model is *given* (bindings, grounding, retrieval). Here: adequate. The
  spec pinned the go-redis calls; the grounding was thin (see below) but sufficient.
- **Generation** - what the model *produces*. Here: correct. The go-redis code was right, import included.
- **Post-processing** - what deterministic tools do to that correct output (prune, gofmt, tag, stage).
  **Both bugs were here** - tools mangling or under-processing output the model had already gotten right.

1. **`prune_imports` deleted the model's valid go-redis import.** It infers a package name from the import
   path's last element - `github.com/redis/go-redis/v9` -> `v9` - finds `v9.` unused, and removes the
   import *after* generation. But the package is named `redis`; NO path element is `redis`, so
   path-inference can't win. Fix: `prune_imports` now only prunes **stdlib** imports (first path segment has
   no dot) and never touches third-party ones. Silently fine until the first dependency showed up.
2. **compose's `stage_build` skipped the prune/tag pass.** The tdd staging strips unused imports and stamps
   the keyword tags; the compose path did neither, so the model's `server.go` unused-`time` slip wasn't
   auto-fixed (its one repair couldn't remove it). Fix: `stage_build` now runs `prune_imports` + `code_tags`
   before vet, matching the tdd path.

So this build had **no context-assembly hole at all.** The thin `deps`-KB grounding (2500 chars of a 1.1 MB
`go doc` - just the package head, not the sorted-set methods) was a *latent* weakness in the assembled
context, but it never bit: the pinned spec backstopped it. The canonical *context*-assembly hole this
effort produced was elsewhere - the dispatcher's `module_api` omitting constants and interface bodies, so
the next unit's context was genuinely missing real names and it guessed `EventTypeEnqueued` for
`EventEnqueued`. That is what a context hole looks like; these two rate-limiter bugs are not that.

With both post-processing bugs fixed, compose re-ran and **all six units cohered, built, and vetted clean
with the real dependency** - and the in-code `// kw:` tags were auto-stamped along the way.

## End-to-end against live Redis

Server run with `-limit=3 -window=60s -redis=localhost:6379`, hammered from one client:

```
request 1 -> 200      request 4 -> 429
request 2 -> 200      request 5 -> 429
request 3 -> 200      request 6 -> 429
ZCARD ratelimit:<client> = 6
```

The `ZCARD = 6` is the proof that matters: every request hit **real Redis** and landed in the sorted set.
Had Redis been unreachable, the breaker would have failed OPEN and all six would have been 200 - so the
429s plus a populated sorted set confirm the actual sliding-window path executed against the live service.

## Lessons

1. **Spec-pinning extends to third-party dependencies.** Correct go-redis from a pinned spec + thin
   grounding - the same result the stdlib KB gave, now for an external API.
2. **The walls were in POST-PROCESSING, not the model - and the stateless pipeline made that pinpointable.**
   Both failures traced to one deterministic step run *after* generation: `prune_imports` mis-deriving a
   package name and deleting a valid import, and `stage_build` skipping the prune/tag pass. The model's
   context was adequate and its output correct; the bug was downstream. Because each pipeline step is a
   discrete, inspectable transform, "the build failed" became "this exact tool mangled correct output" -
   not "the model is confused." Both fixes now benefit every future build that uses a dependency.
3. **One-shot compose works once the pipeline holes are closed - and layering is still the next step.** This
   built breadth-first in one pass. The same service is the ideal worked example for the layered (depth)
   methodology: L0 interface+memory -> L1 add Redis -> L2 wrap in the breaker -> L3 harden, each a small
   stateless pass over the verified workspace, each locating its edit target via the code-search/tags.

## Teardown

The Redis container was removed after verification; the workspace (git-ignored scratch) holds the generated
service. No state persists beyond the artifact - consistent with the stateless-connection model: build the
context, verify, release.
