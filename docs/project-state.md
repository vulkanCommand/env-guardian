# Env Guardian - Project State

## Project Goal
Env Guardian is a Go CLI tool to validate, lint, analyze, secure, encrypt, and diagnose environment variables before they break applications.

---

## Current Version
v0.1.11

---

## Current Status
## v0.1.11 COMPLETE

The project now includes core validation, multi-environment checks, workflow tooling, codebase env usage analysis, security scanning, log exposure protection, encryption, DevOps/runtime validation, JSON output, GitHub Actions automation, VS Code command palette integration, launch-ready UX polish, green CLI status colors, animated installation, and VS Code Marketplace packaging.

The CLI is stable for the completed backend roadmap blocks through Final UX polish.

---

## Working Features

### Validation
- `validate` command
- detects missing keys
- detects duplicate keys
- detects unused keys (as warnings)
- detects invalid typed values using schema
- compares `.env` with `.env.example`
- supports:
  - `envguard validate`
  - `envguard validate --file`
  - `envguard validate --example`
  - `envguard validate --file --example`
- prints formatted summary with errors and warnings
- returns proper exit codes

---

### Multi-Environment Validation
- `envguard validate --all`
- supports:
  - `.env.dev`
  - `.env.prod`
  - `.env.test`
- runs validation per environment
- performs cross-environment consistency check
- identifies missing keys across environments

---

### Type Validation
- schema file: `examples/.env.types`
- optional (does not fail if missing)
- supported types:
  - boolean
  - number
  - URL

---

### Linting
- `lint` command
- detects:
  - invalid syntax
  - missing `=`
  - empty keys
- ignores comments and empty lines

---

### Analysis
- `analyze` command
- outputs:
  - total key count
  - empty values
  - potential sensitive keys

---

### Doctor
- `doctor` command
- checks:
  - `.env` existence
  - `.env.example` existence
  - missing keys in `.env` compared to example
  - tracked `.env` warning

---

### Codebase Analysis
- `scan-code` command
- scans Go, JavaScript, TypeScript, and Python source files
- detects env variables used in code but missing in the env file
- detects env variables present in the env file but unused in code
- detects likely variable naming mismatches

---

### Security
- `security` command
- detects secret-looking values in env files
- scans repository files for common leaked secrets
- scans git history for common leaked secrets
- warns when the target env file is tracked by git
- supported patterns include AWS, OpenAI, Stripe, GitHub, Slack, and private key blocks

---

### Log Exposure Protection
- `log-scan` command
- scans source code for direct logging of env variable values
- scans `.log` files for common leaked secrets
- scans `.log` files for sensitive key/value pairs

---

### Encryption
- `encrypt` command
- `decrypt` command
- reads encryption key from `ENVGUARD_KEY`
- encrypts env files with AES-GCM
- writes encrypted payloads in Env Guardian v1 format
- decrypts Env Guardian encrypted files back to plaintext env files

---

### DevOps / Runtime
- `docker` command
- `ci` command
- `run` command
- validates Dockerfile `ARG`, `ENV`, `$KEY`, and `${KEY}` references against an env file
- runs fail-fast CI checks for linting, required keys, duplicates, and typed values
- validates env configuration before starting an application command

---

### Developer Experience
- `--json` output for report-style commands
- supported by:
  - `envguard validate --json`
  - `envguard validate --all --json`
  - `envguard lint --json`
  - `envguard analyze --json`
  - `envguard doctor --json`
  - `envguard scan-code --json`
  - `envguard security --json`
  - `envguard log-scan --json`
  - `envguard docker --json`
  - `envguard ci --json`
- `.github/workflows/envguard.yml` runs Go tests, builds the CLI, and runs Env Guardian checks in GitHub Actions
- `vscode-extension/` provides command palette actions that run the existing `envguard` executable
- VS Code settings support executable path, env file, example file, root directory, and JSON output

---

### Workflow Commands

#### Generate Example
- `envguard generate-example`
- creates `.env.example` from `.env`
- preserves keys
- sets empty values

#### Sync Example
- `envguard sync-example`
- adds missing keys from `.env` into `.env.example`
- does not overwrite existing keys
- safe incremental update

---

### CLI / UX
- `version`
- `help`
- command-specific help
- strict flag validation
- consistent output formatting
- green ASCII root title card
- grouped command overview
- quick-start examples
- support email and GitHub Issues link
- green pass labels
- red error labels
- yellow warning labels
- `NO_COLOR` support

---

### Open Source Launch
- MIT license
- contributing guide
- security policy
- changelog
- install script
- release packaging script
- tag-based release workflow
- animated installer build step
- VSIX packaging workflow
- VS Code Marketplace metadata
- public docs for architecture, commands, errors, output format, roadmap, and schema

---

## Current Commands

### Core
- `envguard help`
- `envguard version`

### Validation
- `envguard validate`
- `envguard validate --all`
- `envguard validate --file`
- `envguard validate --example`

### Lint
- `envguard lint`
- `envguard lint --file`

### Analyze
- `envguard analyze`
- `envguard analyze --file`

### Doctor
- `envguard doctor`
- `envguard doctor --file --example`

### Codebase Analysis
- `envguard scan-code`
- `envguard scan-code --dir`
- `envguard scan-code --file`
- `envguard scan-code --dir --file`

### Security
- `envguard security`
- `envguard security --dir`
- `envguard security --file`
- `envguard security --dir --file`

### Log Exposure
- `envguard log-scan`
- `envguard log-scan --dir`

### Encryption
- `envguard encrypt`
- `envguard encrypt --file --out`
- `envguard decrypt`
- `envguard decrypt --file --out`

### DevOps / Runtime
- `envguard docker`
- `envguard docker --dockerfile --file`
- `envguard ci`
- `envguard ci --file --example`
- `envguard run -- <command>`
- `envguard run --file --example -- <command>`

### Workflow
- `envguard generate-example`
- `envguard sync-example`

---

## Architecture
- `cmd/envguard` - CLI entry point
- `internal/parser` - env parsing + schema loading
- `internal/validator` - validation + type checks
- `internal/linter` - syntax checks
- `internal/analyzer` - env insights
- `internal/codebase` - codebase env usage scanning
- `internal/security` - env, repository, git history security checks
- `internal/logscan` - log exposure scanning
- `internal/encryption` - env file encryption/decryption
- `internal/runtimecheck` - Docker and runtime environment checks
- `internal/doctor` - diagnostics
- `internal/version` - version constant

---

## Feature Coverage vs Vision

### Completed
- env validation
- syntax linting
- duplicate detection
- missing key detection
- unused key detection
- type validation (optional schema)
- multi-environment validation
- cross-environment consistency check
- example generation
- example sync
- codebase env usage analysis
- code/env missing key detection
- env/code unused key detection
- env variable naming mismatch detection
- secret leak detection in `.env`
- repository secret scanner
- git history secret scanning
- warn if `.env` is tracked by git
- env log scan
- accidental logging detection
- log file secret detection
- env encrypt
- env decrypt
- secure key-based encryption for environment secrets
- Docker validation
- CI/CD validation mode
- pre-start validation wrapper
- JSON output
- GitHub Action workflow
- VS Code extension
- final CLI title card
- colored status output
- public launch docs
- install and release scripts
- release artifact workflow
- VSIX artifact workflow
- VS Code Marketplace preparation
- open-source support files

---

## Remaining Work (Next Phases)

### Next Feature Block - Open Source Launch
- create first GitHub release
- attach release artifacts
- publish announcement

---

## Development Strategy
1. Core Validation (DONE)
2. Codebase Analysis (DONE)
3. Security (DONE)
4. Log Exposure Protection (DONE)
5. Encryption (DONE)
6. DevOps (DONE)
7. Developer Experience (DONE)
8. Final UX polish (DONE)

No jumping ahead.

---

## Git Status
- branch: main
- version: v0.1.11
- CLI stable
- ready for open-source launch

---

## Next Step
Launch **Env Guardian** as open source
