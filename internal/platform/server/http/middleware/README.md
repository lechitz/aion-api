# Platform HTTP Middleware Layer

**Path:** `internal/platform/server/http/middleware`

## Purpose

This package owns cross-cutting HTTP middleware applied before bounded-context handlers.

## Current Runtime Usage

| Middleware | Current scope | What it does |
| --- | --- | --- |
| `requestid.New()` | global router and health endpoints | normalize or generate `X-Request-ID`, store it in context, echo it in the response |
| `recovery.New(genericHandler)` | global router only | catch panics and delegate sanitized response handling to the generic handler |
| `cors.New()` | global router and health endpoints | apply the current browser-origin and credential policy |
| `servicetoken.New(cfg, log)` | GraphQL mount only | validate trusted S2S headers and optionally inject service-account user context |

## Effective Order In `composer.go`

1. `requestid.New()`
2. `recovery.New(genericHandler)`
3. `cors.New()`
4. `servicetoken.New(cfg, log)` only around the GraphQL mount

This order is the current truth and takes precedence over older recommendations elsewhere.

## Boundary Rules

- this README owns `cors`, `recovery`, and `requestid` behavior
- those leaf READMEs are intentionally removed to reduce drift
- `servicetoken` keeps its own README because it is a distinct trust boundary

## Validate

```bash
go test ./internal/platform/server/http/middleware/...
go test ./internal/platform/server/http/...
```

## Risks And Compatibility Notes

- health routes intentionally skip the global recovery and `otelhttp` wrappers, so middleware ordering changes can alter diagnostics behavior
- request-id and CORS behavior are transport contracts that affect every HTTP consumer

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../README.md)
<!-- doc-nav:end -->
