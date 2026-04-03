# API Entrypoint (`cmd/api`)

## Purpose

`cmd/api` is the process entrypoint for the `aion-api` binary.
It owns bootstrap, lifecycle, and process-level metadata only.

## Current Runtime Flow

1. `main.go` invokes `run`.
2. `bootstrap_config.go` resolves bootstrap start and stop timeouts.
3. `bootstrap_fx.go` builds the Fx application graph through the platform composition root.
4. `bootstrap_runtime.go` loads config, starts the app, waits for OS signals, and performs graceful shutdown.
5. `swagger.go` exposes the Swagger annotation block and injects runtime metadata used by the published REST docs.

## Boundary Rules

- no domain or usecase logic belongs in `cmd/api`
- route registration, GraphQL construction, and server wiring belong in `internal/platform`
- bootstrap-only env knobs stay here; durable config sections stay under `internal/platform/config`

## Validate

```bash
go run ./cmd/api
make dev
make verify
```

## Risks And Compatibility Notes

- bootstrap timeout changes and module selection here affect every runtime path, not just local startup ergonomics
