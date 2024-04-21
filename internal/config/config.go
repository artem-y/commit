package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/artem-y/commit/internal/helpers"
)

type CommitConfig struct {
	IssueRegex         string  `json:"issueRegex"`
	OutputIssuePrefix  *string `json:"outputIssuePrefix"`
	OutputIssueSuffix  *string `json:"outputIssueSuffix"`
	OutputStringPrefix *string `json:"outputStringPrefix"`
	OutputStringSuffix *string `json:"outputStringSuffix"`
}

// Reads config at the file path and unmarshals it into commitConfig struct
func ReadCommitConfig(configFilePath string) CommitConfig {

	var cfg CommitConfig

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
		defaultIssuePrefix := helpers.DEFAULT_OUTPUT_ISSUE_PREFIX
		cfg.OutputIssuePrefix = &defaultIssuePrefix
	}

	if cfg.OutputIssueSuffix == nil {
		defaultIssueSuffix := helpers.DEFAULT_OUTPUT_ISSUE_SUFFIX
		cfg.OutputIssueSuffix = &defaultIssueSuffix
	}

	if cfg.OutputStringPrefix == nil {
		defaultStringPrefix := helpers.DEFAULT_OUTPUT_STRING_PREFIX
		cfg.OutputStringPrefix = &defaultStringPrefix
	}

	if cfg.OutputStringSuffix == nil {
		defaultStringSuffix := helpers.DEFAULT_OUTPUT_STRING_SUFFIX
		cfg.OutputStringSuffix = &defaultStringSuffix
	}

	return cfg
}
