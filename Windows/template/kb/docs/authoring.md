# Authoring a ratchet (add a flow, a tool, or a KB)

A new capability is almost always a new flow, tool, or KB entry, not an engine change. The engine is a
domain-agnostic harness; the domain lives in the ratchet.

Keywords: author, add, build, flow, tool, kb, validate, doctor, ratchet.json, manifest, capability

- **Add a tool:** drop a script in `tools/`, declare it in `tools/manifest.json` (command, inputSchema,
  optional stdin, timeout). Make it exit 0 on pass, non-zero on fail; print diagnostics on fail so a
  repair step can use them.
- **Add a flow:** create `flows/<name>/chain.json` (entry, nodes, budgets) and a directory per node
  under `actions/<node>/` (`action.json`, plus `prompt.md` for a generate node). Wire repair as a
  failure edge back to the generate node, then forward to the gate again.
- **Add knowledge:** put markdown topics under `kb/<lib>/`, run `ratchet index <dir>/kb/<lib>`, and
  register the library in `ratchet.json` `knowledgeBases` (and `kb/catalog.json` where used).
- **Fix mechanical errors deterministically, not with prompts.** If a gate failure is about layout,
  imports, or a missing entry point, add a pipeline step that fixes the class. Reserve prompts for
  intent.
- **Verify your work:** `ratchet validate-flow <dir>` lints every chain (node kinds, fields, unknown
  tools, reachability); `ratchet doctor <dir>` preflights the declared toolchain.
