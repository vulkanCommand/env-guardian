# Env Guardian VS Code Extension

Run Env Guardian from the VS Code command palette.

## Commands

- `Env Guardian: Validate`
- `Env Guardian: Validate All Environments`
- `Env Guardian: Run CI Check`
- `Env Guardian: Run Security Scan`
- `Env Guardian: Run Log Exposure Scan`
- `Env Guardian: Show Version`

## Settings

- `envGuardian.executablePath` - path to the `envguard` executable
- `envGuardian.envFile` - target env file
- `envGuardian.exampleFile` - example env file
- `envGuardian.rootDirectory` - root directory for scan commands
- `envGuardian.useJson` - pass `--json` to supported commands

## Local Use

Build the CLI first:

```bash
go build -o envguard ./cmd/envguard
```

Then set `envGuardian.executablePath` to the built binary path if `envguard` is not available on `PATH`.
