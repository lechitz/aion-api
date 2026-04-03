# HTTP Service Token Middleware

**Path:** `internal/platform/server/http/middleware/servicetoken`

## Purpose

This package provides service-to-service authentication through HTTP headers.
It validates a trusted service key and can optionally inject a user identity into request context.

## Package Scope

| Area | Responsibility |
| --- | --- |
| S2S key validation | validate `X-Service-Key` against configured service key |
| Optional user impersonation | parse `X-Service-User-Id` and inject into context when valid |
| Context enrichment | set `ctxkeys.ServiceAccount` and optional `ctxkeys.UserID` |
| Fast rejection | return `401 Unauthorized` on invalid or misconfigured service key |

## Public API Reference

| Function / Constant | Description |
| --- | --- |
| `New(cfg, log)` | returns middleware that validates S2S calls when the service key header is present |
| `HeaderServiceKey` | header name for service credential: `X-Service-Key` |
| `HeaderServiceUser` | optional user ID header: `X-Service-User-Id` |
| `ErrServiceTokenInvalid` | error message for unauthorized S2S attempts |

## Runtime Behavior

1. Read `X-Service-Key` from the request.
2. If the header is absent, pass the request through unchanged.
3. If the header exists but the configured service key is empty, respond `401`.
4. If the header exists and does not match the configured key, respond `401`.
5. On success, set `ctxkeys.ServiceAccount = true`.
6. If `X-Service-User-Id` is present and parseable as `uint64`, set `ctxkeys.UserID`.
7. Continue with the enriched context.

## Boundary Rules

- this is a trust boundary, not an authorization boundary
- S2S authentication is transport-level behavior; domain authorization remains outside this package
- shared constants are used for log keys and context keys to avoid string drift

## Validate

```bash
go test ./internal/platform/server/http/middleware/servicetoken/...
go test ./internal/platform/server/http/...
```

## Risks And Compatibility Notes

- the middleware is permissive when no service key header is present and becomes authoritative only when S2S headers are supplied
- the package already has focused unit coverage for missing, invalid, and valid header combinations; keep those tests aligned with any header or context-contract change
- response bodies are currently plain `http.Error` strings, so any move to JSON output must be coordinated with callers that inspect transport behavior

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../../README.md)
<!-- doc-nav:end -->
