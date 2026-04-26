#!/usr/bin/env sh
set -eu

VERSION="${VERSION:-$(go run ./cmd/envguard version)}"
DIST_DIR="${DIST_DIR:-dist}"

printf 'Building Env Guardian release %s\n' "$VERSION"

rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

build_binary() {
  os="$1"
  arch="$2"
  ext="$3"
  output="$DIST_DIR/envguard_${VERSION}_${os}_${arch}${ext}"

  printf '  %s/%s\n' "$os" "$arch"
  GOOS="$os" GOARCH="$arch" go build -o "$output" ./cmd/envguard
}

build_binary linux amd64 ""
build_binary linux arm64 ""
build_binary darwin amd64 ""
build_binary darwin arm64 ""
build_binary windows amd64 ".exe"

printf '%s\n' ''
printf 'Release artifacts written to %s\n' "$DIST_DIR"
printf '%s\n' 'Attach these files to the GitHub release.'
