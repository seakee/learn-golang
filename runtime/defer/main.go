package main

import "fmt"

func function(index int, value int) int {

	fmt.Println(index)

	return index
}

func main() {
	defer function(1, function(3, 0))
	defer function(2, function(4, 0))
}
