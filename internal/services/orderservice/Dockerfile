# syntax=docker/dockerfile:1

# Build context is the repository root because the service module replaces
# `.../internal/pkg` with the local `../../pkg` directory:
#
#   docker build -f internal/services/orderservice/Dockerfile .

FROM golang:1.27.4-alpine AS builder

WORKDIR /src

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOFLAGS=-mod=readonly

COPY internal/pkg/go.mod internal/pkg/go.sum ./internal/pkg/
COPY internal/services/orderservice/go.mod internal/services/orderservice/go.sum ./internal/services/orderservice/

RUN --mount=type=cache,target=/go/pkg/mod \
    cd internal/services/orderservice && go mod download

COPY internal/pkg ./internal/pkg
COPY internal/services/orderservice ./internal/services/orderservice

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd internal/services/orderservice && \
    go build -trimpath -buildvcs=false -ldflags="-s -w" -o /out/app ./cmd/app

FROM alpine:3.22 AS runtime

RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 appuser

WORKDIR /app

COPY --from=builder /out/app /app/app
COPY internal/services/orderservice/config/config.development.json /app/config/config.development.json

ENV APP_ENV=development \
    CONFIG_PATH=/app/config

USER appuser

EXPOSE 8000 6005

ENTRYPOINT ["/app/app"]
