Write a COMPLETE Go file for ONE unit of a larger system. Output ONLY the Go source - no prose, no
markdown fences.

The file you are writing is: {{ path }}

Rules:
- This unit is part of a single `package main` at the module root. Put `package main` at the top.
  Every other unit is in the SAME package, so call their types and functions directly - do NOT import
  other units and do NOT redeclare them.
- If and only if this file is `main.go`, give it `func main()`. Any other file must NOT define `func
  main` (it would be a duplicate).
- Include every standard-library `import` you use, and use every import and variable (Go rejects
  unused ones). Standard library only unless the spec says otherwise.

CRITICAL - call the EXISTING names EXACTLY as listed in "API ALREADY IN THE MODULE" below: use those
type, function, and method names with their exact signatures. Do NOT invent or rename them. If the spec
names something differently than the API, the API WINS.

## API ALREADY IN THE MODULE (the authoritative surface you may call - verbatim)
{{ api }}

## THIS UNIT'S SPEC (what to build; defer to the API above for any name it lists)
{{ spec }}

## THE MODULE SO FAR (PROJECT.md, go.mod, file tree)
{{ project }}

## Reference (retrieved for this spec; may be empty - use only if relevant)
### Go standard library
{{ stdlib_refs }}
### Third-party module (already in this workspace's go.mod)
{{ dep_refs }}
### Design pattern
{{ pattern_refs }}
