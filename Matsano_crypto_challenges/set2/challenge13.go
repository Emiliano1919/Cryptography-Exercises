package main

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

var users int = 0
var stableKey []byte

const blockSize int = 16

func parsingRoutine(s string) ([]byte, error) {
	dat := make(map[string]interface{}) // Initialize to be able to address
	pairs := strings.Split(s, "&")
	for _, p := range pairs {
		field := strings.SplitN(p, "=", 2)
		if len(field) != 2 {
			return nil, fmt.Errorf("invalid pair: %s", p)
		}
		dat[field[0]] = field[1]
	}
	jsonOut, err := json.Marshal(dat)
	if err != nil {
		return nil, err
	}
	return jsonOut, nil
}


func init() {
	stableKey = make([]byte, 16)
	rand.Read(stableKey)
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



func printByBlocks(bytes []byte) {
	println("\n------ START --------\n")
	for i := 0; i < len(bytes); i += blockSize {
		fmt.Printf("%q | ", bytes[i:i+blockSize])
	}
	println("\n------ END --------\n")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Very helpful ( I need to improve how I Wrap and display errors)

	fmt.Printf("This is the Frankestein cipher Output:  %s\n", FinalOutput)

}
