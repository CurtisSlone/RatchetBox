# pulsehook - run transcript (operator console, slash commands)

Running the composed `pulsehook` low-latency webhook server through the ratchet's run capability
(phase 4). Captured from the operator console (`ratchet <this-dir>`) so it shows the actual slash
commands. The app itself was built earlier from specs only - see `pulsehook-build.md`.

The run capability is deterministic (running needs no model):
- `run_app` tool - build the workspace's `main` and run it, capturing stdout/stderr/exit; a blocking
  server is stopped after a timeout and reported as still-running. Invoke with `/do run_app <proj>`.
- `run` flow - a thin one-action wrapper so it is operator-friendly with a workspace: `/flow run`
  after `/ws switch <proj>` (console), or `ratchet flow . run --ws <proj> ""` (CLI).
- Running OBSERVES behavior; the build + tests remain the oracle (compose's `module_check`).

## Console session (verbatim)

```text
$ ratchet ../RatchetBox/Linux/go        # open the operator console

ratchet operator console - 'go'
  dispatch seat: phi3:mini   generate seat: qwen3-coder:latest   ollama: http://172.18.160.1:11434

ratchet >  /flows

Authored flows (/route can match these, or run with /flow <name>):
  add_file - Add a new Go file to a workspace from a request; ... verified with go vet + go test ...
  add_unit - Generate one composed Go unit against the accumulated module ...
  compose  - Compose a Go module from a dir of specs: plan ... then build and test the whole thing.
  edit_file - Apply a change to an existing Go file in a workspace ...
  go       - Generate Go code for a focused task, verify it with `go build`, repair once ...
  run      - Build and run a workspace's main, capturing stdout/stderr/exit ...
  test     - Generate a Go implementation plus a test ... verify with `go vet` + `go test` ...

ratchet >  /do run_app pulsehook

built pulsehook; running for up to 3s...
--- program output ---
2026/06/26 00:30:07 Listening on :8080
----------------------
still running after 3s; stopped (normal for a server/long-running process).

ratchet >  /ws switch pulsehook

active workspace: pulsehook

ratchet >  /flow run
  - step 1: run.exec (action)
  - step 2: run.done (exit)

built .../workspaces/pulsehook; running for up to 5s...
--- program output ---
2026/06/26 00:30:10 Listening on :8080
----------------------
still running after 5s; stopped (normal for a server/long-running process).

ratchet >  /exit
bye
```

Note: in the console the active workspace is set with `/ws switch <proj>`; `--ws` is the CLI form
(`ratchet flow . run --ws pulsehook ""`), which produces the same `run.exec -> run.done` result.

## Serving under load (built binary + real curl)

The server logs `Listening on :8080` and accepts webhooks at sub-millisecond latency, because the
handler enqueues to a buffered channel and returns `202` immediately while a 4-worker pool drains the
queue asynchronously:

```text
$ curl -X POST localhost:8080/webhook -d event-1   ->  HTTP 202  in 0.000878s
$ curl -X POST localhost:8080/webhook -d event-2   ->  HTTP 202  in 0.000437s
$ curl -X POST localhost:8080/webhook -d event-3   ->  HTTP 202  in 0.000405s
$ curl     localhost:8080/webhook                  ->  HTTP 405
time-to-first-byte on a POST: 0.000355s
```

And the behavior is pinned by the generated test (`go test ./...`):

```text
=== RUN   TestWebhookAcceptsAndProcesses
--- PASS: TestWebhookAcceptsAndProcesses (0.00s)   # 202 + Processed()==1 (async work done)
=== RUN   TestWebhookRejectsGet
--- PASS: TestWebhookRejectsGet (0.00s)            # GET -> 405
PASS
ok  	pulsehook	0.002s
```
