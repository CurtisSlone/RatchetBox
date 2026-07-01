#!/usr/bin/env bash
# evolve_ctx.sh - assemble the generation context for one evolve (layer) pass. For every unit the diff
# marks ADDED or CHANGED between spec/<from> and spec/<to>, emit a block with: its target .go path, its
# NEW spec (from spec/<to>), and its CURRENT file contents (so the model EDITS rather than reinvents -
# and so it never re-declares types that live in other, unchanged files). A change layer is a
# cross-cutting edit: all changed units are regenerated TOGETHER and staged atomically (stage_files),
# which is why this context lists them as one set, not one-at-a-time. Args: proj, layers ("<from> <to>").
set -u
proj="${1:?usage: evolve_ctx <proj> \"<from> <to>\"}"
layers="${2:?usage: evolve_ctx <proj> \"<from> <to>\"}"
root="workspaces/$proj"; [ -d "$root" ] || root="$proj"
[ -d "$root" ] || { echo "no such workspace: $proj"; exit 1; }
set -- $layers; frm="${1:-}"; to="${2:-}"
[ -n "$frm" ] && [ -n "$to" ] || { echo "layers must be \"<from> <to>\""; exit 1; }

# spec_diff yields the changed worklist: "<path> <specfile>" per unit, in dependency order.
work="$(bash tools/spec_diff.sh "$proj" "$layers")"
[ -n "$work" ] || { echo "no changed units between $frm and $to"; exit 1; }

echo "===== LAYER $frm -> $to: regenerate EXACTLY these files together ====="
echo "$work" | awk '{print "  - " $1 "   (from spec " $2 ")"}'
echo
echo "$work" | while read -r path spec; do
  [ -n "$path" ] || continue
  echo "########## UNIT: $path ##########"
  if [ -f "$root/$path" ]; then echo "(status: CHANGED - edit the current file below to satisfy the new spec)"; else echo "(status: NEW - create this file)"; fi
  echo
  echo "----- NEW spec ($spec) -----"
  cat "$root/spec/$to/$spec" 2>/dev/null || echo "(spec not found)"
  echo
  if [ -f "$root/$path" ]; then
    echo "----- CURRENT $path (edit this) -----"
    cat "$root/$path"
  else
    echo "----- CURRENT $path: none (new file) -----"
  fi
  echo
done

# Units DROPPED in this layer (their files are deleted by the prune step). Tell the model so it stops
# referencing them - a lingering reference would fail the whole-module verify.
removed="$(bash tools/spec_diff.sh "$proj" "$layers" removed 2>/dev/null)"
if [ -n "$removed" ]; then
  echo "########## REMOVED in this layer (do NOT reference these; their files are deleted) ##########"
  echo "$removed" | awk '{print "  - " $1 "   (was spec " $2 ")"}'
  echo
fi
exit 0
