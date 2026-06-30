name: CacheTest
role: test
intent: Test the TTL cache: a hit while fresh, a miss after expiry, the zero/negative-ttl edge, replacement, and concurrency safety.
api:
  - func TestCache(t *testing.T)
behavior:
  - "Hit: after Set(\"k\", \"v\", time.Minute), Get(\"k\") returns (\"v\", true)."
  - "Miss: Get on a key never set returns (\"\", false)."
  - "Expiry: after Set(\"k\", \"v\", d) with a very small d (e.g. time.Millisecond), wait past d (e.g. time.Sleep(5*time.Millisecond)), then Get(\"k\") returns (\"\", false)."
  - "Edge: Set(\"k\", \"v\", 0) and Set(\"k\", \"v\", -time.Second) must make Get(\"k\") return (\"\", false) - never panic."
  - "Replace: Set(\"k\",\"a\",time.Minute) then Set(\"k\",\"b\",time.Minute); Get(\"k\") returns (\"b\", true)."
  - "Concurrency: run many goroutines doing Set and Get concurrently and assert it does not panic or race (run under -race)."
constraints: package: main
