#!/usr/bin/env bash
# read_file.sh - print one workspace file verbatim, to ground a surgical edit on its exact current
# contents. Args: proj path. Exits 1 if the file does not exist.
set -u
proj="${1:?usage: read_file <proj> <path>}"
path="${2:?missing path}"
root="workspaces/$proj"; [ -d "$root" ] || root="$proj"
f="$root/$path"
[ -f "$f" ] || { echo "no such file: $path (in workspace $proj)"; exit 1; }
cat "$f"
exit 0
