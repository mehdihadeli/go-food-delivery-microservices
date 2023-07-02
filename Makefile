.PHONY: install-tools
install-tools:
	./scripts/install-tools.sh

.PHONY: run
run:
	@./scripts/run.sh  catalogs_write
	@./scripts/run.sh  catalogs_read
	@./scripts/run.sh  orders


.PHONY: build
build:
	@./scripts/build.sh  catalogs_write
	@./scripts/build.sh  catalogs_read
	@./scripts/build.sh  orders


.PHONY: docker
docker-compose_infra_up:
	@echo Starting infrastructure docker-compose
	docker-compose -f deployments/docker-compose/docker-compose.infrastructure.yaml up --build

docker-compose_infra_down:
	@echo Stoping infrastructure docker-compose
	docker-compose -f deployments/docker-compose/docker-compose.infrastructure.yaml down

.PHONY: openapi
openapi:
	@./scripts/openapi.sh catalogs_write
	@./scripts/openapi.sh catalogs_read
	@./scripts/openapi.sh orders

.PHONY: proto
proto:
	@./scripts/proto.sh catalogs_write
	@./scripts/proto.sh catalogs_read
	@./scripts/proto.sh orders

.PHONY: unit_test
unit_test:
	@./scripts/test.sh catalogs_write unit
	@./scripts/test.sh catalogs_read unit
	@./scripts/test.sh  orders unit
	@./scripts/test.sh  pkg unit

.PHONY: integration_test
integration_test:
	@./scripts/test.sh catalogs_write integration
	@./scripts/test.sh catalogs_read integration
	@./scripts/test.sh  orders integration
	@./scripts/test.sh  pkg integration

.PHONY: e2e_test
e2e_test:
	@./scripts/test.sh catalogs_write e2e
	@./scripts/test.sh catalogs_read e2e
	@./scripts/test.sh  orders e2e

.PHONY: load_test
load_test:
	@./scripts/test.sh catalogs_write load-test
	@./scripts/test.sh catalogs_read load-test
	@./scripts/test.sh  orders load-test

.PHONY: format
format:
	@./scripts/format.sh catalogs_write
	@./scripts/format.sh catalogs_read
	@./scripts/format.sh orders
	@./scripts/format.sh pkg

.PHONY: lint
lint:
	@./scripts/lint.sh catalogs_write
	@./scripts/lint.sh catalogs_read
	@./scripts/lint.sh orders
	@./scripts/lint.sh pkg

.PHONY: c4
c4:
	cd tools/c4 && go mod tidy && sh generate.sh

# https://medium.com/yemeksepeti-teknoloji/mocking-an-interface-using-mockery-in-go-afbcb83cc773
# https://vektra.github.io/mockery/latest/running/
# https://amitshekhar.me/blog/test-with-testify-and-mockery-in-go
.PHONY: pkg-mocks
pkg-mocks:
	cd internal/pkg/messaging && mockery --output mocks --all
	cd internal/pkg/es && mockery --output mocks --all
	cd internal/pkg/core && mockery --output mocks --all

.PHONY: services-mocks
services-mocks:
	cd internal/services/catalogs_write && mockery --output mocks --all
	cd internal/services/catalogs_read && mockery --output mocks --all
	cd internal/services/orders && mockery --output mocks --all
