#!/usr/bin/env bash
# bump.sh - change a dependency or toolchain version in a workspace, VERIFIED with rollback. A version
# bump is a change, so the oracle gates it: snapshot go.mod/go.sum, apply the change, `go mod tidy`, then
# run the full go_quality gate (build, test -race, vet, staticcheck, govulncheck). Keep the change only
# if everything is still clean; otherwise restore go.mod/go.sum to the pre-bump state. Args: proj, change
# (one of: "<pkg>@<ver|latest>", "go=<ver>", "-u" to upgrade all, "tidy"). Run: /do bump <proj> <change>.
set -u
proj="${1:?usage: bump <proj> <change>}"
change="${2:?missing change: <pkg>@<ver> | go=<ver> | -u | tidy}"
root="$proj"; [ -d "$root" ] || root="workspaces/$proj"
[ -d "$root" ] || { echo "no such workspace: $proj"; exit 1; }
[ -f "$root/go.mod" ] || { echo "no go.mod in $proj"; exit 1; }

snap="$(mktemp -d)"; trap 'rm -rf "$snap"' EXIT
cp "$root/go.mod" "$snap/go.mod"
[ -f "$root/go.sum" ] && cp "$root/go.sum" "$snap/go.sum"
restore() {
  cp "$snap/go.mod" "$root/go.mod"
  if [ -f "$snap/go.sum" ]; then cp "$snap/go.sum" "$root/go.sum"; else rm -f "$root/go.sum"; fi
}

echo "== bump $proj: $change =="
applied=0
case "$change" in
  go=*)   ( cd "$root" && go mod edit -go="${change#go=}" ) || applied=1 ;;
  -u|all) o="$(cd "$root" && GOFLAGS=-mod=mod go get -u ./... 2>&1)"; st=$?; printf '%s\n' "$o"; [ "$st" -eq 0 ] || applied=1 ;;
  tidy)   : ;;
  *)      o="$(cd "$root" && GOFLAGS=-mod=mod go get "$change" 2>&1)"; st=$?; printf '%s\n' "$o"; [ "$st" -eq 0 ] || applied=1 ;;
esac
if [ "$applied" -ne 0 ]; then echo "CHANGE FAILED to apply (e.g. no such version); go.mod restored."; restore; exit 1; fi
( cd "$root" && GOFLAGS=-mod=mod go mod tidy ) >/dev/null 2>&1

echo "-- verifying with go_quality --"
if bash tools/go_quality.sh "$proj"; then
  echo "OK: '$change' verified clean and KEPT."
  exit 0
else
  echo "VERIFICATION FAILED - rolling back go.mod/go.sum to the pre-bump state."
  restore
  echo "REVERTED: '$change' did not pass the quality gate; the workspace is unchanged."
  exit 1
fi
