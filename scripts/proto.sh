#!/bin/bash
set -e

readonly service="$1"

protoc \
  --proto_path=api/protobuf "api/proto/$service.proto" \
  "--go_out=internal/services/$service/proto/service_clients/" --go_opt=paths=source_relative \
  --go-grpc_opt=require_unimplemented_servers=false \
  "--go-grpc_out=internal/services/$service/proto/service_clients/" --go-grpc_opt=paths=source_relativ
