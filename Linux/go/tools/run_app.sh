#!/usr/bin/env bash
# run_app.sh - build and RUN a workspace's main, capturing stdout/stderr/exit (the `go run` / cpp
# run_app analog). A program that exits on its own reports its real exit code; a server that blocks is
# stopped after the timeout and that is reported as a normal "still running" result. Running OBSERVES
# behavior - the build and tests remain the oracle. Args: proj [timeout_secs] (default 3).
# Run: /do run_app <proj>
set -u
proj="${1:?usage: run_app <proj> [timeout_secs]}"
secs="${2:-3}"
# Tolerate a missing/unsubstituted/non-numeric secs (e.g. an unfilled "{secs}" placeholder) -> default 3.
case "$secs" in ''|*[!0-9]*) secs=3;; esac
root="workspaces/$proj"; [ -d "$root" ] || root="$proj"
[ -d "$root" ] || { echo "no such workspace: $proj"; exit 1; }

tmp="$(mktemp -d)"; bin="$tmp/app"; trap 'rm -rf "$tmp"' EXIT

build="$(cd "$root" && GOFLAGS=-mod=mod go build -o "$bin" . 2>&1)"
if [ $? -ne 0 ]; then echo "BUILD FAILED:"; printf '%s\n' "$build"; exit 1; fi
echo "built $proj; running for up to ${secs}s..."

out="$tmp/out"
timeout "$secs" "$bin" >"$out" 2>&1 </dev/null
status=$?

echo "--- program output ---"
cat "$out"
echo "----------------------"
if [ "$status" -eq 124 ]; then
  echo "still running after ${secs}s; stopped (normal for a server/long-running process)."
  exit 0
fi
echo "exited with status $status."
exit "$status"
