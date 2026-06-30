# Context Binding (each step sees only what it is given)

Context Binding is the rule that every chain step receives only its declared inputs, never a growing
transcript. It is the single biggest reliability lever for a weak model.

Keywords: context, binding, inputs, scoped, isolation, slot, grounding, search, ref, snippet

- A step's prompt is assembled from named **input bindings**, and nothing else. Each binding names a
  source: a prior step's output (`from`), a fixed reference entry (`ref`), or a retrieval query against
  a knowledge library (`search`).
- Bindings are capped in size (`max_chars`) and bound into named slots the prompt template fills.
- The model never sees a cumulative tape, prior prompts, or engine state. Every call gets a small,
  clean, known context. A weak model cannot be confused by accumulated noise it never receives.
- Because a step sees only its declared inputs, what it sends to a model is a controlled snippet, not
  the whole workspace. This is also what bounds what could ever leave the machine.
- This is the opposite of "prompt injection." Context Binding is a containment mechanism: the author
  controls exactly what reaches each step.
