package main

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

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

func encrypted_profile_for(s string) ([]byte, error) {
	jsonOut, _, err := profile_for(s)
	if err != nil {
		return nil, err
	} else {
		return encryptECB(string(stableKey), jsonOut), nil
	}
}

func main() {

}
