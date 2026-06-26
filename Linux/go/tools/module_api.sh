#!/usr/bin/env bash
# module_api.sh - print the existing API surface of a composed module: every top-level type/func and
# method declaration, grouped by directory (package), each labeled with how to reach it. For a sub-
# package it shows the exact import path ("<module>/<dir>"); for the root it notes same-package access.
# This is the authoritative list of names a new unit must call VERBATIM (Go gives it for free). Arg:
# proj (workspace name). Always exits 0 (grounding, not an oracle).
set -u
proj="${1:?usage: module_api <proj>}"
root="workspaces/$proj"; [ -d "$root" ] || root="$proj"
[ -d "$root" ] || { echo "(no such workspace: $proj)"; exit 0; }
module="$(awk '/^module /{print $2; exit}' "$root/go.mod" 2>/dev/null)"

n=0
while IFS= read -r f; do
  [ -e "$f" ] || continue
  sigs="$(grep -nE '^(type |func )' "$f" 2>/dev/null)"
  [ -z "$sigs" ] && continue
  rel="${f#"$root"/}"
  dir="$(dirname "$rel")"
  if [ "$dir" = "." ]; then
    note="package main - same package, call these names directly (no import)"
  else
    note="package $(basename "$dir") - import \"$module/$dir\" then call as $(basename "$dir").Name (exported names only)"
  fi
  printf '## %s  [%s]\n%s\n\n' "$rel" "$note" "$sigs"
  n=$((n + 1))
done < <(find -L "$root" -name '*.go' -not -name '*_test.go' | sort)
[ "$n" -eq 0 ] && echo "(no Go units yet - this is the first unit)"
exit 0
