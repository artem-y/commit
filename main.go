package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, red("Commit message cannot be empty"))
		os.Exit(1)
	}

	var commitMessage string
	commitMessage = args[0]

	// Open the repository in the current directory
	repo, err := git.PlainOpen(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, red("Failed to open repository: %v\n"), err)
		os.Exit(1)
	}

	// Get current HEAD
	ref, err := repo.Head()
	if err != nil {
		fmt.Fprintf(os.Stderr, red("Failed to read current HEAD: %v\n"), err)
		os.Exit(1)
	}

	// Read branch name or HEAD
	if ref.Name().IsBranch() {

		branchName := ref.Name().Short()

		// TODO: Add support for custom regex
		matchesArray := findIssueNumbersInBranch(`#[0-9]+`, branchName)

		if len(matchesArray) > 0 {
			joinedIssues := strings.Join(matchesArray, ", ")
			commitMessage = fmt.Sprintf("[%s]: %s", joinedIssues, commitMessage)
		}

		worktree, err := repo.Worktree()
		if err != nil {
			fmt.Fprintf(os.Stderr, red("Failed to open worktree: %v\n"), err)
			os.Exit(1)
		}

		usr := getUser(*repo)
		commitOptions := makeCommitOptions(usr)

		_, err = worktree.Commit(commitMessage, &commitOptions)
		if err != nil {
			fmt.Fprintf(os.Stderr, red("Failed to commit: %v\n"), err)
			os.Exit(1)
		}

		fmt.Println(commitMessage)

	} else if ref.Name().IsTag() {
		fmt.Printf("HEAD is a tag: %v\n", ref.Name().Short())
	} else {
		fmt.Printf("Detached HEAD at %v\n", ref.Hash())
	}
}

// Wraps the message string in red color
func red(msg string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", msg)
}

// Searches the branch name for issue numbers matching the given regex
func findIssueNumbersInBranch(rgxRaw string, branchName string) []string {
	// TODO: Don't use MustCompile, handle errors
	rgx := regexp.MustCompile(rgxRaw)

	matches := rgx.FindAllStringSubmatch(branchName, -1)
	matchesArray := []string{}

	for _, match := range matches {
		matchesArray = append(matchesArray, match...)
	}

	return matchesArray
}

// Creates commit options with the author information
func makeCommitOptions(usr user) git.CommitOptions {
	return git.CommitOptions{
		Author: &object.Signature{
			Name:  usr.name,
			Email: usr.email,
			When:  time.Now(),
		},
		AllowEmptyCommits: false,
		Amend:             false,
	}
}
