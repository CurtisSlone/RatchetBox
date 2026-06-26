# AGENTS.md - the template ratchet (for AI agents)

The skeleton to copy when starting a new ratchet for the
[Ratchet](https://github.com/scanset/Ratchet) engine. Humans: see [README.md](README.md). Layout:
[STRUCTURE.md](STRUCTURE.md).

## To create a new ratchet

1. Copy this directory; rename it.
2. Edit `ratchet.json`: `name`, `domain`, `models` (must exist in Ollama), and `requirements` (what
   `ratchet doctor` validates).
3. Implement the `CHANGE_ME` tools (`build_project`/`new_project`/`project_api`) and fill the
   `<LANGUAGE>` placeholders in the flow prompts for your domain. See `RatchetBox/dotnet4-x` and
   `RatchetBox/cpp` for full implementations.

## The buckets

- `flows/` - action chains. Ships six: **add_file** / **edit_file** (the file lifecycle),
  **compose** / **add_unit** (multi-spec composition), **generate** (a one-off snippet: generate ->
  single-file `check` -> repair) and **answer** (grounded Q&A from the kb). A chain is `flows/<chain>/chain.json` +
  `actions/<node>/{action.json, prompt.md}`. Node kinds: `action`, `generate` (add `output_schema` for
  structured JSON), `ai_branch`, `summarizer`, `foreach` (sequential sub-chain fan-out, used by
  compose), `exit`. Each node sees only its declared, scoped inputs (Context Binding).
- `tools/` - scripts declared in `tools/manifest.json` (command, inputSchema, optional `stdin`, timeout);
  exit code is the Oracle verdict. GENERIC (work as-is): `read_specs`, `plan_units`, `read_project`,
  `register_file`, `log_edit`, `read_file`. **CHANGE_ME** (your domain): `build_project` + `check` (the
  whole-project and single-file oracles), `new_project` (scaffold), `project_api` (existing-units API),
  + the knobs in `plan_units`/`stage_and_build`.
- `kb/` - one topic per markdown file under `kb/<lib>/`; `ratchet index <dir>\kb` builds the manifest;
  register libraries in `ratchet.json` `knowledgeBases[]`.
- `schemas/` + `samples/` - optional TSV oracle (a `<table>.json` schema + `<table>.txt` rows).

## Conventions

- Keep domain specifics here, never in the engine - new capability = a new flow or tool, no host change.
- stdin payloads to tools arrive with a UTF-8 BOM; strip it (read raw bytes).
- Verify: `ratchet validate-flow <dir>` (lint chains) and `ratchet doctor <dir>` (preflight tools).

Full authoring contract: the engine's `docs/how-to/build-a-ratchet.md`, `docs/how-to/author-flows.md`,
`docs/how-to/author-tools.md`, `docs/concepts/context-binding.md`, and the composition guide
`docs/how-to/compose-from-specs.md`.
