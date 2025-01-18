package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

func main() {

	bytes := flag.Bool("c", false, "Count bytes")
	lines := flag.Bool("l", false, "Count lines")
	words := flag.Bool("w", false, "Count words")
	chars := flag.Bool("m", false, "Count characters")

	w := flag.Bool("word", false, "Count words")
	flag.Parse()

	noFlagSet := !(*bytes || *lines || *words || *chars)
	countLines := *lines || noFlagSet
	countWords := *words || noFlagSet
	countBytes := *bytes || noFlagSet
	countChars := *chars && !noFlagSet

	_ = countBytes
	_ = countLines
	_ = countWords
	_ = countChars

	filePath := flag.Arg(0)
	if *bytes {
		opened_file, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error:", err)
		}

		fmt.Println(len(opened_file), " ", filePath)
	}

	if *w {
		opened_file, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error:", err)
		}

		words := strings.Fields(string(opened_file))
		fmt.Println(len(words), " ", filePath)
	}

	if *lines {
		opened_file, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error:", err)
		}

		words := strings.Split(string(opened_file), "\n")
		fmt.Println(len(words), " ", filePath)
	}

	// numBytes, numLines, numWords, numChars, err := get_values(filePath, countBytes, countLines, countWords, countChars)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// } else {
	// 	display_clean_values(numBytes, numLines, numWords, numChars, filePath)
	// }

}

func get_values(filePath string, bytes, lines, words, chars bool) (numBytes, numLines, numWords, numChars int, err error) {
	var reader io.Reader

	if filePath == "" {
		reader = os.Stdin
	} else {
		opened_file, err := os.Open(filePath)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("ccwc: failed to open file: %v", err)
		}
		defer opened_file.Close()
		reader = opened_file
	}

	bufReader := bufio.NewReader(reader)
	inWord := false
	for {
		r, b, err := bufReader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, 0, 0, 0, fmt.Errorf("ccwc: failed to read: %v", err)
		}

		if bytes && numBytes < 10 {
			numBytes += b
		}

		if lines && r == '\n' {
			numLines++
		}

		if words {
			if unicode.IsSpace(r) {
				inWord = false
			} else if !inWord {
				numWords++
				inWord = true
			}
		}

		if chars {
			numChars++
		}

	}

	return numBytes, numLines, numWords, numChars, nil
}

func display_clean_values(numBytes, numLines, numWords, numChars int, filePath string) {
	result := ""
	if numLines > 0 {
		result += fmt.Sprintf("%d", numLines)
	}
	if numWords > 0 {
		result += fmt.Sprintf("%d", numWords)
	}
	if numChars > 0 {
		result += fmt.Sprintf("%d", numChars)
	} else if numBytes > 0 {
		result += fmt.Sprintf("%d", numBytes)
	}
	if filePath != "" {
		result += fmt.Sprintf(" %s", filePath)
	}

	fmt.Println(result)
}
