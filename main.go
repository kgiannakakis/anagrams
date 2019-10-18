package main

import (
	"fmt"
	"log"

	"db"
	"db/mongo"
	"words"
)

// Anagram maps a word with its anagrams
type Anagram struct {
	SortedWord string
	Word       string
	Anagrams   []string
	Count      int
}

func main() {

	wordsList, err := words.LoadWords("el_GR.dic")
	if err != nil {
		log.Fatal(err)
	}

	config := mongo.Config{URI: "mongodb://localhost:27017", Database: "test", Collection: "anagrams"}

	var saver db.DataSaver
	saver = &mongo.Saver{}
	err = saver.Connect(config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected!")

	wordsMap := make(map[string]int)
	maxCount := 0

	var wordsToAdd []interface{}

	for i, word := range wordsList {

		if i%1000 == 0 && i > 0 {
			fmt.Printf("Handling word %d\n", i)
		}

		sortedWord := words.SortString(word)
		_, exists := wordsMap[sortedWord]
		if exists {
			continue
		}

		anagrams := words.FindAnagrams(word, wordsList)
		count := len(anagrams)

		if count > 1 {
			wordsToAdd = append(wordsToAdd,
				Anagram{Word: word, SortedWord: sortedWord, Anagrams: anagrams, Count: count})
			if len(wordsToAdd) >= 100 {
				fmt.Println("Inserting 100 words")
				saver.InsertMany(wordsToAdd)
				wordsToAdd = nil
			}

			wordsMap[sortedWord] = count
			if count > maxCount {
				maxCount = count
				fmt.Printf("Word %s: Count %d\n", word, count)
			}
		}
	}

	if len(wordsToAdd) > 0 {
		fmt.Printf("Inserting %d words", len(wordsToAdd))
		saver.InsertMany(wordsToAdd)
		wordsToAdd = nil
	}

	err = saver.Disconnect()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection closed.")
}
