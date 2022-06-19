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
# Go migrate postgresql https://github.com/golang-migrate/migrate

DB_NAME = catalogs.service
DB_HOST = localhost
DB_USER = postgres
DB_PASS = postgres
DB_HOST = localhost
DB_PORT = 5432
SSL_MODE = disable

# go the last successful version, which is 1 here
# https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md#forcing-your-database-version
# https://github.com/golang-migrate/migrate/issues/282#issuecomment-530743258
# https://github.com/golang-migrate/migrate/issues/35
# https://github.com/golang-migrate/migrate/issues/21
# https://dev.to/techschoolguru/how-to-write-run-database-migration-in-golang-5h6g

postgres:
    docker run --name postgres -p $(DB_PORT)\:$(DB_PORT) -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASS) -d postgres:11.1-alpine

create_db:
	docker exec -it postgres createdb -U $(DB_USER) -O $(DB_USER) $(DB_NAME)

drop_db:
	docker exec -it postgres dropdb -U $(DB_USER) $(DB_NAME)

force_db:
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -verbose -path migrations force 1

version_db:
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -verbose -path migrations version

migrate_up:
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -verbose -path migrations up

migrate_down:
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -verbose -path migrations down


# ==============================================================================
# Swagger

swagger:
	@echo Starting swagger generating
	swag init -g **/**/*.go

# ==============================================================================
# Proto

proto_product_kafka_messages:
	@echo Generating products kafka messages proto
	protoc --go_out=./services/catalogs/internal/products/contracts/grpc/kafka_messages --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./services/catalogs/internal/products/contracts/grpc/kafka_messages api_docs/catalogs/protobuf/products/kafka_messages/product_kafka_messages.proto

proto_product_service:
	@echo Generating product_service client proto
	protoc --go_out=./services/catalogs/internal/products/contracts/grpc/service_clients --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./services/catalogs/internal/products/contracts/grpc/service_clients api_docs/catalogs/protobuf/products/service_clients/product_service_client.proto

