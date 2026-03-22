package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
)

var ErrPaddingSize = errors.New("Incorrect Padding Size")
var ErrPadding = errors.New("Incorrect Padding")
var ErrVoid = errors.New("Incorrect empty text")

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
func pad(plaintext []byte, blockSize int) []byte {
	cSize := len(plaintext)
	remainder := cSize % blockSize
	var result []byte
	if remainder == 0 { //Edge case add a whole other block
		// Useful info: https://crypto.stackexchange.com/questions/78187/how-do-i-find-out-whether-a-message-in-cbc-has-padding/80481#80481
		result = padByteVersion(plaintext, cSize+blockSize)
	} else {
		nextSize := blockSize - remainder + cSize
		result = padByteVersion(plaintext, nextSize)
	}
	return result
}

// func isValidPadding(paddedText []byte, blockSize int) error {
// 	fmt.Printf("\nResult input Version: %q\n", paddedText)
// 	if blockSize > 32 {
// 		return MultiplicityError
// 	}
// 	cSize := len(paddedText)
// 	if cSize%blockSize != 0 {
// 		return PaddingSizeError
// 	}
// 	if bytes.ContainsAny(paddedText, UnprintableASCII) {
// 		first := bytes.IndexAny(paddedText, UnprintableASCII)
// 		withoutPadding := make([]byte, first)
// 		copy(withoutPadding, paddedText[:first])
// 		fmt.Printf("Result without padding Version: %q\n", withoutPadding)
// 		generatedPadding := pad(withoutPadding, blockSize)
// 		fmt.Printf("Result Generated Version: %q\n", generatedPadding)
// 		fmt.Printf("Result Original Version: %q\n", paddedText)
// 		if bytes.Equal(generatedPadding, paddedText) {
// 			return nil
// 		} else {
// 			return PaddingError
// 		}

// 	} else {
// 		return nil
// 	}
// }

func isValidPadding(paddedText []byte, blockSize int) error {
	fmt.Printf("\nResult input Version: %q\n", paddedText)
	cSize := len(paddedText)
	paddingSize := int(paddedText[len(paddedText)-1])
	if len(paddedText) == 0 || paddingSize == 0 {
		return ErrVoid
	}
	if cSize%blockSize != 0 || (paddingSize > blockSize) {
		return ErrPaddingSize
	}

	for i := 0; i < paddingSize; i++ {
		if paddedText[cSize-i-1] != byte(paddingSize) {
			return ErrPadding
		}
	}
	return nil

}

func main() {
	test := "YELLOW SUBMARINE"
	fmt.Printf("Result Byte Size: %d\n", len(test))
	test2 := []byte("ICE ICE BABY\x04\x04\x04\x04")
	test3 := []byte("ICE ICE BABY\x05\x05\x05\x05")
	test4 := []byte("ICE ICE BABY\x01\x02\x03\x04")

	new1 := padByteVersion([]byte(test), 20)

	fmt.Printf("Result Byte Version: %q\n", new1)
	fmt.Printf("Result Byte Size: %d\n", len(new1))

	err := isValidPadding([]byte(new1), 16)
	if err != nil {
		println("New1 kapput")
		log.Println(err)
	}
	new2 := pad([]byte(test), 16)

	fmt.Printf("Result Byte Version: %q\n", new2)
	fmt.Printf("Result Byte Size: %d\n", len(new2))
	err2 := isValidPadding(new2, 16)
	if err2 != nil {
		println("New2 kapput")
		log.Println(err2)
	}

	err3 := isValidPadding(test2, 16)
	if err3 != nil {
		println("Test2 kapput")
		log.Println(err3)
	}
	err4 := isValidPadding(test3, 16)
	if err4 != nil {

		println("Test3 kapput")
		log.Println(err4)
	}
	err5 := isValidPadding(test4, 16)
	if err5 != nil {
		println("Test4 kapput")
		log.Println(err5)
	}
}
