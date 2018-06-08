package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]map[string]int)
	var fileNames []string

	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, "os.Stdin", counts)
	} else {
		for _, filename := range files {
			if contains(fileNames, filename) {
				continue
			}
			fileNames = append(fileNames, filename)

			f, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup: %v\n", err)
				continue
			}
			countLines(f, filename, counts)
			f.Close()
		}
	}

	for line, filenames := range counts {
		fileCount := len(filenames)

		if fileCount == 1 {
			total := 0
			for _, count := range filenames {
				total += count
			}

			if total <= 1 {
				continue
			}
		}

		fmt.Printf("[Found duplicates in %d file(s)]:\t%s\n", fileCount, line)
		for name, count := range filenames {
			fmt.Printf("\t%d hit(s) in %s\n", count, name)
		}
	}
}

func contains(fileNames []string, filename string) bool {
	for _, value := range fileNames {
		if value == filename {
			return true
		}
	}
	return false
}

func countLines(f *os.File, filename string, counts map[string]map[string]int) {
	input := bufio.NewScanner(f)

	for input.Scan() {
		line := input.Text()

		if counts[line] == nil {
			counts[line] = make(map[string]int)
		}

		counts[line][filename]++
	}
}
