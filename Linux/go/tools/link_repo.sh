#!/usr/bin/env bash
# link_repo.sh - work on an EXISTING Go module that lives outside workspaces/. Symlinks it in as
# workspaces/<name>, so every lifecycle flow (add_file, edit_file, harden, run, module_api, ...) operates
# on the real repo unchanged - the engine's --ws <name> resolves to the symlink. The module's own package
# layout is respected (the lifecycle derives a file's package from its path). The build/test/vet oracles
# protect the real code: a change that breaks the module is rejected before it sticks. Args: name,
# path-to-module (must contain go.mod). Run: /do link_repo <name> <path>.
set -u
name="${1:?usage: link_repo <name> <path-to-go-module>}"
path="${2:?missing path to the module}"
[ -d "$path" ] || { echo "no such directory: $path"; exit 1; }
[ -f "$path/go.mod" ] || { echo "not a Go module (no go.mod): $path"; exit 1; }
abs="$(cd "$path" && pwd)"
mkdir -p workspaces
target="workspaces/$name"
if [ -e "$target" ] && [ ! -L "$target" ]; then
  echo "workspaces/$name already exists and is NOT a symlink - refusing to overwrite a real workspace"; exit 1
fi
ln -sfn "$abs" "$target"
echo "linked workspaces/$name -> $abs"
echo "module: $(awk '/^module /{print $2; exit}' "$abs/go.mod")"
echo "now drive it: ratchet flow . add_file --ws $name \"<path> <request>\"   |   /flow harden --ws $name"
echo "(unlink later with: rm workspaces/$name  - removes only the symlink, not your repo)"
exit 0
