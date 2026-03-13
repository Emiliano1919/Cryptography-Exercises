package main

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
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

func randomKey16Bytes() []byte {
	key := make([]byte, 16)
	rand.Read(key)
	return key
}

func encryption_oracle_ECB(plaintext []byte) []byte {
	prefix := make([]byte, prefixSize)
	rand.Read(prefix)
	plaintext = append(prefix, plaintext...)
	lastBytes := []byte(lastB64)
	plaintext = append(plaintext, lastBytes...)
	plaintext = padByteToNextMultipleOf(plaintext, 16)
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

const blockSize = 16 // The block size for AES (We find the size using techniques from previous challenges but I skip it for this problem)
var lastB64 []byte
var prefixSize int
var stableKey []byte
var alphabet = []byte{
	'A', 'B', 'C', 'D', 'E', 'F', 'G',
	'H', 'I', 'J', 'K', 'L', 'M', 'N',
	'O', 'P', 'Q', 'R', 'S', 'T',
	'U', 'V', 'W', 'X', 'Y', 'Z',
	'a', 'b', 'c', 'd', 'e', 'f', 'g',
	'h', 'i', 'j', 'k', 'l', 'm', 'n',
	'o', 'p', 'q', 'r', 's', 't',
	'u', 'v', 'w', 'x', 'y', 'z',
}

const lastPlain = `Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkg
aGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBq
dXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUg
YnkK`

func init() {
	stableKey = make([]byte, 16)
	rand.Read(stableKey)
	lastB64, _ = base64.StdEncoding.DecodeString(lastPlain)
	prefixSize = mrand.Intn(128)
}

func main() {
	// This detection depends on Chosen Plaintext Attack (CPA)
	twoBlock := make([]byte, 48) //48 bytes
	output := encryption_oracle_ECB([]byte(twoBlock))
	set := make(map[string]int)
	var initialIndex int
	for i := 0; i < len(output); i += blockSize {
		current := string(output[i : i+blockSize])
		println(i)
		println(current)
		if initial, exists := set[current]; exists {
			fmt.Println("ECB mode here")
			fmt.Printf("Starting index is: %d \n", initial)
			initialIndex = initial
		} else {
			set[current] = i
		}
	}

	base := string(make([]byte, 15)) //15 so 15 bytes as base to build dictionary
	dictionary := make(map[string]string)

	for _, l := range alphabet {
		lbase := base + string(l)
		out := encryption_oracle_ECB([]byte(lbase))
		block := string(out[initialIndex : initialIndex+blockSize])
		dictionary[block] = string(l)
	}

	println("\n-------Same result?----------\n")
	for _, x := range lastB64 {
		xbase := base + string(x)
		out := encryption_oracle_ECB([]byte(xbase))
		fmt.Print(dictionary[string(out[initialIndex:initialIndex+blockSize])])
	}

}
