# Outbox Publisher Entrypoint (`cmd/outbox-publisher`)

## Purpose

`cmd/outbox-publisher` boots the dedicated background process that reads pending rows from `aion_api.event_outbox` and publishes them to Kafka.

This entrypoint exists so publication cadence and failure handling can evolve independently from the main API server process.

## Current Runtime Flow

1. `main.go` invokes `run`.
2. `bootstrap_config.go` resolves bootstrap start and stop timeouts.
3. `bootstrap_fx.go` builds an Fx app with `fxapp.InfraModule` and `fxapp.OutboxPublisherModule`.
4. `bootstrap_runtime.go` starts the app, waits for process shutdown signals, and stops it gracefully.
5. `fxapp.OutboxPublisherModule` wires the repository, Kafka publisher, and periodic `PublishPending` loop.

## Boundary Rules

- no HTTP, GraphQL, or route registration belongs here
- business contexts decide when to enqueue events; this process only publishes pending rows
- durable configuration still comes from `internal/platform/config`

## Validate

```bash
go run ./cmd/outbox-publisher
make dev
make logs-api-publisher
```

## Risks And Compatibility Notes

- outbox worker behavior is operationally separate from the API process, so startup success of one does not prove health of the other
