# Tools (declared scripts; the exit code is the verdict)

A tool is a deterministic script the host runs. Tools do the work code is good at: compile, parse,
scaffold, build, validate. An `action` node runs a tool and reads its exit code as the Oracle verdict.

Keywords: tool, script, command, exit code, stdin, manifest, action, validate, oracle, doctor

- Tools live in `tools/` and are declared in `tools/manifest.json`: a `name`, a `description`, a
  `command` (the argv to run), an `inputSchema` (its arguments), an optional `stdin` flag, and a
  `timeout`.
- The model never invents a command. A tool's command is authored in the ratchet; the model only fills
  declared arguments. Which tool runs is decided by the chain, not by the model.
- The host picks the interpreter by extension and OS: `.sh` runs under bash on Linux/macOS/WSL, `.ps1`
  under PowerShell on Windows, `.py` under python. So a ratchet runs wherever its declared tools exist.
- Exit 0 means pass (Oracle pass); non-zero means fail, and the tool's output is fed back for repair.
  Large payloads (a generated file) are passed on stdin when the tool declares `stdin`.
- `ratchet doctor <dir>` preflights the tools a ratchet declares it needs.
