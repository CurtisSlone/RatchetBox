#!/usr/bin/env bash
# add_dep.sh - add a third-party Go module to a workspace AND ingest its API docs into the `deps` KB,
# so the lifecycle/compose flows ground generation on the module's real surface (the same offline
# `go doc` trick the stdlib library uses). Runs `go get` + `go mod tidy` in the workspace, then writes
# `go doc -all <pkg>` to kb/deps/<slug>.md and re-indexes. Needs network for `go get`.
# Run: /do add_dep <workspace> <module-path[@version]>
set -u
proj="${1:?usage: add_dep <workspace> <module-path[@version]>}"
mod="${2:?missing module path (e.g. github.com/google/uuid or ...@v1.6.0)}"
root="workspaces/$proj"; [ -d "$root" ] || root="$proj"
[ -d "$root" ] || { echo "no such workspace: $proj"; exit 1; }
pkg="${mod%@*}"   # import path without any @version, for go doc

if ! ( cd "$root" && GOFLAGS=-mod=mod go get "$mod" 2>&1 ); then
  echo "go get failed for $mod (network? bad module path?)"; exit 1
fi
# NOTE: do NOT `go mod tidy` here - the dep is not imported by any file yet, so tidy would strip it
# straight back out of go.mod (and then go doc could not find it). It stays as a require until a unit
# imports it; the operator can tidy later. An unused require is not a build/vet error.

doc="$(cd "$root" && go doc -all "$pkg" 2>/dev/null)"
if [ -z "$doc" ]; then
  echo "added $mod to $proj, but go doc found no package at '$pkg' (a multi-package module? pass the exact import path)."
  exit 0
fi

mkdir -p kb/deps
slug="$(printf '%s' "$pkg" | tr -c 'A-Za-z0-9' '_' | sed 's/_*$//')"
syn="$(cd "$root" && go doc "$pkg" 2>/dev/null | awk 'NR>=3 && NF{print; exit}')"
{
  printf '# %s (third-party Go module)\n\n' "$pkg"
  [ -n "$syn" ] && printf '%s\n\n' "$syn"
  printf 'Import path: %s   Added to workspace: %s\n\n' "$pkg" "$proj"
  printf '%s\n' "$doc"
} > "kb/deps/$slug.md"
echo "ingested $pkg -> kb/deps/$slug.md"

if command -v ratchet >/dev/null 2>&1; then
  ratchet index kb/deps >/dev/null 2>&1 && echo "reindexed kb/deps"
else
  echo "(ratchet not on PATH; run: ratchet index kb/deps)"
fi
echo "OK: $mod added to workspace $proj and grounded in the deps KB"
exit 0
