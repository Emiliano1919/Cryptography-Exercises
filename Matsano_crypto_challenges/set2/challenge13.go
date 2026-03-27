package main

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
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

func profile_for(s string) ([]byte, string, error) {
	email := regexp.QuoteMeta(s) // Escape any special character
	uid := users
	users++
	role := "user"
	stringOut := fmt.Sprintf("email=%s&uid=%d&role=%s", email, uid, role)
	dat := map[string]interface{}{
		"email": email,
		"uid":   uid,
		"role":  role,
	}
	jsonOut, err := json.Marshal(dat)
	if err != nil {
		return nil, "", err
	}
	return jsonOut, stringOut, nil
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

func encrypted_profile_for(s string) ([]byte, error) {
	_, stringOut, err := profile_for(s)
	stringOut = string(padByteToNextblockSize([]byte(stringOut), blockSize))
	if err != nil {
		return nil, err
	} else {
		return encryptECB(stableKey, []byte(stringOut)), nil
	}
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

func main() {
	test := "DOG@doggy.com"
	cipher, err := encrypted_profile_for(string(test))
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%s\n", cipher)
	plain := decryptECB(string(stableKey), cipher)
	plainPar, err := parsingRoutine(string(plain))
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%s\n", plainPar)

}
