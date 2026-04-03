# Secondary Adapters Layer

**Path:** `internal/adapter/secondary`

## Purpose

This layer holds shared outbound implementations reused across the application.
It exists for vendor-facing concerns that are truly cross-context, not for convenience moves out of a bounded context.

## Current Packages

| Package | Role |
| --- | --- |
| `cache/redis/` | shared Redis client boundary reused by context-local cache adapters |
| `db/postgres/` | shared PostgreSQL connection and transaction bootstrap |
| `graph/gremlin/` | optional graph client integration |
| `contextlogger/` | structured logger implementation |
| `crypto/` | key generation helpers |
| `hasher/` | bcrypt password hashing |
| `token/` | JWT token provider and expiry helpers |

## Boundary Rules

- context-owned DB, cache, HTTP, Kafka, and storage adapters stay inside the owning bounded context
- packages here implement reusable platform or security concerns shared across contexts
- business semantics must remain in core/usecase layers and output ports
- a package belongs here only if at least two boundaries could reasonably reuse the same implementation contract

## Validate

```bash
go test ./internal/adapter/secondary/...
make verify
```

## Risks And Compatibility Notes

- changes here often ripple across multiple contexts, even when the exported surface looks small
- vendor-specific behavior should be normalized here before it can leak into context code
- if a reusable adapter becomes effectively single-context, prefer moving it back to the owning context instead of growing this layer without clear reuse

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
