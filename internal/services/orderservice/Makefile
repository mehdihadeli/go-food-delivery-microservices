GOPATH:=$(shell go env GOPATH)

.PHONY: init
init:
	@go get -u google.golang.org/protobuf/proto
	@go install github.com/golang/protobuf/protoc-gen-go@latest

.PHONY: run_orders_service
run_orders_service:
	go run ./cmd/main.go

.PHONY: build_orders_service
build_orders_service:
	go build ./cmd/main.go

.PHONY: test_orders_service
test_orders_service:
	go test -v  ./... -cover

.PHONY: lint
lint:
	revive -config revive-config.toml -formatter friendly ./...
	staticcheck ./...
	golangci-lint run ./...

.PHONY: format
format:
	golines -m 120 -w --ignore-generated .
	gci write --skip-generated -s standard -s "prefix(github.com/mehdihadeli/go-food-delivery-microservices)" -s default -s blank -s dot --custom-order  .
	gofumpt -l -w .

.PHONY: update
update:
	@go get -u

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: deps-reset
deps-reset:
	git checkout -- go.mod
	go mod tidy

.PHONY: deps-upgrade
deps-upgrade:
	go get -u -t -d -v ./...
	go mod tidy

.PHONY: deps-cleancache
deps-cleancache:
	go clean -modcache

# ==============================================================================
# Linters https://golangci-lint.run/usage/install/

.PHONY: run-linter
run-linter:
	@echo Starting linters
	golangci-lint run ./...

.PHONY: docker
docker:
	@docker build -t go-orders:latest .

