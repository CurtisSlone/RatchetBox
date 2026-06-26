name: URLStore
role: data
intent: An in-memory store mapping short codes to long URLs
api:
  - func NewURLStore() *URLStore
  - func (s *URLStore) Put(code string, url string) error
  - func (s *URLStore) Get(code string) (string, bool)
behavior:
  - Store should map short codes to long URLs
  - Get should return false for non-existent codes
  - Put should overwrite existing codes
constraints: Use sync.RWMutex for concurrent access; package: main
