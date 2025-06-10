package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var re *regexp.Regexp = regexp.MustCompile(`^["',.!?:;]+|["',.!?:;]+$`)

type keyValue struct {
	Key   string
	Value int
}

func Top10(str string) []string {
	words := strings.Fields(str)
	wordFrequencies := make(map[string]int)

	if len(str) == 0 {
		return nil
	}

	for _, w := range words {
		word := strings.ToLower(w)
		if word == "-" {
			continue
		}
		reWord := re.ReplaceAllString(word, "")

		_, ok := wordFrequencies[reWord]
		if ok {
			wordFrequencies[reWord]++
		} else {
			wordFrequencies[reWord] = 1
		}
	}

	sortedFrequeencies := getSortedFrequeencies(wordFrequencies, 10)

	result := make([]string, 0, len(sortedFrequeencies))

	for _, pair := range sortedFrequeencies {
		result = append(result, pair.Key)
	}

	return result
}

// returns "amount" elements if available (or amount = length of result slice).
func getSortedFrequeencies(wordFrequencies map[string]int, amount int) []keyValue {
	result := make([]keyValue, 0, len(wordFrequencies))

	for key, value := range wordFrequencies {
		result = append(result, keyValue{key, value})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Value == result[j].Value {
			return result[i].Key < result[j].Key
		}
		return result[i].Value > result[j].Value
	})

	return result[:min(amount, len(result))]
}
