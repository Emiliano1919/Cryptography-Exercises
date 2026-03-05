package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
)

const blockSize = 16 // The block size for AES

func counterOfRepeat(cipher []byte) int {

	seen := make(map[string]struct{})
	totalBlocks := 0

	for i := 0; i+blockSize <= len(cipher); i += blockSize {
		block := string(cipher[i : i+blockSize])
		seen[block] = struct{}{}
		totalBlocks++
	}

	return totalBlocks - len(seen)
}

func main() {
	data, err := os.ReadFile("Challenge8.txt")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	maxNumberOfRepeats := 0
	bestIndex := -1
	for j, line := range lines {
		bytes, err := hex.DecodeString(line)
		if err != nil {
			log.Fatal(err)
		}
		current := counterOfRepeat(bytes)
		if maxNumberOfRepeats < current {
			bestIndex = j
			maxNumberOfRepeats = current
		}
		// This are pieces of my original solution, it worked but not really good for future reference
		// set := make(map[[]byte]struct{}) sad does not work because not comparable, cant use []byte
		// set := make(map[string]struct{})

		// for i := 0; i < len(bytes); i += blockSize {
		// 	block := string(bytes[i : i+blockSize])

		// 	if _, ok := set[block]; ok {
		// 		println(j)
		// 		println("bingo")
		// 		println(block)
		// 	} else {
		// 		set[block] = struct{}{}
		// 	}

		// }
	}
	fmt.Printf("The index with the most repeats is : %d and the number of repeats is: %d \n", bestIndex, maxNumberOfRepeats)
}
