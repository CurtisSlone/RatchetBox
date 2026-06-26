# Working on an existing repo + type-safe refactor

The lifecycle flows assumed `workspaces/<proj>`. This lets the ratchet drive a Go module that already
exists ANYWHERE on disk, and adds a type-safe rename.

- **`link_repo`** (`/do link_repo <name> <path>`) symlinks an external module in as
  `workspaces/<name>`, so every flow (`add_file`, `edit_file`, `harden`, `run`, `module_api`, `refactor`)
  operates on the real repo unchanged. Its own package layout is respected (the lifecycle derives a
  file's package from its path). The build/test/vet oracles protect the real code - a change that breaks
  the module is rejected before it sticks.
- **`refactor`** flow (`ratchet flow . refactor --ws <proj> "rename X to Y"`) - the model reads the
  module API, maps the request to the exact target, and `gorename` renames every reference across files
  BY TYPE (not text), then verifies the build.

## Driving a real external repo (at /tmp/extmod, a `Ledger` module)

```text
$ /do link_repo extmod /tmp/extmod
  linked workspaces/extmod -> /tmp/extmod   (module: extmod)

$ /do go_quality extmod                      # harden the repo as-is
  PRODUCTION-CLEAN: all available gates passed.
```

Add a feature to the real repo, grounded on its existing API (`module_api` showed `Ledger.Deposit`,
`Ledger.Balance`):

```text
$ ratchet flow . add_file --ws extmod "transfer.go a Transfer(l *Ledger, from, to string, amount int)
                                        error using the existing Ledger methods; error if insufficient"
```
```go
// landed at /tmp/extmod/transfer.go - calls the REAL Ledger API
func Transfer(l *Ledger, from, to string, amount int) error {
	if l.Balance(from) < amount {
		return errors.New("insufficient balance")
	}
	l.Deposit(from, -amount)
	l.Deposit(to, amount)
	return nil
}
```

Type-safe rename across the whole module (definition + every caller, including the just-added file):

```text
$ ratchet flow . refactor --ws extmod "rename the Credit method to Deposit"
  - refactor.api -> refactor.plan -> refactor.apply -> refactor.done
  Renamed 4 occurrences in 3 files in 1 package.
  OK: renamed Ledger.Credit -> Deposit across the module; it still builds.
```

## The takeaway

`link_repo` turns the whole ratchet - generate, grow, harden, refactor - onto code you already have,
with the oracles guarding every change. `gorename` does renames by type, so it cannot rename the wrong
symbol, and the post-rename build is verified. This is the daily-driver step: point it at your repo and
work.
