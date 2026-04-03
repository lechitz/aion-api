# Primary Adapters Layer

**Path:** `internal/adapter/primary`

## Purpose

This layer owns shared inbound transport infrastructure.
Today that means the central GraphQL boundary consumed across multiple contexts.

## Current Surface

| Subpackage | Role |
| --- | --- |
| `graphql/` | gqlgen config, schema composition, shared resolvers, directives, and server bootstrap |

## Boundary Rules

- most REST handlers remain in the owning bounded context under `internal/<ctx>/adapter/primary/http`
- shared primary adapter code belongs here only when it coordinates multiple contexts
- resolvers and transport glue stay orchestration-only; usecase behavior remains in context services

## Validate

```bash
go test ./internal/adapter/primary/...
make graphql.validate
```

## Risks And Compatibility Notes

- shared transport changes here can affect every context exposed through GraphQL
- if a new shared primary transport appears, it should have a clear multi-context ownership reason before being added under this layer

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
