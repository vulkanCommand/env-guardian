# Open-Source Launch Guide

Use this guide to launch Env Guardian publicly.

## Repository Setup

- make the repository public
- set description to `CLI tool to validate, secure, encrypt, and diagnose environment variables`
- enable Issues
- add topics: `go`, `cli`, `dotenv`, `env`, `security`, `devops`, `ci`, `secrets`
- keep the README install path curl-first

## Local Verification

```bash
go test ./...
go build -o envguard ./cmd/envguard
./envguard version
./envguard validate
./envguard security
./envguard log-scan
./envguard ci --json
```

## Build Release Artifacts

```bash
sh scripts/deploy.sh
```

This writes release binaries to `dist/`.

## GitHub Release

Create a release for:

```text
Tag: v0.1.13
Title: Env Guardian v0.1.13
```

Use `docs/release-v0.1.13.md` as the release notes and attach the files from `dist/`.

## Announcement Copy

```text
Env Guardian v0.1.13 is live.

A Go CLI to validate, secure, encrypt, and diagnose environment variables before they break your application.

Install:
curl -fsSL https://raw.githubusercontent.com/vulkanCommand/env-guardian/main/scripts/install.sh | sh

GitHub:
https://github.com/vulkanCommand/env-guardian
```
