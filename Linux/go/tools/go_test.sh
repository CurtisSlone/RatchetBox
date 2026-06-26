#!/usr/bin/env bash
# go_test.sh - the BEHAVIOR oracle for the generate-test-repair loop.
# Reads a payload on stdin that carries TWO Go files, each introduced by a marker line:
#   // === solution.go ===
#   ...implementation (package solution, no func main)...
#   // === solution_test.go ===
#   ...package solution; func TestXxx(t *testing.T){...}
# Writes them into a throwaway module, normalizes with gofmt, then runs `go vet` and `go test`.
# Exit 0 iff vet is clean AND the tests pass. Otherwise prints the diagnostics and exits 1.
# This ratchets on behavior, not just "it compiles" - the headline of the go ratchet.
# Cross-platform (bash on Linux/WSL/macOS); add a .ps1 sibling for native Windows later.
set -u

payload="$(cat)"

# Models sometimes wrap output in Markdown code fences - or emit a dangling fence with no opener.
# A line that is only a fence marker (```), ```go, etc.) is never valid Go, so drop every such line.
payload="$(printf '%s\n' "$payload" | sed '/^[[:space:]]*```/d')"

if ! command -v go >/dev/null 2>&1; then
  echo "go not found on PATH (install Go / check PATH)"
  exit 1
fi

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

# Split the payload into files on `// === <path> ===` markers. Content before the first marker (if any)
# goes to solution.go, so a bare implementation with only a trailing test marker still lands correctly.
cur="$tmp/solution.go"
: > "$cur"
have_marker=0
while IFS= read -r line || [ -n "$line" ]; do
  if printf '%s' "$line" | grep -qE '^[[:space:]]*//[[:space:]]*===[[:space:]]*.+[[:space:]]*===[[:space:]]*$'; then
    rel="$(printf '%s' "$line" | sed -E 's@^[[:space:]]*//[[:space:]]*===[[:space:]]*(.+)[[:space:]]*===[[:space:]]*$@\1@' | sed 's/[[:space:]]*$//')"
    case "$rel" in
      */*|..*|/*) echo "rejecting unsafe file path in marker: $rel"; exit 1 ;;
    esac
    cur="$tmp/$rel"
    : > "$cur"
    have_marker=1
    continue
  fi
  printf '%s\n' "$line" >> "$cur"
done <<EOF
$payload
EOF

# A behavior oracle is meaningless without a test file. `go test` exits 0 on a package with no tests,
# which would be a silent false pass - so require at least one *_test.go.
if ! ls "$tmp"/*_test.go >/dev/null 2>&1; then
  echo "NO TEST: payload has no *_test.go file. Emit two files with markers:"
  echo "  // === solution.go ===   (the implementation, package solution)"
  echo "  // === solution_test.go ===   (package solution; func TestXxx(t *testing.T))"
  exit 1
fi

cd "$tmp" || { echo "cannot enter temp dir"; exit 1; }
go mod init snippet >/dev/null 2>&1
gofmt -w ./*.go 2>/dev/null   # cosmetic normalize; real syntax errors surface in vet/test below

# Race detector gates concurrency bugs go test alone misses - enable it where supported (CGO + a C
# compiler), skip it cleanly otherwise so the oracle still runs.
RACE=""
if [ "$(go env CGO_ENABLED 2>/dev/null)" != "0" ] && { command -v gcc >/dev/null 2>&1 || command -v clang >/dev/null 2>&1 || command -v cc >/dev/null 2>&1; }; then RACE="-race"; fi

vet_out="$(GOFLAGS=-mod=mod go vet ./... 2>&1)"
if [ $? -ne 0 ]; then
  echo "VET FAILED:"
  printf '%s\n' "$vet_out"
  exit 1
fi

test_out="$(GOFLAGS=-mod=mod go test $RACE ./... 2>&1)"
status=$?
if [ "$status" -eq 0 ]; then
  echo "OK: vet clean and tests pass${RACE:+ (with -race)} with $(go version | awk '{print $3}')"
  printf '%s\n' "$test_out"
else
  echo "TEST FAILED${RACE:+ (-race)}:"
  printf '%s\n' "$test_out"
fi
exit "$status"
