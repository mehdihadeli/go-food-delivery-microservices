.PHONY:

run_catalogs_service:
	go run ./cmd/main.go -config=./config/config.development.yaml

build_catalogs_service:
	go build ./cmd/main.go

test_catalogs_service:
	go test -cover ./...

# ==============================================================================
# Golang Helpers

tidy:
	go mod tidy

deps-reset:
	git checkout -- go.mod
	go mod tidy

deps-upgrade:
	go get -u -t -d -v ./...
	go mod tidy

deps-cleancache:
	go clean -modcache

# ==============================================================================
# Linters https://golangci-lint.run/usage/install/

run-linter:
	@echo Starting linters
	golangci-lint run ./...

# ==============================================================================
