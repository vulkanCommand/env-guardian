# Env Guardian Commands

Env Guardian is command-first. Running `envguard` with no arguments prints the title card, core command list, quick-start examples, and support links.

## Core

```bash
envguard
envguard help
envguard help <command>
envguard version
```

## Validation

```bash
envguard validate
envguard validate --all
envguard validate --file .env.prod
envguard validate --file .env.prod --example .env.example.prod
envguard validate --json
```

## Quality Checks

```bash
envguard lint
envguard analyze
envguard doctor
envguard scan-code
```

## Security

```bash
envguard security
envguard log-scan
envguard encrypt --file .env.prod --out .env.prod.enc
envguard decrypt --file .env.prod.enc --out .env.prod
```

`encrypt` and `decrypt` read the encryption key from `ENVGUARD_KEY`.

## DevOps

```bash
envguard docker
envguard ci
envguard run -- <command>
```

## Workflow

```bash
envguard generate-example
envguard sync-example
```

## Support

- Email: `gdkalyan2109@gmail.com`
- Issues: `https://github.com/vulkanCommand/env-guardian/issues`
