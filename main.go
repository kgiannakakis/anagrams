package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Anagram maps a word with its anagrams
type Anagram struct {
	SortedWord string
	Word       string
	Anagrams   []string
	Count      int
}

func main() {

	words, err := loadWords("el_GR.dic")
	if err != nil {
		panic(err)
	}

	client, err := connect("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	collection := client.Database("anagrams").Collection("anagrams")

	wordsMap := make(map[string]int)
	maxCount := 0

	var wordsToAdd []interface{}

	for i, word := range words {

		if i%1000 == 0 && i > 0 {
			fmt.Printf("Handling word %d\n", i)
		}

		sortedWord := SortString(word)
		_, exists := wordsMap[sortedWord]
		if exists {
			continue
		}

		anagrams := findAnagrams(word, words)
		count := len(anagrams)

		if count > 1 {
			wordsToAdd = append(wordsToAdd,
				Anagram{Word: word, SortedWord: sortedWord, Anagrams: anagrams, Count: count})
			if len(wordsToAdd) >= 100 {
				fmt.Println("Inserting 100 words")
				insertItems(collection, wordsToAdd)
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
		insertItems(collection, wordsToAdd)
		wordsToAdd = nil
	}

	err = disconnect(client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func connect(uri string) (client *mongo.Client, err error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	return
}

func disconnect(client *mongo.Client) (err error) {
	err = client.Disconnect(context.TODO())
	return
}

func insertItem(collection *mongo.Collection, item interface{}) {
	// Insert a single item
	_, err := collection.InsertOne(context.TODO(), item)
	if err != nil {
		log.Fatal(err)
	}
}

func insertItems(collection *mongo.Collection, items []interface{}) {
	// Insert a single item
	_, err := collection.InsertMany(context.TODO(), items)
	if err != nil {
		log.Fatal(err)
	}
}
