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

func inputFunction(s string) []byte {
	strings.ReplaceAll(s, ";", `";"`)
	strings.ReplaceAll(s, "=", `"="`)
	prefix := `comment1=cooking%20MCs;userdata=`
	suffix := `;comment2=%20like%20a%20pound%20of%20bacon`
	completeMsg := padByteToNextblockSize([]byte(prefix+s+suffix), blockSize)
	iv := make([]byte, 16) // Just empty (I don't think it matters much)
	cipher := encryptCBC(iv, string(stableKey), completeMsg)
	return cipher
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

func tester(s string) {
	cipher, err := encrypted_profile_for(string(s))
	if err != nil {
		log.Println(err)
	}
	printByBlocks([]byte(cipher))
	plain := decryptECB(string(stableKey), cipher)
	printByBlocks([]byte(plain))
	plainPar, err := parsingRoutine(string(plain))
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Actual Output of Tester: %s\n-----------------------\n", plainPar)
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
	// You can calculate the amount of text needed by just calculating the characters
	// But this program is more visual (For understanding)
	admingBlockGetter := "DOG@doggy.com1234567893456" // Fills the first to blocks to get just admin and padding on the last one
	admin := "admin"
	adminPadded := padByteToNextblockSize([]byte(admin), blockSize) // We need to pad in the same way
	admingBlockGetter += string(adminPadded)                        // We will get it on the 3rd block
	initialPartGetter := "perro12345.com"                           // Get .com&uid=2&role= in a single block (2nd block)
	tests := []string{admingBlockGetter, initialPartGetter}
	for i, v := range tests {
		fmt.Printf("\n----------Test %d -------------\n", i)
		tester(v)
	}
	var FullInput []byte
	initialPart, err := encrypted_profile_for(initialPartGetter)
	if err != nil {
		log.Println(err)
	}
	initialPart = initialPart[0 : 2*blockSize]
	adminPart, err := encrypted_profile_for(admingBlockGetter)
	if err != nil {
		log.Println(err)
	}
	adminPart = adminPart[2*blockSize : 3*blockSize]
	FullInput = append(FullInput, initialPart...) // This is the solution (We basically just grab each block individually assamble them and it is done)
	FullInput = append(FullInput, adminPart...)
	// This part bellow is just for testing
	DecryptedInput := decryptECB(string(stableKey), FullInput)
	FinalOutput, err := parsingRoutine(DecryptedInput)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("This is the Frankestein cipher Output:  %s\n", FinalOutput)

}
