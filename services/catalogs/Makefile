.PHONY:

run_catalogs_service:
	go run ./cmd/main.go -config=./config/config.yaml

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
