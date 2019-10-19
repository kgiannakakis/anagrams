package main

import (
	"db"
	"db/mongo"
	"fmt"
	"log"
	"time"

	"words"
)

// Anagram maps a word with its anagrams
type Anagram struct {
	SortedWord string
	Anagrams   []string
	Count      int
}

func main() {

	wordsList, err := words.LoadWords("el_GR.dic")
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	wordsMap := make(map[string]Anagram)
	maxCount := 0
	var maxWord string
	for _, word := range wordsList {
		sortedWord := words.SortString(word)
		_, exists := wordsMap[sortedWord]
		if exists {
			a := wordsMap[sortedWord]
			count := a.Count + 1
			newAnagrams := append(a.Anagrams, word)
			wordsMap[sortedWord] = Anagram{SortedWord: sortedWord, Anagrams: newAnagrams, Count: count}
			if count >= maxCount {
				maxCount = count
				maxWord = word
				fmt.Printf("%s %d\n", word, count)
			}
		} else {
			wordsMap[sortedWord] = Anagram{SortedWord: sortedWord, Anagrams: []string{word}, Count: 1}
		}
	}

	anagrams := words.FindAnagrams(maxWord, wordsList)

	for _, a := range anagrams {
		fmt.Println(a)
	}

	elapsed := time.Since(start)
	fmt.Printf("Took %s\n", elapsed)

	start = time.Now()

	var aa []interface{}
	for _, v := range wordsMap {
		aa = append(aa, v)
	}

	err = storeWords(aa)
	if err != nil {
		log.Fatal(err)
	}

	elapsed = time.Since(start)
	fmt.Printf("Took %s\n", elapsed)
}

func storeWords(aa []interface{}) (err error) {
	config := mongo.Config{URI: "mongodb://localhost:27017", Database: "anagrams", Collection: "anagrams"}

	var saver db.DataSaver
	saver = &mongo.Saver{}
	err = saver.Connect(config)
	if err != nil {
		return
	}
	fmt.Println("Connected!")

	saver.InsertMany(aa)

	err = saver.Disconnect()
	if err != nil {
		return
	}
	fmt.Println("Connection closed.")

	return
}
