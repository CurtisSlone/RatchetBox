#!/usr/bin/env bash
# log_edit.sh - append an 'edited' changelog line to a workspace's PROJECT.md (best-effort, always
# exits 0 - bookkeeping must not fail the chain). Args: proj path summary.
set -u
proj="${1:?usage: log_edit <proj> <path> <summary>}"
path="${2:?missing path}"
summary="${3:-edited}"
root="workspaces/$proj"; [ -d "$root" ] || root="$proj"
[ -d "$root" ] || { echo "(no such workspace: $proj)"; exit 0; }
printf -- '- edited `%s`: %s\n' "$path" "$summary" >> "$root/PROJECT.md"
echo "logged edit to $path"
exit 0
