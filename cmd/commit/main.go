package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/artem-y/commit/internal/config"
	"github.com/artem-y/commit/internal/helpers"
	"github.com/artem-y/commit/internal/message_generator"
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
		"",
		"Path to the config json file",
	)

	var dryRun bool
	flag.BoolVar(
		&dryRun,
		"dry-run",
		false,
		"Prints the commit message without making the actual commit",
	)

	flag.Parse()

	commitMessage := getCommitMessage()
	repo := openRepo()
	worktree := openWorktree(repo)

	if configFilePath == "" {
		configFilePath = filepath.Join(
			worktree.Filesystem.Root(),
			helpers.DEFAULT_CONFIG_FILE_PATH,
		)
	}

	headRef := getCurrentHead(repo)

	// Read branch name or HEAD
	if headRef.Name().IsBranch() {

		fileReader := config.FileReader{}
		cfg, err := config.ReadCommitConfig(fileReader, configFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, helpers.Red("Failed to read config: %v\n"), err)
			os.Exit(1)
		}

		messageGenerator := message_generator.MessageGenerator{
			BranchName:  headRef.Name().Short(),
			UserMessage: commitMessage,
			Config:      cfg,
		}
		commitMessage = messageGenerator.GenerateMessage()

		if !dryRun {
			commitChanges(repo, worktree, commitMessage)
		}

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

// Opens the current repository
func openRepo() *git.Repository {

	options := git.PlainOpenOptions{DetectDotGit: true}

	repo, err := git.PlainOpenWithOptions(".", &options)
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
func commitChanges(repo *git.Repository, worktree *git.Worktree, commitMessage string) {

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
		if status.Staging != git.Unmodified && status.Staging != git.Untracked {
			return
		}
	}

	fmt.Fprintln(os.Stderr, helpers.Red("No staged changes to commit"))
	os.Exit(1)
}
