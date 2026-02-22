package main

import (
	"encoding/hex"
	"log"
)

func xorOn2(a string, b string) string {
	newa, err := hex.DecodeString(a)
	if err != nil {
		log.Fatal(err)
	}
	newb, err := hex.DecodeString(b)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]byte, len(newa))
	for i, _ := range newa {
		result[i] = newa[i] ^ newb[i]
	}
	return string(result)
}

func main() {
	val1 := "1c0111001f010100061a024b53535009181c"
	val2 := "686974207468652062756c6c277320657965"
	output := xorOn2(val1, val2)
	println(output)
	println(hex.EncodeToString([]byte(output)))
}
