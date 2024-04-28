package config_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/artem-y/commit/internal/config"
	"github.com/artem-y/commit/internal/config/mocks"
	"github.com/artem-y/commit/internal/helpers"
)

func Test_ReadCommitConfig_WhenFileDoesNotExist_ReturnsDefaultConfig(t *testing.T) {
	// Arrange
	var mock *mocks.FileReadingMock = &mocks.FileReadingMock{}
	err := errors.New("file does not exist")
	mock.Results.Stat.Error = err
	mock.Results.ReadFile.Error = err

	defaultConfig := makeDefaultConfig()

	// Act
	cfg, err := config.ReadCommitConfig(mock, "some/path")

	// Assert
	for _, invocation := range mock.Invocations {
		if invocation == mocks.InvocationReadFile {
			t.Error("Didn't expect trying to read file when the file does not exist")
		}
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(cfg, defaultConfig) {
		t.Errorf("Expected default config, got %s", describe(cfg))
	}
}

func Test_ReadCommitConfig_WhenFilledWithValidSettings_LoadsAllValuesFromConfig(t *testing.T) {
	// Arrange
	var mock *mocks.FileReadingMock = &mocks.FileReadingMock{}
	issueRegex := "XY[0-9]+"
	outputIssuePrefix := "("
	outputIssueSuffix := ")"
	outputStringPrefix := "[ "
	outputStringSuffix := " ]"

	expectedConfig := config.CommitConfig{
		IssueRegex:         issueRegex,
		OutputIssuePrefix:  &outputIssuePrefix,
		OutputIssueSuffix:  &outputIssueSuffix,
		OutputStringPrefix: &outputStringPrefix,
		OutputStringSuffix: &outputStringSuffix,
	}
	configJson := makeJSONFromConfig(expectedConfig)

	mock.Results.ReadFile.Success = []byte(configJson)

	// Act
	cfg, err := config.ReadCommitConfig(mock, "some/path")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(cfg, expectedConfig) {
		t.Errorf(
			"Expected `%s`, got `%s`",
			describe(expectedConfig),
			describe(cfg),
		)
	}
}

func Test_ReadCommitConfig_WhenInvalidJson_ReturnsError(t *testing.T) {
	// Arrange
	var mock *mocks.FileReadingMock = &mocks.FileReadingMock{}
	mock.Results.ReadFile.Success = []byte("{invalid json}")

	// Act
	_, err := config.ReadCommitConfig(mock, "some/path")

	// Assert
	if err == nil {
		t.Error("Expected an error, got `nil`")
	}
}

func Test_ReadCommitConfig_WhenFailedToReadFile_ReturnsError(t *testing.T) {
	// Arrange
	var mock *mocks.FileReadingMock = &mocks.FileReadingMock{}
	mock.Results.ReadFile.Error = errors.New("failed to read file")

	// Act
	_, err := config.ReadCommitConfig(mock, "some/path")

	// Assert
	if err == nil {
		t.Errorf("Expected error 'failed to read file', got '%v'", err)
		return
	}
	if err.Error() != "failed to read file" {
		t.Errorf("Expected error 'failed to read file', got '%v'", err)
	}
}

func Test_ReadCommitConfig_WhenOnlyRegexInConfix_ReturnsConfigWithRegex(t *testing.T) {
	// Arrange
	var mock *mocks.FileReadingMock = &mocks.FileReadingMock{}
	expectedRegex := "ABC-[0-9]+"
	configJson := fmt.Sprintf("{\"issueRegex\":\"%s\"}", expectedRegex)
	mock.Results.ReadFile.Success = []byte(configJson)

	expectedConfig := makeDefaultConfig()
	expectedConfig.IssueRegex = expectedRegex

	// Act
	cfg, err := config.ReadCommitConfig(mock, "some/path")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if cfg.IssueRegex != expectedRegex {
		t.Errorf(
			"Expected regex ('%s') in config, got '%s'",
			expectedRegex,
			cfg.IssueRegex,
		)
	}

	if !reflect.DeepEqual(cfg, expectedConfig) {
		t.Errorf(
			"Expected config:\n'%s'\nActual config:\n'%s'",
			describe(expectedConfig),
			describe(cfg),
		)
	}
}

// Helper function to create a default config
func makeDefaultConfig() config.CommitConfig {
	outputIssuePrefix := helpers.DEFAULT_OUTPUT_ISSUE_PREFIX
	outputIssueSuffix := helpers.DEFAULT_OUTPUT_ISSUE_SUFFIX
	outputStringPrefix := helpers.DEFAULT_OUTPUT_STRING_PREFIX
	outputStringSuffix := helpers.DEFAULT_OUTPUT_STRING_SUFFIX

	return config.CommitConfig{
		IssueRegex:         helpers.DEFAULT_ISSUE_REGEX,
		OutputIssuePrefix:  &outputIssuePrefix,
		OutputIssueSuffix:  &outputIssueSuffix,
		OutputStringPrefix: &outputStringPrefix,
		OutputStringSuffix: &outputStringSuffix,
	}
}

// Helper function to make a human-readable description for a config
func describe(cfg config.CommitConfig) string {
	return fmt.Sprintf(
		`{
			IssueRegex: '%s', 
			OutputIssuePrefix: '%s', 
			OutputIssueSuffix: '%s', 
			OutputStringPrefix: '%s', 
			OutputStringSuffix: '%s'
		}`,
		cfg.IssueRegex,
		*cfg.OutputIssuePrefix,
		*cfg.OutputIssueSuffix,
		*cfg.OutputStringPrefix,
		*cfg.OutputStringSuffix,
	)
}

// Helper function to create a JSON string from a config
func makeJSONFromConfig(cfg config.CommitConfig) string {
	return fmt.Sprintf(
		`{
			"issueRegex": "%s",
			"outputIssuePrefix": "%s",
			"outputIssueSuffix": "%s",
			"outputStringPrefix": "%s",
			"outputStringSuffix": "%s"
		}`,
		cfg.IssueRegex,
		*cfg.OutputIssuePrefix,
		*cfg.OutputIssueSuffix,
		*cfg.OutputStringPrefix,
		*cfg.OutputStringSuffix,
	)
}
