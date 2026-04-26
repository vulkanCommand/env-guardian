# Output Format

Env Guardian defaults to human-readable terminal output.

## Title Card

Running `envguard` with no arguments prints:

```text
+------------------------------------------------------------+
| Env Guardian                                               |
| Version 0.1.10                                             |
| Validate, secure, encrypt, and ship env files safely.      |
+------------------------------------------------------------+
```

The command overview appears below the card with quick-start examples and support links.

## Report Output

Report commands use consistent sections:

```text
<Report Name>
-------------
Target file: .env

[PASS] ...

Summary: 0 error(s), 0 warning(s)
```

Errors start with `[ERROR]`, warnings start with `[WARN]` or `[WARNING]`, and successful checks start with `[PASS]`.

## JSON Output

Report-style commands support `--json`:

```bash
envguard validate --json
envguard lint --json
envguard analyze --json
envguard doctor --json
envguard scan-code --json
envguard security --json
envguard log-scan --json
envguard docker --json
envguard ci --json
```

JSON reports include:

- `command`
- `status`
- file or directory fields when relevant
- `errors`
- `warnings`
- `summary`
- `details`

`status` is one of:

- `pass`
- `warning`
- `fail`

## Exit Codes

- `0` means the command completed successfully.
- `1` means Env Guardian found blocking issues or could not complete the command.
- `run` returns the wrapped command exit code after validation passes.
