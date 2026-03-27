package main

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
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

func randomKey16Bytes() []byte {
	key := make([]byte, 16)
	rand.Read(key)
	return key
}

func encryption_oracle_ECB(plaintext []byte) []byte {
	var lastB64 []byte
	const lastPlain = `Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkg
aGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBq
dXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUg
YnkK`
	lastB64, _ = base64.StdEncoding.DecodeString(lastPlain)
	lastBytes := []byte(lastB64)
	plaintext = append(plaintext, lastBytes...)
	plaintext = padByteToNextblockSize(plaintext, 16)
	return []byte(encryptECB(stableKey, plaintext))
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

func isECB(cipher []byte) {
	if counterOfRepeat(cipher) > 0 {
		println("\n It is ------ECB-------\n")
	}
}

const blockSize = 16 // The block size for AES

var stableKey []byte

func init() {
	stableKey = make([]byte, 16)
	rand.Read(stableKey)
}
func buildDictionary(base string, blockStartIndex int) map[string]string {
	dictionary := make(map[string]string)
	// Full 256 ASCII
	for l := 0; l < 256; l++ {
		lbase := base + string(l)
		//println(lbase)
		out := encryption_oracle_ECB([]byte(lbase))
		block := string(out[blockStartIndex : blockSize+blockStartIndex])
		dictionary[block] = string(l)
	}
	return dictionary
}

func byteByByteDecryption() {
	var decryptedLetters []byte
	size := encryption_oracle_ECB(nil)
	for decryptedBytes := 0; decryptedBytes < len(size); decryptedBytes++ {
		baseSize := (blockSize - 1) - (decryptedBytes % blockSize)
		blockNumber := decryptedBytes / blockSize
		blockIndex := blockNumber * blockSize
		base := bytes.Repeat([]byte("A"), baseSize)    // This remains unmuted to decrypt
		baseDic := bytes.Repeat([]byte("A"), baseSize) // This one needs the decrypted bytes to build the dictionary
		baseDic = append(baseDic, decryptedLetters...)
		currentDictionary := buildDictionary(string(baseDic), blockIndex)
		currentOut := encryption_oracle_ECB(base)
		target := string(currentOut[blockIndex : blockIndex+blockSize])
		decryptedByte := currentDictionary[target]
		decryptedLetters = append(decryptedLetters, []byte(decryptedByte)...)
		//fmt.Printf("THE LETTER IS: %s \n", decryptedByte)
	}
	fmt.Printf("The entire text is:\n-------------\n%s\n-------------\n", string(decryptedLetters))

}

func main() {
	// This detection depends on Chosen Plaintext Attack (CPA)
	var test string
	for i := 0; i < 17; i++ {
		test += "A"
		bytesTest := []byte(test)
		encryptedTest := encryption_oracle_ECB(bytesTest)
		isECB(encryptedTest)
		println(i)
		println(len(bytesTest))
		fmt.Printf("Result encryption test1: \n%q\n", encryptedTest)
	}
	// The decryption
	println("-----------------Let's try this--------------------\n")
	byteByByteDecryption()

}
