# greetmod - multi-package scaffolding transcript

Tests that the ratchet can scaffold a real Go **directory layout** (sub-packages) from specs, not just a flat `package main`. A spec declares `package: <name>`; the plan carries it as `pkg`; `plan_units` maps it to `<pkg>/<name>.go`; `module_api` advertises the import path; and the generate prompt derives the package from the path and imports siblings via `<module>/<pkg>`.

- Generated: 2026-06-26
- Command: `ratchet flow . compose --ws greetmod ""`  (console: `/flow compose --ws greetmod`)
- Result: built first try, no repairs; `go test ./...` -> `ok greetmod/greeter`.

## Input specs

### `app.spec`

```yaml
name: App
role: behavior
intent: the entry point - uses the greeter package from another directory
behavior:
  - in func main, call greeter.Greet("world") and print the result with fmt.Println
  - import the greeter package using the module path from go.mod followed by "/greeter"
constraints: package main at the module ROOT (main.go); imports + calls the greeter package; func main only here
```

### `greeter.spec`

```yaml
name: Greeter
role: component
package: greeter
intent: a tiny greeting library that lives in its OWN package/subdirectory
api:
  - func Greet(name string) string   // returns the string "Hello, " + name + "!"
constraints: standard library only; this file is greeter/greeter.go in `package greeter`; exported func
```

### `greeter_test.spec`

```yaml
name: GreeterTest
role: test
package: greeter
intent: prove Greet formats correctly
behavior:
  - a Go test file in `package greeter` (greeter/greeter_test.go) with func TestGreet(t *testing.T)
  - assert that Greet("world") == "Hello, world!"
constraints: standard library only (testing); same package as greeter.go (no import of greeter needed)
```

## Build plan (note the per-unit `pkg`)

```json
{"contracts":[{"name":"Greeter","normalizedFrom":"","type":"func"}],"units":[{"dependsOn":[],"name":"Greeter","pkg":"greeter","role":"component","spec":"greeter.spec"},{"dependsOn":["Greeter"],"name":"GreeterTest","pkg":"greeter","role":"test","spec":"greeter_test.spec"},{"dependsOn":["Greeter"],"name":"App","pkg":"","role":"behavior","spec":"app.spec"}]}
```

## Worklist (subdirectory paths)

```text
greeter/greeter.go greeter.spec
greeter/greeter_test.go greeter_test.spec
main.go app.spec
```

## Per-unit generation

### `greeter/greeter.go greeter.spec`

```go
package greeter

func Greet(name string) string {
	return "Hello, " + name + "!"
}
```

stage_build (go vet ./...): pass


### `greeter/greeter_test.go greeter_test.spec`

```go
package greeter

import "testing"

func TestGreet(t *testing.T) {
	got := Greet("world")
	want := "Hello, world!"
	if got != want {
		t.Errorf("Greet(\"world\") = %q, want %q", got, want)
	}
}
```

stage_build (go vet ./...): pass


### `main.go app.spec`

```go
package main

import (
	"fmt"
	"greetmod/greeter"
)

func main() {
	result := greeter.Greet("world")
	fmt.Println(result)
}
```

stage_build (go vet ./...): pass


## Final module_check + layout

```text
OK: module builds and tests pass with go1.26.4
?   	greetmod	[no test files]
ok  	greetmod/greeter	0.001s
```

Resulting directory layout (a real package structure, not flat):

```text
greetmod/
  go.mod        (module greetmod)
  main.go       package main  -> import "greetmod/greeter"
  greeter/
    greeter.go       package greeter  (exports Greet)
    greeter_test.go  package greeter  (TestGreet)
```

```text
$ go run .
Hello, world!
```