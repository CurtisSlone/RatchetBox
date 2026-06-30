# Driving a ratchet (the console and one-shot commands)

You drive a ratchet from its operator console, or run a single flow in one shot. The local model does
the generating; you (or a frontier model over MCP) drive.

Keywords: drive, console, command, slash, chat, search, flow, do, ws, route, runs, rollback, mcp, one-shot

- Open the console: `ratchet <dir>` (or `ratchet chat <dir>`). Inside it:
  - plain text - ordinary ungrounded chat with the model.
  - `/search [library] <question>` - a grounded answer retrieved from a knowledge base.
  - `/flow <name> [input]` - run an action chain by name.
  - `/do <tool [arg] | shell command>` - run a declared tool, or a command you paste.
  - `/ws switch|create <name>` - switch or create the active workspace (the session focus).
  - `/route <request>` - let the model pick the best flow (you confirm before it runs).
  - `/runs [n]` - list recent runs (the audit log); `/rollback [id]` - restore the workspace to a run's
    pre-state; `/snapshot` - save a manual restore point.
  - `/flows`, `/tools`, `/note`, `/help`, `/clear`, `/quit`.
- One-shot, no console: `ratchet flow <dir> <name> [--ws <workspace>] [input...]`.
- Over MCP: `ratchet mcp <dir>` exposes the same flows so a frontier model can drive the engine.
