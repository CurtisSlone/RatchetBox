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
ratchet doctor <path-to>/go                       # preflight: go on PATH + models pulled
ratchet validate-flow <path-to>/go                # lint the flow (model-free)
ratchet flow <path-to>/go go "a function that reverses a UTF-8 string"   # one-shot
ratchet <path-to>/go                              # operator console; then: /flow go <task>
```

## The oracle
`tools/go_check.sh` reads a complete Go file on stdin, writes it into a throwaway module, and runs
`go build` as a library package (`package solution`, no `func main` - the analog of C#'s
`csc -target:library`). Exit 0 = type-checks. To run natively on Windows too, add a `go_check.ps1`
and point the tool's `script` (or add a second tool entry) at it.

## Layout
- `flows/go/` - the action chain: generate -> build -> done | fix -> rebuild -> done | fail.
- `tools/` - `go_build` (the oracle) + its manifest.
- `kb/` - Go idioms the generate step grounds on (rebuild the index with `ratchet index`).
