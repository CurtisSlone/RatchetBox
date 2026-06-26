#!/usr/bin/env bash
# read_files.sh - print several workspace files, each under a "=== path ===" marker, as the current state
# for a coordinated multi-file edit (the coedit flow). Args: proj, a comma-separated list of file paths
# (relative to the workspace). A path that does not exist yet is marked as a NEW FILE.
set -u
proj="${1:?usage: read_files <proj> <comma-separated-files>}"
files="${2:?missing comma-separated file list}"
root="$proj"; [ -d "$root" ] || root="workspaces/$proj"
[ -d "$root" ] || { echo "no such workspace: $proj"; exit 1; }
IFS=',' read -ra arr <<< "$files"
n=0
for raw in "${arr[@]}"; do
  f="$(printf '%s' "$raw" | xargs)"   # trim whitespace
  [ -z "$f" ] && continue
  printf '=== %s ===\n' "$f"
  if [ -f "$root/$f" ]; then cat "$root/$f"; else printf '(NEW FILE - does not exist yet)\n'; fi
  printf '\n\n'
  n=$((n + 1))
done
[ "$n" -eq 0 ] && { echo "no files given"; exit 1; }
exit 0
