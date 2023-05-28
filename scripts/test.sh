#!/bin/bash
set -e

# https://blog.devgenius.io/go-golang-testing-tools-tips-to-step-up-your-game-4ed165a5b3b5

readonly service="$1"
readonly type="$2"

if [ "$service" = "pkg" ]; then
      cd "./internal/pkg/$service"
else
    cd "./internal/services/$service"
fi

if [ "$type" = "load-test" ]; then
    # go run ./cmd/app/main.go
  	k6 run ./load_tests/script.js --insecure-skip-tls-verify
else
    go test -tags="$type" -count=1 -p=8 -parallel=8 -race ./...
fi


