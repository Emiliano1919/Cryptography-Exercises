package main

// This one i did not code it, just leave it here to have it as note

import (
	"encoding/hex"
	"fmt"
	"strings"
)

var freq = map[byte]float64{
	' ': 13,
	'e': 12.7, 't': 9.1, 'a': 8.2, 'o': 7.5, 'i': 7.0, 'n': 6.7,
	's': 6.3, 'h': 6.1, 'r': 6.0, 'd': 4.3, 'l': 4.0, 'u': 2.8,
}

func scoreEnglish(text []byte) float64 {
	score := 0.0
	for _, c := range text {
		c = byte(strings.ToLower(string(c))[0])
		if val, ok := freq[c]; ok {
			score += val
		}
	}
	return score
}

func xorWithKey(data []byte, key byte) []byte {
	out := make([]byte, len(data))
	for i := range data {
		out[i] = data[i] ^ key
	}
	return out
}

func main() {
	hexCipher := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"

	cipherBytes, err := hex.DecodeString(hexCipher)
	if err != nil {
		panic(err)
	}

	var bestKey byte
	var bestPlain []byte
	var bestScore float64

	for key := 0; key < 256; key++ {
		// Try all 256, because we want to try all bytes possible there are 2^8
		// So we try them all
		plain := xorWithKey(cipherBytes, byte(key))
		score := scoreEnglish(plain) // We put a score on their legibility
		if score > bestScore || key == 0 {
			bestScore = score // the one with the best score is probably the correct one
			bestKey = byte(key)
			bestPlain = plain
		}
	}

	fmt.Printf("Best key: %q (%d)\n", bestKey, bestKey)
	fmt.Printf("Decrypted: %s\n", bestPlain)
}
