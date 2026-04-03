# GraphQL Documentation Artifacts

**Path:** `docs/graphql`

## Purpose

This folder stores generated GraphQL artifacts used by consumers and tooling.
It complements, but does not replace, the live schema modules under `internal/adapter/primary/graphql/schema/modules/` and the shared operation set under `contracts/graphql/`.

## Files

| File | Purpose |
| --- | --- |
| `schema.graphql` | flattened SDL snapshot generated from the current schema modules |

## Related Sources

- schema modules: `internal/adapter/primary/graphql/schema/modules/`
- shared operations: `contracts/graphql/queries/` and `contracts/graphql/mutations/`
- contract manifest: `contracts/graphql/manifest.json`
- ownership map: `.github/DOCUMENTATION_OWNERSHIP.md`

## Public Contract Notes

These artifacts currently include the main read and dashboard surfaces consumed by the workspace, including:

- legacy record reads such as `recordsLatest`
- derived projection reads such as `recordProjectionById`, `recordProjections`, and `recordProjectionsLatest`
- dashboard reads such as `dashboardSnapshot`, `insightFeed`, `analyticsSeries`, `metricDefinitions`, `dashboardViews`, `dashboardView`, `dashboardWidgetCatalog`, and `suggestMetricDefinitions`
- dashboard mutations such as `createDashboardView`, `setDefaultDashboardView`, `upsertDashboardWidget`, `reorderDashboardWidgets`, `deleteDashboardWidget`, and `createMetricAndWidget`

`dashboardSnapshot` also carries a checklist-oriented sub-payload inside `metrics[].checklist` for count-backed checklist widgets.

For the widget system specifically, the published GraphQL contract is intentionally coarse for v1:

- `dashboardWidgetCatalog` exposes canonical widget types, coarse sizes, and large-card limits
- widget records expose persisted `configJson`, but `aion-api` currently treats the richer visual layout grammar inside that JSON as a dashboard-owned concern rather than a server-validated schema

## Validate

```bash
make graphql.schema
make graphql.queries
make graphql.manifest
make graphql.validate
```

## Boundary Rules

- treat files here as generated artifacts
- regenerate after schema changes in the same PR
- keep this folder aligned with consumer tooling expectations, but treat `contracts/graphql/` as the reusable public operation surface

## Risks And Compatibility Notes

- if this folder drifts from live schema modules or shared operations, the live modules win
- codegen consumers and static docs readers often depend on this snapshot without reading backend code, so stale generation is a real contract risk

## Related Contract Sources

- `contracts/graphql/README.md`
- `contracts/graphql/queries/README.md`
- `contracts/graphql/mutations/README.md`

---

<!-- doc-nav:start -->
## Navigation
- [Back to docs index](../index.md)
<!-- doc-nav:end -->
