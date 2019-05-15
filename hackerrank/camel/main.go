package main

import (
	"fmt"
	"unicode"
)

func main() {
	var input string
	fmt.Scanln(&input)

	count := 0
	for _, c := range input {
		if unicode.IsUpper(c) {
			count++
		}
	}
	fmt.Println(count + 1)
}