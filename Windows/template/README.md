# template - a starter ratchet (lifecycle + composition)

The skeleton to copy when starting a new ratchet for the
[Ratchet](https://github.com/scanset/Ratchet) engine. It ships the GENERIC machinery working and the
DOMAIN-SPECIFIC pieces as clearly-marked `CHANGE_ME` stubs - implement those for your language and the
flows light up. `ratchet validate-flow` passes out of the box.

## What's inside

Four working flow graphs:

- **add_file** - generate ONE new file -> build the whole project (the oracle) -> repair x2 -> register.
- **edit_file** - read a file -> rewrite it -> build -> repair x2 -> log the edit.
- **compose** - read `<ws>/specs/*.spec` -> plan ordered units + contracts -> `foreach` **add_unit** -> build.
- **add_unit** - generate one composed unit against the project's existing API (`project_api`) -> build ->
  repair x2 -> register. (Single-file-per-unit model; for a header+source language, emit a pair - see
  `RatchetBox/cpp`.)

Tools (`tools/`, declared in `tools/manifest.json`):

- GENERIC, work as-is: `read_specs`, `plan_units`, `read_project`, `register_file`, `log_edit`, `read_file`.
- **CHANGE_ME** (your domain): `build_project` (THE ORACLE - the whole-project build, exit 0 = built),
  `new_project` (the scaffold), `project_api` (emit existing units' public API so new units call them
  exactly), plus the source-extension/entry knobs at the top of `plan_units` and the path-guard in
  `stage_and_build`.

Prompts (`flows/*/actions/*/prompt.md`) carry the generic structure with a `<LANGUAGE>` placeholder and
`CHANGE_ME` notes - set your language, its constraints, and your module/namespace convention.

## Instantiate a new ratchet

1. Copy this folder and rename it.
2. Edit `ratchet.json`: `name`, `domain`, `models` (point at models you have in Ollama), `requirements`
   (what `ratchet doctor` checks - add your toolchain).
3. Implement the `CHANGE_ME` tools: `build_project` (exit 0 = built), `new_project` (scaffold
   project.json + src/ + an entry), `project_api` (the existing-units API surface). Set the source
   extension/entry in `plan_units` and the path-guard in `stage_and_build`.
4. Fill the `<LANGUAGE>`/`CHANGE_ME` placeholders in the flow prompts.
5. (Optional) drop reference docs into `kb/<lib>/`, then `ratchet index <dir>\kb`.
6. Verify: `ratchet validate-flow <dir>` (structure) and `ratchet doctor <dir>` (toolchain).
7. Use it: `ratchet flow <dir> new_project <proj>` style, write `<proj>/specs/*.spec`, then
   `ratchet flow <dir> compose --ws <proj> ""`.

Working references that implement every `CHANGE_ME` piece: `RatchetBox/dotnet4-x` (C#, in-box csc) and
`RatchetBox/cpp` (C++, MSVC cl). Copy from whichever is closest to your domain.

## The idea

A ratchet is just data; the engine never changes per domain. The model proposes into a chain's slots, a
deterministic Oracle (here `build_project`) accepts or rejects, and the chain advances only on a pass -
like a ratchet that turns one way. Composition extends this: a directory of `.spec` files becomes an
ordered build plan, and each unit is generated against the real API of the units already built.

See `STRUCTURE.md` for the full layout and the engine's docs (in the Ratchet repo:
`docs/how-to/build-a-ratchet.md`, `docs/how-to/author-flows.md`, `docs/how-to/compose-from-specs.md`)
for the authoring contract.
