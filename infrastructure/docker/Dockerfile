# syntax=docker/dockerfile:1.4
# Use Go 1.25 to satisfy go.mod requirement (go 1.25.4 locally).
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Allow controlling ldflags via build-arg. For dev we omit strips, for prod we strip symbols.
ARG BUILD_LDFLAGS="-s -w"
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Use build cache for go build and module cache to speed up repeated builds.
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -ldflags="${BUILD_LDFLAGS}" -o aion-api ./cmd/api && \
    go build -ldflags="${BUILD_LDFLAGS}" -o aion-api-outbox-publisher ./cmd/outbox-publisher

# Separate stage: build the golang-migrate CLI with the postgres driver.
# Pinned version keeps distribution reproducible and isolated from the API modules.
FROM golang:1.25 AS migrate-builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1

FROM alpine:3.19.1

WORKDIR /app

COPY --from=builder /app/aion-api /usr/local/bin/aion-api
COPY --from=builder /app/aion-api-outbox-publisher /usr/local/bin/aion-api-outbox-publisher
COPY --from=migrate-builder /go/bin/migrate /usr/local/bin/migrate

COPY infrastructure/db/migrations /app/migrations
COPY infrastructure/docker/scripts/entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
