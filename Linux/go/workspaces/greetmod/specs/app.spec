name: App
role: behavior
intent: the entry point - uses the greeter package from another directory
behavior:
  - in func main, call greeter.Greet("world") and print the result with fmt.Println
  - import the greeter package using the module path from go.mod followed by "/greeter"
constraints: package main at the module ROOT (main.go); imports + calls the greeter package; func main only here
