#!/usr/bin/env sh
set -eu

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="${BINARY_NAME:-envguard}"

printf '%s\n' 'Env Guardian installer'
printf '%s\n' '----------------------'

if ! command -v go >/dev/null 2>&1; then
  printf '%s\n' 'Error: Go is required to build envguard.'
  exit 1
fi

mkdir -p "$INSTALL_DIR"

printf 'Building %s...\n' "$BINARY_NAME"
go build -o "$INSTALL_DIR/$BINARY_NAME" ./cmd/envguard

printf 'Installed: %s\n' "$INSTALL_DIR/$BINARY_NAME"
printf '%s\n' ''
printf '%s\n' 'Run:'
printf '  %s\n' "$INSTALL_DIR/$BINARY_NAME"
printf '  %s validate\n' "$INSTALL_DIR/$BINARY_NAME"
