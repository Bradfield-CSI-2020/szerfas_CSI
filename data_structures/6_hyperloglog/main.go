package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"log"
	"math/bits"
	"os"
)

func loadWords(path string) ([][]byte, error) {
//func loadWords(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var result [][]byte
	//var result []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//result = append(result, scanner.Text())
		result = append(result, scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func getLeadingZeroes() int {
	// declare maxCount
	var maxCount = 0
	// read usr/share/dict/words
	words, err := loadWords("/usr/share/dict/words")
	if err != nil {
		log.Fatal("error loading words")
	}
	hash := fnv.New32()
	// loop through each word
	for _, word := range words {
		fmt.Println(word)
		//hash.Write([]byte(word))
		hash.Write(word)
		hashUint32 := hash.Sum32()
		fmt.Println(hashUint32)
		hash.Reset()
		count := bits.LeadingZeros32(hashUint32)
		if maxCount < count {
			maxCount = count
		}
		fmt.Printf("maxCount: %d\n", maxCount)
	}
	return maxCount
}

func main() {
	count := getLeadingZeroes()
	fmt.Printf("max leading zeroes in usr/share/dict/words: %d\n", count)
}
