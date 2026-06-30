name: Limiter
role: component
intent: The rate-limiter abstraction. A Limiter decides whether a given client may make a request right now under a sliding-window policy. Code depends on this interface so the Redis-backed and in-memory implementations are interchangeable (and the in-memory one makes tests run without Redis).
api:
  - type Limiter interface { Allow(ctx context.Context, client string) (bool, error) }
behavior:
  - "Allow reports whether the request from `client` is permitted now. It returns (true, nil) if the request is within the limit, (false, nil) if the client has exceeded the limit in the current window, and (false-or-true, err) only on an infrastructure error (the concrete impl decides the fail-open vs fail-closed policy)."
  - "Allow is called on every request and must be safe for concurrent use."
constraints: package main; standard library only (this file declares just the interface)
