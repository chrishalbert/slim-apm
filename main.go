// main uses the slimapm objects to analyze healthchecks and report on metrics
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Healthchecks struct has a version and an associated slim metric
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
		return
	}

	oms := NewSlimApp()
	for _, healthcheck := range healthchecks {
		err = oms.AddVersionMetric(healthcheck.Version, healthcheck.SlimMetric)
		if err != nil {
			fmt.Printf("ERROR ADDING %v: %v", healthcheck.Version, err)
		}
	}

	fmt.Println("\nDeliverable 1: Aggregates By Version")
	fmt.Println("-------------------------------------------------")
	for _, version := range oms.GetVersions() {
		fmt.Println(version)
	}

	fmt.Println("\nDeliverable 2: Release Overview")
	fmt.Println("------------------------------------------------------")
	fmt.Println(oms)

	fmt.Println("\nDeliverable 3: Release History")
	fmt.Println("------------------------------------------------------")
	for _, version := range *oms.GetReleaseHistory() {
		fmt.Printf("%v => %v\n", time.Unix(int64(version.start), 0), version.hash)
	}
}

func getFileBytes(fileName string) (*[]byte, error) {
	f, err := os.Open(filepath.Clean(fileName))
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
