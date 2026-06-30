# Transcript: Cache — a TTL key-value cache (local model, DONE — with one honest caveat)

**Plain-language summary:** We asked the small local model to build an in-memory cache where each entry
expires after a time-to-live (TTL). It passed every gate on the first implementation attempt. But it also
shipped a **subtle concurrency flaw the tests never exercised** — a useful reminder that the gates are only
as strong as the test that drives them.

---

## 1. The input task (the spec) — and why it's written this way

```
name: Cache   role: component
api:
  - func NewCache() *Cache
  - func (c *Cache) Set(key string, value string, ttl time.Duration)
  - func (c *Cache) Get(key string) (string, bool)
behavior:
  - After Set(key,value,ttl) with ttl>0, an immediate Get returns (value, true).
  - EXPIRY: once time since Set exceeds ttl, Get returns ("", false). Store an expiry
    instant (now.Add(ttl)); fresh while time.Now is before expiry.
  - A ttl of zero or negative means already expired: Get returns ("", false). Never panic.
  - CONCURRENCY: Set and Get may be called from many goroutines; guard the map (no race under -race).
```

**Input reasoning:** TTL caches have two classic traps — (a) the zero/negative-TTL edge (does it panic or
treat as expired?) and (b) concurrent access to the map. The spec pins both explicitly, and even tells the
model *how* to represent freshness (store `now.Add(ttl)`, compare to `time.Now`). This is "pin the
contract" again: the previously observed negative-TTL bug is turned into a stated, checkable rule.

## 2. What the model produced (output)

**Test (written first) — covers hit, miss, expiry, the zero/negative edge, replace, concurrency, + a fuzz target:**

```go
t.Run("Expiry", func(t *testing.T) {
    c := NewCache(); c.Set("k", "v", time.Millisecond)
    time.Sleep(5 * time.Millisecond)
    if got, ok := c.Get("k"); got != "" || ok { t.Errorf(...) }   // expired -> ("", false)
})
t.Run("Edge", func(t *testing.T) {
    c := NewCache()
    c.Set("k", "v", 0);            if _, ok := c.Get("k"); ok { t.Errorf(...) }   // ttl 0
    c.Set("k", "v", -time.Second); if _, ok := c.Get("k"); ok { t.Errorf(...) }  // ttl < 0
})

func FuzzCache(f *testing.F) {
    f.Fuzz(func(t *testing.T, key, value string, ttlMs int64) {
        c := NewCache(); ttl := time.Duration(ttlMs) * time.Millisecond
        c.Set(key, value, ttl)
        got, ok := c.Get(key)
        if ttl <= 0 { if got != "" || ok { t.Errorf(...) } }        // negative-TTL edge, fuzzed
        else        { if got != value || !ok { t.Errorf(...) } }
    })
}
```

The fuzz target derives a `time.Duration` from a fuzzed `int64` (the idiom from the grounding) and asserts
the right behavior on *both* sides of the zero boundary — exactly the edge that bit an earlier run.

**Implementation (green rung):**

```go
type Cache struct { mu sync.RWMutex; items map[string]item }
type item struct { value string; expiry time.Time }

func (c *Cache) Set(key, value string, ttl time.Duration) {
    c.mu.Lock(); defer c.mu.Unlock()
    c.items[key] = item{value, time.Now().Add(ttl)}   // ttl<=0 -> expiry in the past -> already expired
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock(); defer c.mu.RUnlock()
    it, ok := c.items[key]
    if !ok { return "", false }
    if time.Now().After(it.expiry) { delete(c.items, key); return "", false }  // <-- see caveat
    return it.value, true
}
```

The zero/negative-TTL handling is elegant: `now.Add(ttl)` with a non-positive `ttl` lands in the past, so
the freshness check treats it as expired with no special-casing. The functional behavior is correct.

## 3. The oracle's verdicts

```
green  : OK: cache.go staged; module vets and tests pass (-race)
fuzz   : all targets clean (5s each)
harden : PRODUCTION-CLEAN — all gates pass
```

Flow path: `reset → read → stub → stubwrite → test → red → impl → green → fuzz → harden → DONE`, 11 steps,
first implementation attempt.

## 4. The honest caveat — a latent bug the gates missed

`Get` calls `delete(c.items, key)` — a **map write** — while holding only `RLock` (a *read* lock). Under a
`sync.RWMutex`, multiple readers share an `RLock`, so writing the map there is a data race (and Go's runtime
can panic on a concurrent map write). A correct version would take a full `Lock` in `Get`, or not delete in
the read path.

**Why no gate caught it:** the race detector only flags races that *actually happen during the test*. The
test's concurrency case used `time.Minute` TTLs, so entries never expired during the run and the
`delete`-under-`RLock` path never executed concurrently. The bug is real but un-exercised.

**The lesson (and the point of capturing it):** the oracle is exactly as strong as the test that drives it.
This is precisely the kind of signal the workflow should surface — not "the model is bad," but "the test
left a path uncovered, route a senior (or a stronger model) to add the expiring-concurrent case." It is a
*test-coverage* gap, not a capability wall.

## 5. Bottom line

A functionally-correct TTL cache (including the nasty zero/negative-TTL edge), race-clean on the tests as
written, first attempt, **local model only**. The residual `delete`-under-`RLock` flaw is a clean,
actionable signal for review — the kind of finding the workflow is meant to localize.
