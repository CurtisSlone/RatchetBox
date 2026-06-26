# ROADMAP - building out the go ratchet

Working plan for growing `Linux/go` from a single generate-build-repair MVP into the full
cross-platform reference ratchet. Plain and grounded; verify each step, don't assert. Domain logic
only - no engine changes (those live in the Ratchet repo).

## Where it stands today

A clean MVP that proves the loop end to end:

- one flow: `flows/go/` - generate -> `go build` -> done | fix -> rebuild -> done | fail (bounded,
  unrolled repair, no cycle).
- one oracle tool: `tools/go_check.sh` - reads a Go file on stdin, writes it into a throwaway module,
  type-checks it as a library package (`package solution`, no `func main`). Exit 0 = type-checks.
- one KB file: `kb/go-idioms.md` (single manifest entry; under-routed).
- `ratchet.json` seats: generate `qwen3-coder`, dispatch `phi3:mini`, embed `nomic-embed-text`.

The mature targets to grow toward are the Windows ratchets: `Windows/dotnet4-x` (7 flows, 16 tools,
77-file KB across 5 libraries, full project lifecycle, spec-to-code `compose`) and `Windows/cpp`
(plan-routed retrieval over 7 KBs, lifecycle, native run, vcpkg). That is the parity gap.

This is a Linux-first ratchet: every tool is a POSIX shell script (`.sh`) dispatched to bash. The
`.ps1` siblings are a later, optional add for native Windows; they are not the focus. The bash oracle
already runs on Linux, WSL, and macOS, which is the cross-platform claim that matters.

## 0. Toolbelt check (do this first, like cpp's doctor)

`Windows/cpp` ships a domain probe (`tools/doctor_cl.ps1`) wired into `ratchet.json` `requirements`
as `{ "name": "MSVC cl", "tool": "doctor_cl" }`, so `ratchet doctor` reports whether the real
toolchain is reachable before any flow runs. The go ratchet needs the same: a `tools/doctor_go.sh`
that probes `go` + the optional linters and reports what is available, so flows can gate on it
(e.g. skip the staticcheck oracle when staticcheck is absent).

Scan of this box (Linux, go1.26.4, 2026-06-25) - what the toolbelt actually has:

| Tool | Status | Notes |
|---|---|---|
| `go` | present | go1.26.4 linux/amd64 at /usr/local/go/bin |
| `gofmt` | present | ships with go - `go_fmt` oracle works now |
| `go vet` | present | subcommand - `go_vet` oracle works now |
| `go test` | present | subcommand - `go_test` oracle works now |
| `go build` | present | already the `go_build` oracle |
| bash / sh | present | tool dispatch target |
| ollama | reachable | OLLAMA_URL points at the Windows-host gateway; models pulled |
| `goimports` | MISSING | `go install golang.org/x/tools/cmd/goimports@latest` (optional, import-fixing) |
| `staticcheck` | MISSING | `go install honnef.co/go/tools/cmd/staticcheck@latest` (optional stricter oracle) |
| `golangci-lint` | MISSING | optional meta-linter |
| `dlv`, `gopls` | MISSING | dev tooling, not needed by the ratchet |

Takeaway: the entire behavior-ratchet stack (`gofmt` -> `go vet` -> `go build` -> `go test`) is
available today with **zero extra installs**. The optional tools above are nice-to-have, gated by
`doctor_go.sh`, not blockers. Build phase 1 against the stdlib toolchain; treat staticcheck as opt-in.

## The thesis: ratchet on behavior, not just "it compiles"

The C# and C++ ratchets stop their oracle at "it compiles." Go's toolchain is richer and pure-Go
(no cgo, fits the engine's `CGO_ENABLED=0` rule), so the go ratchet can go further:

`gofmt` (canonicalize) -> `go vet` (suspicious-but-buildable) -> `go build` (type-check) ->
`go test` (behavior) -> `go run` (run-verify).

`go test` is the headline: generate solution + tests and ratchet until the tests pass. That is a
capability the existing ratchets cannot cleanly express. Build the roadmap around it.

## 1. Tools (oracles) - layer them

Every tool is a POSIX `.sh` script (Linux-first), dispatched to bash by extension. `.ps1` siblings
are deferred (see section 4). All pure-Go toolchain, no external deps beyond `go` (staticcheck
optional, ships as a static binary).

| Tool | Verdict (exit 0 = pass) | Available now | Why |
|---|---|---|---|
| `doctor_go` | probe: `go` + linters present | yes | Toolbelt check (section 0); gates optional oracles |
| `go_fmt` | `gofmt -l` empty / emits canonical source | yes | Cheap pre-pass; normalizes model output and kills the fence problem at the source instead of `sed`-stripping |
| `go_vet` | `go vet ./...` clean | yes | Catches printf-verb, shadowing, unreachable, struct-tag bugs `go build` accepts |
| `go_build` | `go build ./...` (have it) | yes | Type-check, library package |
| `go_test` | `go test ./...` green | yes | Behavior - the real ratchet |
| `go_run` | `go run` a `package main`, capture stdout | yes | Run-verify an executable (analog of cpp `run_app`) |
| `go_mod` | `go mod tidy` / add a dep | yes | Deliberate opt-in to third-party modules |
| `go_staticcheck` | `staticcheck ./...` clean | opt-in | Stricter lint; gate on `doctor_go` since the binary is not installed here yet |

Note on `go_check.sh` today: it `sed`-strips markdown fences because the model emits them. A `go_fmt`
pre-pass is the cleaner fix - canonical source or a non-zero verdict, no regex.

## 2. Flows - climb the lifecycle ladder

Ordered by leverage:

1. **`go_test`** - generate solution + a test (or accept a provided one) -> `go test` -> repair.
   The capability the other ratchets lack. Start here.
2. **`gofix`** - take broken/ugly Go -> vet + build + repair -> return canonical. Pure-repair flow.
3. **Project lifecycle** (mirror dotnet4-x): `new_project` (scaffold a real module with `go.mod`) ->
   `add_file` / `edit_file`, each rebuilding the whole module via `go build ./...` (+ `go test ./...`).
   Go's module model is cleaner here than the C# `response.rsp` build.
4. **`compose`** - spec-dir -> plan build order -> `foreach` generate each file against the
   accumulated package -> `go build ./...` + `go test ./...`. Go packages accumulate more naturally
   than C# files do.
5. **structured/`generate` with `output_schema`** - JSON-shaped tasks (struct + tags), oracle =
   build + a marshalling round-trip test.

Reuse the existing unrolled-repair shape (generate -> verify -> done | fix -> reverify -> done | fail)
so chains lint clean and terminate; no cycles.

## 3. Knowledge bases - cpp-style ingested libraries

PROGRESS (2026-06-25): adopted the cpp multi-library model (one `kb/<lib>/` per source, each
registered separately and indexed to its own manifest), not the single recursive KB. Shipped:

- `tools/kb_ingest_godoc.sh` - deterministic, offline, no-model ingest of the Go standard library via
  `go doc -all`, one markdown per package. The Go analog of cpp's `cppref` ingest; `go doc` generates
  it locally so there is nothing to scrape and it matches the installed toolchain.
- `kb/stdlib/` - 176 package reference files (2.7 MB), registered as the `stdlib` KB, indexed (176
  manifest entries). Retrieval verified: `/search stdlib <heap query>` routes to `container/heap` and
  grounds a correct `heap.Interface` priority queue - the exact API run 10 got wrong.
- `kb/idioms/` - the former single `go-idioms.md`, now its own library (the `kb` default), so the
  default KB no longer sits at the kb root and does not absorb `stdlib`.
- `kb/README.md` - the library index (sources, build commands, how to add one).
- `tools/kb_ingest_patterns.sh` + `kb/patterns/` - the 23 GoF design patterns in compiling Go (impl +
  runnable example per pattern), ingested from a Category/Pattern repo (the Go analog of cpp's
  `ingest-md.ps1 -Mode patterns`). Registered as the `patterns` KB, indexed (23 entries; keywords pulled
  from the fenced Go code). Retrieval verified: a "notify many subscribers" query routes to
  `behavioral_observer.md`.

WIRED (2026-06-25): both flows (`go`, `test`) now open with a `plan` node (cpp-style plan-routing) -
the model emits one search query per library (`stdlib_q`/`patterns_q`/`idioms_q`), an empty query is
skipped, and `generate` routes each non-empty query into that library's `search`. Verified live: the
run-10 heap task now plans `stdlib_q="container/heap"` (patterns skipped), the `container/heap` doc
reaches the prompt, and the task PASSES (it failed even after repair pre-grounding); the `go` flow
plans `stdlib_q="encoding/json Marshal"` for a JSON task. On a malformed plan, the node routes to
`generate` and proceeds ungrounded (never aborts).

REMAINING libraries (the routed taxonomy below). Author `patterns/` (done - 23 ingested) is in; add
`pitfalls/` from the run evidence next (cheap, high-signal); `guidelines/` + `spec/` are split-ingests
of external text (Effective Go, Code Review Comments, Go Proverbs, the language spec) - confirm
sourcing before fetch. Re-running the phase-2 batch now (with grounding on) would quantify the lift.

### The routed taxonomy

`kb/go-idioms.md` is a single entry today, so routing has almost nothing to discriminate on. Mirror
the dotnet `reference/patterns/recipes/scaffold/snippets` taxonomy, Go-flavored, one topic per file
with a sharp one-line `summary` (routing sees only the summary), then `ratchet index <go>/kb`:

- **`reference/`** - stdlib facts the model gets wrong: error wrapping (`%w`), `context` rules,
  slice/map aliasing, goroutine/channel semantics, generics constraints, `go.mod` versioning.
- **`patterns/`** - complete units: worker pool, `errgroup`, table-driven test, `io.Reader` pipeline,
  functional options, HTTP handler.
- **`recipes/`** - playbooks: CLI with flag parsing + tests, JSON API client, fan-out/fan-in.
- **`pitfalls/`** - builds-but-wrong: loop-var capture, nil-map write, unbuffered-channel deadlock,
  `defer` in a loop. Pairs with the `go_vet` / `go_test` oracles.
- **`scaffold/`** - `go.mod` + package-layout templates for `new_project` / `compose`.

Migrate the current `go-idioms.md` content into `reference/` and `patterns/` rather than deleting it.

## 4. Cross-platform finish (deferred - Linux-first)

We author and run on Linux, so every tool is a `.sh` script and that is the supported path. The
bash oracle already covers Linux, WSL, and macOS. A later, optional pass can add a `.ps1` sibling
per oracle for native Windows; it is not in scope for the phases below and should not slow them down.

## Third-party modules (DONE 2026-06-25)

`tools/add_dep.sh` (`/do add_dep <ws> <module[@ver]>`) `go get`s a module into a workspace and ingests
its `go doc -all` into the `deps` KB (kb/deps/<slug>.md), then re-indexes - the user's insight that a
new dependency should bring its docs as grounding, using the same offline `go doc` trick as `stdlib`.
Registered `deps` as a 5th KB; the workspace-generate nodes (`add_file`, `edit_file`, `add_unit`) now
search it keyed on the request/spec. Gotcha handled: NO `go mod tidy` in add_dep (it would strip the
not-yet-imported require straight back out). Verified live: added `github.com/google/uuid` to a
workspace, then `add_file` generated `uuid.New().String()` grounded on the deps KB and passed `go vet`
+ `go test` first try (vet resolves the third-party import; `stage_check` uses vet not build, so no
`func main` is required).

## KB depth (DONE 2026-06-25/26)

Three new libraries, wired into the `go`/`test` plan-routing (5 query fields now: stdlib, patterns,
guidelines, pitfalls, idioms; `patterns_q` broadened to cover named algorithms):
- `guidelines/` (12) - idiomatic-style entries authored from the canonical sources fetched via WebFetch
  (Effective Go, Code Review Comments, Go Proverbs), real code in fences; attribution in
  `kb/guidelines/ATTRIBUTION.md` (CC BY 3.0 / BSD).
- `pitfalls/` (5) - builds-but-wrong traps authored from run evidence: unused imports, nil-map write,
  loop-var capture, channel deadlock, slice aliasing.
- `patterns/` += authored `algo_two_heap_streaming_median.md` (verified correct via the go_test oracle).

Result: the two-heap median (run 10), which FAILED even after repair in both earlier batches, now passes
FIRST TRY. Confirmed it was the grounding: the plan routed `patterns_q="two-heap streaming median"` +
`stdlib_q="container/heap"`, both reached the prompt. This is the evidence-driven KB loop closing -
an observed oracle failure became a KB entry that removed the failure.

## Phased sequence

1. DONE (2026-06-25) - `doctor_go.sh` toolbelt probe (wired into `ratchet.json` requirements) +
   the behavior flow + oracle. Shipped: `tools/doctor_go.sh`, `tools/go_test.sh` (gofmt-normalize ->
   `go vet` -> `go test`, requires a `*_test.go` so an empty package is not a silent false pass),
   `flows/test/` (generate impl+test -> check -> done | fix -> recheck -> fail), `flows/manifest.json`.
   Verified: `validate-flow` clean (go + test), `doctor` green incl. the toolbelt probe, the oracle
   unit-tested on 4 payloads (pass / behaviorally-wrong / no-test / fenced+vet-bug all correct), and a
   live end-to-end run (FizzBuzz) generated a real table-driven test and passed first try
   (run `20260625-170206-296`). `gofmt`/`go vet`/`go test` are folded into the one `go_test` oracle;
   split them into standalone oracle nodes later only if the run evidence shows it helps attribution.
2. Split + grow the KB into the routed taxonomy.
3. DONE (2026-06-25) - project lifecycle on a persistent workspace. Shipped `tools/stage_check.sh`
   (vet + test whole-module oracle; uses `go vet` not `go build` so it works with or without a
   `func main`), `tools/read_file.sh`, `tools/log_edit.sh`, and the `add_file` + `edit_file` flows
   (request-driven, grounded on `module_api` + plan-routing-lite, repair once). `new_module` (from
   phase 4) scaffolds the workspace. Verified live on the `ledger` module: `add_file` added `double.go`
   (its first attempt had an unused `strconv` import, which the oracle caught -> repair removed it);
   `edit_file` then added `math.MaxInt` saturation to `Double`, grounded on the stdlib search, preserving
   the rest of the file, first try. Daily-use loop now works: new_module -> add_file/edit_file -> compose.
4. DONE (2026-06-25) - `compose` from specs. Shipped 8 Go-native tools (`new_module`, `read_specs`,
   `module_api`, `stage_build`, `read_module`, `plan_units`, `register_file`, `module_check`) + the
   `add_unit` and `compose` flows. Go-native simplifications over the C# version: every unit shares
   `package main` at the module root (no cross-imports - the analog of "everything in namespace App");
   `module_api` is grep over `^type|^func` (no parser); the per-unit oracle is `go vet ./...`
   (type-checks without linking, so intermediate units pass before the entry's `func main` is staged);
   the final `module_check` is `go build ./...` + `go test ./...`, so composition is BEHAVIOR-verified
   (stronger than the C# compose, which only gates compilation). Plan-routing wired into `add_unit`
   too (stdlib + patterns keyed on the spec). Verified end to end: the 3-spec `ledger` module composed,
   built, and its concurrent `TestCounter` (50 goroutines x 100 = 5000) passed under `go test`. Full
   transcript: `transcripts/phase2-grounded-and-compose.md`.
5. `.ps1` oracle siblings for native Windows parity.

## Run capability (DONE 2026-06-26)

`tools/run_app.sh` (`/do run_app <proj>`) builds a workspace's `main` and runs it, capturing
stdout/stderr/exit; a blocking server is stopped after a timeout (default 3s) and reported as
still-running (exit 0). A thin `run` flow wraps it for `--ws`/`/ws switch` ergonomics
(`/flow run` after `/ws switch <proj>`, or `ratchet flow . run --ws <proj> ""`). Running needs no
model, so `run` is a single deterministic action node, not a generate/repair chain - the build + tests
remain the oracle. Verified by running the composed `pulsehook` server (logs `Listening on :8080`,
stopped cleanly); transcript with the live console slash-command session in
`transcripts/pulsehook-run.md`. NOTE: `--ws` is the CLI flag; in the console set the workspace with
`/ws switch <proj>` first (the `--ws` token is not parsed inside a console `/flow`).

## Open decisions (settle before building each phase)

- **Test sourcing for `go_test`**: model generates the test too, or operator supplies it, or both
  (generate when absent). Self-generated tests can be gamed by weak models; consider an oracle that
  rejects trivial/empty tests.
- **Library vs main package**: `go_build` uses `package solution` (library). `go_run` and lifecycle
  need `package main`. Decide per flow; keep the library default for one-shot snippets.
- **Module identity in lifecycle**: one module per workspace (`go.mod` at `workspaces/<name>`), build
  with `go build ./...`. Confirm the engine's workspace/`--ws` plumbing matches dotnet4-x's.
- **Dispatch seat**: `phi3:mini` for dispatch today; confirm it is strong enough to route over a
  5-library KB, or move dispatch to `qwen3-coder`.
- **staticcheck**: include as an optional stricter oracle, or keep to the stdlib `go vet` only.

## Verify after each change

```
ratchet validate-flow <go>     # lint every chain (node kinds, fields, unknown tools, reachability)
ratchet doctor <go>            # preflight declared tools (go on PATH, models pulled)
ratchet index <go>/kb          # rebuild the KB routing index after editing kb/
ratchet flow <go> <chain> "<task>"   # end-to-end against Ollama
```

"Builds" is not "behaves" - once `go_test` lands, prove behavior with a real failing-then-passing run.
