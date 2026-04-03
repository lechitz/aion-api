# Platform HTTP Generic Layer

**Path:** `internal/platform/server/http/generic`

## Purpose

This package contains platform-level generic HTTP components shared by all contexts.
It provides health and fallback handlers with consistent response envelopes, logging, and tracing.

## Package Scope

| Area | Responsibility |
| --- | --- |
| Health endpoint | return service status and runtime metadata |
| Router fallback handlers | standardize `404`, `405`, and router-level error behavior |
| Panic recovery integration | convert recovered panics into safe HTTP error responses |
| Shared transport DTOs | define transport-only payloads used by generic handlers |

## Subpackages

| Subpackage | Role |
| --- | --- |
| `dto/` | generic HTTP DTOs, currently the health response payload |
| `handler/` | generic handler implementation and tracing/logging constants |

## Main Components

| Component | Description |
| --- | --- |
| `handler.New(logger, generalCfg)` | creates `*handler.Handler` with logger and general app metadata |
| `(*Handler).HealthCheck` | handles `/health` and returns app metadata plus healthy status |
| `(*Handler).NotFoundHandler` | standardized JSON `404` |
| `(*Handler).MethodNotAllowedHandler` | standardized JSON `405` |
| `(*Handler).ErrorHandler` | standardized JSON `500` for router-level errors |
| `(*Handler).RecoveryHandler` | handles panic recovery payloads and emits telemetry |

## Integration Flow

1. HTTP composer creates `genericHandler := handler.New(log, cfg.General)`.
2. Router wiring connects fallback handlers and health routes to that handler.
3. Recovery middleware delegates panic handling to `genericHandler.RecoveryHandler`.

## Boundary Rules

- this layer is transport/platform only and must not contain business rules
- generic handlers provide the baseline fallback behavior for all bounded contexts
- response shapes and status semantics must stay aligned with `httpresponse` and `sharederrors`

## Validate

```bash
go test ./internal/platform/server/http/generic/...
go test ./internal/platform/server/http/...
```

## Risks And Compatibility Notes

- health, `404`, `405`, router-level error, and panic-recovery responses are shared transport contracts
- this package already has handler coverage in `handler/generic_handler_test.go`; if response shapes change, update those tests in the same PR
- fallback behavior should stay centralized here so individual context handlers do not drift in transport semantics

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../README.md)
<!-- doc-nav:end -->
