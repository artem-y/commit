package message_generator_test

import (
	"testing"

	"github.com/artem-y/commit/internal/config"
	"github.com/artem-y/commit/internal/message_generator"
)

func Test_Generate_WhenNoMatchesInBranchName_ReturnsUserCommitMessageUnchanged(t *testing.T) {
	// Arrange
	outputIssuePrefix := "("
	outputIssueSuffix := ")"
	outputStringPrefix := "[ "
	outputStringSuffix := " ]"

	expectedMessage := "Test commit message"

	sut := message_generator.MessageGenerator{
		BranchName:  "main",
		UserMessage: expectedMessage,
		Config: config.CommitConfig{
			IssueRegex:         "XY[0-9]+",
			OutputIssuePrefix:  &outputIssuePrefix,
			OutputIssueSuffix:  &outputIssueSuffix,
			OutputStringPrefix: &outputStringPrefix,
			OutputStringSuffix: &outputStringSuffix,
		},
	}

	// Act
	message := sut.GenerateMessage()

	// Assert
	if message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, message)
	}
}

func Test_Generate_WhenSingleMatchFound_AddsIssueToTheMessage(t *testing.T) {
	// Arrange
	outputIssuePrefix := "("
	outputIssueSuffix := ")"
	outputStringPrefix := ""
	outputStringSuffix := " "

	userMessage := "Add validation service"
	expectedMessage := "(CD-13) Add validation service"

	sut := message_generator.MessageGenerator{
		BranchName:  "feature/CD-13-implement-login-screen-validation",
		UserMessage: userMessage,
		Config: config.CommitConfig{
			IssueRegex:         "CD-[0-9]+",
			OutputIssuePrefix:  &outputIssuePrefix,
			OutputIssueSuffix:  &outputIssueSuffix,
			OutputStringPrefix: &outputStringPrefix,
			OutputStringSuffix: &outputStringSuffix,
		},
	}

	// Act
	message := sut.GenerateMessage()

	// Assert
	if message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, message)
	}
}

func Test_Generate_WhenFoundMultipleMatches_AddsCommaSeparatedIssuesToTheMessage(t *testing.T) {
	// Arrange
	outputIssuePrefix := "#"
	outputIssueSuffix := ""
	outputStringPrefix := "["
	outputStringSuffix := "]: "

	userMessage := "Prepare mocks for core unit tests"
	expectedMessage := "[#27, #30]: Prepare mocks for core unit tests"

	sut := message_generator.MessageGenerator{
		BranchName:  "add-unit-tests-for-issues-27-and-30",
		UserMessage: userMessage,
		Config: config.CommitConfig{
			IssueRegex:         "[0-9]+",
			OutputIssuePrefix:  &outputIssuePrefix,
			OutputIssueSuffix:  &outputIssueSuffix,
			OutputStringPrefix: &outputStringPrefix,
			OutputStringSuffix: &outputStringSuffix,
		},
	}

	// Act
	message := sut.GenerateMessage()

	// Assert
	if message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, message)
	}
}

func Test_Generate_WhenAllPrefixesAndSuffixesEmpty_AddsIssueWithoutWrapping(t *testing.T) {
	// Arrange
	outputIssuePrefix := ""
	outputIssueSuffix := ""
	outputStringPrefix := ""
	outputStringSuffix := ""

	userMessage := "chore: regenerate localisation files"
	expectedMessage := "#210chore: regenerate localisation files"

	sut := message_generator.MessageGenerator{
		BranchName:  "#210-implement-login-screen-validation",
		UserMessage: userMessage,
		Config: config.CommitConfig{
			IssueRegex:         "(#)?[0-9]+",
			OutputIssuePrefix:  &outputIssuePrefix,
			OutputIssueSuffix:  &outputIssueSuffix,
			OutputStringPrefix: &outputStringPrefix,
			OutputStringSuffix: &outputStringSuffix,
		},
	}

	// Act
	message := sut.GenerateMessage()

	// Assert
	if message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, message)
	}
}
