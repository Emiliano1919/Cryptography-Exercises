package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
)

var PaddingSizeError = errors.New("Incorrect Padding Size")
var PaddingError = errors.New("Incorrect Padding")
var MultiplicityError = errors.New("Multiplicity above 32")

const UnprintableASCII = "" +
	"\x00\x01\x02\x03\x04\x05\x06\x07" +
	"\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f" +
	"\x10\x11\x12\x13\x14\x15\x16\x17" +
	"\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f"

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
	if multipleOf > 32 {
		return MultiplicityError
	}
	currentSize := len(paddedText)
	if currentSize%multipleOf != 0 {
		return PaddingSizeError
	}
	if bytes.ContainsAny(paddedText, UnprintableASCII) {
		first := bytes.IndexAny(paddedText, UnprintableASCII)
		withoutPadding := paddedText[0:first]
		generatedPadding := padByteToNextMultipleOf(withoutPadding, multipleOf)
		if string(generatedPadding) == string(paddedText) {
			return nil
		} else {
			return PaddingError
		}

	} else {
		return nil
	}
}

func main() {
	test := "YELLOW SUBMARINE"

	new1 := padByteVersion([]byte(test), 20)

	fmt.Printf("Result Byte Version: %q\n", new1)
	fmt.Printf("Result Byte Size: %d\n", len(new1))

	err := isValidPaddingMultipleOf([]byte(new1), 16)
	if err != nil {
		log.Println(err)
	}
	new2 := padByteToNextMultipleOf([]byte(test), 16)

	fmt.Printf("Result Byte Version: %q\n", new2)
	fmt.Printf("Result Byte Size: %d\n", len(new2))
	err2 := isValidPaddingMultipleOf(new2, 16)
	if err2 != nil {
		log.Println(err2)
	}
}
