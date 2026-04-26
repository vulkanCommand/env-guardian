# Env Guardian

Env Guardian is a Go CLI tool to validate, lint, analyze, and diagnose environment variables before they break your application.

---

## Current Version

v0.1.3

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

v0.1.3 is complete.

Completed in this version:
- schema-based type validation
- support for boolean validation
- support for number validation
- support for URL validation
- optional schema support
- parser test coverage for missing schema fallback
- validator test coverage for empty schema behavior
- validate help updated to reflect optional schema behavior
- multi-environment validation via `--all`
- codebase env usage scan command
- example generation and sync workflow commands

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

### Next
- secret detection
- warn if `.env` is committed
- log exposure checks

### Later
- CI/CD integration
- GitHub Action
- VS Code extension

---

## Goal

Make environment configuration safe, predictable, and production-ready.
