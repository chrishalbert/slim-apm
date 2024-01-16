package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Println("Hello world")
}

func getFileBytes(fileName string) []byte {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println("error opening file", err)
	}

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("error reading file", err)
	}

	return fileBytes
}
