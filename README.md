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
- invalid typed values (optional schema)
- potential sensitive keys

---

## Commands

### Validate

```bash
envguard validate
envguard validate --file .env.prod
envguard validate --file .env.prod --example .env.example.prod
