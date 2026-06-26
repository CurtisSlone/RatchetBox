# Recipes + framework profiles (third-party frameworks, end to end)

Two things landed here:

- **`recipes/` KB** - playbooks per app *type* (a JSON HTTP API, a worker pool, a flag CLI) plus
  framework *profiles* (chi web service, cobra CLI). Each is a short "when to use + the API shape + real
  code" doc. Wired into the `spec` drafter and `add_unit`, so describing "a web API" pulls the recipe and
  the units come out shaped right.
- **Framework support rides `add_dep`.** The ratchet was already framework-agnostic: `add_dep` `go get`s
  any module and ingests its `go doc` into the `deps` KB. A "profile" is just the recipe + that ingest.

## Proof: a chi web service, built from scratch and hardened

```text
$ /do new_module chiweb
$ /do add_dep chiweb github.com/go-chi/chi/v5      # go get + ingest chi's go doc into the deps KB
  ingested github.com/go-chi/chi/v5 -> kb/deps/github_com_go_chi_chi_v5.md
```

The model then wrote a router grounded on the ingested chi docs + the chi recipe - using the REAL chi
API, not a hallucination:

```go
// add_file router.go  (first try, vet-clean)
package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
)

func Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK); w.Write([]byte("ok"))
	})
	r.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")        // exact chi API, from the ingested docs
		w.Write([]byte("Hello, " + name + "!"))
	})
	return r
}
```

A `main` was added (production `*http.Server` + graceful shutdown, grounded on the recipe), and the
complete service builds and passes the full gate:

```text
$ go build ./...                 -> BUILD OK
$ /flow harden --ws chiweb
ok   (gofmt) (go vet) (go build) (go test -race) (staticcheck)
ok   (govulncheck (known CVEs))   # scanned the chi dependency - no known CVEs
PRODUCTION-CLEAN
```

## The takeaway

Supporting a framework is not new plumbing - it is `add_dep` (ingest its docs) + a `recipes/` profile
(the API shape). The local model then writes correct framework code because it is grounded on the real
package surface, and `harden`/`govulncheck` vets the dependency for known vulnerabilities. Add `cobra`,
`pgx`, `grpc`, etc. the same way.
