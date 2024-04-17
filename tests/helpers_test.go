package tests

import (
	"fmt"
	"testing"

	"github.com/artem-y/commit/internal/helpers"
)

func Test_Red_WrapsCommitMessage_InRedColor(t *testing.T) {
	errorMessage := "Error: Something went wrong"
	actualOutput := helpers.Red(errorMessage)
	expectedOutput := fmt.Sprintf("\033[31m%s\033[0m", errorMessage)

	if actualOutput != expectedOutput {
		t.Errorf("got %q, want %q", actualOutput, expectedOutput)
	}
}
