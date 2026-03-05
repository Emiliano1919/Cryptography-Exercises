package main

import (
	"crypto/aes"
	"encoding/base64"
	"log"
	"os"
	"strings"
)

func main() {
	data, err := os.ReadFile("Challenge7.txt")
	if err != nil {
		log.Fatal(err)
	}
	bigStr := strings.ReplaceAll(string(data), "\n", "")
	bytes, err := base64.StdEncoding.DecodeString(bigStr)
	if err != nil {
		log.Fatal(err)
	}
	key := "YELLOW SUBMARINE"
	keyBytes := []byte(key)
	println(bytes)
	println(len(keyBytes))
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]byte, len(bytes))
	blockSize := block.BlockSize() // 16 bytes at a time
	for i := 0; i < len(bytes); i += blockSize {
		block.Decrypt(result[i:i+blockSize], bytes[i:i+blockSize]) // We have to decrypt it 16 bytes at a time
	}
	println(string(result))
}
