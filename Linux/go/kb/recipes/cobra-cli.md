# Recipe / framework profile: a cobra CLI

`github.com/spf13/cobra` is the standard framework for rich CLIs (nested subcommands, flags, help,
completion) - it powers kubectl, hugo, gh. Use it when the stdlib `flag` recipe outgrows a single
level of commands.

Profile (one-time per workspace):

```
/do add_dep <workspace> github.com/spf13/cobra
```

Key API: a root `*cobra.Command`; child commands added with `root.AddCommand(child)`; each command has
`Use`, `Short`, and a `Run`/`RunE func(cmd *cobra.Command, args []string)`; flags via
`cmd.Flags().StringVar(...)`; execute with `root.Execute()` in `main`.

```go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var name string

	root := &cobra.Command{Use: "tool", Short: "a demo CLI"}

	greet := &cobra.Command{
		Use:   "greet",
		Short: "print a greeting",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Hello, %s!\n", name)
			return nil
		},
	}
	greet.Flags().StringVar(&name, "name", "world", "who to greet")

	root.AddCommand(greet)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

Keep `main` thin (build the command tree, `Execute`); put command logic in `RunE` functions that return
errors rather than calling `os.Exit` deep in the tree.
