<p align="center"><img src="ratchet.png" alt="Ratchet" width="320"></p>

# RatchetBox

A collection of ready-made **ratchets** for the [Ratchet](https://github.com/scanset/Ratchet)
engine - the cross-platform host (Windows, Linux, macOS) that runs a small local model as a
*constrained proposer*, gated by a deterministic Oracle. Ratchet is domain-agnostic; the domain lives
in a ratchet. This repo is a box of them.

> **New here?** A *ratchet* is a self-contained directory (`ratchet.json` + `flows/`, `tools/`, `kb/`)
> that you point the engine at: `ratchet <dir>`. The engine and how it works live in the
> [Ratchet repo](https://github.com/scanset/Ratchet).

## Layout

Ratchets are grouped by the platform their **toolchain** targets (the oracle a ratchet runs - csc,
cl, go build - decides where it runs). The engine itself is cross-platform.

```
Windows/    ratchets whose tools are PowerShell / Windows toolchains
  dotnet4-x   Windows C# / PowerShell (in-box csc, C# 5)
  cpp         C++ with MSVC cl (vcvars32; STL + Win32)
  template    skeleton to copy (PowerShell-tooled)
Linux/      ratchets whose tools are POSIX / cross-platform
  go          Go code generation verified with `go build` (bash oracle; runs on Linux/WSL/macOS)
```

## The ratchets

| Ratchet | Platform | What it builds | Highlights |
| --- | --- | --- | --- |
| **`Windows/dotnet4-x`** | Windows | C# / PowerShell, in-box `csc` (C# 5) | `csharp`/`winforms`/`powershell` one-offs + a project lifecycle (`new_project` -> `add_file`/`edit_file` -> build) + **spec-to-code composition** (`compose` a whole multi-file system from a dir of `.spec` files); `transcripts/` has end-to-end build transcripts |
| **`Windows/cpp`** | Windows | C++ with MSVC `cl` (vcvars32; STL + Win32) | plan-routed retrieval over 7 knowledge bases, generate-compile-repair, a full project lifecycle, native-app linking, vcpkg packages + **spec-to-code composition**; `transcripts/` has a worked build |
| **`Windows/template`** | Windows | (nothing - a skeleton to copy) | the lifecycle + composition skeleton: working flows with the domain-specific tools as `CHANGE_ME` stubs to fill in |
| **`Linux/go`** | Linux / macOS | Go, verified with `go build` / `go vet` / `go test` | generate (`go`) + **behavior** generate-with-tests (`test`), a project lifecycle (`new_module` -> `add_file`/`edit_file` -> `run`), **spec-to-code composition** (`compose`), third-party modules with auto-ingested docs (`add_dep`), and plan-routed retrieval over 7 KBs (stdlib via `go doc`, GoF + algorithm patterns, Effective Go / Code Review Comments guidelines, builds-but-wrong pitfalls, deps, idioms); `transcripts/` has worked builds incl. a low-latency webhook server. The cross-platform reference (the bash oracles run on Linux, WSL, and macOS) |

## Quick start

```
# 1. Get the engine (see its README for build + Ollama setup)
git clone https://github.com/scanset/Ratchet

# 2. Get this box of ratchets, alongside it
git clone https://github.com/scanset/RatchetBox

# 3. Open one (from the Ratchet repo)
cd Ratchet
ratchet ../RatchetBox/Linux/go               # Linux/macOS/WSL
.\ratchet.cmd ..\RatchetBox\Windows\dotnet4-x   # Windows (in-memory launcher)

# 4. Sanity-check a ratchet's toolchain first (optional)
ratchet doctor ../RatchetBox/Linux/go
```

Each ratchet has its own `README.md` (what it does + how to drive it) and `AGENTS.md` (orientation for
AI agents). Start there. `Windows/dotnet4-x`, `Windows/cpp`, and `Linux/go` carry a `transcripts/`
folder with real end-to-end build transcripts (the `Linux/go` set includes a low-latency webhook server
built from specs and run) - the fastest way to see what driving the ratchet looks like.

## Make your own

For a Windows toolchain, copy `Windows/template`, rename it, edit its `ratchet.json`, implement the
`CHANGE_ME` tools, and fill `kb/` + the flow prompts. For a cross-platform (POSIX) ratchet, `Linux/go`
is the smallest working reference to start from. The full authoring contract is in the engine's docs
(`docs/how-to/build-a-ratchet.md`, `docs/how-to/author-flows.md`, `docs/how-to/author-tools.md`,
`docs/how-to/compose-from-specs.md`).

## License

Apache 2.0 - see [LICENSE](LICENSE) and [NOTICE](NOTICE).
