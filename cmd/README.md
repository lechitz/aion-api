# Application Entrypoints (`cmd`)

## Purpose

`cmd` owns process entrypoints and bootstrap-only concerns.

## Current Flow

| Entrypoint | Role |
| --- | --- |
| `cmd/api` | main API server process |
| `cmd/outbox-publisher` | dedicated background publisher for pending outbox rows |

## Boundary Rules

- business logic does not belong in `cmd`
- runtime composition belongs in `internal/platform/fxapp`
- bootstrap-only parsing and signal handling may live here when they are process-specific

## Validate

```bash
go run ./cmd/api
go run ./cmd/outbox-publisher
```

## Risks And Compatibility Notes

- entrypoint drift is easy to hide because the app may still compile while booting the wrong module set

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../README.md)
<!-- doc-nav:end -->
