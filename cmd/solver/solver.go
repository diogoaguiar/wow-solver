package main

import (
	"fmt"
	"os"
)

const (
	LANG            = "pt-pt"
	MIN_WORD_LENGTH = 3
)

func main() {
	lang := LANG

	letters, err := parseInputLetters(os.Args)

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
		return nil, fmt.Errorf("Usage: %s <letters>", args[0])
	}

	letters := args[1]
	if len(letters) < MIN_WORD_LENGTH {
		return nil, fmt.Errorf("At least %d letters are required", MIN_WORD_LENGTH)
	}

	letters = normalizeWord(letters)

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
	for _, word := range words {
		word = normalizeWord(word)

		if len(word) < MIN_WORD_LENGTH {
			continue
		}

		if len(word) > len(letters) {
			continue
		}

		if isWordInLetters(word, letters) {
			filtered = append(filtered, word)
		}
	}

	return filtered
}

func normalizeWord(word string) string {
	var normalized []rune
	for _, letter := range word {
		if letter >= 'A' && letter <= 'Z' {
			letter += 32
		}

		if letter >= 'a' && letter <= 'z' {
			normalized = append(normalized, letter)
		}
	}

	return string(normalized)
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
