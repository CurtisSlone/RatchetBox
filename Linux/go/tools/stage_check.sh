#!/usr/bin/env bash
# stage_check.sh - the project-lifecycle oracle (add_file / edit_file): write a generated Go file into a
# workspace, then type-check AND test the WHOLE module with `go vet ./...` + `go test ./...`. Exit 0 iff
# both pass. Uses `go vet` (not `go build`) so it works whether or not the module has a `func main` yet
# (a bare `go build` would link-fail without one); behavior is covered by `go test`. Args: proj path;
# the Go source on stdin.
set -u
proj="${1:?usage: stage_check <proj> <path> (code on stdin)}"
path="${2:?missing path}"
root="workspaces/$proj"; [ -d "$root" ] || root="$proj"
[ -d "$root" ] || { echo "no such workspace: $proj"; exit 1; }
case "$path" in /*|*..*) echo "unsafe path: $path"; exit 1;; esac

code="$(cat)"
code="$(printf '%s\n' "$code" | sed '/^[[:space:]]*```/d')"
mkdir -p "$(dirname "$root/$path")"
printf '%s\n' "$code" > "$root/$path"
gofmt -w "$root/$path" 2>/dev/null

v="$(cd "$root" && GOFLAGS=-mod=mod go vet ./... 2>&1)"; vs=$?
if [ "$vs" -ne 0 ]; then echo "VET FAILED after staging $path:"; printf '%s\n' "$v"; exit 1; fi

t="$(cd "$root" && GOFLAGS=-mod=mod go test ./... 2>&1)"; ts=$?
if [ "$ts" -ne 0 ]; then echo "TEST FAILED after staging $path:"; printf '%s\n' "$t"; exit 1; fi

echo "OK: staged $path; module vets and tests pass"
printf '%s\n' "$t"
exit 0
