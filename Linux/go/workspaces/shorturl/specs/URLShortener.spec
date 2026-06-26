name: URLShortener
role: component
intent: Combines store and encoder to provide URL shortening functionality
api:
  - func NewURLShortener() *URLShortener
  - func (s *URLShortener) Shorten(url string) (string, error)
  - func (s *URLShortener) Expand(code string) (string, bool)
behavior:
  - Shorten should return a unique short code for each URL
  - Expand should return the original URL for valid codes
  - Should use the store and encoder components
constraints: Use atomic.Int64 for counter; package: main
