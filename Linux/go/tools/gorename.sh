#!/usr/bin/env bash
# gorename.sh - TYPE-SAFE rename across a workspace module (golang.org/x/tools/cmd/gorename), then verify
# the module still builds. It renames every reference across files by type - not text substitution - so
# it cannot rename the wrong thing. Args: proj, target (a package-level name, or Type.Method), newname.
# The module path is read from go.mod and prepended automatically. Run: /do gorename <proj> <target> <new>.
set -u
proj="${1:?usage: gorename <proj> <target> <newname>}"
target="${2:?missing target (a name or Type.Method)}"
new="${3:?missing new name}"
root="$proj"; [ -d "$root" ] || root="workspaces/$proj"
[ -d "$root" ] || { echo "no such workspace: $proj"; exit 1; }
command -v gorename >/dev/null 2>&1 || { echo "gorename not installed: go install golang.org/x/tools/cmd/gorename@latest"; exit 1; }
cd "$root" || { echo "cannot enter $root"; exit 1; }
module="$(awk '/^module /{print $2; exit}' go.mod 2>/dev/null)"
[ -z "$module" ] && { echo "no module path in go.mod"; exit 1; }

out="$(GOFLAGS=-mod=mod gorename -from "\"$module\".$target" -to "$new" 2>&1)"; st=$?
printf '%s\n' "$out"
if [ "$st" -ne 0 ]; then echo "RENAME FAILED (nothing changed)"; exit 1; fi

b="$(GOFLAGS=-mod=mod go build ./... 2>&1)"
if [ $? -ne 0 ]; then echo "POST-RENAME BUILD FAILED:"; printf '%s\n' "$b"; exit 1; fi
echo "OK: renamed $target -> $new across the module; it still builds."
exit 0
