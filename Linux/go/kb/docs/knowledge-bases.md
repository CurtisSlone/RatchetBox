# Knowledge bases (grounding the model in your conventions)

A knowledge base (KB) is a folder of plain-markdown topics the model retrieves from per step. Grounding
is how the output comes back in your conventions instead of generic training-set patterns.

Keywords: knowledge base, kb, grounding, search, retrieve, index, manifest, catalog, BM25, embedding, keywords

- A KB library is `kb/<name>/`, one topic per markdown file. Lead each file with a title and a one-line
  summary, and add a `Keywords:` line of the terms the model actually searches for.
- `ratchet index <dir>/kb/<name>` builds that library's `manifest.json` (id, title, summary, keywords)
  from file content. The high-level registry of all libraries is `kb/catalog.json` (name, path,
  doc-count, one-line subject), which the engine reads as its knowledge registry.
- Register libraries in `ratchet.json` under `knowledgeBases` (a `name`, a `path`, and `default: true`
  for the one used when no library is named).
- A step grounds on a KB with a `search` binding: a query (often a `{{ slot }}`), a library name, and
  `k` (how many top hits to inject). Retrieval is BM25 plus an embedding rerank.
- A KB helps **known-pattern recall** ("the model forgot the idiom"). It does not make the model reason
  or restructure; a wall that survives grounding is a capability or spec problem, not a missing-doc one.
