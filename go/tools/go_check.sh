#!/usr/bin/env bash
# go_check.sh - the Go oracle for the generate-verify-repair loop.
# Reads a complete Go source file from stdin, writes it into a throwaway module, and type-checks it
# with `go build` (as a library package - no func main required, the analog of csc -target:library).
# Prints "OK: ..." and exits 0 iff it builds; otherwise prints the build diagnostics and exits 1.
# Does not run anything. Cross-platform (bash on Linux/WSL/macOS); add a .ps1 sibling for native Windows.
set -u

code="$(cat)"

# Models sometimes wrap output in Markdown code fences - or emit a dangling fence with no opener.
# A line that is only a fence marker (```), ```go, etc.) is never valid Go, so drop every such line.
code="$(printf '%s\n' "$code" | sed '/^[[:space:]]*```/d')"

if ! command -v go >/dev/null 2>&1; then
  echo "go not found on PATH (install Go / check PATH)"
  exit 1
fi

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
printf '%s\n' "$code" > "$tmp/snippet.go"

out="$(cd "$tmp" && go mod init snippet >/dev/null 2>&1; GOFLAGS=-mod=mod go build ./... 2>&1)"
status=$?

if [ "$status" -eq 0 ]; then
  echo "OK: builds with $(go version | awk '{print $3}')"
else
  printf '%s\n' "$out"
fi
exit "$status"
