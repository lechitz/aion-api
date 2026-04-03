# Shared GraphQL Mutations

## Purpose

This folder contains the shared GraphQL mutation documents published for consumers and contract validation.

Unlike ad-hoc consumer queries, files here are backend-owned artifacts and should stay aligned with:

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

- Dashboard mutations are owned by the backend even when the dashboard drives the UX.
- There is no shared GraphQL chat mutation surface here; chat write flows remain HTTP-driven in the current runtime.
- Query and mutation inventories are sibling public surfaces; consumers should not assume one README supersedes the other.

## Validate

```bash
make graphql.queries graphql.manifest graphql.validate
```

See [`../queries/README.md`](../queries/README.md) for the complementary query inventory and the broader governance rules.
