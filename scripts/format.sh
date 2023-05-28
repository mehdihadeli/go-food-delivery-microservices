#!/bin/bash

#https://blog.devgenius.io/sort-go-imports-acb76224dfa7

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

# https://github.com/incu6us/goimports-reviser
goimports-reviser -rm-unused -set-alias -format -recursive ./...

# # https://pkg.go.dev/golang.org/x/tools/cmd/goimports
# goimports  . -l -w

# https://github.com/segmentio/golines
golines .  -m 120 -w --ignore-generated

