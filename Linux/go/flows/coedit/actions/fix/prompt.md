Your previous multi-file edit failed the whole-module verify (go vet / go test -race). Return the
CORRECTED set. Fix exactly what the diagnostics report - usually a missed call site or import after a
signature change. Output ONLY the `=== path ===` marker blocks, every file complete.

## Requested change
{{ request }}

## Diagnostics (whole-module vet/test over the staged set, then rolled back)
{{ errors }}

## Original current files
{{ current }}

## Your previous attempt
{{ prev }}
