#!/bin/bash
set -e

readonly service="$1"

if [ "$service" = "pkg" ]; then
      cd "./internal/pkg"
# Check if input is not empty or null
elif [ -n "$service"  ]; then
    cd "./internal/services/$service"
fi

# https://golangci-lint.run/usage/linters/
# https://golangci-lint.run/usage/configuration/
# https://golangci-lint.run/usage/quick-start/
golangci-lint run
