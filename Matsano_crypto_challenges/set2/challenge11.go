package main

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"fmt"
	"log"
	mrand "math/rand"
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
func encryptECB(key []byte, bytes []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]byte, len(bytes))
	blockSize := block.BlockSize() // 16 bytes at a time
	for i := 0; i < len(bytes); i += blockSize {
		block.Encrypt(result[i:i+blockSize], bytes[i:i+blockSize]) // We have to encrypt it 16 bytes at a time
	}
	return result
}

func padByteVersion(plaintext []byte, size int) []byte {
	padding := size - len([]byte(plaintext))

	bpad := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(plaintext, bpad...)
}

func padByteToNextblockSize(plaintext []byte, blockSize int) []byte {
	currentSize := len(plaintext)
	remainder := currentSize % blockSize
	var result []byte
	if remainder == 0 {
		result = padByteVersion(plaintext, currentSize)
	} else {
		nextSize := blockSize - remainder + currentSize
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

func encryptCBC(iv []byte, key []byte, bytes []byte) []byte {
	block, err := aes.NewCipher(key)
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

func randomKey16Bytes() []byte {
	key := make([]byte, 16)
	rand.Read(key)
	return key
}

func encryption_oracle(plaintext []byte) []byte {
	key := randomKey16Bytes()
	rand1Size := mrand.Intn(6) + 5
	rand2Size := mrand.Intn(6) + 5
	first := make([]byte, rand1Size)
	rand.Read(first)
	last := make([]byte, rand2Size)
	rand.Read(last)
	plaintext = append(first, plaintext...)
	plaintext = append(plaintext, last...)
	plaintext = padByteToNextblockSize(plaintext, 16)
	if mrand.Intn(2) == 0 {
		println("ECB")
		return []byte(encryptECB(key, plaintext))
	} else {
		println("CBC")
		iv := make([]byte, 16)
		rand.Read(iv)
		return encryptCBC(iv, key, plaintext)
	}
}
func counterOfRepeat(cipher []byte) int {

	seen := make(map[string]struct{})
	totalBlocks := 0

	for i := 0; i+blockSize <= len(cipher); i += blockSize {
		block := string(cipher[i : i+blockSize])
		seen[block] = struct{}{}
		totalBlocks++
	}

	return totalBlocks - len(seen)
}

const blockSize = 16 // The block size for AES

func isECB(cipher []byte) {
	if counterOfRepeat(cipher) > 0 {
		println("\n It is ------ECB-------\n")
	} else {
		println("\n It is ------CBC-------\n")
	}
}
func main() {
	// This detection depends on Chosen Plaintext Attack (CPA)
	test := "QUE CLASE de perro es este? Un perritoooooo :)oooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo"
	bytesTest := []byte(test)
	fmt.Printf("Result padding plaintext test: %q\n", test)
	encryptedTest := encryption_oracle(bytesTest)
	fmt.Printf("Result encryption test1: %q\n", encryptedTest)
	isECB(encryptedTest)
	encryptedTest = encryption_oracle(bytesTest)
	fmt.Printf("Result encryption test2: %q\n", encryptedTest)
	isECB(encryptedTest)

}
