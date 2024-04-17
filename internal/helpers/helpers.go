package helpers

import "fmt"

// Wraps the message string in red color
func Red(msg string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", msg)
}

