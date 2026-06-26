# Pitfall: unused imports and variables

Go REJECTS an unused import or an unused local variable - these are compile errors, not warnings. The
local model often adds an import while planning code it then does not write (observed: it imported
`strconv`/`errors` it never used, costing a repair round).

- Import a package only when you actually call it; add the import when you write the call, not before.
- Remove any local you stop using. Use `_` only as a deliberate discard, never to silence dead code.
- Blank-import `_ "pkg"` only for intended side effects (e.g. registering a driver).

```go
// WRONG - does not compile: "strconv" imported and not used
import (
	"fmt"
	"strconv" // never referenced below
)

func Greet(name string) string { return fmt.Sprintf("hi %s", name) }

// RIGHT - import only what you use
import "fmt"

func Greet(name string) string { return fmt.Sprintf("hi %s", name) }

// Deliberate discard of one return value:
_, err := fmt.Println("x")
```
