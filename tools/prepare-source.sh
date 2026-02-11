#!/bin/sh

set -eux

# Pin Go and toolbox versions at a reasonable version
go get go@1.24.12 toolchain@1.24.12

# Update go.mod and go.sum:
go mod tidy
go mod vendor

# Generate all sources (skip vendor/):
go generate ./cmd/... ./weldr/...

# Format all sources (skip vendor/):
go fmt ./cmd/... ./weldr/...
