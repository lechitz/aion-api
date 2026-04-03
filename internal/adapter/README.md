# Shared Adapter Layer

**Path:** `internal/adapter`

## Purpose

`internal/adapter` holds adapter code that is shared across more than one bounded context.

## Current Split

| Subpackage | Role |
| --- | --- |
| `primary/` | shared inbound transport infrastructure |
| `secondary/` | shared outbound infrastructure implementations |

## Boundary Rules

- context-specific HTTP, DB, cache, storage, or provider adapters should stay inside the owning context
- only cross-context adapter infrastructure belongs here
- adapters translate transport or vendor behavior; they do not own business orchestration

## Validate

```bash
go test ./internal/adapter/...
make verify
```

## Risks And Compatibility Notes

- overusing this layer weakens bounded-context ownership
- if an adapter is shared only because of temporary migration pressure, document the exit path instead of treating it as permanent architecture

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
