# The Oracle (verify, then advance)

The Oracle is the deterministic check that decides whether a step's output is allowed to advance the
chain. It is what makes a small local model reliable.

Keywords: oracle, gate, verify, deterministic, repair, won't break, validate, pass, fail

- The model **proposes** output into a step; a non-LLM check (a compiler, a linter, a parser, a test, a
  schema validator) **accepts or rejects** it. The chain advances only on a pass.
- The verdict is deterministic: the same output always gets the same answer. The model never grades its
  own work.
- On a failure, the exact errors are fed back for a **bounded repair** (one or a few re-attempts), then
  the run fails clean rather than looping forever.
- Hard limit, stated honestly: an Oracle pass means **"won't break," not "is correct."** It enforces
  form (it compiles, it parses, the test passes), not intent or quality. A clean compile is not a proof
  of behavior. A human still reviews intent.
- The Oracle is only as strong as the check you wire in. "It builds" is the floor; lint + tests +
  audits raise the bar.
