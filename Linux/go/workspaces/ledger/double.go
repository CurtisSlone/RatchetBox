package main

import "math"

func Double(n int) int {
	if n > math.MaxInt/2 {
		return math.MaxInt
	}
	return n * 2
}

func TestDouble() {
	result := Double(21)
	expected := 42
	if result != expected {
		panic("TestDouble failed")
	}
}
