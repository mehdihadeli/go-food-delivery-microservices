#!/bin/bash

set -e

readonly service="$1"

cd "./internal/services/$service"

make run
