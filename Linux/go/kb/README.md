# kb/ - the knowledge libraries

Indexed, routable knowledge for the go ratchet, organized cpp-style: one library per subdirectory,
each registered separately in `ratchet.json` `knowledgeBases[]` and indexed into its own
`manifest.json` (id/title/summary/keywords, extracted deterministically by `ratchet index`). Chains
ground on a library via a `search`/`ref` input binding; `/search <lib> <query>` answers from it.

## Libraries

| Dir | Name | Source | Built by | Status |
|---|---|---|---|---|
| `idioms/`  | `kb` (default) | hand-written Go idioms for buildable code | author | seed |
| `stdlib/`  | `stdlib`       | the Go standard library API reference (176 packages) | `tools/kb_ingest_godoc.sh` (offline `go doc -all`) | built |
| `patterns/`| `patterns`     | 23 GoF patterns (ingested) + authored algorithm patterns (e.g. two-heap median) | `tools/kb_ingest_patterns.sh` + authored | built (24) |
| `deps/`    | `deps`         | third-party module API (go doc per added module) | `tools/add_dep.sh` (`/do add_dep <ws> <module>`) | built (auto-populated) |
| `guidelines/`| `guidelines` | idiomatic-style entries from Effective Go / Code Review Comments / Go Proverbs (real code) | authored from canonical sources (see `guidelines/ATTRIBUTION.md`) | built (12) |
| `pitfalls/`  | `pitfalls`   | builds-but-wrong traps (unused imports, nil-map write, loop-var capture, channel deadlock, slice aliasing) | authored from `runs/` evidence | built (5) |
| `spec/`      | (planned) | the Go language spec, split by section | ingest (split) | planned |

The `stdlib` library is the Go analog of the cpp ratchet's `cppref` - except `go doc` generates it
locally from the installed toolchain, so there is nothing to scrape and it always matches your Go
version. Regenerate after a toolchain upgrade.

## Build / rebuild

```
tools/kb_ingest_godoc.sh kb/stdlib                    # (re)generate stdlib, offline, deterministic
tools/kb_ingest_patterns.sh <patterns-repo> kb/patterns   # (re)generate the design-patterns library
ratchet index kb/stdlib                # (re)build a library's routing manifest (after any kb/<lib> edit)
ratchet index kb/patterns
ratchet index kb/idioms
```

Add a library: create `kb/<lib>/` with one-topic-per-file markdown, register it in `ratchet.json`
`knowledgeBases[]` (and `requirements[]`), then `ratchet index kb/<lib>`. Wire a flow to it with a
`{ "search": "<lib>", "query": "{{ ... }}", "k": N, "as": "refs" }` binding (multi-library plan-routed
retrieval, cpp-style, is the next phase - see `../ROADMAP.md`).
