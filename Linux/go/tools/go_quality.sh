#!/usr/bin/env bash
# go_quality.sh - the production gate: run the full quality/security suite over a workspace and FAIL if
# any gate fails. gofmt -> go vet -> go build -> go test -race -> staticcheck -> govulncheck -> gosec.
# Unlike the per-unit oracles (where staticcheck is advisory), here EVERY present tool gates. Tools that
# are not installed are skipped with a note, so it runs anywhere and gets stricter as you install more.
# govulncheck needs network. Arg: proj. Run: /do go_quality <proj>, or via the `harden` flow.
set -u
proj="${1:?usage: go_quality <proj>}"
root="$proj"; [ -d "$root" ] || root="workspaces/$proj"
[ -d "$root" ] || { echo "no such workspace: $proj"; exit 1; }
cd "$root" || { echo "cannot enter $root"; exit 1; }

fail=0
gate() { # label  output  status
  if [ "$3" -ne 0 ]; then echo "FAIL ($1):"; printf '%s\n' "$2"; fail=1; else echo "ok   ($1)"; fi
}

RACE=""
if [ "$(go env CGO_ENABLED 2>/dev/null)" != "0" ] && { command -v gcc >/dev/null 2>&1 || command -v clang >/dev/null 2>&1 || command -v cc >/dev/null 2>&1; }; then RACE="-race"; fi

echo "== go_quality: $root =="
o="$(gofmt -l . 2>&1)"; if [ -n "$o" ]; then echo "FAIL (gofmt): not gofmt-clean:"; printf '%s\n' "$o"; fail=1; else echo "ok   (gofmt)"; fi
o="$(GOFLAGS=-mod=mod go vet ./... 2>&1)";            gate "go vet" "$o" $?
o="$(GOFLAGS=-mod=mod go build ./... 2>&1)";          gate "go build" "$o" $?
o="$(GOFLAGS=-mod=mod go test $RACE ./... 2>&1)";     gate "go test${RACE:+ -race}" "$o" $?
if command -v staticcheck >/dev/null 2>&1; then o="$(staticcheck ./... 2>&1)"; gate "staticcheck" "$o" $?; else echo "--   staticcheck absent (skipped)"; fi
if command -v govulncheck >/dev/null 2>&1; then o="$(govulncheck ./... 2>&1)"; gate "govulncheck (known CVEs)" "$o" $?; else echo "--   govulncheck absent (skipped)"; fi
if command -v gosec >/dev/null 2>&1; then o="$(gosec -quiet ./... 2>&1)"; gate "gosec (SAST)" "$o" $?; else echo "--   gosec absent (skipped)"; fi

echo "================================="
if [ "$fail" -eq 0 ]; then echo "PRODUCTION-CLEAN: all available gates passed."; else echo "NOT CLEAN: a gate failed (see above)."; fi
exit "$fail"
