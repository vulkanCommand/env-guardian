# Env Guardian

Env Guardian is a CLI tool to validate, lint, and analyze environment variables before they break your application.

## Planned Commands

- envguard validate
- envguard lint
- envguard analyze
- envguard doctor

## Current Status

Initial CLI scaffold is working.

## Local Development

Build the binary

go build -o envguard ./cmd/envguard

Run version

./envguard version

Run help

./envguard

## Roadmap

v1
- validate
- lint
- analyze
- doctor

v1.1
- secret leak scan
- log exposure scan
- repository secret detection

v2
- env encryption
- env decryption
- runtime environment verification
