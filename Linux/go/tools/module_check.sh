#!/usr/bin/env bash
# module_check.sh - the whole-module FINAL oracle for compose: go build ./... then go test ./... over the
# entire module. Exit 0 iff both pass. This is where behavior (not just compilation) of the composed
# system is verified. Arg: proj (workspace name).
set -u
proj="${1:?usage: module_check <proj>}"
root="workspaces/$proj"; [ -d "$root" ] || root="$proj"
[ -d "$root" ] || { echo "no such workspace: $proj"; exit 1; }
cd "$root" || { echo "cannot enter $root"; exit 1; }

RACE=""
if [ "$(go env CGO_ENABLED 2>/dev/null)" != "0" ] && { command -v gcc >/dev/null 2>&1 || command -v clang >/dev/null 2>&1 || command -v cc >/dev/null 2>&1; }; then RACE="-race"; fi

b="$(GOFLAGS=-mod=mod go build ./... 2>&1)"; bs=$?
if [ "$bs" -ne 0 ]; then echo "MODULE BUILD FAILED:"; printf '%s\n' "$b"; exit 1; fi

t="$(GOFLAGS=-mod=mod go test $RACE ./... 2>&1)"; ts=$?
if [ "$ts" -ne 0 ]; then echo "MODULE TEST FAILED${RACE:+ (-race)}:"; printf '%s\n' "$t"; exit 1; fi

echo "OK: module builds and tests pass${RACE:+ (-race)} with $(go version | awk '{print $3}')"
printf '%s\n' "$t"
if command -v staticcheck >/dev/null 2>&1; then
  sc="$(staticcheck ./... 2>&1)"   # already in $root (cd'd above)
  [ -n "$sc" ] && { echo "STATICCHECK (advisory):"; printf '%s\n' "$sc"; }
fi
exit 0
