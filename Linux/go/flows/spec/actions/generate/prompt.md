Draft Go build SPEC(S) for the request below. A spec is a short, structured prompt that the `compose`
flow later turns into code - you are writing the SPEC, not the code. Output ONLY marker blocks.

If the request is a single function/type, emit ONE spec. If it is a system, DECOMPOSE it into multiple
specs - ONE spec per distinct TYPE or component, plus the entry point and a test for the core behavior.
Keep a type's definition AND its methods/behavior in the SAME spec - NEVER split one type across a
separate "data" spec and a "component" spec (that produces duplicate declarations downstream). Build order
does not matter here (compose plans it).

Emit each unit as a marker block in EXACTLY this shape:

=== <name>.spec ===
name: <PascalCase unit name>
role: <one of: data, interface, component, behavior, gui, test>
intent: <one line - what this unit is for>
api:
  - <exact Go signature(s) to implement, for a type/component>
behavior:
  - <what it should do or assert, for the entry or a test>
constraints: <stdlib limits, concurrency notes, `package main`, or `package: <subdir>` for a sub-package>

Rules:
- Every spec needs `name`, a valid `role`, and at least one of `api`/`behavior`/`intent`.
- Exactly ONE unit has role `behavior` or `gui` (the entry point -> main.go). Include a `role: test`
  unit that exercises the core behavior.
- Keep names consistent ACROSS specs so the units fit together (one unit's `api` is another's call).
- Use the reference material to name the RIGHT concerns: e.g. `atomic.Int64` for a concurrent counter,
  `http.Server` timeouts + `http.MaxBytesReader` + graceful shutdown for a server, a non-blocking
  channel send for low latency. Bake those into the specs' `api`/`constraints`.
- Output ONLY the marker blocks - no prose before/after, no code fences.

## Request
{{ desc }}

## Reference (retrieved from the KB(s) selected for this task; any section may be empty)
### {{ kb1 }}
{{ refs1 }}
### {{ kb2 }}
{{ refs2 }}
### {{ kb3 }}
{{ refs3 }}
