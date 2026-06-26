# AGENTS.md - the cpp ratchet (for AI agents)

A C++ ratchet for the [Ratchet](https://github.com/scanset/Ratchet) engine: the local model writes
C++, MSVC `cl` is the Oracle. Humans: see [README.md](README.md). Operating rules: [SYSTEM.md](SYSTEM.md).
Layout: [STRUCTURE.md](STRUCTURE.md).

## Flows

- **`cpp`** - one-off generate -> compile-check (`cl /Zs`) -> repair x2. No project.
- **`add_file`** - plan -> generate a new file -> stage + build the whole project -> repair x2 -> register.
- **`edit_file`** - read the file -> plan -> whole-file rewrite -> stage + build -> repair x2 -> log.
- **`answer`** - plan -> retrieve -> explain. Grounded Q&A, no build.
- **`compose`** - multi-spec composition: `read_specs` -> infer an ordered build plan (`compose.plan`,
  schema-forced) -> `plan_units` worklist -> `foreach add_unit` -> whole-project build. Specs live in
  `<ws>/specs/`.
- **`add_unit`** - the compose sub-flow: generate one unit against its spec + the existing units'
  headers (`project_api`) -> stage + build -> repair x2 -> register.

The four generating one-off/lifecycle flows start with a **plan** node (a `generate` with an
`output_schema` emitting one search query per KB, empty = skip); downstream `search` bindings pull each
via the binding `path`, routing retrieval to exactly the KBs the task needs.

## Knowledge bases

`cppref` (stdlib/language), `guidelines` (Core Guidelines), `patterns` (design patterns), `errors`
(MSVC diagnostics, used in repair), `cppdocs` (MSVC compiler/CRT/build), `win32` (curated desktop +
DirectX/MediaFoundation/Winsock), `howto` (using this ratchet - adding packages). After editing a KB,
`ratchet index <this>\kb\<lib>`.

## Tools

`cl_check` (single-unit oracle), `new_project` (console|win32 scaffold), `build_project` (compile +
link -> exe; honors `subsystem`/`libs`/`includes`/`libdirs`; applocal-deploys vcpkg DLLs; compiles
`/utf-8`), `stage_and_build` (write + build, the add/edit oracle), `stage_unit` + `project_api` +
`read_specs` + `plan_units` (composition: stage a unit's .h/.cpp, emit the existing headers as the API,
extract the worklist), `transcript` (assemble a run's prompt->response transcript), `link_libs`,
`bootstrap_vcpkg` + `add_package` (vcpkg, x86-windows triplet), `doctor_cl`, plus the language-agnostic
`read_file`/`read_project`/`register_file`/`log_edit`/`run_app`.

## Composition

A `.spec` is a structured prompt (`name`/`intent`/`behavior`/`constraints`), not a parsed grammar.
`compose` builds a multi-file program from a dir of specs, in dependency order, each unit generated
against the units already built. The C++ specifics:

- A component is a `src/core/<Name>.h` **declaration** + a `src/core/<Name>.cpp` **definition**, emitted
  together via `// === <path> ===` markers and split by `stage_unit`. The entry is `src/main.cpp`. `cl`
  globs every `.cpp`, so each component compiles + links on its own (real separate compilation).
- `project_api` surfaces the prior units' **headers** - in C++ the header *is* the contract, so this
  closes the multi-reference frontier (a unit calling several others gets their exact signatures).
- An interface (pure abstract base) is `.h`-only; the generate prompt emits no `.cpp` for it.

Run it: `ratchet flow <this> compose --ws <project> ""`. Worked example (compose + a concurrent app +
an author-correction pass): `transcripts/warehouse_transcript.md`.

## Hard constraints (these bite)

- **C++17, MSVC `cl`** (32-bit via `vcvars32`). Generated code must compile; the repair loop feeds `cl`
  errors back. Modern libs (fmt) need `/utf-8` - `build_project`/`cl_check` already pass it.
- **stdin to tools** arrives with a UTF-8 BOM - strip it (read raw bytes), and feed a child process's
  stdin from a temp file, not the managed pipe (it garbles native input). See `run_app`.
- **PowerShell 5.1:** don't pipe native `git`/`cmd` through `2>&1` under `ErrorActionPreference=Stop`
  (wraps stderr as a terminating error); a `[type]`-constrained param coerces later same-named
  assignments (name parsed objects distinctly).
- **Smart App Control** gates *running* native exes per-binary; the **build** (`build_project` exit 0)
  is the deterministic verification.

## Verify

```
ratchet validate-flow <this>     # lint all chains (one-offs + lifecycle + compose/add_unit + answer)
ratchet doctor <this>            # preflight cl, SDK, Ollama+models, KB integrity (requirements in ratchet.json)
```
