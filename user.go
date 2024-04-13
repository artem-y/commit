package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

// Simple representation of a git user
type user struct {
	name  string
	email string
}

// Returns the user name and email from the local repository or the global git config
func getUser(repo git.Repository) user {
	cfg, err := repo.Config()
	if err != nil {
		fmt.Fprintf(os.Stderr, red("Error loading local config: %v\n"), err)
		os.Exit(1)
	}

	usr := user{
		name:  cfg.User.Name,
		email: cfg.User.Email,
	}

	globalCfg, err := config.LoadConfig(config.GlobalScope)
	if err != nil {
		fmt.Fprintf(os.Stderr, red("Error loading global config: %v\n"), err)
		os.Exit(1)
	}

	if usr.email == "" {
		usr.email = globalCfg.User.Email
	}

	if usr.name == "" {
		usr.name = globalCfg.User.Name
	}

	return usr
}
