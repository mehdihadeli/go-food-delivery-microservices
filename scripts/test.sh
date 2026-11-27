#!/bin/bash

# https://blog.devgenius.io/go-golang-testing-tools-tips-to-step-up-your-game-4ed165a5b3b5

# In a bash script, set -e is a command that enables the "exit immediately" option. When this option is set, the script will terminate immediately if any command within the script exits with a non-zero status (indicating an error).
set -e

readonly service="${1:-}"
readonly type="${2:-unit}"

if [ -z "$service" ]; then
    echo "usage: $(basename "$0") <pkg|service> [unit|integration|e2e|load-test]" >&2
    exit 1
fi

if [ "$service" = "pkg" ]; then
    readonly module_dir="./internal/pkg"
else
    readonly module_dir="./internal/services/$service"
fi

if [ ! -d "$module_dir" ]; then
    echo "no module at $module_dir" >&2
    exit 1
fi

cd "$module_dir"

if [ "$type" = "load-test" ]; then
    # go run ./cmd/app/main.go
  	k6 run ./load_tests/script.js --insecure-skip-tls-verify
else
    go test -tags="$type" -timeout=30m -v -count=1 -p=1 -parallel=1 ./...
fi


