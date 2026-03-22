package main

import (
	"bytes"
	"fmt"
)

func padByteVersion(plaintext []byte, size int) []byte {
	padding := size - len(plaintext)

	bpad := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(plaintext, bpad...)
}
func padStringVersion(plaintext string, size int) string {
	padding := size - len(plaintext)

	bpad := bytes.Repeat([]byte{byte(padding)}, padding)
	plaintext += string(bpad)

	return plaintext
}
func padByteToNext(plaintext []byte, blockSize int) []byte {
	cSize := len(plaintext)
	remainder := cSize % blockSize
	var result []byte
	if remainder == 0 { //Edge case add a whole other block
		result = padByteVersion(plaintext, cSize+blockSize)
	} else {
		nextSize := blockSize - remainder + cSize
		result = padByteVersion(plaintext, nextSize)
	}
	return result
}

func main() {
	test := "YELLOW SUBMARINE"

	new1 := padByteVersion([]byte(test), 20)

	fmt.Printf("Result Byte Version: %q\n", new1)
	new2 := padStringVersion(test, 20)

	fmt.Printf("Result String Version: %q\n", new2)
}
