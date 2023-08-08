#!/bin/bash

# In a bash script, set -e is a command that enables the "exit immediately" option. When this option is set, the script will terminate immediately if any command within the script exits with a non-zero status (indicating an error).
set -e

readonly service="$1"

swag init --parseDependency --parseInternal --parseDepth 1  -g ./cmd/app/main.go  -d "./internal/services/$service/" -o "./internal/services/$service/docs"
swag init --parseDependency --parseInternal --parseDepth 1  -g ./cmd/app/main.go  -d "./internal/services/$service/" -o "./api/openapi/$service/"
