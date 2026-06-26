package main

import (
	"fmt"
	"greetmod/greeter"
)

func main() {
	result := greeter.Greet("world")
	fmt.Println(result)
}
