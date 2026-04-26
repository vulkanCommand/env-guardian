# Env Guardian - Project State

## Project Goal
Env Guardian is a Go CLI tool to validate, lint, analyze, and diagnose environment variables before they break applications.

---

## Current Version
v0.1.3

---

## Current Status
## v0.1.3 COMPLETE

The project has progressed beyond basic validation and now includes multi-environment validation, workflow tooling, and codebase env usage analysis.

The CLI is stable and feature-complete for core environment validation.

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
  - potential sensitive keys (e.g., SECRET, JWT)

---

### Doctor
- `doctor` command
- checks:
  - `.env` existence
  - `.env.example` existence
  - missing keys in `.env` compared to example

---

### Codebase Analysis
- `scan-code` command
- scans Go, JavaScript, TypeScript, and Python source files
- detects env variables used in code but missing in the env file
- detects env variables present in the env file but unused in code
- detects likely variable naming mismatches

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

---

## Remaining Work (Next Phases)

### Next Feature Block - Security
- secret pattern detection
- repo secret scanner
- git history scanning
- warn if `.env` committed

---

### Log Exposure Protection
- env log scan
- detect accidental logging of secrets
- scan logs and code for exposed environment variables

---

### Encryption
- env encrypt
- env decrypt
- secure key-based encryption for environment secrets

---

### DevOps / Runtime
- Docker validation
- CI/CD validation mode
- pre-start validation wrapper

---

### Developer Experience
- JSON output
- GitHub Action
- VS Code extension

---

## Development Strategy
1. Core Validation (DONE)
2. Codebase Analysis (DONE)
3. Security (NEXT)
4. Log Exposure Protection
5. Encryption
6. DevOps
7. Developer Experience
8. Final UX polish

No jumping ahead.

---

## Git Status
- branch: main
- version: v0.1.3
- CLI stable
- ready for security feature block

---

## Next Step
Start **Security** feature block
