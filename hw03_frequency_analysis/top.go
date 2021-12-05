package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordCalc struct {
	word  string
	count int
}

func Top10(text string) []string {
	if len(text) == 0 {
		return []string{}
	}

	words := strings.Fields(text)
	calculated := map[string]int{}

	for _, word := range words {
		calculated[word]++
	}

	wordSlice := []WordCalc{}

	for word, count := range calculated {
		wordSlice = append(wordSlice, WordCalc{word, count})
	}

	sort.Slice(wordSlice, func(i, j int) bool {
		return wordSlice[i].count > wordSlice[j].count
	})

	topSlice := wordSlice[:10]
	result := []string{}

	sort.Slice(topSlice, func(i, j int) bool {
		if wordSlice[i].count == wordSlice[j].count {
			return wordSlice[i].word < wordSlice[j].word
		}

		return wordSlice[i].count > wordSlice[j].count
	})

	for _, wordCalc := range topSlice {
		result = append(result, wordCalc.word)
	}

	return result
}
