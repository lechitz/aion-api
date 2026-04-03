# Record Bounded Context

**Path:** `internal/record`

## Purpose

`internal/record` owns user-scoped record lifecycle, derived read models, dashboard semantics, and the v1 insight and analytics layer built on top of record truth.

## Current Responsibilities

| Area | Responsibility |
| --- | --- |
| Record lifecycle | create, read, update, and soft-delete records |
| Query surfaces | date, tag, category, user, and search-driven retrieval |
| Derived models | record projection and graph projection shaping |
| Dashboard semantics | metric definitions, goal templates, widget catalog rules, and dashboard snapshot assembly |
| Intelligence | deterministic `insightFeed` and narrow `analyticsSeries` aggregation |

## Current Shape

| Area | Responsibility |
| --- | --- |
| `core/domain` | canonical record, projection, dashboard, and filter models |
| `core/usecase` | lifecycle orchestration, projection reads, dashboard semantics, insights, and analytics |
| `adapter/primary/graphql` | GraphQL transport mapping for record-facing contracts |
| `adapter/secondary/db` | authoritative persistence and derived read-model queries |
| `adapter/secondary/cache` | cache boundary for hot record reads |

## Boundary Rules

- record lifecycle ownership stays here even when category, tag, dashboard, or chat flows consume derived record semantics
- dashboard and insight meaning for v1 is backend-owned here, not invented downstream by UI state
- graph projections and dashboard projections are derived surfaces; they must not redefine the authoritative record model
- transport adapters must keep mapping, pagination, and filter decoding thin; lifecycle and scope semantics belong in core usecases

## Validation

```bash
go test ./internal/record/...
go test ./internal/adapter/primary/graphql/...
make graphql.validate
make verify
```

## Performance Readiness

This context is one of the main candidates for future real benchmarks, but today the reliable checks are boundary-level:

```bash
make record-projection-smoke
make record-projection-page-smoke
make realtime-record-smoke
make load-test-record-projections
make load-test-dashboard-snapshot
```

Observe:

- projection materialization health
- pagination stability on derived reads
- dashboard and analytics query latency in Grafana
- whether outbox or realtime issues are inflating perceived record latency
- whether `recordProjectionsLatest` stays inside the committed local load thresholds
- whether `dashboardSnapshot` stays inside the committed local load thresholds

Do not document per-function numbers here until the repo has committed benchmarks for record-heavy hot paths.

## Risks And Compatibility Notes

- widget layout remains a two-tier v1 contract:
  - `aion-api` owns widget types, coarse sizes, persisted ordering, and large-card limits
  - `aion-web` owns the richer visual grammar stored in `configJson`
- `dashboardSnapshot` checklist payloads and `analyticsSeries` semantics are consumer-visible contracts and should evolve additively
- scope semantics must remain aligned across read surfaces:
  - `window`
  - optional `date`
  - optional `timezone`
  - optional `categoryId`
  - optional `tagIds`
- graph projection exports are useful diagnostics, but they are not the authority over record truth

## Related Docs

- [`../category/README.md`](../category/README.md)
- [`../tag/README.md`](../tag/README.md)
- [`../../contracts/graphql/README.md`](../../contracts/graphql/README.md)

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
