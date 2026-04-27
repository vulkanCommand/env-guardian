# Architecture

Env Guardian keeps CLI routing in `cmd/envguard` and domain logic in focused internal packages.

## Layout

- `cmd/envguard` - CLI entry point, flag parsing, terminal output, JSON output
- `internal/parser` - env parsing and optional type schema loading
- `internal/validator` - required keys, duplicates, unused keys, and typed values
- `internal/linter` - syntax checks
- `internal/analyzer` - empty values and sensitive-looking key names
- `internal/codebase` - env usage scanning in source files
- `internal/security` - env, repository, and git history secret scanning
- `internal/logscan` - accidental log exposure scanning
- `internal/encryption` - AES-GCM env encryption and decryption
- `internal/runtimecheck` - Dockerfile env reference validation
- `internal/doctor` - diagnostics
- `internal/version` - version constant
- `vscode-extension` - VS Code command palette wrapper around the CLI
- `.github/workflows/vscode-extension.yml` - VSIX packaging workflow

## Design Rules

- The Go CLI is the source of truth.
- The VS Code extension shells out to `envguard` instead of duplicating validation logic.
- JSON output is additive and does not replace human-readable output.
- Type validation is optional when `examples/.env.types` is missing.
- No command should print secret values in findings.

## Release Surface

- CLI binary from `go build -o envguard ./cmd/envguard`
- release archives from `scripts/deploy.sh`
- GitHub Actions workflow from `.github/workflows/envguard.yml`
- VS Code extension scaffold from `vscode-extension/`
- VSIX artifact workflow for tagged releases
