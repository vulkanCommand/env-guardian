#!/usr/bin/env sh
set -eu

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="${BINARY_NAME:-envguard}"

if [ -t 1 ]; then
  green="$(printf '\033[32m')"
  red="$(printf '\033[31m')"
  cyan="$(printf '\033[36m')"
  reset="$(printf '\033[0m')"
else
  green=""
  red=""
  cyan=""
  reset=""
fi

run_step() {
  label="$1"
  shift

  printf '%s' "$label "

  "$@" >/tmp/envguard-install.log 2>&1 &
  pid="$!"
  frames='|/-\'
  index=0

  while kill -0 "$pid" 2>/dev/null; do
    index=$((index + 1))
    frame=$(printf '%s' "$frames" | cut -c "$(((index % 4) + 1))")
    printf '\r%s %s' "$label" "$frame"
    sleep 0.1
  done

  if wait "$pid"; then
    printf '\r%s %s\n' "$label" "${green}done${reset}"
    rm -f /tmp/envguard-install.log
    return 0
  fi

  printf '\r%s %s\n' "$label" "${red}failed${reset}"
  cat /tmp/envguard-install.log
  rm -f /tmp/envguard-install.log
  return 1
}

printf '%s\n' 'Env Guardian installer'
printf '%s\n' '----------------------'

if ! command -v go >/dev/null 2>&1; then
  printf '%s\n' 'Error: Go is required to build envguard.'
  exit 1
fi

mkdir -p "$INSTALL_DIR"

run_step "Building $BINARY_NAME" go build -o "$INSTALL_DIR/$BINARY_NAME" ./cmd/envguard

printf '%sInstalled:%s %s\n' "$green" "$reset" "$INSTALL_DIR/$BINARY_NAME"
printf '%s\n' ''
printf '%s\n' "${cyan}Run:${reset}"
printf '  %s\n' "$INSTALL_DIR/$BINARY_NAME"
printf '  %s validate\n' "$INSTALL_DIR/$BINARY_NAME"
