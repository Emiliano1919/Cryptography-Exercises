package main

import (
	"bytes"
	"fmt"
)

func padByteVersion(plaintext []byte, size int) []byte {
	padding := size - len([]byte(plaintext))

	bpad := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(plaintext, bpad...)
}
func padStringVersion(plaintext string, size int) string {
	padding := size - len([]byte(plaintext))

	bpad := bytes.Repeat([]byte{byte(padding)}, padding)
	plaintext += string(bpad)

	return plaintext
}

func main() {
	test := "YELLOW SUBMARINE"

	new1 := padByteVersion([]byte(test), 20)

	fmt.Printf("Result Byte Version: %q\n", new1)
	new2 := padStringVersion(test, 20)

	fmt.Printf("Result String Version: %q\n", new2)
}
