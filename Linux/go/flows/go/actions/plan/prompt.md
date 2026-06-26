Decide which reference libraries would help write the Go task below. Return a JSON object with one
search query per library - or an empty string `""` to skip a library that is not relevant. Skipping is
normal and preferred: most simple tasks need none.

- `stdlib_q`: a Go standard-library package or symbol query, when the task needs a specific stdlib API.
  Examples: "container/heap priority queue", "sort.Slice stable", "encoding/json struct tags",
  "strings.Builder", "regexp MatchString". Empty if plain language/arithmetic suffices.
- `patterns_q`: a design-pattern OR named-algorithm query, when the task maps to a known structure.
  Examples: "factory create product", "two-heap streaming median", "worker pool". Empty for ordinary code.
- `guidelines_q`: an idiomatic-style query from Effective Go / Code Review Comments, when style matters.
  Examples: "value vs pointer receiver", "error strings", "accept interfaces return structs". Empty if not.
- `pitfalls_q`: a builds-but-wrong trap query, when the task risks one. Examples: "nil map write",
  "loop variable capture goroutine", "channel deadlock", "slice aliasing append". Empty if not at risk.
- `idioms_q`: a general Go-idioms query. Examples: "error wrapping %w", "slice preallocate append".
  Empty if not needed.

Pick the narrowest queries that retrieve the right reference. Do not invent a need; prefer empty.

## Task
{{ task }}

Output ONLY the JSON object.
