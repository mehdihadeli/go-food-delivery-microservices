#!/bin/bash
readonly service="$1"
readonly project_name="$2"

docker build -t "gcr.io/$project_name/$service" "./internal" -f "./docker/app-prod/Dockerfile" --build-arg "SERVICE=$service"
docker push "gcr.io/$project_name/$service"