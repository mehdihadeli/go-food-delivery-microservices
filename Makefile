.PHONY:

run_catalogs_write_service:
	cd services/catalogs/write_service/ && 	go run ./cmd/main.go

run_catalogs_read_service:
	cd services/catalogs/read_service/ && go run ./cmd/main.go


# Docker Compose TASKS
docker-compose_infra_up:
	@echo Starting infrastructure docker-compose
	docker-compose -f deployments/docker-compose/docker-compose.infrastructure.yaml up --build

docker-compose_infra_down:
	@echo Stoping infrastructure docker-compose
	docker-compose -f deployments/docker-compose/docker-compose.infrastructure.yaml down


# DOCKER TASKS
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

# docker-build: ## [DOCKER] Build given container. Example: `make docker-build BIN=user`
	# docker build -f cmd/$(BIN)/Dockerfile --no-cache --build-arg BIN=$(BIN) --build-arg VERSION=$(VERSION) --build-arg GIT_COMMIT=$(GIT_COMMIT) -t go-api-boilerplate-$(BIN) .

# docker-run: ## [DOCKER] Run container on given port. Example: `make docker-run BIN=user PORT=3000`
	# docker run -i -t --rm -p=$(PORT):$(PORT) --name="go-api-boilerplate-$(BIN)" go-api-boilerplate-$(BIN)

# docker-stop: ## [DOCKER] Stop docker container. Example: `make docker-stop BIN=user`
	# docker stop go-api-boilerplate-$(BIN)

# docker-rm: docker-stop ## [DOCKER] Stop and then remove docker container. Example: `make docker-rm BIN=user`
	# docker rm go-api-boilerplate-$(BIN)

# docker-publish: docker-tag-latest docker-tag-version docker-publish-latest docker-publish-version ## [DOCKER] Docker publish. Example: `make docker-publish BIN=user REGISTRY=https://your-registry.com`

# docker-publish-latest:
	# @echo 'publish latest to $(REGISTRY)'
	# docker push $(REGISTRY)/go-api-boilerplate-$(BIN):latest

# docker-publish-version:
	# @echo 'publish $(VERSION) to $(REGISTRY)'
	# docker push $(REGISTRY)/go-api-boilerplate-$(BIN):$(VERSION)

# docker-tag: docker-tag-latest docker-tag-version ## [DOCKER] Tag current container. Example: `make docker-tag BIN=user REGISTRY=https://your-registry.com`

# docker-tag-latest:
	# @echo 'create tag latest'
	# docker tag go-api-boilerplate-$(BIN) $(REGISTRY)/go-api-boilerplate-$(BIN):latest

# docker-tag-version:
	# @echo 'create tag $(VERSION)'
	# docker tag go-api-boilerplate-$(BIN) $(REGISTRY)/go-api-boilerplate-$(BIN):$(VERSION)

# docker-release: docker-build docker-publish ## [DOCKER] Docker release - build, tag and push the container. Example: `make docker-release BIN=user REGISTRY=https://your-registry.com`

# # TELEPRESENCE TASKS
# telepresence-swap-local: ## [TELEPRESENCE] Replace the existing deployment with the Telepresence proxy for local process. Example: `make telepresence-swap-local BIN=user PORT=3000 DEPLOYMENT=api`
	# go build -o cmd/$(BIN)/$(BIN) cmd/$(BIN)/main.go
	# telepresence \
	# --swap-deployment $(DEPLOYMENT) \
	# --expose 3000 \
	# --run ./cmd/$(BIN)/$(BIN) \
	# --port=$(PORT) \
	# --method vpn-tcp

# telepresence-swap-docker: ## [TELEPRESENCE] Replace the existing deployment with the Telepresence proxy for local docker image. Example: `make telepresence-swap-docker BIN=user PORT=3000 DEPLOYMENT=api`
	# telepresence \
	# --swap-deployment $(DEPLOYMENT) \
	# --docker-run -i -t --rm -p=$(PORT):$(PORT) --name="$(BIN)" $(BIN):latest


# Linters TASK
run-linter:
	@echo Starting linters
	golangci-lint run ./...


# HTTPS TASK
key: ## [HTTP] Generate key
	openssl genrsa -out server.key 2048
	openssl ecparam -genkey -name secp384r1 -out server.key

cert: ## [HTTP] Generate self signed certificate
	openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650


# PPROF TASK
pprof_heap:
	go tool pprof -http :8006 http://localhost:6060/debug/pprof/heap?seconds=10

pprof_cpu:
	go tool pprof -http :8006 http://localhost:6060/debug/pprof/profile?seconds=10

pprof_allocs:
	go tool pprof -http :8006 http://localhost:6060/debug/pprof/allocs?seconds=10


# Proto Catalogs Write Service
proto_catalogs_write_product_kafka_messages:
	@echo Generating products kafka messages proto
	protoc --go_out=./services/catalogs/write_service/internal/products/contracts/proto/kafka_messages --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./services/catalogs/write_service/internal/products/contracts/proto/kafka_messages api_docs/catalogs/write_service/protobuf/products/kafka_messages/product_kafka_messages.proto

proto_catalogs_write_product_service:
	@echo Generating product_service client proto
	protoc --go_out=./services/catalogs/write_service/internal/products/contracts/proto/service_clients --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./services/catalogs/write_service/internal/products/contracts/proto/service_clients api_docs/catalogs/write_service/protobuf/products/service_clients/products_service_client.proto



# Proto Catalogs Read Service
proto_catalogs_read_product_kafka_messages:
	@echo Generating products kafka messages proto
	protoc --go_out=./services/catalogs/read_service/internal/products/contracts/proto/kafka_messages --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./services/catalogs/read_service/internal/products/contracts/proto/kafka_messages api_docs/catalogs/read_service/protobuf/products/kafka_messages/product_kafka_messages.proto

proto_catalogs_read_product_service:
	@echo Generating product_service client proto
	protoc --go_out=./services/catalogs/read_service/internal/products/contracts/proto/service_clients --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./services/catalogs/read_service/internal/products/contracts/proto/service_clients api_docs/catalogs/read_service/protobuf/products/service_clients/products_service_client.proto

# Proto Orders Service
proto_orders_order_service:
	@echo Generating order_service client proto
	protoc --go_out=./services/orders/internal/orders/contracts/proto/service_clients --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./services/orders/internal/orders/contracts/proto/service_clients/ api_docs/orders/protobuf/orders/service_clients/orders_service_client.proto


# Swagger Catalogs Write Service  #https://github.com/swaggo/swag/issues/817
swagger_catalogs_write:
	@echo Starting swagger generating
	swag init --parseDependency --parseInternal --parseDepth 1 -g ./cmd/main.go -d ./services/catalogs/write_service/ -o ./services/catalogs/write_service/docs
	swag init --parseDependency --parseInternal --parseDepth 1 -g ./cmd/main.go -d ./services/catalogs/write_service/ -o ./api_docs/catalogs/write_service/openapi/


# Swagger Catalogs Read Service  #https://github.com/swaggo/swag/issues/817
swagger_catalogs_read:
	@echo Starting swagger generating
	swag init --parseDependency --parseInternal --parseDepth 1  -g ./cmd/main.go  -d ./services/catalogs/read_service/ -o ./services/catalogs/read_service/docs
	swag init --parseDependency --parseInternal --parseDepth 1  -g ./cmd/main.go  -d ./services/catalogs/read_service/ -o ./api_docs/catalogs/read_service/openapi/

# Swagger Orders Service
swagger_orders:
	@echo Starting swagger generating
	swag init --parseDependency --parseInternal --parseDepth 1  -g ./cmd/main.go  -d ./services/orders/ -o ./services/orders/docs
	swag init --parseDependency --parseInternal --parseDepth 1  -g ./cmd/main.go  -d ./services/orders/ -o ./api_docs/orders/openapi/


## Generate Load Test Client for Catalogs Write Service  # #https://craftbakery.dev/testing-rest-api-using-k6/
generate_load_test_client_catalogs_write_service:
	@echo Generating load test client for catalogs write service
	docker run --rm -v ${PWD}:/local  openapitools/openapi-generator-cli generate --skip-validate-spec -i  local/api_docs/catalogs/write_service/openapi/swagger.json -g k6 -o local/performance_tests/catalogs/write_service/k6-test/


## Generate Load Test Client for Catalogs Read Service  # #https://craftbakery.dev/testing-rest-api-using-k6/
generate_load_test_client_catalogs_read_service:
	@echo Generating load test client for catalogs write service
	docker run --rm -v ${PWD}:/local  openapitools/openapi-generator-cli generate --skip-validate-spec -i  local/api_docs/catalogs/read_service/openapi/swagger.json -g k6 -o local/performance_tests/catalogs/read_service/k6-test/

## Execute k6 for catalog write service
execute_k6_catalogs_write_service:
	@echo Executing k6 for catalogs write service
	cd services/catalogs/write_service/ && 	go run ./cmd/main.go
	k6 run ./performance_tests/catalogs/write_service/script.js --insecure-skip-tls-verify

## Execute k6 for catalog read_position service
execute_k6_catalogs_read_service:
	@echo Executing k6 for catalogs read service
	cd services/catalogs/read_service/ && 	go run ./cmd/main.go
	k6 run ./performance_tests/catalogs/read_service/script.js --insecure-skip-tls-verify