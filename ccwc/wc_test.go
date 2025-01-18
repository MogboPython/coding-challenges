package main

import (
	"os"
	"testing"
)

func TestGetFileSize(t *testing.T) {
	testFilePath := "test2.txt"
	os.Args = []string{"main", "-c", testFilePath}
	main()
	os.Args = os.Args[:1]
}
