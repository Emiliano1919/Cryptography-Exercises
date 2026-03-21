package main

import (
	"bytes"
	"errors"
	"fmt"
)

const PaddingError = errors.New("Incorrect Padding")

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
func padByteToNextMultipleOf(plaintext []byte, multipleOf int) []byte {
	currentSize := len(plaintext)
	remainder := currentSize % multipleOf
	var result []byte
	if remainder == 0 {
		result = padByteVersion(plaintext, currentSize)
	} else {
		nextSize := multipleOf - remainder + currentSize
		result = padByteVersion(plaintext, nextSize)
	}
	return result
}

func isValidPaddingMultipleOf(paddedText []byte, multipleOf int) error {
	currentSize := len(paddedText)
	if currentSize%multipleOf != 0 {
		return PaddingError
	} 
	
}

func main() {
	test := "YELLOW SUBMARINE"

	new1 := padByteVersion([]byte(test), 20)

	fmt.Printf("Result Byte Version: %q\n", new1)
	new2 := padStringVersion(test, 20)

	fmt.Printf("Result String Version: %q\n", new2)
}
