GOPATH:=$(shell go env GOPATH)

.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: build
build:
	go build ./cmd/main.go

.PHONY: lint
lint:
	revive -config revive-config.toml -formatter friendly ./...
	staticcheck ./...
	golangci-lint run ./...

.PHONY: format
format:
	goimports-reviser -company-prefixes "github.com/mehdihadeli" -project-name "github.com/mehdihadeli/go-ecommerce-microservices" -imports-order "std,general,company,project" -recursive ./...
	golines -m 120 -w --ignore-generated .
	gofumpt -l -w .

.PHONY: test
test:
	go test -cover ./...

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
	@docker build -t go-catalogs-read:latest .
