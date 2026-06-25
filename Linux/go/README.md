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

## The two flows
- **`go`** - generate -> `go build` -> repair once. Verifies the code *type-checks* (library package).
- **`test`** - generate an implementation AND a test -> `go vet` + `go test` -> repair once. Verifies
  *behavior*, not just that it compiles. This is what `go build` alone cannot do; it is the headline of
  this ratchet. The generate step emits two marker-separated files (`// === solution.go ===` and
  `// === solution_test.go ===`) that the oracle splits and runs.

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
- `flows/go/`, `flows/test/` - the action chains (`flows/manifest.json` indexes them).
- `tools/` - `go_build`, `go_test`, `doctor_go` + their manifest.
- `kb/` - Go idioms the generate step grounds on (rebuild the index with `ratchet index`).
- `ROADMAP.md` - the build-out plan and phase status.
