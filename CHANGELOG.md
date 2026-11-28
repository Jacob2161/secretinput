# Changelog

This module follows semantic versioning. Before v1, a minor version may
change the API; each such change says what callers must do.

## v0.2.1 (2026-11-28)

- CI runs every step in bash, so the Windows runners read
  `-modfile=tools/go.mod` whole. No API change.
- The README says the module requires Go 1.27.

## v0.2.0 (2026-11-28)

- Requires Go 1.27, the newest Go release; Go 1.26 is no longer
  supported.
- CI runs staticcheck and govulncheck at the versions pinned in
  `tools/go.mod`, a module of its own that importers never see.

## v0.1.1 (2026-11-27)

- `Prompt` is tested on a real pseudo-terminal on Linux: the secret is not
  echoed, and echo is restored afterwards. No API or behavior change.
- `golang.org/x/sys` is now a direct requirement, for that test; it was
  already required through `golang.org/x/term`.
- CI also checks formatting, `go mod tidy -diff`, and `govulncheck`.

## v0.1.0 (2026-11-27)

First release.

- `ReadFile` reads a secret from a file without its trailing line endings.
- `Prompt` reads a secret from a terminal without echoing it, and refuses
  anything that is not a terminal with `ErrNotTerminal`, so scripts fail
  instead of hanging.
