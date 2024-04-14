package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type commitConfig struct {
	IssueRegex string `json:"issueRegex"`
}

// Reads .commit.json file from current directory and unmarshals it into commitConfig struct
func readCommitConfig() commitConfig {

	configFileName := ".commit.json"

	_, err := os.Stat(configFileName)
	if os.IsNotExist(err) {
		return commitConfig{}
	}

	file, err := os.ReadFile(configFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, red("Error reading %s file: %v\n"), err, configFileName)
		os.Exit(1)
	}

	var cfg commitConfig
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, red("Error unmarshalling %s file: %v\n"), err, configFileName)
		os.Exit(1)
	}

	return cfg
}
