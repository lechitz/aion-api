# Shared GraphQL Mutations

**Path:** `contracts/graphql/mutations`

## Purpose

This folder contains the shared GraphQL mutation documents published for consumers and contract validation.

Unlike ad-hoc consumer mutations, files here are backend-owned artifacts and should stay aligned with:

- the live schema modules in `internal/adapter/primary/graphql/schema/modules/`
- the generated schema artifact in `docs/graphql/schema.graphql`
- the shared manifest under `contracts/graphql/manifest.json`

## Current Inventory

| Area | Documents |
| --- | --- |
| `categories/` | `create`, `update`, `delete` |
| `tags/` | `create`, `update`, `delete` |
| `records/` | `create`, `update`, `delete`, `delete-all` |
| `dashboard/` | `create-view`, `update-view`, `set-default-view`, `delete-view`, `upsert-widget`, `reorder-widgets`, `create-metric-and-widget`, `upsert-metric-definition`, `upsert-goal-template`, `delete-widget`, `delete-goal-template` |

## Current Contract Notes

- dashboard mutations are owned by the backend even when the dashboard drives the UX
- there is no shared GraphQL chat mutation surface here; chat write flows remain HTTP-driven in the current runtime
- query and mutation inventories are sibling public surfaces; shared governance for both lives in [`../README.md`](../README.md)

## Validate

```bash
make graphql.queries
make graphql.manifest
make graphql.validate
```

## Risks And Compatibility Notes

- mutation names and top-level argument meaning are compatibility-sensitive for generated clients and persisted workflows
- dashboard mutations remain backend-owned even when UI layers own richer local layout presentation
- if a mutation changes behavior materially, update the matching schema, shared operation document, and consumer-facing narrative docs together

See [`../queries/README.md`](../queries/README.md) for the complementary query inventory.
