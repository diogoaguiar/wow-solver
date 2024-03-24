package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	LANG = "pt-pt"
)

func main() {
	lang := LANG

	letters, err := parseInputLetters(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}

	words, err := loadDictionary(lang)
	if err != nil {
		fmt.Println("Error loading dictionary:", err)
		return
	}

	words = filterWords(words, letters)

	for _, word := range words {
		fmt.Println(word)
	}
}

func parseInputLetters(args []string) ([]rune, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("usage: %s <letters>", args[0])
	}

	letters := args[1]

	letters = strings.ToLower(letters)

	return []rune(letters), nil
}

func loadDictionary(lang string) ([]string, error) {
	file, err := os.Open("dicts/" + lang)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	for {
		var word string
		_, err := fmt.Fscanln(file, &word)
		if err != nil {
			break
		}
		words = append(words, word)
	}

	return words, nil
}

func filterWords(words []string, letters []rune) []string {
	var filtered []string
	lettersLen := len(letters)
	for _, word := range words {
		if len(word) > lettersLen {
			continue
		}

		if isWordInLetters(word, letters) {
			filtered = append(filtered, word)
		}
	}

	return filtered
}

func isWordInLetters(word string, letters []rune) bool {
	available := make(map[rune]int)
	for _, letter := range letters {
		available[letter]++
	}

	for _, letter := range word {
		if available[letter] == 0 {
			return false
		}
		available[letter]--
	}

	return true
}
