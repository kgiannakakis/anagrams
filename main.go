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
	Word     string
	Anagrams []string
	Count    int
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

	for _, word := range words {

		sortedWord := SortString(word)
		_, exists := wordsMap[sortedWord]
		if exists {
			continue
		}

		anagrams := findAnagrams(word, words)
		count := len(anagrams)

		if count > 1 {
			insertItem(collection, Anagram{Word: sortedWord, Anagrams: anagrams, Count: count})
			wordsMap[sortedWord] = count
			if count > maxCount {
				maxCount = count
				fmt.Printf("Word %s: Count %d\n", word, count)
			}
		}
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

func insertItem(collection *mongo.Collection, anagram Anagram) {
	// Insert a single item
	_, err := collection.InsertOne(context.TODO(), anagram)
	if err != nil {
		log.Fatal(err)
	}
}
