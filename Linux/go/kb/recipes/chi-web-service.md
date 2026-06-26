# Recipe / framework profile: a chi web service

`github.com/go-chi/chi/v5` is a lightweight, idiomatic HTTP router (stdlib-compatible `http.Handler`).
Use it when you want middleware and URL params without a heavy framework.

Profile (one-time per workspace): pull the dependency and ingest its docs so generation grounds on the
real API:

```
/do add_dep <workspace> github.com/go-chi/chi/v5
```

Then build handlers against `chi.Router`. Key API: `chi.NewRouter()`, `r.Use(middleware...)`,
`r.Get/Post/Put/Delete(pattern, handlerFunc)`, URL params via `chi.URLParam(r, "id")`, sub-routers via
`r.Route("/prefix", func(r chi.Router){...})`. A `chi.Router` is an `http.Handler`, so wrap it in the
production `*http.Server` (timeouts + graceful shutdown).

```go
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Routes(s *Server) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	r.Route("/items", func(r chi.Router) {
		r.Post("/", s.createItem)
		r.Get("/{id}", s.getItem) // id := chi.URLParam(r, "id")
	})
	return r
}
```

Then serve it with the production-http-server pattern (an `*http.Server` with timeouts +
`signal.NotifyContext` + `Shutdown`). Verify with the `harden` flow (govulncheck will flag any known
CVE in the chi version you pulled).
