package main

import (
	"fmt"
	"os"
    "regexp"
    "strings"

	"github.com/go-git/go-git/v5"
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

    // Output branch name or HEAD
    if ref.Name().IsBranch() {

        branchName := ref.Name().Short()
        fmt.Printf("branch: %v\n", branchName)

        rgxp := regexp.MustCompile(`#[0-9]+`)
        matches := rgxp.FindAllStringSubmatch(branchName, -1)
        matchesArray := []string{}

        for _, match := range matches {
            for _, matchItem := range match {
                matchesArray = append(matchesArray, matchItem)
            }
        }

        if len(matchesArray) > 0 {
            joinedIssues := strings.Join(matchesArray, ", ")
            commitMessage = fmt.Sprintf("[%s]: %s", joinedIssues, commitMessage)
        }

        fmt.Println(commitMessage)

    } else if ref.Name().IsTag() {
        fmt.Printf("HEAD is a tag: %v\n", ref.Name().Short())
    } else {
        fmt.Printf("Detached HEAD at %v\n", ref.Hash())
    }
}

func red(msg string) string {
    return fmt.Sprintf("\033[31m%s\033[0m", msg)
}
