#!/usr/bin/env bash
# write_pitfall.sh - the oracle + writer for the `learn` flow: validate a drafted kb/pitfalls entry
# (markdown on stdin) and write it. Well-formed = an H1 `# ...` title AND at least one ```go example
# block (the bad->good code). Slug derived from the title. Reindexes kb/pitfalls. Exit 1 (-> repair) if
# malformed. Args: none; the entry on stdin.
set -u
content="$(cat)"
# Strip a wrapping markdown fence if the model added one around the whole doc.
content="$(printf '%s\n' "$content" | sed '1{/^```markdown$/d}; ${/^```$/d}')"

title="$(printf '%s\n' "$content" | grep -m1 '^# ' | sed 's/^# *//')"
if [ -z "$title" ]; then echo "INVALID: no '# Title' heading"; exit 1; fi
if ! printf '%s\n' "$content" | grep -q '```go'; then echo "INVALID: needs a \`\`\`go example block (show the wrong code and the fix)"; exit 1; fi

slug="$(printf '%s' "$title" | tr '[:upper:]' '[:lower:]' | sed 's/^pitfall: *//; s/[^a-z0-9]\{1,\}/-/g; s/^-//; s/-$//')"
[ -z "$slug" ] && slug="pitfall-entry"
out="kb/pitfalls/$slug.md"
printf '%s\n' "$content" > "$out"
echo "wrote $out"
if command -v ratchet >/dev/null 2>&1; then ratchet index kb/pitfalls >/dev/null 2>&1 && echo "reindexed kb/pitfalls"; fi
echo "OK: pitfall '$title' added to the KB"
exit 0
