# HTTP Cookie Utilities

**Path:** `internal/platform/server/http/utils/cookies`

## Purpose

This package centralizes HTTP cookie mechanics used by auth-facing handlers.

## Current Behavior

| Concern | Current rule |
| --- | --- |
| Access cookie name | `commonkeys.AuthTokenCookieName` |
| Refresh cookie name | literal `"refresh_token"` |
| Write policy | `HttpOnly`, config-driven `Secure`, config-driven `SameSite`, config path and domain |
| Refresh lifetime | `cfg.MaxAge * 7` |
| Clear policy | empty value, `MaxAge=-1`, expired timestamp, `SameSite=Strict` |

## Public Helpers

| Function | Responsibility |
| --- | --- |
| `SetAuthCookie`, `ClearAuthCookie` | write and expire the access cookie |
| `SetRefreshCookie`, `ClearRefreshCookie` | write and expire the refresh cookie |
| `ExtractAuthToken`, `ExtractRefreshToken` | read cookies from the incoming request |

## Boundary Rules

- this package owns cookie attributes and extraction only
- token issuance, session validation, and auth policy stay in the auth boundary
- any change to cookie names or lifetimes must stay aligned with dashboard and auth runtime behavior

## Validate

```bash
go test ./internal/platform/server/http/utils/cookies/...
make verify
```

## Risks And Compatibility Notes

- cookie names, max-age rules, and clear behavior are compatibility-sensitive for every authenticated client
- the refresh cookie name still uses a literal string, so any rename must be coordinated carefully across producers and consumers

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../../README.md)
<!-- doc-nav:end -->
