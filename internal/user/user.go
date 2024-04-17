package user

import (
	"fmt"
	"os"

    "github.com/artem-y/commit/internal/helpers"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

// Simple representation of a git user
type User struct {
	Name  string
	Email string
}

// Returns the user name and email from the local repository or the global git config
func GetUser(repo git.Repository) User {
	cfg, err := repo.Config()
	if err != nil {
		fmt.Fprintf(os.Stderr, helpers.Red("Error loading local config: %v\n"), err)
		os.Exit(1)
	}

	usr := User{
		Name:  cfg.User.Name,
		Email: cfg.User.Email,
	}

	globalCfg, err := config.LoadConfig(config.GlobalScope)
	if err != nil {
		fmt.Fprintf(os.Stderr, helpers.Red("Error loading global config: %v\n"), err)
		os.Exit(1)
	}

	if usr.Email == "" {
		usr.Email = globalCfg.User.Email
	}

	if usr.Name == "" {
		usr.Name = globalCfg.User.Name
	}

	return usr
}
