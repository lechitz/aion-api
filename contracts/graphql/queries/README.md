# Shared GraphQL Queries

**Path:** `contracts/graphql/queries`

## Purpose

This folder stores reusable shared read operations aligned with the backend schema.
It is the query half of the published GraphQL contract surface consumed by Aion clients and tools.

## Current Areas

| Folder | Scope |
| --- | --- |
| `categories/` | category read operations |
| `tags/` | tag read operations |
| `records/` | record read operations |
| `chat/` | chat read operations |
| `user/` | user read operations |
| `dashboard/` | dashboard, insight, and analytics reads |

## Query Inventory

### Categories
- `queries/categories/list.graphql`
- `queries/categories/by-id.graphql`
- `queries/categories/by-name.graphql`

### Tags
- `queries/tags/list.graphql`
- `queries/tags/by-id.graphql`
- `queries/tags/by-name.graphql`
- `queries/tags/by-category-id.graphql`

### Records
- `queries/records/list.graphql`
- `queries/records/by-id.graphql`
- `queries/records/latest.graphql`
- `queries/records/projection-by-id.graphql`
- `queries/records/projections.graphql`
- `queries/records/projections-latest.graphql`
- `queries/records/by-tag.graphql`
- `queries/records/by-category.graphql`
- `queries/records/by-day.graphql`
- `queries/records/until.graphql`
- `queries/records/between.graphql`
- `queries/records/search.graphql`
- `queries/records/stats.graphql`

### Chat
- `queries/chat/history.graphql`
- `queries/chat/context.graphql`
- `queries/chat/data-pack.graphql`

### User
- `queries/user/stats.graphql`

### Dashboard
- `queries/dashboard/snapshot.graphql`
- `queries/dashboard/insight-feed.graphql`
- `queries/dashboard/analytics-series.graphql`
- `queries/dashboard/metric-definitions.graphql`
- `queries/dashboard/views.graphql`
- `queries/dashboard/view.graphql`
- `queries/dashboard/widget-catalog.graphql`
- `queries/dashboard/suggest-metric-definitions.graphql`

## Contract Notes

Canonical v1 intelligence queries:

- `InsightFeed`
- `AnalyticsSeries`

Shared scope model:

- `window`
- optional `date`
- optional `timezone`
- optional `categoryId`
- optional `tagIds`

Dashboard widget contract note:

- `queries/dashboard/widget-catalog.graphql` is the canonical coarse catalog for v1, not a full UI-layout schema
- the backend owns widget types, coarse sizes, and `maxLargeWidgets`
- the richer layout grammar stored in `configJson` remains dashboard-owned for now
- `INSIGHT_FEED` remains the only widget type that is valid without `metricDefinitionId`

## Boundary Rules

- keep query documents stable for observability, typed-client generation, and cache behavior
- selection-set changes should stay additive unless all downstream consumers are updated together
- shared governance for the whole GraphQL surface lives in [`../README.md`](../README.md); this README should stay query-focused

## Validate

```bash
make graphql.queries
make graphql.manifest
make graphql.validate
```

## Risks And Compatibility Notes

- `queries/dashboard/insight-feed.graphql` and `queries/dashboard/analytics-series.graphql` are backend-owned public contracts for v1 intelligence surfaces
- record projection queries are the preferred shared read surface for derived dashboard and chat consumers while compatibility migration completes
- if a shared query lags behind live behavior, correct the document in the same PR instead of compensating in consumer-specific docs

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
