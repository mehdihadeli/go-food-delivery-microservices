#!/bin/bash

# In a bash script, set -e is a command that enables the "exit immediately" option. When this option is set, the script will terminate immediately if any command within the script exits with a non-zero status (indicating an error).
set -e

readonly service="$1"
readonly project_name="$2"

docker build -t "gcr.io/$project_name/$service" "./internal" -f "./docker/app-prod/Dockerfile" --build-arg "SERVICE=$service"
docker push "gcr.io/$project_name/$service"
