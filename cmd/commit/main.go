package main

import (
	"flag"
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
	var configFilePath string

	flag.StringVar(
		&configFilePath,
		"config-path",
		helpers.DEFAULT_CONFIG_FILE_PATH,
		"Path to the config json file",
	)
	flag.Parse()

	commitMessage := getCommitMessage()
	repo := openRepo()
	headRef := getCurrentHead(repo)

	// Read branch name or HEAD
	if headRef.Name().IsBranch() {

		cfg := config.ReadCommitConfig(configFilePath)
		branchName := headRef.Name().Short()
		matches := findIssueMatchesInBranch(cfg.IssueRegex, branchName)

		if len(matches) > 0 {
			commitMessage = generateCommitMessageWithMatches(matches, cfg, commitMessage)
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
	args := flag.Args()

	if len(args) < 1 || args[0] == "" {
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

// Generates a commit message with the issue number matches and config settings
func generateCommitMessageWithMatches(matches []string, cfg config.CommitConfig, commitMessage string) string {
	mappedMatches := make([]string, len(matches))

	for index, match := range matches {
		wrappedIssueNumber := fmt.Sprintf(
			"%s%s%s",
			*cfg.OutputIssuePrefix,
			match,
			*cfg.OutputIssueSuffix,
		)
		mappedMatches[index] = wrappedIssueNumber
	}

	joinedIssues := strings.Join(mappedMatches, ", ")
	return fmt.Sprintf(
		"%s%s%s%s",
		*cfg.OutputStringPrefix,
		joinedIssues,
		*cfg.OutputStringSuffix,
		commitMessage,
	)
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

	checkStagedChanges(worktree)

	usr := user.GetUser(*repo)
	commitOptions := makeCommitOptions(usr)

	_, err := worktree.Commit(commitMessage, &commitOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, helpers.Red("Failed to commit: %v\n"), err)
		os.Exit(1)
	}
}

// Checks if there are any staged changes to commit
func checkStagedChanges(worktree *git.Worktree) {
	fileStatuses, err := worktree.Status()
	if err != nil {
		fmt.Fprintf(os.Stderr, helpers.Red("Failed to read the status of the worktree: %v\n"), err)
		os.Exit(1)
	}

	for _, status := range fileStatuses {
		if status.Staging != git.Unmodified {
			return
		}
	}

	fmt.Fprintln(os.Stderr, helpers.Red("No staged changes to commit"))
	os.Exit(1)
}
