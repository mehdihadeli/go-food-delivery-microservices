#!/bin/bash

# Builds one microservice image locally with the same Dockerfile the pipeline
# uses. Build context is the repository root: each service module replaces
# .../internal/pkg with the local ../../pkg directory.
#
#   ./scripts/build-docker.sh <service> [tag]
set -e

readonly service="$1"
readonly tag="${2:-dev}"

if [ -z "$service" ]; then
  echo "usage: $0 <catalogreadservice|catalogwriteservice|orderservice> [tag]" >&2
  exit 1
fi

docker build \
  -t "ghcr.io/mehdihadeli/go-food-delivery-microservices/${service}:${tag}" \
  -f "./internal/services/${service}/Dockerfile" \
  .
