# Type Schema

Env Guardian can validate typed env values using:

```text
examples/.env.types
```

The schema file is optional. If it is missing, type validation is skipped silently and normal validation continues.

## Format

```text
KEY=type
```

Example:

```text
DEBUG=boolean
PORT=number
API_URL=url
```

## Supported Types

- `boolean`
- `number`
- `url`

## Behavior

- schema keys that do not exist in the target env file are ignored
- invalid schema lines fail validation
- missing schema file does not fail validation

## Example

```bash
envguard validate
envguard validate --json
```
