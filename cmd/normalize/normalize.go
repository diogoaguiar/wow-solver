package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	ORIG_DICTS_LOC  = "dicts/raw/"
	DICTS_LOC       = "dicts/"
	MAX_WORD_LENGTH = 10
	MIN_WORD_LENGTH = 3
)

func main() {
	dictLoc := os.Args[1]

	dict := loadDictionary(dictLoc)

	dict = toLowerCase(dict)
	dict = convertSpecialCharacters(dict)
	dict = removeWordsWithHyphen(dict)
	dict = removeWordsWithNumbers(dict)
	dict = removeDuplicates(dict)
	dict = removeWordLongerThan(dict, MAX_WORD_LENGTH)
	dict = removeWordShorterThan(dict, MIN_WORD_LENGTH)
	dict = sortWords(dict)

	saveDictionary(dictLoc, dict)
}

func loadDictionary(dictLoc string) []string {
	file, err := os.Open(ORIG_DICTS_LOC + dictLoc)
	if err != nil {
		fmt.Println("Error loading dictionary:", err)
		os.Exit(1)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return words
}

func toLowerCase(words []string) []string {
	var lowerWords []string
	for _, word := range words {
		lowerWords = append(lowerWords, strings.ToLower(word))
	}
	return lowerWords
}

func convertSpecialCharacters(words []string) []string {
	var convertedWords []string
	for _, word := range words {
		convertedWords = append(convertedWords, normalizeWord(word))
	}
	return convertedWords
}

func normalizeWord(word string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, word)
	return result
}

func removeWordsWithHyphen(words []string) []string {
	var filtered []string
	for _, word := range words {
		if !strings.Contains(word, "-") {
			filtered = append(filtered, word)
		}
	}
	return filtered
}

func removeWordsWithNumbers(words []string) []string {
	var filtered []string
	for _, word := range words {
		if !containsNumber(word) {
			filtered = append(filtered, word)
		}
	}
	return filtered
}

func containsNumber(word string) bool {
	for _, r := range word {
		if unicode.IsNumber(r) {
			return true
		}
	}
	return false
}

func removeDuplicates(words []string) []string {
	var uniqueWords []string
	unique := make(map[string]bool)
	for _, word := range words {
		if !unique[word] {
			uniqueWords = append(uniqueWords, word)
			unique[word] = true
		}
	}
	return uniqueWords
}

func removeWordLongerThan(words []string, length int) []string {
	var filtered []string
	for _, word := range words {
		if len(word) <= length {
			filtered = append(filtered, word)
		}
	}
	return filtered
}

func removeWordShorterThan(words []string, length int) []string {
	var filtered []string
	for _, word := range words {
		if len(word) >= length {
			filtered = append(filtered, word)
		}
	}
	return filtered
}

func sortWords(words []string) []string {
	slices.Sort(words)
	return words
}

func saveDictionary(dictLoc string, words []string) {
	file, err := os.Create(DICTS_LOC + dictLoc)
	if err != nil {
		fmt.Println("Error saving dictionary:", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, word := range words {
		fmt.Fprintln(writer, word)
	}
	writer.Flush()
}
