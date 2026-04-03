# Platform HTTP Utilities Layer

**Path:** `internal/platform/server/http/utils`

## Purpose

This package aggregates shared HTTP utility modules used by platform and context adapters.
It provides a consistent foundation for response writing, semantic error handling, and auth cookie operations.

## Subpackages

| Subpackage | Responsibility |
| --- | --- |
| `sharederrors/` | HTTP-facing semantic errors and error-to-status mapping |
| `httpresponse/` | standard response envelope and response-writer helpers |
| `cookies/` | auth and refresh cookie lifecycle and extraction helpers |

## Typical Request Or Response Path

1. Handler or usecase returns success payload or semantic error.
2. `sharederrors` defines the semantic category.
3. `httpresponse` maps status and writes the standardized JSON response.
4. `cookies` helpers set, clear, or extract auth cookies when needed.

## Boundary Rules

- keep all utilities in this layer transport-focused and reusable across contexts
- domain or business rules must not be added to these utility packages
- subpackage READMEs contain the detailed contracts; this README is the high-level integration view

## Validate

```bash
go test ./internal/platform/server/http/utils/...
go test ./internal/platform/server/http/...
```

## Risks And Compatibility Notes

- response shape, cookie behavior, and semantic error mapping are shared transport contracts
- fragmenting those rules across handlers is a fast path to drift, so prefer extending these utilities rather than duplicating behavior

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../README.md)
<!-- doc-nav:end -->
