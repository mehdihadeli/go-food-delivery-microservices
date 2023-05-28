#!/bin/bash

set -e

readonly service="$1"

if [ "$service" = "pkg" ]; then
    cd "./internal/pkg"
    go build ./...
# Check if input is not empty or null
elif [ -n "$service"  ]; then
    cd "./internal/services/$service"
    make build
fi

