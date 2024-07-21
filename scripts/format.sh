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

# https://github.com/segmentio/golines
# # will do `gofmt` internally
golines -m 120 -w --ignore-generated .


# # https://pkg.go.dev/golang.org/x/tools/cmd/goimports
# goimports -l -w .

# https://github.com/incu6us/goimports-reviser
# https://github.com/incu6us/goimports-reviser/issues/118
# https://github.com/incu6us/goimports-reviser/issues/88
# https://github.com/incu6us/goimports-reviser/issues/104
# will do `gofmt` internally if we use -format
# -rm-unused, -set-alias have some errors ---> goimports-reviser -rm-unused -set-alias -format -recursive ./...
# goimports-reviser -company-prefixes "github.com/mehdihadeli" -project-name "github.com/mehdihadeli/go-food-delivery-microservices" -rm-unused -set-alias -imports-order "std,general,company,project,blanked,dotted" -recursive ./...

gci write --skip-generated -s standard -s "prefix(github.com/mehdihadeli/go-food-delivery-microservices)" -s default -s blank -s dot --custom-order  .

# https://golang.org/cmd/gofmt/
# gofmt -w .

# https://github.com/mvdan/gofumpt
# will do `gofmt` internally
gofumpt -l -w .
