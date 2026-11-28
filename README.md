# secretinput

[![Go Reference](https://pkg.go.dev/badge/github.com/jacob2161/secretinput.svg)](https://pkg.go.dev/github.com/jacob2161/secretinput)
[![CI](https://github.com/jacob2161/secretinput/actions/workflows/ci.yml/badge.svg)](https://github.com/jacob2161/secretinput/actions/workflows/ci.yml)

Read secrets for command line programs from a file or a terminal prompt that does not echo, never
from arguments.

```bash
go get github.com/jacob2161/secretinput
```

## Usage

```go
password, err := secretinput.ReadFile(passwordFile) // trailing newline removed
// or, interactively:
password, err := secretinput.Prompt(os.Stdin, os.Stderr, "Password: ")
if errors.Is(err, secretinput.ErrNotTerminal) {
	// a script: ask for --password-file instead
}
```

- Secrets in arguments leak through shell history and the process list; a file or a prompt does not.
- `Prompt` refuses to read from a pipe, so scripts fail loudly instead of hanging.

Requires Go 1.27. Licensed under the MIT license.
