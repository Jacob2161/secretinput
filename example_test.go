package secretinput_test

import (
	"errors"
	"fmt"
	"os"

	"github.com/jacob2161/secretinput"
)

// readPassword takes the password from a file when one is named, and asks
// on the terminal otherwise.
func readPassword(passwordFile string) (string, error) {
	if passwordFile != "" {
		return secretinput.ReadFile(passwordFile)
	}
	password, err := secretinput.Prompt(os.Stdin, os.Stderr, "Password: ")
	if errors.Is(err, secretinput.ErrNotTerminal) {
		return "", fmt.Errorf("pass --password-file when not on a terminal: %w", err)
	}
	return password, err
}

func Example() {
	_, _ = readPassword(os.Getenv("APP_PASSWORD_FILE"))
}
