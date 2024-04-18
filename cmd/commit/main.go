package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/artem-y/commit/internal/config"
	"github.com/artem-y/commit/internal/helpers"
	"github.com/artem-y/commit/internal/user"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {
	commitMessage := getCommitMessage()

	issueRegex := helpers.DEFAULT_ISSUE_REGEX
	outputIssuePrefix := helpers.DEFAULT_OUTPUT_ISSUE_PREFIX
	outputIssueSuffix := helpers.DEFAULT_OUTPUT_ISSUE_SUFFIX

	commitCfg := config.ReadCommitConfig()
	if commitCfg.IssueRegex != "" {
		issueRegex = commitCfg.IssueRegex
	}

	if commitCfg.OutputIssuePrefix != "" {
		outputIssuePrefix = commitCfg.OutputIssuePrefix
	}

	if commitCfg.OutputIssueSuffix != "" {
		outputIssueSuffix = commitCfg.OutputIssueSuffix
	}

	repo := openRepo()
	headRef := getCurrentHead(repo)

	// Read branch name or HEAD
	if headRef.Name().IsBranch() {

		branchName := headRef.Name().Short()
		matches := findIssueMatchesInBranch(issueRegex, branchName)

		if len(matches) > 0 {
			joinedIssues := strings.Join(matches, ", ")
			commitMessage = fmt.Sprintf("%s%s%s%s", outputIssuePrefix, joinedIssues, outputIssueSuffix, commitMessage)
		}

		commitChanges(repo, commitMessage)

		fmt.Println(commitMessage)

	} else if headRef.Name().IsTag() {
		fmt.Printf("HEAD is a tag: %v\n", headRef.Name().Short())
	} else {
		fmt.Printf("Detached HEAD at %v\n", headRef.Hash())
	}
}

// Reads commit message from command line arguments
func getCommitMessage() string {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, helpers.Red("Commit message cannot be empty"))
		os.Exit(1)
	}

	return args[0]
}

// Opens the repository in current directory
func openRepo() *git.Repository {
	repo, err := git.PlainOpen(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, helpers.Red("Failed to open repository: %v\n"), err)
		os.Exit(1)
	}
	return repo
}

// Reads the current HEAD reference
func getCurrentHead(repo *git.Repository) *plumbing.Reference {
	headRef, err := repo.Head()
	if err != nil {
		fmt.Fprintf(os.Stderr, helpers.Red("Failed to read current HEAD: %v\n"), err)
		os.Exit(1)
	}
	return headRef
}

// Opens worktree
func openWorktree(repo *git.Repository) *git.Worktree {
	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Fprintf(os.Stderr, helpers.Red("Failed to open worktree: %v\n"), err)
		os.Exit(1)
	}
	return worktree
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

// Creates commit options with the author information
func makeCommitOptions(usr user.User) git.CommitOptions {
	return git.CommitOptions{
		Author: &object.Signature{
			Name:  usr.Name,
			Email: usr.Email,
			When:  time.Now(),
		},
		AllowEmptyCommits: false,
		Amend:             false,
	}
}

// Commits changes with provided message
func commitChanges(repo *git.Repository, commitMessage string) {
	worktree := openWorktree(repo)
	usr := user.GetUser(*repo)
	commitOptions := makeCommitOptions(usr)

	_, err := worktree.Commit(commitMessage, &commitOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, helpers.Red("Failed to commit: %v\n"), err)
		os.Exit(1)
	}
}
