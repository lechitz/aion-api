# Platform HTTP Server Layer

**Path:** `internal/platform/server/http`

## Purpose

This package composes the complete HTTP transport surface for `aion-api`.
It owns router creation, middleware wiring, Swagger/docs mounts, REST registrar composition, GraphQL mount, and health endpoints.

## Current Composition

| Concern | Current behavior |
| --- | --- |
| Main router | `chi` adapter created in `ComposeHandler` |
| Global middleware order | `requestid` -> `recovery` -> `cors` |
| Fallback handlers | `NotFound`, `MethodNotAllowed`, and generic error callback are wired from `generic/handler` |
| Swagger/docs | mounted under `cfg.ServerHTTP.Context` with alias redirect to `.../swagger/index.html` |
| REST routes | mounted under `cfg.ServerHTTP.Context + cfg.ServerHTTP.APIRoot` |
| GraphQL | mounted inside the API root at `cfg.ServerGraphql.Path`, wrapped by `servicetoken` |
| OTel HTTP wrapper | wraps the main router, not the dedicated health mux |
| Health routes | exposed separately at both `${context}${health}` and `${context}${apiRoot}${health}` |

## Conditional Route Registration

`registerDomainRoutes` mounts REST modules only when their dependencies are present:

- auth
- user
- admin
- chat
- audit
- realtime

GraphQL handler construction is also dependency-driven and mounted only after successful composition.

## Health Exception

Health endpoints intentionally bypass the instrumented main router.
They currently run through:

- `requestid`
- `cors`

They do not go through `otelhttp` or the main router fallback chain.

## Key Files

| File | Purpose |
| --- | --- |
| `composer.go` | route/middleware composition and mount logic |
| `server.go` | `http.Server` construction from config |
| `http_constants.go` | default route and mount constants |

## Boundaries

- No domain usecase logic belongs here.
- Bounded contexts register handlers through adapters; this package only composes them.
- Shared transport primitives live in `middleware/`, `generic/`, `router/`, `ports/`, and `utils/`.

## Validate

```bash
go test ./internal/platform/server/http/...
make verify
```

## Performance Readiness

Current transport-level performance validation is dashboard-driven plus versioned local load scenarios.

Use this path:

```bash
make dev
./infrastructure/observability/scripts/setup-improvements.sh
make load-test-auth-login
make load-test-record-projections
make load-test-dashboard-snapshot
make load-test-realtime-record-created
```

Then confirm in Grafana:

- throughput is populated
- error rate stays flat
- p95 latency remains within the expected local baseline for the exercised routes
- top slow endpoints remain explainable from traces

The committed load scenarios exercise real transport paths:

- `load-test-auth-login` covers the public auth HTTP boundary
- `load-test-record-projections` covers the authenticated GraphQL read path used by derived record consumers
- `load-test-dashboard-snapshot` covers the authenticated GraphQL dashboard aggregation path
- `load-test-realtime-record-created` covers the authenticated SSE path until a projection-ready event is delivered

For the full protocol, see [`../../../../docs/performance-readiness.md`](../../../../docs/performance-readiness.md).

## Risks And Compatibility Notes

- middleware order, route composition, and GraphQL mount behavior can shift latency across the whole transport surface even when no bounded-context code changes
- health routes intentionally bypass part of the main instrumented path, so they are not representative latency probes for the full API

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../README.md)
<!-- doc-nav:end -->
