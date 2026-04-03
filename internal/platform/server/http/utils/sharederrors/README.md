# HTTP Shared Errors

**Path:** `internal/platform/server/http/utils/sharederrors`

## Purpose

This package centralizes HTTP-facing semantic errors and status-code mapping for handlers and adapter boundaries.
It keeps transport error behavior consistent without leaking infrastructure details.

## Package Scope

| Area | Responsibility |
| --- | --- |
| Typed errors | define explicit error types such as `ValidationError` and `UnauthorizedError` |
| Sentinel errors | provide reusable `errors.Is` targets for common conflicts and validation failures |
| Status mapping | convert known errors to stable HTTP status codes |
| Message constants | keep shared error messages consistent and reusable |

## Public API Reference

### Constructors And Helpers

| Function | Returns | Notes |
| --- | --- | --- |
| `ErrMissingUserID()` | `error` | missing user id in context |
| `ErrUnauthorized(reason)` | `error` | unauthorized with optional reason |
| `ErrForbidden(reason)` | `error` | forbidden with optional reason |
| `NewValidationError(field, reason)` | `error` | validation error with field context |
| `NewAuthenticationError(reason)` | `error` | authentication failure, mapped as unauthorized |
| `AtLeastOneFieldRequired(fields...)` | `error` | validation helper for partial update commands |
| `MissingFields(fields...)` | `error` | validation helper for required field checks |
| `MapErrorToHTTPStatus(err)` | `int` | main entrypoint used by HTTP response helpers |

### Typed Errors

| Type | Typical status |
| --- | --- |
| `ValidationError` | `400 Bad Request` |
| `UnauthorizedError` | `401 Unauthorized` |
| `ForbiddenError` | `403 Forbidden` |
| `MissingUserIDError` | `401 Unauthorized` |
| `AuthenticationError` | `401 Unauthorized` |

### Sentinel Errors

| Sentinel | Typical status |
| --- | --- |
| `ErrParseUserID` | `400 Bad Request` |
| `ErrNoFieldsToUpdate` | `400 Bad Request` |
| `ErrUsernameExists` | `409 Conflict` |
| `ErrEmailExists` | `409 Conflict` |
| `ErrDomainConflict` | `409 Conflict` |

## HTTP Mapping Summary

| Condition | Status code |
| --- | --- |
| `nil` error | `200 OK` |
| Validation errors | `400` |
| Unauthorized/authentication errors | `401` |
| Forbidden errors | `403` |
| `httperrors.ErrResourceNotFound` | `404` |
| `httperrors.ErrMethodNotAllowed` | `405` |
| Conflict errors | `409` |
| Unknown errors | `500` |

## Usage Example

```go
if err != nil {
    status := sharederrors.MapErrorToHTTPStatus(err)
    httpresponse.WriteError(w, err, "request failed", logger)
    _ = status // useful for metrics or structured logs
    return
}
```

## Boundary Rules

- keep this package transport-semantic only; domain orchestration does not belong here
- prefer `errors.As` for typed errors and `errors.Is` for sentinels
- if a new semantic error changes public HTTP behavior, update response-layer tests in the same PR

## Validate

```bash
go test ./internal/platform/server/http/utils/sharederrors/...
go test ./internal/platform/server/http/utils/httpresponse/...
```

## Risks And Compatibility Notes

- HTTP status mapping is shared behavior for all HTTP adapters, so small changes can create repo-wide regressions
- prefer adding new typed or sentinel errors over broadening the meaning of existing ones
- if a semantic error becomes part of a consumer-visible API contract, document it at the nearest boundary that emits it

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../../README.md)
<!-- doc-nav:end -->
