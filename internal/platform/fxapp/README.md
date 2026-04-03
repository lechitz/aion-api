# Fx Application Wiring

**Path:** `internal/platform/fxapp`

## Purpose

This package is the Uber Fx composition root for `aion-api`.
It wires infrastructure providers, application services, HTTP runtime, realtime consumption, and the dedicated outbox publisher process.

## Current Modules

| Module | Role |
| --- | --- |
| `InfraModule` | logger, config, cache, DB, HTTP client, and observability init |
| `ApplicationModule` | compose repositories, usecases, and `app.Dependencies` |
| `ServerModule` | compose HTTP handler, build server, and manage lifecycle |
| `RealtimeModule` | start the Kafka projection consumer when realtime is enabled |
| `OutboxPublisherModule` | start the periodic Kafka outbox publisher loop |

## Runtime Use

- `cmd/api` boots `InfraModule`, `ApplicationModule`, `RealtimeModule`, and `ServerModule`
- `cmd/outbox-publisher` boots `InfraModule` and `OutboxPublisherModule`

## Boundary Rules

- Fx wiring belongs here, not in the bounded contexts
- providers should expose stable contracts and delegate behavior to owning packages
- if startup behavior changes, update the matching command README and runtime docs in the same PR

## Validate

```bash
go test ./internal/platform/fxapp/...
make verify
```

## Risks And Compatibility Notes

- degraded startup behavior and module ordering are runtime contracts even when they are not directly user-facing
- this package is where cross-cutting startup failures usually become visible first; keep validation commands truthful

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
