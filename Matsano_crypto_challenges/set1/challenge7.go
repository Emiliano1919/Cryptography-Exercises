package main

import (
	"crypto/aes"
	"crypto/cipher"
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
	aes.NewCipher()
	cipher.Block.Decrypt()
}
