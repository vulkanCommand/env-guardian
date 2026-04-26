# Env Guardian

Env Guardian is a Go CLI tool to validate, lint, analyze, and diagnose environment variables before they break your application.

---

## Current Version

v0.1.7

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

---

## Commands

### Validate

```bash
envguard validate
envguard validate --all
envguard validate --file .env.prod
envguard validate --example .env.example.prod
envguard validate --file .env.prod --example .env.example.prod
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
```

Outputs:
- total keys
- empty values
- potential sensitive keys

### Doctor

```bash
envguard doctor
envguard doctor --file .env.prod --example .env.example.prod
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

## Current Status

v0.1.7 is complete.

Completed in this version:
- Dockerfile environment validation
- CI/CD validation mode
- pre-start validation wrapper
- runtime command tests

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

### Next
- JSON output

### Later
- GitHub Action
- VS Code extension

---

## Goal

Make environment configuration safe, predictable, and production-ready.
