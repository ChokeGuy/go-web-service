package main

import "strings"

func mostWordsFound(sentences []string) int {
	maxLen := 0
	for _, sentence := range sentences {
		words := strings.Split(sentence, " ")
		maxLen = max(maxLen, len(words))
	}

	return maxLen
}

func main() {
}
