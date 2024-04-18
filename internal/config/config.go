package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/artem-y/commit/internal/helpers"
)

type commitConfig struct {
	IssueRegex        string  `json:"issueRegex"`
	OutputIssuePrefix *string `json:"outputIssuePrefix"`
	OutputIssueSuffix *string `json:"outputIssueSuffix"`
}

// Reads config at the file path and unmarshals it into commitConfig struct
func ReadCommitConfig(configFilePath string) commitConfig {

	var cfg commitConfig

	_, err := os.Stat(configFilePath)
	if err == nil {

		file, err := os.ReadFile(configFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, helpers.Red("Error reading %s file: %v\n"), err, configFilePath)
			os.Exit(1)
		}

		err = json.Unmarshal(file, &cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, helpers.Red("Error unmarshalling %s file: %v\n"), err, configFilePath)
			os.Exit(1)
		}
	}

	if cfg.IssueRegex == "" {
		cfg.IssueRegex = helpers.DEFAULT_ISSUE_REGEX
	}

	if cfg.OutputIssuePrefix == nil {
		*cfg.OutputIssuePrefix = helpers.DEFAULT_OUTPUT_ISSUE_PREFIX
	}

	if cfg.OutputIssueSuffix == nil {
		*cfg.OutputIssueSuffix = helpers.DEFAULT_OUTPUT_ISSUE_SUFFIX
	}

	return cfg
}
