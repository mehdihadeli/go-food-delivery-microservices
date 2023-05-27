#!/bin/bash
set -e

readonly service="$1"

cd "./internal/services/$service"
go test -count=1 -p=8 -parallel=8 -race ./...