# Pitfall: top-level statements

Go REJECTS non-declaration statements at package scope - `go build` fails with a syntax error because only declarations (var, const, type, func) are allowed at package level; statements like assignments or function calls must be inside a function.

- Only declarations (var, const, type, func) are allowed at package scope.
- All executable code must live inside a function body.
- Package-level code can only contain top-level declarations and comments.

```go
// WRONG - syntax error: non-declaration statement outside function body
package main

import "fmt"
fmt.Println("hello") // not allowed at package scope

// RIGHT - move the statement inside a function
package main

import "fmt"

func main() {
	fmt.Println("hello")
}
