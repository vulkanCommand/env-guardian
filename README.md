# Env Guardian

[![CI](https://github.com/vulkanCommand/env-guardian/actions/workflows/envguard.yml/badge.svg)](https://github.com/vulkanCommand/env-guardian/actions/workflows/envguard.yml)
[![Release](https://img.shields.io/github/v/release/vulkanCommand/env-guardian?label=release)](https://github.com/vulkanCommand/env-guardian/releases)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go](https://img.shields.io/github/go-mod/go-version/vulkanCommand/env-guardian?label=go)](go.mod)

Env Guardian is a Go CLI tool to validate, lint, analyze, and diagnose environment variables before they break your application.

---

## Current Version

v0.1.13

---

## What It Does

Env Guardian helps you catch environment configuration issues early:

- missing variables
- duplicate variables
- unused variables
- invalid env syntax
- invalid typed values through an optional schema
- potential sensitive keys
- codebase env usage mismatches
- secret leaks in env files, repository files, and git history
- accidental logging of environment secrets
- encrypted environment files
- Docker and runtime environment checks
- JSON output for automation and CI
- VS Code command palette integration
- green launch-ready CLI help and support links
- colored pass, warning, and error output
- curl-based one-command installer

---

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/vulkanCommand/env-guardian/main/scripts/install.sh | sh
```

The installer builds Env Guardian from GitHub, installs it into a user-local bin directory, and shows an animated build step. Go is required.

Run after install:

```bash
envguard
envguard validate
envguard security
envguard ci --json
```

---

## Commands

When opened without arguments, Env Guardian shows a green ASCII title card, a command table, quick-start examples, and support links.

### Validate

```bash
envguard validate
envguard validate --all
envguard validate --file .env.prod
envguard validate --example .env.example.prod
envguard validate --file .env.prod --example .env.example.prod
envguard validate --json
```

Checks:
- missing keys compared to the example file
- duplicate keys in the target env file
- unused keys not present in the example file
- typed validation from `examples/.env.types` when the schema file exists

`--all`:
- validates `.env.dev`, `.env.prod`, and `.env.test`
- skips environments that do not exist
- runs validation per environment with grouped output

### Lint

```bash
envguard lint
envguard lint --file .env.prod
envguard lint --json
```

Checks:
- invalid syntax
- malformed lines
- missing `=`
- empty keys

### Analyze

```bash
envguard analyze
envguard analyze --file .env.prod
envguard analyze --json
```

Outputs:
- total keys
- empty values
- potential sensitive keys

### Doctor

```bash
envguard doctor
envguard doctor --file .env.prod --example .env.example.prod
envguard doctor --json
```

Checks:
- env file existence
- example file existence
- missing required keys
- tracked env file warning

### Security

```bash
envguard security
envguard security --dir .
envguard security --file .env.prod
envguard security --dir . --file .env.prod
envguard security --json
```

Checks:
- secret-looking values in the env file
- secret-looking values in repository files
- secret-looking values in git history
- tracked env files in git

### Log Scan

```bash
envguard log-scan
envguard log-scan --dir .
envguard log-scan --json
```

Checks:
- source code that logs env variable values
- log files containing secret-looking values
- log files containing sensitive key/value pairs

### Encryption

```bash
envguard encrypt
envguard encrypt --file .env.prod --out .env.prod.enc
envguard decrypt
envguard decrypt --file .env.prod.enc --out .env.prod
```

Checks:
- uses `ENVGUARD_KEY` for encryption and decryption
- encrypts env files with AES-GCM
- decrypts Env Guardian encrypted files
- writes output to the selected file

### DevOps / Runtime

```bash
envguard docker
envguard docker --dockerfile Dockerfile --file .env.prod
envguard ci
envguard ci --file .env.prod --example .env.example.prod
envguard ci --json
envguard run -- go run ./cmd/envguard
envguard run --file .env.prod --example .env.example.prod -- ./app
```

Checks:
- Dockerfile `ARG`, `ENV`, `$KEY`, and `${KEY}` references
- fail-fast CI validation for lint, required keys, duplicates, and typed values
- pre-start validation before running an application command

### Scan Code

```bash
envguard scan-code
envguard scan-code --dir .
envguard scan-code --dir . --file .env.prod
envguard scan-code --json
```

Checks:
- env variables used in code but missing in the env file
- env variables present in the env file but not used in code
- likely variable naming mismatches

Supported patterns include Go, JavaScript, TypeScript, and Python env access.

### Workflow

```bash
envguard generate-example
envguard sync-example
```

Checks:
- generate `.env.example` from `.env`
- sync missing keys from `.env` into `.env.example`

---

## Developer Experience

Machine-readable output:

```bash
envguard validate --json
envguard lint --json
envguard analyze --json
envguard doctor --json
envguard scan-code --json
envguard security --json
envguard log-scan --json
envguard docker --json
envguard ci --json
```

GitHub Actions:
- `.github/workflows/envguard.yml` runs tests, builds the CLI, prepares `.env` from `.env.example`, and runs CI/security/log exposure checks with JSON output.

VS Code:
- `vscode-extension/` contains a lightweight extension that runs the existing `envguard` executable from the command palette.
- commands include Validate, Validate All Environments, CI Check, Security Scan, Log Exposure Scan, and Show Version.
- settings allow configuring executable path, target env file, example env file, root directory, and JSON output.
- marketplace packaging metadata is included in `vscode-extension/package.json`.
- package a `.vsix` with `cd vscode-extension && npm run package`.

---

## Type Validation

Type validation uses:

```text
examples/.env.types
```

Example:

```text
DEBUG=boolean
PORT=number
API_URL=url
```

Supported types:
- boolean
- number
- url

If `examples/.env.types` is missing, validation still runs normally and type checks are skipped.

---

## Local Development

Build:

```bash
go build -o envguard ./cmd/envguard
```

Run:

```bash
./envguard
./envguard version
./envguard help
./envguard help validate
```

Tests:

```bash
go test ./...
```

---

## Support

- Email: `gdkalyan2109@gmail.com`
- Issues: `https://github.com/vulkanCommand/env-guardian/issues`

---

## Open-Source Launch

v0.1.13 is the first launch-ready version of Env Guardian.

- install with one curl command
- run locally, in CI, or before app startup
- use JSON output for automation
- report bugs and feature requests through GitHub Issues

Release notes: `docs/release-v0.1.13.md`

---

## Current Status

v0.1.13 is complete.

Completed in this version:
- clearer `ENV GUARDIAN` title banner
- curl-based one-command install docs
- launch README badges
- first public release notes

---

## Roadmap

### Current Phase
- core validation
- linting
- analysis
- doctor
- schema-based type validation
- optional schema support
- codebase env usage analysis
- team workflow commands
- security scanning
- log exposure protection
- encryption
- DevOps/runtime validation
- Developer Experience JSON output
- GitHub Action
- VS Code extension
- Final UX polish
- VS Code Marketplace packaging

### Next
- Open-source launch

### Later
- npm wrapper and VS Code marketplace publishing

---

## Goal

Make environment configuration safe, predictable, and production-ready.
