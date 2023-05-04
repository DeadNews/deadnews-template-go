package main

import (
	"fmt"
)

// Get answer to the Meaning of Life, the Universe, and Everything.
func GetAnswer() int {
	result := 40 + 2
	return result
}

func main() {
	fmt.Println(GetAnswer())
	// Output: 42
}
