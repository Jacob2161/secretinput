package secretinput

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestReadFile reads a secret without its trailing line ending, and keeps
// other whitespace.
func TestReadFile(t *testing.T) {
	directory := t.TempDir()
	for contents, want := range map[string]string{
		"hunter2\n":    "hunter2",
		"hunter2\r\n":  "hunter2",
		" spaced \n\n": " spaced ",
		"no newline":   "no newline",
	} {
		path := filepath.Join(directory, "secret")
		if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
		got, err := ReadFile(path)
		if err != nil || got != want {
			t.Errorf("ReadFile of %q = %q, %v; want %q", contents, got, err, want)
		}
	}
	if _, err := ReadFile(filepath.Join(directory, "missing")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ReadFile(a missing file) = %v, want os.ErrNotExist", err)
	}
}

// TestPromptRefusesNonTerminal checks that a pipe is refused without
// writing the prompt.
func TestPromptRefusesNonTerminal(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	defer writer.Close()
	var prompt strings.Builder
	if _, err := Prompt(reader, &prompt, "Password: "); !errors.Is(err, ErrNotTerminal) {
		t.Fatalf("Prompt(a pipe) = %v, want ErrNotTerminal", err)
	}
	if prompt.Len() != 0 {
		t.Fatalf("the prompt was written to a non-terminal: %q", prompt.String())
	}
}
