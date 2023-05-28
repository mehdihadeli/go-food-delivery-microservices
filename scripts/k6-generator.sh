#!/bin/bash

# https://craftbakery.dev/testing-rest-api-using-k6/
# https://k6.io/blog/load-testing-your-api-with-swagger-openapi-and-k6/

set -e

readonly service="$1"

@echo Generating load test client for catalogs write service
docker run --rm -v ${PWD}:/local  openapitools/openapi-generator-cli generate --skip-validate-spec -i  "local/api/$service/openapi/swagger.json" -g k6 -o "local/internal/services/$service/test/load_tests/"
