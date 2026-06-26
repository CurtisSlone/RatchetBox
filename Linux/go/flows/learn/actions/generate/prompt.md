Write a knowledge-base entry for the `pitfalls` library describing the Go failure class below. A
pitfall entry teaches a "builds-but-wrong" trap or a thing the compiler rejects that the model keeps
getting wrong. Match the house style of the examples below. Output ONLY the markdown.

Required shape:
- An H1 title: `# Pitfall: <short title>`
- One or two sentences naming the trap: what goes wrong, and why (does it fail at build, vet, or run?).
- 2-4 bullet points with the rule(s) to follow.
- A single ```go block showing the WRONG code (commented `// WRONG - <why>`) and the RIGHT fix
  (commented `// RIGHT - <why>`). Keep it minimal and self-contained.

## Failure class to document
{{ klass }}

## House style (existing pitfalls entries - match their shape)
{{ style_refs }}

## Relevant stdlib (if any)
{{ stdlib_refs }}
