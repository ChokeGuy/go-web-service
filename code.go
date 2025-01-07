package main

import "strings"

func stringMatching(words []string) []string {
	result := make([]string, 0)
	for _, value := range words {
		for _, value2 := range words {
			if value == value2 || len(value) > len(value2) {
				continue
			} else if strings.Contains(value2, value) {
				result = append(result, value)
				break
			}
		}
	}
	return result
}

func main() {

}
