package message_generator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/artem-y/commit/internal/config"
)

type MessageGenerator struct {
	BranchName  string
	UserMessage string
	Config      config.CommitConfig
}

// Generates a commit message based on the branch name and user message
func (generator *MessageGenerator) GenerateMessage() string {
	cfg := generator.Config
	branchName := generator.BranchName
	commitMessage := generator.UserMessage

	matches := findIssueMatchesInBranch(cfg.IssueRegex, branchName)

	if len(matches) > 0 {
		commitMessage = generateCommitMessageWithMatches(matches, cfg, commitMessage)
	}

	return commitMessage
}

// Searches the branch name for issue numbers matching the given regex
func findIssueMatchesInBranch(rgxRaw string, branchName string) []string {
	rgx := regexp.MustCompile(rgxRaw)

	allSubmatches := rgx.FindAllStringSubmatch(branchName, -1)

	matches := []string{}
	for _, submatches := range allSubmatches {
		matches = append(matches, submatches[0])
	}

	return matches
}

// Generates a commit message with the issue number matches and config settings
func generateCommitMessageWithMatches(matches []string, cfg config.CommitConfig, commitMessage string) string {
	mappedMatches := make([]string, len(matches))

	for index, match := range matches {
		wrappedIssueNumber := fmt.Sprintf(
			"%s%s%s",
			cfg.OutputIssuePrefix,
			match,
			cfg.OutputIssueSuffix,
		)
		mappedMatches[index] = wrappedIssueNumber
	}

	joinedIssues := strings.Join(mappedMatches, ", ")
	return fmt.Sprintf(
		"%s%s%s%s",
		cfg.OutputStringPrefix,
		joinedIssues,
		cfg.OutputStringSuffix,
		commitMessage,
	)
}
