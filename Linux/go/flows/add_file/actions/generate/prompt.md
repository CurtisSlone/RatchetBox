Write a COMPLETE Go file to ADD to an existing module. Output ONLY the Go source - no prose, no
markdown fences.

The file you are writing is: {{ path }}

Rules:
- The module is a single `package main` at the root. Put `package main` at the top. Other files are in
  the SAME package - call their types and functions directly, do NOT import or redeclare them.
- Only `main.go` may define `func main`; any other file must not.
- Include every standard-library `import` you use, and use every import and variable (Go rejects unused
  ones). Standard library only unless the request says otherwise.
- Call the EXISTING names in "API ALREADY IN THE MODULE" verbatim - do not invent or rename them.

## API ALREADY IN THE MODULE (call these verbatim)
{{ api }}

## THE MODULE SO FAR (PROJECT.md, go.mod, file tree)
{{ project }}

## Request (what to add)
{{ request }}

## Reference (retrieved for this request; may be empty - use only if relevant)
### Go standard library
{{ stdlib_refs }}
### Third-party module (already in this workspace's go.mod)
{{ dep_refs }}
### Design pattern
{{ pattern_refs }}
