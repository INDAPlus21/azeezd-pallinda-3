package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const DataFile = "loremipsum.txt"

// Return the word frequencies of the text argument.
func WordCount(text string) map[string]int {
	text = func() string { // Removes non-alphabetical characters and converts upper to low
		var cleaned strings.Builder
		cleaned.Grow(len(text))
		for i := 0; i < len(text); i++ {
			c := text[i]
			if ('a' <= c && c <= 'z') || c == ' ' || c == '\n' {
				cleaned.WriteByte(c)
			} else if 'A' <= c && c <= 'Z' {
				cleaned.WriteByte(c + 32) // ASCII difference between upper and lower
			}
		}
		return cleaned.String()
	}()

	freqs := make(map[string]int)
	words := strings.Fields(text)

	// Update map on each occurance of a word
	for _, w := range words {
		freqs[w]++
	}
	return freqs
}

// Benchmark how long it takes to count word frequencies in text numRuns times.
//
// Return the total time elapsed.
func benchmark(text string, numRuns int) int64 {
	start := time.Now()
	for i := 0; i < numRuns; i++ {
		WordCount(text)
	}
	runtimeMillis := time.Since(start).Nanoseconds() / 1e6

	return runtimeMillis
}

// Print the results of a benchmark
func printResults(runtimeMillis int64, numRuns int) {
	fmt.Printf("amount of runs: %d\n", numRuns)
	fmt.Printf("total time: %d ms\n", runtimeMillis)
	average := float64(runtimeMillis) / float64(numRuns)
	fmt.Printf("average time/run: %.2f ms\n", average)
}

func main() {
	// read in DataFile as a string called data
	dat, _ := os.ReadFile(DataFile)
	data := string(dat)

	fmt.Printf("%#v", WordCount(string(data)))

	numRuns := 100
	runtimeMillis := benchmark(string(data), numRuns)
	printResults(runtimeMillis, numRuns)
}
