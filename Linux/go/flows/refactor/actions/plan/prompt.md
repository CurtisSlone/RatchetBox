Map the rename request to a target name and a new name. Output ONLY the JSON object.

- `target`: the EXISTING identifier to rename, written exactly as it appears in the module API below.
  For a method, use `Type.Method` (e.g. `Ledger.Balance`); for a package-level function/type, use the
  bare name (e.g. `Transfer`). Do not include the package path.
- `newname`: the new identifier (a valid Go name).

## Request
{{ request }}

## Module API (the existing names you may rename)
{{ api }}
