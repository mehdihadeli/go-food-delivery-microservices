.PHONY:

run_catalogs_service:
	go run ./services/catalogs/cmd/main.go -config=./services/catalogs/config/config.yaml

# ==============================================================================
# Docker Compose
docker-compose_infra_up:
	@echo Starting infrastructure docker-compose
	docker-compose -f deployments/docker-compose/docker-compose.infrastructure.yaml up --build

docker-compose_infra_down:
	@echo Stoping infrastructure docker-compose
	docker-compose -f deployments/docker-compose/docker-compose.infrastructure.yaml down

# ==============================================================================
# Docker

FILES := $(shell docker ps -aq)

docker_path:
	@echo $(FILES)

docker_down:
	docker stop $(FILES)
	docker rm $(FILES)

docker_clean:
	docker system prune -f

docker_logs:
	docker logs -f $(FILES)


# ==============================================================================
# Linters https://golangci-lint.run/usage/install/

run-linter:
	@echo Starting linters
	golangci-lint run ./...

# ==============================================================================
# PPROF

pprof_heap:
	go tool pprof -http :8006 http://localhost:6060/debug/pprof/heap?seconds=10

pprof_cpu:
	go tool pprof -http :8006 http://localhost:6060/debug/pprof/profile?seconds=10

pprof_allocs:
	go tool pprof -http :8006 http://localhost:6060/debug/pprof/allocs?seconds=10

# ==============================================================================
# Proto

proto_product_kafka_messages:
	@echo Generating products kafka messages proto
	protoc --go_out=./services/catalogs/internal/products/contracts/grpc/kafka_messages --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./services/catalogs/internal/products/contracts/grpc/kafka_messages api_docs/catalogs/protobuf/products/kafka_messages/product_kafka_messages.proto

proto_product_service:
	@echo Generating product_service client proto
	protoc --go_out=./services/catalogs/internal/products/contracts/grpc/service_clients --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./services/catalogs/internal/products/contracts/grpc/service_clients api_docs/catalogs/protobuf/products/service_clients/product_service_client.proto


# ==============================================================================
# Swagger

swagger_catalogs:
	@echo Starting swagger generating
	swag init -g ./services/catalogs/cmd/main.go -o ./services/catalogs/docs
	swag init -g ./services/catalogs/cmd/main.go -o ./api_docs/catalogs/openapi/
