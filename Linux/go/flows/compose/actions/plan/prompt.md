You are given several `.spec` files (each a short structured prompt) that together describe ONE Go
program. They are in NO particular order and may use slightly inconsistent names for the same thing.

Produce a BUILD PLAN as structured data:

- `units`: every unit, in DEPENDENCY ORDER - data/foundation types first, then things that depend on
  them, and the program entry LAST. For each unit give its `name`, its `role` (one of: data, interface,
  component, behavior, gui), `dependsOn` (the names of earlier units it needs; empty for foundations),
  and `spec` (the source `.spec` filename it came from), and `pkg`. Exactly ONE unit should have role
  `behavior` or `gui` - the program entry point (it becomes `main.go` with `func main` at the module
  root); the rest are data/interface/component files. A spec that describes a Go test (it asks for
  `func TestXxx`) gets role `test` and becomes a `_test.go` file - INCLUDE every such spec as its own
  unit, ordered LAST (after the code it tests).
  - `pkg`: the subdirectory/package this unit belongs in. COPY the spec's `package:` line verbatim if it
    has one; otherwise leave it EMPTY. An empty pkg puts the file at the module root in `package main`
    (the default, flat layout). Units in the same pkg call each other directly; a unit in one pkg calls
    another pkg's EXPORTED names via `import "<module>/<pkg>"`. The entry (main) is always at the root.
    Prefer empty/flat unless the specs clearly call for separate packages.
- `contracts`: the canonical shared names/types every unit must agree on. Give the single canonical
  `name` (prefer the name from the spec that DEFINES the thing over one that merely mentions it) and a
  short `type`. Set `normalizedFrom` ONLY when the specs actually used a different name for the same
  thing; leave it empty when they were already consistent.

Infer roles and order from intent and the cross-references between specs.

SPECS:
{{ specs }}
