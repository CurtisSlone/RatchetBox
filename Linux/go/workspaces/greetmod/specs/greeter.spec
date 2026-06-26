name: Greeter
role: component
package: greeter
intent: a tiny greeting library that lives in its OWN package/subdirectory
api:
  - func Greet(name string) string   // returns the string "Hello, " + name + "!"
constraints: standard library only; this file is greeter/greeter.go in `package greeter`; exported func
