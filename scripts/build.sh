#!/bin/bash

# In a bash script, set -e is a command that enables the "exit immediately" option. When this option is set, the script will terminate immediately if any command within the script exits with a non-zero status (indicating an error).
set -e

readonly service="$1"

echo "start building $service"

if [ "$service" = "pkg" ]; then
    cd "./internal/pkg" && go build ./...
# Check if input is not empty or null
elif [ -n "$service"  ]; then
    cd "./internal/services/$service" && go build ./...
fi

