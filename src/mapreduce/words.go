package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const DataFile = "loremipsum.txt"

// Return the word frequencies of the text argument.
//
// Split load optimally across processor cores.
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

	freq := make(map[string]int)
	words := strings.Fields(text)
	w_amount := len(words)

	// Determine amount of work that takes an appropiate amount of time to calculate
	work_size := 9e8 / w_amount
	if work_size < 1 {
		work_size = 1
	}

	var wg sync.WaitGroup

	// Buffered channel where goroutines store their work at
	// There should be approximately words/units amount of goroutines therefore creating 1 extra size to channel becoming filled
	subfreq_chan := make(chan map[string]int, w_amount/work_size+1)

	for i, j := 0, work_size; i < w_amount; i, j = j, j+work_size {
		if j > w_amount {
			j = w_amount
		}

		wg.Add(1)
		go func(i, j int) {
			// Similar to Single Worker but doing it on a given range within the Fields
			sub_freq := make(map[string]int)
			for w := i; w < j; w++ {
				sub_freq[words[w]]++
			}
			subfreq_chan <- sub_freq // Send work to channel
			wg.Done()
		}(i, j)
	}

	wg.Wait()
	close(subfreq_chan)

	// Go through all subfrequency maps received from the goroutines and add their frequencies to the main frequency map
	for sf := range subfreq_chan {
		for k, v := range sf {
			freq[k] += v
		}
	}

	return freq
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

	numRuns := 100
	runtimeMillis := benchmark(string(data), numRuns)
	printResults(runtimeMillis, numRuns)
}
