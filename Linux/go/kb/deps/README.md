# deps - third-party module API reference

Auto-populated by `tools/add_dep.sh` (`/do add_dep <workspace> <module>`): when you add a third-party
Go module to a workspace, its `go doc -all` is ingested here as one markdown file per package, so the
lifecycle/compose flows can ground generation on the module's real API (the same offline `go doc` trick
the `stdlib` library uses). Re-indexed automatically after each add.

This file is just the seed so the library is non-empty before the first dependency is added.
