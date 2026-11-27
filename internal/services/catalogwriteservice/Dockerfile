# syntax=docker/dockerfile:1

# Build context is the repository root because the service module replaces
# `.../internal/pkg` with the local `../../pkg` directory:
#
#   docker build -f internal/services/catalogwriteservice/Dockerfile .

FROM golang:1.27.4-alpine AS builder

WORKDIR /src

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOFLAGS=-mod=readonly

COPY internal/pkg/go.mod internal/pkg/go.sum ./internal/pkg/
COPY internal/services/catalogwriteservice/go.mod internal/services/catalogwriteservice/go.sum ./internal/services/catalogwriteservice/

RUN --mount=type=cache,target=/go/pkg/mod \
    cd internal/services/catalogwriteservice && go mod download

COPY internal/pkg ./internal/pkg
COPY internal/services/catalogwriteservice ./internal/services/catalogwriteservice

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd internal/services/catalogwriteservice && \
    go build -trimpath -buildvcs=false -ldflags="-s -w" -o /out/app ./cmd/app && \
    go build -trimpath -buildvcs=false -ldflags="-s -w" -o /out/migration ./cmd/migration

FROM alpine:3.22 AS runtime

RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 appuser

WORKDIR /app

COPY --from=builder /out/app /app/app
COPY --from=builder /out/migration /app/migration
COPY internal/services/catalogwriteservice/config/config.development.json /app/config/config.development.json
COPY internal/services/catalogwriteservice/db/migrations /app/db/migrations

ENV APP_ENV=development \
    CONFIG_PATH=/app/config

USER appuser

EXPOSE 7000 6003

ENTRYPOINT ["/app/app"]
