# Recipe: a JSON HTTP API (standard library)

A small JSON-over-HTTP service using only `net/http` + `encoding/json`, with the production-server
shape (timeouts, body cap, graceful shutdown). Build it as: a store/component, handlers on a server
type, and a thin `main` that wires + serves.

- Decode requests with `json.NewDecoder(r.Body).Decode(&v)`; cap the body with `http.MaxBytesReader`.
- Encode responses with a small helper that sets `Content-Type: application/json` and `json.NewEncoder`.
- Route with `http.ServeMux` (Go 1.22+ supports method+path patterns like `"POST /items"` and
  `r.PathValue("id")`).
- Always use an `*http.Server` with timeouts and `signal.NotifyContext` + `Shutdown` (see the
  production-http-server pattern).

```go
package main

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type CreateItem struct {
	Name string `json:"name"`
}

func (s *Server) createItem(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var in CreateItem
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	id := s.store.Put(in.Name)
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /items", s.createItem)            // Go 1.22+ method+path
	mux.HandleFunc("GET /items/{id}", s.getItem)           // r.PathValue("id")
	return mux
}
```
