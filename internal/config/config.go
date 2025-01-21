package config

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/artem-y/commit/internal/helpers"
)

// Representation of settings that can be set from the config file
type CommitConfig struct {
	IssueRegex         string // Regex for issue numbers in branch name
	OutputIssuePrefix  string // Prefix before each issue number in the commit message
	OutputIssueSuffix  string // Suffix after each issue number in the commit message
	OutputStringPrefix string // Prefix before the list of issues in the commit message
	OutputStringSuffix string // Suffix after the list of issues and before the user's message
}

// DTO for unmarshalling JSON config and safe-guarding against nil values
type commitConfigDTO struct {
	IssueRegex         *string `json:"issueRegex"`
	OutputIssuePrefix  *string `json:"outputIssuePrefix"`
	OutputIssueSuffix  *string `json:"outputIssueSuffix"`
	OutputStringPrefix *string `json:"outputStringPrefix"`
	OutputStringSuffix *string `json:"outputStringSuffix"`
}

// Reads config at the file path and unmarshals it into commitConfig struct
func ReadCommitConfig(fileReader FileReading, configFilePath string) (CommitConfig, error) {
	var cfgDto commitConfigDTO

	_, err := fileReader.Stat(configFilePath)
	if err == nil {

		file, err := fileReader.ReadFile(configFilePath)
		if err != nil {
			return CommitConfig{}, err
		}

		err = json.Unmarshal(file, &cfgDto)
		if err != nil {
			return CommitConfig{}, err
		}
	}

	cfg := makeConfig(cfgDto)

	if err := validateRegex(cfg.IssueRegex); err != nil {
		return CommitConfig{}, err
	}

	return cfg, nil
}

// Helper function to create a default config
func MakeDefaultConfig() CommitConfig {
	return CommitConfig{
		IssueRegex:         helpers.DEFAULT_ISSUE_REGEX,
		OutputIssuePrefix:  helpers.DEFAULT_OUTPUT_ISSUE_PREFIX,
		OutputIssueSuffix:  helpers.DEFAULT_OUTPUT_ISSUE_SUFFIX,
		OutputStringPrefix: helpers.DEFAULT_OUTPUT_STRING_PREFIX,
		OutputStringSuffix: helpers.DEFAULT_OUTPUT_STRING_SUFFIX,
	}
}

// MARK: - Private

func makeConfig(cfgDto commitConfigDTO) CommitConfig {
	cfg := CommitConfig{}

	if cfgDto.IssueRegex != nil {
		if *cfgDto.IssueRegex == "" {
			return cfg
		}
		cfg.IssueRegex = *cfgDto.IssueRegex
	}

	if cfgDto.OutputIssuePrefix != nil {
		cfg.OutputIssuePrefix = *cfgDto.OutputIssuePrefix
	}

	if cfgDto.OutputIssueSuffix != nil {
		cfg.OutputIssueSuffix = *cfgDto.OutputIssueSuffix
	}

	if cfgDto.OutputStringPrefix != nil {
		cfg.OutputStringPrefix = *cfgDto.OutputStringPrefix
	}

	if cfgDto.OutputStringSuffix != nil {
		cfg.OutputStringSuffix = *cfgDto.OutputStringSuffix
	}

	if (cfg == CommitConfig{}) {
		return MakeDefaultConfig()
	}

	return cfg
}

// Validate the issue regex
func validateRegex(rawString string) error {
	if rawString == "" {
		return errors.New("Issue regex can't be empty. Please update the config file.")
	}

	_, err := regexp.Compile(rawString)
	return err
}
