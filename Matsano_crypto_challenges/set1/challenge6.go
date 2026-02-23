package main

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

func main() {
	test1 := `this is a test`
	test2 := `wokka wokka!!!`
	println(HammingDistance([]byte(test1), []byte(test2)))

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
