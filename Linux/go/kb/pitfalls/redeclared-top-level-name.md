# Pitfall: redeclared top-level name

Go REJECTS redeclaration of any top-level name (like `func main`) when adding a new file to an existing package — this causes a compile error at build time because all top-level declarations must be unique across all files in the package.

- A package's top-level names (funcs, vars, types) must be declared only once across all files in that package.
- When adding a file to a package that already has `func main`, do not declare another `func main`.
- Do not redeclare any other top-level names the package already defines.

```go
// WRONG - build error: func main redeclared in this block
package main

func main() {
    println("hello from main.go")
}

// RIGHT - only add new functions/types, not main
package main

func main() {
    println("hello from main.go")
}

func helper() {
    println("helper function")
}
