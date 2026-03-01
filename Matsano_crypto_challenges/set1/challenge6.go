package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func BrianKernighan(v byte) int {
	c := 0
	for v > 0 {
		v &= v - 1 // clear the least significant set bit
		c++
	}
	return c
}

func HammingDistance(a []byte, b []byte) int {
	if len(a) != len(b) {
		println("Error")
	}
	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	count := 0
	for _, t := range result {
		count += BrianKernighan(t)
	}
	return count
}

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

type Result struct {
	KSize int
	Score float64
}

func breakSingleByteXOR(cipher []byte) (byte, []byte, float64) {
	var bestKey byte
	var bestPlain []byte
	var bestScore float64

	for key := 0; key < 256; key++ {
		plain := xorWithKey(cipher, byte(key))
		score := scoreEnglish(plain)

		if score > bestScore || key == 0 {
			bestScore = score
			bestKey = byte(key)
			bestPlain = plain
		}
	}

	return bestKey, bestPlain, bestScore
}

func main() {
	test1 := `this is a test`
	test2 := `wokka wokka!!!`
	println(HammingDistance([]byte(test1), []byte(test2)))
	data, err := os.ReadFile("Challenge6.txt")
	if err != nil {
		log.Fatal(err)
	}
	bigStr := strings.TrimSpace(string(data))
	bytes, err := base64.StdEncoding.DecodeString(bigStr)
	if err != nil {
		log.Fatal(err)
	}
	results := []Result{}
	for i := 2; i < 40; i++ {
		inter := HammingDistance(bytes[0:i], bytes[i:2*i])
		final := float64(inter) / float64(i)
		results = append(results, Result{i, final})
	}

	sort.Slice(results, func(a, b int) bool {
		return results[a].Score < results[b].Score
	})

	best4 := results[0:4]
	sum := 0.0
	for _, r := range best4 {
		sum += r.Score
	}

	avg := sum / float64(len(best4))
	keySize := int(avg)
	println(avg)
	println(keySize)

	var KeysizeList [][]byte
	for j := 0; j < len(bytes); j += keySize {
		end := j + keySize
		if end > len(bytes) {
			end = len(bytes)
		}
		KeysizeList = append(KeysizeList, bytes[j:end])
	}

	blocksByKey := make([][]byte, keySize)
	for k := 0; k < len(KeysizeList); k++ {
		for t := 0; t < keySize; t++ {
			if t < len(KeysizeList[k]) { // in case last block
				blocksByKey[t] = append(blocksByKey[t], KeysizeList[k][t])
			}
		}
	}

	key := make([]byte, keySize)
	for i := range blocksByKey {
		k, _, _ := breakSingleByteXOR(blocksByKey[i])
		key[i] = k
		println(k)
	}
	plaintext := make([]byte, len(bytes))
	for i := range bytes {
		plaintext[i] = bytes[i] ^ key[i%len(key)]
	}

	fmt.Println(string(plaintext))

	// bplain := ([]byte(plain))
	// word := "ICE"
	// bword := ([]byte(word))
	// result := make([]byte, len(plain))
	// for index, b := range bplain {
	// 	letterIndex := index % 3
	// 	letter := bword[letterIndex]
	// 	result[index] = b ^ letter
	// }
	// println(hex.EncodeToString(result))
}
