# go - the cross-platform Go reference ratchet

Generates Go code for a focused task, verifies it with `go build`, and repairs once if it fails.
The companion to `dotnet4-x` (C#) and `cpp`, but toolchain-portable: the oracle is a `bash` script
the engine dispatches by extension, so it runs on Linux, WSL, and macOS out of the box.

## Requirements
- The Go toolchain on PATH (`go version`).
- Ollama reachable with the model seats pulled: `qwen3-coder` (generate/dispatch) and
  `nomic-embed-text` (embed). On WSL talking to Ollama on the Windows host, set
  `OLLAMA_URL=http://<windows-host>:11434` (the engine honors it; the config's localhost is a default).

## Build the KB index (once, after editing kb/)
```
ratchet index <path-to>/go/kb
```

## Use it
```
ratchet doctor <path-to>/go                       # preflight: go + toolbelt + models pulled
ratchet validate-flow <path-to>/go                # lint the flows (model-free)
ratchet flow <path-to>/go go   "a function that reverses a UTF-8 string"   # one-shot: build-only
ratchet flow <path-to>/go test "a function FizzBuzz(n int) string ..."     # one-shot: behavior
ratchet <path-to>/go                              # operator console; then: /flow go|test <task>
```

## The flows

**One-shot generation** (throwaway module, prints the code). Both open with a **plan** step (cpp-style
plan-routing): the model emits one search query per library (`stdlib`, `patterns`, `idioms`), an empty
query skips it, so generation is grounded only on what the task needs - a heap task pulls
`container/heap`, a FizzBuzz pulls nothing.
- **`go`** - plan -> generate -> `go build` -> repair once. Verifies the code *type-checks*.
- **`test`** - plan -> generate an implementation AND a test -> `go vet` + `go test` -> repair once.
  Verifies *behavior*, not just compilation - the headline of this ratchet. Emits two marker-separated
  files (`// === solution.go ===` / `// === solution_test.go ===`) the oracle splits and runs.

**Versions** (verified bumps, with rollback):
- **`upgrade`** - `ratchet flow . upgrade --ws <proj> "<pkg>@<ver> | go=<ver> | -u | tidy"` (or
  `/do bump <proj> <change>`): a version change is a change, so the oracle gates it - it snapshots
  go.mod/go.sum, applies the bump, `go mod tidy`s, runs the full `go_quality` gate (incl. `govulncheck`
  on the new version), and keeps it only if still clean, else rolls back.

**Frameworks** (any third-party module, grounded):
- `/do add_dep <proj> <module>` ingests a module's `go doc` into the `deps` KB, so generation writes
  correct API against it (proven with `github.com/go-chi/chi/v5` - the model used the real chi router
  API and the service built + passed `harden`). The `recipes` KB carries profiles for common ones (chi
  web service, cobra CLI) plus app-type playbooks (JSON API, worker pool, flag CLI), pulled by `spec`
  and `add_unit`.

**Spec authoring** (the front-end to compose):
- **`spec`** - `ratchet flow . spec --ws <proj> "<description>"` (after `new_module <proj>`): drafts
  well-formed `.spec` file(s) from a free-text description - one unit or a whole system (decomposed).
  Validated by the `spec_check` oracle and written into `workspaces/<proj>/specs/`. Grounded on
  patterns/pitfalls so the specs name the right concerns (e.g. it adds `atomic.Int64`, `http.Server`
  timeouts) before any code exists. Review the specs, then `compose`.

**Existing repos + refactor** (drive code you already have):
- `/do link_repo <name> <path>` symlinks any external Go module in as `workspaces/<name>`, so every flow
  below operates on your real repo (the oracles guard each change). **`refactor`** -
  `ratchet flow . refactor --ws <proj> "rename X to Y"` - type-safe rename across files via `gorename`,
  build-verified.

**Project lifecycle** (a persistent module under `workspaces/<proj>` you keep and grow):
- `/do new_module <proj>` - scaffold `workspaces/<proj>` (go.mod, PROJECT.md, specs/). All units live
  in one `package main` at the module root.
- **`add_file`** - `ratchet flow . add_file --ws <proj> "<path> <request>"`: generate a new file,
  grounded on the module's existing API, verified with `go vet` + `go test` over the whole module,
  repaired once, recorded.
- **`edit_file`** - `ratchet flow . edit_file --ws <proj> "<path> <request>"`: read the file's exact
  contents, apply the change (whole file rewritten, rest preserved), re-verify, repair once, log it.

**Composition from specs** (build a whole module at once):
- **`compose`** - `ratchet flow . compose --ws <proj> ""`: read `.spec` files in
  `workspaces/<proj>/specs`, plan the build order, generate each unit in dependency order against the
  module so far (`add_unit` per unit), then `go build ./...` + `go test ./...` the whole thing.

**Harden** (the production gate; deterministic, no model):
- **`harden`** - `ratchet flow . harden --ws <proj> ""` (or `/do go_quality <proj>`): run the full
  quality/security suite over a workspace - gofmt, `go vet`, `go build`, **`go test -race`**,
  `staticcheck`, **`govulncheck`** (known CVEs), `gosec` (SAST) - and pass only if every installed tool
  is clean. (The per-unit oracles already run `-race` and surface `staticcheck`; `harden` gates the lot.)

**Self-improvement** (the KB learns from the oracle):
- `/do mine_runs` scans `runs/` for recurring oracle failures and flags which are not yet covered by
  `pitfalls/`; **`learn`** (`ratchet flow . learn "<class>"`) drafts the missing pitfalls entry and adds
  it. The ratchet teaches itself from its own compiler - the more you run it, the lower the repair rate.

**Reference / learning** (grounded Q&A, no code generation, no oracle):
- **`explain`** - `ratchet flow . explain "<question>"` (or `/route "explain ..."`): answers a Go
  question in **prose**, grounded in the KB (Effective Go / Code Review Comments, stdlib, patterns,
  pitfalls). Use this for "how does X work / what's idiomatic" - unlike `go`/`test` it does not emit a
  `package solution` file. (Plain chat and `/search` go through the code-biased seat; `explain` is the
  conversational path.)

**Run** (observe a built program; deterministic, no model):
- `/do run_app <proj>` (or **`run`** flow: `/ws switch <proj>` then `/flow run`, or
  `ratchet flow . run --ws <proj> ""`) - build the workspace's `main` and run it, capturing
  stdout/stderr/exit. A blocking server is stopped after a few seconds. The build + tests stay the oracle.

## Worked examples (built by this ratchet, with transcripts)

These were generated by the local model under this ratchet - the `.spec` files and prompts were the
only hand-written input; the model wrote every line of Go, and the oracles verified it. The
`transcripts/` folder records the exact prompts/specs, the generated code, and the oracle verdicts -
the fastest way to see what driving the ratchet looks like.

- **`workspaces/pulsehook/`** - a low-latency webhook server (HTTP handler -> non-blocking enqueue ->
  worker-pool drain; 202 in well under a millisecond). Composed from five specs, then a behavior test
  was added with `add_file`, and it was run with the `run` flow.
  - [`transcripts/pulsehook-build.md`](transcripts/pulsehook-build.md) - the five input specs, the build
    plan, the per-unit generated code, and `go build ./...` + `go test ./...` passing.
  - [`transcripts/pulsehook-run.md`](transcripts/pulsehook-run.md) - the operator-console session (slash
    commands) running it, plus the live `curl` latency numbers.
- **`workspaces/pulsehook2/`** - the same webhook server, REBUILT from hardened specs after a code
  review of the first one. Shows the review -> KB -> better-code loop: the review's four findings became
  KB entries (`production-http-server` pattern, `atomic-int64-alignment` pitfall) and the rebuild is
  born with `atomic.Int64`, `http.MaxBytesReader`, `http.Server` timeouts, and graceful shutdown.
  [`transcripts/pulsehook2-hardened.md`](transcripts/pulsehook2-hardened.md) - the before/after diffs.
- **`workspaces/greetmod/`** - a multi-PACKAGE module (a `greeter/` sub-package imported by the root
  `main.go`), proving the ratchet scaffolds a real directory layout from specs (a spec declares
  `package:`), not just a flat `package main`.
  [`transcripts/greetmod-multipackage.md`](transcripts/greetmod-multipackage.md).
- **`workspaces/ledger/`** - a goroutine-safe counter module (compose) later extended with `add_file`
  and `edit_file`. Its compose run (specs -> plan -> per-unit generation -> whole-module test) is
  Part B of [`transcripts/phase2-grounded-and-compose.md`](transcripts/phase2-grounded-and-compose.md).
- **Single-function generation** - 10 focused tasks run through `go`/`test`, captured with the full
  rendered prompts, model output, and oracle verdicts:
  [`transcripts/phase2-evidence-batch.md`](transcripts/phase2-evidence-batch.md) (ungrounded) and
  [`transcripts/phase2-grounded-and-compose.md`](transcripts/phase2-grounded-and-compose.md) (grounded
  via plan-routing - shows how KB grounding changed the outcomes).

## The oracles (tools)
- `go_build` (`tools/go_check.sh`) - reads a Go file on stdin and runs `go build` as a library package
  (`package solution`, no `func main`). Exit 0 = type-checks.
- `go_test` (`tools/go_test.sh`) - reads the two marker-separated files on stdin, normalizes with
  `gofmt`, then runs `go vet` and `go test`. Exit 0 iff vet is clean and the tests pass. Requires a
  `*_test.go` (a package with no tests would be a silent false pass).
- `doctor_go` (`tools/doctor_go.sh`) - the toolbelt probe `ratchet doctor` runs; `/do doctor_go` prints
  the full report (go/gofmt/vet/test + optional goimports/staticcheck/golangci-lint).

All oracles are POSIX `.sh` (Linux/WSL/macOS); add a `.ps1` sibling per tool later for native Windows.

## Layout
- `flows/` - the action chains: `go`, `test`, `add_file`, `edit_file`, `compose`, `add_unit`, `run`
  (`flows/manifest.json` indexes them).
- `tools/` - oracles (`go_build`, `go_test`, `stage_check`, `module_check`), lifecycle/composition
  helpers (`new_module`, `read_file`, `read_module`, `module_api`, `stage_build`, `register_file`,
  `log_edit`, `plan_units`, `read_specs`), the `doctor_go` probe, and offline KB ingest
  (`kb_ingest_godoc`, `kb_ingest_patterns`) + their manifest.
- `kb/` - the routed libraries: `stdlib` (176 pkgs), `patterns` (23 GoF), `idioms` (see `kb/README.md`;
  rebuild an index with `ratchet index kb/<lib>`).
- `workspaces/` - persistent modules built by the lifecycle/compose flows.
- `ROADMAP.md` - the build-out plan and phase status; `transcripts/` - real end-to-end runs.
