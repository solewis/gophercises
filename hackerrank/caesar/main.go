package main

import (
	"fmt"
	"strings"
)

func main() {
	var k int
	fmt.Scanln(&k)

	var str string
	fmt.Scanln(&str)

	fmt.Printf("Running caesar cipher with offset %d on string %s\n", k, str)

	var ciphered strings.Builder

	for _, c := range str {
		if c >= 'A' && c <= 'Z' {
			cipheredRune := ((c - 'A' + int32(k)) % 26) + 'A'
			ciphered.WriteRune(cipheredRune)
		} else if c >= 'a' && c <= 'z' {
			cipheredRune := ((c - 'a' + int32(k)) % 26) + 'a'
			ciphered.WriteRune(cipheredRune)
		} else {
			ciphered.WriteRune(c)
		}
	}

	fmt.Println(ciphered.String())
}
