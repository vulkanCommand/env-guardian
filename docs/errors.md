# Error Guide

Env Guardian separates blocking errors from warnings.

## Blocking Errors

Blocking errors return exit code `1`.

Common examples:

- missing required key
- duplicate env key
- invalid typed value
- invalid env syntax
- missing Docker runtime key
- secret-looking value found by `security`
- log exposure found by `log-scan`
- unreadable input file

## Warnings

Warnings do not always fail the command.

Common examples:

- unused env key
- potential sensitive key name from `analyze`
- tracked env file warning from `doctor`
- git history scan skipped outside a git repository

## Troubleshooting

Run command-specific help:

```bash
envguard help validate
envguard help security
envguard help ci
```

Use JSON output for automation:

```bash
envguard ci --json
```

## Support

- Email: `gdkalyan2109@gmail.com`
- Issues: `https://github.com/vulkanCommand/env-guardian/issues`
