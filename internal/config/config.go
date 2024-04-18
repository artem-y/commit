package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/artem-y/commit/internal/helpers"
)

type commitConfig struct {
	IssueRegex        string `json:"issueRegex"`
	OutputIssuePrefix string `json:"outputIssuePrefix"`
	OutputIssueSuffix string `json:"outputIssueSuffix"`
}

// Reads .commit.json file from current directory and unmarshals it into commitConfig struct
func ReadCommitConfig() commitConfig {

	configFilePath := helpers.DEFAULT_CONFIG_FILE_PATH

	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		return commitConfig{}
	}

	file, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, helpers.Red("Error reading %s file: %v\n"), err, configFilePath)
		os.Exit(1)
	}

	var cfg commitConfig
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, helpers.Red("Error unmarshalling %s file: %v\n"), err, configFilePath)
		os.Exit(1)
	}

	return cfg
}
