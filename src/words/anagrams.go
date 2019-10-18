package words

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

// LoadWords loads a list of words from a file in path locatoin
func LoadWords(path string) (words []string, err error) {
	var text *os.File
	text, err = os.Open(path)
	if err != nil {
		return
	}
	defer text.Close()

	wordsMap := make(map[string]int)
	words = make([]string, 0, 100000)
	s := bufio.NewScanner(text)
	for s.Scan() {
		word := s.Text()
		_, exists := wordsMap[word]
		if !exists {
			words = append(words, word)
			wordsMap[word] = 1
		}
	}
	return
}

// FindAnagrams find anagrams of a word. Copied from https://stackoverflow.com/a/54881186/24054
func FindAnagrams(find string, words []string) []string {
	find = strings.ToUpper(find)
	findSum := 0
	findRunes := []rune(find)
	j := 0
	for i, r := range findRunes {
		if r != ' ' {
			findSum += int(r)
			if i != j {
				findRunes[j] = r
			}
			j++
		}
	}
	findRunes = findRunes[:j]
	sort.Slice(findRunes, func(i, j int) bool { return findRunes[i] < findRunes[j] })
	findStr := string(findRunes)

	anagrams := []string{find}
	for _, word := range words {
		wordSum := 0
		wordRunes := []rune(word)
		j := 0
		for i, r := range wordRunes {
			if r != ' ' {
				wordSum += int(r)
				if i != j {
					wordRunes[j] = r
				}
				j++
			}
		}
		wordRunes = wordRunes[:j]
		if len(wordRunes) != len(findRunes) {
			continue
		}
		if wordSum != findSum {
			continue
		}
		sort.Slice(wordRunes, func(i, j int) bool { return wordRunes[i] < wordRunes[j] })
		if string(wordRunes) == findStr {
			if word != find {
				anagrams = append(anagrams, word)
			}
		}
	}
	return anagrams
}
