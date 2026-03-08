package main

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"log"
	"os"
	"strings"
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
func xorBytes(a, b []byte) []byte {
	if len(a) != len(b) {
		log.Fatal("Not same length")
	}

	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result
}

// a is the in place
func xorBytesInPlace(a, b []byte) {
	if len(a) != len(b) {
		panic("Nos same length")
	}
	for i := range a {
		a[i] ^= b[i]
	}
}

func encryptCBC(iv []byte, key string, bytes []byte) []byte {
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]byte, len(bytes))
	blockSize := block.BlockSize() // 16 bytes at a time
	xorBytesInPlace(bytes[0:blockSize], iv)
	for i := 0; i < len(bytes); i += blockSize {
		if i != 0 {
			// Here we do the XOR
			xorBytesInPlace(bytes[i:i+blockSize], result[i-blockSize:i])
		}
		block.Encrypt(result[i:i+blockSize], bytes[i:i+blockSize])
	}
	return result
}

func decryptCBC(iv []byte, key string, bytes []byte) []byte {
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]byte, len(bytes))
	blockSize := block.BlockSize() // 16 bytes at a time
	for i := 0; i < len(bytes); i += blockSize {
		block.Decrypt(result[i:i+blockSize], bytes[i:i+blockSize])
		if i == 0 {
			xorBytesInPlace(result[i:i+blockSize], iv)
		} else {
			xorBytesInPlace(result[i:i+blockSize], bytes[i-blockSize:i])
		}
	}
	return result
}

func main() {
	key := "YELLOW SUBMARINE"
	test := "QUE CLASE de perro es este? Un perritoooooo :)"
	bytesTest := []byte(test)
	paddedByteTest := padByteToNextMultipleOf(bytesTest, 16) // We need to pad it to an acceptable length multiple of 16
	fmt.Printf("Result padding plaintext test: %q\n", paddedByteTest)
	encryptedTest := encryptECB(key, paddedByteTest)
	fmt.Printf("Result encryption test: %q\n", encryptedTest)
	decryptedTest := decryptECB(key, []byte(encryptedTest))
	fmt.Printf("Result decryption test: %q\n", decryptedTest)
	data, err := os.ReadFile("Challenge10.txt")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

}
