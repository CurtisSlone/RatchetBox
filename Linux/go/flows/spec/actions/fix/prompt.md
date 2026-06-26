Your previous spec draft was not well-formed. Return the CORRECTED spec(s). Fix exactly what the
validator reported. Output ONLY the marker blocks - no prose, no code fences.

Each spec must have, under its `=== <name>.spec ===` marker: a `name:` line, a `role:` line whose value
is one of data/interface/component/behavior/gui/test, and at least one of `api:`/`behavior:`/`intent:`.
Exactly one unit is role behavior/gui (the entry).

## Request
{{ desc }}

## Validator errors
{{ errors }}

## Previous draft
{{ prev }}
