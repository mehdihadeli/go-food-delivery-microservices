#!/bin/bash

# ref: https://blog.devgenius.io/sort-go-imports-acb76224dfa7
# https://yolken.net/blog/cleaner-go-code-golines

# In a bash script, set -e is a command that enables the "exit immediately" option. When this option is set, the script will terminate immediately if any command within the script exits with a non-zero status (indicating an error).
set -e

readonly service="$1"

if [ "$service" = "pkg" ]; then
      cd "./internal/pkg"
# Check if input is not empty or null
elif [ -n "$service"  ]; then
    cd "./internal/services/$service"
fi

# https://github.com/mvdan/gofumpt
gofumpt -l -w .

# https://golang.org/cmd/gofmt/
# gofmt -w .

# # https://pkg.go.dev/golang.org/x/tools/cmd/goimports
# goimports  . -l -w

# https://github.com/incu6us/goimports-reviser
# will do `gofmt` and `goimports` internally
goimports-reviser -rm-unused -set-alias -format -recursive ./...

# https://github.com/segmentio/golines
golines .  -m 120 -w --ignore-generated

