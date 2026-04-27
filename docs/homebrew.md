# Homebrew

Env Guardian includes a Homebrew formula at `Formula/envguard.rb`.

## Install

```bash
brew tap vulkanCommand/env-guardian https://github.com/vulkanCommand/env-guardian
brew install envguard
```

One line:

```bash
brew tap vulkanCommand/env-guardian https://github.com/vulkanCommand/env-guardian && brew install envguard
```

## Upgrade

```bash
brew update
brew upgrade envguard
```

## Test

```bash
envguard version
envguard
envguard validate
```

## Notes

The clean command `brew install envguard` works after tapping this repository.

To make `brew install envguard` work without a tap, Env Guardian would need to be accepted into `homebrew/core`.
