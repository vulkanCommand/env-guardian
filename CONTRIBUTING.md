# Contributing to Env Guardian

Thanks for helping improve Env Guardian.

## Development Setup

```bash
git clone https://github.com/vulkanCommand/env-guardian
cd env-guardian
go test ./...
go build -o envguard ./cmd/envguard
```

## Before Opening a Pull Request

Run:

```bash
go test ./...
go build -o envguard ./cmd/envguard
./envguard validate
./envguard security
./envguard log-scan
```

## Contribution Guidelines

- keep changes focused
- do not commit real secrets or local `.env` files
- add tests when changing validation, scanning, encryption, or runtime behavior
- keep CLI output stable unless the change is specifically about UX
- prefer JSON output additions that do not break existing fields

## Issues

Use GitHub Issues for bugs, feature requests, and questions:

```text
https://github.com/vulkanCommand/env-guardian/issues
```

For direct support, email:

```text
gdkalyan2109@gmail.com
```
