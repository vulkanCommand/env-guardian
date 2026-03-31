# Env Guardian — Project State

## Project Goal
Env Guardian is a Go CLI tool to validate, lint, analyze, and diagnose environment variables before they break applications.

---

## Current Version
v0.1.2

---

## Current Status
## v1.2 COMPLETE

The project has moved beyond v1.1 and now includes the full v1.2 checkpoint.

Completed in this version:
- schema-based type validation
- support for boolean validation
- support for number validation
- support for URL validation
- validator unit tests for typed values
- CLI help updated to reflect typed validation
- feature committed and pushed to GitHub

The project is now ready for the next feature block.

---

## Working Features

### Validation
- `validate` command
- detects missing keys
- detects duplicate keys
- detects unused keys as warnings
- detects invalid typed values using schema
- compares target env file with example env file
- supports:
  - `envguard validate`
  - `envguard validate --file .env.prod`
  - `envguard validate --example .env.example.prod`
  - `envguard validate --file .env.prod --example .env.example.prod`
- rejects invalid, unknown, missing, or duplicate flags
- prints formatted summary with error and warning counts
- returns proper exit codes

### Type Validation
- schema file path: `examples/.env.types`
- schema format: `KEY=type`
- supported types:
  - `boolean` → `true` or `false`
  - `number` → numeric values using ParseFloat
  - `url` → valid URL with scheme and host
- type validation runs only when:
  - the key exists in the schema
  - the key exists in `.env`
- CLI output format for invalid types:
  - `[ERROR] Invalid type: KEY expected <type> but got "<value>"`

### Linting
- `lint` command
- detects invalid lines
- detects missing `=`
- detects empty keys
- validates env syntax
- ignores comments and empty lines
- supports:
  - `envguard lint`
  - `envguard lint --file .env.prod`
- rejects invalid, unknown, missing, or duplicate flags
- prints formatted summary
- returns proper exit codes

### Analysis
- `analyze` command
- counts total keys
- detects empty values
- detects potential sensitive keys
- supports:
  - `envguard analyze`
  - `envguard analyze --file .env.prod`
- rejects invalid, unknown, missing, or duplicate flags
- prints formatted summary
- returns proper exit codes

### Doctor
- `doctor` command
- checks target env file existence
- checks example env file existence
- detects missing keys from example file
- supports:
  - `envguard doctor`
  - `envguard doctor --file .env.prod --example .env.example.prod`
- rejects invalid, unknown, missing, or duplicate flags
- prints formatted summary
- returns proper exit codes

### CLI / UX
- `version` command
- root help works:
  - `envguard`
  - `envguard help`
  - `envguard --help`
- topic help works:
  - `envguard help validate`
  - `envguard help lint`
  - `envguard help analyze`
  - `envguard help doctor`
- subcommand help works:
  - `envguard validate --help`
  - `envguard lint --help`
  - `envguard analyze --help`
  - `envguard doctor --help`
- `envguard help validate` now includes typed validation info and references `examples/.env.types`

---

## Completed In This Chat

### v0.1.2 Features Completed
- created schema loader in `internal/parser/types.go`
- added schema loading into validate flow
- threaded schema into validator layer
- added `InvalidTypeValues` to validation result
- implemented boolean validation
- implemented number validation
- implemented URL validation
- added CLI output for invalid type errors
- added validator unit tests
- fixed test structure to match `models.EnvFile`
- verified success and failure scenarios manually in CLI
- renamed schema file to `.env.types`
- updated validate help output
- bumped version to `v0.1.2`
- committed and pushed changes to GitHub

---

## Current Commands

### Root
- `envguard help`
- `envguard version`

### Validation
- `envguard validate`
- `envguard validate --file .env.prod`
- `envguard validate --example .env.example.prod`
- `envguard validate --file .env.prod --example .env.example.prod`

### Lint
- `envguard lint`
- `envguard lint --file .env.prod`

### Analyze
- `envguard analyze`
- `envguard analyze --file .env.prod`

### Doctor
- `envguard doctor`
- `envguard doctor --file .env.prod --example .env.example.prod`

### Help Topics
- `envguard help validate`
- `envguard help lint`
- `envguard help analyze`
- `envguard help doctor`

### Tests
- `go test ./internal/validator`

---

## Architecture
- `cmd/envguard` → CLI entry point
- `internal/parser` → parses env files and loads `.env.types`
- `internal/validator` → validation logic and type validation
- `internal/linter` → lint logic
- `internal/analyzer` → analysis logic
- `internal/doctor` → doctor diagnostics
- `internal/version` → version constant
- `internal/reporter` → future output formatting layer

---

## Feature Coverage vs Vision

### Implemented

#### Core Validation
- env validation
- syntax linting
- missing variable detection
- duplicate variable detection
- unused variable detection
- `.env` vs `.env.example` comparison
- custom target file selection
- custom example file selection
- schema-based typed value validation

#### Analysis
- total key count
- empty value detection
- potential sensitive key detection

#### Doctor
- env presence checks
- example presence checks
- missing required key checks
- summary output
- proper exit codes

#### CLI
- version command
- help command
- command-specific help
- stricter argument handling

#### Type Validation Coverage
- boolean
- number
- URL

---

## Remaining Work

### Immediate Next Feature Block
#### v0.1.3 — Make Schema Optional
Current issue:
- CLI fails if `examples/.env.types` is missing

Next change:
- if schema file is missing:
  - do not fail validation
  - skip type validation
  - continue normal validation flow

Expected behavior:
- graceful fallback
- non-breaking validate command
- typed validation only when schema exists

### Core Environment Validation
- optional schema support
- multi-environment support across `.env.dev`, `.env.prod`, `.env.test`
- environment consistency check across multiple environment files
- sync `.env` with `.env.example`
- generate `.env.example` automatically

### Codebase Analysis
- scan codebase for env usage
- detect variables used in code but missing in env files
- detect variables present in env files but unused in code
- detect variable naming mismatches

### Security
- real secret detection with patterns
- repository secret scanner
- git history secret scanning
- warn if `.env` is committed

### Log Exposure Protection
- env log scan
- detect accidental logging of secrets
- detect exposed tokens, keys, and passwords in logs or code

### Encryption
- env encrypt
- env decrypt
- key-based secret encryption

### DevOps / Runtime
- Docker environment validation
- CI/CD validation mode
- pre-start validation wrapper

### Developer Experience
- JSON output
- GitHub Action integration
- VS Code extension

### Meta Tool
- full aggregated env doctor report

---

## Development Strategy
Proceed strictly in this order:
1. Core Validation Improvements
2. Codebase Analysis
3. Security
4. DevOps
5. Developer Experience
6. Final UX polish

No jumping ahead.

---

## Release Notes

### v0.1.0
- initial working CLI
- validate
- lint
- analyze
- doctor
- parser and validator foundations

### v0.1.1
- validate file and example flags
- lint file flag
- analyze file flag
- doctor file and example flags
- unused variable detection
- improved summaries
- safer CLI flag parsing
- command help support

### v0.1.2
- added schema-based type validation
- added support for boolean, number, and URL types
- added `examples/.env.types`
- added `InvalidTypeValues` validation result handling
- added CLI invalid type output
- added validator unit tests
- updated validate help text
- bumped version to `v0.1.2`

---

## Git Status
- branch: `main`
- current version: `v0.1.2`
- latest commit: `feat: add schema-based type validation for env variables (boolean, number, url)`
- pushed to GitHub: yes

---

## Important Notes For Next Chat
- `scripts/deploy.sh` is currently empty and not part of the working release flow yet
- use direct git commands for now
- no refactors unless explicitly needed
- keep the one-step-at-a-time workflow
- do not jump ahead to JSON output, security, or CI/CD yet
- the exact next implementation target is: make schema optional without breaking validate

---

## Next Chat Start Line
Continue v0.1.3 — make schema optional