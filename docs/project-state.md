# Env Guardian — Project State

## Project Goal

Env Guardian is a Go CLI tool to validate, lint, analyze, and diagnose environment variables before they break applications.

---

## Current Status

## v1.1 COMPLETE

The project has moved beyond the original v1 core and now includes the full v1.1 checkpoint:

- CLI file flag support
- unused variable detection
- improved doctor aggregation and summaries
- consistent command help and safer CLI argument handling

The project is now ready to start the next real feature block in a new chat.

---

## Current Version

- `v0.1.1`

---

## Working Features

### Validation

- `validate` command
- detects missing keys
- detects duplicate keys
- detects unused keys as warnings
- compares target env file with example env file
- supports:
  - `envguard validate`
  - `envguard validate --file .env.prod`
  - `envguard validate --example .env.example.prod`
  - `envguard validate --file .env.prod --example .env.example.prod`
- rejects invalid, unknown, missing, or duplicate flags
- prints formatted summary with error and warning counts
- returns proper exit codes

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

---

## Completed In This Chat

### v1.1 Features Completed

- added `validate --file`
- added `validate --example`
- added unused variable detection
- added warning handling in validation
- improved validation summaries
- added safer validation flag parsing
- added `lint --file`
- added `analyze --file`
- added `doctor --file --example`
- improved doctor aggregation and summaries
- added command-specific help support
- added topic help via `envguard help <command>`
- improved root help output consistency
- verified each change step-by-step with local build and command tests

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

---

## Architecture

- `cmd/envguard` → CLI entry point
- `internal/parser` → parses env files
- `internal/validator` → validation logic
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

---

## Remaining Work

### Core Environment Validation
- type validation for values such as boolean, number, URL
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

## Correct Next Step

The next real feature is:

## Start type validation

Suggested first scope:
- boolean validation
- integer / number validation
- URL validation

This should be implemented as the next feature block in the new chat.

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

---

## Important Notes For Next Chat

- `scripts/deploy.sh` is currently empty and not part of the working release flow yet
- use direct git commands for now
- do not spend more time on CLI polish
- start immediately with type validation
- keep the one-step-at-a-time workflow
- no refactors unless explicitly needed