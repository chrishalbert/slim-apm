package main

import (
	"fmt"
	"io/ioutil"
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

	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("error reading file", err)
	}

	return fileBytes
}
