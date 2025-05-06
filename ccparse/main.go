package main

import (
	"flag"
	"fmt"
	"os"
)

//TODO:Read each string item, start with ", end with ", demacate with :, repeat for value, check EOF too

func main() {
	flag.Parse()

	filePath := flag.Arg(0)

	file := fmt.Sprintf("tests/step4/%s", filePath)

	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	jsonData, err := check_validity(content)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println(jsonData)
	os.Exit(0)
}

// func isNumber(char rune) bool {
// 	return '0' <= char && char <= '9' || char == '.' || char == '-'
// }

// func isLetter(char rune) bool {
// 	return 'a' <= char && char <= 'z'
// }
