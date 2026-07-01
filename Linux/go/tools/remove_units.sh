#!/usr/bin/env bash
# remove_units.sh - the evolve PRUNE step: delete the .go files for units DROPPED between spec/<from> and
# spec/<to> (present in <from>, absent in <to>), so a layer can shrink the system, not only grow it. The
# changed units are regenerated afterwards to stop referencing the removed ones; the transactional
# stage_files verify then proves nothing dangling remains (a still-referenced removed type fails the whole
# module build and rolls back). Rollback of the deletion itself is covered by the engine's per-run
# workspace snapshot (/rollback). Never fails the pass: exits 0 whether or not anything was removed.
# Args: proj, layers ("<from> <to>").
set -u
proj="${1:?usage: remove_units <proj> \"<from> <to>\"}"
layers="${2:?usage: remove_units <proj> \"<from> <to>\"}"
root="workspaces/$proj"; [ -d "$root" ] || root="$proj"
[ -d "$root" ] || { echo "no such workspace: $proj"; exit 0; }

removed="$(bash tools/spec_diff.sh "$proj" "$layers" removed 2>/dev/null)"
if [ -z "$removed" ]; then
  echo "prune: no units removed between $layers"
  exit 0
fi
n=0
echo "$removed" | while read -r path spec; do
  [ -n "$path" ] || continue
  case "$path" in /*|*..*) continue;; esac      # never escape the workspace
  if [ -f "$root/$path" ]; then
    rm -f "$root/$path"; echo "prune: removed $path (dropped spec $spec)"; n=$((n+1))
  else
    echo "prune: $path already absent (dropped spec $spec)"
  fi
done
exit 0
