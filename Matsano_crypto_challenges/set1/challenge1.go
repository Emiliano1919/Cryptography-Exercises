package main

import (
	b64 "encoding/base64"
	"encoding/hex"
	"log"
)

func main() {
	cipher := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	bytes, err := hex.DecodeString(cipher)
	if err != nil {
		log.Fatal(err)
	}
	plaintext := b64.StdEncoding.EncodeToString([]byte(bytes))
	println(plaintext)
}
