name: GreeterTest
role: test
package: greeter
intent: prove Greet formats correctly
behavior:
  - a Go test file in `package greeter` (greeter/greeter_test.go) with func TestGreet(t *testing.T)
  - assert that Greet("world") == "Hello, world!"
constraints: standard library only (testing); same package as greeter.go (no import of greeter needed)
