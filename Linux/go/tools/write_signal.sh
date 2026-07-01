#!/usr/bin/env bash
# write_signal.sh - the evolve ESCALATION step. When a layer still fails after the bounded repair, the
# problem is usually residual-4 (a structural move the small model cannot make even when shown) - so
# instead of a bare abort, package everything a human or a stronger model needs to finish the ONE stuck
# node: the layer, the changed units' specs + current files, the whole-module verify errors, and the last
# rolled-back attempt. Written to workspaces/<proj>/signal/layer-<from>-<to>.md. This turns a dead-end
# into an actionable hand-off (the "signal for a senior" the methodology calls for). Args: proj, layers;
# the assembled report on stdin. Always exits 0 (it is a side-effect; the pass still fails downstream).
set -u
proj="${1:?usage: write_signal <proj> \"<from> <to>\" (report on stdin)}"
layers="${2:?usage: write_signal <proj> \"<from> <to>\" (report on stdin)}"
root="workspaces/$proj"; [ -d "$root" ] || root="$proj"
[ -d "$root" ] || { echo "no such workspace: $proj"; exit 0; }
set -- $layers; frm="${1:-from}"; to="${2:-to}"

dir="$root/signal"; mkdir -p "$dir"
out="$dir/layer-${frm}-${to}.md"
{
  echo "# STUCK LAYER: $frm -> $to  ($proj)"
  echo
  echo "The evolve pass for this layer failed the whole-module verify (go vet + go test -race) even after"
  echo "the bounded repair, and was rolled back. This is the escalation hand-off: the single stuck change,"
  echo "with everything needed to finish it by hand or with a stronger model. Fix the file(s) below to"
  echo "satisfy the specs AND clear the errors, drop them into the workspace, and re-run:"
  echo "    ratchet flow . evolve --ws $proj \"$frm $to\""
  echo
  cat
} > "$out"
echo "escalation signal written: workspaces/$proj/signal/layer-${frm}-${to}.md"
exit 0
