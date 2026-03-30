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
var iv = make([]byte, 16)

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

func inputFunction(s string) []byte {
	s = strings.ReplaceAll(s, ";", "\\;")
	s = strings.ReplaceAll(s, "=", "\\=")
	prefix := `comment1=cooking%20MCs;userdata=`
	suffix := `;comment2=%20like%20a%20pound%20of%20bacon`
	completeMsg := padByteToNextblockSize([]byte(prefix+s+suffix), blockSize)
	printByBlocks(completeMsg)

	cipher := encryptCBC(iv, string(stableKey), completeMsg)
	printByBlocks(cipher)
	return cipher
}

func validateAnswer(bytes []byte) bool {
	plaintext := string(decryptCBC(iv, string(stableKey), bytes))
	printByBlocks([]byte(plaintext))
	if strings.Contains(plaintext, ";admin=true;") {
		return true
	} else {
		return false
	}

}

func init() {
	stableKey = make([]byte, 16)
	rand.Read(stableKey)
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
	// It was nice being able to finally solve this one.
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Very helpful ( I need to improve how I Wrap and display errors)
	output := inputFunction(";admin.true")       // We make our work easy by chosing the plaintext and getting the resulting cipher
	printByBlocks(output)
	// We calculate where the point is and modify its equivalent position in the previous block
	// Since this modification will be passed directly (without the diffusion because it has not been decrypted yet) to XOR with the decrypted text of our target block
	// We are targeting the . which is 00101110 and we need to get = which is 00111101
	// We are therefore calculating	00111101 = x XOR 00101110
	// So we need to XOR the previous block at equivalent position with 00010011 to get the = that we want
	output[blockSize+7] ^= 0b00010011
	printByBlocks(output)
	if validateAnswer(output) {
		fmt.Printf("\n--------Success--------\n")
	} else {
		fmt.Printf("\n-------Failure---------\n")
	}
}
