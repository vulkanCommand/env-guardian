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
go test ./internal/validator
go test ./internal/parser
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

---

## Roadmap

### Current Phase
- core validation
- linting
- analysis
- doctor
- schema-based type validation
- optional schema support

### Next
- environment consistency checks across multiple env files
- sync `.env` with `.env.example`
- generate `.env.example` automatically

### Later
- codebase env usage analysis
- secret detection
- CI/CD integration
- GitHub Action
- VS Code extension

---

## Goal

Make environment configuration safe, predictable, and production-ready.