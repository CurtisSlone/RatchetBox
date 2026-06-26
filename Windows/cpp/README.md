# cpp - a C++ ratchet

Generate, edit, and build **C++** with the local model, gated by the MSVC compiler. The model proposes;
`cl` decides. Grounded in curated knowledge (the standard library, the C++ Core Guidelines, design
patterns, MSVC errors, and a large slice of the Win32 API), so the small model writes idiomatic,
compiling code - and when it doesn't, the compiler's errors drive a bounded repair.

Open it from a [Ratchet](https://github.com/scanset/Ratchet) checkout:

```
.\ratchet.cmd ..\RatchetBox\cpp
```

## First: check your toolchain

```
.\ratchet.cmd doctor ..\RatchetBox\cpp
```

This validates everything the ratchet needs: MSVC `cl` (via `vcvars32`), the Windows SDK, Ollama + the
configured models, and the knowledge bases. `git`/`vcpkg` are optional (only for third-party packages).

## What it can do

| Flow | What it does |
| --- | --- |
| `cpp` | One-off: describe a task, get a complete C++ program, syntax-checked by `cl` and repaired up to twice. |
| `add_file` | Add a new file to the active project; build the whole project; repair; record it. |
| `edit_file` | Surgically edit an existing file (refactor, fix, extend); rebuild; repair. |
| `answer` | Grounded Q&A over the knowledge bases - no build, just an explained answer with sources. |

Each generating flow first **plans** which knowledge bases the task needs, then retrieves only those -
so a `std::filesystem` task pulls cppreference, a window app pulls Win32, a factory pulls the patterns
KB, and nothing irrelevant is dragged in.

## Drive it

One-off generation:

```
/flow cpp "read a CSV from stdin and print column means, using the STL"
```

A real project (lives under `workspaces/`, built to a real `.exe`):

```
/do new_project Stats              # scaffold workspaces/Stats (or 'win32' for a GUI app)
/ws switch Stats
/flow add_file src/core/stats.cpp  "mean, median, variance over a vector<double>"
/flow edit_file src/main.cpp       "read numbers from stdin and print the three stats"
/do build_project Stats            # compile + link -> dist/Stats.exe
```

Ask a question instead of building:

```
/flow answer "what is the difference between std::vector reserve and resize?"
```

## Adding libraries

- **Windows SDK libs** (GUI, graphics, audio, networking - all already on the box):
  `/do link_libs -Proj Stats -Libs "d2d1.lib dwrite.lib user32.lib"`, then `#include` and build.
- **Third-party** (fmt, nlohmann-json, sqlite, ...): one-time `/do bootstrap_vcpkg`, then
  `/do add_package -Proj Stats -Name fmt -Libs "fmt.lib"`. It provisions the library and wires it into
  the build; `build_project` even copies the needed DLLs next to your exe.

See `kb/howto/adding-packages.md` (or ask the `answer` flow) for the full guide.

## Good to know

- **C++17, MSVC `cl`.** The Oracle is the compiler, so generated code is checked for real, not guessed.
- **Smart App Control** may block *running* a freshly built `.exe` (it doesn't block the build). If a
  run is blocked, that's SAC, not the ratchet - the build is the verification.
- Operating rules and hard limits are in `SYSTEM.md`; the layout is in `STRUCTURE.md`.
