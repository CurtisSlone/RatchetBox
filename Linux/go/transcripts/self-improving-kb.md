# Self-improving KB (the compounding loop, automated)

We closed the evidence loop by hand twice (the atomic-alignment pitfall, the two-heap pattern). This
automates it: the ratchet mines its OWN run history for what the oracle keeps rejecting, finds which
failure classes the KB does not yet cover, and drafts the missing entry.

- **`mine_runs`** (`/do mine_runs`) - scans `runs/` for failed oracle steps, clusters the diagnostics
  into error classes, counts them, and checks each against `kb/pitfalls`. Recurring + uncovered = a KB
  candidate.
- **`learn`** flow (`ratchet flow . learn "<class>"`) - drafts a `pitfalls/` entry for that class
  (grounded on the existing entries' house style), validated by `write_pitfall` (needs an H1 title + a
  ```go example) and written + reindexed. Repairs once if malformed.

## The loop, run live

Step 1 - mine the runs. The top failures are already covered (the KB is doing its job); one class is not:

```text
$ /do mine_runs
class               count  covered?   example
unused-import           9  yes        ./main.go:10:2: "os" imported and not used
unused-var              7  yes        vet: ./main.go:108:2: declared and not used: encoder
wrong-behavior          2  yes        TEST FAILED:
syntax-error            1  NO         ./snippet.go:11:1: syntax error: non-declaration statement ...
KB CANDIDATES (recurring + not covered): syntax-error
```

Step 2 - learn the gap:

```text
$ ratchet flow . learn "syntax error: non-declaration statement outside function body ..."
  - step 1: learn.generate   - step 2: learn.write   - step 3: learn.done
wrote kb/pitfalls/top-level-statements.md ; reindexed kb/pitfalls
```

The drafted entry (model-written, house-style, first try):

```markdown
# Pitfall: top-level statements
Go REJECTS non-declaration statements at package scope ... only declarations (var/const/type/func) are
allowed at package level; statements must be inside a function.
...
// WRONG - syntax error: non-declaration statement outside function body
package main
import "fmt"
fmt.Println("hello")        // not allowed at package scope
// RIGHT - move the statement inside a function
func main() { fmt.Println("hello") }
```

Step 3 - re-mine: the loop has closed.

```text
$ /do mine_runs
syntax-error            1  yes
All recurring failure classes are already covered by kb/pitfalls. The loop is closed.
```

## Why this is the moat

Every other tool gets stale; this one gets BETTER the more you run it. Each oracle rejection is a data
point; `mine_runs` turns the accumulation into a coverage report; `learn` fills the gap. Over time the
generate step sees grounding for exactly the mistakes this model-on-this-codebase actually makes, so the
repair rate falls. The ratchet teaches itself from its own compiler.
