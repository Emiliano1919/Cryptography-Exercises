package main

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"log"
)

func decryptECB(key string, bytes []byte) string {
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]byte, len(bytes))
	blockSize := block.BlockSize() // 16 bytes at a time
	for i := 0; i < len(bytes); i += blockSize {
		block.Decrypt(result[i:i+blockSize], bytes[i:i+blockSize]) // We have to Decrypt it 16 bytes at a time
	}
	return string(result)
}
func encryptECB(key string, bytes []byte) string {
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]byte, len(bytes))
	blockSize := block.BlockSize() // 16 bytes at a time
	for i := 0; i < len(bytes); i += blockSize {
		block.Encrypt(result[i:i+blockSize], bytes[i:i+blockSize]) // We have to encrypt it 16 bytes at a time
	}
	return string(result)
}

func padByteVersion(plaintext []byte, size int) []byte {
	padding := size - len([]byte(plaintext))

	bpad := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(plaintext, bpad...)
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

func main() {
	key := "YELLOW SUBMARINE"
	test := "QUE CLASE de perro es este? Un perritoooooo :)"
	bytesTest := []byte(test)
	paddedByteTest := padByteToNextMultipleOf(bytesTest, 16) // We need to pad it to an acceptable length multiple of 16
	fmt.Printf("Result Byte Version: %q\n", paddedByteTest)
	encryptedTest := encryptECB(key, paddedByteTest)
	println(encryptedTest)
	decryptedTest := decryptECB(key, []byte(encryptedTest))
	fmt.Printf("Result Byte Version: %q\n", decryptedTest)
}
