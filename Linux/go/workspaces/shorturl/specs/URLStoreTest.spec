name: URLStoreTest
role: test
intent: Test the core behavior of the URL store
api:
  - func TestURLStore(t *testing.T)
behavior:
  - Store should correctly map codes to URLs
  - Get should return false for non-existent codes
  - Put should overwrite existing entries
  - Concurrent access should be safe
constraints: Use sync.RWMutex for concurrent access; package: main
