# Flows and action chains (the node graph)

A flow is an action chain: a fixed graph of steps on disk that the engine walks one node at a time. The
model fills one slot per step; deterministic code sequences the steps.

Keywords: flow, chain, action, node, graph, generate, ai_branch, summarizer, foreach, exit, on_success, on_failure

- A flow lives at `flows/<chain>/chain.json` (the graph: `id`, `entry`, `nodes`, `budgets`) plus one
  directory per node under `flows/<chain>/actions/<node>/` (an `action.json`, and a `prompt.md` for AI
  nodes).
- Each node declares `inputs` (the bindings it sees), and edges `on_success` / `on_failure` (or
  `transitions` for a branch). The graph is lint-checked before anything runs.
- Node kinds:
  - `generate` - the model proposes text/code into a slot (add `output_schema` for structured JSON).
  - `action` - run a declared tool; the tool's exit code is the Oracle verdict.
  - `ai_branch` - the model picks one value from a fixed enum, which routes the chain.
  - `summarizer` - deterministically concatenate prior outputs.
  - `foreach` - run a sub-chain once per item in a list (fan-out).
  - `exit` - terminal node with an outcome.
- Repair is wired as edges: a gate's `on_failure` routes back to a generate node with the verdict bound
  in, then forward to the gate again. The model never picks the next step; the graph does.
