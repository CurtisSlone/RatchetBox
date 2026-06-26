# Quality + security gates (deepening the oracle)

The ratchet's power is "the oracle gates every step." This deepens that gate. Two layers:

- **Per-unit oracles** (`go_test`, `stage_check`, `module_check`) now run `go test -race` (gating - a
  data race fails the build) and `staticcheck` (advisory - findings surface but do not block, so flows
  stay reliable). Race is enabled where CGO + a C compiler exist, skipped cleanly otherwise.
- **The `harden` flow** (`go_quality` tool) is the production GATE: it runs the whole suite over a
  workspace - gofmt, go vet, go build, `go test -race`, staticcheck, govulncheck (known CVEs), gosec
  (SAST) - and passes only if EVERY installed tool is clean. Absent tools are skipped, so it runs
  anywhere and gets stricter as you install more. `doctor_go` reports which are present.

Run: `ratchet flow . harden --ws <proj> ""` (console: `/ws switch <proj>` then `/flow harden`), or
`/do go_quality <proj>`.

## A clean module passes the whole stack

`pulsehook2` (the hardened webhook server, concurrent worker pool):

```text
$ /do go_quality pulsehook2
== go_quality: workspaces/pulsehook2 ==
ok   (gofmt)
ok   (go vet)
ok   (go build)
ok   (go test -race)        # the concurrent dispatcher is race-clean
ok   (staticcheck)
ok   (govulncheck (known CVEs))
--   gosec absent (skipped)
PRODUCTION-CLEAN: all available gates passed.        # exit 0
```

## A broken module is caught, loudly

`shorturl` (the partial compose, with an unused import + variable):

```text
$ ratchet flow . harden --ws shorturl ""
  - step 1: harden.check (action)
  - step 2: harden.fail (exit)
FAIL (go vet):
FAIL (go build):     ./main.go:10:2: "os" imported and not used
FAIL (go test -race):
FAIL (staticcheck):
FAIL (govulncheck (known CVEs)):
NOT CLEAN: a gate failed (see above).                # exit 1 -> harden.fail
```

## Why two layers (gate vs advisory)

A small local model can fix a compile/test/race failure reliably, so those GATE in the per-unit oracles
(repair loop). staticcheck/gosec have occasional false positives and harder fixes, so they are advisory
in the inner loop (keeps generation flowing) and only GATE in `harden`, which you run deliberately when
you want production-clean - and fix findings with `edit_file`. govulncheck (network, time-sensitive) is
harden-only, not on every compose.
