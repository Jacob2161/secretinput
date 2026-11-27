// Package secretinput reads secrets for command line programs from a file or
// a terminal prompt that does not echo, so they never appear in arguments,
// shell history, or the process list.
package secretinput

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// ErrNotTerminal means Prompt was asked to read from something other than a
// terminal; a script should pass the secret by file or environment instead.
var ErrNotTerminal = errors.New("not a terminal; pass the secret by file or environment")

// ReadFile reads a secret from a file, without the carriage returns and line
// feeds that end it: "secret\n\n" reads as "secret".
func ReadFile(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(contents), "\r\n"), nil
}

// Prompt writes prompt to w, reads one line from the terminal in without
// echoing it, and then ends the prompt's line on w. It returns ErrNotTerminal
// when in is not a terminal. A typical call is
//
//	secret, err := secretinput.Prompt(os.Stdin, os.Stderr, "Password: ")
func Prompt(in *os.File, w io.Writer, prompt string) (string, error) {
	descriptor := int(in.Fd())
	if !term.IsTerminal(descriptor) {
		return "", ErrNotTerminal
	}
	if _, err := io.WriteString(w, prompt); err != nil {
		return "", err
	}
	secret, err := term.ReadPassword(descriptor)
	_, lineErr := io.WriteString(w, "\n")
	if err != nil {
		return "", fmt.Errorf("reading from the terminal: %w", err)
	}
	if lineErr != nil {
		return "", lineErr
	}
	return string(secret), nil
}
