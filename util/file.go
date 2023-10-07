package util

import (
	"io"
	"log"
	"os"
	"strings"
)

func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		return false
	}
	return true
}

func ReadFile(f string) *string {
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	fileContent := strings.Trim(string(b), "\r\n")
	return &fileContent
}
