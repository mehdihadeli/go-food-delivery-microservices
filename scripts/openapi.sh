#!/bin/bash
set -e

readonly service="$1"

swag init --parseDependency --parseInternal --parseDepth 1  -g ./cmd/app/main.go  -d "./internal/services/$service/" -o "./internal/services/$service/docs"
swag init --parseDependency --parseInternal --parseDepth 1  -g ./cmd/app/main.go  -d "./internal/services/$service/" -o "./api/openapi/$service/"
