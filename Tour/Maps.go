package main

import (
	"code.google.com/p/go-tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	output := make(map[string]int)
	fields := strings.Fields(s)
	for _, field := range fields {
		output[field] += 1
	}
	return output
}

func main() {
	wc.Test(WordCount)
}
