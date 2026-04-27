# Env Guardian v0.1.13

Env Guardian v0.1.13 is the first launch-ready open-source release.

## Highlights

- clear `ENV GUARDIAN` CLI title banner
- green pass, yellow warning, and red error output
- curl-based one-command installer
- validation, linting, analysis, doctor, CI, Docker, and runtime checks
- codebase env usage scanning
- secret scanning for env files, repository files, and git history
- log exposure scanning
- env file encryption and decryption
- JSON output for automation
- GitHub Actions workflow for tests and Env Guardian checks
- release packaging script for Linux, macOS, and Windows binaries

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/vulkanCommand/env-guardian/main/scripts/install.sh | sh
```

## Quick Test

```bash
envguard
envguard version
envguard validate
envguard security
envguard ci --json
```

## Support

- Issues: https://github.com/vulkanCommand/env-guardian/issues
- Email: gdkalyan2109@gmail.com

## Release Checklist

- verify `go test ./...`
- verify `go build -o envguard ./cmd/envguard`
- verify `envguard validate`
- verify `envguard security`
- verify `envguard log-scan`
- create GitHub release for tag `v0.1.13`
- attach release artifacts from `dist/`
