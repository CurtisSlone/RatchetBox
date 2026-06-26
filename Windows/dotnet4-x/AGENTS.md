# AGENTS.md - the dotnet ratchet (for AI agents)

A C# / PowerShell ratchet for the [Ratchet](https://github.com/scanset/Ratchet) engine: the local
model writes C# / PowerShell, the in-box `csc` / PS parser is the Oracle. Humans: see
[README.md](README.md). Operating rules: [SYSTEM.md](SYSTEM.md). Layout: [STRUCTURE.md](STRUCTURE.md).

## Flows

- **`csharp`** - one-off generate -> `csc` compile-check -> repair.
- **`winforms`** - one-off WinForms generate -> compile -> repair.
- **`powershell`** - one-off generate -> parse-check -> repair.
- **`add_file`** / **`edit_file`** - project lifecycle: generate/rewrite -> stage + build the whole
  project (the Oracle) -> repair -> register/log.
- **`compose`** - multi-spec composition: `read_specs` -> infer an ordered build plan (units +
  dependencies + canonical contracts + each unit's `module`, schema-forced, in the `compose.plan` node)
  -> `plan_units` worklist -> `foreach add_unit` (each unit generated against the project so far via
  `read_project`, so contracts accumulate) -> whole-project build. Specs live in `<ws>/specs/`; the entry
  unit -> `src/Program.cs`, others -> `src/<module>/<Name>.cs` (the spec's `module:`, default `Core`).
- **`add_unit`** - the compose sub-flow (called by `compose`'s foreach, not usually run directly):
  generate one unit against its spec + the project so far -> stage+build -> repair x2 -> register.

## Tools

`csc_check` (compile oracle), `new_project` (console|winforms scaffold: `src/`, `response.rsp`,
`build.ps1`, `project.json`), `build_project` (csc via the response file -> `dist/<name>.exe`),
`stage_and_build` (write + build, the add/edit oracle), `make_launcher` (SAC-safe in-memory `.cmd`),
`read_specs` (read a dir or one `.spec` file), `plan_units` (plan JSON -> the ordered `foreach`
worklist), `project_api` (the existing units' public signatures, bound into each composed unit),
`transcript` (assemble a run's prompt->response transcript), plus the language-agnostic
`read_file`/`read_project`/`register_file`/`log_edit`/`run_app`.

## Composition

A `.spec` is a **structured prompt** (`name`/`intent`/`behavior`/`constraints`, + optional `module:` =
the folder under `src/` the unit lands in, default `Core`), not a parsed grammar. Folders organize files;
the namespace stays flat `App` (so cross-unit refs need no extra `using`). `compose` plans the build
order from a dir of specs and generates each unit in dependency order, each seeing the project built so
far, with `project_api` binding the already-built units' exact signatures (this closes the
multi-reference frontier - a unit that must coordinate several others' signatures at once). Interfaces
help further: they collapse N concrete contracts into one. Run it:

```
ratchet flow <this> compose --ws <project> ""
```

Worked end-to-end example (compose + a concurrent app + an author-correction pass):
`transcripts/warehouse_transcript.md`.

## Hard constraints (these bite)

- **C# 5 only** (pre-Roslyn `csc`): NO string interpolation (`$"..."`), NO `?.`, NO expression-bodied
  members, NO tuples. The flow prompts forbid these because the local model reaches for them; the Oracle
  rejects what slips through.
- **stdin to tools** carries a UTF-8 BOM - strip it (read raw bytes).
- **Smart App Control** blocks unsigned exes - run built apps via `make_launcher`'s `.cmd`, never the
  bare exe; never disable SAC.

## Verify

```
ratchet validate-flow <this>     # lint all chains (one-offs + project lifecycle + compose/add_unit)
ratchet doctor <this>            # preflight csc, Ollama + models, KB integrity
```

Real end-to-end transcripts live in `transcripts/` - read them to see the exact commands, generated code,
and oracle/build results, including where the model under-implemented a spec and the author-correction
pass that fixed it.
