package main

import (
	"encoding/hex"
)

func main() {
	plain := `Burning 'em, if you ain't quick and nimble
I go crazy when I hear a cymbal`
	bplain := ([]byte(plain))
	word := "ICE"
	bword := ([]byte(word))
	result := make([]byte, len(plain))
	for index, b := range bplain {
		letterIndex := index % 3
		letter := bword[letterIndex]
		result[index] = b ^ letter
	}
	println(hex.EncodeToString(result))
}
