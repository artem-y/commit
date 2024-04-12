package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

func main() {
	// Open the repository in the current directory
	repo, err := git.PlainOpen(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open repository: %v\n", err)
		os.Exit(1)
	}
    
    // Get current HEAD
    ref, err := repo.Head()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to read current HEAD: %v\n", err)
        os.Exit(1)
    }

    // Output branch name or HEAD
    if ref.Name().IsBranch() {
        fmt.Printf("%v\n", ref.Name().Short())
    } else if ref.Name().IsTag() {
        fmt.Printf("HEAD is a tag: %v\n", ref.Name().Short())
    } else {
        fmt.Printf("Detached HEAD at %v\n", ref.Hash())
    }
}

