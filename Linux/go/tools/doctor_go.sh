#!/usr/bin/env bash
# doctor_go.sh - the toolbelt probe for `ratchet doctor` (the analog of cpp's doctor_cl.ps1).
# Reports what the Go toolbelt has. Exits 0 if the REQUIRED toolchain (`go`) is present, 1 if not.
# Optional linters (goimports, staticcheck, golangci-lint) are reported as info and never fail the
# check - flows gate on them via this report, they are not blockers. Referenced from ratchet.json
# requirements as { "name": "go toolbelt", "tool": "doctor_go" }. Run directly with `/do doctor_go`.
set -u

ok=0

probe_required() { # name cmd
  if command -v "$2" >/dev/null 2>&1; then
    printf '  [ok]   %-14s %s\n' "$1" "$($2 version 2>/dev/null | head -n1)"
  else
    printf '  [MISS] %-14s required - %s\n' "$1" "$3"
    ok=1
  fi
}

probe_sub() { # label  (go subcommand availability)
  if go help "$1" >/dev/null 2>&1; then
    printf '  [ok]   go %-11s available\n' "$1"
  else
    printf '  [MISS] go %-11s missing\n' "$1"
    ok=1
  fi
}

probe_optional() { # name cmd installhint
  if command -v "$2" >/dev/null 2>&1; then
    printf '  [ok]   %-14s (optional) present\n' "$1"
  else
    printf '  [--]   %-14s (optional) absent - %s\n' "$1" "$3"
  fi
}

echo "Go toolbelt:"
probe_required "go"     "go"     "install Go (https://go.dev/dl) and put it on PATH"
if command -v go >/dev/null 2>&1; then
  probe_required "gofmt" "gofmt" "ships with the Go toolchain; check your install"
  probe_sub "vet"
  probe_sub "test"
  probe_sub "build"
fi
echo "Optional linters (gate flows, never block):"
probe_optional "goimports"     "goimports"     "go install golang.org/x/tools/cmd/goimports@latest"
probe_optional "staticcheck"   "staticcheck"   "go install honnef.co/go/tools/cmd/staticcheck@latest"
probe_optional "golangci-lint" "golangci-lint" "https://golangci-lint.run install"

if [ "$ok" -eq 0 ]; then
  echo "OK: required Go toolchain present (build/vet/test oracle ready)."
else
  echo "MISSING: required Go toolchain incomplete - see [MISS] lines above."
fi
exit "$ok"
