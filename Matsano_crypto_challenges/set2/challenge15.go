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
	fmt.Printf("\nResult input Version: %q\n", paddedText)
	if multipleOf > 32 {
		return MultiplicityError
	}
	currentSize := len(paddedText)
	if currentSize%multipleOf != 0 {
		return PaddingSizeError
	}
	if bytes.ContainsAny(paddedText, UnprintableASCII) {
		first := bytes.IndexAny(paddedText, UnprintableASCII)
		withoutPadding := make([]byte, first)
		copy(withoutPadding, paddedText[:first])
		fmt.Printf("Result without padding Version: %q\n", withoutPadding)
		generatedPadding := padByteToNextMultipleOf(withoutPadding, multipleOf)
		fmt.Printf("Result Generated Version: %q\n", generatedPadding)
		fmt.Printf("Result Original Version: %q\n", paddedText)
		if bytes.Equal(generatedPadding, paddedText) {
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
	test2 := []byte("ICE ICE BABY\x04\x04\x04\x04")
	test3 := []byte("ICE ICE BABY\x05\x05\x05\x05")
	test4 := []byte("ICE ICE BABY\x01\x02\x03\x04")

	new1 := padByteVersion([]byte(test), 20)

	fmt.Printf("Result Byte Version: %q\n", new1)
	fmt.Printf("Result Byte Size: %d\n", len(new1))

	err := isValidPaddingMultipleOf([]byte(new1), 16)
	if err != nil {
		println("New1 kapput")
		log.Println(err)
	}
	new2 := padByteToNextMultipleOf([]byte(test), 16)

	fmt.Printf("Result Byte Version: %q\n", new2)
	fmt.Printf("Result Byte Size: %d\n", len(new2))
	err2 := isValidPaddingMultipleOf(new2, 16)
	if err2 != nil {
		println("New2 kapput")
		log.Println(err2)
	}

	err3 := isValidPaddingMultipleOf(test2, 16)
	if err3 != nil {
		println("Test2 kapput")
		log.Println(err3)
	}
	err4 := isValidPaddingMultipleOf(test3, 16)
	if err4 != nil {

		println("Test3 kapput")
		log.Println(err4)
	}
	err5 := isValidPaddingMultipleOf(test4, 16)
	if err5 != nil {
		println("Test4 kapput")
		log.Println(err5)
	}
}
