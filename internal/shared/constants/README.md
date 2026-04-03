# Shared Constants

**Path:** `internal/shared/constants`

## Purpose

This package centralizes the stable cross-cutting keys reused by multiple parts of the runtime.

## Current Namespaces

| Namespace | Role |
| --- | --- |
| `claimskeys/` | JWT claim names |
| `commonkeys/` | log fields, request keys, cookie names, and other shared labels |
| `ctxkeys/` | typed context keys |
| `roles/` | shared role names |
| `tracingkeys/` | legacy HTTP and request tracing attributes kept for compatibility |

## Boundary Rules

- keep only constants and minimal helper types here
- tracing and status strings specific to one bounded context should stay local to that context
- adding a constant here means the value is intentionally shared across boundaries

## Validate

```bash
make verify
```

## Risks And Compatibility Notes

- renaming shared keys is a compatibility change because it can affect auth, logs, traces, dashboards, and clients at once
- if a value is not truly cross-cutting, keep it local instead of growing this package

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
