You are given several `.spec` files (each a short structured prompt) that together describe ONE Go
program. They are in NO particular order and may use slightly inconsistent names for the same thing.

Produce a BUILD PLAN as structured data:

- `units`: every unit, in DEPENDENCY ORDER - data/foundation types first, then things that depend on
  them, and the program entry LAST. For each unit give its `name`, its `role` (one of: data, interface,
  component, behavior, gui), `dependsOn` (the names of earlier units it needs; empty for foundations),
  and `spec` (the source `.spec` filename it came from). Exactly ONE unit should have role `behavior`
  or `gui` - the program entry point (it becomes `main.go` with `func main`); the rest are
  data/interface/component files. A spec that describes a Go test (it asks for `func TestXxx`) gets
  role `test` and becomes a `_test.go` file - INCLUDE every such spec as its own unit, ordered LAST
  (after the code it tests). Every unit lives in the same `package main` at the module root, so units
  call each other directly with no imports.
- `contracts`: the canonical shared names/types every unit must agree on. Give the single canonical
  `name` (prefer the name from the spec that DEFINES the thing over one that merely mentions it) and a
  short `type`. Set `normalizedFrom` ONLY when the specs actually used a different name for the same
  thing; leave it empty when they were already consistent.

Infer roles and order from intent and the cross-references between specs.

SPECS:
{{ specs }}
