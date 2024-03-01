package main

import "runtime/debug"

func main() {
	slice := make([]string, 2, 4)
	Example(slice, "hello", 10)
}

func Example(slice []string, str string, i int) {
	debug.PrintStack()
}
