.PHONY:

run_catalogs_write_service:
	cd services/catalogs/write_service/ && 	go run ./cmd/main.go

run_catalogs_read_service:
	cd services/catalogs/read_service/ && go run ./cmd/main.go

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
# Proto Catalogs Write Service

proto_catalogs_write_product_kafka_messages:
	@echo Generating products kafka messages proto
	protoc --go_out=./services/catalogs/write_service/internal/products/contracts/grpc/kafka_messages --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./services/catalogs/write_service/internal/products/contracts/grpc/kafka_messages api_docs/catalogs/write_service/protobuf/products/kafka_messages/product_kafka_messages.proto

proto_catalogs_write_product_service:
	@echo Generating product_service client proto
	protoc --go_out=./services/catalogs/write_service/internal/products/contracts/grpc/service_clients --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./services/catalogs/write_service/internal/products/contracts/grpc/service_clients api_docs/catalogs/write_service/protobuf/products/service_clients/product_service_client.proto


# ==============================================================================
# Proto Catalogs Read Service

proto_catalogs_read_product_kafka_messages:
	@echo Generating products kafka messages proto
	protoc --go_out=./services/catalogs/read_service/internal/products/contracts/grpc/kafka_messages --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./services/catalogs/read_service/internal/products/contracts/grpc/kafka_messages api_docs/catalogs/read_service/protobuf/products/kafka_messages/product_kafka_messages.proto

proto_catalogs_read_product_service:
	@echo Generating product_service client proto
	protoc --go_out=./services/catalogs/read_service/internal/products/contracts/grpc/service_clients --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./services/catalogs/read_service/internal/products/contracts/grpc/service_clients api_docs/catalogs/read_service/protobuf/products/service_clients/product_service_client.proto


# ==============================================================================
# Swagger Catalogs Write Service  #https://github.com/swaggo/swag/issues/817

swagger_catalogs_write:
	@echo Starting swagger generating
	swag init --parseDependency --parseInternal --parseDepth 1 -g ./services/catalogs/write_service/cmd/main.go -o ./services/catalogs/write_service/docs
	swag init --parseDependency --parseInternal --parseDepth 1 -g ./services/catalogs/write_service/cmd/main.go -o ./api_docs/catalogs/write_service/openapi/


# ==============================================================================
# Swagger Catalogs Read Service  #https://github.com/swaggo/swag/issues/817
swagger_catalogs_read:
	@echo Starting swagger generating
	swag init --parseDependency --parseInternal --parseDepth 1 -g ./services/catalogs/read_service/cmd/main.go -o ./services/catalogs/read_service/docs
	swag init --parseDependency --parseInternal --parseDepth 1 -g ./services/catalogs/read_service/cmd/main.go -o ./api_docs/catalogs/read_service/openapi/
