package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Healthchecks struct {
	Version string `json:"version"`
	SlimMetric
}

func main() {

	fileBytes, err := getFileBytes("./events.json")
	if err != nil {
		fmt.Printf("FAILED TO OPEN LOG FILE: %v", err)
		return
	}

	var healthchecks []Healthchecks
	if err = json.Unmarshal(*fileBytes, &healthchecks); err != nil {
		fmt.Printf("FAILED TO PARSE JSON: %v", err)
	}

	oms := NewSlimApp()
	for _, healthcheck := range healthchecks {
		oms.AddVersionMetric(healthcheck.Version, healthcheck.SlimMetric)
	}

}

func getFileBytes(fileName string) (*[]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println("error opening file", err)
		return nil, err
	}

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("error reading file", err)
		return nil, err
	}

	return &fileBytes, nil
}