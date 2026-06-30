name: Cache
role: component
intent: An in-memory key-value cache with per-entry time-to-live (TTL). A stored value is returned only while it is still fresh; once its TTL has elapsed it is treated as absent. The cache is used concurrently, so it must be safe for use by multiple goroutines.
api:
  - func NewCache() *Cache
  - func (c *Cache) Set(key string, value string, ttl time.Duration)
  - func (c *Cache) Get(key string) (string, bool)
behavior:
  - "NewCache returns an empty cache. Get on a missing key returns (\"\", false)."
  - "After Set(key, value, ttl) with ttl > 0, an immediate Get(key) returns (value, true)."
  - "EXPIRY: once the time since the Set exceeds ttl, Get(key) returns (\"\", false). Use time.Now and store an expiry instant (now.Add(ttl)); an entry is fresh while time.Now is before its expiry."
  - "A ttl of zero or negative means the entry is already expired: Get returns (\"\", false). Never panic on a zero or negative ttl."
  - "Set on an existing key replaces both its value and its expiry."
  - "CONCURRENCY: Set and Get may be called from many goroutines at once; guard the map with a sync.Mutex or sync.RWMutex so there is no data race under go test -race."
constraints: package: main
