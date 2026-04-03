# Platform Server Layer

**Path:** `internal/platform/server`

## Purpose

This layer owns server runtime composition.
Today it is HTTP-only.

## Current Flow

| Area | Role |
| --- | --- |
| `http/` | compose handlers, router adapters, middleware, generic fallbacks, and the `http.Server` instance |

`fxapp.ServerModule` consumes this layer to compose the runtime handler and start the server lifecycle.

## Boundary Rules

- context packages register routes through the HTTP port abstractions, not directly against concrete routers
- protocol-specific behavior belongs in `http/`; only cross-protocol server composition should live here in the future
- startup and shutdown ownership stays with `fxapp`, not with the bounded contexts

## Validate

```bash
go test ./internal/platform/server/...
make verify
```

## Risks And Compatibility Notes

- protocol expansion should happen here only when a new server boundary is truly cross-context
- handler composition and lifecycle ownership must stay aligned with `fxapp` and the command entrypoints

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
