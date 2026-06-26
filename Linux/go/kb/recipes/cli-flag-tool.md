# Recipe: a CLI tool (standard library flag)

A command-line tool using only `flag` + `os` - no dependencies. For subcommands, give each its own
`flag.FlagSet` and switch on `os.Args[1]`.

- Define flags, call `flag.Parse()`, read positional args from `flag.Args()`.
- Print usage to stderr and exit non-zero on bad input; keep `main` thin (parse, then call a function).
- For subcommands, one `flag.NewFlagSet(name, flag.ExitOnError)` per command.

```go
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: tool <greet|count> [flags]")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "greet":
		fs := flag.NewFlagSet("greet", flag.ExitOnError)
		name := fs.String("name", "world", "who to greet")
		_ = fs.Parse(os.Args[2:])
		fmt.Printf("Hello, %s!\n", *name)
	case "count":
		fs := flag.NewFlagSet("count", flag.ExitOnError)
		to := fs.Int("to", 5, "count up to")
		_ = fs.Parse(os.Args[2:])
		for i := 1; i <= *to; i++ {
			fmt.Println(i)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		os.Exit(2)
	}
}
```

For richer CLIs (nested commands, completion), use the `cobra` framework - see the cobra-cli recipe.
