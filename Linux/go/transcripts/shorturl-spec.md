# spec - scaffolding specs from a description (the front-end to compose)

The `spec` flow turns a free-text description into well-formed `.spec` files - the only artifact you hand-write becomes a one-liner. One request becomes one spec (a single unit) or many (a system, decomposed). The oracle is well-formedness (`spec_check`): `write_specs` validates and writes atomically, repairing once. Grounded on patterns/guidelines/pitfalls so the specs name the right concerns BEFORE any code exists.

- Generated: 2026-06-26
- Command: `ratchet flow . spec --ws shorturl "<description>"` (after `new_module shorturl`)
- Result: 5 specs drafted first try, no repair, all well-formed.

## The request (one line)

```text
a URL shortener: an in-memory store mapping short codes to long URLs, a base62 encoder for ids, an HTTP API with POST /shorten returning a short code and GET /{code} doing a 302 redirect, and a test of the store
```

## The model's drafted output (marker-separated specs)

```text
=== URLStore.spec ===
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

=== Base62Encoder.spec ===
name: Base62Encoder
role: component
intent: Encodes integers to base62 strings for short URL codes
api:
  - func NewBase62Encoder() *Base62Encoder
  - func (e *Base62Encoder) Encode(n int64) string
behavior:
  - Encode should produce base62 strings using characters [0-9a-zA-Z]
  - Should handle zero and negative numbers
  - Should produce consistent output for same input
constraints: Use math/big for large numbers if needed; package: main

=== URLShortener.spec ===
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

=== HTTPServer.spec ===
name: HTTPServer
role: behavior
intent: HTTP API with POST /shorten and GET /{code} redirect
api:
  - func NewHTTPServer(shortener *URLShortener) *HTTPServer
  - func (s *HTTPServer) Start(addr string) error
  - func (s *HTTPServer) Shutdown() error
behavior:
  - POST /shorten should accept a JSON body with "url" field and return a short code
  - GET /{code} should redirect to the original URL with 302 status
  - Should handle request size limits and timeouts
  - Should gracefully shutdown on SIGINT/SIGTERM
constraints: Use http.Server with ReadHeaderTimeout, ReadTimeout, WriteTimeout, IdleTimeout; use http.MaxBytesReader; package: main

=== URLStoreTest.spec ===
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
```

## Oracle + write (well-formedness gate)

```text
wrote 5 spec(s) to /home/local/RatchetBox/Linux/go/workspaces/shorturl/specs: URLStore.spec, Base62Encoder.spec, URLShortener.spec, HTTPServer.spec, URLStoreTest.spec
OK: specs written; review them, then: ratchet flow . compose --ws shorturl ""
```

## Why this matters: the specs are hardened BY DEFAULT

The grounding (patterns + pitfalls) made the drafter name production concerns the request never mentioned - so the downstream `compose` builds them in without anyone remembering to ask:

- `URLStore` -> `constraints: Use sync.RWMutex for concurrent access`
- `URLShortener` -> `constraints: Use atomic.Int64 for counter` (the atomic-int64-alignment pitfall)
- `HTTPServer` -> `constraints: Use http.Server with ReadHeaderTimeout/ReadTimeout/WriteTimeout/IdleTimeout; use http.MaxBytesReader; graceful shutdown on SIGINT/SIGTERM` (the production-http-server pattern)

## Closing the loop - honest result

```text
$ ratchet flow . spec    --ws shorturl "a URL shortener: ..."   # describe -> 5 .spec files  [worked, first try]
$ # (review/edit the specs)
$ ratchet flow . compose --ws shorturl ""                        # specs -> code             [PARTIAL]
```

The spec drafting is solid; composing this particular (non-trivial) system was only partial:

| unit | result |
|---|---|
| urlstore.go | first try |
| urlshortener.go | first try |
| base62encoder.go | 1 repair, ok |
| main.go (HTTPServer) | FAILED after one repair (`"os" imported and not used`, unused `encoder`) |
| urlstore_test.go | FAILED after one repair |

Two real lessons (roadmap evidence, not a spec-flow bug):
- **Entry/component conflation.** The drafter made `HTTPServer` both the entry (role behavior -> main.go)
  AND a type with `Start`/`Shutdown`. compose then had to put a type definition AND `func main` in one
  file; the model left an unused import/var the single repair did not clean. Fix: the spec prompt should
  make the entry a THIN `main` that only wires the other components, with the server its own component.
- **One repair is thin for complex units.** compose could allow two repairs (like the C# `add_file`).

Cheap to recover: `/flow edit_file --ws shorturl "main.go remove the unused os import and unused encoder
variable"`. Takeaway: `spec` reliably drafts grounded specs; end-to-end auto-compose of a *larger* system
is at the edge of this local model and wants a thin-entry prompt + a second repair round.