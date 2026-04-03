# Tag Bounded Context

**Path:** `internal/tag`

## Purpose

`internal/tag` owns user-scoped tag lifecycle, category association, and lookup flows used by record and dashboard surfaces.

## Current Surface

Current transport exposure is GraphQL-driven through the tag controller surface:

- create tag
- update tag
- soft-delete tag
- get by id
- get by name
- get by category
- list all tags for the authenticated user

There is no dedicated REST tag surface in the current runtime.

## Runtime Contract

- tags remain owned by one user and one category
- core usecases enforce uniqueness, ownership, and category relation rules
- cache adapters support id, name, category, and list lookups
- DB persistence is authoritative and backs derived fields such as `usageCount` and `lastUsedAt`

## Boundary Rules

- transport controllers should only map GraphQL types and authenticated user context
- tag and category consistency belongs in the tag and category cores plus their repositories
- record queries may consume tag semantics, but tag lifecycle ownership stays here

## Validate

```bash
go test ./internal/tag/...
go test ./internal/adapter/primary/graphql/...
make verify
```

## Risks And Compatibility Notes

- category relation rules are compatibility-sensitive for record filters and dashboard grouping
- if tag ownership or soft-delete semantics change, update related category and record docs in the same PR

## Related Docs

- [`../category/README.md`](../category/README.md)
- [`../record/README.md`](../record/README.md)

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
