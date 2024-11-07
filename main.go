package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func createUnigramIndex() {
	unigramFile, err := os.Create("unigram_index.txt")
	if err != nil {
		fmt.Println("Error creating unigram index file:", err)
		return
	}
	defer unigramFile.Close()
	unigramMap := make(map[string]map[string]int)
	err = filepath.Walk("./data/fulldata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			docID := strings.TrimSuffix(info.Name(), ".txt")
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				words := preprocessText(line)

				for _, word := range words {
					if _, exists := unigramMap[word]; !exists {
						unigramMap[word] = make(map[string]int)
					}
					unigramMap[word][docID]++
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking through files:", err)
		return
	}
	for word, docMap := range unigramMap {
		var docIDCounts []string
		for docID, count := range docMap {
			docIDCounts = append(docIDCounts, fmt.Sprintf("%s:%d", docID, count))
		}
		unigramFile.WriteString(fmt.Sprintf("%s\t%s\n", word, strings.Join(docIDCounts, " ")))
	}
}

func createBigramIndex() {
	bigramFile, err := os.Create("selected_bigram_index.txt")
	if err != nil {
		fmt.Println("Error creating bigram index file:", err)
		return
	}
	defer bigramFile.Close()
	bigramTargets := []string{
		"computer science",
		"information retrieval",
		"power politics",
		"los angeles",
		"bruce willis",
	}
	bigramMap := make(map[string]map[string]int)
	err = filepath.Walk("./data/devdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			docID := strings.TrimSuffix(info.Name(), ".txt")
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				words := preprocessText(line)
				for i := 0; i < len(words)-1; i++ {
					bigram := words[i] + " " + words[i+1]
					if _, exists := bigramMap[bigram]; !exists {
						bigramMap[bigram] = make(map[string]int)
					}
					bigramMap[bigram][docID]++
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking through files:", err)
		return
	}
	for _, target := range bigramTargets {
		if docMap, exists := bigramMap[target]; exists {
			var docIDCounts []string
			for docID, count := range docMap {
				docIDCounts = append(docIDCounts, fmt.Sprintf("%s:%d", docID, count))
				fmt.Println("docID 1: ", docID, "Count: ", count)
			}
			fmt.Println("Bigram: ", docIDCounts)
			bigramFile.WriteString(fmt.Sprintf("%s\t%s\n", target, strings.Join(docIDCounts, " ")))
		}
	}
}

func preprocessText(text string) []string {
	re := regexp.MustCompile(`[^\w\s]`)
	text = re.ReplaceAllString(text, " ")
	text = strings.ToLower(text)
	words := strings.Fields(text)
	return words
}

func main() {
	createUnigramIndex()
	createBigramIndex()
	fmt.Println("Unigram and bigram indexing completed.")
}
