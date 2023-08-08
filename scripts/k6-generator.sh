#!/bin/bash

# https://craftbakery.dev/testing-rest-api-using-k6/
# https://k6.io/blog/load-testing-your-api-with-swagger-openapi-and-k6/

# In a bash script, set -e is a command that enables the "exit immediately" option. When this option is set, the script will terminate immediately if any command within the script exits with a non-zero status (indicating an error).
set -e

readonly service="$1"

@echo Generating load test client for catalogs write service
docker run --rm -v ${PWD}:/local  openapitools/openapi-generator-cli generate --skip-validate-spec -i  "local/api/$service/openapi/swagger.json" -g k6 -o "local/internal/services/$service/test/load_tests/"
