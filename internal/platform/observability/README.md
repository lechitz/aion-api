# Platform Observability

**Path:** `internal/platform/observability`

## Purpose

This package owns observability bootstrap helpers for `aion-api`.
It initializes trace and metric exporters, normalizes OTLP settings, and defines the shared resource metadata attached to telemetry.

## Structure

| Path | Responsibility |
| --- | --- |
| `tracer/` | OTLP HTTP trace exporter bootstrap and global tracer provider |
| `metric/` | OTLP HTTP metric exporter bootstrap and global meter provider |
| `helpers.go` | shared header parsing and endpoint normalization helpers |

## Current Runtime Behavior

- tracing and metrics both bootstrap from `cfg.Observability`
- OTLP endpoints accept `host:port` or full `http(s)://...` values and are normalized before exporter creation
- resource attributes include service, environment, host, and instance identity
- tracing installs W3C `TraceContext` and `Baggage` propagators globally
- if exporter initialization fails, the app degrades gracefully and returns a no-op cleanup function instead of aborting startup

## Boundary Rules

- instrumentation points live in handlers, controllers, usecases, repositories, and runtime boundaries outside this package
- collector and container wiring belongs to `infrastructure/observability`
- sampling and exporter behavior remain config-driven; do not hardcode environment-specific endpoints elsewhere

## Validate

```bash
go test ./internal/platform/observability/...
make verify
```

## Risks And Compatibility Notes

- degraded startup is intentional, so observability regressions can hide behind a still-running service if validation is skipped
- label and resource-attribute changes can break dashboards and telemetry queries even when code compiles

## Related Docs

- [`../../../infrastructure/observability/README.md`](../../../infrastructure/observability/README.md)
- [`../config/README.md`](../config/README.md)

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
