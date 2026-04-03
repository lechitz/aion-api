# HTTP Response Utilities

**Path:** `internal/platform/server/http/utils/httpresponse`

## Purpose

This package standardizes HTTP JSON responses across handlers, including success envelopes, error envelopes, and span-aware error helpers.
It is the response boundary for HTTP adapters and relies on semantic errors from `sharederrors`.

## Package Scope

| Area | Responsibility |
| --- | --- |
| JSON writing | encode payloads and set HTTP headers consistently |
| Success responses | build normalized success envelopes |
| Error responses | map semantic errors to status and write normalized error envelopes |
| Span-aware helpers | record OTel span status and attributes before returning HTTP errors |

## Public API Reference

### Core Writers

| Function | Behavior |
| --- | --- |
| `WriteJSON(w, status, payload, headers...)` | writes raw JSON payload and skips body for `204 No Content` |
| `WriteSuccess(w, status, result, message, headers...)` | writes a standardized success envelope |
| `WriteError(w, err, message, log, headers...)` | maps error status, logs, and writes a standardized error envelope |
| `WriteDecodeError(w, err, log, headers...)` | shortcut for malformed body responses |
| `WriteAuthError(w, err, log, headers...)` | shortcut for auth failures |
| `WriteNoContent(w, headers...)` | writes `204` with optional headers |

### Span-aware Writers

| Function | Trace behavior |
| --- | --- |
| `WriteAuthErrorSpan(...)` | records the error, sets span status, writes auth error |
| `WriteDecodeErrorSpan(...)` | records decode error and returns `400` |
| `WriteValidationErrorSpan(...)` | records validation error and returns a validation-facing response |
| `WriteDomainErrorSpan(...)` | records domain error and maps status via `sharederrors.MapErrorToHTTPStatus` |

## Status Mapping Behavior

`WriteError` and domain span helpers delegate status mapping to `sharederrors.MapErrorToHTTPStatus(err)`.
This keeps status semantics centralized and consistent across all HTTP adapters.

## Tested Behaviors

| Behavior | Verified by |
| --- | --- |
| `204` responses do not write body or content-type | `TestWriteJSON`, `TestWriteNoContent` |
| Success envelopes keep code, message, result, and date | `TestWriteSuccess` |
| Error envelopes keep mapped status without leaking raw internals | `TestWriteError`, `TestWriteAuthAndDecodeError` |
| Custom response headers are preserved | `TestWriteError_WithCustomHeaders` |
| Span helpers set `codes.Error` and HTTP status attributes | `TestSpanErrorResponses` |

## Boundary Rules

- keep handlers thin by delegating response shape and status mapping to this package
- keep transport concerns here; domain usecases should not depend on this package
- use `Write*Span` helpers at adapter boundaries where tracing is active

## Validate

```bash
go test ./internal/platform/server/http/utils/httpresponse/...
go test ./internal/platform/server/http/utils/sharederrors/...
```

## Risks And Compatibility Notes

- response envelopes are shared HTTP contracts, so changes here ripple across many handlers at once
- `details` stays omitted by default to avoid leaking internal error data; relaxing that rule should be treated as a compatibility and security decision

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../../README.md)
<!-- doc-nav:end -->
