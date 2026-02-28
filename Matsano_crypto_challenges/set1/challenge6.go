package main

import (
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

type Result struct {
	KSize int
	Score float64
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
	bytes := []byte(bigStr)
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
	println(int(avg))

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
