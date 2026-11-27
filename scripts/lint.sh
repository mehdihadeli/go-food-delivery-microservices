#!/bin/bash

# Lints one module. This is the single entry point used by `mise run lint:<module>`
# and by CI, so a developer and the pipeline run exactly the same checks.
#
# Usage:
#   ./scripts/lint.sh <catalogreadservice|catalogwriteservice|orderservice|pkg>
#   ./scripts/lint.sh pkg --new-from-rev=origin/main
#
# Extra arguments are forwarded to golangci-lint, which is where the ratchet
# flags live: --new-from-rev, --new-from-merge-base, --new-from-patch.
#
# The rules come from the shared .golangci.yml at the repository root.
# golangci-lint discovers it from any module directory, so there is one place to
# edit and no per-module copy that can drift.

set -euo pipefail

readonly service="${1:-}"
if [ -n "$service" ]; then
  shift
fi

if [ -z "$service" ]; then
  echo "usage: $0 <catalogreadservice|catalogwriteservice|orderservice|pkg> [golangci-lint args...]" >&2
  exit 1
fi

if [ "$service" = "pkg" ]; then
  readonly module="internal/pkg"
else
  readonly module="internal/services/${service}"
fi

if [ ! -d "$module" ]; then
  echo "error: '${service}' is not a known module, '${module}' does not exist" >&2
  exit 1
fi

cd "$module"

# https://golangci-lint.run/docs/linters/
# One run covers the linters the standalone tools below provide, and it is the
# only one that can be limited to the lines a change touches.
golangci-lint run "$@"

# golangci-lint runs revive, staticcheck and errcheck with its own settings. These
# standalone runs are kept because they additionally honour the per-module
# revive-config.toml and staticcheck.conf, which golangci-lint does not read.
#
# They are skipped when a tool is missing instead of failing the whole run, and
# CI disables the block with LINT_EXTRA_TOOLS=0: neither tool can be limited to a
# diff, so on a repository with existing findings they would always fail.
if [ "${LINT_EXTRA_TOOLS:-1}" != "0" ]; then
  # https://github.com/mgechev/revive
  if command -v revive >/dev/null 2>&1; then
    revive -config revive-config.toml -formatter friendly ./...
  else
    echo "note: revive is not installed, skipping (mise run install-tools)" >&2
  fi

  # https://github.com/dominikh/go-tools
  if command -v staticcheck >/dev/null 2>&1; then
    staticcheck ./...
  else
    echo "note: staticcheck is not installed, skipping (mise run install-tools)" >&2
  fi

  # https://github.com/kisielk/errcheck
  if command -v errcheck >/dev/null 2>&1; then
    errcheck ./...
  else
    echo "note: errcheck is not installed, skipping (mise run install-tools)" >&2
  fi
fi
