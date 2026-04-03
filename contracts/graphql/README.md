# Shared GraphQL Contract Surface

**Path:** `contracts/graphql`

## Purpose

This folder owns the shared GraphQL operation set published for consumers, validation, and cross-repo contract review.
It is the public operation layer that complements, but does not replace, the live schema modules and resolver behavior inside the backend.

## Current Files

| Path | Responsibility |
| --- | --- |
| `manifest.json` | deterministic manifest of published operation documents and checksums |
| `queries/` | shared read operations grouped by capability |
| `mutations/` | shared write operations grouped by capability |

## Authority Order

For GraphQL behavior, use this order:

1. live schema modules and resolver/controller behavior under `internal/adapter/primary/graphql`
2. generated schema artifacts under `docs/graphql`
3. shared operation documents and `manifest.json` in this folder
4. narrative docs that summarize the public surface

## Boundary Rules

- keep only backend-owned shared operation documents here; ad-hoc consumer experiments do not belong in this folder
- operation names, top-level variables, and selection-set meaning are compatibility-sensitive and should change additively whenever possible
- shared governance belongs here at the folder root; inventories inside `queries/` and `mutations/` should stay focused on their own surface

## Validate

```bash
make graphql.queries
make graphql.manifest
make graphql.validate
make graphql.check-dirty
```

## Risks And Compatibility Notes

- `insightFeed` and `analyticsSeries` are backend-owned public contracts used across multiple repos and tools
- record projection queries are the preferred shared read surface for derived record consumers while compatibility migrations complete
- drift between live schema behavior and the published shared documents must be corrected in the same PR; if they conflict, the live backend remains the source of truth

## Related Docs

- [`queries/README.md`](./queries/README.md)
- [`mutations/README.md`](./mutations/README.md)
- [`../../docs/graphql/README.md`](../../docs/graphql/README.md)
- [`../../.github/DOCUMENTATION_OWNERSHIP.md`](../../.github/DOCUMENTATION_OWNERSHIP.md)

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
