# Env Guardian VS Code Extension

Run Env Guardian validation, CI, security, and log exposure checks from the VS Code command palette.

Env Guardian for VS Code uses the `envguard` CLI installed on your machine.

Marketplace identifier:

```text
vulkanCommand.envguard-cli
```

## Install the CLI

```bash
curl -fsSL https://raw.githubusercontent.com/vulkanCommand/env-guardian/main/scripts/install.sh | sh
```

After installing, confirm VS Code can find the CLI:

```bash
envguard version
```

## Commands

- `Env Guardian: Validate`
- `Env Guardian: Validate All Environments`
- `Env Guardian: Run CI Check`
- `Env Guardian: Run Security Scan`
- `Env Guardian: Run Log Exposure Scan`
- `Env Guardian: Show Version`

## Settings

- `envGuardian.executablePath` - path to the `envguard` executable, defaults to `envguard`
- `envGuardian.envFile` - target env file
- `envGuardian.exampleFile` - example env file
- `envGuardian.rootDirectory` - root directory for scan commands
- `envGuardian.useJson` - pass `--json` to supported commands

## If VS Code Cannot Find Env Guardian

If VS Code shows that Env Guardian CLI is not installed, either install it with:

```bash
curl -fsSL https://raw.githubusercontent.com/vulkanCommand/env-guardian/main/scripts/install.sh | sh
```

Or set `envGuardian.executablePath` to the full binary path.

Example:

```json
{
  "envGuardian.executablePath": "C:/Users/gdkal/.local/bin/envguard"
}
```

## Support

- Issues: https://github.com/vulkanCommand/env-guardian/issues
- Email: gdkalyan2109@gmail.com
