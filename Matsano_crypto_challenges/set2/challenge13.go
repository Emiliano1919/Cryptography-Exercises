package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

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

func main() {

}
