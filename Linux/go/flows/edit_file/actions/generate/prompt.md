Apply the requested change to an existing Go file and return the COMPLETE updated file. Output ONLY the
Go source - no prose, no markdown fences.

The file you are editing is: {{ path }}

Rules:
- Return the WHOLE file, not a diff or a fragment. Preserve everything not affected by the request -
  keep the existing package clause, types, and functions unless the request says to change them.
- Keep `package main`. Other files in the module are in the same package; call their names directly, do
  not import or redeclare them. Do not add or remove `func main` unless the request asks.
- Include every import you use; remove any import or variable the change leaves unused.

## Current contents of {{ path }}
{{ current }}

## THE MODULE (PROJECT.md, go.mod, file tree)
{{ project }}

## Requested change
{{ request }}

## Reference (retrieved for this request; may be empty - use only if relevant)
### Go standard library
{{ stdlib_refs }}
### Third-party module (already in this workspace's go.mod)
{{ dep_refs }}
### Pitfalls to avoid (builds-but-wrong / redeclaration)
{{ pitfall_refs }}
